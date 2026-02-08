# Build Your Own SMTP Server

A fully functional SMTP (Simple Mail Transfer Protocol) server built in Go, implementing RFC 5321.

## Overview

This project implements an SMTP server that can:
- Accept incoming email messages
- Validate SMTP command sequences
- Handle multiple concurrent connections
- Store received emails in memory

## What is SMTP?

**SMTP** (Simple Mail Transfer Protocol) is the standard protocol for sending emails across the Internet. It's a text-based, request-response protocol that operates over TCP.

### How Email Works

```
┌─────────────┐     SMTP      ┌─────────────┐     SMTP      ┌─────────────┐
│   Sender's  │ ────────────► │   Sender's  │ ────────────► │ Recipient's │
│   Email     │   (Submit)    │   Mail      │   (Relay)     │   Mail      │
│   Client    │               │   Server    │               │   Server    │
└─────────────┘               └─────────────┘               └─────────────┘
                                                                   │
                                                                   │ IMAP/POP3
                                                                   ▼
                                                            ┌─────────────┐
                                                            │ Recipient's │
                                                            │   Email     │
                                                            │   Client    │
                                                            └─────────────┘
```

### Key Terms

| Term | Description |
|------|-------------|
| **MUA** | Mail User Agent - Email client (Outlook, Gmail) |
| **MTA** | Mail Transfer Agent - Routes mail between servers |
| **MDA** | Mail Delivery Agent - Delivers to recipient's mailbox |
| **Envelope** | MAIL FROM + RCPT TO addresses used for routing |
| **Headers** | From, To, Subject displayed to user |

---

## Project Architecture

```
smtp/
├── main.go                 # Entry point, wires components together
├── config/
│   └── config.go           # Server configuration
├── server/
│   └── server.go           # TCP server, connection management
├── session/
│   ├── session.go          # SMTP state machine
│   └── handlers.go         # Command handlers
├── command/
│   ├── command.go          # Command types & constants
│   └── parser.go           # Command parsing logic
├── email/
│   ├── email.go            # Email struct
│   └── store.go            # Storage interface & implementation
└── go.mod
```

### Package Responsibilities

| Package | Responsibility |
|---------|----------------|
| `main` | Application entry point, dependency injection |
| `config` | Configuration management |
| `server` | TCP listener, connection lifecycle |
| `session` | SMTP protocol state machine |
| `command` | Command parsing and types |
| `email` | Email data model and storage |

---

## SMTP Protocol

### The State Machine

SMTP follows a strict state machine. Commands must be issued in order:

```
┌──────────────┐
│   CONNECTED  │ ← Initial state after TCP connection
└──────┬───────┘
       │ HELO/EHLO
       ▼
┌──────────────┐
│   GREETED    │ ← Ready to receive MAIL FROM
└──────┬───────┘
       │ MAIL FROM
       ▼
┌──────────────┐
│    MAIL      │ ← Sender specified, need recipients
└──────┬───────┘
       │ RCPT TO (can repeat)
       ▼
┌──────────────┐
│    RCPT      │ ← At least one recipient, ready for DATA
└──────┬───────┘
       │ DATA
       ▼
┌──────────────┐
│    DATA      │ ← Reading message body until "."
└──────┬───────┘
       │ "." (end of message)
       ▼
┌──────────────┐
│   GREETED    │ ← Ready for another message
└──────────────┘
```

### SMTP Commands

| Command | Syntax | Description |
|---------|--------|-------------|
| `HELO` | `HELO <domain>` | Identify client (basic) |
| `EHLO` | `EHLO <domain>` | Identify client (extended) |
| `MAIL FROM` | `MAIL FROM:<address>` | Specify sender |
| `RCPT TO` | `RCPT TO:<address>` | Specify recipient (repeatable) |
| `DATA` | `DATA` | Start message body |
| `QUIT` | `QUIT` | End session |
| `RSET` | `RSET` | Reset current transaction |
| `NOOP` | `NOOP` | No operation (keep-alive) |

### Response Codes

| Code | Category | Meaning |
|------|----------|---------|
| 2xx | Success | Action completed |
| 3xx | Intermediate | More input needed |
| 4xx | Temporary Failure | Try again later |
| 5xx | Permanent Failure | Don't retry |

Common codes:
- `220` - Service ready (greeting)
- `221` - Closing connection
- `250` - OK
- `354` - Start mail input
- `500` - Syntax error
- `503` - Bad sequence of commands

### Example SMTP Session

```
S: 220 localhost SMTP ready           ← Server greeting
C: HELO myclient.com                  ← Client identifies
S: 250 Hello myclient.com             ← Server acknowledges
C: MAIL FROM:<alice@sender.com>       ← Specify sender
S: 250 OK
C: RCPT TO:<bob@recipient.com>        ← Specify recipient
S: 250 OK
C: DATA                               ← Start message
S: 354 Start mail input; end with <CRLF>.<CRLF>
C: From: alice@sender.com             ← Message headers
C: To: bob@recipient.com
C: Subject: Hello!
C:                                    ← Blank line
C: Hi Bob,                            ← Message body
C: How are you?
C: .                                  ← End of message
S: 250 OK: message queued as 123456
C: QUIT                               ← End session
S: 221 Bye
```

---

## Key Concepts Learned

### 1. Text-Based Protocol Parsing

SMTP uses simple text commands terminated by CRLF (`\r\n`):

```go
// Reading a command
line, err := reader.ReadString('\n')
line = strings.TrimRight(line, "\r\n")

// Sending a response
fmt.Fprintf(conn, "%d %s\r\n", code, message)
```

### 2. State Machine Pattern

Managing protocol state with explicit states and transitions:

```go
type State int

const (
    StateConnected State = iota
    StateGreeted
    StateMail
    StateRcpt
    StateData
)

func (s *Session) handleMailFrom(cmd *Command) error {
    // Validate state FIRST
    if s.state != StateGreeted {
        return s.reply(503, "Bad sequence of commands")
    }
    
    // Process command
    s.mailFrom = cmd.Args[0]
    s.state = StateMail
    
    return s.reply(250, "OK")
}
```

### 3. Dot-Stuffing

Lines starting with `.` in the message body are "stuffed" with an extra dot:

```go
// When receiving (unstuffing)
if strings.HasPrefix(line, ".") {
    if line == "." {
        break  // End of message
    }
    line = line[1:]  // Remove stuffed dot
}
```

### 4. Separation of Concerns

Clean architecture with single-responsibility packages:

```go
// Each package has one job
server/    → TCP connections
session/   → Protocol logic
command/   → Parsing
email/     → Data model & storage
config/    → Configuration
```

### 5. Interface-Based Design

Using interfaces for extensibility:

```go
// Store interface allows different implementations
type Store interface {
    Save(email *Email) error
    Get(id string) (*Email, error)
    List() ([]*Email, error)
}

// Can implement MemoryStore, FileStore, DatabaseStore, etc.
type MemoryStore struct { ... }
type FileStore struct { ... }
```

### 6. Graceful Shutdown

Handling signals for clean termination:

```go
sigCh := make(chan os.Signal, 1)
signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)

go func() {
    <-sigCh
    fmt.Println("Shutting down...")
    server.Close()
    os.Exit(0)
}()
```

### 7. Concurrent Connection Handling

Each connection runs in its own goroutine:

```go
for {
    conn, err := listener.Accept()
    if err != nil {
        continue
    }
    go s.handleConnection(conn)  // Concurrent handling
}
```

---

## Usage

### Running the Server

```bash
# Build and run
go build -o smtp-server .
./smtp-server

# Or run directly
go run main.go
```

Output:
```
🚀 Starting SMTP Server...
📬 SMTP server listening on 0.0.0.0:2525
```

### Testing with Netcat

```bash
nc localhost 2525
```

### Testing with Telnet

```bash
telnet localhost 2525
```

### Sending a Test Email

```
HELO test
MAIL FROM:<sender@example.com>
RCPT TO:<recipient@example.com>
DATA
From: sender@example.com
To: recipient@example.com
Subject: Test Email

This is a test message.
.
QUIT
```

---

## Configuration

Default configuration in `config/config.go`:

```go
Host:     "0.0.0.0"   // Listen on all interfaces
Port:     "2525"      // Non-privileged port
Hostname: "localhost" // Server hostname for greeting
```

To use a different port, modify `config/config.go` or extend to read from environment variables.

---

## Extending the Server

### Adding File-Based Storage

```go
// email/file_store.go
type FileStore struct {
    directory string
}

func (s *FileStore) Save(email *Email) error {
    filename := filepath.Join(s.directory, email.ID+".eml")
    return os.WriteFile(filename, []byte(email.Content), 0644)
}
```

### Adding EHLO Capabilities

```go
func (s *Session) handleEhlo(cmd *command.Command) error {
    s.state = StateGreeted
    s.resetEnvelope()
    
    // Multi-line response
    s.conn.Write([]byte("250-localhost Hello\r\n"))
    s.conn.Write([]byte("250-SIZE 10485760\r\n"))
    s.conn.Write([]byte("250-8BITMIME\r\n"))
    s.conn.Write([]byte("250 HELP\r\n"))
    
    return nil
}
```

---

## Ports Reference

| Port | Name | Usage |
|------|------|-------|
| 25 | SMTP | Server-to-server relay |
| 465 | SMTPS | SMTP over TLS (deprecated) |
| 587 | Submission | Client-to-server (recommended) |
| 2525 | Alternative | When 587 is blocked |

We use **2525** to avoid needing root privileges (ports < 1024 require root).

---

## Security Considerations

This is a **learning implementation**. For production use, you would need:

- [ ] **TLS/STARTTLS** - Encrypt connections
- [ ] **Authentication** - AUTH PLAIN, AUTH LOGIN
- [ ] **SPF/DKIM/DMARC** - Email authentication
- [ ] **Rate Limiting** - Prevent abuse
- [ ] **Input Validation** - Sanitize addresses
- [ ] **Relay Restrictions** - Prevent open relay

---

## References

- [RFC 5321 - SMTP](https://tools.ietf.org/html/rfc5321) - Current SMTP specification
- [RFC 5322 - Internet Message Format](https://tools.ietf.org/html/rfc5322) - Email format
- [RFC 821 - Original SMTP](https://tools.ietf.org/html/rfc821) - Historical reference
- [Coding Challenges - SMTP](https://codingchallenges.fyi/challenges/challenge-smtp)

---

## What's Next?

Possible enhancements:
- Implement STARTTLS for encryption
- Add AUTH command for authentication
- Create a web interface to view received emails
- Implement email forwarding/relay
- Add persistent storage with SQLite

