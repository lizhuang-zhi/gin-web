package main

func ProvideEmailService(cfg *Config) *EmailService {
	return NewEmailService(cfg.EmailUsername, cfg.EmailPassword)
}

func ProvideWeChatService(cfg *Config) *WeChatService {
	return NewWeChatService(cfg.WeChatAccountName, cfg.WeChatAccountPassword)
}

func ProvideMessageService(cfg *Config, emailService *EmailService, weChatService *WeChatService) MessageService {
	if cfg.messageChoose == "email" {
		return emailService
	} else {
		return weChatService
	}
}
