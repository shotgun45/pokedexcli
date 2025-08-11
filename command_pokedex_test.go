package main

import (
	"bytes"
	"os"
	"strings"
	"testing"
)

func TestCommandPokedex_Empty(t *testing.T) {
	cfg := &config{pokedex: make(map[string]Pokemon)}

	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	err := commandPokedex(cfg)
	w.Close()
	os.Stdout = old
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	var buf bytes.Buffer
	buf.ReadFrom(r)
	output := buf.String()
	if !strings.Contains(output, "Your Pokedex is empty.") {
		t.Errorf("expected empty message, got: %s", output)
	}
}

func TestCommandPokedex_WithPokemon(t *testing.T) {
	cfg := &config{pokedex: map[string]Pokemon{
		"pidgey":   {Name: "pidgey"},
		"caterpie": {Name: "caterpie"},
	}}

	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	err := commandPokedex(cfg)
	w.Close()
	os.Stdout = old
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	var buf bytes.Buffer
	buf.ReadFrom(r)
	output := buf.String()
	if !strings.Contains(output, "Your Pokedex:") {
		t.Errorf("expected pokedex header, got: %s", output)
	}
	if !strings.Contains(output, "- pidgey") || !strings.Contains(output, "- caterpie") {
		t.Errorf("expected pokemon names in output, got: %s", output)
	}
}
