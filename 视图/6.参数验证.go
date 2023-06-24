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

		gin的内置验证器:
		枚举:
		1. oneof=girl boy  枚举: 只能是girl或boy

		字符串:
		2. contains=f      包含: 必须包含f
		3. excludes=f      不包含: 必须不包含f
		4. startswith=f    开头: 必须f开头
		5. endswith=f      结尾: 必须f结尾

		数组:
		6. dive            dive后面的验证,针对数组中的每个元素

		网络严重
		ip                 验证ip是否正确
		uri                验证url  例如: http://localhost/123
		url                验证uri  例如: http://localhost/123 或 /123

		日期格式验证(这个时间是固定的: 1月2号下午三点四分五秒06年)
		datetime           验证日期格式

	*/
	Name       string   `json:"name" binding:"required"`                           // 用户名
	Age        int      `json:"age" binding:"gt=18,lt=100"`                        // 年龄
	Password   string   `json:"password" binding:"min=4,max=7"`                    // 密码
	RePassword string   `json:"re_password" binding:"eqfield=Password"`            // 确认密码
	Gender     string   `json:"gender" binding:"oneof=girl boy"`                   // 性别
	LikeList   []string `json:"like_list" binding:"required,dive,startswith=like"` // 爱好
	IP         string   `json:"ip" binding:"ip"`                                   // ip地址
	Url        string   `json:"url" binding:"url"`                                 // url地址
	Uri        string   `json:"uri" binding:"uri"`                                 // uri地址
	Date       string   `json:"date" binding:"datetime=2006-01-02 15:04:05"`       // 日期                                 // 日期
}

// 通过验证的示例
/*
	{
		"name": "leo li",
		"age": 28,
		"password": "12348",
		"re_password": "12348",
		"gender": "boy",
		"like_list": [
			"like basketball"
		],
		"ip": "255.255.255.0",
		"url": "http://localhost/123",
		"uri": "/123",
		"date": "2022-02-04 12:01:24"
	}
*/

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
