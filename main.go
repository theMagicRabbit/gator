package main

import (
	"context"
	"database/sql"
	"embed"
	"fmt"
	"os"

	"github.com/pressly/goose/v3"
	"github.com/theMagicRabbit/gator/internal/cli"
	"github.com/theMagicRabbit/gator/internal/config"
	"github.com/theMagicRabbit/gator/internal/database"
	"github.com/theMagicRabbit/gator/internal/state"

	_ "github.com/lib/pq"
)

func middlewareLoggedIn(handler func(s *state.State, cmd cli.Command, user database.User) error) func(*state.State, cli.Command) error {
	 return func(s *state.State, cmd cli.Command) error {
		user, err := s.Db.GetUser(context.Background(), s.Config.Current_user_name)
		if err != nil {
			return err
		}
		err = handler(s, cmd, user)
		if err != nil {
			return err
		}
		return nil
	}
}

//go:embed sql/schema/*.sql
var embededMigrations embed.FS


func main() {
	goose.SetBaseFS(embededMigrations)
	conf, err := config.Read()
	if err != nil {
		os.Exit(1)
	}
	db, err := sql.Open("postgres", conf.Db_url)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	err = goose.SetDialect("postgres")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	err = goose.Up(db, "sql/schema")
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
	commands.Register("addfeed", middlewareLoggedIn(cli.HandlerAddFeed))
	commands.Register("agg", cli.HandlerAgg)
	commands.Register("browse", middlewareLoggedIn(cli.HandlerBrowse))
	commands.Register("feeds", cli.HandlerFeeds)
	commands.Register("follow", middlewareLoggedIn(cli.HandlerFollow))
	commands.Register("following", middlewareLoggedIn(cli.HandlerFollowing))
	commands.Register("login", cli.HandlerLogin)
	commands.Register("register", cli.HandlerRegister)
	commands.Register("reset", cli.HandlerReset)
	commands.Register("unfollow", middlewareLoggedIn(cli.HandlerUnfollow))
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

