package main

import (
	"flag"
	"fmt"
	"github.com/novikov-ai/load-balancing/internal/usecases"
	"net"
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

func main() {
	flag.Parse()

	if *backends == "" {
		fmt.Println("Please provide at least one backend")
		os.Exit(1)
		return
	}

	backendsAddresses := strings.Split(*backends, ",")

	fmt.Println("Load balancer L4 started with backends:", backendsAddresses)

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

		handleConnL4(clientConn, backendConn)

		backendIndex = (backendIndex + 1) % len(backendsAddresses)
	}
}

func handleConnL4(conn, backendConn net.Conn) {
	defer conn.Close()
	defer backendConn.Close()

	data := usecases.ReadAllData(conn)
	backendConn.Write(data)
}
