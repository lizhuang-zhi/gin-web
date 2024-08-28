//go:build wireinject
// +build wireinject

package main

import (
	"github.com/google/wire"
)

func InitializeApp(cfg *Config) (*App, error) {
	wire.Build(ProvideEmailService, ProvideWeChatService, ProvideMessageService, NewApp)
	return &App{}, nil
}
