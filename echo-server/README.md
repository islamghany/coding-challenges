# Echo Server

A TCP/UDP Echo Server implementation following [RFC 862](https://datatracker.ietf.org/doc/html/rfc862), built as part of the [Coding Challenges](https://codingchallenges.fyi/challenges/challenge-echo).

## What is an Echo Server?

An Echo Server is one of the simplest network services defined in internet standards. Its purpose is straightforward:

> **Whatever data a client sends to the server, the server sends back (echoes) the exact same data.**

Think of it like talking into a canyon - your voice bounces back exactly as you spoke it.

### Why Does It Exist?

Echo servers are **diagnostic tools** used by network engineers to:

1. **Test network connectivity** - "Can I reach this server?"
2. **Measure round-trip time (RTT)** - "How long does data take to travel there and back?"
3. **Verify data integrity** - "Is the data being transmitted correctly?"
4. **Debug network applications** - A baseline for testing client implementations

It's also the **"Hello World" of network programming** - understanding echo servers means understanding the fundamentals of socket programming.

---

## RFC 862 - The Echo Protocol

The Echo Protocol is defined in [RFC 862](https://datatracker.ietf.org/doc/html/rfc862), published in 1983. Key points:

- **Standard Port**: 7
- **Supports both TCP and UDP**
- **TCP**: Connection-oriented, echoes until client disconnects
- **UDP**: Connectionless, echoes each datagram independently

---

## TCP vs UDP

This implementation supports both protocols. Here's how they differ:

| Aspect | TCP | UDP |
|--------|-----|-----|
| **Connection** | Connection-oriented (3-way handshake) | Connectionless |
| **Reliability** | Guaranteed delivery, ordered packets | Best-effort, may lose packets |
| **Data Model** | Stream of bytes | Individual datagrams |
| **State** | Server maintains connection state | Stateless |
| **Use Case** | When reliability matters | When speed matters |
| **Go Function** | `net.Listen()` → `Accept()` | `net.ListenPacket()` |
| **Read/Write** | `conn.Read()` / `conn.Write()` | `ReadFrom()` / `WriteTo()` |

### TCP Echo Flow
```
Client                    Server
   |                         |
   |------- SYN ------------>|  (Connection setup)
   |<------ SYN-ACK ---------|
   |------- ACK ------------>|
   |                         |
   |------- "Hello" -------->|  (Data exchange)
   |<------ "Hello" ---------|  (Echo)
   |------- "World" -------->|
   |<------ "World" ---------|
   |                         |
   |------- FIN ------------>|  (Connection teardown)
   |<------ ACK -------------|
```

### UDP Echo Flow
```
Client                    Server
   |                         |
   |------- "Hello" -------->|  (Single datagram)
   |<------ "Hello" ---------|  (Echo - includes sender address)
   |                         |
   |------- "World" -------->|  (Another independent datagram)
   |<------ "World" ---------|
```

---

## Implementation Details

### Buffer Sizes

```go
const (
    TCPBufferSize = 4096   // 4KB - good balance for TCP streams
    UDPBufferSize = 65535  // 64KB - max UDP datagram size
)
```

**Why these sizes?**

- **TCP (4KB)**: TCP is a stream protocol. Data arrives continuously, and we read it in chunks. 4KB aligns with common OS page sizes and provides a good balance between memory usage and syscall efficiency.

- **UDP (65KB)**: UDP is datagram-based. Each `ReadFrom()` returns exactly one complete datagram. If our buffer is smaller than the incoming datagram, **data is lost**. 65535 bytes is the maximum UDP datagram size, so we can handle any valid UDP packet.

### Concurrency Model

**TCP**: One goroutine per connection
```go
for {
    conn, err := ln.Accept()
    // ...
    go handleConnection(ctx, conn)  // New goroutine for each client
}
```

**UDP**: Single loop (stateless)
```go
for {
    n, addr, err := ln.ReadFrom(buf)
    // ...
    ln.WriteTo(buf[:n], addr)  // Echo back to sender
}
```

UDP doesn't need goroutines for basic operation because there's no connection state - each datagram is independent.

### Graceful Shutdown

The server implements graceful shutdown using:

1. **Signal Handling**: Catches `SIGINT` (Ctrl+C) and `SIGTERM`
2. **Context Cancellation**: Propagates shutdown signal to all goroutines
3. **WaitGroup**: Tracks active connections and waits for them to finish
4. **Read Deadlines**: Allows periodic checking of context cancellation

```go
// Signal handler
quit := make(chan os.Signal, 1)
signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

go func() {
    <-quit
    cancel()  // Cancel context → signals all handlers to stop
}()
```

**The Deadline Pattern**: Since `Read()` is a blocking call, we can't check `ctx.Done()` while blocked. The solution:

```go
conn.SetReadDeadline(time.Now().Add(20 * time.Second))
n, err := conn.Read(buf)
if err != nil {
    if netErr, ok := err.(net.Error); ok && netErr.Timeout() {
        select {
        case <-ctx.Done():
            return  // Shutdown requested
        default:
            continue  // Keep waiting for data
        }
    }
}
```

This pattern:
1. Sets a short read deadline
2. When timeout occurs, checks if shutdown was requested
3. If yes → exit gracefully
4. If no → continue waiting for data

---

## Usage

### Build
```bash
go build -o echoserver
```

### Run TCP Server (default)
```bash
./echoserver
# or
go run main.go
```

### Run UDP Server
```bash
./echoserver -udp
# or
go run main.go -udp
```

### Testing with netcat

**TCP:**
```bash
nc localhost 7777
Hello, World!
# Server echoes: Hello, World!
```

**UDP:**
```bash
nc -u localhost 7777
Hello, UDP!
# Server echoes: Hello, UDP!
```

### Graceful Shutdown
Press `Ctrl+C` to initiate graceful shutdown. The server will:
1. Stop accepting new connections
2. Wait for active connections to finish (or timeout)
3. Clean up resources
4. Exit

---

## Project Structure

```
echo-server/
├── main.go          # Main implementation
├── go.mod           # Go module file
└── README.md        # This file
```

---

## Key Concepts Learned

1. **Socket Programming**: Using Go's `net` package to create TCP/UDP servers
2. **Protocol Implementation**: Reading and implementing RFC specifications
3. **Concurrency Patterns**: Goroutines, WaitGroups, channels
4. **Context-based Cancellation**: Propagating shutdown signals across goroutines
5. **Graceful Shutdown**: Properly cleaning up resources on termination
6. **Network I/O**: Buffer management, read deadlines, error handling
7. **TCP vs UDP**: Understanding when to use each protocol

---

## References

- [RFC 862 - Echo Protocol](https://datatracker.ietf.org/doc/html/rfc862)
- [Coding Challenges - Echo Server](https://codingchallenges.fyi/challenges/challenge-echo)
- [Go net package documentation](https://pkg.go.dev/net)

