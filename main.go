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

	for {
		fmt.Print("pokedex > ")
		scanner.Scan()
		input := scanner.Text()

		input_string := ""
		split := strings.Split(input, " ")
		input = split[0]

		if len(split) == 2 {
			input_string = split[1]

		}

		if command, ok := commands[input]; ok {
			command.callback(input_string, &conf, cache)
		} else {
			fmt.Print("Invalid command. ")
			fmt.Println("Type 'help' to get a list of the available commands\n")
		}
	}
}
