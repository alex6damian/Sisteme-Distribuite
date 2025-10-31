package main

import (
	"fmt"
	"net/http"
)

// IV.2
func httpServerExample(w http.ResponseWriter, req *http.Request) {
	fmt.Fprintf(w, "Hello everyone! This is a simple HTTP server in Go.")
}

func main() {
	http.HandleFunc("/hello", httpServerExample)
	fmt.Println("Starting server at port 8080...")
	http.ListenAndServe(":8080", nil)
}
