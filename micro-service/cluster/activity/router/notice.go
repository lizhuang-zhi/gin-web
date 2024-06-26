package router

import (
	"booking-app/micro-service/cluster/activity/api"
	"booking-app/micro-service/cluster/common/core"

	"github.com/gin-gonic/gin"
)

func InitNoticeRouter(router *gin.RouterGroup, opts *core.Options) {
	NoticeRouter := router.Group("notice")
	NoticeRouter.GET("/query", WrapHandle(api.QueryActivity, opts))
	NoticeRouter.POST("/insert", WrapHandle(api.InsertActivity, opts))
	NoticeRouter.POST("/udpate", WrapHandle(api.UpdateActivity, opts))
}
