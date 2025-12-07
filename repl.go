package main

import (
	"fmt"
	"net/http"
	"os"
	"strings"
	"encoding/json"
)

func cleanInput(text string) []string {
	return strings.Fields(strings.ToLower(text))
}

type cliCommand struct {
	name        string
	description string
	callback    func() error
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

func commandExit() error {
	fmt.Println("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return nil
}

func help() error {
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

func areaLocationPrinter(URL string) error {
	res, err := http.Get(URL)
	if err != nil {
		fmt.Printf("Can't display locations. Error: %v", err)
		return err
	}
	defer res.Body.Close()
	decoder := json.NewDecoder(res.Body)
	var locations locationArea
	if err := decoder.Decode(&locations); err != nil {
		fmt.Printf("Can't display locations. Error: %v", err)
		return err
	}
	for _, place := range(locations.Results) {
		fmt.Println(place.Name)
	}
	nextMapURL = locations.Next
	previousMapURL = locations.Previous
	return nil
}

func nextLocationAreas() error {
	URL := nextMapURL
	if URL == "" {
		fmt.Println("Wow you exhausted all the location areas!")
		return fmt.Errorf("no new area")
	}
	return areaLocationPrinter(URL)
}

func previousLocationAreas() error {
	URL := previousMapURL
	if URL == "" {
		fmt.Println("No previous areas!")
		return fmt.Errorf("no previous area")
	}
	return areaLocationPrinter(URL)
}