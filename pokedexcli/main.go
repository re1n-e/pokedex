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
)

type config struct {
	nextLocationUrl string
	prevLocationUrl string
	mapCache        pokecache.Cache
	pokemonCache    pokecache.Cache
	otherArg        string
}

const mapCacheTime = 5
const pokemonCacheTime = 7

func configNew() config {
	return config{
		nextLocationUrl: pokeapi.PokeApi + "location-area/",
		prevLocationUrl: "",
		mapCache:        *pokecache.NewCache(mapCacheTime * time.Second),
		pokemonCache:    *pokecache.NewCache(pokemonCacheTime * time.Second),
		otherArg:        "",
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
	enc, ok := config.pokemonCache.Get(config.otherArg)
	encounters := pokeapi.EncounterNew()
	if ok {
		if err := json.Unmarshal(enc, &encounters); err != nil {
			return fmt.Errorf("-Failed to unmarshal struct: %v", err)
		}
	} else {
		if err := encounters.GetPokemons(pokeapi.PokeApi + "location-area/" + config.otherArg); err != nil {
			return fmt.Errorf("-Error fetching pokemons: %v", err)
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
	config.pokemonCache.Add(config.otherArg, bytes)
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
