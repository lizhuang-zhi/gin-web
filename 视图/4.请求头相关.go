package main

import (
	"fmt"
	"strings"

	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()
	// 请求头的各种获取方式
	router.GET("/", func(c *gin.Context) {
		// 不区分大小写
		fmt.Println(c.GetHeader("User-Agent")) // client
		// c.Request.Header 是一个 map[string][]string
		fmt.Println(c.Request.Header) // map[Accept:[*/*] Accept-Encoding:[gzip, deflate, br] Connection:[keep-alive] Content-Length:[39] Content-Type:[application/json] User-Agent:[client]]
		// 使用 Get方法 或者 .GetHeader 可以不用区分大小写, 并且返回第一个value
		fmt.Println(c.Request.Header.Get("User-Agent")) // client
		// 如果是用map的取值方式, 请注意大小写问题
		fmt.Println(c.Request.Header["User-Agent"]) // [client]

		// 自定义请求头
		fmt.Println(c.Request.Header.Get("Token"))
		fmt.Println(c.Request.Header.Get("token"))
	})

	// 爬虫和用户的区别对待
	router.GET("/index", func(c *gin.Context) {
		userAgent := c.GetHeader("User-Agent")
		// 判断userAgent是否包含python子串
		if strings.Contains(userAgent, "python") {
			// 检测到爬虫程序
			c.JSON(0, gin.H{"data": "返回爬虫对应的数据"})
			return 
		}
		c.JSON(0, gin.H{"data": "返回用户对应的数据"})
	})

	// 设置响应头
	router.GET("/res", func (c *gin.Context)  {
		// 设置响应头token
		c.Header("Token", "ssdfkskaoals99888888")
		// 设置响应头Content-Type为text文本
		c.Header("Content-Type", "application/text; charset=utf-8")
		c.JSON(0, gin.H{"data": "查看响应头"})
	})

	router.Run(":80")
}
