use chrono::{DateTime, Duration, Utc};
use sqlx::{Error as SqlxError, PgPool};
use uuid::Uuid;

#[derive(Debug, sqlx::FromRow)]
pub struct RefreshToken {
    pub id: Uuid,
    pub user_id: Uuid,
    pub token_hash: String,
    pub expires_at: DateTime<Utc>,
    pub created_at: DateTime<Utc>,
}

#[derive(Clone)]
pub struct RefreshTokenRepository {
    pool: PgPool,
}

impl RefreshTokenRepository {
    pub fn new(pool: PgPool) -> Self {
        Self { pool }
    }

    /// Store a hashed refresh token
    pub async fn create(
        &self,
        user_id: Uuid,
        raw_token: &str, // now we pass the raw token
        expires_in: Duration,
    ) -> Result<Uuid, SqlxError> {
        let token_hash = crate::utils::hash::hash_refresh_token(raw_token);
        let expires_at = Utc::now() + expires_in;
        let row: (Uuid,) = sqlx::query_as(
            "INSERT INTO refresh_tokens (user_id, token_hash, expires_at) VALUES ($1, $2, $3) RETURNING id",
        )
        .bind(user_id)
        .bind(token_hash)
        .bind(expires_at)
        .fetch_one(&self.pool)
        .await?;
        Ok(row.0)
    }

    // In find_by_hash: accept raw token, hash it, then search
    pub async fn find_by_token(&self, raw_token: &str) -> Result<Option<RefreshToken>, SqlxError> {
        let token_hash = crate::utils::hash::hash_refresh_token(raw_token);
        sqlx::query_as::<_, RefreshToken>(
            "SELECT * FROM refresh_tokens WHERE token_hash = $1 AND expires_at > NOW()",
        )
        .bind(token_hash)
        .fetch_optional(&self.pool)
        .await
    }

    // Delete by raw token
    pub async fn delete_by_token(&self, raw_token: &str) -> Result<(), SqlxError> {
        let token_hash = crate::utils::hash::hash_refresh_token(raw_token);
        sqlx::query("DELETE FROM refresh_tokens WHERE token_hash = $1")
            .bind(token_hash)
            .execute(&self.pool)
            .await?;
        Ok(())
    }
}
