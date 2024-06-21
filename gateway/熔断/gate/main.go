// gateway.go

package main

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"

	"github.com/afex/hystrix-go/hystrix"
	"github.com/gin-gonic/gin"
)

// 服务映射
var GatewayMap = map[string]string{
	"activity": "http://localhost:8081",
	"mail":     "http://localhost:8082",
}

func main() {
	// 配置熔断器
	hystrix.ConfigureCommand("my_command", hystrix.CommandConfig{
		Timeout:                1000, // 命令超时时间设置为1秒
		MaxConcurrentRequests:  100,  // 最大并发请求数设置为100
		RequestVolumeThreshold: 10,   // 熔断器请求里的阈值，用于确定是否需要熔断
		ErrorPercentThreshold:  50,   // 错误比率阈值，超过则触发熔断
		SleepWindow:            5000, // 熔断后进行半开状态尝试恢复的时间窗口
	})

	r := gin.Default()

	// 遍历服务映射，为每个服务创建反向代理
	for service, serviceUrl := range GatewayMap {
		relativePath := fmt.Sprintf("/%s/*proxyPath", service)
		reverseProxy := createReverseProxy(serviceUrl)
		r.Any(relativePath, reverseProxy)
	}

	log.Println("Gateway is running on port 8080...")
	r.Run(":8080")
}

func createReverseProxy(target string) gin.HandlerFunc {
	targetUrl, err := url.Parse(target)
	if err != nil {
		log.Fatal(err)
	}

	proxy := httputil.NewSingleHostReverseProxy(targetUrl)

	return func(c *gin.Context) {
		var proxyError error

		customWriter := &responseCaptureWriter{ResponseWriter: c.Writer, statusCode: http.StatusOK}

		err := hystrix.Do("my_command", func() error {
			c.Request.URL.Path = c.Param("proxyPath")
			proxy.ServeHTTP(customWriter, c.Request)
			if customWriter.statusCode >= 400 {
				// 如果状态码表示错误，设置proxyError以在降级逻辑中使用
				proxyError = fmt.Errorf("upstream service returned with status code: %d", customWriter.statusCode)
			}
			return proxyError
		}, func(err error) error {
			/*
				正式项目中：
				1. 降级页面：返回一个简化的页面，告诉用户服务不可用“系统维护中，请稍后再试”
				2. 缓存数据：如果有缓存数据，可以返回缓存数据，提高用户体验
				3. 简化处理：简化处理流程，跳过非核心的计算与逻辑，返回一个简化的结果
			*/
			c.JSON(http.StatusInternalServerError, gin.H{
				"code":    http.StatusInternalServerError,
				"message": "服务熔断，执行降级操作",
			})
			c.Abort() // 终止请求
			return nil
		})

		if err != nil {
			fmt.Println("Hystrix错误或代理错误：", err)
			return
		}
	}
}

// responseCaptureWriter 用于捕获状态码和允许后续读取
type responseCaptureWriter struct {
	gin.ResponseWriter
	statusCode int
}

func (w *responseCaptureWriter) WriteHeader(code int) {
	w.statusCode = code
	w.ResponseWriter.WriteHeader(code)
}
