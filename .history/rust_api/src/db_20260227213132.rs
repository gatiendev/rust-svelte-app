use crate::config::Config;
use sqlx::postgres::{PgPool, PgPoolOptions};

pub async fn create_pool(config: &Config) -> anyhow::Result<PgPool> {
    let pool = PgPoolOptions::new()
        .max_connections(5)
        .connect(&config.database_url)
        .await?;
    Ok(pool)
}
