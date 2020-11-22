package ocl

import "github.com/opolobot/Opolo/ocl/args"

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

// Args adds argument for the command.
func (bldr *Builder) Args(arguments ...*args.Argument) {
	bldr.cmd.args = append(bldr.cmd.args, arguments...)
}

// PermissionLevel sets the default permission level for the command.
func (bldr *Builder) PermissionLevel(lvl PermissionLevel) {
	bldr.cmd.PermissionLevel = lvl
}

// Use adds a middleware to the stack.
func (bldr *Builder) Use(middleware Middleware) {
	bldr.cmd.stack = append(bldr.cmd.stack, middleware)
}

// New creates a command builder instance.
func New() *Builder {
	return &Builder{
		cmd: &Command{},
	}
}
