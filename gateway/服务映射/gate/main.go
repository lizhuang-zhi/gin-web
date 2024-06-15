// gateway.go

package main

import (
	"fmt"
	"log"
	"net/http/httputil"
	"net/url"

	"github.com/gin-gonic/gin"
)

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
		c.Request.URL.Path = c.Param("proxyPath")

		log.Printf("proxyPath is %s\n", c.Param("proxyPath")) // proxyPath is /create[/update, /send, /accept]

		proxy.ServeHTTP(c.Writer, c.Request)
	}
}
