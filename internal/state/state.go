package state

import (
	"github.com/theMagicRabbit/gator/internal/config"
	"github.com/theMagicRabbit/gator/internal/database"
)

type State struct {
	Config *config.Config;
	Db *database.Queries;
}

