## Challenge 1: Concurrent Package Processor

Build a program that processes package-delivery requests concurrently.

Each package contains:

```go
type Package struct {
	ID       int
	Address  string
	Distance int
}
```

Create three worker goroutines. Each worker receives packages through a shared jobs channel.

Processing time should simulate travel:

```go
time.Sleep(time.Duration(pkg.Distance) * 100 * time.Millisecond)
```

After processing, the worker sends a result through a results channel:

```go
type DeliveryResult struct {
	PackageID int
	WorkerID  int
	Status    string
}
```

__Requirements__

- [ ] Create at least six packages.
- [ ] Start exactly three worker goroutines.
- [ ] Send packages into a shared jobs channel.
- [ ] Close the jobs channel after all packages are submitted.

__Each worker must:__

- [ ] Receive packages until the jobs channel closes.
- [ ] Simulate processing time.
- [ ] Send a DeliveryResult to the results channel.
- [ ] Close the results channel only after all workers finish.
- [ ] Print results as they arrive.
- [ ] Do not assume results will arrive in package-ID order.

__Example output__
```go
Worker 2 delivered package 2
Worker 1 delivered package 1
Worker 3 delivered package 3
Worker 2 delivered package 4
Worker 1 delivered package 6
Worker 3 delivered package 5
All packages processed
```

The order can vary.

Concepts exercised
* Worker pools
* Shared receive-only job channels
* Result channels
* sync.WaitGroup
* Closing channels safely
* Reading with range

__Starter code__

```go
package main

import (
	"fmt"
	"sync"
	"time"
)

type Package struct {
	ID       int
	Address  string
	Distance int
}

type DeliveryResult struct {
	PackageID int
	WorkerID  int
	Status    string
}

func worker(
	workerID int,
	jobs <-chan Package,
	results chan<- DeliveryResult,
	wg *sync.WaitGroup,
) {
	// Implement this.
}

func main() {
	packages := []Package{
		{ID: 1, Address: "100 Main Street", Distance: 3},
		{ID: 2, Address: "200 Oak Avenue", Distance: 1},
		{ID: 3, Address: "300 Pine Road", Distance: 4},
		{ID: 4, Address: "400 Lake Drive", Distance: 2},
		{ID: 5, Address: "500 Hill Street", Distance: 5},
		{ID: 6, Address: "600 River Road", Distance: 2},
	}

	jobs := make(chan Package)
	results := make(chan DeliveryResult)

	var wg sync.WaitGroup

	// Start three workers.

	// Submit packages.

	// Close results after all workers finish.

	// Print results.
	fmt.Println("All packages processed")
}
```

__Bonus__

Add a Failed bool field to DeliveryResult. Consider any package with a distance greater than four to have failed delivery.