// backend-service.go

package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/create", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "创建活动成功！- activity service 2")
	})

	http.HandleFunc("/update", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "修改活动成功！- activity service 2")
	})

	log.Println("activity services are running on port 8083...")
	log.Fatal(http.ListenAndServe(":8083", nil))
}
