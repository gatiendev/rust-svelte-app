mod config;
mod db;
mod handlers;
mod models;
mod utils;

use actix_identity::IdentityMiddleware;
use actix_session::{storage::CookieSessionStore, SessionMiddleware};
use actix_web::{cookie::Key, web, App, HttpServer};
use handlers::auth::AppState;
use sqlx::migrate::Migrator;
use std::path::Path;

env_logger::init_from_env(env_logger::Env::new().default_filter_or("info"));


#[actix_web::main]
async fn main() -> anyhow::Result<()> {
    // Load config
    let config = config::Config::from_env()?;

    // Create database pool
    let pool = db::create_pool(&config).await?;

    // Run migrations automatically
    let migrator = Migrator::new(Path::new("./migrations")).await?;
    migrator.run(&pool).await?;

    // Build app state
    let user_repo = models::user::UserRepository::new(pool.clone());
    let app_state = web::Data::new(AppState { user_repo });

    // Session key from env (must be 64 bytes for Key::from)
    let session_key = Key::from(config.session_secret.as_bytes());

    println!("Starting server at {}:{}", config.host, config.port);

    HttpServer::new(move || {
        App::new()
            .app_data(app_state.clone())
            .wrap(Logger::new("%r %s %Dms")) 
            .wrap(
                SessionMiddleware::builder(CookieSessionStore::default(), session_key.clone())
                    .cookie_name("auth-session".to_string())
                    .cookie_secure(false) // Set true in production
                    .cookie_http_only(true)
                    .build(),
            )
            .wrap(IdentityMiddleware::default())
            .service(handlers::auth::register)
            .service(handlers::auth::login)
            .service(handlers::auth::logout)
        // add profile handler if needed
    })
    .bind((config.host.as_str(), config.port))?
    .run()
    .await?;

    Ok(())
}
