package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func _string(c *gin.Context) {
	c.String(http.StatusOK, "Hello")
}

func _json(c *gin.Context) {
	// json响应结构体
	// type UserInfo struct {
	// 	UserName string `json:"user_name"`
	// 	Age      int `json:"age"`
	// 	Password string `json:"-"`   // 设置“-”不进行json序列化, 返回结果不显示
	// }
	// user := UserInfo{
	// 	"Leo Li",
	// 	22,
	// 	"123456",
	// }
	// c.JSON(http.StatusOK, user)

	// json响应map(注意: map输出顺序是随机的)
	// userMap := map[string]string {
	// 	"user_name": "Leo Li",
	// 	"age": "22",
	// }
	// c.JSON(http.StatusOK, userMap)

	// 直接响应json
	c.JSON(http.StatusOK, gin.H{"username": "Leo Li", "age": 28})
}

func main() {
	router := gin.Default()
	router.GET("/", _string)
	router.GET("/json", _json)
	router.Run(":80")
}
