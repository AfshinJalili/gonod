# Ticket 0008: Logout Endpoint and Session Revocation

## Goal
Implement `POST /logout` so active sessions can be revoked intentionally.

## Scope
- Add logout endpoint requiring authentication context.
- Revoke current session and active refresh tokens for that session.
- Define idempotent logout behavior (repeat calls should remain safe).
- Optional: support "logout all sessions for user" behind explicit request flag.

## Out of Scope
- Frontchannel/backchannel logout protocols.
- Federated logout with external providers.

## Suggested Implementation Checklist
- [ ] Add route and handler for logout.
- [ ] Extract session ID from auth context.
- [ ] Mark session revoked and revoke related refresh tokens in one transaction.
- [ ] Return standard success response (no sensitive details).
- [ ] Ensure access tokens already issued naturally expire; no blacklist needed yet unless design requires.
- [ ] Add audit log event `session.logout`.

## Files Likely Touched
- `internal/server/routes.go`
- `internal/handler/*`
- `internal/service/*`
- `internal/repository/*`

## Test Requirements
- [ ] Unit tests for revoke service logic.
- [ ] Integration tests:
  - logout revokes session
  - refresh after logout fails
  - repeated logout remains idempotent
- [ ] Negative tests for missing auth context.

## Manual QA Requirements
- [ ] Login, then logout, then try refresh -> verify rejection.
- [ ] Call logout twice -> both calls should be safe and predictable.

## Dependencies
- Depends on `0006` and `0007`.

## Estimated Size
- Medium (1 day).

## Definition of Done
- `/logout` revokes current session and blocks further refresh token use.
- Endpoint is idempotent and has stable response semantics.
- Revoke behavior is validated by integration tests.
- Audit logging exists for logout events.
