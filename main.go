package main

import (
	"fmt"
	"godrop/cli"
	"godrop/network"
)

func main() {
	config := cli.ParseConfig()

	fmt.Println("Mode:", config.Mode)
	fmt.Println("Port:", config.Port)
	fmt.Println("Path:", config.Path)

	switch config.Mode {
	case "send":
		network.SendFile(config.Port, config.Path)
	case "receive":
		network.ReceiveFile(config.Port, config.Path)
	}
}

