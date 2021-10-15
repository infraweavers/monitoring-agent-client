package main

import (
	"bytes"
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"
)

func main() {
	var executableArgs executableArgsType
	_ = flag.String("template", "", "pnp4nagios template")
	hostname := flag.String("host", "", "hostname or ip")
	port := flag.Int("port", 0, "port number")
	cacert := flag.String("cacert", "", "CA certificate")
	certificate := flag.String("certificate", "", "certificate file")
	key := flag.String("key", "", "key file")
	username := flag.String("username", "", "username")
	password := flag.String("password", "", "password")
	executable := flag.String("executable", "", "executable path")
	flag.Var(&executableArgs, "executableArg", "executable arg for multiple specify multiple times")
	script := flag.String("script", "", "script location")
	timeout := flag.String("timeout", "10s", "timeout (e.g. 10s)")
	flag.Parse()

	timeoutDuration, _ := time.ParseDuration(*timeout)

	time.AfterFunc(timeoutDuration, func() {
		fmt.Print("Timeout reached.")
		os.Exit(3)
	})

	scriptContent, err := ioutil.ReadFile(*script) // the file is inside the local directory
	if err != nil {
		fmt.Println("Err")
	}

	sigFile := fmt.Sprintf("%s%s", *script, ".minisig")

	scriptSignatureContent, err := ioutil.ReadFile(sigFile) // the file is inside the local directory
	if err != nil {
		fmt.Println("Err")
	}

	input := map[string]interface{}{
		"path":            executable,
		"args":            executableArgs,
		"stdinsignature":  string(scriptSignatureContent),
		"stdin":           string(scriptContent),
		"scriptarguments": flag.Args(),
		"timeout":         timeout,
	}

	byteArray, _ := json.Marshal(input)
	byteArrayBuffer := bytes.NewBuffer(byteArray)

	url := fmt.Sprintf("https://%s:%d/v1/runscriptstdin", *hostname, *port)

	client := &http.Client{
		Timeout: timeoutDuration,
	}

	certificateToLoad, _ := tls.LoadX509KeyPair(*certificate, *key)

	certificatesCollection := []tls.Certificate{certificateToLoad}

	caCert, err := ioutil.ReadFile(*cacert)
	if err != nil {
		log.Fatal(err)
	}
	caCertPool := x509.NewCertPool()
	caCertPool.AppendCertsFromPEM(caCert)

	transport := new(http.Transport)
	transport.TLSClientConfig = &tls.Config{
		InsecureSkipVerify: false,
		Certificates:       certificatesCollection,
		RootCAs:            caCertPool,
	}

	client.Transport = transport

	req, err := http.NewRequest(http.MethodPost, url, byteArrayBuffer)
	if err != nil {
		fmt.Println(fmt.Errorf("got error %s", err.Error()))
	}
	req.SetBasicAuth(*username, *password)

	response, err := client.Do(req)
	if err != nil {
		fmt.Println(fmt.Errorf("got error %s", err.Error()))
	}

	defer response.Body.Close()

	if response.StatusCode != 200 {
		fmt.Println("Response code: " + response.Status)
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
