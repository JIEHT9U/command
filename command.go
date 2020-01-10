package command

import (
	"context"

	"gopkg.in/alecthomas/kingpin.v2"
)

type Command interface {
	FlagParse(clause *kingpin.CmdClause)
	Exec(ctx context.Context, subCommands *Sub) error
}
