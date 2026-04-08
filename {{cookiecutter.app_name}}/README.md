# {{cookiecutter.app_name}}

{{cookiecutter.project_short_description}}

## Getting started

This project requires Go to be installed.

On OS X with Homebrew you can just run `brew install go`.

On Linux use your package manager to install go or use the [go documentation](https://go.dev/doc/install)

Running it then should be as simple as:

```console
$ make run
```

## Makefile

The Makefile contains a number of useful commands to help you get started. Here are some of the most useful ones:
- `make help` - Prints the help
- `make test` - Runs the tests
- `make bench` - Runs the benchmarks
- `make lint` - Runs the linter
- `make run` - Runs the application
- `make runj` - Runs the application with json logs parsing with jq
- `make loadtest` - Runs gRPC load test ([ghz](https://ghz.sh)) against the running service
- `make build` - Builds the application
- `make generate` - Generates the code

## Local Development Stack

Start infrastructure with docker-compose, then run the app natively:

```console
$ make local-stack                                # Start default services
$ make local-stack PROFILES="postgres kafka obs"  # Override with specific services
$ make run                                        # Run the app (fast native build)
```

Available profiles: `postgres`, `mysql`, `cockroachdb`, `mongodb`, `redis`, `valkey`, `memcached`, `kafka`, `nats`, `elasticsearch`, `ministack`, `dynamodb`, `spanner`, `pubsub`, `bigtable`, `firestore`, `adminer`, `obs`

```console
$ make local-stack-down    # Stop infrastructure
$ make local-stack-reset   # Reset infrastructure
$ make local-exec SVC=postgres CMD="psql -U postgres"  # Exec into any service
```

## Docker

This project also contains a Dockerfile to help you get started with Docker. To build the image, run:

```console
$ make build-docker
```

To run the image, run:

```console
$ make run-docker
```

## CI/CD

The project includes ready-to-use CI pipelines for both GitHub and GitLab. Delete whichever platform you don't use.

### GitHub Actions (`.github/workflows/go.yml`)

Runs on push to `main`/`master` and on pull requests. Four parallel jobs: **build**, **test** (race detector + coverage), **benchmark**, and **lint** (govulncheck + golangci-lint). Each job has concurrency control to cancel duplicate runs.

### GitLab CI (`.gitlab-ci.yml`)

Three jobs in a single `test` stage: **unit-test** (with Cobertura coverage report), **lint** (golangci-lint + govulncheck), and **benchmark**. Go module caching is enabled.

## Adding a new endpoint to the API

Our service is grpc first. We use [grpc-gateway] to automatically map HTTP requests to gRPC requests. This means that you can add a new endpoint to the API by adding a new rpc to `service {{cookiecutter.service_name}}` in `proto/{{cookiecutter.app_name|lower}}.proto` file. Then, you can run `make generate` to generate grpc/http endpoints.

The file `service/service.go` contains the implementation of the API and serves as the entrypoint for the app. You can add your business logic there or any other package.

### HTTP to gRPC mapping

We use [grpc-gateway] to automatically map HTTP requests to gRPC requests. You can find the mapping in the `proto/{{cookiecutter.app_name|lower}}.proto` file. This server is generated according to [custom options](https://cloud.google.com/service-infrastructure/docs/service-management/reference/rpc/google.api#http) in your gRPC definition.  You can find more information about the mapping [here](https://grpc-ecosystem.github.io/grpc-gateway/docs/tutorials/adding_annotations/)

## Application configuration

This project uses [envconfig] to manage configuration as environment variables. You can find the configuration struct in `config/config.go`. You can also find the default values in the `config/config.go` file.

### Environment variables for local development

You can find the environment variables for local development in the `local.env` file (copy from `local.env.example`). This file is used when you run `make run` or `make run-docker`.

### Coldbrew configuration options

A large number of configuration options are powered by [Coldbrew] and used as environment variables. You can find the list of environment variables [here](https://pkg.go.dev/github.com/go-coldbrew/core/config#Config).

## Logging

This project uses `go-coldbrew/log` to manage logging. You can find documentation [here](https://pkg.go.dev/github.com/go-coldbrew/log).

---
[envconfig]: https://github.com/kelseyhightower/envconfig
[grpc-gateway]: https://grpc-ecosystem.github.io/grpc-gateway/
[Coldbrew]: https://docs.coldbrew.cloud
