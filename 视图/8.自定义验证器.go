package main

import (
	"reflect"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

// _GetValidMsg: 返回结构体中的msg参数
func _GetValidMsg(err error, obj interface{}) string {
	getObj := reflect.TypeOf(obj)
	// 将err接口断言为具体类型
	if errs, ok := err.(validator.ValidationErrors); ok {
		// 断言成功
		for _, e := range errs {
			// 循环每个错误信息
			// 根据报错字段名, 获取结构体的具体字段
			if field, exist := getObj.Elem().FieldByName(e.Field()); exist {
				msg := field.Tag.Get("msg")
				return msg
			}
		}
	}
	return err.Error()
}

type User struct {
	Name string `json:"name" binding:"required,sign" msg:"用户名校验错误"`
	Age  int    `json:"age" binding:"required" msg:"年龄校验错误"`
}

// name的自定义验证器
func signValid(fl validator.FieldLevel) bool {
	var nameList []string = []string{"Leo", "Tom", "Jim"}
	for _, nameStr := range nameList {
		name := fl.Field().Interface().(string)
		if name == nameStr {
			return false
		}
	}
	return true
}

func main() {
	router := gin.Default()

	// 自定义验证器
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation("sign", signValid)
	}
	router.POST("/", func(c *gin.Context) {
		user := User{}
		err := c.ShouldBindJSON(&user)
		if err != nil {
			c.JSON(200, gin.H{"msg": _GetValidMsg(err, &user)})
			return
		}
		c.JSON(0, user)
	})
	router.Run(":80")
}
