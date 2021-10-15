package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"ioutil"
	"net/http"
	"time"
)

func main() {
	fmt.Println("Hello, world.")
	//template := flag.String("template", "", "pnp4nagios template")
	hostname := flag.String("hostname", "", "hostname or ip")
	port := flag.Int("port", 0, "port number")
	//cacert := flag.String("cacert", "", "CA certificate")
	//certificate := flag.String("certificate", "", "certificate file")
	//key := flag.String("key", "", "key file")
	username := flag.String("username", "", "username")
	password := flag.String("password", "", "password")
	executable := flag.String("executable", "", "executable path")
	executableArg := flag.String("executableArg", "", "executable arg for multiple specify multiple times")
	script := flag.String("script", "", "script location")
	timeout := flag.String("timeout", "", "timeout (e.g. 10s)")
	flag.Parse()

	scriptContent, err := ioutil.ReadFile(script) // the file is inside the local directory
	if err != nil {
		fmt.Println("Err")
	}

	input := map[string]interface{}{
		"path":            executable,
		"args":            executableArg,
		"stdin":           scriptContent,
		"scriptarguments": flag.Args(),
		"timeout":         timeout,
	}

	byteArray, _ := json.Marshal(input)
	byteArrayBuffer := bytes.NewBuffer(byteArray)

	url := fmt.Sprintf("https://%s:%s/v1/runscriptstdin", hostname, port)

	client := &http.Client{
		Timeout: time.Second * 10,
	}
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
}
