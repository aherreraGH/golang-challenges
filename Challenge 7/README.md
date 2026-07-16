## Challenge #7: Distributed Library Inventory Kiosk

### Scenario

A small closed-network library uses barcode scanners at self-service kiosks.

When a barcode is scanned:

- [ ] The kiosk checks the local inventory service.
- [ ] If the book is available, it displays the book information.
- [ ] If the book is checked out, the kiosk contacts a second service to determine:
    - [ ] who checked it out,
    - [ ] when it was checked out,
    - [ ] when it is due,
    - [ ] how many days it has been checked out,
    - [ ] whether it is overdue.
- [ ] At shutdown, the application prints a final inventory report.

The inventory begins in a pristine state:

- [ ] all books are on the shelf,
- [ ] no checkout records exist,
- [ ] all users are already registered.

### Starter assets

I generated three real Code 128 barcode images and the initial JSON datasets:


The archive contains:
```bash
go-library-inventory-challenge/
├── README.md
├── barcodes/
│   ├── LIB-GO-001-A7.png
│   ├── LIB-MATH-88-Z3.png
│   └── LIB-SYS-204-K9.png
└── data/
    ├── inventory.json
    ├── users.json
    └── checkouts.json
```
The barcode values are:
```bash
LIB-GO-001-A7
LIB-SYS-204-K9
LIB-MATH-88-Z3
```

### Recommended architecture

Build two separate Go applications.
```bash
cmd/
├── kiosk/
│   └── main.go
└── circulation-server/
    └── main.go
```

### Suggested project structure:

```bash
library-inventory/
├── cmd/
│   ├── kiosk/
│   │   └── main.go
│   └── circulation-server/
│       └── main.go
├── internal/
│   ├── inventory/
│   │   ├── service.go
│   │   └── store.go
│   ├── circulation/
│   │   ├── client.go
│   │   ├── server.go
│   │   └── store.go
│   ├── barcode/
│   │   └── decoder.go
│   ├── report/
│   │   └── report.go
│   └── types/
│       └── types.go
├── data/
│   ├── inventory.json
│   ├── users.json
│   └── checkouts.json
├── barcodes/
├── go.mod
└── README.md
```

### Application 1: Kiosk

The kiosk is the primary application.

It should support commands such as:

```bash
scan barcodes/LIB-GO-001-A7.png
checkout LIB-GO-001-A7 USR-1001
return LIB-GO-001-A7
inventory
report
exit
```

You may begin by accepting the barcode value directly:
```bash
scan LIB-GO-001-A7
```
Then add image decoding as the final feature.

#### Scan behavior

For an available book:
```bash
Barcode: LIB-GO-001-A7
Title: The Go Programming Language
Status: AVAILABLE
Shelf age: 179 days
```
For a checked-out book:
```bash
Barcode: LIB-GO-001-A7
Title: The Go Programming Language
Status: CHECKED OUT
Borrower: Avery Stone
Checked out: 2026-07-01
Due date: 2026-07-15
Days checked out: 12
Days remaining: 2
Overdue: No
```
For an unknown barcode:
```bash
Barcode not found: LIB-UNKNOWN-999
```

### Application 2: Circulation server

The circulation service owns:

- [ ] user records,
- [ ] checkout records,
- [ ] due dates,
- [ ] borrower lookup,
- [ ] overdue calculations.

Use either `net/http` or `net/rpc`.

Because your previous challenge already used net/rpc, I recommend using JSON HTTP APIs this time.

Example endpoints:
```bash
GET  /checkouts/{barcodeID}
POST /checkouts
POST /returns
GET  /users/{userID}
GET  /health
```
Example checkout request:
```json
{
  "barcode_id": "LIB-GO-001-A7",
  "user_id": "USR-1001",
  "checkout_days": 14
}
```
Example response:

```json
{
  "barcode_id": "LIB-GO-001-A7",
  "user_id": "USR-1001",
  "checked_out_at": "2026-07-13T15:30:00Z",
  "due_at": "2026-07-27T15:30:00Z"
}
```

### Core data structures
```go
type Book struct {
	BarcodeID   string     `json:"barcode_id"`
	Title       string     `json:"title"`
	Author      string     `json:"author"`
	Category    string     `json:"category"`
	AcquiredAt  time.Time  `json:"acquired_date"`
	Status      BookStatus `json:"status"`
	CheckedOutBy *string   `json:"checked_out_by"`
	CheckedOutAt *time.Time `json:"checked_out_at"`
	DueAt        *time.Time `json:"due_date"`
}
```

You may find it cleaner to keep checkout information separate from Book:

```go
type Checkout struct {
	BarcodeID   string     `json:"barcode_id"`
	UserID      string     `json:"user_id"`
	CheckedOutAt time.Time  `json:"checked_out_at"`
	DueAt        time.Time  `json:"due_at"`
	ReturnedAt   *time.Time `json:"returned_at,omitempty"`
}
```

That separation better resembles a real inventory system.

### The math portion

The report should compute the following values.

#### Basic counts

```bash
Total books
Books available
Books checked out
Books overdue
```

Required invariant:
```bash
total books = available books + checked-out books
```
Your application should return an error if that invariant is ever violated.

#### Inventory percentages

Calculate:

```bash
availability percentage
checkout percentage
overdue percentage
```

Formula:

```bash
availability percentage = available / total × 100
```

Example:

```bash
Total books: 3
Available: 2
Checked out: 1

Availability rate: 66.67%
Checkout rate: 33.33%
```

Be careful to avoid integer division:
```go
percentage := float64(available) / float64(total) * 100
```
#### Shelf age


For an available book:

```bash
shelf age = current date - acquired date
```

Use calendar-day calculations rather than blindly rounding hours when possible.

#### Checkout duration

For an active checkout:
```bash
days checked out = current date - checkout date
```

#### Overdue duration
```bash
if current date > due date:
    overdue days = current date - due date
else:
    overdue days = 0
```

#### Average checkout duration

For all active checkouts:

```bash
average checkout duration =
    total days checked out / number of checked-out books
```

Handle the zero-checkout case without dividing by zero.

#### Longest checkout

Determine:

- [ ] which book has been checked out the longest,
- [ ] who has it,
- [ ] how many days it has been checked out.

This is the interview-style algorithmic component.

A single traversal should be sufficient:
```go
var longest Checkout
var maxDuration time.Duration

for _, checkout := range checkouts {
	duration := now.Sub(checkout.CheckedOutAt)

	if duration > maxDuration {
		maxDuration = duration
		longest = checkout
	}
}
```

### Required behavior

Checkout validation

Reject a checkout when:

- [ ] the barcode does not exist,
- [ ] the user does not exist,
- [ ] the user is inactive,
- [ ] the book is already checked out,
- [ ] the requested checkout duration is invalid.

Example errors:
```bash
book already checked out
unknown user
inactive library account
checkout period must be between 1 and 30 days
```

#### Return validation

Reject a return when:

- [ ] the barcode does not exist,
- [ ] the book is not currently checked out,
- [ ] no active checkout record exists.

#### Persistence

After every successful checkout or return:

- [ ] update in-memory state,
- [ ] save the updated checkout data,
- [ ] write using a temporary file,
- [ ] rename the temporary file over the original.

That gives you an atomic-style file replacement:
```bash
checkouts.json.tmp
    ↓ rename
checkouts.json
```
Avoid directly truncating the live JSON file before a successful write completes.

#### Concurrency requirement

Assume multiple kiosks may contact the circulation server.

Protect shared state with:

```go
sync.RWMutex
```

Example:

```go
type Store struct {
	mu        sync.RWMutex
	checkouts map[string]Checkout
}
```

Use:

`RLock` for lookup and report operations,
`Lock` for checkout and return operations.

The important part is making the availability check and state update atomic:

```go
s.mu.Lock()
defer s.mu.Unlock()

if _, exists := s.checkouts[barcodeID]; exists {
	return ErrAlreadyCheckedOut
}

s.checkouts[barcodeID] = checkout
```

Do not perform the existence check under one lock and the insertion under another.

### Context requirement

Every HTTP request should accept and propagate context.Context.

For example:
```go
func (c *Client) GetCheckout(
	ctx context.Context,
	barcodeID string,
) (Checkout, error)
```

Build requests with:
```go
req, err := http.NewRequestWithContext(
	ctx,
	http.MethodGet,
	url,
	nil,
)
```

Give kiosk lookups a timeout:
```go
ctx, cancel := context.WithTimeout(
	context.Background(),
	2*time.Second,
)
defer cancel()
```

The kiosk should distinguish between:

```bash
book not checked out
circulation server unavailable
request timed out
invalid server response
```

#### Final-state report

When the user enters report or exits the kiosk, display something similar to:

```bash
LIBRARY INVENTORY REPORT
Generated: 2026-07-13 15:45:00

Inventory
---------
Total books:        3
Available:          1
Checked out:        2
Overdue:            1

Rates
-----
Availability rate:  33.33%
Checkout rate:      66.67%
Overdue rate:       50.00% of checked-out books

Checkout Statistics
-------------------
Average active checkout: 9.50 days
Longest active checkout: 12 days

Longest Checkout
----------------
Book:       The Go Programming Language
Barcode:    LIB-GO-001-A7
Borrower:   Avery Stone
Checked out: 2026-07-01
Due:         2026-07-10
Overdue by:  3 days

Available Inventory
-------------------
LIB-MATH-88-Z3 | Concrete Mathematics | Shelf age: 125 days
Testing requirements
```

At minimum, write tests for:

- [ ] checking out an available book,
- [ ] rejecting a duplicate checkout,
- [ ] returning a checked-out book,
- [ ] rejecting a return for an available book,
- [ ] calculating overdue days,
- [ ] calculating inventory percentages,
- [ ] finding the longest checkout,
- [ ] handling an empty checkout list,
- [ ] rejecting an unknown user,
- [ ] handling circulation-server timeout.

Useful table-driven test structure:

```go
func TestCalculateOverdueDays(t *testing.T) {
	tests := []struct {
		name string
		now  time.Time
		due  time.Time
		want int
	}{
		{
			name: "not overdue",
			now:  time.Date(2026, 7, 10, 0, 0, 0, 0, time.UTC),
			due:  time.Date(2026, 7, 12, 0, 0, 0, 0, time.UTC),
			want: 0,
		},
		{
			name: "three days overdue",
			now:  time.Date(2026, 7, 13, 0, 0, 0, 0, time.UTC),
			due:  time.Date(2026, 7, 10, 0, 0, 0, 0, time.UTC),
			want: 3,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := CalculateOverdueDays(tt.now, tt.due)

			if got != tt.want {
				t.Errorf("got %d, want %d", got, tt.want)
			}
		})
	}
}
```

Inject the current time instead of calling time.Now() throughout your business logic:

```go
type Clock interface {
	Now() time.Time
}
```

Or pass now `time.Time` into calculation functions. This makes the math deterministic and easy to test.

### Barcode implementation

For the first iteration, your scanner may extract the ID from the filename:
```bash
barcodes/LIB-GO-001-A7.png
```
Then implement real Code 128 decoding as an enhancement.

That staging keeps the challenge focused on Go architecture before adding a third-party imaging dependency.

A clean interface would be:
```go
type Decoder interface {
	Decode(ctx context.Context, imagePath string) (string, error)
}
```
You can initially implement:
```go
type FilenameDecoder struct{}
```
Then later add:
```go
type ImageDecoder struct{}
```

This also gives you an easy fake decoder for tests.

### Stretch goals

After the core implementation works:

- [ ] Add an audit log for every checkout and return.
- [ ] Add maximum checkout limits per user.
- [ ] Prevent users with overdue books from checking out another book.
- [ ] Add reservation support.
- [ ] Add graceful server shutdown.
- [ ] Add request IDs and structured logging.
- [ ] Add atomic counters for scans, checkouts, returns, and failures.
- [ ] Add a /metrics endpoint.
- [ ] Run simultaneous checkout attempts against the same barcode and prove only one succeeds.
- [ ] Add an HTTP test server using httptest.Server.
- [ ] Add a checksum or version number to persisted inventory state.

### Completion commands

Your final validation should include:

```bash
gofmt -w .
go vet ./...
go test ./...
go run ./cmd/circulation-server
go run ./cmd/kiosk
```

### Difficulty

I would rate the core challenge as upper-intermediate Go, with senior-level opportunities in:

* service boundaries,
* concurrency correctness,
* deterministic time calculations,
* persistence safety,
* interface design,
* HTTP error handling,
* testability,
* invariant enforcement.

The central algorithmic problem is fully solvable: calculate inventory statistics and identify the longest active checkout in one traversal. The larger challenge tests whether you can turn that calculation into a reliable distributed application.