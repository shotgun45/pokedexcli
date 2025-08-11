package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// CacheInterface allows mocking for tests
type CacheInterface interface {
	Get(key string) ([]byte, bool)
	Add(key string, val []byte)
}

// Struct for the relevant part of the location-area response
// Only the pokemon_encounters field is needed for this command
type locationAreaDetail struct {
	PokemonEncounters []struct {
		Pokemon struct {
			Name string `json:"name"`
		} `json:"pokemon"`
	} `json:"pokemon_encounters"`
}

func commandExplore(cfg *config, args ...string) error {
	if len(args) < 1 {
		return fmt.Errorf("usage: explore <location-area>")
	}
	area := args[0]
	url := "https://pokeapi.co/api/v2/location-area/" + area
	var cache CacheInterface = cfg.cache
	if cache == nil {
		return fmt.Errorf("cache not initialized")
	}
	fmt.Printf("Exploring %s...\n", area)
	body, err := getOrFetchURLExplore(cache, url)
	if err != nil {
		return fmt.Errorf("failed to fetch location area: %w", err)
	}
	var data locationAreaDetail
	if err := json.Unmarshal(body, &data); err != nil {
		return fmt.Errorf("failed to parse response: %w", err)
	}
	if len(data.PokemonEncounters) == 0 {
		fmt.Println("No Pokemon found in this area.")
		return nil
	}
	fmt.Println("Found Pokemon:")
	for _, encounter := range data.PokemonEncounters {
		fmt.Printf(" - %s\n", encounter.Pokemon.Name)
	}
	return nil
}

// getOrFetchURLExplore for explore command, using CacheInterface
func getOrFetchURLExplore(cache CacheInterface, url string) ([]byte, error) {
	if data, ok := cache.Get(url); ok {
		return data, nil
	}
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("API returned status: %s", resp.Status)
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	cache.Add(url, body)
	return body, nil
}
