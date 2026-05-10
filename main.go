package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"time"
)

// FuelData represents the structure of our cached data
type FuelData struct {
	Region    string  `json:"region"`
	Price     float64 `json:"price"`
	Currency  string  `json:"currency"`
	UpdatedAt string  `json:"updated_at"`
}

// Global Cache (Safe for concurrent use)
var (
	cache      = make(map[string]FuelData)
	cacheMutex sync.RWMutex
)

func main() {
	// Örnek veriyi önbelleğe ekle
	seedCache()

	http.HandleFunc("/api/v1/fuel", getFuelPrice)

	fmt.Println("🚀 Geo-Fuel API Gateway starting on :8080...")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		fmt.Printf("Error starting server: %s\n", err)
	}
}

func getFuelPrice(w http.ResponseWriter, r *http.Request) {
	region := r.URL.Query().Get("region")
	if region == "" {
		region = "global"
	}

	// Thread-safe read from cache
	cacheMutex.RLock()
	data, exists := cache[region]
	cacheMutex.RUnlock()

	if !exists {
		http.Error(w, "Region data not found in cache", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data)
}

func seedCache() {
	cache["tr-ist"] = FuelData{
		Region:    "Istanbul, TR",
		Price:     42.55,
		Currency:  "TRY",
		UpdatedAt: time.Now().Format(time.RFC3339),
	}
    cache["de-ber"] = FuelData{
		Region:    "Berlin, DE",
		Price:     1.78,
		Currency:  "EUR",
		UpdatedAt: time.Now().Format(time.RFC3339),
	}
}
