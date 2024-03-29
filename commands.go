package main

import (
	"fmt"
	"os"

	"pokedexcli/internal/pokeapi"
	"pokedexcli/internal/pokecache"
)

type cliCommand struct {
	name        string
	description string
	callback    func(*pokeapi.Config, *pokecache.Cache) error
}

func commandHelp() error {
	fmt.Println("\nAvailable commands:\n")
	for _, command := range getCommands() {
		fmt.Printf("%v: %v\n", command.name, command.description)
	}
	fmt.Println()
	return nil
}

func commandExit() error {
	os.Exit(0)
	return nil
}

func commandMapForward(conf *pokeapi.Config, cache *pokecache.Cache) error {
	for i := conf.Next; i < conf.Next+20; i++ {
		url := fmt.Sprintf("https://pokeapi.co/api/v2/location-area/%d/", i+1)
		pokeapi.GetLocation(url, cache)
	}
	fmt.Println()

	conf.Previous = conf.Next - 20
	conf.Next += 20
	return nil
}

func commandMapBack(conf *pokeapi.Config, cache *pokecache.Cache) error {
	if conf.Next < 21 {
		fmt.Println("No previous locations to display\n")
		return nil
	}
	for i := conf.Previous; i < conf.Previous+20; i++ {
		url := fmt.Sprintf("https://pokeapi.co/api/v2/location-area/%d/", i+1)
		pokeapi.GetLocation(url, cache)
	}
	fmt.Println()

	conf.Next = conf.Previous + 20
	conf.Previous -= 20
	return nil
}

func getCommands() map[string]cliCommand {
	return map[string]cliCommand{
		"help": {
			name:        "help",
			description: "Displays a help message",
			callback:    func(c *pokeapi.Config, ch *pokecache.Cache) error { commandHelp(); return nil },
		},
		"exit": {
			name:        "exit",
			description: "Exit the program",
			callback:    func(c *pokeapi.Config, ch *pokecache.Cache) error { commandExit(); return nil },
		},
		"map": {
			name:        "map",
			description: "Display the next 20 location areas in the Pokemon world",
			callback:    func(c *pokeapi.Config, ch *pokecache.Cache) error { commandMapForward(c, ch); return nil },
		},
		"mapb": {
			name:        "mapb",
			description: "Display the previous 20 location areas in the Pokemon world",
			callback:    func(c *pokeapi.Config, ch *pokecache.Cache) error { commandMapBack(c, ch); return nil },
		},
	}
}
