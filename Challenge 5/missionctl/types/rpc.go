package types

import "time"

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

type RPCClient interface {
	Call(
		serviceMethod string,
		args any,
		reply any,
	) error

	Close() error
}

type RPCDialer func(
	address string,
	timeout time.Duration,
) (RPCClient, error)
