package main

import (
	"io"
	"log"
	"net/http"
	"time"
)

func callPartnerService() {
	log.Printf("INFO: Request to server")
	url := "http://localhost:8080/data"
	start := time.Now()
	response, err := http.Get(url)
	if err != nil {
		log.Printf("ERROR: Response error when request to server, latency=%s, message=%s, ", time.Since(start), err.Error())
		return
	}
	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	if err != nil {
		log.Printf("ERROR: Response error reading response body: %s", err.Error())
		return
	}

	log.Printf("INFO: Response from server, latency=%s, message=%s", time.Since(start), string(body))
}

func main() {
	for {
		go callPartnerService()
		time.Sleep(10 * time.Millisecond)
	}
}
