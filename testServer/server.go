package main

import (
	"flag"
	"fmt"
	"log"
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

func main() {
	port := flag.String("port", "", "Port to run server on")
	flag.Parse()

	if *port == "" {
		log.Fatal("You must specify a port to run server on using -port flag")

	}

	http.HandleFunc("/", handler(*port))
	err := http.ListenAndServe(":"+*port, nil)
	if err != nil {
		log.Fatal(err)
	}
}
