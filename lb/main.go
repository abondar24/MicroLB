package main

import (
	"flag"
	"log"
)

func main() {

	backendsFile := flag.String("backends", "", "Backends to redirect")
	tcpPort := flag.String("tcpPort", "9090", "TCP LB Port")
	httpPort := flag.String("httpPort", "8080", "HTTP LB Port")

	flag.Parse()

	if *backendsFile == "" {
		log.Fatal("You must specify a config file with backends location")

	}

	backends, err := LoadBackends(*backendsFile)
	if err != nil {
		log.Fatal(err)
	}

	if len(backends) == 0 {
		log.Fatal("No backends found")
	}

	lb := &LoadBalancer{
		backends: backends,
	}

	go lb.StartHttpLoadBalancer(*httpPort)
	go lb.StartTcpLoadBalancer(*tcpPort)

	select {}
}
