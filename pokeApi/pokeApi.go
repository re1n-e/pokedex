package pokeapi

import (
	"encoding/json"
	"fmt"
	"net/http"
)

const LocationApi = "https://pokeapi.co/api/v2/location-area/"

type Location struct {
	Count    int       `json:"count"`
	Next     *string   `json:"next"`
	Previous *string   `json:"previous"`
	Results  []Results `json:"results"`
}

type Results struct {
	Name string `json:"name"`
	Url  string `json:"url"`
}

func LocationNew() Location {
	return Location{}
}

func (r *Location) fetchData(apiUrl string) error {
	resp, err := http.Get(apiUrl)
	if err != nil {
		return fmt.Errorf("-Unable to fetch request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode > 299 {
		return fmt.Errorf("-Bad status code: %v", resp.StatusCode)
	}

	if err = json.NewDecoder(resp.Body).Decode(r); err != nil {
		return fmt.Errorf("-Unable to decode JSON to struct: %v", err)
	}

	return nil
}

func (r *Location) getNext() error {
	if r.Next == nil {
		return r.fetchData(LocationApi)
	}
	return r.fetchData(*r.Next)
}

func (r *Location) GetNextCities() ([]string, error) {
	if err := r.getNext(); err != nil {
		return nil, err
	}

	var results []string
	for _, result := range r.Results {
		results = append(results, result.Name)
	}
	return results, nil
}

func (r *Location) getPrevious() error {
	if r.Previous == nil {
		return fmt.Errorf("-There are no previous cities")
	}
	return r.fetchData(*r.Previous)
}

func (r *Location) GetPreviousCities() ([]string, error) {
	if err := r.getPrevious(); err != nil {
		return nil, err
	}

	var results []string
	for _, result := range r.Results {
		results = append(results, result.Name)
	}
	return results, nil
}
