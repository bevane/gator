package main

import (
	"database/sql"
	"fmt"
	"os"

	"github.com/bevane/gator/internal/config"
	"github.com/bevane/gator/internal/database"

	_ "github.com/lib/pq"
)

func main() {
	appState := state{}
	cfg, err := config.Read()
	db, err := sql.Open("postgres", cfg.DbUrl)
	if err != nil {
		fmt.Printf("Error opening db: %v\n", err)
		os.Exit(1)
	}
	appState.db = database.New(db)
	if err != nil {
		fmt.Printf("Error reading config file: %v\n", err)
		os.Exit(1)
	}
	appState.cfg = &cfg

	appCommands := commands{nameToCommand: make(map[string]func(*state, command) error)}
	appCommands.register("login", handlerLogin)
	appCommands.register("register", handlerRegister)
	appCommands.register("reset", handlerReset)
	appCommands.register("users", handlerUsers)

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
