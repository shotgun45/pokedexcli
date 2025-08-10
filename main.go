package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func main() {
	commandRegistry := getCommands()
	scanner := bufio.NewScanner(os.Stdin)
	cfg := &config{}
	for {
		fmt.Print("Pokedex > ")
		if !scanner.Scan() {
			break
		}
		input := scanner.Text()
		words := cleanInput(input)
		if len(words) > 0 {
			cmdName := strings.ToLower(words[0])
			if cmd, ok := commandRegistry[cmdName]; ok {
				err := cmd.callback(cfg, words[1:]...)
				if err != nil {
					fmt.Printf("Error: %v\n", err)
				}
			} else {
				fmt.Println("Unknown command:", cmdName)
			}
		}
	}
}

func cleanInput(text string) []string {
	fields := strings.Fields(strings.TrimSpace(text))
	return fields
}
