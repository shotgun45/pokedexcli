package main

import (
	"bytes"
	"os"
	"testing"
)

type fakeCache struct {
	store map[string][]byte
}

func TestCommandExplore(t *testing.T) {
	pokemonJSON := `{"pokemon_encounters": [
		{"pokemon": {"name": "pikachu"}},
		{"pokemon": {"name": "bulbasaur"}}
	]}`
	cache := &fakeCache{store: map[string][]byte{
		"https://pokeapi.co/api/v2/location-area/test-area": []byte(pokemonJSON),
	}}
	cfg := &config{cache: cache}

	// Capture os.Stdout
	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	err := commandExplore(cfg, "test-area")
	w.Close()
	os.Stdout = old

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	var buf bytes.Buffer
	buf.ReadFrom(r)
	output := buf.String()
	if !contains(output, "pikachu") || !contains(output, "bulbasaur") {
		t.Errorf("output missing expected pokemon: %s", output)
	}
}

func contains(s, substr string) bool {
	return bytes.Contains([]byte(s), []byte(substr))
}
