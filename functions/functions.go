package functions

import (
	"fmt"
	"os"
)

type cliCommand struct {
	name        string
	description string
	Callback    func() error
}

func GetCommands() map[string]cliCommand {
	return map[string]cliCommand{
		"help": {
			name:        "help",
			description: "Displays a help message",
			Callback:    commandHelp,
		},
		"exit": {
			name:        "exit",
			description: "Exits the pokedex",
			Callback:    commandExit,
		},
	}
}

func commandHelp() error {
	fmt.Println("Here are the available commands:")
	commands := GetCommands()
	for _, command := range commands {
		fmt.Println(command.name, ":", command.description)
	}
	fmt.Println("")
	return nil
}

func commandExit() error {
	fmt.Println("Goodbye!")
	os.Exit(0)
	return nil
}
