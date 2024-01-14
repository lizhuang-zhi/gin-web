package main

import (
	"booking-app/Prometheus/log_model"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func main() {
	log_model.RecordMetrics() // 记录指标

	router := gin.Default()

	router.GET("/success", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "success",
		})
		// 构建成功日志
		log_model.Producer(&log_model.LogInfo{
			Content: "success",
			Level:   "info",
		})
	})

	router.GET("/error", func(c *gin.Context) {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "error",
		})
		// 构建失败日志
		log_model.Producer(&log_model.LogInfo{
			Content: "error",
			Level:   "error",
		})
	})

	// 开启消费日志
	go log_model.Consumer()

	// Create a new ServeMux.
	mux := http.NewServeMux()

	// Register the Prometheus handler.
	mux.Handle("/metrics", promhttp.Handler())

	// Register the gin router with the ServeMux.
	mux.Handle("/", router)

	// Start the HTTP server using the ServeMux.
	http.ListenAndServe(":2112", mux)
}
