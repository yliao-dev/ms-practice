package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
)

var flags = map[string]bool {
	"new-checkout": true,
	"dark-mode": false,
}


type FeatureResponse struct {
	Feature string 	`json:"feature"`
	Enabled bool `json:"enabled"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}

func FeatureFlagHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", "application/json")
	featureName := strings.TrimPrefix(r.URL.Path, "/features/")
	enabled, exist := flags[featureName]
	if !exist {
		w.WriteHeader(http.StatusNotFound)
		response := ErrorResponse{Error: fmt.Sprintf("feature not found: %s", featureName)}
		json.NewEncoder(w).Encode(response)
		return
	}
	
	response := FeatureResponse{Feature: featureName, Enabled: enabled}
	json.NewEncoder(w).Encode(response)

}

func main() {
	http.HandleFunc("/feature/", FeatureFlagHandler)
	fmt.Println("Feature service serving port on 80801")
	log.Fatal(http.ListenAndServe(":8081", nil))
}