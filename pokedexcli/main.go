package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"

	pokeapi "github.com/re1n-e/pokeApi"
	pokecache "github.com/re1n-e/pokeCache"
	pokemon "github.com/re1n-e/pokemon"
)

type config struct {
	nextLocationUrl      string
	prevLocationUrl      string
	mapCache             pokecache.Cache
	pokemonLocationCache pokecache.Cache
	pokemonCache         pokecache.Cache
	pokeDex              map[string]pokemon.Pokemon
	otherArg             string
}

const mapCacheTime = 5
const pokemonLocationCacheTime = 7

func configNew() config {
	return config{
		nextLocationUrl:      pokeapi.PokeApi + "location-area/",
		prevLocationUrl:      "",
		mapCache:             *pokecache.NewCache(mapCacheTime * time.Second),
		pokemonLocationCache: *pokecache.NewCache(pokemonLocationCacheTime * time.Second),
		pokemonCache:         *pokecache.NewCache(pokemonLocationCacheTime * time.Second),
		pokeDex:              make(map[string]pokemon.Pokemon),
		otherArg:             "",
	}
}

type cliCommand struct {
	name        string
	description string
	callback    func(*config) error
}

func cleanInput(text string) []string {
	cleanTextArr := strings.Fields(text)
	for i, text := range cleanTextArr {
		cleanTextArr[i] = strings.ToLower(text)
	}
	return cleanTextArr
}

func commandHelp(cfg *config) error {
	fmt.Println()
	fmt.Println("Welcome to the Pokedex!")
	fmt.Println("Usage:")
	fmt.Println()
	for _, cmd := range getCommands() {
		fmt.Printf("%s: %s\n", cmd.name, cmd.description)
	}
	fmt.Println()
	return nil
}

func commandExit(*config) error {
	fmt.Println("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return nil
}

func commandMap(config *config) error {
	loc, ok := config.mapCache.Get(config.nextLocationUrl)
	location := pokeapi.LocationNew()
	if ok {
		if err := json.Unmarshal(loc, &location); err != nil {
			return fmt.Errorf("-Failed to unmarshal struct: %v", err)
		}
	} else {
		if err := location.GetNext(config.nextLocationUrl); err != nil {
			return err
		}
	}
	if location.Next == nil {
		config.nextLocationUrl = ""
	} else {
		config.nextLocationUrl = *location.Next
	}
	if location.Previous == nil {
		config.prevLocationUrl = ""
	} else {
		config.prevLocationUrl = *location.Previous
	}
	for _, city := range location.Results {
		fmt.Println(city.Name)
	}
	bytes, err := json.Marshal(location)
	if err != nil {
		return fmt.Errorf("-Failed to marshal struct: %v", err)
	}
	config.mapCache.Add(config.prevLocationUrl, bytes)
	return nil
}

func commandMapB(config *config) error {
	loc, ok := config.mapCache.Get(config.prevLocationUrl)
	location := pokeapi.LocationNew()
	if ok {
		if err := json.Unmarshal(loc, &location); err != nil {
			return fmt.Errorf("-Failed to unmarshal struct: %v", err)
		}
	} else {
		if err := location.GetPrevious(config.prevLocationUrl); err != nil {
			return err
		}
	}
	if location.Next == nil {
		config.nextLocationUrl = ""
	} else {
		config.nextLocationUrl = *location.Next
	}
	if location.Previous == nil {
		config.prevLocationUrl = ""
	} else {
		config.prevLocationUrl = *location.Previous
	}
	for _, city := range location.Results {
		fmt.Println(city.Name)
	}
	return nil
}

func commandExplore(config *config) error {
	if config.otherArg == "" {
		return fmt.Errorf("-Usage: explore <city-name>")
	}
	enc, ok := config.pokemonLocationCache.Get(config.otherArg)
	encounters := pokeapi.EncounterNew()
	if ok {
		if err := json.Unmarshal(enc, &encounters); err != nil {
			return fmt.Errorf("-Failed to unmarshal struct: %v", err)
		}
	} else {
		if err := encounters.GetPokemons(pokeapi.PokeApi + "location-area/" + config.otherArg); err != nil {
			return fmt.Errorf("-Error fetching pokemons location: %v", err)
		}
	}
	fmt.Printf("Exploring %s...\n", config.otherArg)
	fmt.Println("Found Pokemon:")
	for _, pokemon := range encounters.PokemonEncounter {
		fmt.Printf(" - %s\n", pokemon.Pokemon.Name)
	}
	bytes, err := json.Marshal(encounters)
	if err != nil {
		return fmt.Errorf("-Failed to marshal struct: %v", err)
	}
	config.pokemonLocationCache.Add(config.otherArg, bytes)
	return nil
}

func commandCatch(config *config) error {
	if config.otherArg == "" {
		return fmt.Errorf("-Usage: catch <pokemon-name>")
	}
	pokeData, ok := config.pokemonCache.Get(config.otherArg)
	pokemon := pokemon.NewPokemon()
	if ok {
		if err := json.Unmarshal(pokeData, &pokemon); err != nil {
			return fmt.Errorf("-Failed to unmarshall struct: %v", err)
		}
	} else {
		if err := pokemon.FetchPokemon(config.otherArg); err != nil {
			return fmt.Errorf("-Error fetching pokemons")
		}
	}

	if pokemon.TryCatch(config.otherArg) {
		config.pokeDex[config.otherArg] = pokemon
		fmt.Printf("%s was caught!\n", config.otherArg)
	} else {
		fmt.Printf("%s escaped!\n", config.otherArg)
	}

	bytes, err := json.Marshal(pokemon)
	if err != nil {
		return fmt.Errorf("-Failed to marshal struct: %v", err)
	}
	config.pokemonCache.Add(config.otherArg, bytes)
	return nil
}

func commandInspect(config *config) error {
	if config.otherArg == "" {
		return fmt.Errorf("-Usage: inspect <pokemon-name>")
	}
	pokemon, ok := config.pokeDex[config.otherArg]
	if !ok {
		return fmt.Errorf("you have not caught that pokemon")
	}
	fmt.Printf("Name: %s\n", pokemon.Name)
	fmt.Printf("Height: %d\n", pokemon.Height)
	fmt.Printf("Weight: %d\n", pokemon.Weight)
	fmt.Println("Stats:")
	for _, stat := range pokemon.Stats {
		fmt.Printf(" -%s: %d\n", stat.Name.Name, stat.Base)
	}
	fmt.Println("Types:")
	for _, ty := range pokemon.Type {
		fmt.Printf(" - %s\n", ty.TypeName.Name)
	}
	fmt.Println("Ability:")
	for _, ability := range pokemon.Abilities {
		fmt.Printf(" - %s\n", ability.Ability.Name)
	}
	fmt.Println("Moves:")
	for _, move := range pokemon.Moves {
		fmt.Printf(" - %s\n", move.Move.Name)
	}
	return nil
}

func commandPokeDex(config *config) error {
	fmt.Println("Your Pokedex:")
	for pokemon := range config.pokeDex {
		fmt.Printf(" - %s\n", pokemon)
	}
	return nil
}

func getCommands() map[string]cliCommand {
	return map[string]cliCommand{
		"help": {
			name:        "help",
			description: "Displays a help message",
			callback:    commandHelp,
		},
		"explore": {
			name:        "explore <location_name>",
			description: "Explore a location",
			callback:    commandExplore,
		},
		"map": {
			name:        "map",
			description: "Get the next page of locations",
			callback:    commandMap,
		},
		"mapb": {
			name:        "mapb",
			description: "Get the previous page of locations",
			callback:    commandMapB,
		},
		"exit": {
			name:        "exit",
			description: "Exit the Pokedex",
			callback:    commandExit,
		},
		"inspect": {
			name:        "inspect",
			description: "Inspect a pokemon that you have caught",
			callback:    commandInspect,
		},
		"pokedex": {
			name:        "pokedex",
			description: "Shows all caught pokemon",
			callback:    commandPokeDex,
		},
		"catch": {
			name:        "catch <pokemon_name>",
			description: "Catches a pokemon",
			callback:    commandCatch,
		},
	}
}

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	fmt.Print("Pokedex > ")
	config := configNew()
	cmds := getCommands()
	for scanner.Scan() {
		line := scanner.Text()
		if len(line) > 0 {
			cleanLine := cleanInput(line)
			cmd, ok := cmds[cleanLine[0]]
			if !ok {
				fmt.Println("Unknown Command")
			} else {
				if len(cleanLine) > 1 {
					config.otherArg = cleanLine[1]
				}
				if err := cmd.callback(&config); err != nil {
					fmt.Println(err)
				}
			}
		}
		fmt.Print("Pokedex > ")
	}
}
