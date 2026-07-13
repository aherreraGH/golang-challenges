package main

import (
	"fmt"
	"sync"
	"time"
)

type Job struct {
	ID    int
	Value int
}

type Result struct {
	JobID int
	Value int
}

func worker(id int, jobs <-chan Job, results chan<- Result, wg *sync.WaitGroup) {
	defer wg.Done()

	for job := range jobs {
		time.Sleep(500 * time.Millisecond)

		fmt.Printf("worker %d processing job %d\n", id, job.ID)

		results <- Result{
			JobID: job.ID,
			Value: job.Value * 2,
		}
	}
}

func main() {
	jobs := make(chan Job)
	results := make(chan Result)

	var wg sync.WaitGroup

	const workerCount = 3
	const jobCount = 8

	for workerID := 1; workerID <= workerCount; workerID++ {
		wg.Add(1)
		go worker(workerID, jobs, results, &wg)
	}

	go func() {
		for jobID := 1; jobID <= jobCount; jobID++ {
			jobs <- Job{
				ID:    jobID,
				Value: jobID * 10,
			}
		}
		close(jobs)
	}()

	go func() {
		wg.Wait()
		close(results)
	}()

	total := 0

	for result := range results {
		fmt.Printf(
			"job %d produced %d\n",
			result.JobID,
			result.Value,
		)

		total += result.Value
	}

	fmt.Printf("total: %d\n", total)
}
