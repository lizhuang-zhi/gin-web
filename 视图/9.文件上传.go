package main

import (
	"fmt"
	"io"
	"os"

	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()

	// 文件传递通过 form-data
	router.POST("/upload", func(ctx *gin.Context) {
		file, _ := ctx.FormFile("pictureFile")
		fmt.Println(file.Filename)                      // 文件名
		fmt.Println(file.Size / 1024)                   // file.Size 单位是字节(Byte)
		ctx.SaveUploadedFile(file, "../uploads/12.png") // 将文件存储到对应文件夹
		ctx.JSON(200, gin.H{"msg": "上传成功"})
	})

	// 读取上传文件
	router.POST("/upload/read", func(ctx *gin.Context) {
		file, _ := ctx.FormFile("pictureFile")
		// 读取文件中的数据, 返回文件对象
		readerFile, _ := file.Open()
		data, _ := io.ReadAll(readerFile)
		fmt.Printf(string(data))
		ctx.JSON(200, gin.H{"msg": "上传成功"})
	})

	// create + copy
	router.POST("/upload/copy", func(ctx *gin.Context) {
		file, _ := ctx.FormFile("pictureFile")
		// 读取请求中的文件数据
		readerFile, _ := file.Open()
		// 创建文件
		writerFile, _ := os.Create("../uploads/new.png")
		// 关闭文件
		defer writerFile.Close()
		// 将请求中的文件数据copy到创建的文件中
		n, _ := io.Copy(writerFile, readerFile)
		fmt.Println(n) // 返回文件大小
		ctx.JSON(200, gin.H{"msg": "上传成功"})
	})

	// 上传多个文件
	router.POST("/upload/multi", func(ctx *gin.Context) {
		form, _ := ctx.MultipartForm()
		files := form.File["upload[]"]
		for _, file := range files {
			ctx.SaveUploadedFile(file, "../uploads/"+file.Filename)
		}
		ctx.JSON(200, gin.H{"msg": fmt.Sprintf("上传成功 %d 个文件", len(files))})
	})

	router.Run(":8080")
}
