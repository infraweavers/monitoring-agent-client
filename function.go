package main

import (
	"errors"
	"os"
)

func (i *executableArguments) String() string {
	// change this, this is just can example to satisfy the interface
	return "my string representation"
}

func (i *executableArguments) Set(value string) error {
	*i = append(*i, value)
	return nil
}

type MAResponse struct {
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
