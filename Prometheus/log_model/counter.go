package log_model

import (
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

func RecordMetrics() {
	go func() {
		for {
			opsProcessed.Inc()
			// 管道总数
			ChanTotal.Set(200)
			// 管道中当前的日志数量
			ChanCurrent.Set(float64(len(logChan)))
			// 每隔2秒钟记录一次
			time.Sleep(3 * time.Second)
		}
	}()
}

var (
	opsProcessed = promauto.NewCounter(prometheus.CounterOpts{
		Name: "myapp_processed_ops_total",
		Help: "The total number of processed events",
	})

	ReqSuccuess = promauto.NewCounter(prometheus.CounterOpts{
		Name: "myapp_request_success_total",
		Help: "The total number of processed events",
	})
	ReqError = promauto.NewCounter(prometheus.CounterOpts{
		Name: "myapp_request_error_total",
		Help: "The total number of processed events",
	})

	// 管道总数
	ChanTotal = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "myapp_chan_total",
		Help: "The total number of processed events",
	})
	// 管道中当前的日志数量
	ChanCurrent = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "myapp_chan_current",
		Help: "The total number of processed events",
	})
)
