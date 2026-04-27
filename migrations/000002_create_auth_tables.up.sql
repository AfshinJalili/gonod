CREATE TABLE IF NOT EXISTS sessions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid (),
    user_id UUID NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    expires_at TIMESTAMPTZ NOT NULL,
    revoked_at TIMESTAMPTZ NULL,
    user_agent TEXT NULL,
    ip INET NULL,
    CONSTRAINT fk_sessions_user FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE NO ACTION,
    CONSTRAINT chk_sessions_expires_after_created CHECK (expires_at > created_at),
    CONSTRAINT chk_sessions_revoked_after_created CHECK (
        revoked_at IS NULL
        OR revoked_at >= created_at
    )
);

CREATE INDEX IF NOT EXISTS idx_sessions_user_revoked_expires ON sessions (
    user_id,
    revoked_at,
    expires_at
);

CREATE INDEX IF NOT EXISTS idx_sessions_active_expires ON sessions (expires_at)
WHERE
    revoked_at IS NULL;

CREATE TABLE IF NOT EXISTS refresh_tokens (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid (),
    session_id UUID NOT NULL,
    token_hash TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    expires_at TIMESTAMPTZ NOT NULL,
    revoked_at TIMESTAMPTZ NULL,
    rotated_at TIMESTAMPTZ NULL,
    replaced_by_token_id UUID NULL,
    CONSTRAINT fk_refresh_tokens_session FOREIGN KEY (session_id) REFERENCES sessions (id) ON DELETE CASCADE,
    CONSTRAINT fk_refresh_tokens_replaced_by FOREIGN KEY (replaced_by_token_id) REFERENCES refresh_tokens (id) ON DELETE SET NULL,
    CONSTRAINT uq_refresh_tokens_token_hash UNIQUE (token_hash),
    CONSTRAINT chk_refresh_tokens_expires_after_created CHECK (expires_at > created_at),
    CONSTRAINT chk_refresh_tokens_revoked_after_created CHECK (
        revoked_at IS NULL
        OR revoked_at >= created_at
    ),
    CONSTRAINT chk_refresh_tokens_rotated_after_created CHECK (
        rotated_at IS NULL
        OR rotated_at >= created_at
    ),
    CONSTRAINT chk_refresh_tokens_not_self CHECK (
        replaced_by_token_id IS NULL
        OR replaced_by_token_id <> id
    )
);

CREATE INDEX IF NOT EXISTS idx_refresh_tokens_session_revoked_expires ON refresh_tokens (
    session_id,
    revoked_at,
    expires_at
);

CREATE INDEX IF NOT EXISTS idx_refresh_tokens_active_expires ON refresh_tokens (expires_at)
WHERE
    revoked_at IS NULL;

CREATE TABLE IF NOT EXISTS signing_keys (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid (),
    kid TEXT NOT NULL,
    alg TEXT NOT NULL,
    public_key_pem TEXT NOT NULL,
    private_key_pem TEXT NOT NULL,
    status TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT uq_signing_keys_kid UNIQUE (kid),
    CONSTRAINT chk_signing_keys_alg CHECK (alg IN ('RS256')),
    CONSTRAINT chk_signing_keys_status CHECK (
        status IN ('active', 'retired')
    )
);

CREATE INDEX IF NOT EXISTS idx_signing_keys_status ON signing_keys (status);

CREATE UNIQUE INDEX IF NOT EXISTS uq_signing_keys_single_active ON signing_keys (status)
WHERE
    status = 'active';