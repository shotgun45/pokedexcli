package main

import (
	"os"
	"os/exec"
	"strings"
	"testing"
)

func TestCommandExit_ExitsProcess(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test in short mode.")
	}

	cmd := exec.Command(os.Args[0], "-test.run=TestHelperProcessExit")
	cmd.Env = append(cmd.Env, "GO_WANT_HELPER_PROCESS_EXIT=1")
	output, err := cmd.CombinedOutput()

	if err != nil {
		t.Fatalf("expected exit code 0 (no error), got %v, output: %s", err, string(output))
	}

	if !strings.Contains(string(output), "Closing the Pokedex... Goodbye!") {
		t.Errorf("expected exit message, got: %s", string(output))
	}
}

func TestHelperProcessExit(t *testing.T) {
	if os.Getenv("GO_WANT_HELPER_PROCESS_EXIT") != "1" {
		return
	}
	cfg := &config{}
	commandExit(cfg)
}
