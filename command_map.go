package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"pokedexcli/internal/pokecache"
)

type locationAreaResponse struct {
	Results []struct {
		Name string `json:"name"`
	} `json:"results"`
	Next     *string `json:"next"`
	Previous *string `json:"previous"`
}

func getOrFetchURL(cache *pokecache.Cache, url string) ([]byte, error) {
	if data, ok := cache.Get(url); ok {
		return data, nil
	}
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	cache.Add(url, body)
	return body, nil
}

func commandMap(cfg *config, args ...string) error {
	url := "https://pokeapi.co/api/v2/location-area?limit=20"
	if cfg.nextLocationsURL != nil && *cfg.nextLocationsURL != "" {
		url = *cfg.nextLocationsURL
	}
	if cfg.cache == nil {
		return fmt.Errorf("cache not initialized")
	}
	cache, ok := cfg.cache.(*pokecache.Cache)
	if !ok {
		return fmt.Errorf("cache is not of type *pokecache.Cache")
	}
	body, err := getOrFetchURL(cache, url)
	if err != nil {
		return fmt.Errorf("failed to fetch location areas: %w", err)
	}
	var data locationAreaResponse
	if err := json.Unmarshal(body, &data); err != nil {
		return fmt.Errorf("failed to parse response: %w", err)
	}
	for _, area := range data.Results {
		fmt.Println(area.Name)
	}
	cfg.nextLocationsURL = data.Next
	cfg.prevLocationsURL = data.Previous
	return nil
}

func commandMapb(cfg *config, args ...string) error {
	if cfg.prevLocationsURL == nil || *cfg.prevLocationsURL == "" {
		fmt.Println("you're on the first page")
		return nil
	}
	if cfg.cache == nil {
		return fmt.Errorf("cache not initialized")
	}
	url := *cfg.prevLocationsURL
	cache, ok := cfg.cache.(*pokecache.Cache)
	if !ok {
		return fmt.Errorf("cache is not of type *pokecache.Cache")
	}
	body, err := getOrFetchURL(cache, url)
	if err != nil {
		return fmt.Errorf("failed to fetch location areas: %w", err)
	}
	var data locationAreaResponse
	if err := json.Unmarshal(body, &data); err != nil {
		return fmt.Errorf("failed to parse response: %w", err)
	}
	for _, area := range data.Results {
		fmt.Println(area.Name)
	}
	cfg.nextLocationsURL = data.Next
	cfg.prevLocationsURL = data.Previous
	return nil
}
