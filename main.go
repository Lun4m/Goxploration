package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"time"

	"pokedexcli/internal/pokeapi"
	"pokedexcli/internal/pokecache"
)

func main() {
	fmt.Println("Welcome to Pokedex CLI!\n")
	scanner := bufio.NewScanner(os.Stdin)
	conf := pokeapi.Config{}
	commands := getCommands()
	cache := pokecache.NewCache(1 * time.Minute)
	pokedex := make(map[string]pokeapi.Pokemon)

	for {
		fmt.Print("pokedex > ")
		scanner.Scan()
		input := scanner.Text()

		split := strings.Fields(input)
		input = split[0]

		if command, ok := commands[input]; ok {
			if command.isValid(split) {
				command.callback(split, pokedex, &conf, cache)
			}
		} else {
			fmt.Print("Invalid command. ")
			fmt.Println("Type 'help' to get a list of the available commands\n")
		}
	}
}
