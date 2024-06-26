package router

import (
	"booking-app/micro-service/cluster/common/core"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Handler func(c *gin.Context, opts *core.Options) (data interface{}, err error)

func Routers(opts *core.Options) *gin.Engine {
	r := gin.Default()
	api := r.Group("")

	InitNoticeRouter(api, opts) // 公告相关路由

	return r
}

func WrapHandle(handler Handler, opts *core.Options) gin.HandlerFunc {
	return func(c *gin.Context) {
		data, err := handler(c, opts)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": err.Error(),
			})
			return
		}
		c.JSON(http.StatusOK, data)
	}
}
