package network

import "fmt"

func sendFile(port int, path string) {
	fmt.Println("Sending file to port: ", port)
}

func receiveFile(port int, path string) {
	fmt.Println("Receiving file from port: ", port)
}