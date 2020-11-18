package ocl

// Category an incubator for commands of a specific type.
type Category struct {
	name  string

	Commands  []*Command
}

// AddCommand adds a command to the category.
func (cat *Category) AddCommand(cmd *Command) {
	cmd.Category = cat
	cat.Commands = append(cat.Commands, cmd)
}
