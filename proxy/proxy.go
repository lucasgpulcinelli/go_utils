package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"sync"
	"time"
)

var (
	port *int = flag.Int("port", 8080,
		"port to listen (that the client will connect to)",
	)
	addr *string = flag.String("addr", "server:8080",
		"address to connect (real server)",
	)

	logger *log.Logger = log.Default()
)

/*
 * HandleProxy handles the logging and transmission of all traffic from a
 * reader to a writer. This function should be called in pairs (to handle
 * bidirectional traffic) as separated goroutines, having a common waitgroup
 * and channel to signal the end of the connection.
 */
func HandleProxy(
	reader, writer net.Conn,
	done chan struct{},
	wg *sync.WaitGroup,
) {
	defer wg.Done()

	msg := make([]byte, 1024)
	for {
		select {
		case <-done:
			return
		default:
		}

		//hardcoded 5 second deadline to reevaluate connection state
		reader.SetReadDeadline(time.Now().Add(time.Second * 5))
		n, err := reader.Read(msg)
		if os.IsTimeout(err) {
			continue
		}
		if err == io.EOF {
			logger.Print("disconnected ", reader.RemoteAddr())
			break
		}
		if err != nil {
			logger.Printf("error at read %s: %v", reader.RemoteAddr(), err)
			break
		}

		fmt.Print(reader.RemoteAddr())
		fmt.Printf(" ---(%d bytes)---> ", n)
		fmt.Println(writer.RemoteAddr())

		fmt.Println("--------------------")
		fmt.Println(string(msg[:n]))
		fmt.Print("--------------------\n\n")

		writer.SetWriteDeadline(time.Now().Add(time.Second * 5))
		_, err = writer.Write(msg[:n])
		if os.IsTimeout(err) {
			logger.Print("timeout error at write ", writer.RemoteAddr())
			continue
		}
		if err != nil {
			logger.Printf("error at write %s : %v", writer.RemoteAddr(), err)
			break
		}
	}
	done <- struct{}{}
}

/*
 * HandleConnection handles the logging and transmission of all traffic
 * between a server and a client, creating two goroutines and blocking
 * execution until the connection is terminated.
 */
func HandleConnection(client, server net.Conn) {
	done := make(chan struct{}, 2)
	defer close(done)

	wg := &sync.WaitGroup{}

	wg.Add(2)
	go HandleProxy(client, server, done, wg)
	go HandleProxy(server, client, done, wg)

	wg.Wait()
	logger.Printf(
		"closing connection with %s and %s",
		client.RemoteAddr(),
		server.RemoteAddr(),
	)
	client.Close()
	server.Close()
}

func main() {
	flag.Parse()

	logger.Printf("Starting proxy at 0.0.0.0:%d", *port)

	listener, err := net.Listen("tcp", fmt.Sprintf("0.0.0.0:%d", *port))
	if err != nil {
		log.Fatal(err)
	}

	for {
		client, err := listener.Accept()
		if err != nil {
			logger.Fatal(err)
		}
		logger.Print("accepted connection with ", client.RemoteAddr())

		server, err := net.Dial("tcp", *addr)
		if err != nil {
			logger.Print("error connecting to server: ", err)
			client.Close()
			continue
		}
		logger.Print("server connected at ", server.RemoteAddr())

		go HandleConnection(client, server)
	}
}
