// backend-service.go

package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/send", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "发送邮件成功！")
	})

	http.HandleFunc("/accept", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "接收邮件成功！")
	})

	log.Println("mail services are running on port 8082...")
	log.Fatal(http.ListenAndServe(":8082", nil))
}
