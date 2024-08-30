package main

import (
	"github.com/stretchr/testify/assert"
	"io"
	"net"
	"net/http"
	"net/http/httptest"
	"testing"
)

func MockHttpServer(response string) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		io.WriteString(w, response)
	}))
}

func MockTcpServer(t *testing.T, response string) (net.Listener, string) {
	listener, err := net.Listen("tcp", ":9090") // Bind to a random available port
	assert.NoError(t, err)

	address := listener.Addr().String() // Get the address of the listener

	go func() {
		conn, _ := listener.Accept()
		defer conn.Close()
		io.WriteString(conn, response)
	}()

	return listener, address
}

func TestGetNextBackend(t *testing.T) {
	backends := []string{"back1", "back2", "back3"}

	lb := &LoadBalancer{
		backends: backends,
	}

	for _, exp := range backends {
		got := lb.getNextBackend()
		assert.Equal(t, exp, got)
	}

}

func TestHttpHandler(t *testing.T) {
	lb := &LoadBalancer{
		backends: []string{"back1", "back2"},
	}

	server1 := MockHttpServer("Hello from server 1")
	defer server1.Close()

	server2 := MockHttpServer("Hello from server 2")
	defer server2.Close()

	lb.backends = []string{server1.URL[7:], server2.URL[7:]}

	assert.Equal(t, "Hello from server 1", performHttpRequest(t, lb, "/"), "Response body should be from server 1 on the first request")
	assert.Equal(t, "Hello from server 2", performHttpRequest(t, lb, "/"), "Response body should be from server 2 on the second request")
}

func performHttpRequest(t *testing.T, lb *LoadBalancer, requestURL string) string {
	req := httptest.NewRequest(http.MethodGet, requestURL, nil)
	w := httptest.NewRecorder()
	lb.httpHandler(w, req)
	resp := w.Result()
	body, _ := io.ReadAll(resp.Body)
	return string(body)
}

func TestTcpHandler(t *testing.T) {
	mockResponse := "Hello from TCP backend"
	listener, backendAddress := MockTcpServer(t, mockResponse)
	defer listener.Close()

	lb := &LoadBalancer{
		backends: []string{backendAddress},
	}

	conn, err := net.Dial("tcp", ":9090") // Ensure that your TCP load balancer is listening on this port
	assert.NoError(t, err)
	defer conn.Close()

	// Handle the connection using the load balancer
	go lb.tcpHandler(conn)

	// Read the response from the load balancer
	buf := make([]byte, 1024)
	n, err := conn.Read(buf)
	assert.NoError(t, err)

	// Use assertions to check the response
	assert.Equal(t, mockResponse, string(buf[:n]), "Response from TCP backend should match")
}
