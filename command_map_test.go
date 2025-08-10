package main

import (
	"bytes"
	"io"
	"net/http"
	"os"
	"strings"
	"testing"
)

type mockRoundTripper struct {
	responseBody string
}

func (m *mockRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	resp := &http.Response{
		StatusCode: 200,
		Body:       io.NopCloser(strings.NewReader(m.responseBody)),
		Header:     make(http.Header),
	}
	return resp, nil
}

func TestCommandMap_PrintsLocationAreas(t *testing.T) {
	// Mock response JSON
	json := `{"results":[{"name":"area-1"},{"name":"area-2"}],"next":null,"previous":null}`

	// Patch http.DefaultClient.Transport
	oldTransport := http.DefaultTransport
	http.DefaultTransport = &mockRoundTripper{responseBody: json}
	defer func() { http.DefaultTransport = oldTransport }()

	cfg := &config{}

	// Capture stdout
	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	err := commandMap(cfg)
	w.Close()
	os.Stdout = old

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	var buf bytes.Buffer
	buf.ReadFrom(r)
	output := buf.String()

	if !strings.Contains(output, "area-1") || !strings.Contains(output, "area-2") {
		t.Errorf("expected output to contain area names, got: %s", output)
	}
}
