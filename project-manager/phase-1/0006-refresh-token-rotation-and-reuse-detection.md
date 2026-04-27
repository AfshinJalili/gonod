# Ticket 0006: Refresh Token Rotation and Reuse Detection

## Goal
Implement `POST /token/refresh` with one-time-use refresh tokens, rotation, and theft/reuse detection.

## Scope
- Add refresh endpoint and service.
- Validate presented refresh token by hash lookup.
- Enforce one-time token use:
  - mark used token rotated/revoked
  - mint new refresh token and access token
- Detect refresh token reuse and revoke related session/token chain.
- Return standardized unauthorized error for invalid/revoked/expired cases.

## Out of Scope
- Global risk engine.
- User notification for suspicious reuse events (can be future ticket).

## Suggested Implementation Checklist
- [ ] Add refresh handler and route.
- [ ] Parse and validate request payload.
- [ ] Hash inbound token and load token record with session.
- [ ] Enforce expiry and revocation checks.
- [ ] Perform rotation in transaction to avoid race conditions.
- [ ] Implement reuse detection path:
  - if already rotated/revoked token appears again, revoke session and descendants
- [ ] Emit structured security/audit log event for reuse detection.

## Files Likely Touched
- `internal/server/routes.go`
- `internal/handler/auth.go` (or new token handler)
- `internal/service/*`
- `internal/repository/*`

## Test Requirements
- [ ] Unit tests:
  - expired token rejected
  - revoked token rejected
  - unknown token rejected
  - reused token triggers session revocation
- [ ] Integration tests:
  - successful rotation returns new tokens
  - old refresh token cannot be reused
  - concurrent refresh attempts allow only one winner
- [ ] Security regression tests for replay attempts.

## Manual QA Requirements
- [ ] Login -> refresh once -> verify success.
- [ ] Reuse old refresh token -> verify rejection.
- [ ] After reuse event, try new refresh token -> verify blocked if session revoked.

## Dependencies
- Depends on `0005`.

## Estimated Size
- Medium/Large (1.5 to 2 days).

## Definition of Done
- Refresh endpoint supports secure one-time rotation.
- Reuse detection is implemented and revokes compromised session path.
- Concurrency-safe behavior is proven by automated tests.
- No sensitive token values are logged.
