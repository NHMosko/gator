package main

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/nhmosko/gator/internal/database"
)

func handlerLogin(s *State, cmd command) error {
	if len(cmd.Args) != 1 {
		return fmt.Errorf("usage: %s <username>", cmd.Name)
	}
	name := cmd.Args[0]

	if _, err := s.db.GetUser(context.Background(), name); err != nil {
		return fmt.Errorf("user not registered: %v", err)
	}

	err := s.PConfig.SetUser(name)
	if err != nil {
		return fmt.Errorf("couldn't set the current user: %w", err)
	}

	fmt.Println("User switched successfully!")
	return nil
}

func handlerRegister(s *State, cmd command) error {
	if len(cmd.Args) != 1 {
		return fmt.Errorf("usage: %s <username>", cmd.Name)
	}
	name := cmd.Args[0]

	_, err := s.db.CreateUser(context.Background(), database.CreateUserParams{
		ID: uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name: name,
	})
	if err != nil {
		return fmt.Errorf("unable to register user: %v", err)
	}
	
	err = s.PConfig.SetUser(name)
	if err != nil {
		return fmt.Errorf("couldn't set the current user: %w", err)
	}

	fmt.Println("User created and switched successfully!")
	return nil
}

func handlerReset(s *State, cmd command) error {
	err := s.db.DeleteUsers(context.Background())
	if err != nil {
		return fmt.Errorf("unable to reset database: %w", err)
	}
	fmt.Println("successfully reset database")
	return nil
}

func handlerUsers(s *State, cmd command) error {
	users, err := s.db.ListUsers(context.Background())
	if err != nil {
		return fmt.Errorf("unable to list users: %w", err)
	}

	for _, user := range users {
		current := ""
		if user.Name == s.PConfig.CurrentUserName {
			current = " (current)"
		}
		fmt.Printf("* %s%s\n", user.Name, current)
	}

	return nil
}


func handlerAgg(s *State, cmd command) error {
	if len(cmd.Args) != 1 {
		return fmt.Errorf("usage: %s <time-between-reqs>", cmd.Name)
	}
	intervalString := cmd.Args[0]
	
	interval, err := time.ParseDuration(intervalString)
	if err != nil {
		return fmt.Errorf("unable to convert %v into a valid time: %w", intervalString, err)
	}

	fmt.Println("Collecting feeds every", interval)

	ticker := time.NewTicker(interval)
	for ;; <-ticker.C {
		scrapeFeeds(s)
	}

	return nil
}

func handlerAddFeed(s *State, cmd command, user database.User) error {
	if len(cmd.Args) != 2 {
		return fmt.Errorf("usage: %s <feed-name> <feed-url>", cmd.Name)
	}
	name := cmd.Args[0]
	url := cmd.Args[1]

	newFeed, err := s.db.CreateFeed(context.Background(), database.CreateFeedParams{
		ID: uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name: name,
		Url: url,
		UserID: user.ID,
	})
	if err != nil {
		return fmt.Errorf("unable to create feed: %w", err)
	}

	_, err = s.db.CreateFeedFollow(context.Background(), database.CreateFeedFollowParams{
		ID: uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		UserID: user.ID,
		FeedID: newFeed.ID,
	})
	if err != nil {
		return fmt.Errorf("unable to follow added feed: %w", err)
	}

	fmt.Println(newFeed.Name, "successfully added")
	return nil
}

func handlerFeeds(s *State, cmd command) error {
	feeds, err := s.db.ListFeeds(context.Background())
	if err != nil {
		return fmt.Errorf("unable to get feeds: %w", err)
	}

	if len(feeds) == 0 {
		fmt.Println("no feeds added, try using 'addfeed <feed-name> <feed-url>'")
		return nil
	}

	fmt.Println("Feeds:")
	for _, feed := range feeds {
		userID, err := s.db.GetUserName(context.Background(), feed.UserID)
		if err != nil {
			return fmt.Errorf("user id: %s, generated error: %w", feed.UserID, err)
		}

		fmt.Printf("%v - %v (added by: %v)\n", feed.Name, feed.Url, userID)
	}

	return nil
}

func handlerFollow(s *State, cmd command, user database.User) error {
	if len(cmd.Args) != 1 {
		return fmt.Errorf("usage: %s <url>", cmd.Name)
	}
	url := cmd.Args[0]

	feed, err := s.db.GetFeed(context.Background(), url)
	if err != nil {
		return fmt.Errorf("unable to find feed: %w", err)
	}

	feedFollowData, err := s.db.CreateFeedFollow(context.Background(), database.CreateFeedFollowParams{
		ID: uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		FeedID: feed.ID,
		UserID: user.ID,
	})

	fmt.Printf("%s now follows feed %s\n", feedFollowData.UserName, feedFollowData.FeedName)
	return nil
}

func handlerFollowing(s *State, cmd command, user database.User) error {
	following, err := s.db.GetFeedFollowsForUser(context.Background(), user.ID)
	if err != nil {
		return fmt.Errorf("unable to get followed feeds: %w", err)
	}

	if len(following) == 0 {
		fmt.Printf("user %s doesn't follow any feeds\ntry running 'follow <feed-url>'\n", user.Name)
		return nil
	}

	for _, feed := range following {
		fmt.Println(feed.FeedName)
	}
	return nil
}

func handlerUnfollow(s *State, cmd command, user database.User) error {
	if len(cmd.Args) != 1 {
		return fmt.Errorf("usage: %s <url>", cmd.Name)
	}
	url := cmd.Args[0]

	feed, err := s.db.GetFeed(context.Background(), url)
	if err != nil {
		return fmt.Errorf("unable to get feed: %w", err)
	}

	err = s.db.DeleteFeedFollow(context.Background(), database.DeleteFeedFollowParams{
		UserID: user.ID,
		FeedID: feed.ID,
	})
	if err != nil {
		return fmt.Errorf("unable to delete follow: %w", err)
	}

	fmt.Println("successfully unfollowed", feed.Name)
	return nil
}
