## Challenge: Concurrent JSON Task Store

Build a small Go service or CLI that manages tasks stored in a JSON file.

Multiple goroutines will concurrently:

- [ ] add tasks,
- [ ] read tasks,
- [ ] mark tasks complete,
- [ ] periodically save the current state to disk.

The program must shut down cleanly when its context.Context is canceled.

Where each topic fits
```go
sync.RWMutex
```
Protect the in-memory task collection:
```go
type Store struct {
    mu    sync.RWMutex
    tasks []Task
}
```
`RLock()` for listing or retrieving tasks.
`Lock()` for adding or modifying tasks.

Also ensure that only one goroutine writes the JSON file at a time.
```go
context.Context
```
Use context for:

- [ ] canceling worker goroutines,
- [ ] stopping the periodic file writer,
- [ ] preventing new operations during shutdown,
- [ ] optionally applying timeouts to store operations.

Workers should use a pattern such as:
```go
select {
case <-ctx.Done():
    return
case job := <-jobs:
    // process job
}
```

### sync/atomic

Track simple counters and flags without acquiring the store mutex:
```go
var reads atomic.Uint64
var writes atomic.Uint64
var dirty atomic.Bool
```
Examples:

- [ ] number of read operations,
- [ ] number of updates,
- [ ] number of successful disk saves,
- [ ] whether unsaved changes exist.

`atomic` should not protect the task slice or replace the mutex.

### Suggested requirements

The JSON structure could be:
```json
[
  {
    "id": 1,
    "title": "Review context cancellation",
    "completed": false
  }
]
```

Implement:
```go
type Task struct {
    ID        int    `json:"id"`
    Title     string `json:"title"`
    Completed bool   `json:"completed"`
}
```

And a store with methods resembling:

```go
func NewStore(filename string) *Store
func (s *Store) Load(ctx context.Context) error
func (s *Store) Add(ctx context.Context, title string) (Task, error)
func (s *Store) Complete(ctx context.Context, id int) error
func (s *Store) List(ctx context.Context) ([]Task, error)
func (s *Store) Save(ctx context.Context) error
func (s *Store) RunAutoSave(ctx context.Context, interval time.Duration)
```

### Concurrency scenario

Your main function should:

- [ ] Load the existing JSON file.
- [ ] Start an autosave goroutine.
- [ ] Start several worker goroutines.
- [ ] Concurrently add, complete, and read tasks.
- [ ] Cancel the context after a controlled period or simulated shutdown.
- [ ] Wait for all goroutines with sync.WaitGroup.
- [ ] Perform one final save if the store is dirty.
- [ ] Print atomic statistics.

#### Example output:
```bash
Loaded 5 tasks
Worker 1 added task 6
Worker 2 listed 6 tasks
Worker 3 completed task 2
Autosaved 6 tasks
Shutdown requested
Final save completed

Reads: 8
Updates: 12
Disk saves: 3
```

### Important correctness requirements

Run successfully with:

```bash
go test -race ./...
```

* Do not return the internal task slice directly. Return a copy so callers cannot mutate shared state after the lock is released.
* Do not hold the task mutex during slow disk I/O. Copy the tasks while locked, unlock, and then marshal/write the copy.
* Save through a temporary file and rename it so cancellation or a crash is less likely to corrupt the primary JSON file:
```bash
tasks.json.tmp → tasks.json
```

* Check `ctx.Err()` before beginning an operation and at sensible interruption points.
* Ensure canceled operations return an appropriate error rather than silently succeeding.

This is a manageable challenge, but it exercises the three concepts in realistic roles rather than forcing them into arbitrary usage. The most interesting parts will be coordinating dirty, avoiding races during snapshots, and making shutdown perform exactly one reliable final save.