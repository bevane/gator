package main

import (
	"fmt"

	"github.com/bevane/gator/internal/config"
	"github.com/bevane/gator/internal/database"
)

type state struct {
	cfg *config.Config
	db  *database.Queries
}

type command struct {
	name string
	args []string
}

type commands struct {
	nameToCommand map[string]func(*state, command) error
}

func (c *commands) register(name string, f func(*state, command) error) {
	c.nameToCommand[name] = f
	return
}

func (c *commands) run(s *state, cmd command) error {
	commandHandler, ok := c.nameToCommand[cmd.name]
	if !ok {
		return fmt.Errorf("%v is not a valid command", cmd.name)
	}
	err := commandHandler(s, cmd)
	if err != nil {
		return err
	}
	return nil
}

func handlerLogin(s *state, cmd command) error {
	if len(cmd.args) == 0 {
		return fmt.Errorf("The login command expects one argument: username")
	}
	username := cmd.args[0]
	err := s.cfg.SetUser(username)
	if err != nil {
		return err
	}
	fmt.Printf("User has been set to: %v\n", username)
	return nil
}
