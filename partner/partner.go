// partner.go
package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"golang.org/x/exp/rand"
)

func PartnerHandler(w http.ResponseWriter, r *http.Request) {
	delay := time.Duration(rand.Intn(11)) // random delay
	log.Printf("INFO: Request from server, delay_set=%d, time=%v", delay, time.Now().Format("2006-01-02 15:04:05.000000"))

	time.Sleep(delay * time.Second) // Simulate a delay
	w.Write([]byte("Hello from Partner Service"))
}

func main() {
	http.HandleFunc("/data", PartnerHandler)
	fmt.Println("Partner service listening on port 8081...")
	http.ListenAndServe(":8081", nil)
}
