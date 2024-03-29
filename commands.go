package main

import (
	"fmt"
	"os"

	"pokedexcli/pokeAPI"
)

type cliCommand struct {
	name        string
	description string
	callback    func(*pokeapi.Config) error
}

// func (self cliCommand) callbackv2() error {
// 	if self.name == "help" {
// 		return commandHelp()
// 	}
// 	if self.name == "exit" {
// 		return commandExit()
// 	}
// 	if self.name == "map" {
// 		if self.url.Next == 0 {
// 			self.url.Previous = 20
// 		}
// 		return commandMap(self.url.Next)
// 	}
// 	if self.name == "mapb" {
// 		if self.url.Previous == 0 {
// 			fmt.Println("No previous locations to display\n")
// 		}
// 		return commandMapBack(self.url.Previous)
// 	}
// 	return nil
// }

func commandHelp(_ *pokeapi.Config) error {
	fmt.Println("\nAvailable commands:\n")
	for _, command := range getCommands() {
		fmt.Printf("%v: %v\n", command.name, command.description)
	}
	fmt.Println()
	return nil
}

func commandExit(_ *pokeapi.Config) error {
	os.Exit(0)
	return nil
}

func commandMapForward(conf *pokeapi.Config) error {
	for i := conf.Next; i < conf.Next+20; i++ {
		url := fmt.Sprintf("https://pokeapi.co/api/v2/location-area/%d/", i+1)
		pokeapi.GetLocation(url)
	}
	conf.Previous = conf.Next - 20
	conf.Next += 20
	return nil
}

func commandMapBack(conf *pokeapi.Config) error {
	if conf.Next < 21 {
		fmt.Println("No previous locations to display\n")
		return nil
	}
	for i := conf.Previous; i < conf.Previous+20; i++ {
		url := fmt.Sprintf("https://pokeapi.co/api/v2/location-area/%d/", i+1)
		pokeapi.GetLocation(url)
	}
	conf.Next = conf.Previous + 20
	conf.Previous -= 20
	return nil
}

func getCommands() map[string]cliCommand {
	return map[string]cliCommand{
		"help": {
			name:        "help",
			description: "Displays a help message",
			callback:    func(c *pokeapi.Config) error { commandHelp(c); return nil },
		},
		"exit": {
			name:        "exit",
			description: "Exit the program",
			callback:    func(c *pokeapi.Config) error { commandExit(c); return nil },
		},
		"map": {
			name:        "map",
			description: "Display the next 20 location areas in the Pokemon world",
			callback:    func(c *pokeapi.Config) error { commandMapForward(c); return nil },
		},
		"mapb": {
			name:        "mapb",
			description: "Display the previous 20 location areas in the Pokemon world",
			callback:    func(c *pokeapi.Config) error { commandMapBack(c); return nil },
		},
	}
}
