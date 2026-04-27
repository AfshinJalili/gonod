# Phase 1 Design Spec: Contracts and Session/Token Rules

Status: Design freeze for Phase 1 implementation (Tickets `0002`-`0010`)  
Version: `1.0`  
Last Updated: `2026-02-27`

## 1. Scope and Goals

This document freezes Phase 1 API contracts and token/session behavior for:

- `POST /login`
- `POST /token/refresh`
- `POST /logout`
- `GET /me`
- `GET /.well-known/jwks.json`

This spec is implementation-ready and is the source of truth for endpoint payloads, JWT claims, token/session lifecycle, and error semantics.

## 2. API Envelope Contract

All Phase 1 endpoints use the existing envelope shape with explicit error codes for token/session behavior.

### 2.1 Success Envelope

```json
{
  "data": {}
}
```

### 2.2 Error Envelope

```json
{
  "error": {
    "code": "ERROR_CODE",
    "message": "Human readable error",
    "status": 401,
    "details": {
      "field": "reason"
    }
  }
}
```

Rules:

- `error.code` is mandatory for Phase 1 auth/token/session errors.
- `error.details` is optional and must not include secrets.
- `error.status` must mirror the HTTP status code.

## 3. Endpoint Contracts

### 3.1 `POST /login`

Authenticates credentials, creates a session, issues an access token and refresh token.

### Request

Headers:

- `Content-Type: application/json`

Body fields:

| Field | Type | Required | Rules |
|---|---|---|---|
| `email` | string | yes | Lowercased and trimmed; valid email format |
| `password` | string | yes | Trimmed; minimum 8 characters |

Example:

```json
{
  "email": "user@example.com",
  "password": "Password123!"
}
```

### Success (`200`)

```json
{
  "data": {
    "message": "Login Successful",
    "access_token": "<jwt>",
    "token_type": "Bearer",
    "expires_in": 900,
    "refresh_token": "<opaque-refresh-token>",
    "refresh_expires_in": 2592000,
    "session_id": "0f4c4c67-a75d-4f72-9a2a-d8f0059a2f3e"
  }
}
```

### Validation Failure (`400`)

```json
{
  "error": {
    "code": "VALIDATION_ERROR",
    "message": "Invalid request payload",
    "status": 400,
    "details": {
      "email": "must be a valid email format"
    }
  }
}
```

### Unauthorized (`401`)

```json
{
  "error": {
    "code": "INVALID_CREDENTIALS",
    "message": "Invalid credentials",
    "status": 401
  }
}
```

### 3.2 `POST /token/refresh`

Consumes a one-time refresh token, rotates it, and returns a new access token + refresh token.

### Request

Headers:

- `Content-Type: application/json`

Body fields:

| Field | Type | Required | Rules |
|---|---|---|---|
| `refresh_token` | string | yes | Opaque token previously returned by `/login` or `/token/refresh` |

Example:

```json
{
  "refresh_token": "<opaque-refresh-token>"
}
```

### Success (`200`)

```json
{
  "data": {
    "access_token": "<jwt>",
    "token_type": "Bearer",
    "expires_in": 900,
    "refresh_token": "<new-opaque-refresh-token>",
    "refresh_expires_in": 2592000,
    "session_id": "0f4c4c67-a75d-4f72-9a2a-d8f0059a2f3e"
  }
}
```

### Validation Failure (`400`)

```json
{
  "error": {
    "code": "VALIDATION_ERROR",
    "message": "Invalid request payload",
    "status": 400,
    "details": {
      "refresh_token": "is required"
    }
  }
}
```

### Unauthorized/Revoked/Expired/Reused (`401`)

Security behavior is intentionally normalized to avoid token-state oracle leaks.

```json
{
  "error": {
    "code": "REFRESH_TOKEN_INVALID",
    "message": "Invalid refresh token",
    "status": 401
  }
}
```

Notes:

- Returned for unknown, expired, revoked, tampered, or reused refresh token.
- On detected reuse, the server must revoke the full session/token chain before returning this response.

### 3.3 `POST /logout`

Revokes the current session and active refresh tokens for that session.

### Request

Headers:

- `Authorization: Bearer <access_token>`

Body:

- Empty body for Phase 1.

### Success (`200`)

```json
{
  "data": {
    "message": "Logout successful"
  }
}
```

### Validation Failure (`400`)

```json
{
  "error": {
    "code": "VALIDATION_ERROR",
    "message": "Invalid request payload",
    "status": 400
  }
}
```

### Unauthorized/Expired (`401`)

```json
{
  "error": {
    "code": "ACCESS_TOKEN_INVALID",
    "message": "Authentication failed",
    "status": 401
  }
}
```

Notes:

- Logout is idempotent. If session is already revoked, return `200`.

### 3.4 `GET /me`

Returns authenticated principal data from validated access token claims.

### Request

Headers:

- `Authorization: Bearer <access_token>`

### Success (`200`)

```json
{
  "data": {
    "user_id": "3af4f5e9-4d55-4f18-b4a4-3985338a37f2",
    "session_id": "0f4c4c67-a75d-4f72-9a2a-d8f0059a2f3e",
    "subject": "3af4f5e9-4d55-4f18-b4a4-3985338a37f2",
    "issuer": "https://auth.gonod.local",
    "audience": "gonod-api",
    "issued_at": 1772160000,
    "expires_at": 1772160900
  }
}
```

### Validation Failure (`400`)

```json
{
  "error": {
    "code": "VALIDATION_ERROR",
    "message": "Invalid Authorization header",
    "status": 400
  }
}
```

### Unauthorized/Expired (`401`)

```json
{
  "error": {
    "code": "ACCESS_TOKEN_INVALID",
    "message": "Authentication failed",
    "status": 401
  }
}
```

### 3.5 `GET /.well-known/jwks.json`

Returns public signing keys for access-token signature verification.

### Request

- No auth required.

### Success (`200`)

```json
{
  "keys": [
    {
      "kty": "RSA",
      "kid": "phase1-active-key",
      "use": "sig",
      "alg": "RS256",
      "n": "<base64url-modulus>",
      "e": "AQAB"
    }
  ]
}
```

### Service Unavailable (`503`)

```json
{
  "error": {
    "code": "JWKS_UNAVAILABLE",
    "message": "JWKS unavailable",
    "status": 503
  }
}
```

## 4. Access Token Contract (JWT)

### 4.1 Header

| Field | Value |
|---|---|
| `alg` | `RS256` |
| `typ` | `JWT` |
| `kid` | Active signing key ID |

### 4.2 Claims

| Claim | Type | Required | Meaning |
|---|---|---|---|
| `iss` | string | yes | Token issuer from config |
| `sub` | string | yes | User ID (UUID string) |
| `aud` | string | yes | Intended API audience |
| `exp` | number | yes | Expiration time (Unix seconds) |
| `iat` | number | yes | Issued-at time (Unix seconds) |
| `jti` | string | yes | Unique token ID (UUID string) |
| `sid` | string | yes | Session ID (UUID string) |

Validation requirements:

- Enforce algorithm allowlist (`RS256` only).
- Validate `iss`, `aud`, signature, and time-based claims (`exp`, `iat`) with configured clock skew.
- Reject unknown `kid`.

## 5. Token and Session Lifetime

| Item | Value | Rationale |
|---|---|---|
| Access token TTL | `15m` (`900` seconds) | Limits blast radius for leaked bearer token |
| Refresh token TTL | `30d` (`2592000` seconds) | Supports practical login persistence |
| Session TTL | `30d` | Session and refresh lifecycle remain aligned |
| Clock skew tolerance | `30s` | Handles small host time drift safely |

Phase 1 rule:

- Access tokens are not actively blacklisted. They naturally expire.
- Session revocation blocks refresh; existing access token remains usable until `exp`.

## 6. Refresh Rotation and Reuse Policy

Refresh tokens are one-time use.

### 6.1 State Model

| State | Condition |
|---|---|
| `active` | token exists, not expired, not revoked, not rotated |
| `rotated` | token consumed and replaced by successor |
| `revoked` | token invalid due to logout/reuse/session revoke |
| `expired` | current time > token `expires_at` |

### 6.2 Rotation Flow

1. Client sends `refresh_token`.
2. Server hashes token and loads token record + session.
3. Server validates active session and token state.
4. In one DB transaction:
   - mark current token as rotated/revoked for use,
   - create successor refresh token record,
   - persist link (`replaced_by_token_id`).
5. Issue new access token and return raw successor refresh token once.

### 6.3 Reuse Detection Policy

- If a rotated/revoked refresh token is presented again, treat as token theft/replay.
- Immediately revoke the session and all descendant refresh tokens.
- Return generic `401 REFRESH_TOKEN_INVALID`.
- Emit security/audit event without sensitive token values.

## 7. Error Code Catalog (Phase 1)

| Code | HTTP | Meaning |
|---|---|---|
| `VALIDATION_ERROR` | `400` | Request schema/field validation failed |
| `INVALID_CREDENTIALS` | `401` | Login email/password invalid |
| `ACCESS_TOKEN_INVALID` | `401` | Missing/malformed/tampered token, unknown `kid`, or signature failure |
| `ACCESS_TOKEN_EXPIRED` | `401` | Access token expired |
| `REFRESH_TOKEN_INVALID` | `401` | Unknown/expired/revoked/reused/tampered refresh token |
| `SESSION_REVOKED` | `401` | Session is revoked (internal mapping where applicable) |
| `JWKS_UNAVAILABLE` | `503` | Signing public keys cannot be served |
| `INTERNAL_ERROR` | `500` | Unhandled server error |

Client-facing message policy:

- Keep error messages generic for token failures.
- Use codes for deterministic client behavior.

## 8. Phase 1 Config Keys

| Key | Required | Default | Notes |
|---|---|---|---|
| `ENVIRONMENT` | no | `development` | Existing key |
| `PORT` | no | `8080` | Existing key |
| `DB_URL` | yes | none | Existing key |
| `AUTH_ISSUER` | yes | none | JWT `iss` |
| `AUTH_AUDIENCE` | yes | none | JWT `aud` |
| `ACCESS_TOKEN_TTL` | no | `15m` | Go duration string |
| `REFRESH_TOKEN_TTL` | no | `720h` | Go duration string (`30d`) |
| `SESSION_TTL` | no | `720h` | Must be >= `REFRESH_TOKEN_TTL` |
| `AUTH_CLOCK_SKEW_SECONDS` | no | `30` | Allowed clock drift |
| `JWT_KEY_SOURCE` | no | `env` | `env` or `db` |
| `JWT_ACTIVE_KID` | yes | none | Required in `env` mode; key id in JWT header |
| `JWT_PRIVATE_KEY_PEM` | yes (env mode) | none | RS256 private key PEM |
| `JWT_PUBLIC_KEY_PEM` | yes (env mode) | none | RS256 public key PEM |
| `JWKS_CACHE_MAX_AGE_SECONDS` | no | `300` | Cache-Control max-age |
| `REFRESH_TOKEN_BYTES` | no | `32` | Random bytes before encoding |

## 9. Backward-Compatibility Impact

Current `/login` returns only:

```json
{
  "data": {
    "message": "Login Successful"
  }
}
```

Phase 1 `/login` expands `data` with token/session fields while preserving `message`.

Impact:

- Existing clients that only check HTTP `200` or `data.message` continue to work.
- Clients that enforce strict exact schema must update to allow additional fields.
- New clients must consume `access_token`, `refresh_token`, and TTL fields.

## 10. Contract Test Scenarios (Must Be Automated in Later Tickets)

1. `CT-01 Login Success`: valid credentials return `200` with `access_token`, `refresh_token`, `expires_in`, `session_id`.
2. `CT-02 Login Validation`: malformed email returns `400 VALIDATION_ERROR` with `error.details.email`.
3. `CT-03 Login Unauthorized`: wrong password returns `401 INVALID_CREDENTIALS`.
4. `CT-04 Refresh Rotation Success`: valid refresh token returns new token pair; old token marked rotated.
5. `CT-05 Refresh Reuse Detection` (negative security): reusing old refresh token returns `401 REFRESH_TOKEN_INVALID` and revokes session chain.
6. `CT-06 Refresh Expired Token` (negative security): expired refresh token returns `401 REFRESH_TOKEN_INVALID`.
7. `CT-07 /me Tampered Access Token` (negative security): altered JWT signature returns `401 ACCESS_TOKEN_INVALID`.
8. `CT-08 Logout Idempotency`: repeated logout calls with valid access token return `200` and do not error.
9. `CT-09 JWKS Contract`: `/.well-known/jwks.json` returns RFC-compliant RSA key fields (`kty`, `kid`, `use`, `alg`, `n`, `e`).

## 11. Security Notes

### 11.1 TTL Rationale

- `15m` access TTL minimizes exposure from stolen bearer tokens.
- `30d` refresh TTL provides practical UX while rotation/reuse controls reduce replay risk.
- Short access TTL plus refresh rotation is the Phase 1 compromise between security and usability.

### 11.2 Rotation and Reuse Rationale

- One-time refresh token rotation limits replay window to a single use.
- Reuse detection acts as compromise signal and triggers immediate session containment.
- Generic refresh failure response avoids exposing token-state details to attackers.

### 11.3 Logging Redaction Rules

Never log:

- Raw passwords
- Raw refresh tokens
- Raw access tokens
- Full `Authorization` header
- Private keys (`JWT_PRIVATE_KEY_PEM`)

Allowed in logs:

- `user_id`, `session_id`, `kid`, `jti`
- Error code and HTTP status
- Sanitized request metadata (`ip`, `user_agent`) for audit trails

## 12. Design QA Checklist (Ticket 0001)

- [x] All five Phase 1 endpoints have explicit request/response contracts.
- [x] Success, validation, and unauthorized/revoked/expired behavior are documented.
- [x] Access-token claim schema and `RS256` constraints are frozen.
- [x] Refresh rotation/reuse behavior is explicit and implementation-ready.
- [x] Config keys and defaults are defined with required/optional semantics.
- [x] Backward-compatibility behavior for current `/login` clients is documented.
- [x] Contract test scenarios include at least five total and at least three negative security cases.
