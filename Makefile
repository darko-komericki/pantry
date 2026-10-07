# Load .env (DATABASE_URL etc.) and export it to every recipe.
-include .env
export

.PHONY: help db-up db-down migrate migrate-down migrate-status migrate-new \
	gen gen-sql gen-api gen-ts seed dev dev-backend dev-frontend \
	test test-backend test-frontend lint lint-backend lint-frontend setup

GOOSE = cd backend && go tool goose -dir db/migrations postgres

help: ## List targets
	@grep -E '^[a-z-]+:.*## ' Makefile | awk -F ':.*## ' '{printf "  %-16s %s\n", $$1, $$2}'

setup: ## First-time setup: .env, frontend deps, Go modules
	@test -f .env || cp .env.example .env
	cd frontend && pnpm install
	cd backend && go mod download

# --- Database ---------------------------------------------------------------

db-up: ## Start Postgres (dev :5432, test :5433)
	docker compose up -d --wait

db-down: ## Stop Postgres (dev data is kept in a volume)
	docker compose down

migrate: ## Apply migrations to dev and test DBs
	$(GOOSE) "$(DATABASE_URL)" up
	$(GOOSE) "$(TEST_DATABASE_URL)" up

migrate-down: ## Roll back the last migration on the dev DB
	$(GOOSE) "$(DATABASE_URL)" down

migrate-status: ## Show migration status of the dev DB
	$(GOOSE) "$(DATABASE_URL)" status

migrate-new: ## Create a migration: make migrate-new name=create_users
	@test -n "$(name)" || (echo "usage: make migrate-new name=<name>" && exit 1)
	$(GOOSE) "" create $(name) sql

# --- Code generation --------------------------------------------------------

gen: gen-sql gen-api gen-ts ## Run all code generators

gen-sql: ## sqlc: db/queries -> internal/store (skipped until a query exists)
	@if ls backend/db/queries/*.sql >/dev/null 2>&1; then \
		cd backend && go tool sqlc generate; \
	else echo "gen-sql: no queries yet, skipping"; fi

gen-api: ## oapi-codegen: api/openapi.yaml -> internal/api
	@if [ -f api/openapi.yaml ]; then \
		mkdir -p backend/internal/api && \
		cd backend && go tool oapi-codegen -config oapi-codegen.yaml ../api/openapi.yaml; \
	else echo "gen-api: no api/openapi.yaml yet, skipping"; fi

gen-ts: ## openapi-typescript: api/openapi.yaml -> frontend/src/api/schema.d.ts
	@if [ -f api/openapi.yaml ]; then \
		cd frontend && pnpm --silent gen:api; \
	else echo "gen-ts: no api/openapi.yaml yet, skipping"; fi

seed: ## Load realistic fake data (not implemented yet)
	@echo "seed: nothing to seed until the first tables exist"

# --- Run ----------------------------------------------------------------------

dev: ## Backend with live reload + Vite dev server
	$(MAKE) -j2 dev-backend dev-frontend

dev-backend:
	cd backend && go tool air

dev-frontend:
	cd frontend && pnpm dev

# --- Quality --------------------------------------------------------------------

test: test-backend test-frontend ## Go and frontend tests

test-backend:
	cd backend && go test ./...

test-frontend:
	cd frontend && pnpm test

lint: lint-backend lint-frontend ## gofmt, go vet, tsc, eslint

lint-backend:
	@out=$$(cd backend && gofmt -l .); \
		if [ -n "$$out" ]; then echo "gofmt needed:"; echo "$$out"; exit 1; fi
	cd backend && go vet ./...

lint-frontend:
	cd frontend && pnpm typecheck && pnpm lint
