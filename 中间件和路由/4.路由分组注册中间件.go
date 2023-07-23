package main

import (
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

func Middleware(ctx *gin.Context) {
	token := ctx.GetHeader("Authorization")
	if token == GlabalToken {
		ctx.Next()
		return
	}
	ctx.JSON(401, Res{1, nil, "身份验证失败"})
	ctx.Abort()
}

func UserRouterInit(router *gin.RouterGroup) {
	// 分组添加中间件校验
	userManager := router.Group("user").Use(Middleware)
	{
		userManager.GET("/index", UserManagerIndex)
	}
}

func main() {
	router := gin.Default()

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
