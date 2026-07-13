package main

import (
	"fmt"
	"sync"
	"time"
)

type Request struct {
	ID       int
	Value    int
	Response chan int
}

func service(requests <-chan Request, shutdown <-chan struct{}, wg *sync.WaitGroup) {
	defer wg.Done()

	for {
		select {
		case <-shutdown:
			fmt.Println("service shutting down")
			return

		case request, ok := <-requests:
			if !ok {
				fmt.Println("Service was shutdown...")
				return
			}
			time.Sleep(500 * time.Millisecond)

			fmt.Printf("processing request %d\n", request.ID)
			request.Response <- request.Value * request.Value
		}
	}
}

/**
Scenario

A dispatcher sends requests to a service goroutine. Each request includes a response channel so the caller can wait for its own response.

The program should process every request and then shut down cleanly.

Lifecycle:

request goroutine:
    send request
    receive response
    repeat
    close requests

service:
    process requests
    detect requests closure
    return
    wg.Done()

main:
    wg.Wait()
    print completion


*/

func main() {
	requests := make(chan Request)
	shutdown := make(chan struct{})

	var wg sync.WaitGroup

	wg.Add(1)
	go service(requests, shutdown, &wg)

	const requestCount = 5

	go func() {
		for requestID := 1; requestID <= requestCount; requestID++ {
			response := make(chan int)

			requests <- Request{
				ID:       requestID,
				Value:    requestID,
				Response: response,
			}

			result := <-response
			fmt.Printf("request %d returned %d\n", requestID, result)
		}

		close(requests)
	}()

	wg.Wait()

	fmt.Println("all requests completed")
}
