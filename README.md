__note:__ Some challenges with golang. Provided by ChatGPT to help with code assessments.

## Engineering Principles

Throughout these exercises I focused on:

- Idiomatic Go
- Readable and maintainable code
- Proper error handling
- Thread-safe concurrent programming
- Graceful shutdown patterns
- Separation of concerns
- Unit testing
- Minimal global state

# Golang Engineering Challenges

This repository contains a collection of practical Go engineering challenges completed to strengthen backend software engineering skills. Each challenge focuses on writing maintainable, testable, and idiomatic Go while exploring concepts commonly encountered in distributed systems, cloud-native applications, and production software.

## Challenges

### Challenge 1 – Concurrent Worker Pool

A multi-worker job processing system demonstrating safe concurrent execution and synchronization.

**Concepts Covered**
- Goroutines
- Channels
- Worker Pools
- WaitGroups
- Graceful Shutdown
- Concurrent Job Processing

---

### Challenge 2 – Sensor Multiplexing

Simulates multiple sensors publishing data concurrently while a central processor consumes and reacts to incoming events.

**Concepts Covered**
- Goroutines
- Channels
- `select`
- `time.Ticker`
- Timeouts
- Graceful Shutdown
- Event Multiplexing

---

### Challenge 3 – Concurrent Job Queue

Implements a producer/consumer job queue while avoiding deadlocks and coordinating worker completion.

**Concepts Covered**
- Producer / Consumer Pattern
- Buffered Channels
- WaitGroups
- Deadlock Prevention
- Channel Lifecycle Management
- Concurrent Synchronization

---

### Challenge 4 – Go RPC Command-Line Client

Implements a CLI application that communicates with an RPC server using Go's `net/rpc` package.

**Concepts Covered**
- `net/rpc`
- HTTP RPC Transport
- Command-Line Interfaces
- Flag Parsing
- Dependency Injection
- Unit Testing
- Mocking
- Error Handling

---

### Challenge 5 – Context Cancellation

Builds a concurrent application that uses `context.Context` to coordinate cancellation and graceful shutdown of running goroutines.

**Concepts Covered**
- `context.Context`
- Cancellation Propagation
- Goroutine Coordination
- Timeouts
- Graceful Shutdown
- Concurrent Design Patterns

---

### Challenge 6 – Thread-Safe Task Store

Implements a concurrent task store supporting safe reads, writes, periodic autosave, and operational metrics.

**Concepts Covered**
- `sync.Mutex`
- `sync.RWMutex`
- Atomic Counters
- `context.Context`
- `time.Ticker`
- Background Workers
- Thread Safety
- Concurrent State Management
- Graceful Shutdown

### Challenge 7 - Library Inventory Management System

**Concepts covered:**

- REST API design using Go's `net/http`
- Client/server architecture
- JSON-based data persistence
- CLI application development
- Barcode generation and image decoding (Code 128)
- Inventory state management
- HTTP request/response handling
- Struct modeling and JSON serialization
- Package organization and separation of concerns
- Date/time calculations and business rules
- State transitions (checkout/checkin workflows)
- Error handling and API error responses
- Concurrency management:
  - `sync.RWMutex` for protecting shared inventory state
  - `sync.Mutex` concepts for exclusive write operations
  - `sync/atomic` for lock-free statistics counters
- Thread-safe service design
- Internal helper methods for lock ownership (`findUnsafe`)
- Data ownership and defensive copying (`cloneBook`)

## Initial structure

__Note:__ This evovled a bit from the initial setup.

```bash
library-inventory/
│
├── cmd/
│   ├── kiosk/
│   │   └── main.go
│   │
│   └── circulation-server/
│       └── main.go
│
├── internal/
│   ├── api/
│   │   ├── handlers.go
│   │   ├── router.go
│   │   └── client.go
│   │
│   ├── inventory/
│   │   ├── store.go
│   │   ├── service.go
│   │   └── persistence.go
│   │
│   ├── checkout/
│   │   ├── service.go
│   │   └── calculations.go
│   │
│   ├── barcode/
│   │   └── decoder.go
│   │
│   ├── report/
│   │   └── report.go
│   │
│   └── models/
│       ├── book.go
│       ├── checkout.go
│       ├── user.go
│       └── responses.go
│
├── data/
│
├── barcodes/
│
├── go.mod
├── go.sum
└── README.md
```

## About Mutex

```bash
RWMutex
-------
Protects inventory state:
- books
- status
- borrowers
- due dates


atomic
------
Protects independent counters:
- requests
- checkouts
- checkins


context
-------
Controls:
- shutdown
- cancellation
- request lifecycle
```


## Future Challenges

Planned topics include:

- YAML parsing and validation
- Dependency graph resolution
- Topological sorting
- HTTP middleware
- OCI image metadata
- File packaging
- Concurrent file processing
- Kubernetes client interactions
