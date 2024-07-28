package gin_server

import (
	"time"

	"github.com/gin-gonic/gin"
)

// InitGinServer 初始化gin服务
func InitGinServer() {
	router := gin.Default()

	router.GET("/index", func(c *gin.Context) {
		c.JSON(0, gin.H{"data": "返回用户对应的数据"})
	})

	router.GET("/res", func(c *gin.Context) {
		c.Header("Token", "ssdfkskaoals99888888")
		c.Header("Content-Type", "application/text; charset=utf-8")
		c.JSON(0, gin.H{"data": "查看响应头"})
	})

	router.POST("/insert", func(c *gin.Context) {
		// 等待5s后返回
		time.Sleep(5 * time.Second)

		// 获取请求体
		body, err := c.GetRawData()
		if err != nil {
			c.JSON(0, gin.H{"data": "获取请求体失败"})
		}

		c.JSON(0, gin.H{"data": string(body)})
	})

	router.Run(":7878")
}
