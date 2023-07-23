package main

import (
	"github.com/gin-gonic/gin"
)

type UserInfo struct {
	Name string `json:"name"`
	Age  int    `json:"age"`
}

type ArticleInfo struct {
	Title string `json:"title"`
	Desc  string `json:"desc"`
}

type Response struct {
	Code int         `json:"code"`
	Data interface{} `json:"data"`
	Msg  string      `json:"msg"`
}

func UserManagerIndex(ctx *gin.Context) {
	userList := []UserInfo{
		{Name: "张三", Age: 18},
		{Name: "李四", Age: 19},
		{Name: "王五", Age: 20},
	}
	ctx.JSON(200, Response{0, userList, "请求成功"})
}

func ArticleManagerIndex(ctx *gin.Context) {
	articleList := []ArticleInfo{
		{Title: "文章1", Desc: "文章1的描述"},
		{Title: "文章2", Desc: "文章2的描述"},
	}
	ctx.JSON(200, Response{0, articleList, "请求成功"})
}

func UserRouterInit(router *gin.RouterGroup) {
	// 用户管理路由分组
	userManager := router.Group("user")
	{
		userManager.GET("/index", UserManagerIndex)
	}
}

func ArticleRouterInit(router *gin.RouterGroup) {
	// 文章管理路由分组
	articleManager := router.Group("article")
	{
		articleManager.GET("/index", ArticleManagerIndex)
	}
}

func main() {
	router := gin.Default()

	// 路由分组
	api := router.Group("api")

	UserRouterInit(api)
	ArticleRouterInit(api)

	router.Run(":8080")
}
