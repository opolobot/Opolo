package cmds

import "github.com/TeamWhiskey/whiskey/args"

// CommandBuilder ensures easy command creation.
type CommandBuilder struct {
	cmd *Command
}

// Name sets the name of the command.
func (cmdBldr *CommandBuilder) Name(name string) {
	cmdBldr.cmd.Name = name
}

// Aliases sets the aliases of the command.
func (cmdBldr *CommandBuilder) Aliases(aliases ...string) {
	cmdBldr.cmd.Aliases = append(cmdBldr.cmd.Aliases, aliases...)
}

// Description sets the description of the command
func (cmdBldr *CommandBuilder) Description(desc string) {
	cmdBldr.cmd.Description = desc
}

// Args sets the arg codecs for the command.
func (cmdBldr *CommandBuilder) Args(argCodecs ...*args.ArgumentCodec) {
	cmdBldr.cmd.argCodecs = argCodecs
}

// Use adds a middleware to the stack.
func (cmdBldr *CommandBuilder) Use(middleware CommandMiddleware) {
	cmdBldr.cmd.stack = append(cmdBldr.cmd.stack, middleware)
}

// Build returns the constructed command.
func (cmdBldr *CommandBuilder) Build() *Command {
	return cmdBldr.cmd
}

// NewCommandBuilder creates a command builder instance.
func NewCommandBuilder() *CommandBuilder {
	return &CommandBuilder{
		cmd: &Command{},
	}
}
