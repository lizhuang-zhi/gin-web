package main

import "github.com/sirupsen/logrus"

func main()  {
	// 设置JSON输出格式(默认为Text格式)
	// logrus.SetFormatter(&logrus.JSONFormatter{})
	logrus.SetFormatter(&logrus.TextFormatter{})

	log := logrus.WithField("app", "study").WithField("service", "logrus")
	// log1 := logrus.WithFields(logrus.Fields{
	// 	"user_id": "21",
	// 	"ip": "127.0.0.1",
	// })
	log1 := log.WithFields(logrus.Fields{
		"user_id": "21",
		"ip": "127.0.0.1",
	})

	log.Errorf("您好")
	log1.Error("Hello")
}