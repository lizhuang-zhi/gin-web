package gin_server

import (
	"time"

	"github.com/gin-gonic/gin"
)

// InitGinServer 初始化gin服务
func InitGinServer() {
	router := gin.Default()

	router.GET("/timeout", func(c *gin.Context) {
		// 等待5s后返回
		time.Sleep(5 * time.Second)

		c.JSON(0, gin.H{"data": "返回用户对应的数据"})
	})

	router.Run(":7878")
}
