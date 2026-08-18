BASE_STACK = docker compose -f ../deploy/docker-compose.yml
INTEGRATION_TEST_STACK = $(BASE_STACK) -f ../deploy/docker-compose-integration-test.yml
ALL_STACK = $(INTEGRATION_TEST_STACK)

# HELP =================================================================================================================
# This will output the help for each task
# thanks to https://marmelab.com/blog/2016/02/29/auto-documented-makefile.html
.PHONY: help

help: ## Display this help screen
	@awk 'BEGIN {FS = ":.*##"; printf "\nUsage:\n  make \033[36m<target>\033[0m\n"} /^[a-zA-Z0-9_-]+:.*?##/ { printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2 } /^##@/ { printf "\n\033[1m%s\033[0m\n", substr($$0, 5) } ' $(MAKEFILE_LIST)

compose-up: ### Run docker compose (without backend and reverse proxy)
	$(BASE_STACK) up --build -d db && docker compose logs -f
.PHONY: compose-up

compose-up-all: ### Run docker compose (with backend and reverse proxy)
	$(BASE_STACK) up --build -d
.PHONY: compose-up-all

compose-up-integration-test: ### Run docker compose with integration test
	$(INTEGRATION_TEST_STACK) up --build --abort-on-container-exit --exit-code-from integration-test
.PHONY: compose-up-integration-test

compose-down: ### Down docker compose
	$(ALL_STACK) down --remove-orphans
.PHONY: compose-down

swag-v1: ### swag init
	swag init -g internal/router/router.go
.PHONY: swag-v1

proto-v1: ### generate source files from proto
	protoc --go_out=. \
		--go_opt=paths=source_relative \
		--go-grpc_out=. \
		--go-grpc_opt=paths=source_relative \
		docs/proto/v1/*.proto
.PHONY: proto-v1

deps: ### deps tidy + verify
	go mod tidy && go mod verify
.PHONY: deps

deps-audit: ### check dependencies vulnerabilities
	govulncheck ./...
.PHONY: deps-audit

format: ### Run code formatter
	gofumpt -l -w .
	gci write . --skip-generated -s standard -s default
.PHONY: format

run: deps swag-v1 ### swag run for API v1
	go mod download && \
	CGO_ENABLED=0 go run ./cmd/app
.PHONY: run

run-with-migrate: deps swag-v1 ### swag run for API v1
	go mod download && \
	CGO_ENABLED=0 go run -tags migrate ./cmd/app
.PHONY: run-with-migrate

docker-rm-volume: ### remove docker volume
	docker volume rm gin-app_pg-data
.PHONY: docker-rm-volume

linter-golangci: ### check by golangci linter
	golangci-lint run
.PHONY: linter-golangci

test: ### run test
	go test -v -race -covermode atomic -coverprofile=coverage.txt ./internal/...
.PHONY: test

test-integration: ### run Referral integration tests against an isolated PostgreSQL container
	docker compose -f ../deploy/docker-compose.test.yml up -d --wait
	trap 'docker compose -f ../deploy/docker-compose.test.yml down -v' EXIT; \
	REFERRAL_TEST_DATABASE_URL='postgres://referral_test:referral_test@localhost:55432/referral_test?sslmode=disable' \
	go test -tags=integration ./internal/modules/referral/client
.PHONY: test-integration

k6-smoke: ### run authenticated backend smoke/load checks (requires ADMIN_EMAIL and ADMIN_PASSWORD)
	k6 run test/k6/backend-smoke.js
.PHONY: k6-smoke

mock: ### run mockgen
	mockgen -source ./internal/modules/user/user.repository.go -package user_test > ./internal/modules/user/mocks_user_repo_test.go
.PHONY: mock

schema-create: ### create ent schema
	go run -mod=mod entgo.io/ent/cmd/ent new '$(word 2,$(MAKECMDGOALS))'
.PHONY: schema-create

ent-gen: ### generate ent code
	go generate ./ent
.PHONY: ent-gen

migrate-create:  ### create new migration
	go run -mod=mod ent/migrate/main.go '$(word 2,$(MAKECMDGOALS))'
.PHONY: migrate-create

migrate-up: ### migration up
	go run ./cmd/migrate up
.PHONY: migrate-up

migrate-down: ### migration down
	go run ./cmd/migrate down
.PHONY: migrate-down

migrate-hash: ### Recalculate the Atlas migration integrity file
	go run ./cmd/migrationhash
.PHONY: migrate-hash

ls-deps: ### list outdated dependencies
	go list -m -u all
.PHONY: ls-deps

upgrade-deps: ### upgrade dependencies to latest minor/patch versions
	go get -u ./...
	go mod tidy
.PHONY: upgrade-deps

bin-deps: ### install tools
	go install tool
	go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate
.PHONY: bin-deps

pre-commit: swag-v1 format linter-golangci ### run pre-commit
.PHONY: pre-commit
