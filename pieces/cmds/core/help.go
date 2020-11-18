package core

import (
	"fmt"

	"github.com/opolobot/opolo/ocl"
	"github.com/opolobot/opolo/ocl/args"
	"github.com/opolobot/opolo/pieces/parsers"
)

func init() {
	cmd := &ocl.Command{
		Name:        "help",
		Aliases:     []string{"h", "cmds", "commands"},
		Description: "Provides help for using the opolo.",
		Arguments:   []*args.Argument{args.Create("[cmd]", &parsers.String{})},
		Stack:       []ocl.Middleware{help},
	}

	Category.AddCommand(cmd)
}

func help(ctx *ocl.Context, _ ocl.Next) {
	cmdName := ctx.Args["cmd"].(string)
	if cmdName != "" {
		ctx.Send(fmt.Sprintf("found cmd arg %v", cmdName))
		return
	}

	ctx.Send("Send full help here.")
}

func individualCmdHelp(cmdName string) {}