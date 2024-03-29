package main

import (
	"bufio"
	"fmt"
	"os"
	"pokedexcli/pokeAPI"
)

func main() {
	fmt.Println("Welcome to Pokedex CLI!\n")
	scanner := bufio.NewScanner(os.Stdin)
	conf := pokeapi.Config{}
	commands := getCommands()

	for {
		fmt.Print("pokedex > ")
		scanner.Scan()
		input := scanner.Text()
		if command, ok := commands[input]; ok {
			command.callback(&conf)
		} else {
			fmt.Print("Invalid command. ")
			fmt.Println("Type 'help' to get a list of the available commands\n")
		}
	}
}
