package main

import (
	"log"

	"github.com/gin-gonic/gin"
)

func main() {
	cfg := &Config{
		EmailUsername:         "myemail@example.com",
		EmailPassword:         "mypassword",
		WeChatAccountName:     "wechatAccount",
		WeChatAccountPassword: "wechatPassword",
		messageChoose:         "wechat", // 使用wechat作为message服务
	}

	app, err := InitializeApp(cfg)
	if err != nil {
		log.Fatal(err)
	}

	// 使用Gin
	r := gin.Default()
	r.GET("/ping", func(c *gin.Context) {
		err := app.Notify("wechat service", "leoli")
		if err != nil {
			c.JSON(500, gin.H{"message": err.Error()})
		} else {
			c.JSON(200, gin.H{"message": "pong"})
		}
	})

	r.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}
