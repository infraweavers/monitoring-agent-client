package main

func (i *executableArgsType) String() string {
	// change this, this is just can example to satisfy the interface
	return "my string representation"
}

func (i *executableArgsType) Set(value string) error {
	*i = append(*i, value)
	return nil
}

type MAResponse struct {
	Output   string `json:"output"`
	Exitcode int    `json:"exitcode"`
}

type executableArgsType []string

const okExitCode = 0
const warningExitCode = 1
const criticalExitCode = 2
const unknownExitCode = 3
