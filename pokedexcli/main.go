package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

type cliCommand struct {
	name        string
	description string
	callback    func() error
}

func cleanInput(text string) []string {
	cleanTextArr := strings.Fields(text)
	for i, text := range cleanTextArr {
		cleanTextArr[i] = strings.ToLower(text)
	}
	return cleanTextArr
}

func commandHelp() error {
	fmt.Println("Welcome to the Pokedex!")
	fmt.Println("Usage:")
	fmt.Println()
	fmt.Println("help: Displays a help message")
	fmt.Println("exit: Exit the Pokedex")
	return nil
}

func commandExit() error {
	fmt.Println("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return nil
}

func getCmdRegistry(cmdName string) (cliCommand, error) {
	switch cmdName {
	case "exit":
		return cliCommand{
			name:        "exit",
			description: "Exit the pokedex",
			callback:    commandExit,
		}, nil
	case "help":
		return cliCommand{
			name:        "help",
			description: "Displays a help message",
			callback:    commandHelp,
		}, nil
	default:
		return cliCommand{}, fmt.Errorf("Unknown command")
	}
}

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	fmt.Print("Pokedex > ")

	for scanner.Scan() {
		line := scanner.Text()
		if len(line) > 0 {
			cleanLine := cleanInput(line)
			cmd, err := getCmdRegistry(cleanLine[0])
			if err != nil {
				fmt.Println(err)
			} else {
				cmd.callback()
			}
		}
		fmt.Print("Pokedex > ")
	}
}
