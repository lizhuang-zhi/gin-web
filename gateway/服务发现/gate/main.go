// gateway.go

package main

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"sync/atomic"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/hashicorp/consul/api"
)

// 服务实例
type ServiceInstance struct {
	URL     *url.URL
	Active  bool   // 标记一个实例是否可用
	Headthy uint32 // 使用原子操作来跟踪健康状态
}

// LoadBalancer 实现轮询负载均衡
type LoadBalancer struct {
	Instances []*ServiceInstance
	counter   uint64
}

func (lb *LoadBalancer) GetNextActiveInstance() *ServiceInstance {
	start := atomic.AddUint64(&lb.counter, 1)
	count := uint64(len(lb.Instances))

	// 尝试找到一个健康的服务实例
	for i := uint64(0); i < count; i++ {
		idx := (start + i) % count
		instance := lb.Instances[idx]
		if atomic.LoadUint32(&instance.Headthy) == 1 { // 使用原子操作来获取健康状态
			return instance
		}
	}

	return nil
}

func (lb *LoadBalancer) SetInstanceHealthy(instance *ServiceInstance, healthy bool) {
	var state uint32
	if healthy {
		state = 1
	}
	atomic.StoreUint32(&instance.Headthy, state) // 使用原子操作来更新健康状态
}

func (lb *LoadBalancer) CheckHealth(instance *ServiceInstance, interval time.Duration) {
	for {
		time.Sleep(interval) // 每隔一段时间检查一次健康状态

		// 简单的健康检查
		resp, err := http.Get(instance.URL.String())
		if err != nil || resp == nil { // 实例宕机
			lb.SetInstanceHealthy(instance, false)
			continue
		}
		lb.SetInstanceHealthy(instance, true)
	}
}

func NewLoadBalancer(urls []string) *LoadBalancer {
	var instances []*ServiceInstance
	for _, urlString := range urls {
		serviceUrl, err := url.Parse(urlString)
		if err != nil {
			log.Fatal(err)
		}

		instances = append(instances, &ServiceInstance{
			URL:     serviceUrl,
			Active:  true, // 假设所有实例起初都是活跃的
			Headthy: 1,    // 假设所有实例起初都是健康的
		})
	}

	return &LoadBalancer{
		Instances: instances,
	}
}

// ReverseProxyHandler 创建反向代理处理器并使用负载均衡器来选择服务实例
func ReverseProxyHandler(lb *LoadBalancer) gin.HandlerFunc {
	return func(c *gin.Context) {
		instance := lb.GetNextActiveInstance()
		if instance == nil {
			c.JSON(http.StatusServiceUnavailable, gin.H{"error": "No active instances"})
			return
		}

		proxy := httputil.NewSingleHostReverseProxy(instance.URL)
		c.Request.URL.Path = c.Param("proxyPath")

		// 修改请求进行错误处理，重新路由请求如果当前实例宕机
		proxy.ModifyResponse = func(response *http.Response) error {
			if response.StatusCode >= 500 { // 实例宕机或者其他错误
				lb.SetInstanceHealthy(instance, false)
				return errors.New("Server error")
			}
			return nil
		}

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

// DiscoverServiceInstances 使用Consul服务发现查询特定服务的所有健康实例
func DiscoverServiceInstances(consulClient *api.Client, serviceName string) ([]*ServiceInstance, error) {
	// Consul健康API, 查询具体服务的健康实例
	services, _, err := consulClient.Health().Service(serviceName, "", true, nil)
	if err != nil {
		return nil, err
	}

	var instances []*ServiceInstance
	for _, service := range services {
		// 创建服务实例并将其添加到实例列表中
		instanceURL, err := url.Parse(fmt.Sprintf("http://%s:%d", service.Service.Address, service.Service.Port))
		if err != nil {
			return nil, err
		}
		instances = append(instances, &ServiceInstance{
			URL:     instanceURL,
			Active:  true, // 假设所有实例起初都是活跃的
			Headthy: 1,    // 假设所有实例起初都是健康的
		})
	}

	return instances, nil
}

func main() {
	// 创建Consul客户端
	consulConfig := api.DefaultConfig()
	consulClient, err := api.NewClient(consulConfig)
	if err != nil {
		log.Fatal("Consul client err: ", err)
	}

	r := gin.Default()

	// 动态服务发现与健康检查
	/*
		这里省略了外部操作, 需要将服务注册到Consul中, 例如:
		1. 本地启动Consul服务: consul agent -dev(可通过http://localhost:8500/ui/dc1/services查看服务注册UI情况)
		2. 注册服务:
			cd gateway/服务发现/consul/json
			consul services register activity-1.json
			consul services register activity-2.json
		3. 运行命令查询Consul中注册的服务:
			consul catalog services

	*/
	activityServiceName := "activity" // 活动服务在Consul中注册的名称: 具体查看consul/json中的"Name": "activity",
	activityInstances, err := DiscoverServiceInstances(consulClient, activityServiceName)
	// 活动服务的负载均衡
	activityLb := &LoadBalancer{Instances: activityInstances} // 动态服务实例(通过Consul发现)
	// activityLb := NewLoadBalancer(ActivityServiceInstances)  // 静态服务实例
	for _, instance := range activityLb.Instances {
		fmt.Println("instance.Headthy:", instance.Headthy)
		go activityLb.CheckHealth(instance, 5*time.Second) // 每隔10秒检查一次健康状态
	}
	r.Any("/activity/*proxyPath", ReverseProxyHandler(activityLb))

	// // 邮件服务的负载均衡(未使用Consul服务发现)
	// // 1. 书写/health接口, 用于检查服务健康状态
	// // 2. 通过Consul书写json文件并注册服务
	// mailLb := NewLoadBalancer(MailServiceInstances)
	// for _, instance := range mailLb.Instances {
	// 	go mailLb.CheckHealth(instance, 10*time.Second) // 每隔10秒检查一次健康状态
	// }
	// r.Any("/mail/*proxyPath", ReverseProxyHandler(mailLb))

	log.Println("Gateway is running on port 7080...")
	r.Run(":7080")
}
