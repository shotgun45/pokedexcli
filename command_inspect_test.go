package main

import (
	"bytes"
	"os"
	"strings"
	"testing"
)

func TestCommandInspect_NotCaught(t *testing.T) {
	cfg := &config{pokedex: make(map[string]Pokemon)}

	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	err := commandInspect(cfg, "pidgey")
	w.Close()
	os.Stdout = old
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	var buf bytes.Buffer
	buf.ReadFrom(r)
	output := buf.String()
	if !strings.Contains(output, "you have not caught that pokemon") {
		t.Errorf("expected not caught message, got: %s", output)
	}
}

func TestCommandInspect_Caught(t *testing.T) {
	cfg := &config{pokedex: map[string]Pokemon{
		"pidgey": {
			Name:   "pidgey",
			Height: 3,
			Weight: 18,
			Stats: []struct {
				BaseStat int `json:"base_stat"`
				Stat     struct {
					Name string `json:"name"`
				} `json:"stat"`
			}{
				{40, struct {
					Name string `json:"name"`
				}{"hp"}},
				{45, struct {
					Name string `json:"name"`
				}{"attack"}},
				{40, struct {
					Name string `json:"name"`
				}{"defense"}},
				{35, struct {
					Name string `json:"name"`
				}{"special-attack"}},
				{35, struct {
					Name string `json:"name"`
				}{"special-defense"}},
				{56, struct {
					Name string `json:"name"`
				}{"speed"}},
			},
			Types: []struct {
				Type struct {
					Name string `json:"name"`
				} `json:"type"`
			}{
				{struct {
					Name string `json:"name"`
				}{"normal"}},
				{struct {
					Name string `json:"name"`
				}{"flying"}},
			},
		},
	}}

	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	err := commandInspect(cfg, "pidgey")
	w.Close()
	os.Stdout = old
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	var buf bytes.Buffer
	buf.ReadFrom(r)
	output := buf.String()
	if !strings.Contains(output, "Name: pidgey") {
		t.Errorf("expected name in output, got: %s", output)
	}
	if !strings.Contains(output, "Height: 3") {
		t.Errorf("expected height in output, got: %s", output)
	}
	if !strings.Contains(output, "Weight: 18") {
		t.Errorf("expected weight in output, got: %s", output)
	}
	if !strings.Contains(output, "-hp: 40") {
		t.Errorf("expected hp stat in output, got: %s", output)
	}
	if !strings.Contains(output, "-attack: 45") {
		t.Errorf("expected attack stat in output, got: %s", output)
	}
	if !strings.Contains(output, "-defense: 40") {
		t.Errorf("expected defense stat in output, got: %s", output)
	}
	if !strings.Contains(output, "-special-attack: 35") {
		t.Errorf("expected special-attack stat in output, got: %s", output)
	}
	if !strings.Contains(output, "-special-defense: 35") {
		t.Errorf("expected special-defense stat in output, got: %s", output)
	}
	if !strings.Contains(output, "-speed: 56") {
		t.Errorf("expected speed stat in output, got: %s", output)
	}
	if !strings.Contains(output, "- normal") || !strings.Contains(output, "- flying") {
		t.Errorf("expected types in output, got: %s", output)
	}
}
