# Ticket 0009: JWKS Endpoint and Key Rotation Baseline

## Goal
Expose signing public keys via JWKS and establish a minimal manual key-rotation procedure.

## Scope
- Add `GET /.well-known/jwks.json` endpoint.
- Publish active (and optionally recently retired) public keys.
- Ensure JWTs include `kid` matching JWKS.
- Document a manual key-rotation playbook for local/staging/prod.

## Out of Scope
- Automated rotation scheduler.
- External KMS-managed key lifecycle.

## Suggested Implementation Checklist
- [ ] Add JWKS handler and route.
- [ ] Wire JWKS response from token/key service.
- [ ] Ensure response caching strategy is defined (`Cache-Control`).
- [ ] Validate token issuance uses currently active `kid`.
- [ ] Create operational runbook:
  - add new key
  - mark key active
  - serve old+new key overlap window
  - retire old key after access token max TTL

## Files Likely Touched
- `internal/server/routes.go`
- `internal/handler/*`
- `internal/service/*`
- `project-manager/phase-1/key-rotation-runbook.md`

## Test Requirements
- [ ] Unit test JWKS serialization format.
- [ ] Integration test endpoint response shape and status.
- [ ] Rotation test:
  - token signed with old key remains verifiable during overlap
  - new token signed with new key uses new `kid`

## Manual QA Requirements
- [ ] Fetch JWKS and verify expected key count and `kid` values.
- [ ] Issue token before and after key switch and validate both during overlap.

## Dependencies
- Depends on `0004`.
- Recommended after `0005` for real token output verification.

## Estimated Size
- Small/Medium (0.5 to 1 day).

## Definition of Done
- `/.well-known/jwks.json` returns valid JWKS payload.
- JWT header `kid` is consistent with published key set.
- Manual rotation procedure is documented and tested in non-prod.
