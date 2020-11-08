package cmds

import (
	"github.com/zorbyte/whiskey/lib"
)

var randomCmds []*lib.Cmd

func init() {
	randomCmds = append(
		randomCmds,
		&lib.Cmd{
			Runner:      test,
			Name:        "test",
			Usage:       "[...args]",
			Aliases:     []string{"t"},
			Description: "A test command for developing the bot",
		},
	)
}

func test(ctx *lib.Ctx) (string, error) {
	return "test", nil
}
