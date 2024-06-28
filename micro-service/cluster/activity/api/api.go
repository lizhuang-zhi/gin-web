package api

import (
	"booking-app/micro-service/cluster/common"
	"booking-app/micro-service/cluster/common/logger"

	"booking-app/micro-service/cluster/activity/model"

	"github.com/gin-gonic/gin"
)

func QueryActivity(c *gin.Context) (data interface{}, err error) {
	notices, err := model.QueryNoticeData(c)
	if err != nil {
		logger.Errorf("query notice data err:%v", err)
		return nil, err
	}

	return common.Response{
		Code:    common.SuccessCode,
		Message: "success",
		Data:    notices,
	}, nil
}

func InsertActivity(c *gin.Context) (data interface{}, err error) {
	params := &model.Notice{}
	if err != c.BindJSON(params) {
		logger.Warnf("bind json err:%v", err)
		return nil, err
	}

	// 存入内存
	params.ID = len(model.GlobalNotice) + 1
	model.GlobalNotice[params.ID] = params

	logger.Infof("insert notice data to memory, id:%d, title:%s, sub_type:%d, content:%s", params.ID, params.Title, params.SubType, params.Content)

	err = model.InsertNoticeData(c, params)
	if err != nil {
		logger.Errorf("insert notice data to mongo err:%v", err)
		return nil, err
	}

	return common.Response{
		Code:    common.SuccessCode,
		Message: "success",
		Data:    nil,
	}, nil
}

func UpdateActivity(c *gin.Context) (data interface{}, err error) {
	return "UpdateActivity", nil
}
