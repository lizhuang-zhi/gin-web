package main

import (
	"log"

	"github.com/gin-gonic/gin"
)

func main() {
	messageService := &EmailService{}
	app, err := InitializeApp(messageService)
	if err != nil {
		log.Fatal(err)
	}

	messageService2 := &WeChatService{}
	app2, err2 := InitializeApp(messageService2)
	if err2 != nil {
		log.Fatal(err)
	}

	// 使用Gin
	r := gin.Default()
	r.GET("/ping", func(c *gin.Context) {
		err := app.Notify("Ping received", "user@example.com")
		if err != nil {
			c.JSON(500, gin.H{"message": err.Error()})
		} else {
			c.JSON(200, gin.H{"message": "pong"})
		}
	})

	r.GET("/ping2", func(c *gin.Context) {
		err := app2.Notify("Ping 2 received", "leo li")
		if err != nil {
			c.JSON(500, gin.H{"message": err.Error()})
		} else {
			c.JSON(200, gin.H{"message": "pong"})
		}
	})
	r.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}
