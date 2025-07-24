package network

import (
	"fmt"
	"godrop/cli"
	"io"
	"net"
	"os"
	"path/filepath"
	"strconv"
)

func SendFile(host string, port int, path string) {
	file, err := os.Open(path)
	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}
	defer file.Close()

	fileInfo, err := file.Stat()
	if err != nil {
		fmt.Println("Error getting file info:", err)
		return
	}
	fileSize := fileInfo.Size()
	fileName := filepath.Base(path)

	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", host, port))
	if err != nil {
		fmt.Println("Error dialing:", err)
		return
	}
	defer conn.Close()

	fmt.Printf("Sending %s (%d bytes) to %s:%d\n", fileName, fileSize, host, port)
	
	_, err = conn.Write([]byte(fmt.Sprintf("%d\n", fileSize)))
	if err != nil {
		fmt.Println("Error sending file size:", err)
		return
	}
	
	_, err = conn.Write([]byte(fileName + "\n"))
	if err != nil {
		fmt.Println("Error sending filename:", err)
		return
	}

	buffer := make([]byte, 32*1024)
	var totalBytes int64
	progressBar := cli.NewProgressBar()

	for {
		n, err := file.Read(buffer)
		if n > 0 {
			_, writeErr := conn.Write(buffer[:n])
			if writeErr != nil {
				fmt.Println("Error sending data:", writeErr)
				return
			}
			totalBytes += int64(n)
			
			progressBar.Update(totalBytes, fileSize)
		}
		if err == io.EOF {
			break
		}
		if err != nil {
			fmt.Println("Error reading file:", err)
			return
		}
	}

	progressBar.Finish(totalBytes, "File sent")
}

func ReceiveFile(port int, path string) {
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		fmt.Println("Error listening on port:", port)
		return
	}
	defer listener.Close()
	fmt.Println("Listening on port:", port)

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Error accepting connection:", err)
			continue
		}
		fmt.Println("Connection accepted")
		
		go handleConnection(conn, path)
	}
}

func handleConnection(conn net.Conn, basePath string) {
	defer conn.Close()
	
	var fileSize int64
	fileSizeBytes := make([]byte, 0)
	buffer := make([]byte, 1)
	
	for {
		n, err := conn.Read(buffer)
		if err != nil {
			fmt.Println("Error reading file size:", err)
			return
		}
		if n > 0 && buffer[0] == '\n' {
			break
		}
		if n > 0 {
			fileSizeBytes = append(fileSizeBytes, buffer[0])
		}
	}
	
	var err error
	fileSize, err = strconv.ParseInt(string(fileSizeBytes), 10, 64)
	if err != nil {
		fmt.Println("Error parsing file size:", err)
		return
	}
	
	var originalFilename string
	filenameBytes := make([]byte, 0)
	
	for {
		n, err := conn.Read(buffer)
		if err != nil {
			fmt.Println("Error reading filename:", err)
			return
		}
		if n > 0 && buffer[0] == '\n' {
			break
		}
		if n > 0 {
			filenameBytes = append(filenameBytes, buffer[0])
		}
	}
	
	originalFilename = string(filenameBytes)
	
	var outputPath string
	if basePath == "received_file" || basePath == "" {
		outputPath = "received_" + originalFilename
	} else {
		dir := filepath.Dir(basePath)
		if dir == "." {
			outputPath = "received_" + originalFilename
		} else {
			outputPath = filepath.Join(dir, "received_" + originalFilename)
		}
	}
	
	fmt.Printf("Receiving %s (%d bytes)\n", originalFilename, fileSize)

	file, err := os.Create(outputPath)
	if err != nil {
		fmt.Println("Error creating file:", err)
		return
	}
	defer file.Close()

	fmt.Println("Starting to receive data...")
	
	receiveBuffer := make([]byte, 32*1024)
	var totalBytes int64
	progressBar := cli.NewProgressBar()

	for totalBytes < fileSize {
		remaining := fileSize - totalBytes
		readSize := int64(len(receiveBuffer))
		if remaining < readSize {
			readSize = remaining
		}
		
		n, err := conn.Read(receiveBuffer[:readSize])
		if err != nil {
			if err != io.EOF {
				fmt.Println("Error receiving data:", err)
			}
			break
		}
		
		if n > 0 {
			_, writeErr := file.Write(receiveBuffer[:n])
			if writeErr != nil {
				fmt.Println("Error writing to file:", writeErr)
				return
			}
			totalBytes += int64(n)
			
			progressBar.Update(totalBytes, fileSize)
		}
	}

	progressBar.Finish(totalBytes, "File received at " + outputPath)
}