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
	//template := flag.String("template", "", "pnp4nagios template")
	hostname := flag.String("hostname", "", "hostname or ip")
	port := flag.Int("port", 0, "port number")
	cacert := flag.String("cacert", "", "CA certificate")
	certificate := flag.String("certificate", "", "certificate file")
	key := flag.String("key", "", "key file")
	username := flag.String("username", "", "username")
	password := flag.String("password", "", "password")
	executable := flag.String("executable", "", "executable path")
	executableArg := flag.String("executableArg", "", "executable arg for multiple specify multiple times")
	script := flag.String("script", "", "script location")
	timeout := flag.String("timeout", "10s", "timeout (e.g. 10s)")
	flag.Parse()

	scriptContent, err := ioutil.ReadFile(*script) // the file is inside the local directory
	if err != nil {
		fmt.Println("Err")
	}

	scriptSignatureContent, err := ioutil.ReadFile(*script + ".minisig") // the file is inside the local directory
	if err != nil {
		fmt.Println("Err")
	}

	input := map[string]interface{}{
		"path":            executable,
		"args":            executableArg,
		"stdinsignature":  scriptSignatureContent,
		"stdin":           scriptContent,
		"scriptarguments": flag.Args(),
		"timeout":         timeout,
	}

	byteArray, _ := json.Marshal(input)
	byteArrayBuffer := bytes.NewBuffer(byteArray)

	url := fmt.Sprintf("https://%s:%d/v1/runscriptstdin", *hostname, *port)

	client := &http.Client{
		Timeout: time.Second * 10,
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
		fmt.Println(fmt.Errorf("Got error %s", err.Error()))
	}
	req.SetBasicAuth(*username, *password)
	response, err := client.Do(req)
	if err != nil {
		fmt.Println(fmt.Errorf("Got error %s", err.Error()))
	}
	defer response.Body.Close()

	if response.StatusCode != 200 {
		fmt.Println("Response code: " + response.Status)
	}

	type MAResponse struct {
		Output   string
		Exitcode int
	}

	var decodedResponse MAResponse
	json.NewDecoder(response.Body).Decode(decodedResponse)
	fmt.Println(decodedResponse.Output)
	if decodedResponse.Exitcode > 3 {
		os.Exit(3)
	}

	os.Exit(decodedResponse.Exitcode)
}
