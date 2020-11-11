package cmds

import "fmt"

// ArgumentParser takes a raw argument string and outputs an interface for use.
type ArgumentParser func(rawArgument string) (interface{}, error)

// ArgumentCodec is information about an argument and how to process it.
type ArgumentCodec struct {
	name      string
	extraInfo string
	parser    ArgumentParser
	greedy    bool
	required  bool
}

// DisplayName is the user friendly appearence of an argument.
func (codec *ArgumentCodec) DisplayName() string {
	var container string
	if codec.required {
		container = "<%v>"
	} else {
		container = "[%v]"
	}

	var containerValue string
	if codec.greedy {
		containerValue = "..." + codec.name
	} else {
		containerValue = codec.name
	}

	// Useful for arguments that appear as <someInt (<=500)>
	if codec.extraInfo != "" {
		containerValue = containerValue + " (" + codec.extraInfo + ")"
	}

	return fmt.Sprintf(container, containerValue)
}
