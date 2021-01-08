module {{cookiecutter.source_path}}/{{cookiecutter.app_name}}

require (

	github.com/bufbuild/buf v0.28.0
	github.com/go-coldbrew/core v0.0.0-20210108141338-e0e9bbc8553d
	github.com/golang/protobuf v1.4.3
	github.com/grpc-ecosystem/grpc-gateway/v2 v2.0.1
	github.com/kelseyhightower/envconfig v1.4.0
	github.com/rakyll/statik v0.1.7
	github.com/sirupsen/logrus v1.4.1
	github.com/spf13/viper v1.7.0
	google.golang.org/genproto v0.0.0-20201021134325-0d71844de594
	google.golang.org/grpc v1.33.1
	google.golang.org/grpc/cmd/protoc-gen-go-grpc v1.0.1
	google.golang.org/protobuf v1.25.0
)

//replace github.com/go-coldbrew/core => /Users/ankurshrivastava/code/ColdBrew/core
