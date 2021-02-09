package service

import (
	"context"
	"encoding/json"

	"google.golang.org/genproto/googleapis/api/httpbody"
	"{{cookiecutter.source_path}}/{{cookiecutter.app_name}}/version"
)

var hc []byte

func init() {
	hc, _ = json.Marshal(version.Get())
}

func getHealthCheck(context.Context) *httpbody.HttpBody {
	return &httpbody.HttpBody{
		ContentType: "application/json",
		Data:        hc,
	}
}
