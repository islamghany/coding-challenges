# Bitcask

A log-structured key-value storage engine implementation following the [Bitcask paper](https://riak.com/assets/bitcask-intro.pdf), built as part of the [Coding Challenges](https://codingchallenges.fyi/challenges/challenge-bitcask).

## What is Bitcask?

**Bitcask** is a log-structured hash table for fast key/value storage, originally created by Basho Technologies as the default storage engine for [Riak](https://riak.com/), a distributed NoSQL database.

The name comes from: **Bit** (data) + **Cask** (barrel/container)

### Core Design Principles

Instead of updating data in place (like traditional databases), Bitcask uses a brilliantly simple approach:

1. **Write**: Always **append** new data to the end of a log file (never modify existing data)
2. **Read**: Use an **in-memory hash table** (KeyDir) that points to where data lives on disk
3. **Delete**: Write a special "tombstone" marker (data isn't actually removed immediately)
4. **Compact**: Periodically clean up old/deleted entries in a background "merge" process

### Design Goals (from the original paper)

| Goal | How Bitcask Achieves It |
|------|------------------------|
| Low latency reads | Hash table lookup → single disk seek → done |
| Low latency writes | Always append (sequential I/O is fast!) |
| High throughput | Sequential writes saturate disk bandwidth |
| Handle data > RAM | Only keys in memory, values on disk |
| Crash recovery | Log files + CRC = verify & rebuild on startup |
| Easy backup | It's just a directory! Copy it. |
| Simple codebase | Clean, understandable implementation |

---

## Architecture

```
┌────────────────────────────────────────────────────────────────┐
│                         BITCASK ARCHITECTURE                    │
├────────────────────────────────────────────────────────────────┤
│                                                                 │
│   ┌──────────────────┐         ┌─────────────────────────────┐ │
│   │  In-Memory       │         │  On-Disk Log File           │ │
│   │  KeyDir          │         │  (Append-Only)              │ │
│   │                  │         │                             │ │
│   │  key1 → metadata─┼────────►│  [entry][entry][entry]...   │ │
│   │  key2 → metadata─┼────────►│                             │ │
│   │  key3 → metadata─┼────────►│                             │ │
│   │                  │         │                             │ │
│   └──────────────────┘         └─────────────────────────────┘ │
│          │                                                      │
│          │ O(1) lookup                                          │
│          ▼                                                      │
│   ┌──────────────────┐                                          │
│   │  Get("key2")     │                                          │
│   │  → Seek to pos   │                                          │
│   │  → Read value    │                                          │
│   │  → Return        │                                          │
│   └──────────────────┘                                          │
└────────────────────────────────────────────────────────────────┘
```

---

## Binary Log Entry Format

Each entry on disk follows this binary format:

```
┌───────────┬────────────┬──────────┬──────────┬─────────┬─────────┐
│ CRC (4B)  │ Timestamp  │ Key Size │ Val Size │   Key   │  Value  │
│  uint32   │   uint32   │  uint32  │  uint32  │  []byte │  []byte │
└───────────┴────────────┴──────────┴──────────┴─────────┴─────────┘
     4B          4B           4B         4B       varlen    varlen

Header: 16 bytes (fixed)
Total:  16 + len(key) + len(value) bytes
```

| Field | Size | Type | Description |
|-------|------|------|-------------|
| CRC | 4 bytes | uint32 | CRC-32 checksum of everything after CRC |
| Timestamp | 4 bytes | uint32 | Unix timestamp (seconds since 1970) |
| Key Size | 4 bytes | uint32 | Length of key in bytes |
| Value Size | 4 bytes | uint32 | Length of value in bytes (0 = tombstone) |
| Key | variable | []byte | The key data |
| Value | variable | []byte | The value data |

---

## KeyDir (In-Memory Index)

The KeyDir is a hash table that maps keys to their location on disk:

```go
type KeyDirEntry struct {
    FileID    uint32  // Which file contains the value
    ValuePos  uint64  // Byte offset where VALUE starts
    ValueSize uint32  // Size of value in bytes
    Timestamp uint32  // When entry was written
}

type KeyDir struct {
    Entries map[string]KeyDirEntry
    mu      sync.RWMutex  // Thread-safe access
}
```

**Performance:**
- **Write**: O(1) hash table update + O(1) file append
- **Read**: O(1) hash lookup + O(1) disk seek + O(1) read
- **Memory**: ~40 bytes per key (metadata only, not values)

---

## Operations

### Set (Write)

```
1. Create Entry with timestamp, key, value
2. Encode to binary format
3. Calculate CRC-32 checksum
4. Append to log file
5. Update KeyDir with new position
6. Sync to disk
```

### Get (Read)

```
1. Lookup key in KeyDir → O(1)
2. If not found → return error
3. Seek to ValuePos in file
4. Read ValueSize bytes
5. Return value
```

### Delete

```
1. Check key exists in KeyDir
2. Write tombstone entry (ValueSize = 0)
3. Remove key from KeyDir
4. Sync to disk
```

### Merge (Compaction)

```
1. Create new merge file
2. For each live key in KeyDir:
   a. Read value from old file
   b. Write entry to merge file
   c. Update KeyDir with new position
3. Delete old file
4. Rename merge file to database file
```

---

## Tombstones (Deletion)

In an append-only log, we can't remove data. Instead, we write a **tombstone**:

```
Normal entry:  [CRC][TS][KeyLen][ValLen=5][Key]["Value"]
Tombstone:     [CRC][TS][KeyLen][ValLen=0][Key][]  ← Empty value
```

During `loadKeyDir()`, tombstones remove keys from the index. The merge process permanently removes tombstones.

---

## Why Append-Only?

| Aspect | In-Place Update | Append-Only |
|--------|-----------------|-------------|
| Write speed | Random I/O (slow) | Sequential I/O (fast) |
| Disk space | Constant | Grows until compaction |
| Crash recovery | Complex (WAL needed) | Simple (CRC check) |
| Implementation | Complex | Simple |

**Sequential writes are 10-100x faster** than random writes on HDDs, and even on SSDs they reduce write amplification.

---

## CRC-32 Checksum

Each entry includes a CRC-32 checksum for data integrity:

```go
import "hash/crc32"

// When writing
payload := timestamp + keySize + valueSize + key + value
checksum := crc32.ChecksumIEEE(payload)

// When reading
if crc != crc32.ChecksumIEEE(payload) {
    return ErrChecksumMismatch  // Data corrupted!
}
```

CRC detects:
- Single bit flips
- Burst errors (consecutive bits)
- Partial writes from crashes

---

## Usage

### Build

```bash
cd bitcask
go build -o ccbitcask
```

### Commands

```bash
# Set a key
./ccbitcask -db ./database set <key> <value>

# Get a key
./ccbitcask -db ./database get <key>

# Delete a key
./ccbitcask -db ./database del <key>

# Compact the database (merge)
./ccbitcask -db ./database merge
```

### Examples

```bash
# Create some data
./ccbitcask -db ./database set name "Islam"
./ccbitcask -db ./database set age "25"
./ccbitcask -db ./database set city "Cairo"

# Read data
./ccbitcask -db ./database get name
# Islam

# Update data (appends new entry)
./ccbitcask -db ./database set name "Ghany"
./ccbitcask -db ./database get name
# Ghany

# Delete data
./ccbitcask -db ./database del city
./ccbitcask -db ./database get city
# Error: key not found

# Compact (removes old entries and tombstones)
./ccbitcask -db ./database merge
# Database merged successfully
```

---

## Project Structure

```
bitcask/
├── main.go              # CLI interface
├── go.mod               # Go module
├── README.md            # This file
└── bitcask/
    ├── bitcask.go       # Core Bitcask implementation
    ├── entry.go         # Binary entry encoding/decoding
    └── keydir.go        # In-memory hash table index
```

---

## Key Concepts Learned

### 1. Binary Encoding

Using `encoding/binary` for efficient data serialization:

```go
binary.LittleEndian.PutUint32(buf, value)  // Write
value := binary.LittleEndian.Uint32(buf)   // Read
```

**Little Endian** stores least significant byte first, matching most modern CPUs (x86, ARM).

### 2. CRC-32 Checksums

Error-detection codes for data integrity:

```go
checksum := crc32.ChecksumIEEE(data)
```

Fast, 4-byte output, detects accidental corruption (not cryptographically secure).

### 3. Append-Only Logs

Never modify existing data - always append:
- Crash-safe (partial writes don't corrupt existing data)
- Sequential I/O (maximum disk throughput)
- Simple recovery (scan from start, CRC validates entries)

### 4. In-Memory Indexing

Keep frequently-accessed metadata in RAM:
- Keys + positions in memory (~40 bytes per key)
- Values on disk (can be arbitrarily large)
- O(1) lookups via hash table

### 5. Tombstones for Deletion

Soft deletion in immutable logs:
- Mark as deleted (don't remove)
- Compaction removes permanently
- Allows recovery of "deleted" data before compaction

### 6. Log Compaction (Merge)

Garbage collection for append-only logs:
- Remove stale entries (superseded by updates)
- Remove tombstones (deleted keys)
- Reclaim disk space
- Rewrite only live data

### 7. File I/O in Go

```go
os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0644)
file.Seek(offset, io.SeekStart)
io.ReadFull(file, buf)
file.Write(data)
file.Sync()  // Flush to disk
```

### 8. Thread Safety

Using `sync.RWMutex` for concurrent access:

```go
mu.RLock()   // Multiple readers OK
mu.RUnlock()

mu.Lock()    // Exclusive write access
mu.Unlock()
```

---

## Bitcask vs Other Storage Engines

| Engine | Type | Pros | Cons |
|--------|------|------|------|
| **Bitcask** | Log-structured hash | Simple, fast, predictable | Keys must fit in RAM |
| **B-Tree** (SQLite) | Tree-based | Range queries, sorted | Write amplification |
| **LSM-Tree** (LevelDB) | Log-structured merge | Great write throughput | Complex compaction |

---

## Limitations

1. **All keys must fit in RAM** - KeyDir holds every key
2. **No range queries** - Hash table doesn't support ordered iteration
3. **Single file** - This implementation uses one file (production uses multiple)

---

## References

- [Bitcask Paper](https://riak.com/assets/bitcask-intro.pdf) - Original design document
- [Coding Challenges - Bitcask](https://codingchallenges.fyi/challenges/challenge-bitcask)
- [Go encoding/binary](https://pkg.go.dev/encoding/binary)
- [Go hash/crc32](https://pkg.go.dev/hash/crc32)

