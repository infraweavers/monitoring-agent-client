package main

import (
	"bytes"
	"flag"
	"monitoring-agent-client/internal/httpclient"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestArgumentParsing(t *testing.T) {
	t.Run("Basic test returns 200 and renders correct exit code", func(t *testing.T) {
		// We manipuate the Args to set them up for the testcases
		// After this test we restore the initial args
		oldArgs := os.Args
		defer func() { os.Args = oldArgs }()

		flag.CommandLine = flag.NewFlagSet("flag", flag.ExitOnError)

		os.Args = []string{"main.exe", "-host", "remotehost", "-username", "thisismyusername", "-password", "thisismypassword", "-executable", "/path/to/executable", "-script", "README.md"}
		httpClient := httpclient.NewMockHTTPClient(`{"output": "Test output", "exitcode": 2}`, 200)

		var buf bytes.Buffer
		actualExit := invokeClient(&buf, httpClient)

		actualOutput := buf.String()

		assert.Equal(t, `{"args":null,"path":"/path/to/executable","scriptarguments":[],"stdin":"# monitoring-agent-client","timeout":"10s"}`, httpClient.RequestBodyContent)
		assert.Equal(t, 10*time.Second, httpClient.Timeout)
		assert.Equal(t, false, httpClient.Transport.TLSClientConfig.InsecureSkipVerify)
		assert.Equal(t, "Basic dGhpc2lzbXl1c2VybmFtZTp0aGlzaXNteXBhc3N3b3Jk", httpClient.RequestHeaders["Authorization"][0])
		assert.Equal(t, "remotehost:9000", httpClient.RequestHost)
		assert.Equal(t, "/v1/runscriptstdin", httpClient.RequestURI.Path)
		assert.Equal(t, "POST", httpClient.RequestVerb)

		assert.Equal(t, 2, actualExit)
		assert.Equal(t, "Test output", actualOutput)
	})
	t.Run("Insecure basic test returns 200 and renders correct exit code", func(t *testing.T) {
		oldArgs := os.Args
		defer func() { os.Args = oldArgs }()
		flag.CommandLine = flag.NewFlagSet("flag", flag.ExitOnError)
		os.Args = []string{"main.exe", "-host", "remotehost", "-username", "thisismyusername", "-password", "thisismypassword", "-executable", "/path/to/executable", "-script", "README.md", "-insecure"}
		httpClient := httpclient.NewMockHTTPClient(`{"output": "Test output", "exitcode": 2}`, 200)
		var buf bytes.Buffer
		actualExit := invokeClient(&buf, httpClient)
		actualOutput := buf.String()
		assert.Equal(t, `{"args":null,"path":"/path/to/executable","scriptarguments":[],"stdin":"# monitoring-agent-client","timeout":"10s"}`, httpClient.RequestBodyContent)
		assert.Equal(t, 10*time.Second, httpClient.Timeout)
		assert.Equal(t, true, httpClient.Transport.TLSClientConfig.InsecureSkipVerify)
		assert.Equal(t, "Basic dGhpc2lzbXl1c2VybmFtZTp0aGlzaXNteXBhc3N3b3Jk", httpClient.RequestHeaders["Authorization"][0])
		assert.Equal(t, "remotehost:9000", httpClient.RequestHost)
		assert.Equal(t, "/v1/runscriptstdin", httpClient.RequestURI.Path)
		assert.Equal(t, "POST", httpClient.RequestVerb)
		assert.Equal(t, 2, actualExit)
		assert.Equal(t, "Test output", actualOutput)
	})
}
