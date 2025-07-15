package network

import (
	"fmt"
	"io"
	"net"
	"os"
)

func SendFile(port int, path string) {
	fmt.Println("Sending file to port: ", port)
}

func ReceiveFile(port int, path string) {
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		fmt.Println("Error listening on port: ", port)
		return
	}
	defer listener.Close()

	fmt.Println("Listening on port: ", port)
	conn, err := listener.Accept()
	if err != nil {
		fmt.Println("Error accepting connection: ", err)
		return
	}
	defer conn.Close()
	fmt.Println("Connection accepted")

	file, err := os.Create(path) 
	if err != nil {
		fmt.Println("Error creating file: ", err)
		return
	}
	defer file.Close()

	fmt.Println("Starting to receive data...")
	n, err := io.Copy(file, conn)
	if err != nil {
		fmt.Println("Error receiving file: ", err)
		return
	}

	fmt.Printf("File received at %s (%d bytes)\n", path, n)
}