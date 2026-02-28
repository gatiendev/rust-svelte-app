use crate::models::user::UserRepository;
use crate::utils::password;
use actix_web::{web, HttpResponse, HttpRequest, Responder, post};   // <-- added `post`
use actix_web::HttpMessage;                                           // <-- added for .extensions()
use actix_identity::Identity;
use serde::Deserialize;

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
    state: web::Data<AppState>,                // renamed from `data` to `state`
    http_req: HttpRequest,
) -> impl Responder {
    let user = match state.user_repo.find_by_username(&req.username).await {
        Ok(Some(u)) => u,
        _ => return HttpResponse::Unauthorized().body("Invalid credentials"),
    };

    if !password::verify(&req.password, &user.password_hash).unwrap_or(false) {
        return HttpResponse::Unauthorized().body("Invalid credentials");
    }

    Identity::login(&http_req.extensions(), user.id.to_string())
        .expect("Failed to set identity");     // or handle error gracefully

    HttpResponse::Ok().body("Logged in")
}

#[post("/logout")]
pub async fn logout(identity: Identity) -> impl Responder {
    identity.logout();
    HttpResponse::Ok().body("Logged out")
}

#[get("/profile")]
pub async fn profile(
    identity: Identity,
    state: web::Data<AppState>,
) -> impl Responder {
    // Get the user ID from the identity
    let user_id = match identity.id() {
        Ok(id_str) => match Uuid::parse_str(&id_str) {
            Ok(uuid) => uuid,
            Err(_) => return HttpResponse::BadRequest().body("Invalid user ID format"),
        },
        Err(_) => return HttpResponse::Unauthorized().body("Not authenticated"),
    };

    // Fetch user from database
    match state.user_repo.find_by_id(user_id).await {
        Ok(Some(user)) => HttpResponse::Ok().json(serde_json::json!({
            "id": user.id,
            "username": user.username,
            "created_at": user.created_at,
        })),
        Ok(None) => HttpResponse::NotFound().body("User not found"),
        Err(_) => HttpResponse::InternalServerError().body("Database error"),
    }
}