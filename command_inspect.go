package main

import (
	"fmt"
	"strings"
)

func commandInspect(cfg *config, args ...string) error {
	if len(args) < 1 {
		return fmt.Errorf("usage: inspect <pokemon>")
	}
	name := strings.ToLower(args[0])
	if cfg.pokedex == nil {
		fmt.Println("you have not caught that pokemon")
		return nil
	}
	pkmn, ok := cfg.pokedex[name]
	if !ok {
		fmt.Println("you have not caught that pokemon")
		return nil
	}
	fmt.Printf("Name: %s\n", pkmn.Name)
	fmt.Printf("Height: %d\n", pkmn.Height)
	fmt.Printf("Weight: %d\n", pkmn.Weight)
	fmt.Println("Stats:")
	for _, stat := range pkmn.Stats {
		fmt.Printf("  -%s: %d\n", stat.Stat.Name, stat.BaseStat)
	}
	fmt.Println("Types:")
	for _, t := range pkmn.Types {
		fmt.Printf("  - %s\n", t.Type.Name)
	}
	return nil
}
