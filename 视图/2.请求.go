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

func main() {
	router := gin.Default()
	router.GET("/query", _query)
	router.GET("/param/:user_id", _param)
	router.GET("/param/:user_id/:book_id", _param)
	router.Run(":80")
}
