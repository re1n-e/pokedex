package pokemon

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"os"
)

const PokemonApi = "https://pokeapi.co/api/v2/pokemon/"

type Pokemon struct {
	Name      string         `json:"name"`
	Abilities []AbilityEntry `json:"abilities"`
	BaseExp   int            `json:"base_experience"`
	Height    int            `json:"height"`
	Weight    int            `json:"weight"`
	Moves     []MoveEntry    `json:"moves"`
	Stats     []Stats        `json:"stats"`
	Type      []Type         `json:"types"`
}

type AbilityEntry struct {
	Ability Ability `json:"ability"`
}

type Ability struct {
	Name string `json:"name"`
}

type MoveEntry struct {
	Move Move `josn:"move"`
}

type Move struct {
	Name string `json:"name"`
}

type Stats struct {
	Base int      `json:"base_stat"`
	Name StatName `json:"stat"`
}

type StatName struct {
	Name string `json:"name"`
}

type Type struct {
	TypeName TypeName `json:"type"`
}

type TypeName struct {
	Name string `json:"name"`
}

func NewPokemon() Pokemon {
	return Pokemon{}
}

func (r *Pokemon) FetchPokemon(name string) error {
	url := PokemonApi + name
	resp, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("-Failed to create a req: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode > 299 {
		return fmt.Errorf("-Bad status code: %v", resp.StatusCode)
	}

	if err = json.NewDecoder(resp.Body).Decode(r); err != nil {
		return fmt.Errorf("-Failed to decode resp body: %v", err)
	}
	return nil
}

func (p *Pokemon) TryCatch(pokemonName string) bool {
	fmt.Printf("Throwing a Pokeball at %s...\n", pokemonName)
	catchChance := 100 - (p.BaseExp / 5)
	if catchChance < 5 {
		catchChance = 5 // Minimum 5% chance
	}
	if catchChance > 95 {
		catchChance = 95 // Maximum 95% chance
	}

	roll := rand.Intn(100) + 1

	message := fmt.Sprintf(
		"Threw a Pokéball at %s... (Catch chance: %d%%, Roll: %d)\n",
		p.Name, catchChance, roll,
	)

	f, err := os.OpenFile("catch.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Printf("Error writing log: %v\n", err)
		return false
	}
	defer f.Close()

	if _, err := f.WriteString(message); err != nil {
		fmt.Printf("Error writing log: %v\n", err)
		return false
	}

	return roll <= catchChance
}
