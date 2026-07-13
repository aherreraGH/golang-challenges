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
