// gateway.go

package main

import (
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
)

func main() {
	// 定义目标服务的URL
	service1URL, _ := url.Parse("http://localhost:8081")
	service2URL, _ := url.Parse("http://localhost:8082")

	// 创建两个反向代理
	proxy1 := httputil.NewSingleHostReverseProxy(service1URL)
	proxy2 := httputil.NewSingleHostReverseProxy(service2URL)

	// 路由和处理函数
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path

		// 根据路径转发到不同服务
		if path == "/service1" {
			proxy1.ServeHTTP(w, r)
		} else if path == "/service2" {
			proxy2.ServeHTTP(w, r)
		} else {
			http.Error(w, "Service not found", http.StatusNotFound)
		}
	})

	log.Println("Gateway is running on port 8080...")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
