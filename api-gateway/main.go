package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"time"
)

const featureServiceURL = "http://localhost:8081"

func GatewayHandler(w http.ResponseWriter, r *http.Request){
	featureName := strings.TrimPrefix(r.URL.Path, "/v1/features/")
	backendReq, err := http.NewRequest(http.MethodGet, fmt.Sprintf("%s/features/%s", featureServiceURL, featureName), nil)
	if err != nil {
		http.Error(w, "failed to create backend request", http.StatusInternalServerError)
		return
	}
	client := &http.Client{Timeout: 5 * time.Second}
	backendRes, err := client.Do(backendReq)
	if err != nil {
		http.Error(w, "failed to fetch backend service", http.StatusServiceUnavailable)
		return
	}
	defer backendRes.Body.Close()
	w.Header().Set("Content-type", backendRes.Header.Get("Content-Type"))
	w.WriteHeader(backendRes.StatusCode)
	io.Copy(w, backendRes.Body)

	
}

func main() {
	http.HandleFunc("/v1/features/", GatewayHandler)
	fmt.Println("API Gateway listening on port 8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}