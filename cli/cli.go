package cli

import (
	"flag"
	"fmt"
	"os"
	"time"
)

type Config struct {
	Mode string
	Port int
	Path string
	Host string
}

type ProgressBar struct {
	startTime    time.Time
	lastProgress int
}

func NewProgressBar() *ProgressBar {
	return &ProgressBar{
		startTime:    time.Now(),
		lastProgress: -1,
	}
}

func (pb *ProgressBar) Update(current, total int64) {
	progress := int((current * 100) / total)
	if progress != pb.lastProgress && progress%5 == 0 {
		elapsed := time.Since(pb.startTime)
		speed := float64(current) / elapsed.Seconds() / 1024

		barWidth := 20
		filled := (progress * barWidth) / 100
		bar := ""
		for i := 0; i < barWidth; i++ {
			if i < filled {
				bar += "█"
			} else {
				bar += "░"
			}
		}

		fmt.Printf("\r[%s] %d%% - %.1f KB/s", bar, progress, speed)
		pb.lastProgress = progress
	}
}

func (pb *ProgressBar) Finish(totalBytes int64, operation string) {
	elapsed := time.Since(pb.startTime)
	avgSpeed := float64(totalBytes) / elapsed.Seconds() / 1024
	fmt.Printf("\n✓ %s successfully! Total: %d bytes in %.2fs (%.1f KB/s)\n",
		operation, totalBytes, elapsed.Seconds(), avgSpeed)
}

func ParseConfig() Config {
	mode := flag.String("mode", "", "mode: send or receive")
	port := flag.Int("port", 8888, "port: starting on port 8888")
	path := flag.String("path", "", "path: path to the file")
	host := flag.String("host", "localhost", "host: target machine IP address")

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %s --mode <send|receive> --port <number> [--path <file>] [--host <ip>]\n", os.Args[0])
		flag.PrintDefaults()
	}

	flag.Parse()

	if *mode == "" {
		fmt.Println("Error: mode is required (send or receive)")
		flag.Usage()
		os.Exit(1)
	}

	if *mode != "send" && *mode != "receive" {
		fmt.Println("Error: mode must be send or receive")
		flag.Usage()
		os.Exit(1)
	}

	if *port <= 0 || *port > 9999 {
		fmt.Println("Error: port must be between 1 and 9999")
		flag.Usage()
		os.Exit(1)
	}

	if *mode == "send" && *path == "" {
		fmt.Println("Error: path is required for send mode")
		flag.Usage()
		os.Exit(1)
	}

	if *mode == "receive" && *path == "" {
		*path = "received_file"
		fmt.Println("Warning: path is not provided for receive mode, using default path: ", *path)
	}

	return Config{*mode, *port, *path, *host}
}