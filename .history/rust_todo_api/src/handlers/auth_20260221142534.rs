use crate::models::user::UserRepository;
use crate::utils::password;
use actix_identity::Identity;
use actix_web::{post, web, HttpResponse, Responder};
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
    state: web::Data<AppState>,
    identity: Identity,
) -> impl Responder {
    let user = match state.user_repo.find_by_username(&req.username).await {
        Ok(Some(u)) => u,
        _ => return HttpResponse::Unauthorized().body("Invalid credentials"),
    };

    if !password::verify(&req.password, &user.password_hash).unwrap_or(false) {
        return HttpResponse::Unauthorized().body("Invalid credentials");
    }

    identity.remember(user.id.to_string());
    HttpResponse::Ok().body("Logged in")
}

#[post("/logout")]
pub async fn logout(identity: Identity) -> impl Responder {
    identity.logout();
    HttpResponse::Ok().body("Logged out")
}
