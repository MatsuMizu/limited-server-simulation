package main

import (
	"context"
	"github.com/joho/godotenv"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

func isServerAlive(url string, timeout time.Duration) bool {
	client := &http.Client{
		Timeout: timeout,
	}

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		log.Println("Error creating request:", err)
		return false
	}

	resp, err := client.Do(req)
	if err != nil {
		log.Println("Error making request:", err)
		return false
	}
	defer func() {
		err := resp.Body.Close()
		if err != nil {
			log.Println("Error closing response body:", err)
		}
	}()

	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		return true
	} else {
		log.Printf("Server returned status code: %d\n", resp.StatusCode)
		return false
	}
}

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
	serverURL := "http://localhost:" + port

	timeout := 5 * time.Second
	checkInterval := 5 * time.Second

	ticker := time.NewTicker(checkInterval)
	defer ticker.Stop()

	for range ticker.C {
		if isAlive := isServerAlive(serverURL, timeout); isAlive {
			log.Printf("Server is alive!")
		} else {
			log.Printf("Server is dead.")
		}
	}
}
