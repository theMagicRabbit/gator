package main

import (
	"fmt"
	"github.com/theMagicRabbit/gator/internal/cli"
	"github.com/theMagicRabbit/gator/internal/config"
	"github.com/theMagicRabbit/gator/internal/state"
	"os"
)


func main() {
	conf, err := config.Read()
	if err != nil {
		os.Exit(1)
	}
	runState := state.State{Config: &conf,}
	commands := cli.Commands{
		Commands: map[string]func(*state.State, cli.Command) error {},
	}
	commands.Register("login", cli.HandlerLogin)
	if len(os.Args) < 2 {
		fmt.Println("No arguments provided")
		os.Exit(1)
	}
	cmdName := os.Args[1]
	args := os.Args[2:]
	cmd := cli.Command{
		Name: cmdName,
		Args: args,
	}
	err = commands.Run(&runState, cmd)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	os.Exit(0)
}

