package main

import (
	"fmt"
	"strconv"

	"github.com/gin-gonic/gin"
)

// 自定义输出日志格式
func LogFormatterParams(params gin.LogFormatterParams) string {
	return fmt.Sprintf(
		"[ Leo Print ] %s | %s | %s %s\n",
		params.TimeStamp.Format("2006-01-02 15:04:05"),
		params.StatusCodeColor()+strconv.Itoa(params.StatusCode)+params.ResetColor(),
		params.MethodColor()+params.Method+params.ResetColor(),
		params.Path,
	)
}

func index(ctx *gin.Context) {}
func main() {
	// 将日志记录在gin.log中
	// file, _ := os.Create("gin.log")
	// gin.DefaultWriter = io.MultiWriter(file, os.Stdout)

	// 定义路由格式
	// gin.DebugPrintRouteFunc = func(httpMethod, absolutePath, handlerName string, nuHandlers int) {
	// 	log.Printf("[ Leo ] %s %s %s %d\n", httpMethod, absolutePath, handlerName, nuHandlers)
	// }

	// 不想看到debug模式, 转为relase模式(会删掉开头的debug打印)
	// gin.SetMode(gin.ReleaseMode)

	router := gin.New()
	// router.Use(gin.LoggerWithFormatter(LogFormatterParams))
	// 等价
	router.Use(gin.LoggerWithConfig(gin.LoggerConfig{
		Formatter: LogFormatterParams,
	}))

	router.GET("/index", index)
	router.POST("/users", func(ctx *gin.Context) {})
	router.POST("/articles", func(ctx *gin.Context) {})

	api := router.Group("api")
	api.DELETE("/articles/:id", func(ctx *gin.Context) {})

	// 打印定义的所有路由
	// fmt.Println(router.Routes())
	// for _, v := range router.Routes() {
	// 	fmt.Println(v.Handler)
	// 	fmt.Println(v.Path)
	// 	fmt.Println(v.HandlerFunc)
	// 	fmt.Println(v.Method)
	// }

	router.Run()
}
