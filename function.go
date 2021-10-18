package main

import (
	"errors"
	"fmt"
	"io"
	"os"
	"time"
)

func (i *executableArguments) String() string {
	// change this, this is just can example to satisfy the interface
	return "my string representation"
}

func (i *executableArguments) Set(value string) error {
	*i = append(*i, value)
	return nil
}

type MonitoringAgentResponse struct {
	Output   string `json:"output"`
	Exitcode int    `json:"exitcode"`
}

type executableArguments []string

const okExitCode = 0
const warningExitCode = 1
const criticalExitCode = 2
const unknownExitCode = 3

func FileExists(name string) bool {
	_, err := os.Stat(name)
	if err == nil {
		return true
	}
	if errors.Is(err, os.ErrNotExist) {
		return false
	}
	return false
}

func enableTimeout(timeout string) time.Duration {
	timeoutDuration, timeoutParseError := time.ParseDuration(timeout)
	if timeoutParseError != nil {
		panic(fmt.Errorf("error parsing timeout value %s", timeoutParseError.Error()))
	}

	time.AfterFunc(timeoutDuration, func() {
		panic(fmt.Sprintf("Client timeout reached: %s\n", timeoutDuration))
	})
	return timeoutDuration
}

func die(stdout io.Writer, message string) int {
	fmt.Fprint(stdout, message)
	return unknownExitCode
}
