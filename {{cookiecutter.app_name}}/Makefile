.PHONY: build build-alpine clean test default deps generate run run-docker fmt help bench build-docker coverage-html dep lint mock runj vulncheck local-stack local-stack-down local-stack-logs local-stack-reset local-psql loadtest

BIN_NAME={{cookiecutter.app_name}}
GOPRIVATE ?= {{cookiecutter.goprivate}}
export GOPRIVATE

SHELL := $(shell which bash)
VERSION := $(shell grep "const Version " version/version.go | sed -E 's/.*"(.+)"$$/\1/')
GIT_COMMIT=$(shell git rev-parse HEAD)
GIT_DIRTY=$(shell test -n "`git status --porcelain`" && echo "+DIRTY" || true)
GIT_BRANCH=$(shell git rev-parse --abbrev-ref HEAD)
BUILD_DATE=$(shell date '+%Y-%m-%d-%H:%M:%S')
IMAGE_NAME := "{{cookiecutter.app_name|lower}}"

default: test

help:
	@echo 'Management commands for {{cookiecutter.app_name}}:'
	@echo
	@echo 'Usage:'
	@echo '    make bench           Run benchmarks.'
	@echo '    make build           Compile the project and generate a binary.'
	@echo '    make build-docker    Build a docker image.'
	@echo '    make clean           Clean the directory tree.'
	@echo '    make coverage-html   Generate test coverage report.'
	@echo '    make dep             Update dependencies.'
	@echo '    make deps            Download Go modules.'
	@echo '    make fmt             Format Go source files.'
	@echo '    make generate        Generate code from proto files.'
	@echo '    make help            Show this message.'
	@echo '    make lint            Run linters and vulnerability check.'
	@echo '    make vulncheck       Run Go vulnerability check.'
	@echo '    make mock            Generate mocks for interfaces.'
	@echo '    make run             Run the project locally.'
	@echo '    make run-docker      Run the project in a docker container.'
	@echo '    make runj            Run the project locally with jq log parsing.'
	@echo '    make test            Run tests.'
	@echo '    make loadtest         Run gRPC load test (ghz) against the running service.'
	@echo '    make local-stack     Start local dev stack (default: deps; PROFILES="deps obs" for full).'
	@echo '    make local-stack-down  Stop local dev stack.'
	@echo '    make local-stack-logs  Follow local dev stack logs.'
	@echo '    make local-stack-reset Reset local dev stack (down + up).'
	@echo '    make local-psql      Open psql shell in local Postgres.'
	@echo

build:
	@echo "building ${BIN_NAME} ${VERSION}"
	@echo "GOPATH=${GOPATH}"
	go build -ldflags "-X {{cookiecutter.source_path}}/{{cookiecutter.app_name}}/version.GitCommit=${GIT_COMMIT}${GIT_DIRTY} -X {{cookiecutter.source_path}}/{{cookiecutter.app_name}}/version.BuildDate=${BUILD_DATE} -X {{cookiecutter.source_path}}/{{cookiecutter.app_name}}/version.Branch=${GIT_BRANCH}" -o bin/${BIN_NAME}

build-alpine:
	@echo "building ${BIN_NAME} ${VERSION}"
	@echo "GOPATH=${GOPATH}"
	CGO_ENABLED=0 go build -ldflags '-w -X {{cookiecutter.source_path}}/{{cookiecutter.app_name}}/version.GitCommit=${GIT_COMMIT}${GIT_DIRTY} -X {{cookiecutter.source_path}}/{{cookiecutter.app_name}}/version.BuildDate=${BUILD_DATE} -X {{cookiecutter.source_path}}/{{cookiecutter.app_name}}/version.Branch=${GIT_BRANCH}' -o bin/${BIN_NAME}

build-docker:
	@echo "building image ${BIN_NAME} ${VERSION} $(GIT_COMMIT)"
	docker build --build-arg VERSION=${VERSION} --build-arg GIT_COMMIT=$(GIT_COMMIT) --build-arg GIT_BRANCH=$(GIT_BRANCH) --build-arg GOPRIVATE=$(GOPRIVATE) -t $(IMAGE_NAME):local .

dep:
	go mod tidy

deps:
	go mod download

generate: deps
	go tool buf generate proto

clean:
	@test ! -e bin/${BIN_NAME} || rm bin/${BIN_NAME}
	go clean ./...

test:
	go test -race -coverpkg=.,./config/...,./service/... -coverprofile cover.out ./...
	go tool cover -func=cover.out

coverage-html:
	go tool cover -html=cover.out -o=cover.html

bench:
	go test -bench=. -run=^$$ -benchtime 10s -benchmem ./...

lint: deps vulncheck
	go tool golangci-lint run --timeout 5m

vulncheck: deps
	go tool govulncheck ./...

mock: deps
	go tool mockery --config .mockery.yaml

fmt:
	gofmt -w .

run: build
	@echo
	@echo "swagger ui available at http://localhost:9091/swagger/"
	@echo
	@set -a; source local.env; ./bin/${BIN_NAME}

runj: build
	@echo
	@echo "swagger ui available at http://localhost:9091/swagger/"
	@echo
	@set -a; source local.env && ./bin/${BIN_NAME} 1> >(jq -R "fromjson? | ." -C)

run-docker: build-docker
	docker run -p 9091:9091 -p 9090:9090 --env-file local.env ${IMAGE_NAME}:local

# Local development stack (docker-compose)
# Usage:
#   make local-stack                        # deps profile (postgres, redis, adminer)
#   make local-stack PROFILES="deps obs"    # + prometheus, grafana, jaeger
#   make local-stack-down                   # tear down
#   DOCKER_COMPOSE="docker compose" make local-stack   # use Compose v2 plugin

DOCKER_COMPOSE ?= docker-compose
COMPOSE_FILE := docker-compose.local.yml
PROFILES ?= deps
PROFILE_FLAGS := $(foreach p,$(PROFILES),--profile $(p))

local-stack:
	$(DOCKER_COMPOSE) -f $(COMPOSE_FILE) $(PROFILE_FLAGS) up -d --wait --remove-orphans
	@echo
	@echo "Local stack is running. Run 'make run' to start the app."
	@echo "Endpoints:"
	@$(DOCKER_COMPOSE) -f $(COMPOSE_FILE) $(PROFILE_FLAGS) port db 5432 2>/dev/null | sed 's|.*:\(.*\)|  Postgres:    localhost:\1|' || true
	@$(DOCKER_COMPOSE) -f $(COMPOSE_FILE) $(PROFILE_FLAGS) port redis 6379 2>/dev/null | sed 's|.*:\(.*\)|  Redis:       localhost:\1|' || true
	@$(DOCKER_COMPOSE) -f $(COMPOSE_FILE) $(PROFILE_FLAGS) port adminer 8080 2>/dev/null | sed 's|.*:\(.*\)|  Adminer:     http://localhost:\1|' || true
	@$(DOCKER_COMPOSE) -f $(COMPOSE_FILE) $(PROFILE_FLAGS) port prometheus 9090 2>/dev/null | sed 's|.*:\(.*\)|  Prometheus:  http://localhost:\1|' || true
	@$(DOCKER_COMPOSE) -f $(COMPOSE_FILE) $(PROFILE_FLAGS) port grafana 3000 2>/dev/null | sed 's|.*:\(.*\)|  Grafana:     http://localhost:\1  (admin/admin)|' || true
	@$(DOCKER_COMPOSE) -f $(COMPOSE_FILE) $(PROFILE_FLAGS) port jaeger 16686 2>/dev/null | sed 's|.*:\(.*\)|  Jaeger:      http://localhost:\1|' || true
	@echo

local-stack-down:
	$(DOCKER_COMPOSE) -f $(COMPOSE_FILE) $(PROFILE_FLAGS) down

local-stack-logs:
	$(DOCKER_COMPOSE) -f $(COMPOSE_FILE) $(PROFILE_FLAGS) logs -f

local-stack-reset: local-stack-down local-stack

local-psql:
	$(DOCKER_COMPOSE) -f $(COMPOSE_FILE) --profile deps exec db sh -c 'PGPASSWORD=postgres psql -U postgres -d {{cookiecutter.app_name}}_dev'

loadtest:
	ghz --config misc/loadtest/echo.json
