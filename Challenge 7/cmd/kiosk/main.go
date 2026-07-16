package main

import (
	"context"
	"fmt"
	"os"
)

var appConfig AppConfig

// get commands from CLI call
// note: if antivirus complains, then use: go build -o .\bin\kiosk.exe .\cmd\kiosk
func main() {
	ctx := context.Background()

	msg, err := appConfig.LoadConfigurations()
	if err != nil {
		fmt.Println(msg)
		os.Exit(1)
	}
	/**
	[0] executable
	[1] scan
	[2:] scan flags
	*/
	cmd := os.Args[1]
	cArgs := os.Args[2:]

	// handle command keywords and args
	handleCommands(ctx, cmd, cArgs)
}
