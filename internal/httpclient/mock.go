package httpclient

import (
	"io"
	"net/http"
	"strings"
	"time"
)

type mock struct {
	// DoFunc will be executed whenever Do function is executed
	// so we'll be able to create a custom response
	DoFunc    func(*http.Request) (*http.Response, error)
	Timeout   time.Duration
	Transport *http.Transport
}

func (H mock) Do(r *http.Request) (*http.Response, error) {
	return &http.Response{
		Body:       io.NopCloser(strings.NewReader(`{"output": "Test output", "exitcode": 2}`)),
		StatusCode: 200,
	}, nil
}

func (H mock) SetTimeout(timeout time.Duration) {
	H.Timeout = timeout
}

func (H mock) SetTransport(transport *http.Transport) {
	H.Transport = transport
}

func NewMockHTTPClient() *mock {
	client := new(mock)
	return client
}
