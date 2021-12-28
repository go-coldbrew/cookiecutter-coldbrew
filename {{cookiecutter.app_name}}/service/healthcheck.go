package service

import (
	"context"
	"encoding/json"
	"sync/atomic"

	"google.golang.org/genproto/googleapis/api/httpbody"
	"{{cookiecutter.source_path}}/{{cookiecutter.app_name}}/version"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var (
	hc    []byte
	ready int32
)

const (
	statusReady    = 1
	statusNotReady = 0
)

func init() {
	hc, _ = json.Marshal(version.Get())
}

func GetHealthCheck(context.Context) *httpbody.HttpBody {
	return &httpbody.HttpBody{
		ContentType: "application/json",
		Data:        hc,
	}
}

func GetReadyState(ctx context.Context) (*httpbody.HttpBody, error) {
	st := atomic.LoadInt32(&ready)
	if st == statusReady {
		return GetHealthCheck(ctx), nil
	}
	return nil, status.Error(codes.Internal, "Not Ready to server traffic")
}

func SetNotReady() {
	atomic.StoreInt32(&ready, statusNotReady)
}

func SetReady() {
	atomic.StoreInt32(&ready, statusReady)
}
