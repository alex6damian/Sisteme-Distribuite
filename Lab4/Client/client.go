package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
)

func main() {
	// reading message from keyboard
	fmt.Print("Send a message for the server: ")
	reader := bufio.NewReader(os.Stdin)
	message, err := reader.ReadString('\n')
	if err != nil {
		log.Fatalf("Error reading from keyboard: %v", err)
	}

	// connecting to the server
	connection, err := net.Dial("tcp", "localhost:8080")
	if err != nil {
		log.Fatalf("Error connecting to server: %v", err)
	}
	// closing the connection at the end
	defer connection.Close()

	// sending message to the server
	_, err = connection.Write([]byte(message))
	if err != nil {
		log.Fatalf("Error sending message to server: %v", err)
	}

}
