package main

import (
	"io"
	"log"
	"net"
	"net/http"
	"sync/atomic"
)

type LoadBalancer struct {
	backends []string
	current  uint32
}

func (lb *LoadBalancer) StartTcpLoadBalancer(port string) {
	listener, err := net.Listen("tcp", ":"+port)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Starting TCP load balancer on port %s\n", port)

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Printf("Error accepting connection: %s", err)
			continue
		}
		go lb.tcpHandler(conn)
	}
}

func (lb *LoadBalancer) tcpHandler(conn net.Conn) {
	backend := lb.getNextBackend()
	log.Printf("Forwarding TCP connection to %s", backend)

	backendConn, err := net.Dial("tcp", backend)
	if err != nil {
		log.Printf("Error forwarding TCP connection to %s: %s\n", backend, err)
		conn.Close()
		return
	}

	go func() {
		io.Copy(backendConn, conn)
		backendConn.Close()
	}()

	io.Copy(conn, backendConn)
	conn.Close()

}

func (lb *LoadBalancer) StartHttpLoadBalancer(port string) {
	http.HandleFunc("/", lb.httpHandler)
	log.Printf("Starting HTTP load balancer on port %s\n", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}

func (lb *LoadBalancer) httpHandler(w http.ResponseWriter, r *http.Request) {
	backend := lb.getNextBackend()

	req := r.Clone(r.Context())
	req.URL.Scheme = "http"
	req.URL.Host = backend
	log.Printf("Forwarding HTTP request to %s\n", req.URL.String())

	resp, err := http.DefaultTransport.RoundTrip(req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadGateway)
		return
	}
	defer resp.Body.Close()

	for k, v := range resp.Header {
		w.Header()[k] = v
	}
	w.WriteHeader(resp.StatusCode)
	io.Copy(w, resp.Body)

}

// round-robin to get next backend
func (lb *LoadBalancer) getNextBackend() string {
	backendIndex := atomic.AddUint32(&lb.current, 1)
	selectedIndex := (int(backendIndex) - 1) % len(lb.backends)
	log.Printf("Selected backend index: %d", selectedIndex)
	return lb.backends[selectedIndex]
}
