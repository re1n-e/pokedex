package pokeapi

import (
	"encoding/json"
	"fmt"
	"net/http"
)

const PokeApi = "https://pokeapi.co/api/v2/"

type Fetcher interface {
	fetchData(apiUrl string) error
}

type Location struct {
	Count    int       `json:"count"`
	Next     *string   `json:"next"`
	Previous *string   `json:"previous"`
	Results  []Results `json:"results"`
}

type Results struct {
	Name string `json:"name"`
}

type Encounter struct {
	PokemonEncounter []PokemonEncounters `json:"pokemon_encounters"`
}

type PokemonEncounters struct {
	Pokemon Pokemon `json:"pokemon"`
}

type Pokemon struct {
	Name string `json:"name"`
}

func LocationNew() Location {
	return Location{}
}

func EncounterNew() Encounter {
	return Encounter{}
}

func fetchJSON(apiUrl string, target interface{}) error {
	resp, err := http.Get(apiUrl)
	if err != nil {
		return fmt.Errorf("-Unable to fetch request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode > 299 {
		return fmt.Errorf("-Bad status code: %v", resp.StatusCode)
	}

	if err := json.NewDecoder(resp.Body).Decode(target); err != nil {
		return fmt.Errorf("-Unable to decode JSON to struct: %v", err)
	}

	return nil
}

func (r *Location) fetchData(apiUrl string) error {
	return fetchJSON(apiUrl, r)
}

func (e *Encounter) fetchData(apiUrl string) error {
	return fetchJSON(apiUrl, e)
}

func (e *Encounter) GetPokemons(apiUrl string) error {
	return e.fetchData(apiUrl)
}

func (r *Location) GetNext(nextUrl string) error {
	if nextUrl == "" {
		return fmt.Errorf("")
	}
	return r.fetchData(nextUrl)
}

func (r *Location) GetPrevious(prevUrl string) error {
	if prevUrl == "" {
		return fmt.Errorf("")
	}
	return r.fetchData(prevUrl)
}
