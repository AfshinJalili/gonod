# Ticket 0002: DB Migrations for Sessions, Refresh Tokens, and Signing Keys

## Goal
Add the minimum durable schema needed for Phase 1 token and session lifecycle.

## Scope
- Add migration(s) for:
  - `sessions`
  - `refresh_tokens`
  - `signing_keys` (or equivalent local key registry table)
- Add foreign keys and indexes needed for runtime performance.
- Add down migration(s) for rollback.

## Out of Scope
- Endpoint handlers.
- Token issuance logic.
- Key rotation automation.

## Suggested Schema Baseline
- `sessions`
  - `id UUID PK`
  - `user_id UUID NOT NULL REFERENCES users(id)`
  - `created_at`, `updated_at`, `expires_at`
  - `revoked_at NULL`
  - `user_agent TEXT NULL`
  - `ip TEXT NULL`
- `refresh_tokens`
  - `id UUID PK`
  - `session_id UUID NOT NULL REFERENCES sessions(id)`
  - `token_hash TEXT NOT NULL UNIQUE`
  - `created_at`, `expires_at`
  - `revoked_at NULL`
  - `rotated_at NULL`
  - `replaced_by_token_id UUID NULL`
- `signing_keys`
  - `id UUID PK`
  - `kid TEXT NOT NULL UNIQUE`
  - `alg TEXT NOT NULL`
  - `public_key_pem TEXT NOT NULL`
  - `private_key_pem TEXT NOT NULL` (temporary local-dev strategy)
  - `status TEXT NOT NULL` (`active`/`retired`)
  - `created_at`, `updated_at`

## Suggested Implementation Checklist
- [ ] Create `migrations/000002_*.up.sql` for new tables.
- [ ] Create matching `migrations/000002_*.down.sql`.
- [ ] Add key indexes:
  - sessions by `user_id`, `revoked_at`, `expires_at`
  - refresh tokens by `session_id`, `token_hash`, `revoked_at`, `expires_at`
  - signing keys by `kid`, `status`
- [ ] Validate FK cascade behavior does not accidentally delete audit-critical data.
- [ ] Document migration rollback caveats in the file header comments.

## Test Requirements
- [ ] Migration smoke test: up then down then up in a clean database.
- [ ] Integration test asserting inserts/selects on all new tables.
- [ ] Test duplicate `token_hash` and duplicate `kid` fail as expected.

## Manual QA Requirements
- [ ] Run migrations locally and inspect resulting schema with `psql`.
- [ ] Validate indexes exist and are named predictably.

## Dependencies
- Depends on `0001` for finalized schema semantics.

## Estimated Size
- Medium (1 day).

## Definition of Done
- New tables and indexes are created by migrations successfully.
- Down migration cleanly rolls back Phase 1 tables.
- Schema supports session revocation and refresh token rotation without ambiguity.
- Migration behavior is verified by automated test or scripted validation.
