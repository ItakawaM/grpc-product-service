include .env

export

MIGRATIONS_PATH := ./migrations

POSTGRES_URL := postgres://$(POSTGRES_USERNAME):$(POSTGRES_PASSWORD)@$(POSTGRES_HOST):$(POSTGRES_PORT)/$(POSTGRES_DATABASE)?sslmode=disable

.DEFAULT_GOAL := help

.PHONY: \
	help \
	db/migrate/up \
	db/migrate/down

db/migrate/up:
	migrate -path=$(MIGRATIONS_PATH) -database="$(POSTGRES_URL)" up

db/migrate/down:
	migrate -path=$(MIGRATIONS_PATH) -database="$(POSTGRES_URL)" down

help:
	@echo "Commands:"
	@echo "  make db/migrate/up          - Applies database migrations"
	@echo "  make db/migrate/down        - Rolls back database migrations"