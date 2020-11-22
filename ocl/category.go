package ocl

// Category an incubator for commands of a specific type.
type Category struct {
	Name  string
	Emoji string

	Commands []*Command
}

// Add adds a command via its builder to the category.
func (cat *Category) Add(cmdBldr *Builder) {
	cat.AddCommand(cmdBldr.cmd)
}

// AddCommand adds a command to the category.
func (cat *Category) AddCommand(cmd *Command) {
	cmd.category = cat
	cat.Commands = append(cat.Commands, cmd)
}

// NewCategory creates a new category.
func NewCategory(name, emoji string) *Category {
	return &Category{name, emoji, make([]*Command, 0)}
}
