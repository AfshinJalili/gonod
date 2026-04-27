# Ticket 0007: Bearer Auth Middleware and `/me` Endpoint

## Goal
Introduce authenticated request handling using access tokens and add `GET /me`.

## Scope
- Add bearer auth middleware:
  - parse `Authorization: Bearer <token>`
  - validate token signature/claims via token service
  - attach principal/session context to request
- Add `GET /me` endpoint returning authenticated user/session claims.
- Add middleware error responses for:
  - missing token
  - malformed token
  - invalid/expired token

## Out of Scope
- Role/permission authorization (Phase 2).
- User profile expansion beyond essential claims.

## Suggested Implementation Checklist
- [ ] Create auth middleware package/file.
- [ ] Define principal context struct (`user_id`, `session_id`, `subject`, `issued_at`, `expires_at`).
- [ ] Add helper functions for context read/write.
- [ ] Register `/me` in routes and protect it with auth middleware.
- [ ] Return deterministic JSON error envelope on auth failures.
- [ ] Ensure middleware does not panic on malformed headers.

## Files Likely Touched
- `internal/middleware/*`
- `internal/server/routes.go`
- `internal/handler/*`

## Test Requirements
- [ ] Unit tests for middleware parsing/validation branches.
- [ ] Integration tests:
  - `/me` succeeds with valid token
  - `/me` fails for missing/invalid/expired token
- [ ] Regression test to ensure auth context is isolated per request.

## Manual QA Requirements
- [ ] Call `/me` without token -> expect 401.
- [ ] Call `/me` with valid token -> expect correct user identity fields.
- [ ] Call `/me` with tampered token -> expect 401.

## Dependencies
- Depends on `0004` and `0005`.

## Estimated Size
- Medium (1 day).

## Definition of Done
- Auth middleware validates bearer tokens and injects principal context.
- `/me` is functional and protected.
- All unauthorized branches return consistent API errors.
- Tests cover success and failure paths thoroughly.
