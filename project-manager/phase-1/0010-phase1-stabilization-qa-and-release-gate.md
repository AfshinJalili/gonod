# Ticket 0010: Phase 1 Stabilization, QA, and Release Gate

## Goal
Finalize Phase 1 with a strict quality gate so token/session core is stable before Phase 2.

## Scope
- Add/complete Phase 1 automated tests across all new paths.
- Add manual QA checklist with expected outcomes.
- Add short release checklist and rollback plan.
- Update docs for new endpoints and env vars.

## Out of Scope
- New feature development.
- RBAC/authorization (Phase 2).

## Suggested Implementation Checklist
- [ ] Add unit tests for:
  - token issuance/verification
  - refresh rotation
  - reuse detection
  - auth middleware branches
- [ ] Add integration tests for:
  - login -> me -> refresh -> logout lifecycle
  - logout blocks subsequent refresh
  - revoked/reused token rejection paths
- [ ] Add API contract tests for response shape consistency.
- [ ] Document manual QA script at:
  - `project-manager/phase-1/phase1-manual-qa.md`
- [ ] Document release and rollback steps at:
  - `project-manager/phase-1/phase1-release-checklist.md`

## Required Commands
- `go test ./...`
- migration up/down smoke command(s)
- optional targeted integration command(s) if split by build tags

## Manual QA Requirements
- [ ] Fresh user login returns access + refresh tokens.
- [ ] `/me` requires valid access token.
- [ ] Refresh rotation invalidates previous refresh token.
- [ ] Reuse detection revokes session chain.
- [ ] Logout revokes active session/tokens.
- [ ] JWKS is reachable and matches issued token `kid`.

## Dependencies
- Depends on completion of `0001` through `0009`.

## Estimated Size
- Medium (1 to 1.5 days).

## Definition of Done
- All Phase 1 tickets are completed and verified.
- Automated tests for Phase 1 pass in local CI command (`go test ./...`).
- Manual QA checklist executed with results recorded.
- Release checklist and rollback plan are documented.
- Phase 1 is ready to hand off to Phase 2 without known critical gaps.
