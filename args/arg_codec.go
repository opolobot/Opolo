package args

import "fmt"

// ArgumentParser takes a raw argument string and outputs an interface for use.
type ArgumentParser func(rawArgument string) (interface{}, error)

// ArgumentCodec is information about an argument and how to process it.
type ArgumentCodec struct {
	Name      string
	ExtraInfo string
	Parser    ArgumentParser
	Greedy    bool
	Required  bool
}

// DisplayName is the user friendly appearence of an argument.
func (codec *ArgumentCodec) DisplayName() string {
	var container string
	if codec.Required {
		container = "<%v>"
	} else {
		container = "[%v]"
	}

	var containerValue string
	if codec.Greedy {
		containerValue = "..." + codec.Name
	} else {
		containerValue = codec.Name
	}

	// Useful for arguments that appear as <someInt (<=500)>
	if codec.ExtraInfo != "" {
		containerValue = containerValue + " (" + codec.ExtraInfo + ")"
	}

	return fmt.Sprintf(container, containerValue)
}
