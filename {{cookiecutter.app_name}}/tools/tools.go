//go:build tools
// +build tools

package tools

// https://github.com/golang/go/wiki/Modules#how-can-i-track-tool-dependencies-for-a-module

import (
	_ "github.com/boumenot/gocover-cobertura"
	_ "github.com/bufbuild/buf/cmd/buf"
	_ "github.com/golangci/golangci-lint/v2/cmd/golangci-lint"
	_ "github.com/vektra/mockery/v2"
)
