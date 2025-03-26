package main

import (
	"fmt"
	"github.com/joho/godotenv"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"sync"
	"time"
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

	var wg sync.WaitGroup
	numWorkers := 2
	requestsPerWorker := 5
	totalRequests := 100
	serverURL := "http://localhost:" + port

	requests := make(chan int, totalRequests)
	requestsCollector := make([][]int, numWorkers)
	for i := range requestsCollector {
		requestsCollector[i] = make([]int, requestsPerWorker)
	}
	for i := 0; i < totalRequests; i++ {
		requests <- i
	}

	rateLimit := time.Second / 5
	throttle := time.Tick(rateLimit)

	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go func() {
			emptyQueue := false
			for !emptyQueue {
				for j := 0; j < requestsPerWorker; j++ {
					val, ok := <-requests
					if ok {
						requestsCollector[i][j] = val
					} else {
						emptyQueue = true
						break
					}
				}
				if !emptyQueue {
					worker(i, serverURL, requestsCollector[i], throttle)
				}
			}
			wg.Done()
		}()
	}
	close(requests)
	wg.Wait()

	fmt.Println("Client finished sending requests.")
}

func worker(id int, serverURL string, batch []int, throttle <-chan time.Time) {
	client := &http.Client{}

	for i := range batch {
		reqNum := batch[i]
		url := serverURL

		<-throttle
		resp, err := client.Post(url, "application/json", nil)
		if err != nil {
			log.Printf("Worker %d: Request %d failed: %v Time: %s", id+2, reqNum, err, time.Now().Format("15:04:05"))
			continue
		}
		resp.Body.Close()

		log.Printf("Worker %d: Request %d complete. StatusCode: %d Time: %s", id+2, reqNum, resp.StatusCode, time.Now().Format("15:04:05"))
		//		time.Sleep(time.Millisecond * 100)
	}
	log.Printf("Worker %d finished block", id+2)
}
