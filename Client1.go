package main

import (
	"fmt"
	"github.com/joho/godotenv"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"sync"
	"time"
)

func main() {
	numClients := 2
	numStatus := 5
	statsCollector := make([][]int, numClients)
	for i := range statsCollector {
		statsCollector[i] = make([]int, numStatus)
	}
	var clientWg sync.WaitGroup
	for i := 0; i < numClients; i++ {
		clientWg.Add(1)
		go client(i, &clientWg, &statsCollector)
	}
	clientWg.Wait()
	fmt.Println("\n\n----------------------\nSTATS\n----------------------\n\n")
	for i := 0; i < numClients; i++ {
		fmt.Printf("Client #%d:\n", i+1)
		sum := 0
		for j := 0; j < numStatus; j++ {
			sum += statsCollector[i][j]
		}
		fmt.Printf("Requests sent: %d\n", sum)
		fmt.Printf("Statuses:")
		fmt.Printf("    %d: %d", http.StatusOK, statsCollector[i][0])
		fmt.Printf("    %d: %d", http.StatusAccepted, statsCollector[i][1])
		fmt.Printf("    %d: %d", http.StatusBadRequest, statsCollector[i][2])
		fmt.Printf("    %d: %d", http.StatusInternalServerError, statsCollector[i][3])
		fmt.Printf("    %d: %d\n\n", http.StatusTooManyRequests, statsCollector[i][4])
	}
}
func client(clientId int, clientWg *sync.WaitGroup, statsCollector *[][]int) {
	defer clientWg.Done()
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
					worker(clientId, i, serverURL, requestsCollector[i], throttle, statsCollector)
				}
			}
			wg.Done()
		}()
	}
	close(requests)
	wg.Wait()

	fmt.Println("Client finished sending requests.")
}

func worker(clientId int, workerId int, serverURL string, batch []int, throttle <-chan time.Time, statsCollector *[][]int) {
	client := &http.Client{}

	for i := range batch {
		reqNum := batch[i]
		url := serverURL

		<-throttle
		req, err := http.NewRequest("POST", url, nil)
		if err != nil {
			log.Printf("Worker %d: Request %d creation failed: %v", workerId+2*clientId, reqNum, err)
			continue
		}
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("X-User-ID", strconv.Itoa(clientId))
		resp, err := client.Do(req)
		if err != nil {
			log.Printf("Client %d worker %d: Request %d sending failed: %v", clientId, workerId, reqNum, err)
			continue
		}
		resp.Body.Close()

		status := resp.StatusCode
		switch status {
		case http.StatusOK:
			(*statsCollector)[clientId][0] += 1
		case http.StatusAccepted:
			(*statsCollector)[clientId][1] += 1
		case http.StatusBadRequest:
			(*statsCollector)[clientId][2] += 1
		case http.StatusInternalServerError:
			(*statsCollector)[clientId][3] += 1
		case http.StatusTooManyRequests:
			(*statsCollector)[clientId][4] += 1
		}
		log.Printf("Client %d worker %d: Request %d complete. StatusCode: %d", clientId, workerId, reqNum, status)
	}
	log.Printf("Client %d worker %d finished block", clientId, workerId)
}
