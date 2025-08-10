package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type locationAreaResponse struct {
	Results []struct {
		Name string `json:"name"`
	} `json:"results"`
	Next     *string `json:"next"`
	Previous *string `json:"previous"`
}

func commandMap(cfg *config, args ...string) error {
	url := "https://pokeapi.co/api/v2/location-area?limit=20"
	if cfg.nextLocationsURL != nil && *cfg.nextLocationsURL != "" {
		url = *cfg.nextLocationsURL
	}
	resp, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("failed to fetch location areas: %w", err)
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read response: %w", err)
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
	url := *cfg.prevLocationsURL
	resp, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("failed to fetch location areas: %w", err)
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read response: %w", err)
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
