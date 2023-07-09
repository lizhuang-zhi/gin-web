package main

import (
	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()

	// 文件显示(存在缓存问题: http://localhost:8080/show?id=1, 改变id值, 或者记得关闭浏览器缓存(F12工具))
	router.GET("/show", func(ctx *gin.Context) {
		// 只是显示图片
		// ctx.File("/Users/mrkleo/Program/project/BackEnd/golang-learn/gin-web/uploads/动漫风景-2.jpeg")
		// 显示go文件
		ctx.File("/Users/mrkleo/Program/project/BackEnd/golang-learn/gin-web/视图/2.请求.go")
		// ctx.JSON(200, gin.H{"msg": "展示成功"})
	})

	// 文件下载
	router.GET("/download", func(ctx *gin.Context) {
		// 设置文件流响应头, 唤起浏览器下载
		ctx.Header("Content-Type", "application/octet-stream")
		// 指定下载下来的文件名
		ctx.Header("Content-Disposition", "attachment; filename="+"牛逼123.png")
		// 表示传输过程中的编码形式, 乱码问题可能就是因为它
		ctx.Header("Content-Transfer-Encoding", "binary")
		// 下载的文件
		ctx.File("../uploads/12.png")
		ctx.JSON(200, gin.H{"msg": "下载成功"})
	})

	router.Run(":8080")
}
