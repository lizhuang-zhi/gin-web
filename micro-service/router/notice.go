package router

import (
	"booking-app/micro-service/api"

	"github.com/gin-gonic/gin"
)

func InitNoticeRouter(router *gin.RouterGroup) {
	NoticeRouter := router.Group("notice")
	NoticeRouter.GET("/query", WrapHandle(api.QueryActivity))
	NoticeRouter.POST("/insert", WrapHandle(api.InsertActivity))
	NoticeRouter.POST("/udpate", WrapHandle(api.UpdateActivity))
}
