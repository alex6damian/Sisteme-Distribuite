package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"os"
	"sync"
	"time"
)

// Constanta pentru numărul de clienți concurenți
const NUM_CLIENTS = 7

// Structurile JSON, trebuie să se potrivească cu cele din server și din fișier
type GenericRequest struct {
	TaskNumber int             `json:"task"`
	Input      json.RawMessage `json:"input"`
	ClientID   int             `json:"client_id"`
}

type GenericResponse struct {
	Status string          `json:"status"`
	Result json.RawMessage `json:"result,omitempty"`
	Error  string          `json:"error,omitempty"`
}

// running client instance
func runSingleClient(clientID int, wg *sync.WaitGroup, requestToSend GenericRequest) {
	// waitgroup done at the end
	defer wg.Done()

	// server connection
	conn, err := net.DialTimeout("tcp", "localhost:8080", 2*time.Second)
	if err != nil {
		fmt.Printf("[Client %d] Error: %v", clientID, err)
		return
	}
	defer conn.Close()

	// reading and ignoring the welcome message
	// reason: we need to clear the buffer before sending our request
	reader := bufio.NewReader(conn)
	welcomeMessage, err := reader.ReadString('\n')
	if err != nil {
		fmt.Printf("[Client %d] Error reading welcome message: %v", clientID, err)
		return
	}
	log.Printf("[Client %d] %s", clientID, welcomeMessage)

	// adding client ID to the request
	requestToSend.ClientID = clientID

	// encoding the request to JSON
	requestJson, err := json.Marshal(requestToSend)
	if err != nil {
		fmt.Printf("[Client %d] Error encoding request: %v", clientID, err)
		return
	}

	// sending the request
	fmt.Fprintf(conn, "%s\n", requestJson)
	fmt.Printf("[Client %d] -> Requested task #%d\n", clientID, requestToSend.TaskNumber)

	// waiting for the response
	responseJson, err := reader.ReadString('\n')
	if err != nil {
		fmt.Printf("[Client %d] <- Error reading response: %v", clientID, err)
		return
	}

	// decode the response
	var resp GenericResponse
	json.Unmarshal([]byte(responseJson), &resp)

	if resp.Status == "success" {
		fmt.Printf("[Client %d] <- Success: %s\n\n", clientID, string(resp.Result))
	} else {
		fmt.Printf("[Client %d] <- Error response: %s\n\n", clientID, resp.Error)
	}
}

func main() {
	// reading the JSON file with all tasks
	jsonFile, err := os.ReadFile("tasks.json")
	if err != nil {
		log.Fatalf("Error reading 'tasks.json': %v", err)
	}

	// parsing all requests
	var allRequests []GenericRequest
	err = json.Unmarshal(jsonFile, &allRequests)
	if err != nil {
		log.Fatalf("Error parsing: %v", err)
	}

	// getting the request for task requestedTaskNumber
	var requestsForTask []GenericRequest
	requestedTaskNumbers := []int{1, 2, 3, 4, 5, 6, 7} // specify the task numbers you want to run

	// finding the requests for the specified task numbers
	for _, requestedTaskNumber := range requestedTaskNumbers {
		found := false
		for _, req := range allRequests {
			if req.TaskNumber == requestedTaskNumber {
				requestsForTask = append(requestsForTask, req)
				found = true
				break
			}
		}

		if !found {
			log.Fatalf("Could not find request with 'task_number: %d' in JSON file.", requestedTaskNumber)
		}
	}

	// displaying the found requests
	for _, req := range requestsForTask {
		fmt.Printf("Task %d found with input: %s\n", req.TaskNumber, string(req.Input))
	}

	// launching multiple clients concurrently
	var wg sync.WaitGroup // waitgroup to wait for all clients to finish
	fmt.Printf("Running %d clients\n\n", NUM_CLIENTS)

	for i := 1; i <= NUM_CLIENTS; i++ {
		wg.Add(1)
		// launching each client in a separate goroutine
		go runSingleClient(i, &wg, requestsForTask[i-1])
		time.Sleep(50 * time.Millisecond) // break to see clearer logs
	}

	log.Println("Waiting for all clients to finish...")
	wg.Wait()
	log.Println("All clients have finished.")
}
