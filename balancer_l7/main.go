package main

import (
	"bufio"
	"flag"
	"fmt"
	"net"
	"net/http"
	"os"
	"strings"
)

/*
ROUND ROBIN:
1) Server List => Backend servers
2) Initial Selection
3) Request Forwarding
4) Updating Order
5) Even Distribution
*/

var backends = flag.String("backends", "", "Comma-separated list of backend addresses (host:port)")

var loadBalancerURL = "http://localhost:8080"

func main() {
	flag.Parse()

	if *backends == "" {
		fmt.Println("Please provide at least one backend")
		os.Exit(1)
		return
	}

	backendsAddresses := strings.Split(*backends, ",")

	fmt.Println("Load balancer L7 started with backends:", backendsAddresses)

	listener, err := net.Listen("tcp", ":8080")
	if err != nil {
		fmt.Println("Error occurred:", err)
		os.Exit(1)
		return
	}

	defer listener.Close()

	backendIndex := 0

	for {
		clientConn, err := listener.Accept()
		if err != nil {
			fmt.Println("Error accepting connections:", err)
			continue
		}

		fmt.Println("Connection established")

		backendConn, err := net.Dial("tcp", backendsAddresses[backendIndex])
		if err != nil {
			fmt.Println("Error connecting to backend:", err)
			clientConn.Close()
			continue
		}

		handleConnL7(clientConn, backendConn)

		backendIndex = (backendIndex + 1) % len(backendsAddresses)
	}
}

func handleConnL7(conn, backendConn net.Conn) {
	defer conn.Close()
	defer backendConn.Close()

	reader := bufio.NewReader(conn)
	request, err := http.ReadRequest(reader)
	if err != nil {
		fmt.Println("Error reading HTTP request:", err)
		conn.Write([]byte("HTTP/1.1 400 Bad Request\n"))
		return
	}

	if request.Method == http.MethodPost {
		conn.Write([]byte("HTTP/1.1 400 Bad Request\n"))
		return
	}

	fmt.Println("Successfully handled GET!")
}
