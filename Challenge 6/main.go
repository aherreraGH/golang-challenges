package main

import (
	"context"
	"errors"
	"fmt"
	"os"
	"sync"
	"time"

	"localpractice6.com/challenges/store"
	"localpractice6.com/challenges/types"
)

var s store.Store

func service(workerID int, ctx context.Context, tasks <-chan types.Task, wg *sync.WaitGroup) {
	defer wg.Done()

	for task := range tasks {
		fmt.Printf(
			"worker %d processing task %d\n",
			workerID,
			task.ID,
		)

		time.Sleep(500 * time.Millisecond)

		if _, err := s.Add(ctx, task.Title); err != nil {
			fmt.Printf(
				"worker %d failed to add task %d: %v\n",
				workerID,
				task.ID,
				err,
			)
		}
	}

	fmt.Printf("worker %d stopped\n", workerID)
}

/*
*
Your main function should:

Load the existing JSON file. [DONE]
Start an autosave goroutine. [DONE]
Start several worker goroutines. []
Concurrently add, complete, and read tasks.
Cancel the context after a controlled period or simulated shutdown.
Wait for all goroutines with sync.WaitGroup.
Perform one final save if the store is dirty.
Print atomic statistics.
*/
func main() {
	ctx, cancel := context.WithCancel(context.Background())

	// initialize data... remove file
	if err := s.Delete(ctx); err != nil && !errors.Is(err, os.ErrNotExist) {
		fmt.Printf("file not present, ignore... all good.")
	}
	// load data
	if err := s.Load(ctx); err != nil {
		fmt.Printf("Load failed: %v\n", err.Error())
		return
	}
	// start 3 workers
	tasks := make(chan types.Task)

	const workerCount = 5

	var workerWG sync.WaitGroup
	workerWG.Add(workerCount)
	var autosaveWG sync.WaitGroup
	autosaveWG.Add(1)

	// save data periodically
	fmt.Println("Periodic store start...")
	go s.RunAutoSave(ctx, time.Duration(5*time.Second), &autosaveWG)

	for workerID := 1; workerID <= workerCount; workerID++ {
		go service(workerID, ctx, tasks, &workerWG)
	}
	for x := 0; x < 10; x++ {
		if x%3 == 0 {
			items, err := s.List(ctx)
			if err != nil {
				fmt.Println("Failed to retrieve list: ", err.Error())
			}
			for i, item := range items {
				fmt.Printf("%d: %v\n", i, item.Title)
			}
		}
		time.Sleep(time.Duration(1 * time.Second))
		fmt.Printf("Submitting task %v...\n", x)
		tasks <- types.Task{
			ID:    x,
			Title: fmt.Sprintf("Task %v", x),
		}
	}
	close(tasks)

	workerWG.Wait()

	// Tell autosave to stop.
	cancel()

	autosaveWG.Wait()

	fmt.Println("all requests completed")

}
