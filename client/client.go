package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"net/http"
	"time"
)

var (
	addr    *string = flag.String("addr", "http://server", "address")
	port    *int    = flag.Int("port", 8080, "port")
	sleep_t *int    = flag.Int("time", 1000,
		"time in ms to sleep between requests",
	)
)

func main() {
	flag.Parse()

	client := &http.Client{}
	log := log.Default()

	for {
		time.Sleep(time.Duration(*sleep_t) * time.Millisecond)

		r, err := client.Get(fmt.Sprintf("%s:%d", *addr, *port))
		if err != nil {
			log.Printf("Get: %v", err)
			continue
		}

		if r.Status != "200 OK" {
			log.Printf("Get: %v", r.Status)
			continue
		}

		scanner := bufio.NewScanner(r.Body)

		fmt.Println("response:\n----------")
		for scanner.Scan() {
			fmt.Println(scanner.Text())
		}
		fmt.Println("----------")

		err = scanner.Err()
		if err != nil {
			log.Printf("Scanner: %v", err)
			continue
		}
	}
}
