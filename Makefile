# Makefile

# Load environment variables from .env file
ifneq (,$(wildcard ./.env))
    include .env
    export
endif

# Path to migration files
MIGRATIONS_PATH = migrations

# Database URL
DATABASE_URL = postgres://$(POSTGRES_USER):$(POSTGRES_PASSWORD)@$(POSTGRES_HOST):$(POSTGRES_PORT)/$(POSTGRES_DB)?sslmode=disable

# Test Database URL
DATABASE_URL_TEST = postgres://$(POSTGRES_USER):$(POSTGRES_PASSWORD)@$(POSTGRES_HOST):$(POSTGRES_PORT)/$(POSTGRES_DB_TEST)?sslmode=disable

# Install migrate tool
install-migrate:
	go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest

# Create a new migration
create-migration:
ifeq ($(OS),Windows_NT)
	@powershell -Command "if ($$name -eq $$null) {$$name = Read-Host 'Enter migration name'}; migrate create -ext sql -dir $(MIGRATIONS_PATH) -seq $$name"
else
	@read -p "Enter migration name: " name; migrate create -ext sql -dir $(MIGRATIONS_PATH) -seq $$name
endif

# Force migrate to a specific version
migrate-fix:
ifeq ($(OS),Windows_NT)
	@powershell -Command "if ($$version -eq $$null) {$$version = Read-Host 'Enter migration version to force to'}; migrate -path $(MIGRATIONS_PATH) -database \"$(DATABASE_URL)\" force $$version"
else
	@read -p "Enter migration version to force to: " version; migrate -path $(MIGRATIONS_PATH) -database "$(DATABASE_URL)" force $$version
endif

# Apply migrations
migrate-up:
	migrate -path $(MIGRATIONS_PATH) -database "$(DATABASE_URL)" up

# Rollback migrations
migrate-down:
	migrate -path $(MIGRATIONS_PATH) -database "$(DATABASE_URL)" down

# Apply test migrations
migrate-up-test:
	migrate -path $(MIGRATIONS_PATH) -database "$(DATABASE_URL_TEST)" up

# Rollback test migrations
migrate-down-test:
	migrate -path $(MIGRATIONS_PATH) -database "$(DATABASE_URL_TEST)" down

restart:
	docker-compose down
	docker-compose up -d --build
	docker image prune -f

minio-access:
	mc alias set myminio http://localhost:9000 $(MINIO_ACCESS_KEY) $(MINIO_SECRET_KEY)
	mc anonymous set download myminio/be-template

.PHONY: install-migrate create-migration migrate-up migrate-down migrate-up-test migrate-down-test migrate-fix restart minio-access
