package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

type ArticleModel struct {
	Title   string `json:"title"`
	Content string `json:"content"`
}

// 返回响应结构体
type Response struct {
	Code    int    `json:"code"`
	Data    any    `json:"data"`
	Message string `json:"message"`
}

// 封装函数
func _bindJson(c *gin.Context, obj any) (err error) {
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

func _getList(c *gin.Context) {
	articleList := []ArticleModel{
		{"成长手册", "内容1"},
		{"js深入学习", "内容2"},
		{"golang的修炼", "内容3"},
	}
	c.JSON(http.StatusOK, Response{0, articleList, "获取成功"})
}
func _getDetail(c *gin.Context) {
	fmt.Printf("文章id: %s\n", c.Param("id"))
	article := ArticleModel{"成长手册", "内容1"}
	c.JSON(http.StatusOK, Response{0, article, "获取成功"})
}
func _create(c *gin.Context) {
	// 接收前端传递的json数据
	/*
		请求: curl -X POST http://localhost/articles -H "Content-Type: application/json" -d '{"title":"test", "content":"123123123"}'
		响应: {"code":0,"data":{"title":"test","content":"123123123"},"message":"添加成功"}
	*/
	article := ArticleModel{}
	err := _bindJson(c, &article)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	c.JSON(http.StatusOK, Response{0, article, "添加成功"})
}
func _update(c *gin.Context) {
	fmt.Printf("文章id: %s\n", c.Param("id"))
	article := ArticleModel{}
	err := _bindJson(c, &article)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	c.JSON(http.StatusOK, Response{0, article, "修改成功"})
}
func _delete(c *gin.Context) {
	fmt.Printf("文章id: %s\n", c.Param("id"))
	c.JSON(http.StatusOK, Response{0, map[string]string{}, "删除成功"})
}
func main() {
	router := gin.Default()
	router.GET("/articles", _getList)       // 文章列表
	router.GET("/articles/:id", _getDetail) // 文章详情
	router.POST("/articles", _create)       // 文章创建
	router.PUT("/articles/:id", _update)    // 文章修改
	router.DELETE("/articles/:id", _delete) // 文章删除
	router.Run(":80")
}
