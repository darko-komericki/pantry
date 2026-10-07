# Decisions

Tradeoff decisions, newest last. Each entry: date, decision, why, alternatives considered.

## 2026-10-07: Go tools pinned via `tool` directive in go.mod

sqlc, goose, oapi-codegen and air are added with `go get -tool` and run with `go tool <name>`.
Versions are pinned per project in `go.mod`, nothing is installed globally.
Alternative: `brew install` (simpler, but unpinned and machine-dependent).

## 2026-10-07: Typed API client with openapi-fetch

Frontend calls go through `openapi-fetch`, typed by `openapi-typescript` output.
Alternative: hand-written `fetch` wrapper (no dependency, but path/param/response typing is
hard to get right by hand).

## 2026-10-07: pnpm for the frontend

Alternative: npm. pnpm is faster and stricter about undeclared dependencies.

## 2026-10-07: Separate Postgres container for tests

`db-test` on port 5433 with tmpfs storage. Store tests never touch dev data and the instance
resets on restart. Alternative: second database in the same container (less isolation).

## 2026-10-07: Server-side sessions in Postgres

Login creates a random token. The browser gets it in an `HttpOnly`, `SameSite=Lax` cookie
(`Secure` in production). The `sessions` table stores only a SHA-256 hash of the token.
CSRF protection uses Go's `http.CrossOriginProtection` (standard library).
Why: same-origin SPA, real logout and "log out everywhere", token not readable by JavaScript.
Alternatives: JWT (cannot revoke before expiry without rebuilding sessions), external
provider (extra service, hides the part we want to learn).

## 2026-10-07: bcrypt for password hashing

`golang.org/x/crypto/bcrypt`, cost 12. Passwords limited to 72 bytes (bcrypt limit),
enforced in the contract.
Alternative: argon2id (stronger against GPU attacks, but we would hand-write the
salt/parameter encoding).

## 2026-10-07: Users join households through a membership table

`household_members (household_id, user_id)`. The data model allows many households per user;
the UI assumes one for now.
Alternative: `users.household_id` (simpler, but a data migration if it ever changes).

## 2026-10-07: Invite codes for joining a household

The owner creates a single-use code; a new user enters it at sign-up.
Alternatives: owner creates the account (owner knows the password), email invite (needs email
sending).

## 2026-10-07: UUIDv7 primary keys

All tables use `uuid` keys with Postgres 18 `uuidv7()` as the default.
Why: time-ordered (indexes nearly as well as bigint), and clients can create IDs offline,
which makes retries in the Phase 2 offline queue safe.
Alternative: `bigint` identity (8 bytes, easier to read, but needs a separate client ID later).

## Known gaps in Phase 1

Deferred on purpose: email verification, password reset, login rate limiting.

## 2026-10-07: `updated_at` maintained by a database trigger

One shared function `set_updated_at()`, attached to each table with a `BEFORE UPDATE`
trigger. No query can forget to set it.
Alternative: set `updated_at = now()` in every UPDATE query (visible in the SQL, but easy to
forget). Note: `now()` is the transaction start time, so all rows changed in one transaction
get the same value.

## 2026-10-07: Expired sessions removed by a cleanup job

Expired rows stay until a periodic delete (added with the session handlers). No index on
`expires_at` until that job exists and its query plan shows a need.
