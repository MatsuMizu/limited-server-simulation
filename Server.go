package main

import (
	"encoding/json"
	"fmt"
	"github.com/joho/godotenv"
	"golang.org/x/time/rate"
	"log"
	"math/rand"
	"net/http"
	"os"
	"path/filepath"
	"sync"
	"time"
)

type Stats struct {
	Total   int `json:"total"`
	Client1 int `json:"client1"`
	Client2 int `json:"client2"`
}

var (
	limiter      = rate.NewLimiter(rate.Limit(5), 5)
	mu           sync.Mutex
	handledStats Stats
	requestCount int
)

func main() {
	environmentPath := filepath.Join("C:\\", "Users", "Matsuri", "Downloads", "WBTypes6", "Server", ".env")
	err := godotenv.Load(environmentPath)
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	port := os.Getenv("PORT")
	if port == "" {
		log.Fatal("PORT environment variable is not set")
	}

	rand.Seed(time.Now().UnixNano())
	fmt.Printf("Server listening on port %s\n", port)
	requestCount = 0
	handledStats = Stats{
		Total:   0,
		Client1: 0,
		Client2: 0,
	}

	http.HandleFunc("/", Handler)
	http.HandleFunc("/stats", statsHandler)
	err = http.ListenAndServe(":"+port, nil)
	if err != nil {
		log.Fatal(err)
	}
}

func statsHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		modifiedJson, err := json.Marshal(handledStats)
		if err != nil {
			http.Error(w, "Error marshaling JSON: "+err.Error(), http.StatusInternalServerError)
			return
		}
		if _, err := w.Write(modifiedJson); err != nil {
			http.Error(w, "Error writing response: "+err.Error(), http.StatusInternalServerError)
			return
		}
		log.Printf("GET request stats:    Total: %d    Client1: %d    Client2: %d\n", handledStats.Total, handledStats.Client1, handledStats.Client2)
	default:
		http.Error(w, "Invalid Method", http.StatusMethodNotAllowed)
	}
}

func Handler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		getHandler(w, r)
	case http.MethodPost:
		postHandler(w, r)
	default:
		http.Error(w, "Invalid Method", http.StatusMethodNotAllowed)
	}
}

func getHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	log.Printf("Showed signs of life")
}

func postHandler(w http.ResponseWriter, r *http.Request) {
	mu.Lock()
	defer mu.Unlock()
	requestCount += 1
	if !limiter.Allow() {
		uniqueApology := fmt.Sprintf("Limit exceeded! POST request #%d cannot be handled by server.", requestCount)
		http.Error(w, uniqueApology, http.StatusTooManyRequests)
		return
	}
	requestID := r.Header.Get("X-User-ID")
	if requestID == "" {
		fmt.Println("X-User-ID header not found")
		return
	}
	switch requestID {
	case "0":
		handledStats.Client1 += 1
	case "1":
		handledStats.Client2 += 1
	}
	handledStats.Total += 1
	ansCode := rand.Intn(20)
	var response int
	if ansCode >= 0 && ansCode <= 6 {
		response = http.StatusOK
	} else if ansCode >= 7 && ansCode <= 13 {
		response = http.StatusAccepted
	} else if ansCode >= 14 && ansCode <= 16 {
		response = http.StatusBadRequest
	} else {
		response = http.StatusInternalServerError
	}

	w.WriteHeader(response)
	log.Printf("Handled from #1: %d   Handled from #2: %d", handledStats.Client1, handledStats.Client2)
}
