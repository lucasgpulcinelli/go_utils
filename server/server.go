package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
)

var (
	port          *int        = flag.Int("port", 8080, "The port attached to the server")
	response      *string     = flag.String("response", "Hello, world!", "The response the server will give")
	logger        *log.Logger = log.Default()
	num_responses int         = 0
)

func main() {
	flag.Parse()

	logger.Printf("Starting server at 0.0.0.0:%d", *port)

	http.HandleFunc("/", HelloHandler)
	err := http.ListenAndServe(fmt.Sprintf("0.0.0.0:%d", *port), nil)

	logger.Fatalf("Server stopped: %s", err)
}

func HelloHandler(w http.ResponseWriter, _ *http.Request) {
	fmt.Fprintf(w, *response)

	num_responses++
	logger.Printf("Responded %d times", num_responses)
}
