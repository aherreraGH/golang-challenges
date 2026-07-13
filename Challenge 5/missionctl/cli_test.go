package main

import (
	"strings"
	"testing"
	"time"

	"localpractice5b.com/challenges/types"
)

func fakeRPCDialer(
	address string,
	timeout time.Duration,
) (types.RPCClient, error) {
	return &FakeClient{}, nil
}

// TestHelloEmpty calls greetings.Hello with an empty string,
// checking for an error.
func TestSuccessfulFlagParsing(t *testing.T) {
	var args []string = []string{"-name", "joe", "-version", "2.0", "-env", "staging"}
	msg, err := Deploy(args, fakeRPCDialer)
	if err != nil {
		t.Errorf(`%s: %v`, msg, err)
	}
}

func TestExpectFailedFlagParsing(t *testing.T) {
	var args []string = []string{"-name", "joe", "-version", "2.0", "-envx", "staging"}
	msg, err := Deploy(args, fakeRPCDialer)
	if err == nil {
		t.Errorf(`%s: %v`, msg, err)
		return
	}
	if !strings.Contains(msg, "Failed to parse arguments") {
		t.Errorf(`Expected to fail: %s`, msg)
	}
}
