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
	callback    func(*pokecache.Cache) error
}

var commands = map[string]cliCommand{
	"exit": {
		name:        "exit",
		description: "Exit the Pokedex",
		callback:    commandExit,
	},
	"help": {
		name:        "help",
		description: "Helps you know some other commands",
		callback:    help,
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
}

// exit command
func commandExit(c *pokecache.Cache) error {
	fmt.Println("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return nil
}

// help command
func help(c *pokecache.Cache) error {
	fmt.Print(`Welcome to the Pokedex! Here are some useful commands:\n
				exit: exits the Pokedex\n
				map: displays the next 20 location areas\n
				mapb: displays the previous 20 location areas\n`)
	return nil
}

var nextMapURL = "https://pokeapi.co/api/v2/location-area"
var previousMapURL = ""

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
func nextLocationAreas(c *pokecache.Cache) error {
	URL := nextMapURL
	if URL == "" {
		fmt.Println("Wow you exhausted all the location areas!")
		return fmt.Errorf("no new area")
	}
	return areaLocationPrinter(URL, c)
}

// mapb command
func previousLocationAreas(c *pokecache.Cache) error {
	URL := previousMapURL
	if URL == "" {
		fmt.Println("No previous areas!")
		return fmt.Errorf("no previous area")
	}
	return areaLocationPrinter(URL, c)
}
