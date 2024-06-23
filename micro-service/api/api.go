package api

import (
	"booking-app/micro-service/common"
	"booking-app/micro-service/model"

	"github.com/gin-gonic/gin"
)

func QueryActivity(c *gin.Context) (data interface{}, err error) {
	notices := make([]*model.Notice, 0)

	for _, v := range model.GlobalNotice {
		notices = append(notices, v)
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
		return nil, err
	}

	params.ID = len(model.GlobalNotice) + 1
	model.GlobalNotice[params.ID] = params

	return common.Response{
		Code:    common.SuccessCode,
		Message: "success",
		Data:    params.ID,
	}, nil
}

func UpdateActivity(c *gin.Context) (data interface{}, err error) {
	return "UpdateActivity", nil
}
