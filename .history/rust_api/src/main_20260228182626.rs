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
use tower_http::trace::TraceLayer;
use tracing_subscriber::{layer::SubscriberExt, util::SubscriberInitExt};

use crate::{
    config::Config,
    db::create_pool,
    handlers::auth::{login, logout, profile, refresh, register, AppState},
    models::{refresh_token::RefreshTokenRepository, user::UserRepository},
};
use tower_http::trace::{DefaultMakeSpan, DefaultOnRequest, DefaultOnResponse, TraceLayer};
use tracing::Level;

use jsonwebtoken::crypto::CryptoProvider;
use rustls;

#[tokio::main]
async fn main() -> anyhow::Result<()> {
    // Initialize tracing (logging)
    tracing_subscriber::fmt::init();

    let trace_layer = TraceLayer::new_for_http()
        .make_span_with(DefaultMakeSpan::new().level(Level::INFO)) // span level
        .on_request(DefaultOnRequest::new().level(Level::INFO)) // request started
        .on_response(DefaultOnResponse::new().level(Level::INFO));

    tracing::info!("Application starting - logging confirmed");

    // let _ = rustls::crypto::ring::default_provider().install_default();

    // Load config
    let config = Config::from_env()?;

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
