package challenge_1

import (
	"fmt"
	"time"
)

func challenge() {
	/**
	 * create channels to exit the loop and quit.
	 */
	quit := make(chan struct{})
	done := make(chan struct{})

	go func() {
		defer close(done)

		/**
		 * Use a ticker to control the looping.
		 */
		ticker := time.NewTicker(100 * time.Millisecond)
		defer ticker.Stop()

		count := 0

		for {
			select {
			// on quit, exit...
			case <-quit:
				fmt.Println("Exiting")
				return
			// using a ticker ensures the loop is not run forever and causing fast looping activity
			case <-ticker.C:
				count++
				fmt.Printf("count: %d\n", count)
			}
		}
	}()

	time.Sleep(500 * time.Millisecond)
	close(quit) // send quit in

	// all done, exit
	<-done
}
