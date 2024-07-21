// server.go
package main

import (
	"encoding/json"
	"io"
	"log"
	"net"
	"net/http"
	"os"
	"runtime"
	"time"
)

const MemoryLimit = 20 * 1024 * 1024 // 20 MB memory allocation
const Timeout = 0                    // change the timeout here

var startTime time.Time

func getMemoryUsage() uint64 {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	return m.Alloc
}

func checkMemoryUsage() {
	alloc := getMemoryUsage()
	log.Printf("INFO: Current memory usage: %d bytes (%.2f MB), memory allocation: 20 MB\n", alloc, float64(alloc)/1024/1024)

	if alloc > MemoryLimit {
		elapsed := time.Since(startTime).Seconds()
		log.Printf("ERROR: Memory usage exceeded limit: %d bytes (%.2f MB) with running time %f\n", alloc, float64(alloc)/1024/1024, elapsed)
		os.Exit(1)
	}
}

func callPartnerService() (string, error) {
	url := "http://localhost:8081/data"
	client := &http.Client{
		Timeout: Timeout * time.Second, // Set a timeout
	}

	response, err := client.Get(url)
	if err != nil {
		return "", err
	}
	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return "", err
	}

	return string(body), nil
}

func DataHandler(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	data, err := callPartnerService()
	alloc := getMemoryUsage()
	if err != nil {
		errMesage := err.Error()

		if netErr, ok := err.(net.Error); ok && netErr.Timeout() {
			errMesage = "timed out"
		}

		log.Printf("ERROR: method=%s url=%s error=%s latency=%s memory_usage=(%.2f MB)\n", r.Method, r.URL.String(), errMesage, time.Since(start), float64(alloc)/1024/1024)
		http.Error(w, errMesage, http.StatusInternalServerError)
		return
	}

	jsonResponse := map[string]string{"response": data}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(jsonResponse)

	log.Printf("INFO: method=%s url=%s latency=%s status=%d memory_usage(%.2f MB)\n", r.Method, r.URL.String(), time.Since(start), http.StatusOK, float64(alloc)/1024/1024)
}

func main() {
	startTime = time.Now()
	ticker := time.NewTicker(1 * time.Second)
	go func() {
		for range ticker.C {
			checkMemoryUsage()
		}
	}()

	http.HandleFunc("/data", DataHandler)
	log.Println("INFO: Server listening on port 8080...")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatalf("FATAL: Server failed to start: %s\n", err)
	}
}
