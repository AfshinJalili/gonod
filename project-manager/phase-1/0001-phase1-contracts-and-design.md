# Ticket 0001: Phase 1 Contracts and Design Freeze

## Goal
Lock the Phase 1 API contracts and token/session rules before implementation so later tickets are straightforward and low-risk.

## Why This Ticket Exists
Current login only returns a message. Phase 1 needs concrete token/session behavior. If contracts are not frozen first, code churn will be high.

## Scope
- Define and document Phase 1 endpoints and payloads:
  - `POST /login`
  - `POST /token/refresh`
  - `POST /logout`
  - `GET /me`
  - `GET /.well-known/jwks.json`
- Define access token and refresh token lifetimes.
- Define JWT claim set for access token (`iss`, `sub`, `aud`, `exp`, `iat`, `jti`, `sid`).
- Define refresh token rotation behavior and reuse-detection policy.
- Define error codes and response envelope for token/session failures.
- Define env/config keys required for Phase 1.

## Out of Scope
- Writing implementation code.
- DB migrations.
- Adding production key management (KMS/HSM).

## Deliverables
- One design doc at:
  - `project-manager/phase-1/phase1-design-spec.md`
- Contract examples for all new/changed endpoints:
  - success response
  - validation failure
  - unauthorized/revoked/expired cases
- Security notes section with:
  - token TTL rationale
  - rotation/reuse rationale
  - logging redaction rules

## Suggested Implementation Checklist
- [ ] Write endpoint contract tables (request/response fields).
- [ ] Specify HTTP status codes per failure case.
- [ ] Specify token claim schema and signing algorithm (`RS256`).
- [ ] Specify refresh token lifecycle state machine.
- [ ] Specify config keys (names, defaults, required/optional).
- [ ] Document backward-compatibility impact for current `/login` clients.

## Test Requirements (Design-Level)
- [ ] Add at least 5 contract test scenarios in the doc that later tickets must automate.
- [ ] Include at least 3 negative security scenarios (tampered token, expired token, reused refresh token).

## Manual QA Requirements
- [ ] Review the spec line-by-line before code starts.
- [ ] Confirm no endpoint ambiguity remains (field names, status codes, error semantics).

## Dependencies
- None.

## Estimated Size
- Small (0.5 to 1 day).

## Definition of Done
- A complete Phase 1 design spec exists in `project-manager/phase-1/phase1-design-spec.md`.
- All five Phase 1 endpoints have explicit request/response contracts with sample JSON.
- Token/session lifecycle rules are unambiguous and implementation-ready.
- Config keys and defaults are documented and approved for implementation.
