package api

import (
	"booking-app/micro-service/cluster/common"
	"booking-app/micro-service/cluster/common/core"

	"booking-app/micro-service/cluster/activity/model"

	"github.com/gin-gonic/gin"
)

func QueryActivity(c *gin.Context, opts *core.Options) (data interface{}, err error) {
	notices, err := model.QueryNoticeData(c, opts)
	if err != nil {
		return nil, err
	}

	return common.Response{
		Code:    common.SuccessCode,
		Message: "success",
		Data:    notices,
	}, nil
}

func InsertActivity(c *gin.Context, opts *core.Options) (data interface{}, err error) {
	params := &model.Notice{}
	if err != c.BindJSON(params) {
		return nil, err
	}

	// 存入内存
	params.ID = len(model.GlobalNotice) + 1
	model.GlobalNotice[params.ID] = params

	err = model.InsertNoticeData(c, params, opts)
	if err != nil {
		return nil, err
	}

	return common.Response{
		Code:    common.SuccessCode,
		Message: "success",
		Data:    nil,
	}, nil
}

func UpdateActivity(c *gin.Context, opts *core.Options) (data interface{}, err error) {
	return "UpdateActivity", nil
}
