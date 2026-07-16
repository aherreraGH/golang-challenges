## Challenge 3: The Results Collector

### Scenario

Several workers process jobs concurrently. Each worker sends a result back to the coordinator.

The coordinator should:

- [ ] Submit all jobs.
- [ ] Collect every result.
- [ ] Print the total.
- [ ] Exit cleanly.

### Broken code

```go
package main

import (
	"fmt"
	"sync"
)

type Job struct {
	ID    int
	Value int
}

type Result struct {
	JobID int
	Value int
}

func worker(
	id int,
	jobs <-chan Job,
	results chan<- Result,
	wg *sync.WaitGroup,
) {
	defer wg.Done()

	for job := range jobs {
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

	for jobID := 1; jobID <= jobCount; jobID++ {
		jobs <- Job{
			ID:    jobID,
			Value: jobID * 10,
		}
	}

	close(jobs)

	wg.Wait()
	close(results)

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
```

### Constraints

Your correction must preserve these characteristics:

- [ ] jobs and results remain unbuffered.
- [ ] There must still be exactly three worker goroutines.
- [ ] sWorkers must send results through the results channel.
- [ ] Do not use time.Sleep.
- [ ] Do not replace the channels with shared slices.
- [ ] Do not add arbitrary channel-buffer capacity to hide the problem.