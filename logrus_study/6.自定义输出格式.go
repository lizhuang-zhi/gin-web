package main

import (
	"bytes"
	"fmt"
	"path"

	"github.com/sirupsen/logrus"
)

const (
	ccGrayBlock   = 0
	ccRedBlock    = 1
	ccGreenBlock  = 2
	ccYellowBlock = 3
)

func printColor(colorCode int, text string, isBackground bool) {
	if !isBackground {
		fmt.Printf("\033[3%dm %s \033[0m\n", colorCode, text)
		return
	}
	fmt.Printf("\033[4%dm %s \033[0m\n", colorCode, text)
}

type MyFormatter struct {
	Prefix string
	TimeFormat string
}

func (f *MyFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	var color int

	// 设置颜色
	switch entry.Level {
	case logrus.ErrorLevel:
		color = ccRedBlock
	case logrus.WarnLevel:
		color = ccYellowBlock
	case logrus.InfoLevel:
		color = ccGreenBlock
	case logrus.DebugLevel:
		color = ccGrayBlock
	default:
		color = 8
	}

	// 设置buffer(缓冲区)
	var b *bytes.Buffer
	if entry.Buffer == nil {
		b = &bytes.Buffer{}
	} else {
		b = entry.Buffer
	}
	// 时间格式化
	formatTime := entry.Time.Format(f.TimeFormat)

	// 对应文件的行号
	// fileVal := fmt.Sprintf("%s:%d", entry.Caller.File, entry.Caller.Line)
	fileVal := fmt.Sprintf("%s:%d", path.Base(entry.Caller.File), entry.Caller.Line)

	// 设置格式
	fmt.Fprintf(b, "[%s] \033[3%dm[%s]\033[0m [%s] %s %s\n", f.Prefix, color, entry.Level, formatTime, fileVal, entry.Message)

	return b.Bytes(), nil
}

func main() {
	logrus.SetReportCaller(true)
	logrus.SetLevel(logrus.DebugLevel)
	logrus.SetFormatter(&MyFormatter{
		Prefix: "GIN-LOGRUS",
		TimeFormat: "2006-01-02 15:04:06",
	})

	logrus.Infof("你好...")
	logrus.Warnln("你好...")
	logrus.Errorln("你好...")
	logrus.Debugf("你好...")
}
