package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"net/http"
)

func handler(port string) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {
		log.Printf("Server is up on port %s", port)
		response := fmt.Sprintf("Server is running on port %s", port)
		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(http.StatusOK)
		_, err := w.Write([]byte(response))
		if err != nil {
			log.Printf("Error writing response: %v", err)
		}
	}
}

func startHttpServer(port string) {
	http.HandleFunc("/", handler(port))
	log.Printf("Starting HTTP server on port %s", port)
	err := http.ListenAndServe(":"+port, nil)
	if err != nil {
		log.Fatal(err)
	}
}

func startTcpServer(port string) {
	listener, err := net.Listen("tcp", ":"+port)
	if err != nil {
		log.Fatal(err)
	}
	defer listener.Close()

	log.Printf("Starting TCP server on port %s\n", port)
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Printf("Error accepting connection: %s\n", err)
			continue
		}
		go handleTcpConnection(conn, port)
	}
}

func handleTcpConnection(conn net.Conn, port string) {
	defer conn.Close()
	_, err := conn.Write([]byte(fmt.Sprintf("Server is running on port %s\n", port)))
	if err != nil {
		log.Printf("Error writing to connection: %v", err)
	}
}

func main() {
	port := flag.String("port", "", "Port to run server on")
	mode := flag.String("mode", "http", "Mode to run server in: 'http' or 'tcp'")

	flag.Parse()

	if *port == "" {
		log.Fatal("You must specify a port to run server on using -port flag")

	}

	if *mode != "http" && *mode != "tcp" {
		log.Fatalf("Invalid mode '%s'. Must be 'http' or 'tcp'", *mode)
	}

	if *mode == "http" {
		startHttpServer(*port)
	} else if *mode == "tcp" {
		startTcpServer(*port)
	}
}
