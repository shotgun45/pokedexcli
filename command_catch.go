package main

import (
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"strings"
)

var randFloat64 = func() float64 { return rand.Float64() }

type Pokemon struct {
	Name           string `json:"name"`
	BaseExperience int    `json:"base_experience"`
	Height         int    `json:"height"`
	Weight         int    `json:"weight"`
	Stats          []struct {
		BaseStat int `json:"base_stat"`
		Stat     struct {
			Name string `json:"name"`
		} `json:"stat"`
	} `json:"stats"`
	Types []struct {
		Type struct {
			Name string `json:"name"`
		} `json:"type"`
	} `json:"types"`
}

func commandCatch(cfg *config, args ...string) error {
	if len(args) < 1 {
		return fmt.Errorf("usage: catch <pokemon>")
	}
	name := strings.ToLower(args[0])
	url := "https://pokeapi.co/api/v2/pokemon/" + name
	fmt.Printf("Throwing a Pokeball at %s...\n", name)

	var pokeData []byte
	if cfg.cache != nil {
		if data, ok := cfg.cache.Get(url); ok {
			pokeData = data
		}
	}
	if pokeData == nil {
		resp, err := http.Get(url)
		if err != nil {
			return fmt.Errorf("failed to fetch pokemon: %w", err)
		}
		defer resp.Body.Close()
		pokeData, err = io.ReadAll(resp.Body)
		if err != nil {
			return fmt.Errorf("failed to read response: %w", err)
		}
		if cfg.cache != nil {
			cfg.cache.Add(url, pokeData)
		}
	}

	var pkmn Pokemon
	if err := json.Unmarshal(pokeData, &pkmn); err != nil {
		return fmt.Errorf("failed to parse pokemon: %w", err)
	}

	// Lower base_experience = easier catch
	chance := 1.0 - float64(pkmn.BaseExperience)/400.0
	if chance < 0.1 {
		chance = 0.1
	}
	if randFloat64() < chance {
		if cfg.pokedex == nil {
			cfg.pokedex = make(map[string]Pokemon)
		}
		cfg.pokedex[pkmn.Name] = pkmn
		fmt.Printf("%s was caught!\n", name)
	} else {
		fmt.Printf("%s escaped!\n", name)
	}
	return nil
}
