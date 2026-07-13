package main

import (
	"flag"
	"fmt"
	"net"
	"net/rpc"
	"os"
	"time"

	"localpractice5b.com/challenges/types"
)

/**
go run . deploy --name sensor-api --version 1.4.2 --environment staging
*/

/**
call DeploymentReceiver.Deploy

main()
  └── run()
       ├── recognizes "deploy" [DONE]
       ├── creates a deploy FlagSet [DONE]
       ├── validates required flags [DONE]
       ├── constructs DeploymentRequest [DONE]
       ├── calls the injected RPC dialer
       ├── performs the RPC call
       ├── prints the response
       └── returns the appropriate exit code
*/

// func run(
// 	args []string,
// 	stdout io.Writer,
// 	stderr io.Writer,
// 	dialer types.RPCDialer,
// ) int

func realRPCDialer(
	address string,
	timeout time.Duration,
) (types.RPCClient, error) {
	connection, err := net.DialTimeout("tcp", address, timeout)
	if err != nil {
		return nil, err
	}

	return rpc.NewClient(connection), nil
}

func Deploy(args []string, dialer types.RPCDialer) (string, error) {
	fs := flag.NewFlagSet("deploy", flag.ContinueOnError)
	name := fs.String("name", "unknown", "sensor name or ID")
	version := fs.String("version", "0.0.0.0", "receiver version to use")
	env := fs.String("env", "none", "environment to work in")
	dryRun := fs.Bool("dryRun", false, "should not process anything live")
	if err := fs.Parse(args); err != nil {
		return "Failed to parse arguments", err
	}
	req := types.DeploymentRequest{
		Name:        *name,
		Version:     *version,
		Environment: *env,
		DryRun:      *dryRun,
	}
	return HandleDial(req, dialer)
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("error: command is required")
		fmt.Println("usage: missionctl deploy [flags]")
		return
	}
	/**
	[0] executable
	[1] deploy
	[2:] deploy flags
	*/
	cmd := os.Args[1]
	cArgs := os.Args[2:]
	switch cmd {
	case "deploy":
		msg, err := Deploy(cArgs, realRPCDialer)
		if err != nil {
			fmt.Println(err.Error())
			os.Exit(1)
		} else {
			fmt.Println(msg)
		}
	default:
		fmt.Printf("error: unknown command - %q\n", cmd)
	}
}
