package functions

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
)

type cliCommand struct {
	name        string
	description string
	Callback    func(configPtr *Config) error
}

type LocationArea struct {
	Name string `json:"name"`
	Url  string `json:"url"`
}

type ApiResponse struct {
	Next     string         `json:"next"`
	Previous string         `json:"previous"`
	Results  []LocationArea `json:"results"`
}

type Config struct {
	Next     string
	Previous string
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
		"map": {
			name:        "map",
			description: "Displays next 20 locations",
			Callback:    commandMap,
		},
		"mapb": {
			name:        "mapb",
			description: "Displays next 20 locations back",
			Callback:    commandMapb,
		},
	}
}

func commandHelp(configPtr *Config) error {
	fmt.Println("Here are the available commands:")
	commands := GetCommands()
	for _, command := range commands {
		fmt.Println(command.name, ":", command.description)
	}
	fmt.Println("")
	return nil
}

func commandExit(configPtr *Config) error {
	fmt.Println("Goodbye!")
	os.Exit(0)
	return nil
}

func commandMap(configPtr *Config) error {
	if configPtr.Next == "" {
		response, err := http.Get("https://pokeapi.co/api/v2/location-area")
		if err != nil {
			fmt.Println(err)
		}
		body, err := io.ReadAll(response.Body)
		if err != nil {
			fmt.Println(err)
		}
		var apiResponse ApiResponse

		// Unmarshal json
		if err := json.Unmarshal([]byte(body), &apiResponse); err != nil {
			fmt.Println("Error in parsing json: %v", err)
		}

		// Set next and back in config
		configPtr.Next = apiResponse.Next
		configPtr.Previous = "https://pokeapi.co/api/v2/location-area"

		for _, location := range apiResponse.Results {
			fmt.Println(location.Name)
		}
		return nil
	} else {
		response, err := http.Get(configPtr.Next)
		if err != nil {
			fmt.Println(err)
		}
		body, err := io.ReadAll(response.Body)
		if err != nil {
			fmt.Println(err)
		}
		var apiResponse ApiResponse

		// Unmarshal json
		if err := json.Unmarshal([]byte(body), &apiResponse); err != nil {
			fmt.Println("Error in parsing json: %v", err)
		}

		// Set next and back in config
		configPtr.Next = apiResponse.Next
		configPtr.Previous = apiResponse.Previous

		for _, location := range apiResponse.Results {
			fmt.Println(location.Name)
		}
		return nil

	}
}

func commandMapb(configPtr *Config) error {
	if configPtr.Previous == "" {
		fmt.Println("Error: no previous request, please use map first")
	} else {
		response, err := http.Get(configPtr.Previous)

		if err != nil {
			return err
		}

		body, err := io.ReadAll(response.Body)
		if err != nil {
			fmt.Println(err)
		}
		var apiResponse ApiResponse

		// Unmarshal json
		if err := json.Unmarshal([]byte(body), &apiResponse); err != nil {
			fmt.Println("Error in parsing json: %v", err)
		}

		// Set next and back in config
		configPtr.Next = apiResponse.Next
		configPtr.Previous = apiResponse.Previous

		for _, location := range apiResponse.Results {
			fmt.Println(location.Name)
		}
		return nil

	}
	return nil
}
