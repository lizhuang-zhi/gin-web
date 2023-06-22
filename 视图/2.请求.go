package main

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

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

func _param(c *gin.Context) {
	// 请求地址: http://localhost/param/sdk/123
	fmt.Println(c.Param("user_id")) // sdk
	fmt.Println(c.Param("book_id")) // 123
}

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
	forms, err := c.MultipartForm()  // 接收所有的参数, 包括文件
	fmt.Println(forms, err) // &{map[addr:[陕西省] name:[leo li Tom Zhang]] map[file:[0xc0004ce000]]} <nil>
}

func main() {
	router := gin.Default()
	router.GET("/query", _query)
	router.GET("/param/:user_id", _param)
	router.GET("/param/:user_id/:book_id", _param)
	router.POST("/form", _form)
	router.Run(":80")
}
