# Ticket 0003: Domain and Repository Layer for Sessions and Refresh Tokens

## Goal
Introduce domain models and repository interfaces/implementations for session and refresh-token persistence.

## Scope
- Add domain models:
  - `Session`
  - `RefreshToken`
  - `SigningKey` (if DB-backed in Phase 1)
- Add repository interfaces under `internal/domain`.
- Add postgres repository implementations under `internal/repository`.
- Add domain-level errors (not found, revoked, expired, reuse detected).

## Out of Scope
- JWT signing.
- HTTP endpoints.
- Middleware.

## Suggested Implementation Checklist
- [ ] Create domain entities with clear field naming and timestamps.
- [ ] Define repository interfaces with minimal methods required by Phase 1:
  - create session
  - revoke session
  - create refresh token
  - find refresh token by hash
  - rotate refresh token
  - revoke token chain/session on reuse
- [ ] Implement postgres repositories with parameterized SQL.
- [ ] Map SQL errors to domain errors consistently.
- [ ] Add context-aware queries with sane timeout strategy from caller side.

## Files Likely Touched
- `internal/domain/*.go`
- `internal/repository/*.go`
- `cmd/api/main.go` (wiring constructors)

## Test Requirements
- [ ] Unit tests for SQL error-to-domain error mapping.
- [ ] Integration tests for:
  - session create/revoke
  - refresh token insert/find/rotate
  - reuse path revoking session/token chain
- [ ] Concurrency test for refresh rotation race (single valid successor behavior).

## Manual QA Requirements
- [ ] Verify DB rows are written as expected for login and refresh operations.
- [ ] Verify revoked records cannot be reused by repository methods.

## Dependencies
- Depends on `0002` migrations.

## Estimated Size
- Medium (1 to 1.5 days).

## Definition of Done
- Domain models and repository interfaces exist with clear contracts.
- Postgres implementations support all persistence operations required by Phase 1.
- Error semantics are deterministic and tested.
- Integration tests prove correctness for normal and edge cases.
