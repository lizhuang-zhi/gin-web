//go:build wireinject
// +build wireinject

package main

import (
	"github.com/google/wire"
)

func InitializeApp(messageService MessageService) (*App, error) {
	wire.Build(NewApp)
	return &App{}, nil
}
