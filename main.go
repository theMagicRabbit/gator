package main

import (
	"database/sql"
	"fmt"
	"os"

	"github.com/theMagicRabbit/gator/internal/cli"
	"github.com/theMagicRabbit/gator/internal/config"
	"github.com/theMagicRabbit/gator/internal/database"
	"github.com/theMagicRabbit/gator/internal/state"

	_ "github.com/lib/pq"
)


func main() {
	conf, err := config.Read()
	if err != nil {
		os.Exit(1)
	}
	db, err := sql.Open("postgres", conf.Db_url)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	dbQueries := database.New(db)
	runState := state.State{
		Config: &conf,
		Db: dbQueries,
	}
	commands := cli.Commands{
		Commands: map[string]func(*state.State, cli.Command) error {},
	}
	commands.Register("addfeed", cli.HandlerAddFeed)
	commands.Register("agg", cli.HandlerAgg)
	commands.Register("feeds", cli.HandlerFeeds)
	commands.Register("follow", cli.HandlerFollow)
	commands.Register("following", cli.HandlerFollowing)
	commands.Register("login", cli.HandlerLogin)
	commands.Register("register", cli.HandlerRegister)
	commands.Register("reset", cli.HandlerReset)
	commands.Register("users", cli.HandlerUsers)

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

