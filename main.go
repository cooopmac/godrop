package main

import (
	"flag"
	"fmt"
	"os"
)


func main() {
	mode := flag.String("mode", "", "mode: send or receive")
	port := flag.Int("port", 8888, "port: starting on port 8888")
	path := flag.String("path", "", "path: path to the file")

	flag.Parse()

	if *mode != "send" && *mode != "receive" {
		fmt.Println("Error: mode must be send or receive")
		os.Exit(1)
	}

	if *port <= 0 || *port > 9999 {
		fmt.Println("Error: port must be between 1 and 9999")
		os.Exit(1)
	}

	if *mode == "send" && *path == "" {
		fmt.Println("Error: path is required for send mode")
		os.Exit(1)
	}

	if *mode == "receive" && *path == "" {
		*path = "received_file.txt"
		fmt.Println("Warning: path is not provided for receive mode, using default path: ", *path)
	}

	if *mode == "send" {
		fmt.Println("Running in SEND mode.")
	}

	if *mode == "receive" {
		fmt.Println("Running in RECEIVE mode.")
	}

	fmt.Println("Mode: ", *mode)
	fmt.Println("Port: ", *port)
	fmt.Println("Path: ", *path)

	switch *mode {
	case "send":
		network.sendFile(*port, *path)
	case "receive":
		network.receiveFile(*port, *path)
	default:
		fmt.Println("Error: invalid mode")
		os.Exit(1)
	}

}

