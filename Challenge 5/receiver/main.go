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
