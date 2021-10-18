package main

import (
	"bytes"
	"flag"
	"monitoring-agent-client/internal/httpclient"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestArgumentParsing(t *testing.T) {

	// We manipuate the Args to set them up for the testcases
	// After this test we restore the initial args
	oldArgs := os.Args
	defer func() { os.Args = oldArgs }()

	// this call is required because otherwise flags panics,
	// if args are set between flag.Parse call
	flag.CommandLine = flag.NewFlagSet("flag", flag.ExitOnError)

	os.Args = []string{"main.exe", "-host", "localhost", "-password", "passwords", "-executable", "/path/to/executable", "-script", "README.md"}
	var buf bytes.Buffer

	// code under test:
	httpClient := httpclient.NewMockHTTPClient()
	actualExit := invokeClient(&buf, httpClient)

	actualOutput := buf.String()

	assert.Equal(t, 3, actualExit)
	assert.Equal(t, "host is not set", actualOutput)
}
