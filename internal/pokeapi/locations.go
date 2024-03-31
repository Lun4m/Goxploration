package pokeapi

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"

	"pokedexcli/internal/pokecache"
)

type Config struct {
	Next     int
	Previous int
}

type Location struct {
	ID                   int    `json:"id"`
	Name                 string `json:"name"`
	GameIndex            int    `json:"game_index"`
	EncounterMethodRates []struct {
		EncounterMethod struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"encounter_method"`
		VersionDetails []struct {
			Rate    int `json:"rate"`
			Version struct {
				Name string `json:"name"`
				URL  string `json:"url"`
			} `json:"version"`
		} `json:"version_details"`
	} `json:"encounter_method_rates"`
	Location struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	} `json:"location"`
	Names []struct {
		Name     string `json:"name"`
		Language struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"language"`
	} `json:"names"`
	PokemonEncounters []struct {
		Pokemon struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"pokemon"`
		VersionDetails []struct {
			Version struct {
				Name string `json:"name"`
				URL  string `json:"url"`
			} `json:"version"`
			MaxChance        int `json:"max_chance"`
			EncounterDetails []struct {
				MinLevel        int   `json:"min_level"`
				MaxLevel        int   `json:"max_level"`
				ConditionValues []any `json:"condition_values"`
				Chance          int   `json:"chance"`
				Method          struct {
					Name string `json:"name"`
					URL  string `json:"url"`
				} `json:"method"`
			} `json:"encounter_details"`
		} `json:"version_details"`
	} `json:"pokemon_encounters"`
}

func GetLocation(url string, cache *pokecache.Cache) (Location, error) {
	location := Location{}

	if val, ok := cache.Get(url); ok {
		err := json.Unmarshal(val, &location)
		if err != nil {
			fmt.Println(err)
			return location, err
		}
		fmt.Println(location.Name)
	} else {
		response, err := http.Get(url)
		if err != nil {
			fmt.Println(err)
			return location, err
		}

		if response.StatusCode > 299 {
			return location, errors.New("Response status code > 299")
		}

		body, err := io.ReadAll(response.Body)
		response.Body.Close()
		if err != nil {
			fmt.Println(err)
			return location, err
		}

		cache.Add(url, body)
		err = json.Unmarshal(body, &location)
		if err != nil {
			fmt.Println(err)
			return location, err
		}

	}
	return location, nil

}
