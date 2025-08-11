package main

import (
	"bytes"
	"io"
	"net/http"
	"os"
	"strings"
	"testing"
	"time"

	"pokedexcli/internal/pokecache"
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

	cfg := &config{
		cache: pokecache.NewCache(time.Second), // or any short duration
	}

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

func TestCommandMapb_PrintsPreviousPage(t *testing.T) {
	json := `{"results":[{"name":"prev-area-1"},{"name":"prev-area-2"}],"next":null,"previous":null}`
	oldTransport := http.DefaultTransport
	http.DefaultTransport = &mockRoundTripper{responseBody: json}
	defer func() { http.DefaultTransport = oldTransport }()

	prevURL := "http://example.com/prev"
	cfg := &config{
		cache:            pokecache.NewCache(time.Second),
		prevLocationsURL: &prevURL,
	}
	cfg.cache.Add(prevURL, []byte(json)) // Pre-populate cache for prevURL

	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	err := commandMapb(cfg)
	w.Close()
	os.Stdout = old

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	var buf bytes.Buffer
	buf.ReadFrom(r)
	output := buf.String()

	if !strings.Contains(output, "prev-area-1") || !strings.Contains(output, "prev-area-2") {
		t.Errorf("expected output to contain previous area names, got: %s", output)
	}
}

func TestCommandMapb_FirstPageMessage(t *testing.T) {
	cfg := &config{
		cache:            pokecache.NewCache(time.Second),
		prevLocationsURL: nil,
	}

	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	err := commandMapb(cfg)
	w.Close()
	os.Stdout = old

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	var buf bytes.Buffer
	buf.ReadFrom(r)
	output := buf.String()

	if !strings.Contains(output, "you're on the first page") {
		t.Errorf("expected output to contain 'you're on the first page', got: %s", output)
	}
}
