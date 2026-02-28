use axum::{extract::State, http::StatusCode, response::IntoResponse, Json};
use serde::{Deserialize, Serialize};
use tower_cookies::{cookie, Cookie, Cookies}; // added cookie for time::Duration
use uuid::Uuid;

use crate::{
    config::Config,
    models::{refresh_token::RefreshTokenRepository, user::UserRepository},
    utils::{jwt, password},
};

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
    cookies.add(access_cookie);

    let mut refresh_cookie = Cookie::new(REFRESH_TOKEN_COOKIE, refresh_token.to_string());
    refresh_cookie.set_http_only(true);
    refresh_cookie.set_secure(true);
    refresh_cookie.set_path("/");
    cookies.add(refresh_cookie);
}

// Helper to remove auth cookies (set max_age to 0)
fn remove_auth_cookies(cookies: &Cookies) {
    let mut access_cookie = Cookie::new(ACCESS_TOKEN_COOKIE, "");
    access_cookie.set_http_only(true);
    access_cookie.set_secure(true);
    access_cookie.set_path("/");
    access_cookie.set_max_age(Some(cookie::time::Duration::seconds(0)));
    cookies.add(access_cookie);

    let mut refresh_cookie = Cookie::new(REFRESH_TOKEN_COOKIE, "");
    refresh_cookie.set_http_only(true);
    refresh_cookie.set_secure(true);
    refresh_cookie.set_path("/");
    refresh_cookie.set_max_age(Some(cookie::time::Duration::seconds(0)));
    cookies.add(refresh_cookie);
}

// ---------- Handlers ----------

pub async fn register(
    State(state): State<AppState>,
    Json(req): Json<AuthRequest>,
) -> impl IntoResponse {
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

    let refresh_token = uuid::Uuid::new_v4().to_string();
    let refresh_token_hash =
        password::hash(&refresh_token).unwrap_or_else(|_| refresh_token.clone());

    let expires_in = chrono::Duration::seconds(state.config.refresh_token_expiration);
    if state
        .refresh_token_repo
        .create(user.id, &refresh_token, expires_in)
        .await
        .is_err()
    {
        return (
            StatusCode::INTERNAL_SERVER_ERROR,
            Json(AuthResponse {
                message: "Failed to store refresh token".to_string(),
            }),
        );
    }

    set_auth_cookies(&cookies, &access_token, &refresh_token);
    (
        StatusCode::OK,
        Json(AuthResponse {
            message: "Logged in".to_string(),
        }),
    )
}

pub async fn logout(cookies: Cookies, State(state): State<AppState>) -> impl IntoResponse {
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

    let refresh_token_hash =
        password::hash(&refresh_token).unwrap_or_else(|_| refresh_token.clone());

    if state
        .refresh_token_repo
        .delete_by_token(&refresh_token)
        .await
        .is_err()
    {
        eprintln!("Failed to delete refresh token from DB");
    }

    remove_auth_cookies(&cookies);
    (
        StatusCode::OK,
        Json(AuthResponse {
            message: "Logged out".to_string(),
        }),
    )
}

pub async fn refresh(cookies: Cookies, State(state): State<AppState>) -> impl IntoResponse {
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

    let refresh_token_hash =
        password::hash(&refresh_token).unwrap_or_else(|_| refresh_token.clone());

    let token_record = match state.refresh_token_repo.find_by_token(&refresh_token).await {
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

    // Delete the used token (rotation)
    if state
        .refresh_token_repo
        .delete_by_token(&refresh_token)
        .await
        .is_err()
    {
        eprintln!("Failed to delete used refresh token");
    }

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

    let new_refresh_token = uuid::Uuid::new_v4().to_string();
    let new_refresh_token_hash =
        password::hash(&new_refresh_token).unwrap_or_else(|_| new_refresh_token.clone());

    let expires_in = chrono::Duration::seconds(state.config.refresh_token_expiration);
    if state
        .refresh_token_repo
        .create(token_record.user_id, &new_refresh_token, expires_in)
        .await
        .is_err()
    {
        return (
            StatusCode::INTERNAL_SERVER_ERROR,
            Json(AuthResponse {
                message: "Failed to store new refresh token".to_string(),
            }),
        );
    }

    set_auth_cookies(&cookies, &new_access_token, &new_refresh_token);
    (
        StatusCode::OK,
        Json(AuthResponse {
            message: "Token refreshed".to_string(),
        }),
    )
}

pub async fn profile(cookies: Cookies, State(state): State<AppState>) -> impl IntoResponse {
    let access_token = match cookies.get(ACCESS_TOKEN_COOKIE) {
        Some(cookie) => cookie.value().to_string(),
        None => {
            return (
                StatusCode::UNAUTHORIZED,
                Json(serde_json::json!({ "error": "Not authenticated" })),
            )
        }
    };

    let user_id = match jwt::validate_access_token(&access_token, &state.config.jwt_secret) {
        Ok(id) => id,
        Err(_) => {
            return (
                StatusCode::UNAUTHORIZED,
                Json(serde_json::json!({ "error": "Invalid token" })),
            )
        }
    };

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
        Json(serde_json::json!({
            "id": user.id,
            "username": user.username,
            "created_at": user.created_at,
        })),
    )
}
