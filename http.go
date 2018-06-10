package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

// HTTPTester represents a HTTP transport tester.
type HTTPTester struct {
	client http.Client
	url    string
	br     *bytes.Reader
}

// NewHTTPTester returns an instantiated HTTP tester.
func NewHTTPTester(addr string) *HTTPTester {
	return &HTTPTester{
		client: http.Client{},
		url:    fmt.Sprintf("http://%s:/key", addr),
	}
}

// Prepare prepares the tester for execution.
func (h *HTTPTester) Prepare(kv map[string]string) error {
	b, err := json.Marshal(kv)
	if err != nil {
		return err
	}
	h.br = bytes.NewReader(b)

	return nil
}

// Once executes a single test request.
func (h *HTTPTester) Once() (time.Duration, error) {
	h.br.Seek(0, io.SeekStart)

	start := time.Now()
	resp, err := h.client.Post(h.url, "application/json", h.br)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return 0, fmt.Errorf("received %s", resp.Status)
	}
	dur := time.Since(start)

	return dur, nil
}
