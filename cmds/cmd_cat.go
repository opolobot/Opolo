package cmds

// CommandCategory an incubator for commands of a specific type.
type CommandCategory struct {
	name  string
	emoji string
	cmds  []*Command
}

// AddCommand adds a command to the category.
func (cmdCat *CommandCategory) AddCommand(cmd *Command) {
	cmdCat.cmds = append(cmdCat.cmds, cmd)
}

func (cmdCat *CommandCategory) displayName() string {
	return cmdCat.emoji + " __**" + cmdCat.name + "**__"
}

// NewCommandCategory creates new command category.
func NewCommandCategory(name, emoji string) *CommandCategory {
	return &CommandCategory{
		name:  name,
		emoji: emoji,
	}
}
