package main

import (
	"encoding/json"
	"fmt"

	"github.com/gin-gonic/gin"
)

// 查询参数
func _query(c *gin.Context) {
	// 请求地址: http://localhost/query?id=&user=leo
	fmt.Println(c.GetQuery("user"))   // leo true
	fmt.Println(c.GetQuery("id"))     //  true
	fmt.Println(c.GetQuery("test"))   //  false
	fmt.Println(c.QueryArray("user")) // [leo]
	fmt.Println(c.Query("user"))      // leo
	// 请求地址: http://localhost/query?id=&user[name]=leo
	fmt.Println(c.QueryMap("user")) // map[name:leo]
}

// 动态参数
func _param(c *gin.Context) {
	// 请求地址: http://localhost/param/sdk/123
	fmt.Println(c.Param("user_id")) // sdk
	fmt.Println(c.Param("book_id")) // 123
}

// 表单参数
func _form(c *gin.Context) {
	/*
		传递的form-data: name: leo li
					    name: Tom Zhang
						addr: 陕西省
						file: test.png(一张图片)
	*/
	fmt.Println(c.PostForm("name"))               // leo li
	fmt.Println(c.PostFormArray("name"))          // [leo li Tom Zhang]
	fmt.Println(c.DefaultPostForm("addr", "四川省")) // 陕西省
	forms, err := c.MultipartForm()               // 接收所有的参数, 包括文件
	fmt.Println(forms, err)                       // &{map[addr:[陕西省] name:[leo li Tom Zhang]] map[file:[0xc0004ce000]]} <nil>
}

// 原始参数
func _raw(c *gin.Context) {
	type User struct {
		Name string `json="name"`
		Age  int    `json="age"`
	}
	user := User{}
	err := bindJson(c, &user)
	if err != nil {
		fmt.Println(err.Error())
	}
	fmt.Println(user)
}

// 封装函数 ( any等价interface{} )
func bindJson(c *gin.Context, obj any) (err error) {
	data, _ := c.GetRawData()
	// 获取请求头信息
	contentType := c.GetHeader("Content-Type")
	switch contentType {
	// 处理json类型的数据
	case "application/json":
		// 将json数据解析为结构体(这个地方写obj或者&obj都可以获得结果)
		err := json.Unmarshal(data, obj)
		if err != nil {
			fmt.Println(err.Error())
			return err
		}
	}
	return nil
}

func main() {
	router := gin.Default()
	router.GET("/query", _query)
	router.GET("/param/:user_id", _param)
	router.GET("/param/:user_id/:book_id", _param)
	router.POST("/form", _form)
	router.POST("/raw", _raw)
	router.Run(":80")
}
