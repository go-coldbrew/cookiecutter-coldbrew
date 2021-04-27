package service

import (
	"context"
	"fmt"

	{{cookiecutter.app_name|lower}} "{{cookiecutter.source_path}}/{{cookiecutter.app_name}}/proto"
	"{{cookiecutter.source_path}}/{{cookiecutter.app_name}}/config"
	"github.com/golang/protobuf/ptypes/empty"
	"google.golang.org/genproto/googleapis/api/httpbody"
	"github.com/go-coldbrew/errors"
	"github.com/go-coldbrew/log"
)

type svc struct {
	prefix string
}

func (s *svc) HealthCheck(ctx context.Context, _ *empty.Empty) (*httpbody.HttpBody, error) {
	return getHealthCheck(ctx), nil
}

func (s *svc) Echo(_ context.Context, req *{{cookiecutter.app_name|lower}}.EchoRequest) (*{{cookiecutter.app_name|lower}}.EchoResponse, error) {
	return &{{cookiecutter.app_name|lower}}.EchoResponse{
		Msg: fmt.Sprintf("%s: %s", s.prefix, req.GetMsg()),
	}, nil
}

func (s *svc) Error(ctx context.Context, req *{{cookiecutter.app_name|lower}}.EchoRequest) (*{{cookiecutter.app_name|lower}}.EchoResponse, error) {
	err := errors.New("This is an Error")
	log.Info(ctx, "error requested")
	return nil, errors.Wrap(err, "endpoint error")
}

// Creates a new Service
func New(cfg config.Config) ({{cookiecutter.app_name|lower}}.{{cookiecutter.service_name}}Server, error) {
	return &svc{
		prefix: cfg.Prefix,
	}, nil
}
