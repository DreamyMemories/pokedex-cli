# Pokedex-CLI

Welcome to Pokedex-CLI, a command-line interface for exploring the vast world of Pokémon! Written in Go, this tool allows users to explore Pokémon data by querying different locations within the Pokémon universe.

## Features

- `explore <location>`: Retrieve a list of Pokémon found in the specified location.
- `map`: Display a list of available locations to explore.
- `mapb`: Display a list of previous location shown by `map`.

## Installation

To install Pokedex-CLI, you need to have Go installed on your system. Follow these steps:

```bash
git clone https://github.com/yourusername/pokedex-cli.git
cd pokedex-cli
go build -o pokedex-cli
./pokedex-cli.exe
```

## Usage
After installation, you can perform the following commands:

- To explore a location and find Pokémon, use the `explore` command followed by the location name. For example:
```bash
explore pallet-town
```

- To view a list of available locations, use the `map` command:
```bash
map
```

## File Structure
- `main.go`: Contains the main logic for the CLI application.
- `functions.go` : Contains the functions for API Calls and data processing.
- `pokecache/`: Caching system for API data
    - `pokecache.go`: Contains the cache logic.
    - `pokecache_test.go`: Tests for the cache logic.
- `types/`: Contains the type definitions for the application.
    - `types.go`: Struct representing pokemon and location.

    