package main

import (
	"github.com/gin-gonic/gin"
)

type UserInfo struct {
	Name   string `json:"name" form:"name" uri:"name"`
	Age    int    `json:"age" form:"age" uri:"age"`
	Gender string `json:"gender" form:"gender" uri:"gender"`
}

func main() {
	router := gin.Default()
	// 请求体Body使用json格式
	/*
		{
		    "name": "leo li",
		    "age": 22,
		    "gender": "男"
		}
	*/
	router.POST("/", func(c *gin.Context) {
		userInfo := UserInfo{}
		err := c.ShouldBindJSON(&userInfo)
		if err != nil {
			c.JSON(200, gin.H{"msg": err.Error()})
			return
		}
		c.JSON(0, userInfo)
	})

	// 需要在结构体中设置 form:"name"
	// 请求: /query?name=leo&age=22&gender=boy
	router.POST("/query", func(c *gin.Context) {
		userInfo := UserInfo{}
		err := c.ShouldBindQuery(&userInfo)
		if err != nil {
			c.JSON(200, gin.H{"msg": err.Error()})
			return
		}
		c.JSON(0, userInfo)
	})

	// 需要在结构体中设置 uri:"name"
	// 请求: /uri/leo/22/boy
	router.POST("/uri/:name/:age/:gender", func(c *gin.Context) {
		userInfo := UserInfo{}
		err := c.ShouldBindUri(&userInfo)
		if err != nil {
			c.JSON(200, gin.H{"msg": err.Error()})
			return
		}
		c.JSON(0, userInfo)
	})

	// 请求为Body为form-data格式
	// 需要在结构体中设置 form:"name"
	// 请求: /form
	router.POST("/form", func(c *gin.Context) {
		userInfo := UserInfo{}
		err := c.ShouldBind(&userInfo)
		if err != nil {
			c.JSON(200, gin.H{"msg": err.Error()})
			return
		}
		c.JSON(0, userInfo)
	})

	router.Run(":80")
}
