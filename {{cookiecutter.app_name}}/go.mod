module {{cookiecutter.source_path}}/{{cookiecutter.app_name}}

go {{cookiecutter.docker_build_image_version}}

require (
	github.com/bufbuild/buf v0.48.2
	github.com/go-coldbrew/core v0.1.8
	github.com/go-coldbrew/errors v0.1.1
	github.com/go-coldbrew/log v0.1.0
	github.com/golang/protobuf v1.5.2
	github.com/grpc-ecosystem/grpc-gateway/v2 v2.7.2
	github.com/kelseyhightower/envconfig v1.4.0
	github.com/rakyll/statik v0.1.7
	google.golang.org/genproto v0.0.0-20211223182754-3ac035c7e7cb
	google.golang.org/grpc v1.43.0
	google.golang.org/grpc/cmd/protoc-gen-go-grpc v1.1.0
	google.golang.org/protobuf v1.27.1
)
