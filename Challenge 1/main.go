package main

import (
	"fmt"
	"sync"
	"time"

	"localpractice1.com/challenges/types"
)

func worker(workerID int, jobs <-chan types.Package, results chan<- types.DeliveryResult, wg *sync.WaitGroup) {
	defer wg.Done()

	for pkg := range jobs {
		time.Sleep(500 * time.Millisecond)

		results <- types.DeliveryResult{
			PackageID: pkg.ID,
			WorkerID:  workerID,
			Status:    fmt.Sprintf("Delivered %d miles to %s", pkg.Distance, pkg.Address),
		}
	}
}

func main() {
	packages := []types.Package{
		{ID: 1, Address: "100 Main Street", Distance: 3},
		{ID: 2, Address: "200 Oak Avenue", Distance: 1},
		{ID: 3, Address: "300 Pine Road", Distance: 4},
		{ID: 4, Address: "400 Lake Drive", Distance: 2},
		{ID: 5, Address: "500 Hill Street", Distance: 5},
		{ID: 6, Address: "600 River Road", Distance: 2},
	}

	jobs := make(chan types.Package)
	results := make(chan types.DeliveryResult)

	var wg sync.WaitGroup

	// Start three workers.
	for x := 0; x < 3; x++ {
		wg.Add(1)
		go worker(x, jobs, results, &wg)
	}
	// Submit packages.
	go func() {
		for _, pkg := range packages {
			jobs <- pkg
		}
		close(jobs)
	}()

	// Close results after all workers finish.
	go func() {
		wg.Wait()
		close(results)
	}()

	for result := range results {
		fmt.Println(result.Status)
	}

	// Print results.
	fmt.Println("All packages processed")
}
