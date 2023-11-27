package main

import "github.com/sirupsen/logrus"

func main() {
	logrus.SetFormatter(&logrus.TextFormatter{
		ForceColors:   true, // 是否强制使用颜色输出
		FullTimestamp: true, // 是否在连接到 TTY 时输出完整的时间戳
		TimestampFormat: "2006-01-02 15:04:05",  // 输出时间格式
	})
	logrus.SetLevel(logrus.DebugLevel)
	logrus.Error("Hello")
	logrus.Warnln("Hello")
	logrus.Info("Hello")
	logrus.Debug("Hello")
	logrus.Println("Hello")
}
