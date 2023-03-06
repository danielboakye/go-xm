// Code generated by Wire. DO NOT EDIT.

//go:generate go run github.com/google/wire/cmd/wire
//go:build !wireinject
// +build !wireinject

package main

import (
	"context"
	"github.com/danielboakye/go-xm/config"
	"github.com/danielboakye/go-xm/handlers"
	"github.com/danielboakye/go-xm/helpers"
	"github.com/danielboakye/go-xm/pkg/kfkp"
	"github.com/danielboakye/go-xm/pkg/postgres"
	"github.com/danielboakye/go-xm/repo"
)

import (
	_ "github.com/jackc/pgx/v4/stdlib"
)

// Injectors from wire.go:

func BuildCompileTime(ctx context.Context) (HTTPServer, error) {
	configurations, err := config.NewConfigurations()
	if err != nil {
		return HTTPServer{}, err
	}
	db, err := postgres.NewConnection(configurations)
	if err != nil {
		return HTTPServer{}, err
	}
	iKafkaConn, err := kfkp.NewConnection(ctx, configurations)
	if err != nil {
		return HTTPServer{}, err
	}
	iRepository := repo.NewRepository(db, iKafkaConn)
	validation, err := helpers.NewValidation()
	if err != nil {
		return HTTPServer{}, err
	}
	handler := handlers.NewHandler(iRepository, validation, configurations)
	ihttpHandler := newHTTPHandler(handler, configurations)
	httpServer := newHTTPServer(ihttpHandler, configurations)
	return httpServer, nil
}
