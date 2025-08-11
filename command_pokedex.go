package main

import (
	"fmt"
	"sort"
)

func commandPokedex(cfg *config, args ...string) error {
	if len(cfg.pokedex) == 0 {
		fmt.Println("Your Pokedex is empty.")
		return nil
	}
	fmt.Println("Your Pokedex:")
	names := make([]string, 0, len(cfg.pokedex))
	for name := range cfg.pokedex {
		names = append(names, name)
	}
	sort.Strings(names)
	for _, name := range names {
		fmt.Printf(" - %s\n", name)
	}
	return nil
}
