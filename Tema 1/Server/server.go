package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"time"
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
		<-semaphore
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

		// TODO: implementing the tasks

		responseText := "Task is not developed yet."

		response := GenericResponse{
			Status: "success",
			Result: json.RawMessage(fmt.Sprintf(`"%s"`, responseText)),
		}
		responseJson, _ := json.Marshal(response)

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
