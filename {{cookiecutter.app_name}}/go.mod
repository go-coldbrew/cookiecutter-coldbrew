module {{cookiecutter.source_path}}/{{cookiecutter.app_name}}

go 1.15

require (
	github.com/bufbuild/buf v0.41.0
	github.com/go-coldbrew/core v0.1.0
	github.com/go-coldbrew/errors v0.1.0
	github.com/go-coldbrew/log v0.1.0
	github.com/golang/protobuf v1.5.2
	github.com/grpc-ecosystem/grpc-gateway/v2 v2.3.0
	github.com/kelseyhightower/envconfig v1.4.0
	github.com/rakyll/statik v0.1.7
	google.golang.org/genproto v0.0.0-20210426193834-eac7f76ac494
	google.golang.org/grpc v1.37.0
	google.golang.org/grpc/cmd/protoc-gen-go-grpc v1.1.0
	google.golang.org/protobuf v1.26.0
)
