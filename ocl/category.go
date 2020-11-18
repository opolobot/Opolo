package ocl

// Category an incubator for commands of a specific type.
type Category struct {
	name string

	Commands []*Command
}

// AddCommand adds a command to the category.
func (cat *Category) AddCommand(cmd *Command) {
	cmd.category = cat
	cat.Commands = append(cat.Commands, cmd)
}

// NewCategory creates a new category.
func NewCategory(name string) *Category {
	return &Category{name, make([]*Command, 0)}
}
