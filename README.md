# Micro Load Balancer (MicroLB)

MicroLB is a simple load balancer based on round-robin algorithm implemented in Go, capable of handling both HTTP and TCP load balancing. This project includes a test server that can run in either HTTP or TCP mode.

## Build and Run

### Build

- To build the MicroLB load balancer, use the following command:

```sh
cd /lb

go build
```

- To build the Test server, use the following command:

```sh
cd /testServer
go build
```

### Run

- Load Balancer
```sh
./lb -backends=../backensd.txt -tcpPort=<tcpPort> -httpPort=<httpPort>
```

- -backends: Path to the configuration file containing the list of backend servers.
- -tcpPort: Port for TCP load balancing (optional, default is 9090).
- -httpPort: Port for HTTP load balancing (optional, default is 8080).


- Test server
```sh
./testServer -port=<portValue> -mode=<modeValue>

```

- -port: Port for the server to listen on
- -mode: Mode to run the server in, either tcp or http.

### Configuration file

The configuration file for the load balancer (backends.txt) should contain a list of backend servers, one per line:
```
localhost:8020
localhost:8030
```