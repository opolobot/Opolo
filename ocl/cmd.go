package ocl

import (
	"strings"

	"github.com/opolobot/opolo/ocl/args"
)

// Next runs the next middleware in the chain.
// Optionally supply an error to cancel the chain and report an error.
type Next func(err ...error)

// Middleware runs a chain of commands.
type Middleware func(ctx *Context, next Next)

// Command is a command
type Command struct {
	Name        string
	Description string
	Aliases     []string
	Arguments   []*args.Argument
	Stack       []Middleware
	Permission  int

	category *Category
	enabled  bool
}

// Category returns the category of the command.
func (cmd *Command) Category() *Category {
	return cmd.category
}

func (cmd *Command) usage() string {
	var usageBldr strings.Builder
	for i, arg := range cmd.Arguments {
		usageBldr.WriteString(arg.ID)
		if i < len(cmd.Arguments)-1 {
			usageBldr.WriteString(" ")
		}
	}

	return usageBldr.String()
}
