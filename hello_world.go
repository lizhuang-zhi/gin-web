package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func Index(ctx *gin.Context) {
	ctx.String(200, "Hello World! Leo is Good")
}

func main() {
	router := gin.Default()

	router.GET("/index", Index)

	// 启动监听：gin会把web服务运行在本机的0.0.0.0:8080端口上
	router.Run(":8080")
	// 用原生http服务的方式, router.Run本质就是的进一步封装（可以看源码）
	http.ListenAndServe(":8080", router)
}
