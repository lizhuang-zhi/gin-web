package main

import (
	"github.com/gin-gonic/gin"
)

type SignUserInfo struct {
	/*
		binding添加验证:

		针对所有:
		1. required: 必须传入字段且不能为空串

		针对字符串:
		2. min: 最小长度
		3. max: 最大长度
		4. len: 长度为len

		针对数字:
		5. lt: 小于
		6. lte: 小于等于
		7. gt: 大于
		8. gte: 大于等于
		9. eq: 等于
		10. ne: 不等于

		针对同结构体字段:
		11. eqfield: 等于其他字段
		12. nefield: 不等于其他字段

		忽略字段:
		13. binding="-"
	*/
	Name       string `json:"name" binding:"required"`        // 用户名
	Age        int    `json:"age" binding:"gt=18,lt=100"`     // 年龄
	Password   string `json:"password" binding:"min=4,max=7"` // 密码
	RePassword string `json:"re_password" binding:"eqfield=Password"` // 确认密码
}

func main() {
	router := gin.Default()
	router.POST("/", func(c *gin.Context) {
		signUserInfo := SignUserInfo{}
		err := c.ShouldBindJSON(&signUserInfo)
		if err != nil {
			c.JSON(200, gin.H{"msg": err.Error()})
			return
		}
		c.JSON(0, signUserInfo)
	})
	router.Run(":80")
}
