.PHONY: build build-alpine clean test default install generate run run-docker fmt doc

BIN_NAME={{cookiecutter.app_name}}

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
	@echo '    make doc             Generate documentation.'
	@echo '    make fmt             Format Go source files.'
	@echo '    make generate        Generate code from proto files.'
	@echo '    make help            Show this message.'
	@echo '    make lint            Run linters on the project.'
	@echo '    make mock            Generate mocks for interfaces.'
	@echo '    make run             Run the project locally.'
	@echo '    make run-docker      Run the project in a docker container.'
	@echo '    make runj            Run the project locally with jq log parsing.'
	@echo '    make test            Run tests.'
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
	docker build --build-arg VERSION=${VERSION} --build-arg GIT_COMMIT=$(GIT_COMMIT) --build-arg GIT_BRANCH=$(GIT_BRANCH) -t $(IMAGE_NAME):local .

dep:
	go mod tidy

install:
	go mod download

generate: install
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

lint: install
	go tool golangci-lint run --timeout 5m

mock: install
	go tool mockery --config .mockery.yaml

fmt:
	gofmt -w .

doc:
	go run github.com/princjef/gomarkdoc/cmd/gomarkdoc@latest ./... > README.md

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
