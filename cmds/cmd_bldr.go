package cmds

// CommandBuilder ensures easy command creation.
type CommandBuilder struct {
	cmd *Command
}

// Name sets the name of the command.
func (cmdBldr *CommandBuilder) Name(name string) {
	cmdBldr.cmd.name = name
}

// Aliases sets the aliases of the command.
func (cmdBldr *CommandBuilder) Aliases(aliases ...string) {
	cmdBldr.cmd.aliases = append(cmdBldr.cmd.aliases, aliases...)
}

// Description sets the description of the command
func (cmdBldr *CommandBuilder) Description(desc string) {
	cmdBldr.cmd.description = desc
}

// Args sets the arg codecs for the command.
func (cmdBldr *CommandBuilder) Args(argCodecs ...*ArgumentCodec) {
	cmdBldr.cmd.argCodecs = argCodecs
}

// Build returns the constructed command.
func (cmdBldr *CommandBuilder) Build() *Command {
	return cmdBldr.cmd
}

// NewCommandBuilder creates a command builder instance.
func NewCommandBuilder() *CommandBuilder {
	return &CommandBuilder{}
}
