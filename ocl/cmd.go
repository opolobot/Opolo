package ocl

import (
	"strings"

	"github.com/opolobot/opolo/ocl/args"
)

// NextFunc runs the next middleware in the chain.
// Optionally supply an error to cancel the chain and report an error.
type NextFunc func(err ...error)

// Middleware runs a chain of commands.
type Middleware func(ctx *Context, next NextFunc)

// Command is a command
type Command struct {
	Name        string
	Aliases     []string
	Description string
	Permission  int

	Category *Category

	args    []*args.Argument
	stack   []Middleware
	enabled bool

	help string
}

// Help builds or returns the pre-build help string for the command.
func (cmd *Command) Help() string {
	if cmd.help == "" {
		cmd.help = generateHelp(cmd)
	}

	return cmd.help
}

func (cmd *Command) usage() string {
	var usageBldr strings.Builder
	for i, arg := range cmd.args {
		usageBldr.WriteString(arg.ID)
		if i < len(cmd.args)-1 {
			usageBldr.WriteString(" ")
		}
	}

	return usageBldr.String()
}

func generateHelp(cmd *Command) string {
	var cmdHelpStrBldr strings.Builder

	// Use formatting directive for prefix interpolation later.
	cmdHelpStrBldr.WriteString("**`%v")

	// Name and usage.
	cmdHelpStrBldr.WriteString(cmd.Name)
	if usage := cmd.usage(); usage != "" {
		cmdHelpStrBldr.WriteString(" ")
		cmdHelpStrBldr.WriteString(usage)
	}

	// The first **` is missing for dynamic prefix support, as
	// including it and styling the dynamic prefix would ruin the formatting
	cmdHelpStrBldr.WriteString("`**")

	// Description.
	if cmd.Description != "" {
		cmdHelpStrBldr.WriteString("\n**desc ~** ")
		cmdHelpStrBldr.WriteString(cmd.Description)
	}

	// Aliases
	if len(cmd.Aliases) > 0 {
		cmdHelpStrBldr.WriteString("\n**aliases ~** ")
		cmdHelpStrBldr.WriteString(strings.Join(cmd.Aliases, ", "))
	}

	return cmdHelpStrBldr.String()
}
