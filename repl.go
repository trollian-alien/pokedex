package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/trollian-alien/pokedex/internal/pokecache"
)

func cleanInput(text string) []string {
	return strings.Fields(strings.ToLower(text))
}

type cliCommand struct {
	name        string
	description string
	callback    func([]string, *pokecache.Cache) error
}

var commands = map[string]cliCommand{
	"exit": {
		name:        "exit",
		description: "Exit the Pokedex",
		callback:    commandExit,
	},
	"help": { 
		name:        "help",
		description: "Why would you ask help for the help command?",
		callback:    nil, //using help instead of nil here causes circular dependency (bad). implent seprate logic for help lol
	},
	"map": {
		name:        "map",
		description: "displays the next 20 locations areas",
		callback:    nextLocationAreas,
	},
	"mapb": {
		name:        "map back",
		description: "displays the previous 20 locations areas",
		callback:    previousLocationAreas,
	},
	"explore": {
		name:        "explore",
		description: "explores the stopulated area",
		callback:    previousLocationAreas,
	},
}

// exit command
func commandExit(args []string, c *pokecache.Cache) error {
	fmt.Println("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return nil
}

// help command; handled seprately so no need for a cache
func help(args []string) error {
	if len(args) == 0 {
		fmt.Println(`Welcome to the Pokedex! Here is a list of the available commands:`)
		for key := range commands {
			fmt.Println(key)
		}
		fmt.Println("For more details about each command, type help <command-name>")
	} else {
		arg := args[0]
		command, ok := commands[arg]
		if !ok {
			fmt.Println("This command does nothing...")
			fmt.Println("because it doesn't exist!")
		} else {
			fmt.Println("Name: " + command.name)
			fmt.Println("Description: " + command.description)
		}
	}
	return nil
}

//global variables for map and mapb
var nextMapURL = "https://pokeapi.co/api/v2/location-area"
var previousMapURL = ""

//json intepreter struct for map and mapb
type locationArea struct {
	Count    int    `json:"count"`
	Next     string `json:"next"`
	Previous string `json:"previous"`
	Results  []struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	} `json:"results"`
}

// map and mapb command helper
func areaLocationPrinter(URL string, c *pokecache.Cache) error {
	//first check if it's in the cache
	jason, ok := c.Get(URL)
	if !ok { //not in the cache, make a new HTTP request
		res, err := http.Get(URL)
		if err != nil {
			fmt.Printf("Can't display locations. Error: %v", err)
			return err
		}
		defer res.Body.Close()
		// read the JSON
		jason, err = io.ReadAll(res.Body)
		if err != nil {
			fmt.Printf("Can't display locations. Error: %v", err)
			return err
		}
		c.Add(URL, jason)
	}

	var locations locationArea
	if err := json.Unmarshal(jason, &locations); err != nil {
		fmt.Printf("Can't display locations. Error: %v", err)
		return err
	}

	for _, place := range locations.Results {
		fmt.Println(place.Name)
	}
	nextMapURL = locations.Next
	previousMapURL = locations.Previous
	return nil
}

// map command
func nextLocationAreas(args []string, c *pokecache.Cache) error {
	URL := nextMapURL
	if URL == "" {
		fmt.Println("Wow you exhausted all the location areas!")
		return fmt.Errorf("no new area")
	}
	return areaLocationPrinter(URL, c)
}

// mapb command
func previousLocationAreas(args []string, c *pokecache.Cache) error {
	URL := previousMapURL
	if URL == "" {
		fmt.Println("No previous areas!")
		return fmt.Errorf("no previous area")
	}
	return areaLocationPrinter(URL, c)
}

// json interpreter struct for explore command
type encounters struct {
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
	GameIndex int `json:"game_index"`
	ID        int `json:"id"`
	Location  struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	} `json:"location"`
	Name  string `json:"name"`
	Names []struct {
		Language struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"language"`
		Name string `json:"name"`
	} `json:"names"`
	PokemonEncounters []struct {
		Pokemon struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"pokemon"`
		VersionDetails []struct {
			EncounterDetails []struct {
				Chance          int   `json:"chance"`
				ConditionValues []any `json:"condition_values"`
				MaxLevel        int   `json:"max_level"`
				Method          struct {
					Name string `json:"name"`
					URL  string `json:"url"`
				} `json:"method"`
				MinLevel int `json:"min_level"`
			} `json:"encounter_details"`
			MaxChance int `json:"max_chance"`
			Version   struct {
				Name string `json:"name"`
				URL  string `json:"url"`
			} `json:"version"`
		} `json:"version_details"`
	} `json:"pokemon_encounters"`
}
