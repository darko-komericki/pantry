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
