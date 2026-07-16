## Challenge 4: The Request Dispatcher

### Scenario

A dispatcher sends requests to a service goroutine. Each request includes a response channel so the caller can wait for its own response.

The program should process every request and then shut down cleanly.

### Broken code

```go
package main

import (
	"fmt"
	"sync"
)

type Request struct {
	ID       int
	Value    int
	Response chan int
}

func service(
	requests <-chan Request,
	shutdown <-chan struct{},
	wg *sync.WaitGroup,
) {
	defer wg.Done()

	for {
		select {
		case <-shutdown:
			fmt.Println("service shutting down")
			return

		case request := <-requests:
			fmt.Printf("processing request %d\n", request.ID)
			request.Response <- request.Value * request.Value
		}
	}
}

func main() {
	requests := make(chan Request)
	shutdown := make(chan struct{})

	var wg sync.WaitGroup

	wg.Add(1)
	go service(requests, shutdown, &wg)

	const requestCount = 5

	for requestID := 1; requestID <= requestCount; requestID++ {
		response := make(chan int)

		requests <- Request{
			ID:       requestID,
			Value:    requestID,
			Response: response,
		}

		defer close(response)

		result := <-response
		fmt.Printf("request %d returned %d\n", requestID, result)
	}

	close(requests)
	close(shutdown)

	wg.Wait()

	fmt.Println("all requests completed")
}
```

### Constraints

Your correction must:

- [ ] Keep the per-request response-channel design.
- [ ] Keep requests unbuffered.
- [ ] Keep a select in the service loop.
- [ ] Allow the service to stop when request submission ends.
- [ ] Avoid time.Sleep.
- [ ] Avoid panic from sending on or closing a closed channel.
- [ ] Avoid an infinite loop after requests is closed.
- [ ] Avoid relying on shutdown winning a random select.