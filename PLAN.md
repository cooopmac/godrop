-   `main.go` already handles CLI parsing, validation, and calls to `network.SendFile` or `network.ReceiveFile`.
-   `network/tcp.go` has stub functions for sending/receiving— we'll flesh these out.
-   We'll emphasize separation: Move code into packages like `internal/network`, `internal/file`, etc., early on.
-   Limit external packages: Use only Go standard library (e.g., no third-party progress bar; we'll make a simple one with `fmt`). If absolutely needed, we'll add minimal ones, but here we stick to built-ins for learning.

Each step has ✅ checklists with exact commands, code snippets, and explanations. No research needed—copy-paste and follow.

Run commands from your project root (`/Users/cooper/Developer/godrop`). Use `go run main.go` to test.

---

### 🔹 PHASE 1: Enhance CLI Setup (Build on Existing)

Your `main.go` is a great start. We'll refactor for better separation.

#### ✅ Step 1: Refactor CLI Logic into a Separate Function

-   [ ] Create a new file: `touch cli.go`
-   [ ] Add this to `cli.go`:

    ```
    package main

    import (
        \"flag\"
        \"fmt\"
        \"os\"
    )

    type Config struct {
        Mode string
        Port int
        Path string
    }

    func ParseConfig() Config {
        mode := flag.String(\"mode\", \"\", \"mode: send or receive\")
        port := flag.Int(\"port\", 8888, \"port: starting on port 8888\")
        path := flag.String(\"path\", \"\", \"path: path to the file\")

        flag.Parse()

        if *mode != \"send\" && *mode != \"receive\" {
            fmt.Println(\"Error: mode must be send or receive\")
            os.Exit(1)
        }

        if *port <= 0 || *port > 9999 {
            fmt.Println(\"Error: port must be between 1 and 9999\")
            os.Exit(1)
        }

        if *mode == \"send\" && *path == \"\" {
            fmt.Println(\"Error: path is required for send mode\")
            os.Exit(1)
        }

        if *mode == \"receive\" && *path == \"\" {
            *path = \"received_file.txt\"
            fmt.Println(\"Warning: path not provided, using default: \", *path)
        }

        return Config{*mode, *port, *path}
    }
    ```

-   [ ] Update `main.go`: Remove the flag parsing/validation, import your package if needed (but since it's same package, ok). In `main()`: `config := ParseConfig()` then use `config.Mode`, etc., to call network functions.
-   [ ] Test: `go run main.go --mode send --path file.txt --port 9000`
-   [ ] Why: Separates concerns; structs hold data neatly.

#### ✅ Step 2: Print Mode Messages (Already Partially Done)

-   [ ] In `main.go`, keep or add prints like \"Running in SEND mode\" before switching.

**Key Go Concepts Learned**:

-   Structs for config.
-   Separating code into files.

---

### 🔹 PHASE 2: Basic Networking (Build on tcp.go)

Flesh out `network/tcp.go` with TCP logic. Use built-in `net`, `os`, `io`.

#### ✅ Step 3: Implement TCP Receiver in tcp.go

-   [ ] Update `network/tcp.go` (replace ReceiveFile):

    ```
    func ReceiveFile(port int, path string) {
        listener, err := net.Listen(\"tcp\", fmt.Sprintf(\":%d\", port))
        if err != nil {
            fmt.Println(\"Error listening:\", err)
            return
        }
        defer listener.Close()
        fmt.Println(\"Listening on port\", port)

        conn, err := listener.Accept()
        if err != nil {
            fmt.Println(\"Error accepting:\", err)
            return
        }
        defer conn.Close()

        file, err := os.Create(path)
        if err != nil {
            fmt.Println(\"Error creating file:\", err)
            return
        }
        defer file.Close()

        _, err = io.Copy(file, conn)
        if err != nil {
            fmt.Println(\"Error receiving data:\", err)
        }
        fmt.Println(\"File received at\", path)
    }
    ```

-   [ ] Add imports at top: `\"net\" \"fmt\" \"os\" \"io\"`
-   [ ] Test: Create `file.txt`, run receive, then use a tool like `nc localhost 8888 < file.txt` to simulate send.

#### ✅ Step 4: Implement TCP Sender in tcp.go

-   [ ] Add to `network/tcp.go` (replace SendFile):

    ```
    func SendFile(port int, path string) {
        file, err := os.Open(path)
        if err != nil {
            fmt.Println(\"Error opening file:\", err)
            return
        }
        defer file.Close()

        conn, err := net.Dial(\"tcp\", fmt.Sprintf(\"localhost:%d\", port))
        if err != nil {
            fmt.Println(\"Error dialing:\", err)
            return
        }
        defer conn.Close()

        _, err = io.Copy(conn, file)
        if err != nil {
            fmt.Println(\"Error sending data:\", err)
        }
        fmt.Println(\"File sent from\", path)
    }
    ```

-   [ ] Test locally as in original plan.

#### ✅ Step 5: Test Locally

-   [ ] Echo \"test\" > file.txt
-   [ ] Terminal 1: `go run main.go --mode receive --port 8888`
-   [ ] Terminal 2: `go run main.go --mode send --path file.txt --port 8888`
-   [ ] Check `received_file.txt`.

**Key Go Concepts Learned**:

-   TCP with `net`.
-   File I/O streaming.

---

### 🔹 PHASE 3: Add Concurrency

Move to `internal/network` for better separation.

#### ✅ Step 6: Refactor to internal/network

-   [ ] `mkdir -p internal/network`
-   [ ] `mv network/tcp.go internal/network/tcp.go`
-   [ ] Update package name in tcp.go to `network`
-   [ ] In main.go, change import to `\"godrop/internal/network\"`
-   [ ] Update calls if needed.

#### ✅ Step 7: Handle Multiple Connections

-   [ ] In ReceiveFile, add loop:
    ```
    for {
        conn, err := listener.Accept()
        if err != nil { continue }
        go handleConnection(conn, path)  // New func
    }
    ```
-   [ ] Add:
    ```
    func handleConnection(conn net.Conn, path string) {
        defer conn.Close()
        file, err := os.Create(path)  // Note: Overwrites; improve later with unique names
        // ... rest of copy logic
    }
    ```

#### ✅ Step 8: Add Logging

-   [ ] In handleConnection: Add fmt.Println for accept/close/size.

**Key Go Concepts Learned**:

-   Goroutines.
-   Package refactoring.

---

### 🔹 PHASE 4: Add Encryption (TLS)

Use built-in `crypto/tls`.

#### ✅ Step 9: Generate Cert and Add TLS

-   [ ] Run: `go run $GOROOT/src/crypto/tls/generate_cert.go --host localhost`
-   [ ] Update ReceiveFile to use `tls.Listen` with config loading cert/key.
-   [ ] Update SendFile to use `tls.Dial` with InsecureSkipVerify.

#### ✅ Step 10: Test

-   [ ] Repeat local tests.

**Key Go Concepts Learned**:

-   TLS encryption.

---

### 🔹 PHASE 5: Improve Experience

Use built-ins only.

#### ✅ Step 11: Simple Progress (No External Pkg)

-   [ ] In SendFile/ReceiveFile, use a loop with `io.ReadAtLeast` or similar, and fmt.Printf(\"%d%% done\\n\", progress) every 10%.

#### ✅ Step 12: Compression with gzip

-   [ ] Import `\"compress/gzip\"`, wrap writers/readers.

#### ✅ Step 13: Integrity with SHA-256

-   [ ] Import `\"crypto/sha256\" \"encoding/hex\"`, compute/send/verify hash.

**Key Go Concepts Learned**:

-   Built-in compression/hashing.

---

### 🔹 PHASE 6: Optional Discovery

-   [ ] Use `net` for UDP broadcast.

---

### 🔹 PHASE 7: Polish

-   [ ] Refactor more (e.g., internal/file for I/O).
-   [ ] Add tests with `testing`.
-   [ ] Build binaries: `GOOS=linux go build`.
-   [ ] Update README.md with usage.

Track progress with checklists!"
