package main

import (
	"bytes"
	"os"
	"strings"
	"testing"
)

func (f *fakeCache) Get(key string) ([]byte, bool) {
	v, ok := f.store[key]
	return v, ok
}
func (f *fakeCache) Add(key string, val []byte) {
	f.store[key] = val
}

func TestCommandCatch_Caught(t *testing.T) {
	// base_experience = 10, so catch chance is very high
	pokeJSON := `{"name": "pikachu", "base_experience": 10}`
	cache := &fakeCache{store: map[string][]byte{
		"https://pokeapi.co/api/v2/pokemon/pikachu": []byte(pokeJSON),
	}}
	cfg := &config{cache: cache, pokedex: make(map[string]Pokemon)}

	// Patch rand.Float64 to always return 0.0 (guaranteed catch)
	oldRand := randFloat64
	randFloat64 = func() float64 { return 0.0 }
	defer func() { randFloat64 = oldRand }()

	// Capture output
	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	err := commandCatch(cfg, "pikachu")
	w.Close()
	os.Stdout = old
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	var outBuf bytes.Buffer
	outBuf.ReadFrom(r)
	output := outBuf.String()
	if !strings.Contains(output, "Throwing a Pokeball at pikachu") {
		t.Errorf("output missing throw message: %s", output)
	}
	if !strings.Contains(output, "was caught!") {
		t.Errorf("output missing caught message: %s", output)
	}
	if _, ok := cfg.pokedex["pikachu"]; !ok {
		t.Errorf("pikachu not in pokedex after catch")
	}
}

func TestCommandCatch_Escaped(t *testing.T) {
	// base_experience = 10, but force rand to always return 1.0 (guaranteed escape)
	pokeJSON := `{"name": "pikachu", "base_experience": 10}`
	cache := &fakeCache{store: map[string][]byte{
		"https://pokeapi.co/api/v2/pokemon/pikachu": []byte(pokeJSON),
	}}
	cfg := &config{cache: cache, pokedex: make(map[string]Pokemon)}

	// Patch rand.Float64 to always return 1.0 (guaranteed escape)
	oldRand := randFloat64
	randFloat64 = func() float64 { return 1.0 }
	defer func() { randFloat64 = oldRand }()

	// Capture output
	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	err := commandCatch(cfg, "pikachu")
	w.Close()
	os.Stdout = old
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	var outBuf bytes.Buffer
	outBuf.ReadFrom(r)
	output := outBuf.String()
	if !strings.Contains(output, "Throwing a Pokeball at pikachu") {
		t.Errorf("output missing throw message: %s", output)
	}
	if !strings.Contains(output, "escaped!") {
		t.Errorf("output missing escaped message: %s", output)
	}
	if _, ok := cfg.pokedex["pikachu"]; ok {
		t.Errorf("pikachu should not be in pokedex after escape")
	}
}
