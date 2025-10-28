package main

import (
	"context"

	"github.com/nhmosko/gator/internal/database"
)

func middlewareLoggedIn(handler func(s *State, cmd command, user database.User) error) func(*State, command) error {
	return func(s *State, cmd command) error {
		user, err := s.db.GetUser(context.Background(), s.PConfig.CurrentUserName)
		if err != nil {
			return err
		}
		return handler(s, cmd, user)
	}
}
