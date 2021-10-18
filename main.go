package main

import (
	"bytes"
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"monitoring-agent-client/internal/httpclient"
	"net/http"
	"os"
	"strings"
)

func main() {
	httpClient := httpclient.NewHTTPClient()
	os.Exit(invokeClient(os.Stdout, httpClient))
}

func invokeClient(stdout io.Writer, httpClient httpclient.Interface) int {
	_ = flag.String("template", "", "pnp4nagios template")

	hostname := flag.String("host", "", "hostname or ip")
	port := flag.Int("port", 9000, "port number")
	username := flag.String("username", os.Getenv("MONITORING_AGENT_USERNAME"), "username")
	password := flag.String("password", os.Getenv("MONITORING_AGENT_PASSWORD"), "password")
	executable := flag.String("executable", "", "executable path")
	script := flag.String("script", "", "script location")

	cacertificateFilePath := flag.String("cacert", os.Getenv("MONITORING_AGENT_CA_CERTIFICATE_PATH"), "CA certificate")
	certificateFilePath := flag.String("certificate", os.Getenv("MONITORING_AGENT_CLIENT_CERTIFICATE_PATH"), "certificate file")
	privateKeyFilePath := flag.String("key", os.Getenv("MONITORING_AGENT_CLIENT_KEY_PATH"), "key file")
	timeoutString := flag.String("timeout", "10s", "timeout (e.g. 10s)")
	makeInsecure := flag.Bool("insecure", false, "ignore TLS Certificate checks")

	var executableArgs executableArguments
	flag.Var(&executableArgs, "executableArg", "executable arg for multiple specify multiple times")

	flag.Parse()

	if *hostname == "" {
		return die(stdout, "hostname is not set")
	}
	if *password == "" {
		return die(stdout, "password is not set")
	}
	if *executable == "" {
		return die(stdout, "executable is not set")
	}
	if *script == "" {
		return die(stdout, "script is not set")
	}

	timeout := enableTimeout(*timeoutString)

	scriptContentByteArray, err := ioutil.ReadFile(*script)
	if err != nil {
		return die(stdout, fmt.Sprintf("error, could not load script file: %s\n", err))
	}
	scriptContent := string(scriptContentByteArray)

	if strings.HasSuffix(*script, ".ps1") {
		if strings.HasSuffix(scriptContent, "\r\n\r\n") == false && strings.HasSuffix(scriptContent, "\n\n") == false {
			return die(stdout, "Invalid powershell script, the script must end with two blank lines")
		}
	}

	restRequest := map[string]interface{}{
		"path":            executable,
		"args":            executableArgs,
		"stdin":           scriptContent,
		"scriptarguments": flag.Args(),
		"timeout":         timeoutString,
	}

	byteArray, _ := json.Marshal(restRequest)
	byteArrayBuffer := bytes.NewBuffer(byteArray)

	url := fmt.Sprintf("https://%s:%d/v1/runscriptstdin", *hostname, *port)

	httpClient.SetTimeout(timeout)

	transport := new(http.Transport)
	transport.TLSClientConfig = &tls.Config{
		InsecureSkipVerify: *makeInsecure,
	}

	if *certificateFilePath != "" && *privateKeyFilePath != "" {
		certificateToLoad, err := tls.LoadX509KeyPair(*certificateFilePath, *privateKeyFilePath)
		if err != nil {
			return die(stdout, fmt.Sprintf("error loading certificate pair %s", err.Error()))
		}
		transport.TLSClientConfig.Certificates = []tls.Certificate{certificateToLoad}
	}

	if *cacertificateFilePath != "" {
		caCertificate, err := ioutil.ReadFile(*cacertificateFilePath)
		if err != nil {
			return die(stdout, fmt.Sprintf("error loading ca certificate %s", err.Error()))
		}
		CACertificatePool := x509.NewCertPool()
		CACertificatePool.AppendCertsFromPEM(caCertificate)
		transport.TLSClientConfig.RootCAs = CACertificatePool
	}

	scriptSignatureFilename := fmt.Sprintf("%s%s", *script, ".minisig")
	if FileExists(scriptSignatureFilename) {
		scriptSignatureContent, err := ioutil.ReadFile(scriptSignatureFilename)
		if err != nil {
			return die(stdout, fmt.Sprintf("error loading script signature: %s", err.Error()))
		}
		restRequest["stdinsignature"] = scriptSignatureContent
	}

	httpClient.SetTransport(transport)

	req, err := http.NewRequest(http.MethodPost, url, byteArrayBuffer)
	if err != nil {
		panic(fmt.Errorf("got http request error %s", err.Error()))
	}
	req.SetBasicAuth(*username, *password)

	response, err := httpClient.Do(req)

	if err != nil {
		return die(stdout, fmt.Sprintf("got httpClient error %s", err.Error()))
	}

	defer response.Body.Close()

	if response.StatusCode != 200 {
		return die(stdout, fmt.Sprintf("Response code: %s\n%#v", response.Status, response.Body))
	}

	var decodedResponse MonitoringAgentResponse

	decoder := json.NewDecoder(response.Body)
	decoder.DisallowUnknownFields()
	decoder.Decode(&decodedResponse)

	fmt.Fprint(stdout, decodedResponse.Output)

	if decodedResponse.Exitcode > unknownExitCode {
		return unknownExitCode
	}

	return decodedResponse.Exitcode
}
