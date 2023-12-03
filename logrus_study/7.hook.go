package main

import (
	"os"

	"github.com/sirupsen/logrus"
)

/*
	type Hook interface {
		Levels() []Level
		Fire(*Entry) error
	}
*/

type MyHook struct{}

// 作用于什么等级的日志
func (h *MyHook) Levels() []logrus.Level {
	// return logrus.AllLevels  // 
	return []logrus.Level{logrus.ErrorLevel}
}

func (h *MyHook) Fire(entry *logrus.Entry) error {
	// entry.Data["app"] = "fengfeng"
	// fmt.Println(entry.Level)

	// 遇到Error的日志, 执行hook,写入文件
	file, _ := os.OpenFile("./err.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	line, _ := entry.String()
	file.Write([]byte(line))
	
	return nil
}


func main() {
	logrus.AddHook(&MyHook{})

	logrus.Warnln("你好")
	logrus.Errorf("你好")
}
