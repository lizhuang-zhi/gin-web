// main.go
package main

import (
	"fmt"
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var (
	upgrader = websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true // 允许跨域
		},
	}
)

func main() {
	r := gin.Default()

	// WebSocket连接管理
	var connections sync.Map

	r.GET("/ws", func(c *gin.Context) {
		// 升级get请求为webSocket协议
		conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
		if err != nil {
			fmt.Println("Upgrade error:", err)
			return
		}
		defer conn.Close()

		// 将连接保存到connections Map中
		connections.Store(conn, true)

		for {
			// 读取客户端发送的消息
			_, msg, err := conn.ReadMessage()
			if err != nil {
				fmt.Println("Read message error:", err)
				break
			}
			fmt.Printf("Received message: %s\n", msg)

			// 向客户端发送消息
			if err := conn.WriteMessage(websocket.TextMessage, []byte("Hello from server")); err != nil {
				fmt.Println("Write message error:", err)
				break
			}
		}

		// 从connections Map中删除连接
		connections.Delete(conn)
	})

	// HTTP请求处理
	r.GET("/test/push", func(c *gin.Context) {
		// 业务计算
		sum := 0
		for i := 1; i <= 100; i++ {
			sum += i
		}

		// 将计算结果转换为字符串
		result := fmt.Sprintf("从1到100的总和为：%d", sum)

		// 向所有连接的WebSocket客户端推送结果
		connections.Range(func(key, value interface{}) bool {
			conn := key.(*websocket.Conn)
			err := conn.WriteMessage(websocket.TextMessage, []byte(result))
			if err != nil {
				fmt.Println("发送 WebSocket 消息出错:", err)
			}
			return true
		})

		// 返回计算结果给HTTP请求的响应
		c.String(http.StatusOK, result)
	})
	r.Run(":8798")
}
