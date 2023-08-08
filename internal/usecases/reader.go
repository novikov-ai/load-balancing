package usecases

import (
	"fmt"
	"io"
	"net"
)

func ReadAllData(src net.Conn) []byte {
	data := make([]byte, 1024)
	for {
		n, err := src.Read(data)
		if err != nil {
			if err != io.EOF {
				fmt.Println("Error reading data:", err)
			}
			break
		}

		received := data[:n]
		return received
	}

	return nil
}
