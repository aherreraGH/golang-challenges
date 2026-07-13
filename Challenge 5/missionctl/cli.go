package main

import (
	"errors"
	"fmt"
	"time"

	"localpractice5b.com/challenges/types"
)

/**
Request ID: REQ-000001
Success: true
Message: Deployment accepted for sensor-api version 1.4.2 in staging
*/

func HandleDial(req types.DeploymentRequest, dialer types.RPCDialer) (string, error) {
	// msg := fmt.Sprintf("Deploying using name: %v, version: %v, environment: %v", *name, *version, *env)
	client, err := dialer("127.0.0.1:9000", 3*time.Second)
	if err != nil {
		return "Failed to connect", err
	}
	defer client.Close()

	reply := types.DeploymentResponse{}

	client.Call("DeploymentReceiver.Deploy", req, &reply)

	if reply.Success {
		return reply.Message, nil
	} else {
		return "Failed somewere", errors.New(fmt.Sprintf("Error: %v", reply.Message))
	}
}
