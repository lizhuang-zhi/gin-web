package main

import "github.com/sirupsen/logrus"

func main()  {
	// 设置日志等级(默认Info以后不输出到控制台)
	logrus.SetLevel(logrus.DebugLevel)
	// 从上到下四种等级
	logrus.Error("出错了")
	logrus.Warnln("警告")
	logrus.Infof("信息")
	logrus.Debug("debug")

	logrus.Println("打印")
}