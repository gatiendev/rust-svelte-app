use crate::models::user::UserRepository;
use crate::utils::password;
use actix_web::{web, HttpResponse, HttpRequest, Responder, post, get};  // <-- added `get`
use actix_web::HttpMessage;                                           // <-- added for .extensions()
use actix_identity::Identity;
use serde::Deserialize;
use uuid::Uuid;
use actix_session::Session;

#[derive(Deserialize)]
pub struct LoginRequest {
    username: String,
    password: String,
}

#[derive(Deserialize)]
pub struct RegisterRequest {
    username: String,
    password: String,
}

// AppState now holds the repository
#[derive(Clone)]
pub struct AppState {
    pub user_repo: UserRepository,
}

#[post("/register")]
pub async fn register(
    req: web::Json<RegisterRequest>,
    state: web::Data<AppState>,
) -> impl Responder {
    // Check if user exists
    if let Ok(Some(_)) = state.user_repo.find_by_username(&req.username).await {
        return HttpResponse::Conflict().body("Username already taken");
    }

    let hash = match password::hash(&req.password) {
        Ok(h) => h,
        Err(_) => return HttpResponse::InternalServerError().body("Password hashing failed"),
    };

    match state.user_repo.create(&req.username, &hash).await {
        Ok(user_id) => HttpResponse::Created().json(serde_json::json!({ "user_id": user_id })),
        Err(_) => HttpResponse::InternalServerError().body("Failed to create user"),
    }
}

#[post("/login")]
pub async fn login(
    req: web::Json<LoginRequest>,
    state: web::Data<AppState>,
    session: Session,                      // <-- add Session extractor
    http_req: HttpRequest,
) -> impl Responder {
    eprintln!("[LOGIN] Starting login for user: {}", req.username);

    let user = match state.user_repo.find_by_username(&req.username).await {
        Ok(Some(u)) => {
            eprintln!("[LOGIN] User found in DB: {}", u.id);
            u
        }
        _ => {
            eprintln!("[LOGIN] User not found");
            return HttpResponse::Unauthorized().body("Invalid credentials");
        }
    };

    if !password::verify(&req.password, &user.password_hash).unwrap_or(false) {
        eprintln!("[LOGIN] Password verification failed");
        return HttpResponse::Unauthorized().body("Invalid credentials");
    }

    eprintln!("[LOGIN] Password verified, calling Identity::login");
    if let Err(e) = Identity::login(&http_req.extensions(), user.id.to_string()) {
        eprintln!("[LOGIN] Identity::login error: {}", e);
        return HttpResponse::InternalServerError().body("Failed to set identity");
    }
    eprintln!("[LOGIN] Identity::login succeeded");

    // Also manually insert into session to force a cookie
    if let Err(e) = session.insert("user_id", user.id.to_string()) {
        eprintln!("[LOGIN] Session insert error: {}", e);
        return HttpResponse::InternalServerError().body("Failed to set session");
    }
    eprintln!("[LOGIN] Session insert succeeded");

    HttpResponse::Ok().body("Logged in")
}


#[post("/logout")]
pub async fn logout(identity: Identity) -> impl Responder {
    identity.logout();
    HttpResponse::Ok().body("Logged out")
}

#[get("/profile")]
pub async fn profile(
    session: Session,
    state: web::Data<AppState>,
) -> impl Responder {
    // Try to get user_id from session (set in login handler)
    let user_id: Option<String> = session.get("user_id").unwrap_or(None);
    let user_id = match user_id {
        Some(id) => id,
        None => return HttpResponse::Unauthorized().body("Not authenticated"),
    };

    let user_uuid = match Uuid::parse_str(&user_id) {
        Ok(uuid) => uuid,
        Err(_) => return HttpResponse::BadRequest().body("Invalid user ID format"),
    };

    match state.user_repo.find_by_id(user_uuid).await {
        Ok(Some(user)) => HttpResponse::Ok().json(serde_json::json!({
            "id": user.id,
            "username": user.username,
            "created_at": user.created_at,
        })),
        Ok(None) => HttpResponse::NotFound().body("User not found"),
        Err(_) => HttpResponse::InternalServerError().body("Database error"),
    }
}