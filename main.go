package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print("Pokedex > ")
		if !scanner.Scan() {
			break
		}
		input := scanner.Text()
		words := cleanInput(input)
		if len(words) > 0 {
			fmt.Printf("Your command was: %s\n", strings.ToLower(words[0]))
		}
	}
}

func cleanInput(text string) []string {
	fields := strings.Fields(strings.TrimSpace(text))
	return fields
}
