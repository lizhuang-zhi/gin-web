// backend-service.go

package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/send", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "发送邮件成功！- mail service 2")
	})

	http.HandleFunc("/accept", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "接收邮件成功！- mail service 2")
	})

	log.Println("mail services are running on port 8084...")
	log.Fatal(http.ListenAndServe(":8084", nil))
}
