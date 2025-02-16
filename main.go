package main

import (
	"fmt"
	"os"

	"github.com/bevane/gator/internal/config"
)

func main() {
	appState := state{}
	cfg, err := config.Read()
	if err != nil {
		fmt.Printf("Error reading config file: %v\n", err)
		os.Exit(1)
	}
	appState.configPtr = &cfg

	appCommands := commands{nameToCommand: make(map[string]func(*state, command) error)}
	appCommands.register("login", handlerLogin)

	userArgs := os.Args
	if len(userArgs) < 2 {
		fmt.Println("At least one argument required to run")
		os.Exit(1)
	}
	userCommand := command{name: userArgs[1], args: userArgs[2:]}
	err = appCommands.run(&appState, userCommand)
	if err != nil {
		fmt.Printf("Error running %v command: %v\n", userCommand.name, err)
		os.Exit(1)
	}
}
