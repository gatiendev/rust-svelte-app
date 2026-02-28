use std::env;

#[derive(Clone)]
pub struct Config {
    pub database_url: String,
    pub jwt_secret: String,            // for signing tokens
    pub access_token_expiration: i64,  // seconds
    pub refresh_token_expiration: i64, // seconds
    pub host: String,
    pub port: u16,
}

impl Config {
    pub fn from_env() -> anyhow::Result<Self> {
        dotenv::dotenv().ok();

        Ok(Config {
            database_url: env::var("DATABASE_URL").expect("DATABASE_URL must be set"),
            jwt_secret: env::var("JWT_SECRET").expect("JWT_SECRET must be set"),
            access_token_expiration: env::var("ACCESS_TOKEN_EXPIRATION")
                .unwrap_or_else(|_| "900".to_string()) // 15 min
                .parse()?,
            refresh_token_expiration: env::var("REFRESH_TOKEN_EXPIRATION")
                .unwrap_or_else(|_| "604800".to_string()) // 7 days
                .parse()?,
            host: env::var("HOST").unwrap_or_else(|_| "127.0.0.1".to_string()),
            port: env::var("PORT")
                .unwrap_or_else(|_| "8080".to_string())
                .parse()?,
        })
    }
}
