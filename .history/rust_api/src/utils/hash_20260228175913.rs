use sha2::{Digest, Sha256};

/// Deterministically hash a refresh token (for storage & lookup)
pub fn hash_refresh_token(token: &str) -> String {
    let mut hasher = Sha256::new();
    hasher.update(token.as_bytes());
    hex::encode(hasher.finalize())
}
