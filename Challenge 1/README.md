## Challenge 1

Starting off with a simple use case, showing how to handle goroutines, channels, and worker pools.

Task is simple: 

1. Using a list of addresses and distance.
2. Submit the list items to be processed by the workers and show concurrency in action.
3. Ensure there are no deadlocks encountered.

## Testing for deadlocks

use:

```bash
go run -race .
```