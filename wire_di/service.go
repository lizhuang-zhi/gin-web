package main

import "fmt"

type MessageService interface {
	Send(message string, receiver string) error
}

type EmailService struct {
	username string
	password string
}

func NewEmailService(username string, password string) *EmailService {
	return &EmailService{
		username: username,
		password: password,
	}
}

func (e *EmailService) Send(message string, receiver string) error {
	fmt.Println("Email sent to ", receiver, " with message ", message)
	return nil
}

type WeChatService struct {
	accountName     string
	accountPassword string
}

func NewWeChatService(accountName string, accountPassword string) *WeChatService {
	return &WeChatService{
		accountName:     accountName,
		accountPassword: accountPassword,
	}
}

func (w *WeChatService) Send(message string, receiver string) error {
	fmt.Println("WeChat sent to ", receiver, " with message ", message)
	return nil
}

type App struct {
	messageService MessageService
}

func NewApp(messageService MessageService) *App {
	return &App{messageService: messageService}
}

func (a *App) Notify(message string, receiver string) error {
	return a.messageService.Send(message, receiver)
}
