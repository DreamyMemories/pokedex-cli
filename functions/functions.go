package functions

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/DreamyMemories/pokedex-cli/pokecache"
	"github.com/DreamyMemories/pokedex-cli/types"
)

type cliCommand struct {
	name        string
	description string
	Callback    func(configPtr *Config, cache *pokecache.Cache) error
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

func commandHelp(configPtr *Config, cache *pokecache.Cache) error {
	fmt.Println("Here are the available commands:")
	commands := GetCommands()
	for _, command := range commands {
		fmt.Println(command.name, ":", command.description)
	}
	fmt.Println("")
	return nil
}

func commandExit(configPtr *Config, cache *pokecache.Cache) error {
	fmt.Println("Goodbye!")
	os.Exit(0)
	return nil
}

func displayItems(response types.ApiResponse) {
	for _, location := range response.Results {
		fmt.Println(location.Name)
	}
}

func commandMap(configPtr *Config, cache *pokecache.Cache) error {
	if configPtr.Next == "" {
		response, err := http.Get("https://pokeapi.co/api/v2/location-area")
		if err != nil {
			fmt.Println(err)
		}
		body, err := io.ReadAll(response.Body)
		if err != nil {
			fmt.Println(err)
		}
		var apiResponse types.ApiResponse

		// Unmarshal json
		if err := json.Unmarshal([]byte(body), &apiResponse); err != nil {
			fmt.Println("Error in parsing json: %v", err)
		}

		// Set next and back in config
		configPtr.Next = apiResponse.Next
		configPtr.Previous = "https://pokeapi.co/api/v2/location-area"
		fmt.Printf("Set cache %v", configPtr.Previous)
		cache.Add(configPtr.Previous, apiResponse)
		displayItems(apiResponse)
		return nil
	} else {
		// Get data from cache
		data, cached := cache.Get(configPtr.Next)
		if cached {
			displayItems(data)
			configPtr.Previous = data.Previous
			configPtr.Next = data.Next
		} else {
			response, err := http.Get(configPtr.Next)
			if err != nil {
				fmt.Println(err)
			}
			body, err := io.ReadAll(response.Body)
			if err != nil {
				fmt.Println(err)
			}
			var apiResponse types.ApiResponse

			// Unmarshal json
			if err := json.Unmarshal([]byte(body), &apiResponse); err != nil {
				fmt.Println("Error in parsing json: %v", err)
			}
			fmt.Printf("Set cache %v", configPtr.Next)
			cache.Add(configPtr.Next, apiResponse)
			// Set next and back in config
			configPtr.Previous = apiResponse.Previous // Updates with current
			configPtr.Next = apiResponse.Next         // Update with next
			displayItems(apiResponse)
		}
		return nil
	}
}

func commandMapb(configPtr *Config, cache *pokecache.Cache) error {
	if configPtr.Previous == "" {
		fmt.Println("Error: no previous request, please use map first")
	} else {
		// Check cache
		fmt.Printf("Getting cache %v", configPtr.Previous)
		data, cached := cache.Get(configPtr.Previous)
		if cached {
			displayItems(data)
			configPtr.Previous = data.Previous
			configPtr.Next = data.Next
		} else {
			response, err := http.Get(configPtr.Previous)

			if err != nil {
				return err
			}

			body, err := io.ReadAll(response.Body)
			if err != nil {
				fmt.Println(err)
			}
			var apiResponse types.ApiResponse

			// Unmarshal json
			if err := json.Unmarshal([]byte(body), &apiResponse); err != nil {
				fmt.Println("Error in parsing json: %v", err)
			}

			// Set next and back in config
			configPtr.Next = apiResponse.Next
			configPtr.Previous = apiResponse.Previous
			displayItems(apiResponse)
		}
	}
	return nil
}
