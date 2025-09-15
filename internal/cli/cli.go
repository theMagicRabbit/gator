package cli

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/lib/pq"
	"github.com/theMagicRabbit/gator/internal/database"
	"github.com/theMagicRabbit/gator/internal/feed"
	"github.com/theMagicRabbit/gator/internal/state"
)

type Command struct {
	Name string;
	Args []string;
}

type Commands struct {
	Commands map[string]func(*state.State, Command) error;
}

func (c Commands) Run(s *state.State, cmd Command) error {
	if s == nil {
		return fmt.Errorf("State is null")
	}
	handler, ok := c.Commands[cmd.Name]
	if !ok {
		return fmt.Errorf("%s is not a known command", cmd.Name)
	}
	err := handler(s, cmd)
	if err != nil {
		return err
	}
	return nil
}

func (c Commands) Register(name string, f func(*state.State, Command) error) {
	if _, keyExists := c.Commands[name]; keyExists {
		return
	}
	c.Commands[name] = f
}

func HandlerAddFeed(s *state.State, cmd Command) error {
	if argLen := len(cmd.Args); argLen < 2 {
		return fmt.Errorf("addfeed requires two argument; zero provided.")
	} else if argLen > 2 {
		return fmt.Errorf("addfeed requires two argument; %d provided.", argLen)
	}
	utcTime := time.Now().UTC()
	user_id, err := s.Db.GetUser(context.Background(), s.Config.Current_user_name)
	if err != nil {
		return err
	}
	params := database.CreateFeedParams{
		ID:  uuid.New(),
		CreatedAt: utcTime,
		UpdatedAt: utcTime,
		Name: cmd.Args[0],
		Url: cmd.Args[1],
		UserID: user_id.ID,
	}
	createFeed, err := s.Db.CreateFeed(context.Background(), params)
	if err != nil {
		return err
	}
	fmt.Printf("%+v\n", createFeed)
	return nil
}

func HandlerAgg(s *state.State, cmd Command) error {
	url := "https://www.wagslane.dev/index.xml"
	feed, err := feed.FetchFeed(context.Background(), url)
	if err != nil {
		return err
	}
	fmt.Println(feed)
	return nil
}

func HandlerFeeds(s *state.State, cmd Command) error {
	feed, err := s.Db.GetAllFeeds(context.Background())
	if err != nil {
		return err
	}
	for i, f := range feed {
		username, err := s.Db.GetUserFromID(context.Background(), f.UserID)
		if err != nil {
			continue
		}
		fmt.Printf("[Feed %d]\nname: %s\nurl: %s\nusername: %s\n", i, f.Name, f.Url, username.Name)
	}
	return nil
}

func HandlerFollow(s *state.State, cmd Command) error {
	if argLen := len(cmd.Args); argLen < 1 {
		return fmt.Errorf("follow requires one argument; zero provided.")
	} else if argLen > 1 {
		return fmt.Errorf("follow requires one argument; %d provided.", argLen)
	}
	utcTime := time.Now().UTC()
	user, err := s.Db.GetUser(context.Background(), s.Config.Current_user_name)
	if err != nil {
		return err
	}
	feed, err := s.Db.GetFeed(context.Background(), cmd.Args[0])
	if err != nil {
		return err
	}
	params := database.CreateFeedFollowsParams{
		ID: uuid.New(),
		CreatedAt: utcTime,
		UpdatedAt: utcTime,
		UserID: user.ID,
		FeedID: feed.ID,
	}
	following, err := s.Db.CreateFeedFollows(context.Background(), params)
	if err != nil {
		return err
	}
	fmt.Println(following)
	
	return nil
}

func HandlerLogin(s *state.State, cmd Command) error {
	if len(cmd.Args) < 1 {
		return fmt.Errorf("Login requires one argument; zero provided.")
	}
	userName := cmd.Args[0]
	var existingUser database.User
	var err error
	if existingUser, err = s.Db.GetUser(context.Background(), userName); errors.Is(err, sql.ErrNoRows) {
		return fmt.Errorf("User '%s' does not exist", userName)
	} else if err != nil {
		return err
	}

	err = s.Config.SetUser(existingUser.Name)
	if err != nil {
		return err
	}
	fmt.Printf("Login success! Current user: %s\n", cmd.Args[0])
	return nil
}

func HandlerRegister(s *state.State, cmd Command) error {
	if argLen := len(cmd.Args); argLen < 1 {
		return fmt.Errorf("Register requires one argument; zero provided.")
	} else if argLen > 1 {
		return fmt.Errorf("Register requires one argument; %d provided.", argLen)
	}
	utcTime := time.Now().UTC()
	newUsername := cmd.Args[0]
	params := database.CreateUserParams {
		ID: uuid.New(),
		CreatedAt: utcTime,
		UpdatedAt: utcTime,
		Name: newUsername,
	}
	createUser, err := s.Db.CreateUser(context.Background(), params)
	if err != nil {
		var pqErr *pq.Error
		if errors.As(err, &pqErr) {
			if pqErr.Code == "23505" {
				return fmt.Errorf("User %s already exists", newUsername)
			}
		}
		return err
	}
	err = s.Config.SetUser(createUser.Name)
	if err != nil {
		return err
	}
	fmt.Printf("User '%s' created: %s\n", createUser.Name, createUser)
	return nil
}

func HandlerReset(s *state.State, cmd Command) error {
	if err := s.Db.DeleteAllUsers(context.Background()); err != nil {
		return err
	}
	return nil
}

func HandlerUsers(s *state.State, cmd Command) error {
	usernames, err := s.Db.GetAllUsers(context.Background())
	if err != nil {
		return err
	}
	for _, name := range usernames {
		if name == s.Config.Current_user_name {
			fmt.Printf("* %s (current)\n", name)
			continue
		}
		fmt.Printf("* %s\n", name)
	}
	return nil
}
