use sqlx::{Error as SqlxError, PgPool};
use uuid::Uuid;

#[derive(Debug, sqlx::FromRow)]
pub struct User {
    pub id: Uuid,
    pub username: String,
    pub password_hash: String,
    pub created_at: chrono::DateTime<chrono::Utc>,
}

#[derive(Clone)]
pub struct UserRepository {
    pool: PgPool,
}

impl UserRepository {
    pub fn new(pool: PgPool) -> Self {
        Self { pool }
    }

    pub async fn find_by_username(&self, username: &str) -> Result<Option<User>, SqlxError> {
        let user = sqlx::query_as::<_, User>("SELECT * FROM users WHERE username = $1")
            .bind(username)
            .fetch_optional(&self.pool)
            .await?;
        Ok(user)
    }
    
    pub async fn find_by_id(&self, id: Uuid) -> Result<Option<User>, sqlx::Error> {
        let user = sqlx::query_as::<_, User>("SELECT * FROM users WHERE id = $1")
            .bind(id)
            .fetch_optional(&self.pool)
            .await?;
        Ok(user)
    }

    pub async fn create(&self, username: &str, password_hash: &str) -> Result<Uuid, SqlxError> {
        let row: (Uuid,) = sqlx::query_as(
            "INSERT INTO users (username, password_hash) VALUES ($1, $2) RETURNING id",
        )
        .bind(username)
        .bind(password_hash)
        .fetch_one(&self.pool)
        .await?;
        Ok(row.0)
    }

    // Add other methods as needed
}
