package main

import (
	"bufio"
	"fmt"
	"net"
)

func handleConnection(connection net.Conn) {
	defer connection.Close()
	_, err := connection.Write([]byte("Hello from server!\n"))
	if err != nil {
		fmt.Println("Error writing to connection:", err)
	}

	message, err := bufio.NewReader(connection).ReadString('\n')
	if err != nil {
		fmt.Println("Error reading from connection:", err)
		return
	}
	fmt.Print("Message received from client:", message)
}

func main() {
	listenAddr := net.JoinHostPort("localhost", "8080")
	listener, err := net.Listen("tcp", listenAddr)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Server is listening on %s\n", listenAddr)

	// accepting connections
	for {
		connection, err := listener.Accept()
		if err != nil {
			panic(err)
		}
		go handleConnection(connection)
	}

}
