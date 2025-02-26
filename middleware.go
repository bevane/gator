package main

import (
	"context"
	"fmt"

	"github.com/bevane/gator/internal/database"
)

func middlewareLoggedIn(handler func(s *state, cmd command, user database.User) error) func(*state, command) error {
	return func(s *state, cmd command) error {
		user, err := s.db.GetUser(context.Background(), s.cfg.CurrentUserName)
		if err != nil {
			return fmt.Errorf("Error getting current user: user %s does not exist: %v", s.cfg.CurrentUserName, err)
		}
		err = handler(s, cmd, user)
		if err != nil {
			return err
		}
		return nil
	}
}
