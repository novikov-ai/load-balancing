package main

import (
	"fmt"
	"net"
	"sync"

	"github.com/novikov-ai/load-balancing/internal/usecases"
)

var (
	backends = map[string]string{
		"Backend #1": ":1111",
		"Backend #2": ":2222",
		"Backend #3": ":3333",
	}
)

func main() {
	wg := sync.WaitGroup{}

	for name, port := range backends {
		wg.Add(1)

		go handleBackend(&wg, name, port)
	}

	wg.Wait()
}

func handleBackend(wg *sync.WaitGroup, name, port string) {
	defer wg.Done()

	address := "127.0.0.1" + port
	listener, err := net.Listen("tcp", address)
	if err != nil {
		fmt.Println("Error occurred:", err)
		return
	}
	defer listener.Close()

	fmt.Println("Backend is listening on:", listener.Addr())

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Error accepting connection:", err)
			continue
		}

		received := usecases.ReadAllData(conn)
		fmt.Printf("Received data from %s (%s): %s", name, address, received)
	}
}
