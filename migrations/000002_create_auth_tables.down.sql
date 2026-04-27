DROP INDEX IF EXISTS uq_signing_keys_single_active;

DROP INDEX IF EXISTS idx_signing_keys_status;

DROP TABLE IF EXISTS signing_keys;

DROP INDEX IF EXISTS idx_refresh_tokens_active_expires;

DROP INDEX IF EXISTS idx_refresh_tokens_session_revoked_expires;

DROP TABLE IF EXISTS refresh_tokens;

DROP INDEX IF EXISTS idx_sessions_active_expires;

DROP INDEX IF EXISTS idx_sessions_user_revoked_expires;

DROP TABLE IF EXISTS sessions;