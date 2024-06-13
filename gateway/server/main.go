// backend-service.go

package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {
	// 服务1
	http.HandleFunc("/service1", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "This is service 1")
	})

	// 服务2
	http.HandleFunc("/service2", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "This is service 2")
	})

	log.Println("Backend services are running on port 8081 and 8082...")
	go func() {
		log.Fatal(http.ListenAndServe(":8081", nil))
	}()
	log.Fatal(http.ListenAndServe(":8082", nil))
}
