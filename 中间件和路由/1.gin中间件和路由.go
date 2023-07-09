package main

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

// 注册中件间
func m1(ctx *gin.Context) {
	// 请求中间件
	fmt.Println("m1...in")
	ctx.Next()
	// 响应中间件
	fmt.Println("m1...out")
}

// 注册中件间
func index(ctx *gin.Context) {
	// 请求中间件
	fmt.Println("index...in")
	ctx.Next()
	ctx.JSON(200, gin.H{"msg": "index的响应"})
	// 响应中间件
	fmt.Println("index...out")
}

// 注册中件间
func m2(ctx *gin.Context) {
	// 请求中间件
	fmt.Println("m2...in")
	// 响应中间件
	fmt.Println("m2...out")
}

func main() {
	router := gin.Default()

	// 后面跟上的函数都可以理解为中件间
	router.GET("/", func(ctx *gin.Context) {
		ctx.JSON(200, gin.H{"msg": "index"})
	}, func(ctx *gin.Context) {
		fmt.Println(1)
	}, func(ctx *gin.Context) {
		fmt.Println(2)
	}, func(ctx *gin.Context) {
		fmt.Println(3)
	})

	// 洋葱模型
	router.GET("/index", m1, index, m2)

	router.Run(":8080")
}
