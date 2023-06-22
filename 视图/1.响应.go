package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func _string(c *gin.Context) {
	c.String(http.StatusOK, "Hello")
}

// 响应json数据 重点
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

// 响应xml
func _xml(c *gin.Context) {
	c.XML(http.StatusOK, gin.H{"username": "Leo Li", "age": 28, "hobbies": gin.H{
		"basketball": true,
		"football":   true,
	}})
}

// 响应yaml
func _yaml(c *gin.Context) {
	c.YAML(http.StatusOK, gin.H{"username": "Leo Li", "age": 28, "hobbies": gin.H{
		"basketball": true,
		"football":   true,
	}})
}

// 响应html
func _html(c *gin.Context) {
	c.HTML(http.StatusOK, "index.html", gin.H{"msg": "Leo Good"})
}

// 重定向
func _redirect(c *gin.Context) {
	// 临时重定向到对应路由
	// c.Redirect(http.StatusFound, "https://www.baidu.com")
	c.Redirect(http.StatusFound, "/html") 
}

func main() { 
	router := gin.Default()
	// 加载模版目录下所有的模版文件(响应html需要配置的模版)
	router.LoadHTMLGlob("../templates/*")
	// 静态文件响应参数: 1. 路由 2. 相对文件路径  请求示例: http://localhost/bg
	router.StaticFile("/bg", "../static/background.jpeg")
	// 文件响应参数: 1. 路由  2. 可访问目录的相对路径   请求示例: http://localhost/static/home.txt
	router.StaticFS("/static", http.Dir("../static/static_real"))

	router.GET("/", _string)
	router.GET("/json", _json)
	router.GET("/xml", _xml)
	router.GET("/yaml", _yaml)
	router.GET("/html", _html)
	router.GET("/baidu", _redirect)
	router.Run(":80")
}
