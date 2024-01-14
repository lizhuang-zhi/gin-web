package log_model

import (
	"time"
)

// 日志生产消费模型
var (
	logChan  chan *LogInfo = make(chan *LogInfo, 200)
	exitChan chan bool     = make(chan bool) // 退出管道
)

// 日志
type LogInfo struct {
	// 日志内容
	Content string
	// 日志级别
	Level string
}

// 生产日志
func Producer(log *LogInfo) {
	logChan <- log
}

// 消费日志
func Consumer() {
	ticker := time.NewTicker(time.Second * 3)
	for {
		select {
		case <-ticker.C:
			// 消费日志
			// 从日志管道中获取日志, 然后判断成功还是失败,记录到Prometheus中
			ConsumeLogAndRecordToPrometheus()
		case <-exitChan:
			// 退出
			ticker.Stop() // 停止定时器
			return
		}
	}
}

// 消费日志并记录到Prometheus中
func ConsumeLogAndRecordToPrometheus() {
	logInfo := <-logChan
	if logInfo.Level == "error" {
		// 记录到Prometheus中
		ReqError.Inc()
	} else {
		// 记录到Prometheus中
		ReqSuccuess.Inc()
	}
	// 记录管道中的日志数量
	ChanCurrent.Set(float64(len(logChan)))
}
