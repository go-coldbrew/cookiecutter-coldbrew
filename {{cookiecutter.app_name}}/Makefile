.PHONY: build build-alpine clean test default install generate run run-docker

BIN_NAME={{cookiecutter.app_name}}

VERSION := $(shell grep "const Version " version/version.go | sed -E 's/.*"(.+)"$$/\1/')
GIT_COMMIT=$(shell git rev-parse HEAD)
GIT_DIRTY=$(shell test -n "`git status --porcelain`" && echo "+DIRTY" || true)
BUILD_DATE=$(shell date '+%Y-%m-%d-%H:%M:%S')
IMAGE_NAME := "{{cookiecutter.app_name|lower}}"

default: test

help:
	@echo 'Management commands for {{cookiecutter.app_name}}:'
	@echo
	@echo 'Usage:'
	@echo '    make build           Compile the project.'
	@echo '    make test            Run tests on a compiled project.'
	@echo '    make clean           Clean the directory tree.'
	@echo '    make dep             Update dependencies.'
	@echo '    make generate        Generate code.'
	@echo '    make run             Run the project.'
	@echo '    make run-docker      Run the project in a docker container.'
	@echo '    make build-docker    Build a docker image.'
	@echo '    make lint            Run linters.'
	@echo '    make mock            Generate mocks.'
	@echo

build:
	@echo "building ${BIN_NAME} ${VERSION}"
	@echo "GOPATH=${GOPATH}"
	go build -ldflags "-X {{cookiecutter.source_path}}/{{cookiecutter.app_name}}/version.GitCommit=${GIT_COMMIT}${GIT_DIRTY} -X {{cookiecutter.source_path}}/{{cookiecutter.app_name}}/version.BuildDate=${BUILD_DATE} -X {{cookiecutter.source_path}}/{{cookiecutter.app_name}}/version.Branch=${GIT_BRANCH}" -o bin/${BIN_NAME}

build-alpine:
	@echo "building ${BIN_NAME} ${VERSION}"
	@echo "GOPATH=${GOPATH}"
	go build -ldflags '-w -linkmode external -extldflags "-static" -X {{cookiecutter.source_path}}/{{cookiecutter.app_name}}/version.GitCommit=${GIT_COMMIT}${GIT_DIRTY} -X {{cookiecutter.source_path}}/{{cookiecutter.app_name}}/version.BuildDate=${BUILD_DATE} -X {{cookiecutter.source_path}}/{{cookiecutter.app_name}}/version.Branch=${GIT_BRANCH} ' -o bin/${BIN_NAME}

build-docker:
	@echo "building image ${BIN_NAME} ${VERSION} $(GIT_COMMIT)"
	docker build --build-arg VERSION=${VERSION} --build-arg GIT_COMMIT=$(GIT_COMMIT) --build-arg GIT_BRANCH=$(GIT_BRANCH) -t $(IMAGE_NAME):local .

deps:
	go mod tidy

install:
	go install \
		google.golang.org/protobuf/cmd/protoc-gen-go \
		google.golang.org/grpc/cmd/protoc-gen-go-grpc \
		github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway \
		github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2 \
		github.com/rakyll/statik \
		github.com/vektra/mockery/v2 \
		github.com/bufbuild/buf/cmd/buf \
		github.com/golangci/golangci-lint/cmd/golangci-lint

generate: install
	buf generate --path proto/*.proto
	# Generate static assets for OpenAPI UI
	statik -m -f -src third_party/OpenAPI/

clean:
	@test ! -e bin/${BIN_NAME} || rm bin/${BIN_NAME}
	go clean ./...

test:
	go test -race -coverpkg=.,./config/...,./service/... -coverprofile cover.out ./...
	go tool cover -func=cover.out

coverage-html:
	go tool cover -html=cover.out -o=cover.html

bench:
	# -run=^B negates all tests
	go test -bench=. -run=^B -benchtime 10s -benchmem ./...

lint: install
	golangci-lint run

mock: install
	mockery --all --keeptree --dir ./service/ --output ./misc/mocks

run: build
	set -a; source local.env; ./bin/${BIN_NAME}

run-docker: build-docker
	docker run -p 9091:9091 -p 9090:9090 --env-file local.env ${IMAGE_NAME}:local
