package cmd

import "github.com/TeamWhiskey/whiskey/arg"

// Builder ensures easy command creation.
type Builder struct {
	cmd *Command
}

// Name sets the name of the command.
func (bldr *Builder) Name(name string) {
	bldr.cmd.Name = name
}

// Aliases sets the aliases of the command.
func (bldr *Builder) Aliases(aliases ...string) {
	bldr.cmd.Aliases = append(bldr.cmd.Aliases, aliases...)
}

// Description sets the description of the command
func (bldr *Builder) Description(desc string) {
	bldr.cmd.Description = desc
}

// Arg adds an argument for the command.
func (bldr *Builder) Arg(arg *arg.Argument) {
	bldr.cmd.args = append(bldr.cmd.args, arg)
}

// Use adds a middleware to the stack.
func (bldr *Builder) Use(middleware Middleware) {
	bldr.cmd.stack = append(bldr.cmd.stack, middleware)
}

// Command returns the constructed command.
func (bldr *Builder) Command() *Command {
	return bldr.cmd
}

// New creates a command builder instance.
func New() *Builder {
	return &Builder{
		cmd: &Command{},
	}
}
