package main

import (
	"context"
	"fmt"
	"time"

	"github.com/bevane/gator/internal/config"
	"github.com/bevane/gator/internal/database"
	"github.com/google/uuid"
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
	_, err := s.db.GetUser(context.Background(), username)
	if err != nil {
		return fmt.Errorf("%s does not exist", username)
	}
	err = s.cfg.SetUser(username)
	if err != nil {
		return err
	}
	fmt.Printf("User has been set to: %v\n", username)
	return nil
}

func handlerRegister(s *state, cmd command) error {
	if len(cmd.args) == 0 {
		return fmt.Errorf("The login command expects one argument: username")
	}
	username := cmd.args[0]
	user, err := s.db.CreateUser(context.Background(), database.CreateUserParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name:      username,
	})
	if err != nil {
		return fmt.Errorf("%s already exists", username)
	}
	fmt.Printf("New user registerd: %v\n", user)
	err = s.cfg.SetUser(username)
	if err != nil {
		return err
	}
	fmt.Printf("User has been set to: %v\n", username)

	return nil
}

func handlerReset(s *state, cmd command) error {
	err := s.db.DeleteUsers(context.Background())
	if err != nil {
		return fmt.Errorf("Failed to delete users")
	}
	fmt.Println("All users have been deleted succesfully")
	return nil
}
