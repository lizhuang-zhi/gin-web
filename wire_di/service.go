package main

import "fmt"

type MessageService interface {
	Send(message string, receiver string) error
}

type EmailService struct {
}

func (e *EmailService) Send(message string, receiver string) error {
	fmt.Println("Email sent to ", receiver, " with message ", message)
	return nil
}

type WeChatService struct {
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
