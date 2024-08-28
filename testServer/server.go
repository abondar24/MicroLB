package main

import (
	"flag"
	"log"
	"net/http"
)

func handler(port string) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {
		log.Printf("Server is up on port %s", port)

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
