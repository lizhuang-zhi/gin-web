package router

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Handler func(c *gin.Context) (data interface{}, err error)

func Routers() *gin.Engine {
	r := gin.Default()
	api := r.Group("")

	InitNoticeRouter(api) // 公告相关路由

	return r
}

func WrapHandle(handler Handler) gin.HandlerFunc {
	return func(c *gin.Context) {
		data, err := handler(c)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": err.Error(),
			})
			return
		}
		c.JSON(http.StatusOK, data)
	}
}
