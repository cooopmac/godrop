package cli

import (
	"flag"
	"fmt"
	"os"
)

type Config struct {
	Mode string
	Port int
	Path string
}

func ParseConfig() Config {
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

	return Config{*mode, *port, *path}
}