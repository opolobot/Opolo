package ocl

import (
	"strings"

	"github.com/opolobot/Opolo/ocl/args"
)

// Next runs the next middleware in the chain.
// Optionally supply an error to cancel the chain and report an error.
type Next func(err ...error)

// Middleware runs a chain of commands.
type Middleware func(ctx *Context, next Next)

// Command is a command
type Command struct {
	Name            string
	Description     string
	Aliases         []string
	PermissionLevel PermissionLevel

	args  []*args.Argument
	stack []Middleware

	category *Category

	enabled bool
}

// Category returns the category of the command.
func (cmd *Command) Category() *Category {
	return cmd.category
}

// Usage describes how to use the command through the IDs of its arguments.
func (cmd *Command) Usage() string {
	var usageBldr strings.Builder
	for i, arg := range cmd.args {
		usageBldr.WriteString(arg.ID)
		if i < len(cmd.args)-1 {
			usageBldr.WriteString(" ")
		}
	}

	return usageBldr.String()
}
