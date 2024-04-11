package types

type ApiResponse struct {
	Next     string         `json:"next"`
	Previous string         `json:"previous"`
	Results  []LocationArea `json:"results"`
}

type LocationArea struct {
	Name string `json:"name"`
	Url  string `json:"url"`
}

type EncounterApiResponse struct {
	PokemonEncounters []PokemonEncounter `json:"pokemon_encounters"`
}

type PokemonEncounter struct {
	Pokemon PokemonInfo `json:"pokemon"`
}

type PokemonInfo struct {
	Name string `json:"name"`
	Url  string `json:"url"`
}
