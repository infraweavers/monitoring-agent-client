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

		os.Args = []string{
			"main.exe",
			"-host", "remotehost",
			"-username", "thisismyusername",
			"-password", "thisismypassword",
			"-executable", "/path/to/executable",
			"-script", "TestScript-Valid.ps1",
		}
		httpClient := httpclient.NewMockHTTPClient(`{"output": "Test output", "exitcode": 2}`, 200)

		var buf bytes.Buffer
		actualExit := invokeClient(&buf, httpClient)
		actualOutput := buf.String()

		assert.Equal(t, `{"args":null,"path":"/path/to/executable","scriptarguments":[],"stdin":"Write-Host \"This is a test script\"\r\n\r\n","timeout":"10s"}`, httpClient.RequestBodyContent)
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

		os.Args = []string{
			"main.exe",
			"-host", "remotehost",
			"-username", "thisismyusername",
			"-password", "thisismypassword",
			"-executable", "/path/to/executable",
			"-script", "TestScript-Valid.ps1",
			"-insecure",
		}

		httpClient := httpclient.NewMockHTTPClient(`{"output": "Test output", "exitcode": 2}`, 200)

		var buf bytes.Buffer
		actualExit := invokeClient(&buf, httpClient)
		actualOutput := buf.String()

		assert.Equal(t, `{"args":null,"path":"/path/to/executable","scriptarguments":[],"stdin":"Write-Host \"This is a test script\"\r\n\r\n","timeout":"10s"}`, httpClient.RequestBodyContent)
		assert.Equal(t, 10*time.Second, httpClient.Timeout)
		assert.Equal(t, true, httpClient.Transport.TLSClientConfig.InsecureSkipVerify)
		assert.Equal(t, "Basic dGhpc2lzbXl1c2VybmFtZTp0aGlzaXNteXBhc3N3b3Jk", httpClient.RequestHeaders["Authorization"][0])
		assert.Equal(t, "remotehost:9000", httpClient.RequestHost)
		assert.Equal(t, "/v1/runscriptstdin", httpClient.RequestURI.Path)
		assert.Equal(t, "POST", httpClient.RequestVerb)

		assert.Equal(t, 2, actualExit)
		assert.Equal(t, "Test output", actualOutput)
	})

	t.Run("Basic test returns 200 and renders correct exit code with a perl script", func(t *testing.T) {
		// We manipuate the Args to set them up for the testcases
		// After this test we restore the initial args
		oldArgs := os.Args
		defer func() { os.Args = oldArgs }()

		flag.CommandLine = flag.NewFlagSet("flag", flag.ExitOnError)

		os.Args = []string{
			"main.exe",
			"-host", "remotehost",
			"-username", "thisismyusername",
			"-password", "thisismypassword",
			"-executable", "/path/to/executable",
			"-script", "TestScript.pl",
		}
		httpClient := httpclient.NewMockHTTPClient(`{"output": "Test output", "exitcode": 2}`, 200)

		var buf bytes.Buffer
		actualExit := invokeClient(&buf, httpClient)
		actualOutput := buf.String()

		assert.Equal(t, `{"args":null,"path":"/path/to/executable","scriptarguments":[],"stdin":"#!/perl\n\nprint \"this is a test script\"\n","stdinsignature":"untrusted comment: signature from minisign secret key\r\nRWTV8L06+shYI3jk77ofKAmdXcat5J7EVM/6JLX3ssHhRFqqIAU1vc49KF9Hn3+kO/+k6bFBND+W40LZM8ae4TtQY2NF6HaBpAI=\r\ntrusted comment: timestamp:1634631414\tfile:TestScript.pl\r\nixE4k+I3rIX1S3aTt/q4rTx9aZUygKYITgPQFkbnq+WPq4TwtW4Q9LmDMr5caG5FlPxWT6ve8rvBjZXxkogHBw==\r\n","timeout":"10s"}`, httpClient.RequestBodyContent)
		assert.Equal(t, 10*time.Second, httpClient.Timeout)
		assert.Equal(t, false, httpClient.Transport.TLSClientConfig.InsecureSkipVerify)
		assert.Equal(t, "Basic dGhpc2lzbXl1c2VybmFtZTp0aGlzaXNteXBhc3N3b3Jk", httpClient.RequestHeaders["Authorization"][0])
		assert.Equal(t, "remotehost:9000", httpClient.RequestHost)
		assert.Equal(t, "/v1/runscriptstdin", httpClient.RequestURI.Path)
		assert.Equal(t, "POST", httpClient.RequestVerb)

		assert.Equal(t, 2, actualExit)
		assert.Equal(t, "Test output", actualOutput)
	})

	t.Run("Test with executable arguments returns 200 and returns executable arguments passed", func(t *testing.T) {
		oldArgs := os.Args
		defer func() { os.Args = oldArgs }()
		flag.CommandLine = flag.NewFlagSet("flag", flag.ExitOnError)

		os.Args = []string{
			"main.exe",
			"-host", "remotehost",
			"-username", "thisismyusername",
			"-password", "thisismypassword",
			"-executable", "/path/to/executable",
			"-script", "TestScript-Valid.ps1",
			"-executableArg", "arg1",
			"-executableArg", "arg2",
		}

		httpClient := httpclient.NewMockHTTPClient(`{"output": "Test output", "exitcode": 2}`, 200)

		var buf bytes.Buffer
		actualExit := invokeClient(&buf, httpClient)
		actualOutput := buf.String()

		assert.Equal(t, `{"args":["arg1","arg2"],"path":"/path/to/executable","scriptarguments":[],"stdin":"Write-Host \"This is a test script\"\r\n\r\n","timeout":"10s"}`, httpClient.RequestBodyContent)
		assert.Equal(t, 10*time.Second, httpClient.Timeout)
		assert.Equal(t, false, httpClient.Transport.TLSClientConfig.InsecureSkipVerify)
		assert.Equal(t, "Basic dGhpc2lzbXl1c2VybmFtZTp0aGlzaXNteXBhc3N3b3Jk", httpClient.RequestHeaders["Authorization"][0])
		assert.Equal(t, "remotehost:9000", httpClient.RequestHost)
		assert.Equal(t, "/v1/runscriptstdin", httpClient.RequestURI.Path)
		assert.Equal(t, "POST", httpClient.RequestVerb)

		assert.Equal(t, 2, actualExit)
		assert.Equal(t, "Test output", actualOutput)
	})

	t.Run("Test with script arguments returns 200 and returns script arguments passed", func(t *testing.T) {
		oldArgs := os.Args
		defer func() { os.Args = oldArgs }()
		flag.CommandLine = flag.NewFlagSet("flag", flag.ExitOnError)

		os.Args = []string{
			"main.exe",
			"-host", "remotehost",
			"-username", "thisismyusername",
			"-password", "thisismypassword",
			"-executable", "/path/to/executable",
			"-script", "TestScript-Valid.ps1",
			"-executableArg", "arg1",
			"-executableArg", "arg2",
			"--", "scriptarg1", "-scriptarg scriptarg2", "-scriptarg", "scriptarg3", "--warning=3",
		}

		httpClient := httpclient.NewMockHTTPClient(`{"output": "Test output", "exitcode": 2}`, 200)

		var buf bytes.Buffer
		actualExit := invokeClient(&buf, httpClient)
		actualOutput := buf.String()

		assert.Equal(t, `{"args":["arg1","arg2"],"path":"/path/to/executable","scriptarguments":["scriptarg1","-scriptarg scriptarg2","-scriptarg","scriptarg3","--warning=3"],"stdin":"Write-Host \"This is a test script\"\r\n\r\n","timeout":"10s"}`, httpClient.RequestBodyContent)
		assert.Equal(t, 10*time.Second, httpClient.Timeout)
		assert.Equal(t, false, httpClient.Transport.TLSClientConfig.InsecureSkipVerify)
		assert.Equal(t, "Basic dGhpc2lzbXl1c2VybmFtZTp0aGlzaXNteXBhc3N3b3Jk", httpClient.RequestHeaders["Authorization"][0])
		assert.Equal(t, "remotehost:9000", httpClient.RequestHost)
		assert.Equal(t, "/v1/runscriptstdin", httpClient.RequestURI.Path)
		assert.Equal(t, "POST", httpClient.RequestVerb)

		assert.Equal(t, 2, actualExit)
		assert.Equal(t, "Test output", actualOutput)
	})

	t.Run("Test the client certificate and key are loaded corrected", func(t *testing.T) {
		oldArgs := os.Args
		defer func() { os.Args = oldArgs }()
		flag.CommandLine = flag.NewFlagSet("flag", flag.ExitOnError)
		os.Args = []string{"main.exe",
			"-host", "remotehost",
			"-username", "thisismyusername",
			"-password", "thisismypassword",
			"-executable", "/path/to/executable",
			"-certificate", "server.crt",
			"-key", "server.key",
			"-script", "TestScript-Valid.ps1",
		}

		httpClient := httpclient.NewMockHTTPClient(`{"output": "Test output", "exitcode": 1}`, 200)

		var buf bytes.Buffer
		actualExit := invokeClient(&buf, httpClient)
		actualOutput := buf.String()

		assert.Equal(t, `{"args":null,"path":"/path/to/executable","scriptarguments":[],"stdin":"Write-Host \"This is a test script\"\r\n\r\n","timeout":"10s"}`, httpClient.RequestBodyContent)
		assert.Equal(t, 10*time.Second, httpClient.Timeout)
		assert.Equal(t, false, httpClient.Transport.TLSClientConfig.InsecureSkipVerify)
		assert.Equal(t, 1, len(httpClient.Transport.TLSClientConfig.Certificates))
		assert.NotNil(t, httpClient.Transport.TLSClientConfig.Certificates[0].Certificate)
		assert.NotNil(t, httpClient.Transport.TLSClientConfig.Certificates[0].PrivateKey)
		assert.Equal(t, "Basic dGhpc2lzbXl1c2VybmFtZTp0aGlzaXNteXBhc3N3b3Jk", httpClient.RequestHeaders["Authorization"][0])
		assert.Equal(t, "remotehost:9000", httpClient.RequestHost)
		assert.Equal(t, "/v1/runscriptstdin", httpClient.RequestURI.Path)
		assert.Equal(t, "POST", httpClient.RequestVerb)

		assert.Equal(t, 1, actualExit)
		assert.Equal(t, "Test output", actualOutput)
	})

	t.Run("Test the CACertificate is loaded correctly", func(t *testing.T) {
		oldArgs := os.Args
		defer func() { os.Args = oldArgs }()
		flag.CommandLine = flag.NewFlagSet("flag", flag.ExitOnError)

		os.Args = []string{"main.exe",
			"-host", "remotehost",
			"-username", "thisismyusername",
			"-password", "thisismypassword",
			"-executable", "/path/to/executable",
			"-cacert", "cacert.pem",
			"-script", "TestScript-Valid.ps1",
		}

		httpClient := httpclient.NewMockHTTPClient(`{"output": "Test output", "exitcode": 1}`, 200)

		var buf bytes.Buffer
		actualExit := invokeClient(&buf, httpClient)
		actualOutput := buf.String()

		assert.Equal(t, `{"args":null,"path":"/path/to/executable","scriptarguments":[],"stdin":"Write-Host \"This is a test script\"\r\n\r\n","timeout":"10s"}`, httpClient.RequestBodyContent)
		assert.Equal(t, 10*time.Second, httpClient.Timeout)
		assert.Equal(t, false, httpClient.Transport.TLSClientConfig.InsecureSkipVerify)
		assert.NotNil(t, httpClient.Transport.TLSClientConfig.RootCAs)
		assert.Equal(t, "Basic dGhpc2lzbXl1c2VybmFtZTp0aGlzaXNteXBhc3N3b3Jk", httpClient.RequestHeaders["Authorization"][0])
		assert.Equal(t, "remotehost:9000", httpClient.RequestHost)
		assert.Equal(t, "/v1/runscriptstdin", httpClient.RequestURI.Path)
		assert.Equal(t, "POST", httpClient.RequestVerb)

		assert.Equal(t, 1, actualExit)
		assert.Equal(t, "Test output", actualOutput)
	})

	t.Run("Test the CA, client certificate and key are loaded corrected", func(t *testing.T) {
		oldArgs := os.Args
		defer func() { os.Args = oldArgs }()
		flag.CommandLine = flag.NewFlagSet("flag", flag.ExitOnError)
		os.Args = []string{"main.exe",
			"-host", "remotehost",
			"-username", "thisismyusername",
			"-password", "thisismypassword",
			"-executable", "/path/to/executable",
			"-certificate", "server.crt",
			"-cacert", "cacert.pem",
			"-key", "server.key",
			"-script", "TestScript-Valid.ps1",
		}

		httpClient := httpclient.NewMockHTTPClient(`{"output": "Test output", "exitcode": 1}`, 200)

		var buf bytes.Buffer
		actualExit := invokeClient(&buf, httpClient)
		actualOutput := buf.String()

		assert.Equal(t, `{"args":null,"path":"/path/to/executable","scriptarguments":[],"stdin":"Write-Host \"This is a test script\"\r\n\r\n","timeout":"10s"}`, httpClient.RequestBodyContent)
		assert.Equal(t, 10*time.Second, httpClient.Timeout)
		assert.Equal(t, false, httpClient.Transport.TLSClientConfig.InsecureSkipVerify)
		assert.Equal(t, 1, len(httpClient.Transport.TLSClientConfig.Certificates))
		assert.NotNil(t, httpClient.Transport.TLSClientConfig.Certificates[0].Certificate)
		assert.NotNil(t, httpClient.Transport.TLSClientConfig.Certificates[0].PrivateKey)
		assert.NotNil(t, httpClient.Transport.TLSClientConfig.RootCAs)
		assert.Equal(t, "Basic dGhpc2lzbXl1c2VybmFtZTp0aGlzaXNteXBhc3N3b3Jk", httpClient.RequestHeaders["Authorization"][0])
		assert.Equal(t, "remotehost:9000", httpClient.RequestHost)
		assert.Equal(t, "/v1/runscriptstdin", httpClient.RequestURI.Path)
		assert.Equal(t, "POST", httpClient.RequestVerb)

		assert.Equal(t, 1, actualExit)
		assert.Equal(t, "Test output", actualOutput)
	})

	t.Run("The timeout is passed to the remote server and set on the HTTP Client", func(t *testing.T) {
		// We manipuate the Args to set them up for the testcases
		// After this test we restore the initial args
		oldArgs := os.Args
		defer func() { os.Args = oldArgs }()

		flag.CommandLine = flag.NewFlagSet("flag", flag.ExitOnError)

		os.Args = []string{
			"main.exe",
			"-host", "remotehost",
			"-username", "thisismyusername",
			"-password", "thisismypassword",
			"-executable", "/path/to/executable",
			"-executableArg", "-s",
			"-script", "TestScript-Valid.ps1",
			"-timeout", "1s",
		}

		httpClient := httpclient.NewMockHTTPClient(`{"output": "Test output", "exitcode": 2}`, 200)
		var buf bytes.Buffer
		actualExit := invokeClient(&buf, httpClient)
		actualOutput := buf.String()

		assert.Equal(t, `{"args":["-s"],"path":"/path/to/executable","scriptarguments":[],"stdin":"Write-Host \"This is a test script\"\r\n\r\n","timeout":"1s"}`, httpClient.RequestBodyContent)
		assert.Equal(t, 1*time.Second, httpClient.Timeout)
		assert.Equal(t, false, httpClient.Transport.TLSClientConfig.InsecureSkipVerify)
		assert.Equal(t, "Basic dGhpc2lzbXl1c2VybmFtZTp0aGlzaXNteXBhc3N3b3Jk", httpClient.RequestHeaders["Authorization"][0])
		assert.Equal(t, "remotehost:9000", httpClient.RequestHost)
		assert.Equal(t, "/v1/runscriptstdin", httpClient.RequestURI.Path)
		assert.Equal(t, "POST", httpClient.RequestVerb)

		assert.Equal(t, 2, actualExit)
		assert.Equal(t, "Test output", actualOutput)
	})

	t.Run("A 400 response should be an UNKNOWN exit code", func(t *testing.T) {
		// We manipuate the Args to set them up for the testcases
		// After this test we restore the initial args
		oldArgs := os.Args
		defer func() { os.Args = oldArgs }()

		flag.CommandLine = flag.NewFlagSet("flag", flag.ExitOnError)

		os.Args = []string{
			"main.exe",
			"-host", "remotehost",
			"-username", "thisismyusername",
			"-password", "thisismypassword",
			"-executable", "/path/to/executable",
			"-executableArg", "-s",
			"-script", "TestScript-Valid.ps1",
		}
		httpClient := httpclient.NewMockHTTPClient(`{"output": "Error", "exitcode": 1}`, 400)

		var buf bytes.Buffer
		actualExit := invokeClient(&buf, httpClient)

		actualOutput := buf.String()

		assert.Equal(t, `{"args":["-s"],"path":"/path/to/executable","scriptarguments":[],"stdin":"Write-Host \"This is a test script\"\r\n\r\n","timeout":"10s"}`, httpClient.RequestBodyContent)
		assert.Equal(t, 10*time.Second, httpClient.Timeout)
		assert.Equal(t, false, httpClient.Transport.TLSClientConfig.InsecureSkipVerify)
		assert.Equal(t, "Basic dGhpc2lzbXl1c2VybmFtZTp0aGlzaXNteXBhc3N3b3Jk", httpClient.RequestHeaders["Authorization"][0])
		assert.Equal(t, "remotehost:9000", httpClient.RequestHost)
		assert.Equal(t, "/v1/runscriptstdin", httpClient.RequestURI.Path)
		assert.Equal(t, "POST", httpClient.RequestVerb)

		assert.Equal(t, 3, actualExit)
		assert.Equal(t, "Response code: 400\n{\"output\": \"Error\", \"exitcode\": 1}", actualOutput)
	})

	t.Run("A 401 response should be an UNKNOWN exit code", func(t *testing.T) {
		// We manipuate the Args to set them up for the testcases
		// After this test we restore the initial args
		oldArgs := os.Args
		defer func() { os.Args = oldArgs }()

		flag.CommandLine = flag.NewFlagSet("flag", flag.ExitOnError)

		os.Args = []string{
			"main.exe",
			"-host", "remotehost",
			"-username", "thisismyusername",
			"-password", "thisismypassword",
			"-executable", "/path/to/executable",
			"-executableArg", "-s",
			"-script", "TestScript-Valid.ps1",
		}
		httpClient := httpclient.NewMockHTTPClient(`{"output": "Error", "exitcode": 1}`, 401)

		var buf bytes.Buffer
		actualExit := invokeClient(&buf, httpClient)
		actualOutput := buf.String()

		assert.Equal(t, `{"args":["-s"],"path":"/path/to/executable","scriptarguments":[],"stdin":"Write-Host \"This is a test script\"\r\n\r\n","timeout":"10s"}`, httpClient.RequestBodyContent)
		assert.Equal(t, 10*time.Second, httpClient.Timeout)
		assert.Equal(t, false, httpClient.Transport.TLSClientConfig.InsecureSkipVerify)
		assert.Equal(t, "Basic dGhpc2lzbXl1c2VybmFtZTp0aGlzaXNteXBhc3N3b3Jk", httpClient.RequestHeaders["Authorization"][0])
		assert.Equal(t, "remotehost:9000", httpClient.RequestHost)
		assert.Equal(t, "/v1/runscriptstdin", httpClient.RequestURI.Path)
		assert.Equal(t, "POST", httpClient.RequestVerb)

		assert.Equal(t, 3, actualExit)
		assert.Equal(t, "Response code: 401\n{\"output\": \"Error\", \"exitcode\": 1}", actualOutput)
	})

	t.Run("Powershell scripts that don't end with 2 newlines should be rejected", func(t *testing.T) {
		oldArgs := os.Args
		defer func() { os.Args = oldArgs }()
		flag.CommandLine = flag.NewFlagSet("flag", flag.ExitOnError)
		os.Args = []string{
			"main.exe",
			"-host", "remotehost",
			"-username", "thisismyusername",
			"-password", "thisismypassword",
			"-executable", "/path/to/executable",
			"-script", "TestScript-Invalid.ps1",
			"-insecure",
		}

		httpClient := httpclient.NewMockHTTPClient(`{}`, 200)

		var buf bytes.Buffer
		actualExit := invokeClient(&buf, httpClient)
		actualOutput := buf.String()

		assert.Equal(t, 3, actualExit)
		assert.Equal(t, "Invalid powershell script, the script must end with two blank lines", actualOutput)
	})

	t.Run("Powershell scripts with 2 unix line endings are just as valid as windows line endings", func(t *testing.T) {
		// We manipuate the Args to set them up for the testcases
		// After this test we restore the initial args
		oldArgs := os.Args
		defer func() { os.Args = oldArgs }()

		flag.CommandLine = flag.NewFlagSet("flag", flag.ExitOnError)

		os.Args = []string{
			"main.exe",
			"-host", "remotehost",
			"-username", "thisismyusername",
			"-password", "thisismypassword",
			"-executable", "/path/to/executable",
			"-script", "TestScript-Valid-UnixLineEndings.ps1",
		}
		httpClient := httpclient.NewMockHTTPClient(`{"output": "Test output", "exitcode": 2}`, 200)

		var buf bytes.Buffer
		actualExit := invokeClient(&buf, httpClient)
		actualOutput := buf.String()

		assert.Equal(t, `{"args":null,"path":"/path/to/executable","scriptarguments":[],"stdin":"Write-Host \"This is a test script\"\n\n","timeout":"10s"}`, httpClient.RequestBodyContent)
		assert.Equal(t, 10*time.Second, httpClient.Timeout)
		assert.Equal(t, false, httpClient.Transport.TLSClientConfig.InsecureSkipVerify)
		assert.Equal(t, "Basic dGhpc2lzbXl1c2VybmFtZTp0aGlzaXNteXBhc3N3b3Jk", httpClient.RequestHeaders["Authorization"][0])
		assert.Equal(t, "remotehost:9000", httpClient.RequestHost)
		assert.Equal(t, "/v1/runscriptstdin", httpClient.RequestURI.Path)
		assert.Equal(t, "POST", httpClient.RequestVerb)

		assert.Equal(t, 2, actualExit)
		assert.Equal(t, "Test output", actualOutput)
	})
}
