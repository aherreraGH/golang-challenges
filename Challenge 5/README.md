## Challenge Objective

Create a CLI named missionctl that submits a deployment request to the RPC receiver.

Example:

```bash
missionctl deploy `
  --name sensor-api `
  --version 1.4.2 `
  --environment staging `
  --server localhost:9000
```

The receiver will validate the request, print the result to its console, and return a structured response containing:

- [ ] Whether the operation succeeded
- [ ] A human-readable message
- [ ] A generated request ID

### Suggested Project Structure

```bash
golang-rpc-cli-challenge/
├── receiver/
│   ├── go.mod
│   └── main.go
└── missionctl/
    ├── go.mod
    ├── main.go
    ├── cli.go
    └── cli_test.go
```

You may organize your CLI differently, but the receiver should remain unchanged.

### Receiver Application

#### receiver/go.mod

```go
module rpc-receiver

go 1.24
```

Use your installed Go version if necessary.

#### receiver/main.go

```go
package main

import (
	"errors"
	"fmt"
	"log"
	"net"
	"net/rpc"
	"strings"
	"sync/atomic"
	"time"
)

const (
	listenAddress = "127.0.0.1:9000"
)

var requestCounter uint64

// DeploymentRequest is sent by the CLI application.
//
// Every field must be exported because net/rpc uses encoding/gob
// to serialize values between the client and server.
type DeploymentRequest struct {
	Name        string
	Version     string
	Environment string
	DryRun      bool
}

// DeploymentResponse is returned to the CLI application.
type DeploymentResponse struct {
	Success   bool
	Message   string
	RequestID string
}

// DeploymentReceiver exposes methods that may be called through net/rpc.
type DeploymentReceiver struct{}

// Deploy validates and processes a deployment request.
//
// Methods exposed through net/rpc must:
//   - Be exported
//   - Have two exported or built-in argument types
//   - Return exactly one error value
func (r *DeploymentReceiver) Deploy(
	request DeploymentRequest,
	response *DeploymentResponse,
) error {
	requestID := createRequestID()

	response.RequestID = requestID

	fmt.Printf(
		"Received deployment request: id=%s name=%q version=%q environment=%q dryRun=%t\n",
		requestID,
		request.Name,
		request.Version,
		request.Environment,
		request.DryRun,
	)

	if err := validateRequest(request); err != nil {
		response.Success = false
		response.Message = err.Error()

		fmt.Printf(
			"Deployment rejected: id=%s reason=%q\n",
			requestID,
			err,
		)

		// This is intentionally returned as a successful RPC call.
		// The CLI must inspect response.Success.
		return nil
	}

	// Simulate a small amount of remote processing.
	time.Sleep(150 * time.Millisecond)

	response.Success = true

	if request.DryRun {
		response.Message = fmt.Sprintf(
			"Dry-run validation succeeded for %s version %s in %s",
			request.Name,
			request.Version,
			request.Environment,
		)
	} else {
		response.Message = fmt.Sprintf(
			"Deployment accepted for %s version %s in %s",
			request.Name,
			request.Version,
			request.Environment,
		)
	}

	fmt.Printf(
		"Deployment processed: id=%s success=%t message=%q\n",
		requestID,
		response.Success,
		response.Message,
	)

	return nil
}

func validateRequest(request DeploymentRequest) error {
	if strings.TrimSpace(request.Name) == "" {
		return errors.New("application name is required")
	}

	if strings.TrimSpace(request.Version) == "" {
		return errors.New("application version is required")
	}

	switch strings.ToLower(strings.TrimSpace(request.Environment)) {
	case "development", "staging", "production":
		return nil
	default:
		return fmt.Errorf(
			"unsupported environment %q: expected development, staging, or production",
			request.Environment,
		)
	}
}

func createRequestID() string {
	number := atomic.AddUint64(&requestCounter, 1)

	return fmt.Sprintf("REQ-%06d", number)
}

func main() {
	receiver := &DeploymentReceiver{}

	if err := rpc.Register(receiver); err != nil {
		log.Fatalf("failed to register RPC receiver: %v", err)
	}

	listener, err := net.Listen("tcp", listenAddress)
	if err != nil {
		log.Fatalf("failed to listen on %s: %v", listenAddress, err)
	}
	defer listener.Close()

	log.Printf("RPC receiver listening on %s", listenAddress)

	for {
		connection, err := listener.Accept()
		if err != nil {
			log.Printf("failed to accept connection: %v", err)
			continue
		}

		go rpc.ServeConn(connection)
	}
}
```

### Running the Receiver

From the receiver directory:

```bash
go run .
```

Expected startup output:

```bash
RPC receiver listening on 127.0.0.1:9000
```

Leave that terminal running while testing your CLI manually.

### Your CLI Requirements

Your CLI must support this command:

```bash
missionctl deploy
```

Required flags:
```bash
--name
--version
--environment
```

Optional flags:
```bash
--server
--dry-run
--timeout
```
Recommended defaults:
```bash
--server      127.0.0.1:9000
--timeout     3s
--dry-run     false
```
Example:
```bash
go run . deploy `
  --name sensor-api `
  --version 1.4.2 `
  --environment staging
```

### Required CLI Behavior

#### Successful request

Command:
```bash
go run . deploy `
  --name sensor-api `
  --version 1.4.2 `
  --environment staging
```

Suggested output:
```bash
Request ID: REQ-000001
Success: true
Message: Deployment accepted for sensor-api version 1.4.2 in staging
```
Exit code:
```bash
0
```

#### Receiver rejects the request

Command:

```bash
go run . deploy `
  --name sensor-api `
  --version 1.4.2 `
  --environment test
```

The RPC call itself succeeds, but the receiver returns:

Success: false

Suggested CLI output:

```bash
Request ID: REQ-000002
Success: false
Message: unsupported environment "test": expected development, staging, or production
```

Exit code:
```bash
1
```

#### Connection failure

Command:
```bash
go run . deploy `
  --name sensor-api `
  --version 1.4.2 `
  --environment staging `
  --server 127.0.0.1:9999
```

Suggested error:

```bash
error: failed to connect to RPC server 127.0.0.1:9999: ...
```

Exit code:
```bash
1
```

### Technical Requirements

Use only the Go standard library for the initial implementation.

Your CLI should use:
```go
net/rpc
flag
context
net
time
io
```

You do not necessarily need all of those packages, but your timeout handling should prevent the program from hanging indefinitely.

Do not place all logic directly inside main().

A recommended structure is:
```go
func main() {
	os.Exit(run(os.Args[1:], os.Stdout, os.Stderr))
}
```
Your testable CLI logic could resemble:

```go
func run(
	args []string,
	stdout io.Writer,
	stderr io.Writer,
) int
```

This lets your tests:

- [ ] Pass command-line arguments directly
- [ ] Capture standard output
- [ ] Capture error output
- [ ] Inspect the returned exit code
- [ ] Avoid calling os.Exit() during tests

### RPC Types

Because the CLI and receiver are separate Go modules, define matching request and response structs in your CLI.

They must match the receiver’s field names and field types:
```go
type DeploymentRequest struct {
	Name        string
	Version     string
	Environment string
	DryRun      bool
}

type DeploymentResponse struct {
	Success   bool
	Message   string
	RequestID string
}
```
Call this RPC method:
```go
DeploymentReceiver.Deploy
```
Example RPC call:
```go
var response DeploymentResponse

err := client.Call(
	"DeploymentReceiver.Deploy",
	request,
	&response,
)
```

#### Timeout Requirement

The CLI must not wait indefinitely when connecting to a server.

One option is to establish the TCP connection yourself:
```go
connection, err := net.DialTimeout(
	"tcp",
	serverAddress,
	timeout,
)
```
Then create the RPC client from that connection:
```go
client := rpc.NewClient(connection)
defer client.Close()
```
For additional credit, also enforce a timeout around the RPC call itself rather than only the initial connection.

### Testing Requirements

You must implement at least two automated tests, and both must pass with:
```go
go test ./...
```
Do not require the real receiver application to be running during unit tests.

#### Test 1: Successful deployment

Verify that a valid request:
```bash
name: sensor-api
version: 1.4.2
environment: staging
```
Produces:

- [ ] Exit code 0
- [ ] Output containing Success: true
- [ ] Output containing the receiver’s success message
- [ ] No unexpected error output

#### Test 2: Expected failure

Choose one of these failure paths.

##### Option A: Receiver rejection

Simulate an RPC response containing:
```go
DeploymentResponse{
	Success:   false,
	Message:   `unsupported environment "test"`,
	RequestID: "REQ-000002",
}
```

Verify:

- [ ] Exit code 1
- [ ] Output or error output contains the rejection message
- [ ] The test passes because the CLI handled the failure correctly
- [ ] Option B: Connection or RPC error

Simulate an error such as:
```go
errors.New("connection refused")
```
Verify:

- [ ] Exit code 1
- [ ] Error output contains a useful error message
- [ ] The test passes because the error was expected and handled

The test should not fail merely because the scenario represents an application failure.

### Recommended Testability Design

Introduce a small interface instead of coupling all CLI logic directly to *rpc.Client.

For example:
```go
type RPCClient interface {
	Call(
		serviceMethod string,
		args any,
		reply any,
	) error

	Close() error
}
```
You can then create:
```go
type RPCDialer func(
	address string,
	timeout time.Duration,
) (RPCClient, error)
```
Production code uses a real RPC dialer. Tests provide a fake implementation.

A fake client might look like:
```go
type fakeRPCClient struct {
	callFunc func(
		serviceMethod string,
		args any,
		reply any,
	) error
}

func (f *fakeRPCClient) Call(
	serviceMethod string,
	args any,
	reply any,
) error {
	return f.callFunc(serviceMethod, args, reply)
}

func (f *fakeRPCClient) Close() error {
	return nil
}
```

This is only an architectural hint. You still need to decide where the interface belongs and how the dependencies are injected.

### Acceptance Criteria

Your submission is complete when:

- [ ] missionctl deploy works against the supplied receiver.
- [ ] Required flags are validated.
- [ ] The server address is configurable.
- [ ] A connection timeout is implemented.
- [ ] Receiver responses with Success: false produce a nonzero exit code.
- [ ] RPC and connection errors produce useful messages.
- [ ] Resources are properly closed.
- [ ] main() remains small.
- [ ] At least two tests are present.
- [ ] Both tests pass with go test ./....
- [ ] Code is formatted with gofmt.
- [ ] go vet ./... reports no problems.

### Optional Enhancements

After meeting the core requirements, you could add:

- [ ] --output text
- [ ] --output json
- [ ] Environment normalization
- [ ] Semantic version validation
- [ ] Retry support
- [ ] Configurable retry delay
- [ ] Signal cancellation
- [ ] Additional subcommands
- [ ] Table-driven flag-validation tests

A JSON output mode could produce:

```json
{
  "requestId": "REQ-000001",
  "success": true,
  "message": "Deployment accepted for sensor-api version 1.4.2 in staging"
}
```

The strongest version of this challenge will keep parsing, RPC transport, result handling, and process exit behavior separated rather than placing them all in main.go.