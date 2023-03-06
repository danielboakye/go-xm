//go:build wireinject
// +build wireinject

package main

import (
	"context"

	"github.com/danielboakye/go-xm/config"
	"github.com/danielboakye/go-xm/handlers"
	"github.com/danielboakye/go-xm/helpers"
	"github.com/danielboakye/go-xm/pkg/kfkp"
	"github.com/danielboakye/go-xm/pkg/postgres"
	"github.com/danielboakye/go-xm/repo"
	"github.com/google/wire"
)

func buildCompileTime(ctx context.Context) (httpServer HTTPServer, err error) {
	wire.Build(
		config.NewConfigurations,
		helpers.NewValidation,
		postgres.NewConnection,
		kfkp.NewConnection,

		repo.NewRepository,
		handlers.NewHandler,

		newHTTPHandler,
		newHTTPServer,
	)

	return
}
