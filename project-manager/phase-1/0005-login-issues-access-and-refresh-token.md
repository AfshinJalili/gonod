# Ticket 0005: Login Issues Access Token and Refresh Token

## Goal
Change `POST /login` from "credentials check only" to full session creation with token issuance.

## Scope
- Update login service flow:
  - validate credentials
  - create session record
  - create refresh token record (store hash only)
  - issue signed access token
  - return access + refresh token payload
- Update handler response format to contract from Ticket `0001`.
- Capture device metadata (`ip`, `user-agent`) on session creation.

## Out of Scope
- Refresh rotation.
- Logout.
- `/me` endpoint.

## Suggested Implementation Checklist
- [ ] Extend service layer with `LoginWithTokens(...)`.
- [ ] Add secure refresh token generator (`crypto/rand` based).
- [ ] Hash refresh token before persistence (do not store raw token).
- [ ] Return raw refresh token once in response body.
- [ ] Ensure old login behavior is removed or versioned cleanly.
- [ ] Add structured audit log event for successful/failed login without sensitive values.

## Files Likely Touched
- `internal/service/user_service.go` (or new auth service)
- `internal/handler/auth.go`
- `internal/repository/*`
- `internal/handler/response.go`

## Test Requirements
- [ ] Unit tests for login flow branches:
  - invalid credentials
  - session create failure
  - token issuance failure
- [ ] Integration tests:
  - successful login writes session + refresh token row
  - response includes token fields and TTL
  - no raw refresh token persisted in DB
- [ ] Regression tests for old error semantics where applicable.

## Manual QA Requirements
- [ ] Perform login with valid credentials and inspect JSON token payload.
- [ ] Verify access token works for protected endpoint once Ticket `0007` lands.
- [ ] Verify invalid login does not create session records.

## Dependencies
- Depends on `0003` and `0004`.

## Estimated Size
- Medium (1 to 1.5 days).

## Definition of Done
- `/login` creates session and returns access + refresh tokens on success.
- Refresh token is stored hashed only.
- Login failure paths are deterministic and tested.
- Integration tests prove persistence and response behavior.
