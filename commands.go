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

func handlerUsers(s *state, cmd command) error {
	users, err := s.db.GetUsers(context.Background())
	if err != nil {
		return fmt.Errorf("Failed to get users")
	}
	var out string
	for _, user := range users {
		if user.Name == s.cfg.CurrentUserName {
			out += user.Name + " (current)\n"
		} else {
			out += user.Name + "\n"
		}
	}
	fmt.Print(out)
	return nil
}

func handlerAgg(s *state, cmd command) error {
	url := "https://www.wagslane.dev/index.xml"
	rssFeed, err := fetchFeed(context.Background(), url)
	if err != nil {
		return fmt.Errorf("Failed to fetch rss feed: %v", err)
	}
	fmt.Println(rssFeed)
	return nil
}

func handlerAddFeed(s *state, cmd command, user database.User) error {
	if len(cmd.args) < 2 {
		return fmt.Errorf("The addfeed command expects two arguments: name, url")
	}
	name := cmd.args[0]
	url := cmd.args[1]
	feed, err := s.db.CreateFeed(context.Background(), database.CreateFeedParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name:      name,
		Url:       url,
		UserID:    user.ID,
	})
	if err != nil {
		return fmt.Errorf("Failed to add rss feed: %v", err)
	}
	_, err = s.db.CreateFeedFollow(context.Background(), database.CreateFeedFollowParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		UserID:    user.ID,
		FeedID:    feed.ID,
	})
	if err != nil {
		return fmt.Errorf("Failed to create feed follow: %v", err)
	}
	fmt.Println(feed)
	return nil
}

func handlerFeeds(s *state, cmd command) error {
	feeds, err := s.db.GetFeedsWithUserName(context.Background())
	if err != nil {
		return fmt.Errorf("Failed to get rss feeds: %v", err)
	}
	for _, feed := range feeds {
		fmt.Println(feed)
	}
	return nil
}

func handlerFollow(s *state, cmd command, user database.User) error {
	if len(cmd.args) < 1 {
		return fmt.Errorf("The follow command expects one argument: url")
	}
	url := cmd.args[0]
	feed, err := s.db.GetFeedByUrl(context.Background(), url)
	if err != nil {
		return fmt.Errorf("Feed does not exist for url: %v", url)
	}
	feedFollow, err := s.db.CreateFeedFollow(context.Background(), database.CreateFeedFollowParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		UserID:    user.ID,
		FeedID:    feed.ID,
	})
	if err != nil {
		return fmt.Errorf("Failed to create feed follow: %v", err)
	}
	fmt.Printf("%v followed %v\n", feedFollow.UserName, feedFollow.FeedName)
	return nil
}

func handlerFollowing(s *state, cmd command, user database.User) error {
	feedFollows, err := s.db.GetFeedFollowsForUser(context.Background(), user.ID)
	if err != nil {
		return fmt.Errorf("Failed to get feed follows for user: %v", err)
	}
	if len(feedFollows) == 0 {
		fmt.Printf("%v has not followd any feeds\n", user.Name)
		return nil
	}
	out := ""
	for _, feedFollow := range feedFollows {
		out += feedFollow.FeedName + "\n"
	}
	fmt.Printf("%v has followed:\n%s", user.Name, out)
	return nil
}

func handlerUnfollow(s *state, cmd command, user database.User) error {
	if len(cmd.args) < 1 {
		return fmt.Errorf("The follow command expects one argument: url")
	}
	url := cmd.args[0]
	err := s.db.DeleteFeedFollowForUserByUrl(context.Background(), database.DeleteFeedFollowForUserByUrlParams{
		UserID: user.ID,
		Url:    url,
	})
	if err != nil {
		return fmt.Errorf("Failed to delete feed follow: %v", err)
	}
	fmt.Printf("%v unfollowed:\n%s", user.Name, url)
	return nil
}
