package main

import (
	"bufio"
	"fmt"
	"os"
	"time"
	"github.com/trollian-alien/pokedex/internal/pokecache"
)

func main() {
	c := pokecache.NewCache(60*time.Second)
	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print("Pokedex > ")
		scanner.Scan()
		input := cleanInput(scanner.Text())
		if len(input) == 0 {
			continue
		}
		commandName := input[0]
		command, ok := commands[commandName]
		if !ok {
			fmt.Println("Unknown command")
		} else {
			command.callback(c)
		}
	}
}
