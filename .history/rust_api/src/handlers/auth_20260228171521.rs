use axum::{
    extract::State,
    http::StatusCode,
    response::{IntoResponse, Json},
    routing::{get, post},
    Router,
};
use serde::{Deserialize, Serialize};
use tower_cookies::{Cookie, Cookies};
use uuid::Uuid;

use crate::{
    config::Config,
    models::{refresh_token::RefreshTokenRepository, user::UserRepository},
    utils::{jwt, password},
};
use chrono::Utc;

// Request/Response structs
#[derive(Deserialize)]
pub struct AuthRequest {
    username: String,
    password: String,
}

#[derive(Serialize)]
pub struct AuthResponse {
    message: String,
}

#[derive(Serialize)]
pub struct ProfileResponse {
    id: Uuid,
    username: String,
    created_at: chrono::DateTime<chrono::Utc>,
}

// App state
#[derive(Clone)]
pub struct AppState {
    pub user_repo: UserRepository,
    pub refresh_token_repo: RefreshTokenRepository,
    pub config: Config,
}

// Cookie names
const ACCESS_TOKEN_COOKIE: &str = "access_token";
const REFRESH_TOKEN_COOKIE: &str = "refresh_token";

// Helper to set auth cookies
fn set_auth_cookies(cookies: &Cookies, access_token: &str, refresh_token: &str) {
    let mut access_cookie = Cookie::new(ACCESS_TOKEN_COOKIE, access_token.to_string());
    access_cookie.set_http_only(true);
    access_cookie.set_secure(true); // set to false in dev without HTTPS
    access_cookie.set_path("/");
    // Optionally set SameSite
    cookies.add(access_cookie);

    let mut refresh_cookie = Cookie::new(REFRESH_TOKEN_COOKIE, refresh_token.to_string());
    refresh_cookie.set_http_only(true);
    refresh_cookie.set_secure(true);
    refresh_cookie.set_path("/");
    cookies.add(refresh_cookie);
}

// Helper to remove auth cookies
fn remove_auth_cookies(cookies: &Cookies) {
    let mut access_cookie = Cookie::new(ACCESS_TOKEN_COOKIE, "");
    access_cookie.set_http_only(true);
    access_cookie.set_secure(true);
    access_cookie.set_path("/");
    access_cookie.set_max_age(time::Duration::seconds(0)); // or .expires(OffsetDateTime::now_utc() - Duration::days(1))
    cookies.add(access_cookie);

    let mut refresh_cookie = Cookie::new(REFRESH_TOKEN_COOKIE, "");
    refresh_cookie.set_http_only(true);
    refresh_cookie.set_secure(true);
    refresh_cookie.set_path("/");
    refresh_cookie.set_max_age(time::Duration::seconds(0));
    cookies.add(refresh_cookie);
}

// ---------- Handlers ----------

pub async fn register(
    State(state): State<AppState>,
    Json(req): Json<AuthRequest>,
) -> impl IntoResponse {
    // Check if user exists
    if let Ok(Some(_)) = state.user_repo.find_by_username(&req.username).await {
        return (
            StatusCode::CONFLICT,
            Json(AuthResponse {
                message: "Username already taken".to_string(),
            }),
        );
    }

    let hash = match password::hash(&req.password) {
        Ok(h) => h,
        Err(_) => {
            return (
                StatusCode::INTERNAL_SERVER_ERROR,
                Json(AuthResponse {
                    message: "Password hashing failed".to_string(),
                }),
            )
        }
    };

    match state.user_repo.create(&req.username, &hash).await {
        Ok(user_id) => (
            StatusCode::CREATED,
            Json(AuthResponse {
                message: format!("User created with id: {}", user_id),
            }),
        ),
        Err(_) => (
            StatusCode::INTERNAL_SERVER_ERROR,
            Json(AuthResponse {
                message: "Failed to create user".to_string(),
            }),
        ),
    }
}

pub async fn login(
    cookies: Cookies,
    State(state): State<AppState>,
    Json(req): Json<AuthRequest>,
) -> impl IntoResponse {
    let user = match state.user_repo.find_by_username(&req.username).await {
        Ok(Some(u)) => u,
        _ => {
            return (
                StatusCode::UNAUTHORIZED,
                Json(AuthResponse {
                    message: "Invalid credentials".to_string(),
                }),
            )
        }
    };

    if !password::verify(&req.password, &user.password_hash).unwrap_or(false) {
        return (
            StatusCode::UNAUTHORIZED,
            Json(AuthResponse {
                message: "Invalid credentials".to_string(),
            }),
        );
    }

    // Generate access token
    let access_token = match jwt::create_access_token(
        user.id,
        &state.config.jwt_secret,
        state.config.access_token_expiration,
    ) {
        Ok(t) => t,
        Err(_) => {
            return (
                StatusCode::INTERNAL_SERVER_ERROR,
                Json(AuthResponse {
                    message: "Failed to generate access token".to_string(),
                }),
            )
        }
    };

    // Generate refresh token (random string)
    let refresh_token = uuid::Uuid::new_v4().to_string();
    // Hash it before storing (avoid storing raw tokens)
    let refresh_token_hash =
        password::hash(&refresh_token).unwrap_or_else(|_| refresh_token.clone()); // fallback; better to use a proper hash function

    // Store hashed refresh token in DB
    let expires_in = chrono::Duration::seconds(state.config.refresh_token_expiration);
    if let Err(_) = state
        .refresh_token_repo
        .create(user.id, &refresh_token_hash, expires_in)
        .await
    {
        return (
            StatusCode::INTERNAL_SERVER_ERROR,
            Json(AuthResponse {
                message: "Failed to store refresh token".to_string(),
            }),
        );
    }

    // Set cookies
    set_auth_cookies(&cookies, &access_token, &refresh_token);

    (
        StatusCode::OK,
        Json(AuthResponse {
            message: "Logged in".to_string(),
        }),
    )
}

pub async fn logout(cookies: Cookies, State(state): State<AppState>) -> impl IntoResponse {
    // Get refresh token from cookie
    let refresh_token = match cookies.get(REFRESH_TOKEN_COOKIE) {
        Some(cookie) => cookie.value().to_string(),
        None => {
            return (
                StatusCode::UNAUTHORIZED,
                Json(AuthResponse {
                    message: "No refresh token".to_string(),
                }),
            )
        }
    };

    // Hash it to find in DB (because we stored the hash)
    let refresh_token_hash =
        password::hash(&refresh_token).unwrap_or_else(|_| refresh_token.clone());

    // Delete from DB
    if let Err(_) = state.refresh_token_repo.delete(&refresh_token_hash).await {
        // Log error but still clear cookies
        eprintln!("Failed to delete refresh token from DB");
    }

    // Clear cookies
    remove_auth_cookies(&cookies);

    (
        StatusCode::OK,
        Json(AuthResponse {
            message: "Logged out".to_string(),
        }),
    )
}

pub async fn refresh(cookies: Cookies, State(state): State<AppState>) -> impl IntoResponse {
    // Get refresh token from cookie
    let refresh_token = match cookies.get(REFRESH_TOKEN_COOKIE) {
        Some(cookie) => cookie.value().to_string(),
        None => {
            return (
                StatusCode::UNAUTHORIZED,
                Json(AuthResponse {
                    message: "No refresh token".to_string(),
                }),
            )
        }
    };

    // Hash it to find in DB
    let refresh_token_hash =
        password::hash(&refresh_token).unwrap_or_else(|_| refresh_token.clone());

    // Validate refresh token in DB
    let token_record = match state
        .refresh_token_repo
        .find_by_hash(&refresh_token_hash)
        .await
    {
        Ok(Some(record)) => record,
        _ => {
            return (
                StatusCode::UNAUTHORIZED,
                Json(AuthResponse {
                    message: "Invalid or expired refresh token".to_string(),
                }),
            )
        }
    };

    // (Optional) Implement refresh token rotation: delete old token and issue new one
    // For simplicity, we just delete the used one
    if let Err(_) = state.refresh_token_repo.delete(&refresh_token_hash).await {
        eprintln!("Failed to delete used refresh token");
    }

    // Generate new access token
    let new_access_token = match jwt::create_access_token(
        token_record.user_id,
        &state.config.jwt_secret,
        state.config.access_token_expiration,
    ) {
        Ok(t) => t,
        Err(_) => {
            return (
                StatusCode::INTERNAL_SERVER_ERROR,
                Json(AuthResponse {
                    message: "Failed to generate new access token".to_string(),
                }),
            )
        }
    };

    // Generate new refresh token
    let new_refresh_token = uuid::Uuid::new_v4().to_string();
    let new_refresh_token_hash =
        password::hash(&new_refresh_token).unwrap_or_else(|_| new_refresh_token.clone());

    // Store new refresh token
    let expires_in = chrono::Duration::seconds(state.config.refresh_token_expiration);
    if let Err(_) = state
        .refresh_token_repo
        .create(token_record.user_id, &new_refresh_token_hash, expires_in)
        .await
    {
        return (
            StatusCode::INTERNAL_SERVER_ERROR,
            Json(AuthResponse {
                message: "Failed to store new refresh token".to_string(),
            }),
        );
    }

    // Set new cookies
    set_auth_cookies(&cookies, &new_access_token, &new_refresh_token);

    (
        StatusCode::OK,
        Json(AuthResponse {
            message: "Token refreshed".to_string(),
        }),
    )
}

pub async fn profile(cookies: Cookies, State(state): State<AppState>) -> impl IntoResponse {
    // Extract access token from cookie
    let access_token = match cookies.get(ACCESS_TOKEN_COOKIE) {
        Some(cookie) => cookie.value().to_string(),
        None => {
            return (
                StatusCode::UNAUTHORIZED,
                Json(serde_json::json!({ "error": "Not authenticated" })),
            )
        }
    };

    // Validate token
    let user_id = match jwt::validate_access_token(&access_token, &state.config.jwt_secret) {
        Ok(id) => id,
        Err(_) => {
            return (
                StatusCode::UNAUTHORIZED,
                Json(serde_json::json!({ "error": "Invalid token" })),
            )
        }
    };

    // Fetch user from DB
    let user = match state.user_repo.find_by_id(user_id).await {
        Ok(Some(u)) => u,
        _ => {
            return (
                StatusCode::NOT_FOUND,
                Json(serde_json::json!({ "error": "User not found" })),
            )
        }
    };

    (
        StatusCode::OK,
        Json(ProfileResponse {
            id: user.id,
            username: user.username,
            created_at: user.created_at,
        }),
    )
}
