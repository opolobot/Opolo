package arg

import "fmt"

// Parser takes a raw argument string and outputs an interface for use.
type Parser func(arg *Argument, rawArg string) (interface{}, error)

// Argument is information about an argument and how to process it.
type Argument struct {
	Parser Parser

	Name        string
	Constraints string

	Greedy   bool
	Required bool
}

// DisplayName is the user friendly appearance of an argument.
func (arg *Argument) DisplayName() string {
	var container string
	if arg.Required {
		container = "<%v>"
	} else {
		container = "[%v]"
	}

	var containerValue string
	if arg.Greedy {
		containerValue = "..." + arg.Name
	} else {
		containerValue = arg.Name
	}

	// Useful for arguments that appear as <someInt (<=500)>
	if arg.Constraints != "" {
		containerValue = containerValue + " (" + arg.Constraints + ")"
	}

	return fmt.Sprintf(container, containerValue)
}

func (arg *Argument) parse(rawArg string) (interface{}, error) {
	return arg.Parser(arg, rawArg)
}
