package main

import (
	"bytes"
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"time"
)

func main() {

	_ = flag.String("template", "", "pnp4nagios template")

	hostname := flag.String("host", "", "hostname or ip")
	port := flag.Int("port", -1, "port number")
	username := flag.String("username", os.Getenv("MONITORING_AGENT_USERNAME"), "username")
	password := flag.String("password", os.Getenv("MONITORING_AGENT_PASSWORD"), "password")
	executable := flag.String("executable", "", "executable path")
	script := flag.String("script", "", "script location")

	cacertificateFilePath := flag.String("cacert", os.Getenv("MONITORING_AGENT_CA_CERTIFICATE_PATH"), "CA certificate")
	certificateFilePath := flag.String("certificate", os.Getenv("MONITORING_AGENT_CLIENT_CERTIFICATE_PATH"), "certificate file")
	privateKeyFilePath := flag.String("key", os.Getenv("MONITORING_AGENT_CLIENT_KEY_PATH"), "key file")
	timeout := flag.String("timeout", "10s", "timeout (e.g. 10s)")
	makeInsecure := flag.Bool("insecure", false, "ignore TLS Certificate checks")

	var executableArgs executableArguments
	flag.Var(&executableArgs, "executableArg", "executable arg for multiple specify multiple times")

	flag.Parse()

	if *hostname == "" {
		panic("hostname is not set")
	}
	if *password == "" {
		panic("password is not set")
	}
	if *executable == "" {
		panic("executable is not set")
	}
	if *script == "" {
		panic("script is not set")
	}

	timeoutDuration, timeoutParseError := time.ParseDuration(*timeout)
	if timeoutParseError != nil {
		panic(fmt.Errorf("error parsing timeout value %s", timeoutParseError.Error()))
	}

	time.AfterFunc(timeoutDuration, func() {
		panic(fmt.Sprintf("Client timeout reached: %s\n", timeoutDuration))
	})

	scriptContent, err := ioutil.ReadFile(*script)
	if err != nil {
		panic(fmt.Sprintf("error, could not load script file: %s\n", err))
	}

	restRequest := map[string]interface{}{
		"path":            executable,
		"args":            executableArgs,
		"stdin":           string(scriptContent),
		"scriptarguments": flag.Args(),
		"timeout":         timeout,
	}

	scriptSignatureFilename := fmt.Sprintf("%s%s", *script, ".minisig")
	if FileExists(scriptSignatureFilename) {
		scriptSignatureContent, err := ioutil.ReadFile(scriptSignatureFilename)
		if err != nil {
			panic(fmt.Sprintf("error loading script signature: %s", err))
		}
		restRequest["stdinsignature"] = scriptSignatureContent
	}

	byteArray, _ := json.Marshal(restRequest)
	byteArrayBuffer := bytes.NewBuffer(byteArray)

	url := fmt.Sprintf("https://%s:%d/v1/runscriptstdin", *hostname, *port)

	client := &http.Client{
		Timeout: timeoutDuration,
	}
	transport := new(http.Transport)
	transport.TLSClientConfig = &tls.Config{
		InsecureSkipVerify: *makeInsecure,
	}

	if *certificateFilePath != "" && *privateKeyFilePath != "" {
		certificateToLoad, err := tls.LoadX509KeyPair(*certificateFilePath, *privateKeyFilePath)
		if err != nil {
			panic(fmt.Errorf("error loading certificate pair %s", err.Error()))
		}
		transport.TLSClientConfig.Certificates = []tls.Certificate{certificateToLoad}
	}

	if *cacertificateFilePath != "" {
		caCertificate, err := ioutil.ReadFile(*cacertificateFilePath)
		if err != nil {
			panic(fmt.Errorf("error loading ca certificate %s", err.Error()))
		}
		CACertificatePool := x509.NewCertPool()
		CACertificatePool.AppendCertsFromPEM(caCertificate)
		transport.TLSClientConfig.RootCAs = CACertificatePool
	}

	client.Transport = transport

	req, err := http.NewRequest(http.MethodPost, url, byteArrayBuffer)
	if err != nil {
		panic(fmt.Errorf("got error %s", err.Error()))
	}
	req.SetBasicAuth(*username, *password)

	response, err := client.Do(req)

	if err != nil {
		panic(fmt.Errorf("got error %s", err.Error()))
	}

	defer response.Body.Close()

	if response.StatusCode != 200 {
		fmt.Printf("Response code: %s\n%#v", response.Status, response.Body)
		os.Exit(unknownExitCode)
	}

	var decodedResponse MAResponse

	decoder := json.NewDecoder(response.Body)
	decoder.DisallowUnknownFields()
	decoder.Decode(&decodedResponse)

	fmt.Print(decodedResponse.Output)

	if decodedResponse.Exitcode > unknownExitCode {
		os.Exit(unknownExitCode)
	}

	os.Exit(decodedResponse.Exitcode)
}
