package main

import (
	"bytes"
	"os"
	"strings"
	"testing"
)

func TestCommandHelp_PrintsHelpMessage(t *testing.T) {
	cfg := &config{}

	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	err := commandHelp(cfg)
	w.Close()
	os.Stdout = old

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	var buf bytes.Buffer
	buf.ReadFrom(r)
	output := buf.String()

	if !strings.Contains(output, "Welcome to the Pokedex!") {
		t.Errorf("help output missing welcome message")
	}
	if !strings.Contains(output, "help: Displays a help message") {
		t.Errorf("help output missing help command description")
	}
	if !strings.Contains(output, "exit: Exit the Pokedex") {
		t.Errorf("help output missing exit command description")
	}
	if !strings.Contains(output, "map: Displays the next 20 location areas in the Pokemon world") {
		t.Errorf("help output missing map command description")
	}
	if !strings.Contains(output, "mapb: Displays the previous 20 location areas in the Pokemon world") {
		t.Errorf("help output missing mapb command description")
	}
}
