// Code generated by Wire. DO NOT EDIT.

//go:generate go run -mod=mod github.com/google/wire/cmd/wire
//go:build !wireinject
// +build !wireinject

package main

// Injectors from wire.go:

func InitializeApp(cfg *Config) (*App, error) {
	emailService := ProvideEmailService(cfg)
	weChatService := ProvideWeChatService(cfg)
	messageService := ProvideMessageService(cfg, emailService, weChatService)
	app := NewApp(messageService)
	return app, nil
}
