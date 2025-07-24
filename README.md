# GoDrop - Cross-Platform File Transfer Tool

A simple, fast file transfer tool built in Go. Works like AirDrop but for any operating system - transfer files between machines on your local network with visual progress bars.

## Features

-   🚀 **Fast transfers** with real-time progress bars
-   🔄 **Cross-platform** - works on Mac, Windows, Linux
-   🌐 **Network transfers** between different machines
-   📁 **Any file type** - preserves original file extensions
-   ⚡ **Concurrent transfers** - handle multiple files simultaneously
-   🎯 **Simple CLI** - just two commands: send and receive

## Installation

1. **Clone the repository:**

    ```bash
    git clone <your-repo-url>
    cd godrop
    ```

2. **Build the application:**

    ```bash
    go build -o godrop
    ```

3. **Or run directly:**
    ```bash
    go run .
    ```

## Usage

### Basic Commands

**Start a receiver:**

```bash
go run . --mode receive --port 8888
```

**Send a file:**

```bash
go run . --mode send --path myfile.txt --port 8888
```

### Command Line Options

-   `--mode` - Operation mode: `send` or `receive` (required)
-   `--port` - Port number (default: 8888)
-   `--path` - File path for sending, or output path for receiving
-   `--host` - Target machine IP address (default: localhost)

## Examples

### Local Transfer (Same Machine)

**Terminal 1 - Start receiver:**

```bash
go run . --mode receive --port 8888
```

**Terminal 2 - Send file:**

```bash
go run . --mode send --path document.pdf --port 8888
```

### Network Transfer (Different Machines)

**Machine A (Receiver) - IP: 192.168.1.100:**

```bash
go run . --mode receive --port 8888
```

**Machine B (Sender):**

```bash
go run . --mode send --path video.mp4 --port 8888 --host 192.168.1.100
```

### File Types

Works with any file type:

```bash
# Images
go run . --mode send --path photo.jpg --port 8888 --host 192.168.1.50

# Videos
go run . --mode send --path movie.mp4 --port 8888 --host 192.168.1.50

# Documents
go run . --mode send --path presentation.pptx --port 8888 --host 192.168.1.50

# Archives
go run . --mode send --path backup.zip --port 8888 --host 192.168.1.50
```

## What You'll See

**Sender output:**

```
Sending video.mp4 (15728640 bytes) to 192.168.1.100:8888
[████████████████████] 100% - 1245.2 KB/s
✓ File sent successfully! Total: 15728640 bytes in 12.65s (1214.8 KB/s)
```

**Receiver output:**

```
Listening on port: 8888
Connection accepted
Receiving video.mp4 (15728640 bytes)
Starting to receive data...
[████████████████████] 100% - 1205.3 KB/s
✓ File received at received_video.mp4 successfully! Total: 15728640 bytes in 12.67s (1212.1 KB/s)
```

## Network Setup

### Find Your IP Address

**macOS/Linux:**

```bash
ifconfig | grep "inet " | grep -v 127.0.0.1
```

**Windows:**

```bash
ipconfig
```

### Firewall

Make sure the port (default 8888) is allowed through your firewall:

**macOS:**

```bash
# Firewall should allow incoming connections automatically
```

**Windows:**

```bash
# Add firewall rule for the port if needed
netsh advfirewall firewall add rule name="GoDrop" dir=in action=allow protocol=TCP localport=8888
```

**Linux:**

```bash
# UFW
sudo ufw allow 8888

# iptables
sudo iptables -A INPUT -p tcp --dport 8888 -j ACCEPT
```

## Tips

-   **Multiple receivers:** Each receiver can handle multiple senders simultaneously
-   **File naming:** Received files are prefixed with `received_` to avoid overwrites
-   **Large files:** Progress bars update every 5% to show smooth progress
-   **Network speed:** Transfer speed depends on your network connection
-   **Same network:** Both machines must be on the same WiFi/Ethernet network

## Troubleshooting

**Connection refused:**

-   Check if receiver is running
-   Verify IP address and port
-   Check firewall settings

**File not received:**

-   Check available disk space
-   Verify file permissions
-   Ensure network connectivity

**Slow transfers:**

-   Check network speed between machines
-   Try a different port if network is congested
-   Close other network-intensive applications

## Building for Different Platforms

```bash
# Windows
GOOS=windows GOARCH=amd64 go build -o godrop.exe

# Linux
GOOS=linux GOARCH=amd64 go build -o godrop-linux

# macOS
GOOS=darwin GOARCH=amd64 go build -o godrop-mac
```

## License

MIT License - feel free to use and modify!
