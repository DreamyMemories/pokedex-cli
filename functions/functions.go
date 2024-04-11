package functions

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/DreamyMemories/pokedex-cli/pokecache"
	"github.com/DreamyMemories/pokedex-cli/types"
)

type cliCommand struct {
	name        string
	description string
	Callback    func(configPtr *Config, cache *pokecache.Cache, arg string) error
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
		"explore": {
			name:        "explore",
			description: "Explores a specific area and returning a list of pokemons, explore <area-name>",
			Callback:    commandExplore,
		},
	}
}

func GetNameAndArg(input string) (commandName string, argument string) {
	args := strings.Fields(input)
	command := args[0]
	var commandArg string
	if len(args) > 1 {
		commandArg = args[1]
	}

	return command, commandArg
}

func commandHelp(configPtr *Config, cache *pokecache.Cache, argument string) error {
	fmt.Println("Here are the available commands:")
	commands := GetCommands()
	for _, command := range commands {
		fmt.Println(command.name, ":", command.description)
	}
	fmt.Println("")
	return nil
}

func commandExit(configPtr *Config, cache *pokecache.Cache, argument string) error {
	fmt.Println("Goodbye!")
	os.Exit(0)
	return nil
}

func commandExplore(configPtr *Config, cache *pokecache.Cache, argument string) error {
	data, cached := cache.Get(argument)
	if cached {
		switch d := data.(type) {
		case types.EncounterApiResponse:
			displayItems(d)
		}
	} else {

		response, err := http.Get("https://pokeapi.co/api/v2/location-area/" + argument)

		if err != nil {
			return fmt.Errorf("HTTP request failed: %v", err)
		}

		body, _ := io.ReadAll(response.Body)
		if response.StatusCode != http.StatusOK {
			return fmt.Errorf("HTTP request unsuccessful: %v", response.StatusCode)
		}

		var apiResponse types.EncounterApiResponse

		if err := json.Unmarshal(body, &apiResponse); err != nil {
			return fmt.Errorf("Error in parsing json: %v", err)
		}
		cache.Add(argument, apiResponse)
		displayItems(apiResponse)
	}
	return nil
}

func displayItems(response interface{}) {
	switch r := response.(type) {
	case types.ApiResponse:
		for _, location := range r.Results {
			fmt.Println(location.Name)
		}
	case types.EncounterApiResponse:
		for _, pokemon := range r.PokemonEncounters {
			fmt.Println(pokemon.Pokemon.Name)
		}
	}

}

func commandMap(configPtr *Config, cache *pokecache.Cache, argument string) error {
	if configPtr.Next == "" {
		response, _ := http.Get("https://pokeapi.co/api/v2/location-area")
		if response.StatusCode != http.StatusOK {
			return fmt.Errorf("bad HTTP Status Code: %v", response.StatusCode)
		}
		body, err := io.ReadAll(response.Body)
		if err != nil {
			return fmt.Errorf("http request unsuccessful body: %v", err)
		}
		var apiResponse types.ApiResponse

		// Unmarshal json
		if err := json.Unmarshal([]byte(body), &apiResponse); err != nil {
			return fmt.Errorf("error in parsing json: %v", err)
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
			switch v := data.(type) {
			case types.ApiResponse:
				displayItems(v)
				configPtr.Previous = v.Previous
				configPtr.Next = v.Next
			}

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

func commandMapb(configPtr *Config, cache *pokecache.Cache, argument string) error {
	if configPtr.Previous == "" {
		fmt.Println("Error: no previous request, please use map first")
	} else {
		// Check cache
		data, cached := cache.Get(configPtr.Previous)
		if cached {
			switch v := data.(type) {
			case types.ApiResponse:
				displayItems(v)
				configPtr.Previous = v.Previous
				configPtr.Next = v.Next
			}
		} else {
			response, _ := http.Get(configPtr.Previous)

			if response.StatusCode != http.StatusOK {
				return fmt.Errorf("Bad HTTP Response Status Code: %v", response.StatusCode)
			}

			body, err := io.ReadAll(response.Body)
			if err != nil {
				return fmt.Errorf("Bad HTTP body: %v", body)
			}
			var apiResponse types.ApiResponse

			// Unmarshal json
			if err := json.Unmarshal([]byte(body), &apiResponse); err != nil {
				return fmt.Errorf("Error in parsing json: %v", err)
			}

			// Set next and back in config
			configPtr.Next = apiResponse.Next
			configPtr.Previous = apiResponse.Previous
			displayItems(apiResponse)
		}
	}
	return nil
}
