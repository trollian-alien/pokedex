package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"math/rand/v2"
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
		callback:    encounter,
	},
	"catch": {
		name:        "catch",
		description: "attempts to catch a pokemon",
		callback:    catch,
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

// map and mapb command helper
func areaLocationPrinter(URL string, c *pokecache.Cache) error {
	jason, err := getJSON(URL, c)
	if err != nil {
		fmt.Printf("Can't get location areas! Error: %v\n", err)
		return err
	}

	var locations locationArea
	if err := json.Unmarshal(jason, &locations); err != nil {
		fmt.Printf("Can't display locations. Error: %v\n", err)
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

// struct to store pokemon name and pokemon data location
type pokemon struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}

//gloabl variable having the latest location's pokemons
var currentPokemons = make([]pokemon, 0)

//encounter command
func encounter(args []string, c *pokecache.Cache) error {
	if len(args) == 0 {
		fmt.Println("You need to specify a location area duh! Usage encounter <location-area>")
		return nil
	}
	
	location := args[0]
	locationURL := "https://pokeapi.co/api/v2/location-area/" + location + "/"
	jason, err := getJSON(locationURL, c)
	if err != nil {
		fmt.Printf("Error getting encounter data. Error: %v\n", err)
		return err
	} else if bytes.Equal(jason,[]byte("Not Found")) {
		fmt.Println("That's not a pokemon location!")
		return nil
	}

	var pokemons encounters
	if err := json.Unmarshal(jason, &pokemons); err != nil {
		fmt.Printf("Error getting encounter data. Error: %v\n", err)
		return err
	}

	fmt.Printf("Exploring %v...\n", args[0])
	fmt.Println("Found Pokemon:")
	for _, pokemon := range pokemons.PokemonEncounters {
		currentPokemons = append(currentPokemons, pokemon.Pokemon)
		fmt.Println(pokemon.Pokemon.Name)
	}
	return nil	
}

//global variable showing the caught pokemon
var caught = make([]pokemon, 0)

func catch(args []string, c *pokecache.Cache) error {
	if len(args) == 0 {
		fmt.Println("You caught nothing! Usage catch <pokemon-you-want-to-catch>")
		return nil
	}
	if len(currentPokemons) == 0 {
		fmt.Println("You need to use the explore command to find some pokemons first!")
		return nil
	}

	pokemonName := args[0]
	pokemonURL := ""
	for _, poke := range caught {
		if poke.Name == pokemonName {
		fmt.Printf("You already caught a %v!", pokemonName)
		return nil
	}
	}
	for _, poke := range currentPokemons {
		if poke.Name == pokemonName {
			pokemonURL = poke.URL
		}
	}
	if pokemonURL == "" {
		fmt.Println("This pokemon does not exist in this location!")
		fmt.Println("Try to catch a pokemon in the current area!")
		return nil
	}

	jason, err := getJSON(pokemonURL, c)
	if err!=nil {
		fmt.Printf("Error Getting Pokemon data! Error: %v\n", err)
		return err
	}

	var stats pokemonStat
	if err := json.Unmarshal(jason, &stats); err != nil {
		fmt.Printf("Error getting Pokemon data. Error: %v\n", err)
		return err
	}

	fmt.Printf("Throwing a Pokeball at %v", pokemonName)
	exp := stats.BaseExperience
	roll := rand.IntN(exp)
	if roll >= 35 {
		fmt.Printf("%v escaped!\n", pokemonName)
	} else {
		caught = append(caught, pokemon{pokemonName,pokemonURL})
		fmt.Printf("%v is caught!\n", pokemonName)
	}

	return nil
}