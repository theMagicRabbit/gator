package cli

import (
	"fmt"
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

func HandlerLogin(s *state.State, cmd Command) error {
	if len(cmd.Args) < 1 {
		return fmt.Errorf("Login requires one argument; zero provided.")
	}
	err := s.Config.SetUser(cmd.Args[0])
	if err != nil {
		return err
	}
	fmt.Printf("Login success! Current user: %s\n", cmd.Args[0])
	return nil
}

