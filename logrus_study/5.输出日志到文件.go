package main

import (
	"io"
	"os"

	"github.com/sirupsen/logrus"
)

func main() {
	// 输出到日志文件(追加形式)
	file, _ := os.OpenFile("./info.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	// file, _ := os.Create("./info.log")   // 底层还是调用的os.OpenFile(覆盖之前的)
	
	// logrus.SetOutput(file)
	logrus.SetOutput(io.MultiWriter(file, os.Stdout))

	logrus.Infof("你好")
	logrus.Error("出错了")
	logrus.Errorf("出错了 %s", "xxxx")
	logrus.Errorln("出错了")
}