package main

import (
	"fmt"
	"math/rand/v2"
	"sync"
	"time"

	"localpractice2.com/challenges/types"
)

func sensor(id string, interval time.Duration, readings chan<- types.Reading, shutdown <-chan struct{}, wg *sync.WaitGroup) {
	defer wg.Done()

	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for {
		select {
		case <-shutdown:
			return

		case <-ticker.C:
			reading := types.Reading{
				SensorID:    id,
				Temperature: rand.Float64() * 100,
			}

			/**
			 * Prevent 2nd deadlock... when shutdown occurs.
			 */
			select {
			case readings <- reading:
				// Reading delivered successfully.

			case <-shutdown:
				// Monitor may already have exited.
				return
			}
		}
	}
}

func monitor(readings <-chan types.Reading, shutdown <-chan struct{}, done chan<- struct{}) {
	// Implement this using select and a timer.
	defer close(done)

	for {
		select {
		case <-shutdown:
			return

		case reading := <-readings:
			message := fmt.Sprintf(
				"Sensor %s: %.1f°F",
				reading.SensorID,
				reading.Temperature,
			)

			if reading.Temperature >= 90.0 {
				message += " - ALERT"
			}

			fmt.Println(message)
		}
	}
}

func main() {
	readings := make(chan types.Reading)
	shutdown := make(chan struct{})
	monitorDone := make(chan struct{})

	var wg sync.WaitGroup

	// Start three sensors with different intervals.
	var intervals []int = []int{500, 800, 3000}
	var names []string = []string{"A", "B", "C"}
	for x := 0; x < 3; x++ {
		wg.Add(1)
		go sensor(names[x], time.Duration(intervals[x])*time.Millisecond, readings, shutdown, &wg)
	}
	// Suggested intervals:
	// Sensor A: 500 milliseconds
	// Sensor B: 800 milliseconds
	// Sensor C: 1200 milliseconds

	go monitor(readings, shutdown, monitorDone)

	// Let the program run for approximately eight seconds.
	time.Sleep(time.Duration(8 * time.Second))

	// Signal all goroutines to stop.
	close(shutdown)

	// Wait for sensors and monitor to finish.
	wg.Wait()
	<-monitorDone
	fmt.Println("Program finished")
}
