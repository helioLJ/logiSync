DB_URL ?= postgres://postgres:postgres@localhost:5432/logisync?sslmode=disable

.PHONY: api worker mock-portal migrate test test-coverage test-coverage-html playwright-install

api:
	go run ./cmd/api

worker:
	go run ./cmd/worker

mock-portal:
	go run ./cmd/mock-portal

playwright-install:
	go run github.com/playwright-community/playwright-go/cmd/playwright install

migrate:
	@for file in internal/db/migrations/*.sql; do \
		printf "applying %s\n" $$file; \
		psql "$(DB_URL)" -v ON_ERROR_STOP=1 -f $$file; \
	done

test:
	go test ./...

test-coverage:
	go test -cover ./...

test-coverage-html:
	go test -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"
