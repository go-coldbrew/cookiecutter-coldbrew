module {{cookiecutter.source_path}}/{{cookiecutter.app_name}}

go 1.15

require (
	github.com/bufbuild/buf v0.38.0
	github.com/go-coldbrew/core v0.0.0-20210204143115-a9cf4658616b
	github.com/go-coldbrew/errors v0.0.0-20210202032732-032991513cc0
	github.com/go-coldbrew/log v0.0.0-20210108160031-027cd625aad3
	github.com/golang/protobuf v1.4.3
	github.com/grpc-ecosystem/grpc-gateway/v2 v2.3.0
	github.com/kelseyhightower/envconfig v1.4.0
	github.com/rakyll/statik v0.1.7
	google.golang.org/genproto v0.0.0-20210226172003-ab064af71705
	google.golang.org/grpc v1.36.0
	google.golang.org/grpc/cmd/protoc-gen-go-grpc v1.1.0
	google.golang.org/protobuf v1.25.1-0.20201208041424-160c7477e0e8
)
