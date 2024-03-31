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
	callback    func(string, *pokeapi.Config, *pokecache.Cache) error
}

func commandHelp() {
	fmt.Println("\nAvailable commands:\n")
	for _, command := range getCommands() {
		fmt.Printf("%v: %v\n", command.name, command.description)
	}
	fmt.Println()
}

func commandExit() {
	os.Exit(0)
}

func commandMapForward(conf *pokeapi.Config, cache *pokecache.Cache) {
	for i := conf.Next; i < conf.Next+20; i++ {
		url := fmt.Sprintf("https://pokeapi.co/api/v2/location-area/%d/", i+1)
		location, err := pokeapi.GetLocation(url, cache)
		if err != nil {
			return
		}
		fmt.Println(location.Name)
	}
	fmt.Println()

	conf.Previous = conf.Next - 20
	conf.Next += 20
}

func commandMapBack(conf *pokeapi.Config, cache *pokecache.Cache) {
	if conf.Next < 21 {
		fmt.Println("No previous locations to display\n")
		return
	}
	for i := conf.Previous; i < conf.Previous+20; i++ {
		url := fmt.Sprintf("https://pokeapi.co/api/v2/location-area/%d/", i+1)
		pokeapi.GetLocation(url, cache)
		location, err := pokeapi.GetLocation(url, cache)
		if err != nil {
			return
		}
		fmt.Println(location.Name)
	}
	fmt.Println()

	conf.Next = conf.Previous + 20
	conf.Previous -= 20
}

func commandExplore(input string, cache *pokecache.Cache) {
	url := fmt.Sprintf("https://pokeapi.co/api/v2/location-area/%s/", input)
	location, err := pokeapi.GetLocation(url, cache)
	if err != nil {
		fmt.Println("Unknown area\n")
		return
	}

	fmt.Println("Found these Pokemon:")
	for _, pokemon := range location.PokemonEncounters {
		fmt.Printf(" - %s\n", pokemon.Pokemon.Name)
	}
	fmt.Println()
}

func getCommands() map[string]cliCommand {
	return map[string]cliCommand{
		"help": {
			name:        "help",
			description: "Displays a help message",
			callback: func(s string, c *pokeapi.Config, ch *pokecache.Cache) error {
				commandHelp()
				return nil
			},
		},
		"exit": {
			name:        "exit",
			description: "Exit the program",
			callback: func(s string, c *pokeapi.Config, ch *pokecache.Cache) error {
				commandExit()
				return nil
			},
		},
		"map": {
			name:        "map",
			description: "Display the next 20 location areas in the Pokemon world",
			callback: func(s string, c *pokeapi.Config, ch *pokecache.Cache) error {
				commandMapForward(c, ch)
				return nil
			},
		},
		"mapb": {
			name:        "mapb",
			description: "Display the previous 20 location areas in the Pokemon world",
			callback: func(s string, c *pokeapi.Config, ch *pokecache.Cache) error {
				commandMapBack(c, ch)
				return nil
			},
		},
		"explore": {
			name:        "explore",
			description: "Display the names of the pokemon present in the region",
			callback: func(s string, c *pokeapi.Config, ch *pokecache.Cache) error {
				commandExplore(s, ch)
				return nil
			},
		},
	}
}
