// gateway.go

package main

import (
	"log"
	"net/http/httputil"
	"net/url"
	"sync/atomic"

	"github.com/gin-gonic/gin"
)

// 服务实例
type ServiceInstance struct {
	URL    *url.URL
	Active bool // 标记一个实例是否可用
}

// LoadBalancer 实现轮询负载均衡
type LoadBalancer struct {
	Instances []*ServiceInstance
	counter   uint64
}

func NewLoadBalancer(urls []string) *LoadBalancer {
	var instances []*ServiceInstance
	for _, urlString := range urls {
		serviceUrl, err := url.Parse(urlString)
		if err != nil {
			log.Fatal(err)
		}

		instances = append(instances, &ServiceInstance{
			URL:    serviceUrl,
			Active: true, // 假设所有实例起初都是活跃的
		})
	}

	return &LoadBalancer{
		Instances: instances,
	}
}

// GetInstance 通过轮询策略获取一个服务实例
func (l *LoadBalancer) GetInstance() *ServiceInstance {
	instanceIndex := atomic.AddUint64(&l.counter, 1) % uint64(len(l.Instances))
	return l.Instances[instanceIndex]
}

// ReverseProxyHandler 创建反向代理处理器并使用负载均衡器来选择服务实例
func ReverseProxyHandler(lb *LoadBalancer) gin.HandlerFunc {
	return func(c *gin.Context) {
		instance := lb.GetInstance()
		if !instance.Active {
			c.JSON(500, gin.H{"message": "No available instance"})
			return
		}

		proxy := httputil.NewSingleHostReverseProxy(instance.URL)
		c.Request.URL.Path = c.Param("proxyPath")
		proxy.ServeHTTP(c.Writer, c.Request)
	}
}

// 活动服务的实例
var ActivityServiceInstances = []string{
	"http://localhost:8081",
	"http://localhost:8083",
}

// 邮件服务的实例
var MailServiceInstances = []string{
	"http://localhost:8082",
	"http://localhost:8084",
}

func main() {
	r := gin.Default()

	// 活动服务的负载均衡
	activityLb := NewLoadBalancer(ActivityServiceInstances)
	r.Any("/activity/*proxyPath", ReverseProxyHandler(activityLb))

	// 邮件服务的负载均衡
	mailLb := NewLoadBalancer(MailServiceInstances)
	r.Any("/mail/*proxyPath", ReverseProxyHandler(mailLb))

	log.Println("Gateway is running on port 7080...")
	r.Run(":7080")
}
