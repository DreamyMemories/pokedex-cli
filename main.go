package main

import (
	"bufio"
	"fmt"
	"os"
	"time"

	"github.com/DreamyMemories/pokedex-cli/functions"
	"github.com/DreamyMemories/pokedex-cli/pokecache"
)

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	commands := functions.GetCommands()
	config := functions.Config{
		Next:     "",
		Previous: "",
	}
	cache := *pokecache.NewCache(5 * time.Minute)
	fmt.Println("Welcome to the Pokedex!")
	fmt.Print("Pokedex > ")
	for scanner.Scan() {
		_, exist := commands[scanner.Text()] // Map property, second value checks
		if !exist {
			fmt.Println("Command not found. Type 'help' for a list of commands.")

		} else {
			commands[scanner.Text()].Callback(&config, &cache)
		}
		fmt.Print("Pokedex > ")
	}

	if err := scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, "reading standard input:", err)
	}
}
