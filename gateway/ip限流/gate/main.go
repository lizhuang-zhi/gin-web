// gateway.go

package main

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"net/http/httputil"
	"net/url"
	"sync"

	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
)

// IP限流器的存储
var (
	ipsLimiter = make(map[string]*rate.Limiter)
	mtx        sync.Mutex
)

const (
	rps   = 1 // 每秒请求次数
	burst = 2 // 桶的容量
)

// 获取IP地址限流器
func getLimiter(ip string) *rate.Limiter {
	mtx.Lock()
	defer mtx.Unlock()

	if limiter, exists := ipsLimiter[ip]; exists {
		return limiter
	}

	limiter := rate.NewLimiter(rate.Limit(rps), burst)
	ipsLimiter[ip] = limiter
	return limiter
}

// 服务映射
var GatewayMap = map[string]string{
	"activity": "http://localhost:8081",
	"mail":     "http://localhost:8082",
}

func main() {
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

	log.Printf("targetUrl is %s\n", targetUrl) // targetUrl is http://localhost:8081

	// 创建反向代理
	proxy := httputil.NewSingleHostReverseProxy(targetUrl)

	return func(c *gin.Context) {
		ip, _, err := net.SplitHostPort(c.Request.RemoteAddr)
		if err != nil {
			log.Printf("获取客户端IP失败: %v\n", err)

			c.JSON(http.StatusInternalServerError, gin.H{
				"code":    http.StatusInternalServerError,
				"message": "获取客户端IP失败",
			})
			c.Abort()
			return
		}

		fmt.Println("client ip is: ", ip)

		limiter := getLimiter(ip)
		if !limiter.Allow() {
			// 拦截请求
			c.JSON(http.StatusTooManyRequests, gin.H{
				"code":    http.StatusTooManyRequests,
				"message": "请求过于频繁，请稍后再试",
			})
			c.Abort()
			return
		}

		c.Request.URL.Path = c.Param("proxyPath")

		log.Printf("proxyPath is %s\n", c.Param("proxyPath")) // proxyPath is /create[/update, /send, /accept]

		proxy.ServeHTTP(c.Writer, c.Request)
	}
}
