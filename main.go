package main

import (
	"bufio"
	"fmt"
	"os"

	"github.com/DreamyMemories/pokedex-cli/functions"
)

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	commands := functions.GetCommands()
	fmt.Println("Welcome to the Pokedex!")
	fmt.Print("Pokedex > ")
	for scanner.Scan() {
		_, exist := commands[scanner.Text()] // Map property, second value checks
		if !exist {
			fmt.Println("Command not found. Type 'help' for a list of commands.")

		} else {
			commands[scanner.Text()].Callback()
		}
		fmt.Print("Pokedex > ")
	}

	if err := scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, "reading standard input:", err)
	}
}
