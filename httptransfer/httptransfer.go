package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"
)

var (
	port *int = flag.Int("port", 8080,
		"The port attached to the server",
	)
	filename_get *string = flag.String("get", "get_file",
		"File the client will receive from a GET request",
	)
	filename_post *string = flag.String("post", "post_file",
		"File the client will write to from a POST request",
	)

	logger  *log.Logger = log.Default()
	fp_get  *os.File
	fp_post *os.File
)

func main() {
	var err error

	flag.Parse()

	fp_get, err = os.Open(*filename_get)
	if err != nil {
		logger.Fatal(err)
	}
	defer fp_get.Close()

	fp_post, err = os.Create(*filename_post)
	if err != nil {
		logger.Fatal(err)
	}
	defer fp_post.Close()

	logger.Printf("Starting server at 0.0.0.0:%d", *port)
	http.HandleFunc("/", FileHandler)

	err = http.ListenAndServe(fmt.Sprintf("0.0.0.0:%d", *port), nil)
	logger.Fatalf("Server stopped: %s", err)
}

func FileGetRequest(w http.ResponseWriter, req *http.Request) {
	_, err := fp_get.Seek(0, io.SeekStart)
	if err != nil {
		logger.Printf("Post request seek fp_post: %v", err)
		w.WriteHeader(500)
		return
	}

	b := make([]byte, 32*1024)

	for {
		n, err := fp_get.Read(b)
		if err == io.EOF {
			break
		}
		if err != nil {
			logger.Printf("Get request read fp_get: %v", err)
			w.WriteHeader(500)
			return
		}

		_, err = w.Write(b[:n])
		if err != nil {
			logger.Printf("Get request write: %v", err)
			w.WriteHeader(500)
			return
		}
	}

	logger.Print("Sent complete file")
}

func FilePostRequest(w http.ResponseWriter, req *http.Request) {
	_, err := fp_post.Seek(0, io.SeekStart)
	if err != nil {
		logger.Printf("Post request seek fp_post: %v", err)
		w.WriteHeader(500)
		return
	}
	err = fp_post.Truncate(0)
	if err != nil {
		logger.Printf("Post request truncate fp_post: %v", err)
		w.WriteHeader(500)
		return
	}

	b := make([]byte, 32*1024)

	last_iter := false

	for !last_iter {
		n, err := req.Body.Read(b)
		if err == io.EOF {
			last_iter = true
		} else if err != nil {
			log.Printf("Post request read: %v", err)
			w.WriteHeader(500)
			return
		}

		_, err = fp_post.Write(b[:n])
		if err != nil {
			log.Printf("Post request write fp_post: %v", err)
			return
		}
	}

	log.Print("Got complete file from POST request")
}

func FileHandler(w http.ResponseWriter, req *http.Request) {
	start := time.Now()

	switch req.Method {
	case "GET":
		FileGetRequest(w, req)
	case "POST":
		FilePostRequest(w, req)
	default:
		log.Print("Got Method ", req.Method)
		w.WriteHeader(405)
	}

	log.Printf("Whole communication took %dms",
		time.Since(start).Milliseconds(),
	)
}
