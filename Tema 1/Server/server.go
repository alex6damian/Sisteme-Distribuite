package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math"
	"net"
	"os"
	"strconv"
	"strings"
	"time"
	"unicode"
)

// request and response structures
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

// configuration structure
type Config struct {
	Host                         string `json:"Host"`
	Port                         string `json:"Port"`
	WelcomeMessage               string `json:"WelcomeMessage"`
	MaxMessageSize               int    `json:"MaxMessageSize"`
	MaxConcurrentConnections     int    `json:"MaxConcurrentConnections"`
	ConnectionIdleTimeoutSeconds int    `json:"ConnectionIdleTimeoutSeconds"`
}

// loadConfig reads the configuration from config.json file
func loadConfig(filename string) (Config, error) {
	// open and read the config file
	var config Config
	configFile, err := os.Open(filename)
	if err != nil {
		return config, err
	}
	defer configFile.Close()

	jsonParser := json.NewDecoder(configFile)
	err = jsonParser.Decode(&config)
	return config, err
}

// sendErrorResponse sends an error response to the client
func sendErrorResponse(conn net.Conn, errorMsg string, timeout time.Duration) {
	response := GenericResponse{Status: "error", Error: errorMsg}
	responseJson, _ := json.Marshal(response)

	conn.SetWriteDeadline(time.Now().Add(timeout))
	_, err := conn.Write(append(responseJson, '\n'))
	if err != nil {
		log.Printf("Error while sending error response: %v", err)
	}
}

// handleConnection processes each client connection
func handleConnection(connection net.Conn, config Config, semaphore chan struct{}) {
	defer func() {
		connection.Close()
		<-semaphore // releasing the semaphore slot
		log.Printf("Connection closed: %s\n\n", connection.RemoteAddr().String())
	}()
	log.Printf("New connection from %s\n", connection.RemoteAddr().String())

	timeoutDuration := time.Duration(config.ConnectionIdleTimeoutSeconds) * time.Second

	// sending the welcome message
	connection.SetWriteDeadline(time.Now().Add(timeoutDuration))
	_, err := connection.Write([]byte(config.WelcomeMessage + "\n"))
	if err != nil {
		log.Printf("Error while sending welcome message: %v", err)
		return
	}

	reader := bufio.NewReader(connection)

	for {
		connection.SetReadDeadline(time.Now().Add(timeoutDuration))

		requestJson, err := reader.ReadString('\n')
		if err != nil {
			if err != io.EOF {
				log.Printf("Error reading from %s: %v", connection.RemoteAddr().String(), err)
			}
			break
		}

		log.Printf("Received from %s: %s", connection.RemoteAddr().String(), requestJson)

		var req GenericRequest
		if err := json.Unmarshal([]byte(requestJson), &req); err != nil {
			log.Printf("Error decoding JSON: %v. Request: %s", err, requestJson)
			sendErrorResponse(connection, "Invalid JSON", timeoutDuration)
			continue
		}

		log.Printf("Processing request #%d with input: %s", req.TaskNumber, string(req.Input))

		// handling the task
		results, err := handleTask(req.TaskNumber, req.Input)
		if err != nil {
			log.Printf("Error handling task: %v", err)
			sendErrorResponse(connection, "Internal server error", timeoutDuration)
			continue
		}

		// sending the response
		response := GenericResponse{
			Status: "success",
			Result: results,
		}
		responseJson, _ := json.Marshal(response)

		// setting write deadline and sending response
		connection.SetWriteDeadline(time.Now().Add(timeoutDuration))
		_, err = connection.Write(append(responseJson, '\n'))
		if err != nil {
			log.Printf("Error while sending response: %v", err)
			break
		}
		log.Printf("Response sent to %s: %s", connection.RemoteAddr().String(), responseJson)
	}
}

func main() {
	// loading configuration
	config, err := loadConfig("config.json")
	if err != nil {
		log.Fatalf("Error loading configuration: %s\n", err.Error())
	}
	fmt.Printf("Configuration loaded: %+v\n", config)

	// semaphore to limit concurrent connections
	semaphore := make(chan struct{}, config.MaxConcurrentConnections)

	// starting the server using config
	listenAddr := net.JoinHostPort(config.Host, config.Port)
	listener, err := net.Listen("tcp", listenAddr)
	if err != nil {
		log.Fatalf("Error starting server: %s\n", err.Error())
	}
	defer listener.Close()
	fmt.Printf("Server listening on %s\n", listenAddr)

	// accepting connections
	for {
		connection, err := listener.Accept()
		if err != nil {
			log.Printf("Error accepting connection: %s\n", err.Error())
			continue
		}

		// selecting a slot in the semaphore
		semaphore <- struct{}{}

		// handling the connection in a new goroutine
		go handleConnection(connection, config, semaphore)
	}
}

func handleTask(taskNumber int, input json.RawMessage) (json.RawMessage, error) {

	var results interface{}
	switch taskNumber {
	case 1:
		var inputs []string
		if err := json.Unmarshal(input, &inputs); err != nil {
			return nil, fmt.Errorf("invalid input format for task 1")
		}
		results = task1(inputs)
	case 2:
		var inputs []string
		if err := json.Unmarshal(input, &inputs); err != nil {
			return nil, fmt.Errorf("invalid input format for task 2")
		}
		results = task2(inputs)
	case 3:
		var inputs []int
		if err := json.Unmarshal(input, &inputs); err != nil {
			return nil, fmt.Errorf("invalid input format for task 3")
		}
		results = task3(inputs)
	case 4:
		var inputs []int
		if err := json.Unmarshal(input, &inputs); err != nil {
			return nil, fmt.Errorf("invalid input format for task 4")
		}
		results = task4(inputs)
	case 5:
		var inputs []string
		if err := json.Unmarshal(input, &inputs); err != nil {
			return nil, fmt.Errorf("invalid input format for task 5")
		}
		results = task5(inputs)
	case 6:
		var inputs []string
		if err := json.Unmarshal(input, &inputs); err != nil {
			return nil, fmt.Errorf("invalid input format for task 6")
		}
		results = task6(inputs)
	default:
		return nil, fmt.Errorf("unknown task number: %d", taskNumber)
	}
	resultsJson, err := json.Marshal(results)
	if err != nil {
		return nil, fmt.Errorf("error encoding results to JSON")
	}
	return json.RawMessage(resultsJson), nil
}

func task1(input []string) []string {
	// casa, masa, trei, tanc, 4321 -> cmtt4, aara3, ssen2, aaic1
	words := input
	words_len := len([]rune(words[0]))
	results := make([]string, 0, words_len)
	// selecting characters by index from each word
	for i := 0; i < words_len; i++ {
		var sb strings.Builder
		for _, word := range words {
			r := []rune(word)
			sb.WriteRune(r[i])
		}
		results = append(results, sb.String())

	}
	return results
}

func task2(input []string) int {
	// abd4g5, 1sdf6fd, fd2fdsf5 -> 2 patrate perfecte (16, 25)
	count := 0
	// extracting the digits
	for _, word := range input {
		var sb strings.Builder
		for _, ch := range word {
			if unicode.IsDigit(ch) {
				sb.WriteRune(ch)
			}
		}
		if sb.Len() > 0 {
			// convert to integer and check perfect square
			num, err := strconv.Atoi(sb.String())
			if err == nil {
				sqrt := math.Sqrt(float64(num))
				if sqrt == math.Trunc(sqrt) {
					count++
				}
			}
		}
	}
	return count
}

func task3(input []int) int {
	// 12, 13, 14 => 21 + 31 + 41 = 93
	sum := 0
	// reversing and summing
	for _, num := range input {
		s := strconv.Itoa(num)
		var sb strings.Builder
		r := []rune(s)
		for i := len(r) - 1; i >= 0; i-- {
			sb.WriteRune(r[i])
		}
		num, _ := strconv.Atoi(sb.String())
		sum += num
	}
	return sum
}

func task4(input []int) int {
	average := 0
	count := 0
	min_bound := input[0]
	max_bound := input[1]
	// calculating average of numbers whose sum of digits is within bounds
	for i := 3; i < len(input); i++ {
		sum := 0
		s := strconv.Itoa(input[i])
		r := []rune(s)
		for j := len(r) - 1; j >= 0; j-- {
			digit, _ := strconv.Atoi(string(r[j]))
			sum += digit
		}
		if sum >= min_bound && sum <= max_bound {
			average += input[i]
			count++
		}
	}
	if count > 0 {
		average /= count
	}
	return average
}

func task5(input []string) []int {
	// 2dasdas, 12, dasdas, 1010, 101 => 10, 5 (1010=10, 101=5)
	results := make([]int, 0)
	for _, word := range input {
		// checking if the word is a binary number
		num, err := strconv.ParseInt(word, 2, 64)
		if err == nil {
			results = append(results, int(num))
		}

	}
	return results
}

// for uppercase, check if upper, convert to lower, do the same and convert back
func task6(input []string) []string {
	// LEFT, 3, abcdef, salut, ceva => xyzabc, pxirq, zbsx
	// extracting direction and number of steps
	direction := input[0]
	steps, _ := strconv.Atoi(input[1])
	results := make([]string, 0)
	for i := 2; i < len(input); i++ {
		var sb strings.Builder
		// building the shifted string
		for _, ch := range input[i] {
			var shifted rune
			base := int(ch)
			switch direction {
			case "LEFT":
				shifted = rune((base-97-steps+26)%26 + 97)
			case "RIGHT":
				shifted = rune((base-97+steps)%26 + 97)
			}
			sb.WriteRune(shifted)
		}
		results = append(results, sb.String())
	}
	return results
}
