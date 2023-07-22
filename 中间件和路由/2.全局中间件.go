package main

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

type User struct {
	Name string `json:"name"`
	Age  int    `json:"age"`
}

func global_1(c *gin.Context) {
	// 被添加到上下文中的值，可以在后续的中间件和处理程序中获取
	c.Set("name", "leo")
	c.Set("user", User{
		Name: "Leo",
		Age:  22,
	})
	c.Next()
}

func main() {
	router := gin.Default()

	// 添加多个全局中间件
	router.Use(global_1)

	router.GET("/", func(ctx *gin.Context) {
		// 获取上下文中的值
		value, isExist := ctx.Get("name")
		if isExist {
			fmt.Println(value) // leo
		}
		_user, _ := ctx.Get("user")
		// 添加断言，获取实例具体属性 use.xxx
		user, ok := _user.(User)
		// 判断是否断言成功（是否是User类型）
		if ok {
			fmt.Println(user.Name) // Leo
			fmt.Println(user.Age)  // 22
		}
		ctx.JSON(200, gin.H{"msg": "m10"})
	})

	router.Run(":8080")
}
