package cmd

// Category an incubator for commands of a specific type.
type Category struct {
	name  string
	emoji string

	Commands  []*Command
}

// AddCommand adds a command to the category.
func (cat *Category) AddCommand(cmd *Command) {
	cmd.Category = cat
	cat.Commands = append(cat.Commands, cmd)
}

// DisplayName provides a string suitable for the help menu.
func (cat *Category) DisplayName() string {
	return "**" + cat.name + "**"
}

// NewCategory creates new command category.
func NewCategory(name, emoji string) *Category {
	return &Category{
		name:  name,
		emoji: emoji,
	}
}
