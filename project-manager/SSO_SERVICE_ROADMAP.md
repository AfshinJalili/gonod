# Gonod SSO Service Roadmap (Keycloak-Like)

**Last Updated:** 2026-02-27  
**Audience:** Founder/Engineering team  
**Current Stage:** Early Auth API (register/login)  
**Target Stage:** Production-grade, multi-tenant IAM/SSO platform

## 1. Current State Review (From Full Repository Read)

### What exists now
- Go HTTP API with clear layers: `handler -> service -> repository -> postgres`.
- Endpoints:
  - `POST /register`
  - `POST /login`
  - `GET /health`
- Postgres setup with migrations and simple `users` table.
- Password hashing with `bcrypt`.
- Basic middleware for logging and panic recovery.
- Structured logging via `slog`.

### Key gaps before real SSO work
- No tests at all (`go test ./...` shows `[no test files]` everywhere).
- No token issuance (no JWT, refresh token, session lifecycle).
- Not OAuth2 / OIDC yet (no authorization endpoint, token endpoint, discovery, JWKS, userinfo).
- No multi-tenant/realm model.
- No client/app registration model.
- No roles/permissions model.
- No admin API/UI.
- Limited error model (validation details not returned to client yet).
- Missing secure defaults for production (timeouts, shutdown, secret management, rotation strategy, audit trail).

## 2. Product Target (What “Keycloak-like” Means)

## Core Product Capabilities
- Act as Identity Provider (IdP) and Authorization Server.
- Support OAuth2.1-style secure flows (especially Authorization Code + PKCE).
- Support OpenID Connect (OIDC) for authentication.
- Issue/validate access tokens, refresh tokens, ID tokens.
- Expose standard OIDC metadata and JWKS.
- Support users, groups, roles, and fine-grained authorization policies.
- Provide admin controls for realms/tenants, clients, users, and authentication flows.
- Provide enterprise integrations (SAML, LDAP/AD, SCIM, social login).

## Non-Functional Goals
- Secure-by-default, observable, testable, scalable, operable.
- Clear upgrade/migration path and strong backward compatibility contracts.

## 3. Roadmap Principles

- Build in **vertical slices**: feature + tests + observability + docs each phase.
- Do not add advanced features on unstable core.
- Every phase must have:
  - Unit tests
  - Integration tests
  - Manual QA checklist
  - Security checks
  - Rollback/migration notes
- Use standards first (RFC-aligned behavior) before custom behavior.

## 4. Multi-Phase Roadmap (Simple -> Enterprise)

## Phase 0: Foundation Hardening (Now -> “Safe Core”)
**Goal:** Make current auth API production-safe enough to build on.

### Features/Changes
- Config validation at startup (fail fast if required env missing).
- HTTP server hardening:
  - Read timeout
  - Write timeout
  - Idle timeout
  - Graceful shutdown with context deadline
- Consistent API error envelope including validation details.
- Request ID/correlation ID middleware.
- Normalize and canonicalize email handling.
- Add readiness endpoint (`/ready`) with DB health check.
- Add migration command path (separate app startup vs migration-only command).
- Remove/guard `cmd/scratch` from production workflow.

### Data/Schema
- Add fields for account status (`is_active`, `email_verified`, `locked_until`).
- Add created/updated trigger strategy or explicit updates.

### Required Tests
- Unit tests:
  - `AuthRequest.Validate`
  - handler responses (status, body schema, error details)
  - middleware behavior (recovery/logging/request-id)
- Integration tests:
  - register/login against test Postgres
  - duplicate email conflict mapping
- Basic API contract tests for JSON response shape.

### Manual QA
- Register valid/invalid payloads.
- Login success/failure.
- DB down -> readiness fails, health still alive (if desired split).
- SIGTERM graceful shutdown behavior.

### Exit Criteria
- `go test ./...` with meaningful coverage on all current layers.
- Stable error contract and basic operability in local/docker.

---

## Phase 1: Token & Session Core
**Goal:** Convert login from “credential check” to real session/token issuance.

### Features/Changes
- Access token + refresh token model.
- Token signing strategy (start with asymmetric keys; publish JWKS early).
- Refresh token rotation + reuse detection.
- Session store (per-device/session metadata, revocation support).
- Logout endpoint (session invalidation).
- Basic `/me` endpoint from token claims.

### Data/Schema
- `sessions` table
- `refresh_tokens` table (hashed token storage)
- `signing_keys` table or secure key source integration

### Required Tests
- Unit tests for token service, expiry logic, rotation behavior.
- Integration tests for login/refresh/logout flows.
- Negative tests for expired, tampered, revoked tokens.

### Manual QA
- Login -> use access token -> refresh -> old refresh rejected.
- Logout invalidates session/token use.

### Exit Criteria
- Stateless access token + stateful refresh/session lifecycle implemented and proven.

---

## Phase 2: Authorization Fundamentals (RBAC)
**Goal:** Add first-class authorization beyond authentication.

### Features/Changes
- Roles and permissions model.
- User-role assignment.
- Route-level authorization middleware/policy check.
- Admin APIs for role management.

### Data/Schema
- `roles`, `permissions`, `user_roles`, `role_permissions`

### Required Tests
- Policy evaluation unit tests.
- API integration tests for protected endpoints by role.
- Regression tests for unauthorized/forbidden semantics.

### Manual QA
- Verify least privilege behavior.
- Verify role updates apply to new sessions/tokens as designed.

### Exit Criteria
- Predictable, tested RBAC with clear admin APIs.

---

## Phase 3: OAuth2 Authorization Server (Developer-Usable)
**Goal:** Become an OAuth2-compliant auth server for applications.

### Features/Changes
- Client registration model:
  - confidential/public clients
  - redirect URIs
  - grant types
  - scopes
- Implement Authorization Code + PKCE.
- Token endpoint with standard responses/errors.
- Introspection + revocation endpoints.
- Consent screen model (minimal first pass).

### Data/Schema
- `clients`, `client_secrets`, `auth_codes`, `consents`, `scopes`

### Required Tests
- Flow tests for auth code + PKCE success/failure.
- Redirect URI validation tests.
- Client auth tests (basic/post/private key JWT if added).
- RFC-compatible error response tests.

### Manual QA
- Use Postman/sample app to complete full OAuth login.
- Test invalid code/verifier/redirect URI/client credentials.

### Exit Criteria
- External app can integrate without custom hacks.

---

## Phase 4: OIDC Provider Core
**Goal:** Become standards-compliant OpenID Connect provider.

### Features/Changes
- Issue ID tokens (claims + nonce handling).
- `/.well-known/openid-configuration`
- `/.well-known/jwks.json`
- `/userinfo`
- Support scopes: `openid`, `profile`, `email`, `offline_access`.
- Subject identifier strategy (`sub`) and claim mapping.

### Required Tests
- ID token claim validation tests (iss, aud, exp, iat, nonce).
- Discovery/JWKS contract tests.
- Userinfo authorization/scope tests.

### Manual QA
- Integrate with an OIDC relying party and validate full login.
- Rotate keys and verify old/new token validation behavior.

### Exit Criteria
- Third-party OIDC clients can integrate reliably.

---

## Phase 5: Multi-Tenant (Realm) Architecture
**Goal:** Move from single-tenant auth API to true IAM platform.

### Features/Changes
- Realm/tenant model with strict data isolation.
- Realm-scoped users, clients, roles, keys, auth flows.
- Realm discovery and domain mapping.
- Admin boundary controls (global admin vs realm admin).

### Data/Schema
- `realms` and realm_id on all security-sensitive entities.

### Required Tests
- Cross-tenant isolation tests (critical).
- Authorization tests for admin boundary enforcement.
- Tenant migration/backfill tests.

### Manual QA
- Verify no cross-tenant user/client/token access paths.

### Exit Criteria
- Strong isolation guarantees with test evidence.

---

## Phase 6: Advanced Authentication
**Goal:** Add stronger auth controls needed by real companies.

### Features/Changes
- MFA (TOTP first, then WebAuthn/passkeys).
- Email verification flow.
- Password reset + account recovery.
- Brute-force protection:
  - rate limiting
  - progressive lockout
  - IP/device risk flags

### Required Tests
- MFA enrollment/challenge tests.
- Recovery and reset token tests.
- Abuse/rate-limit behavior tests.

### Manual QA
- Full auth journeys: normal, locked, recovery, MFA bypass attempts.

### Exit Criteria
- Secure user lifecycle and fraud-resistance baseline.

---

## Phase 7: Identity Federation
**Goal:** Integrate external identity sources.

### Features/Changes
- Social login (Google/GitHub) via OIDC.
- Enterprise IdP brokering (SAML/OIDC).
- Optional LDAP/AD user federation.
- Account linking and conflict resolution policy.

### Required Tests
- Brokered login tests for each provider.
- Mapping and account-linking tests.
- Replay/invalid assertion/token hardening tests.

### Manual QA
- End-to-end login via each configured external IdP.

### Exit Criteria
- Organizations can onboard without replacing existing IdP immediately.

---

## Phase 8: Admin Plane & Developer Experience
**Goal:** Make service operable and usable by others.

### Features/Changes
- Admin API (versioned) + initial admin web console.
- Audit log viewer + export.
- Developer portal/docs:
  - quickstart
  - sample apps
  - SDK snippets
- Client secret rotation UX.

### Required Tests
- Admin permission boundary tests.
- API versioning and backward compatibility tests.
- UI e2e smoke tests (if console exists).

### Manual QA
- Full tenant onboarding via admin API/UI only.

### Exit Criteria
- Another team can self-serve integration.

---

## Phase 9: Security, Compliance, and Governance
**Goal:** Reach enterprise trust requirements.

### Features/Changes
- Key management with rotation policy and KMS integration.
- Immutable audit logs for auth/admin events.
- Security headers, CSRF protection (where relevant), strict CORS policy.
- Threat model and abuse case coverage.
- Compliance controls prep (SOC2/ISO style evidence mapping).

### Required Tests
- Security test suite (authz bypass, token confusion, replay attempts).
- Dependency and SAST scanning in CI.
- Audit event completeness tests.

### Manual QA
- Security checklist walkthrough and incident-response tabletop.

### Exit Criteria
- Documented controls and tested security posture.

---

## Phase 10: Scalability, HA, and SRE Maturity
**Goal:** Run reliably under high traffic and failure scenarios.

### Features/Changes
- Horizontal scaling strategy (stateless API + shared session/token state).
- Caching strategy (JWKS, session lookups, consent, metadata).
- Background jobs/workers (cleanup, key rotation, event delivery).
- Metrics, traces, alerts, SLOs.
- Backup/restore drills and disaster recovery plan.

### Required Tests
- Load/stress tests for auth/token endpoints.
- Chaos/failure injection tests (DB partial outage, cache outage).
- DR restore validation tests.

### Manual QA
- Runbook execution drills with synthetic incidents.

### Exit Criteria
- Known capacity envelope + proven recovery workflows.

---

## Phase 11: Enterprise Extensions
**Goal:** Deliver high-value enterprise IAM capabilities.

### Features/Changes
- SCIM provisioning/deprovisioning.
- Fine-grained authz policies (resource/action/context-based).
- Delegated administration and approval workflows.
- Event hooks/webhooks for IAM lifecycle.
- Data residency options and tenant-level policy packs.

### Required Tests
- SCIM compatibility tests.
- Policy engine correctness and performance tests.
- Webhook reliability tests (retry/idempotency/signature).

### Manual QA
- Enterprise onboarding simulation with provisioning and policy setup.

### Exit Criteria
- Competitive enterprise feature set for serious adoption.

## 5. Recommended Execution Cadence

- **Cadence:** 2-week iterations with one vertical slice each.
- **Release train:** Monthly tagged releases.
- **Quality gate each release:**
  - 0 critical vulnerabilities open
  - all critical paths covered by automated tests
  - manual QA checklist completed
  - migration rollback tested

## 6. Immediate Next 30-Day Plan (Practical Start)

## Week 1
- Finish Phase 0 design doc.
- Add test harness (unit + integration scaffolding).
- Add server timeout + graceful shutdown + config validation.

## Week 2
- Standardize error model + validation detail responses.
- Add readiness endpoint and DB health checks.
- Add first integration tests (register/login).

## Week 3
- Implement token service skeleton (access + refresh schema).
- Introduce signing key abstraction and JWKS endpoint stub.

## Week 4
- Finish login token issuance and refresh endpoint MVP.
- Run manual QA checklist and patch reliability issues.

## 7. Definition of Done (For Every Future Feature)

- Feature implemented with clean architecture boundaries.
- Unit tests + integration tests + negative tests added.
- Manual QA evidence captured.
- Logging/metrics for operational visibility added.
- API docs updated.
- Security and migration impacts reviewed.

## 8. Critical Risks to Manage Early

- Implementing OAuth/OIDC without strict standards validation.
- Skipping tests while adding security-sensitive features.
- Weak tenant isolation design introduced too late.
- No key rotation/revocation model from day one of tokens.
- Growing endpoints without versioning and compatibility policy.

---

This roadmap is intentionally ordered so each phase enables the next safely.  
If you want, next step is creating **Phase 0 execution tickets** (file-by-file implementation checklist for this exact repo).
