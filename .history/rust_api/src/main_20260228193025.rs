mod config;
mod db;
mod handlers;
mod models;
mod utils;

use axum::{
    routing::{get, post},
    Router,
};
use tower_cookies::CookieManagerLayer;
use tracing_subscriber;

use crate::{
    config::Config,
    db::create_pool,
    handlers::auth::{login, logout, profile, refresh, register, AppState},
    models::{refresh_token::RefreshTokenRepository, user::UserRepository},
};
use tower_http::cors::{Any, CorsLayer};
use tower_http::trace::{DefaultMakeSpan, DefaultOnResponse, TraceLayer};
use tracing::Level;

#[tokio::main]
async fn main() -> anyhow::Result<()> {
    // Initialize tracing (logging)
    tracing_subscriber::fmt::init();

    let trace_layer = TraceLayer::new_for_http()
        .make_span_with(DefaultMakeSpan::new().level(Level::INFO)) // span level
        // .on_request(DefaultOnRequest::new().level(Level::INFO)) // request started
        .on_response(DefaultOnResponse::new().level(Level::INFO));

    tracing::info!("Application starting - logging confirmed");

    // let _ = rustls::crypto::ring::default_provider().install_default();

    // Load config
    let config = Config::from_env()?;

    let cors = CorsLayer::new()
        .allow_origin("http://localhost:5173".parse().unwrap())
        .allow_methods(Any)
        .allow_headers(Any)
        .allow_credentials(true);

    // Database pool
    let pool = create_pool(&config).await?;

    // Run migrations
    sqlx::migrate!("./migrations").run(&pool).await?;

    // Repositories
    let user_repo = UserRepository::new(pool.clone());
    let refresh_token_repo = RefreshTokenRepository::new(pool);

    // App state
    let state = AppState {
        user_repo,
        refresh_token_repo,
        config: config.clone(),
    };

    // Build router
    let app = Router::new()
        .route("/register", post(register))
        .route("/login", post(login))
        .route("/logout", post(logout))
        .route("/refresh", post(refresh))
        .route("/profile", get(profile))
        .layer(cors)
        .layer(CookieManagerLayer::new()) // for cookie handling
        .layer(trace_layer)
        // .layer(TraceLayer::new_for_http()) // request logging
        .with_state(state);

    let addr = format!("{}:{}", config.host, config.port);
    tracing::info!("Server listening on {}", addr);

    let listener = tokio::net::TcpListener::bind(addr).await.unwrap();

    axum::serve(listener, app).await.unwrap();

    Ok(())
}
