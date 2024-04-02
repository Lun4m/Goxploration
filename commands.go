package main

import (
	"fmt"
	"math/rand"
	"os"
	"time"

	"pokedexcli/internal/pokeapi"
	"pokedexcli/internal/pokecache"
)

type cliCommand struct {
	name           string
	description    string
	callback       func([]string, map[string]pokeapi.Pokemon, *pokeapi.Config, *pokecache.Cache) error
	requires_input bool
}

func (self *cliCommand) isValid(input []string) bool {
	if !self.requires_input && len(input) > 1 {
		fmt.Printf("'%s' doesn't accept arguments\n\n", self.name)
		return false
	}
	if self.requires_input && len(input) == 1 {
		fmt.Printf("'%s' requires an argument\n\n", self.name)
		return false
	}
	if self.requires_input && len(input) > 2 {
		fmt.Printf("'%s' only accepts a single argument\n\n", self.name)
		return false
	}
	return true
}

func commandHelp() {
	fmt.Println("\nAvailable commands:\n")
	keys := [6]string{"help", "exit", "map", "mapb", "explore", "catch"}

	commands := getCommands()
	for _, key := range keys {
		command := commands[key]
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

func getCatchProbability(baseExp int) float64 {
	if baseExp > 361 {
		return 0.001
	}
	x := float64(baseExp-36) / float64(360-36)
	// Simple interpolation
	return 0.05*x + (1.0-x)*0.75
}

func commandCatch(input string, pokedex map[string]pokeapi.Pokemon, cache *pokecache.Cache) {
	if _, ok := pokedex[input]; ok {
		fmt.Printf("%v is alredy in your pokedex\n\n", input)
		return
	}

	url := fmt.Sprintf("https://pokeapi.co/api/v2/pokemon/%s/", input)
	pokemon, err := pokeapi.GetPokemon(url, cache)
	if err != nil {
		fmt.Println("Unknown pokemon\n")
		return
	}
	fmt.Printf("Throwing a pokeball at %s...\n", input)
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	catchProbability := getCatchProbability(pokemon.BaseExperience)

	fmt.Printf("Probability of catching: %v\n", catchProbability)

	if rng.Float64() < float64(catchProbability) {
		fmt.Printf("%v was caught!\n", input)
		pokedex[input] = pokemon
	} else {
		fmt.Printf("%v escaped!\n", input)
	}
	fmt.Println()
}

func getCommands() map[string]cliCommand {
	return map[string]cliCommand{
		"help": {
			name:           "help",
			description:    "Displays a help message",
			requires_input: false,
			callback: func(s []string, p map[string]pokeapi.Pokemon, c *pokeapi.Config, ch *pokecache.Cache) error {
				commandHelp()
				return nil
			},
		},
		"exit": {
			name:           "exit",
			description:    "Exit the program",
			requires_input: false,
			callback: func(s []string, p map[string]pokeapi.Pokemon, c *pokeapi.Config, ch *pokecache.Cache) error {
				commandExit()
				return nil
			},
		},
		"map": {
			name:           "map",
			description:    "Display the next 20 location areas in the Pokemon world",
			requires_input: false,
			callback: func(s []string, p map[string]pokeapi.Pokemon, c *pokeapi.Config, ch *pokecache.Cache) error {
				commandMapForward(c, ch)
				return nil
			},
		},
		"mapb": {
			name:           "mapb",
			description:    "Display the previous 20 location areas in the Pokemon world",
			requires_input: false,
			callback: func(s []string, p map[string]pokeapi.Pokemon, c *pokeapi.Config, ch *pokecache.Cache) error {
				commandMapBack(c, ch)
				return nil
			},
		},
		"explore": {
			name:           "explore",
			description:    "Display the names of the pokemon present in the region",
			requires_input: true,
			callback: func(s []string, p map[string]pokeapi.Pokemon, c *pokeapi.Config, ch *pokecache.Cache) error {
				arg := s[1]
				commandExplore(arg, ch)
				return nil
			},
		},
		"catch": {
			name:           "catch",
			description:    "Try to catch a pokemon and add it to your pokedex",
			requires_input: true,
			callback: func(s []string, p map[string]pokeapi.Pokemon, c *pokeapi.Config, ch *pokecache.Cache) error {
				arg := s[1]
				commandCatch(arg, p, ch)
				return nil
			},
		},
	}
}
