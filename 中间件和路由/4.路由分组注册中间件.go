package main

import (
	"fmt"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

type Res struct {
	Code int         `json:"code"`
	Data interface{} `json:"data"`
	Msg  string      `json:"msg"`
}

var GlabalToken = "123"

func UserManagerIndex(ctx *gin.Context) {
	type UserInfo struct {
		Name string `json:"name"`
		Age  int    `json:"age"`
	}
	userList := []UserInfo{
		{Name: "张三", Age: 18},
		{Name: "李四", Age: 19},
		{Name: "王五", Age: 20},
	}
	ctx.JSON(200, Res{0, userList, "请求成功"})
}

// 验证中间件(闭包写法)
func AuthMiddleware(msg string) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		token := ctx.GetHeader("Authorization")
		if token == GlabalToken {
			ctx.Next()
			return
		}
		ctx.JSON(401, Res{1, nil, msg})
		ctx.Abort()
	}
}

// 计算请求耗时中间件
func TotalTimeMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		startTime := time.Now()
		ctx.Next()
		since := time.Since(startTime)
		currentFunc := ctx.HandlerName()
		fmt.Printf("===>函数 %s 耗时⌚️ %d ns\n", currentFunc, since)
	}
}

func UserRouterInit(router *gin.RouterGroup) {
	// 分组添加中间件校验
	userManager := router.Group("user").Use(AuthMiddleware("身份验证失败"))
	{
		userManager.GET("/index", UserManagerIndex)
	}
}

func main() {
	router := gin.New()

	// 添加计算耗时中间件
	router.Use(TotalTimeMiddleware(), gin.LoggerWithFormatter(func(params gin.LogFormatterParams) string {
		// 修改日志格式
		codeString := strconv.Itoa(params.StatusCode)
		return "===>新的日志格式😄:    " + codeString + "     " + params.Path + "     " + params.Method + "    " + params.ResetColor() + "\n"
	}), gin.Recovery())

	// 路由分组
	api := router.Group("api")

	// 不需要登录校验
	api.GET("/login", func(ctx *gin.Context) {
		type LoginToken struct {
			Token string `json:"token"`
		}
		dataToken := LoginToken{
			GlabalToken,
		}
		ctx.JSON(200, Res{0, dataToken, "登录成功"})
	})

	UserRouterInit(api)

	router.Run(":8080")
}
