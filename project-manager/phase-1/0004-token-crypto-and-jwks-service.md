# Ticket 0004: Token Crypto Service and JWKS Builder

## Goal
Build a dedicated token service that signs/verifies access tokens and exposes public key material for JWKS.

## Scope
- Introduce `internal/service/token_service.go` (or dedicated package).
- Add signing using `RS256`.
- Add token verification/parsing utility used by middleware.
- Add `kid` handling in JWT header.
- Add JWKS document builder from active key(s).
- Add config for issuer, audience, and token TTL.

## Out of Scope
- Full OIDC discovery document.
- Key rotation scheduler.
- External KMS/HSM integration.

## Suggested Implementation Checklist
- [ ] Define `TokenService` interface:
  - `IssueAccessToken(...)`
  - `ParseAndValidateAccessToken(...)`
  - `CurrentJWKS()`
- [ ] Implement claims model aligned with Ticket `0001`.
- [ ] Add clock abstraction to make expiry testing deterministic.
- [ ] Add key-loading strategy for Phase 1:
  - DB-backed active key or startup config key pair
- [ ] Ensure signature verification checks:
  - algorithm allowlist
  - issuer
  - audience
  - expiration and not-before (if used)
- [ ] Ensure no private key material is ever logged.

## Files Likely Touched
- `internal/service/*`
- `internal/config/config.go`
- `cmd/api/main.go`

## Test Requirements
- [ ] Unit tests:
  - successful issue and verify
  - wrong issuer/audience rejected
  - expired token rejected
  - algorithm confusion rejected
  - unknown `kid` rejected
- [ ] Property-style test for token expiry boundaries (now +/- skew).
- [ ] JWKS output test validates key shape (`kty`, `kid`, `use`, `alg`, `n`, `e`).

## Manual QA Requirements
- [ ] Generate token from login and verify with external JWT tool.
- [ ] Fetch JWKS and verify token signature using returned public key.

## Dependencies
- Depends on `0001`.
- Depends on `0003` if key metadata is stored in DB.

## Estimated Size
- Medium (1.5 to 2 days).

## Definition of Done
- Access token can be issued and validated reliably with RS256.
- JWKS can be generated from active signing key(s).
- Failure modes for invalid/tampered/expired tokens are deterministic.
- Unit tests cover all critical validation branches.
