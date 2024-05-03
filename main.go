package main

import (
	"bufio"
	"fmt"
	"os"
	"time"

	"github.com/DreamyMemories/pokedex-cli/functions"
	"github.com/DreamyMemories/pokedex-cli/pokecache"
	"github.com/DreamyMemories/pokedex-cli/types"
)

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	commands := functions.GetCommands()
	config := functions.Config{
		Next:     "",
		Previous: "",
	}
	cache := *pokecache.NewCache(5 * time.Minute)
	pokemonTeam := make(map[string]types.Pokemon)
	fmt.Println("Welcome to the Pokedex!")
	fmt.Print("Pokedex > ")
	for scanner.Scan() {
		command, arg := functions.GetNameAndArg(scanner.Text())
		_, exist := commands[command] // Map property, second value checks
		if !exist {
			fmt.Println("Command not found. Type 'help' for a list of commands.")

		} else {
			commands[command].Callback(&config, &cache, arg, pokemonTeam)
		}
		fmt.Print("Pokedex > ")
	}

	if err := scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, "reading standard input:", err)
	}
}
