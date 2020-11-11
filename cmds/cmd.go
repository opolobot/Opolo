package cmds

import "strings"

// NextFunc runs the next middleware in the chain.
// Optionally supply an error to cancel the chain and report an error.
type NextFunc func(err ...error)

// CommandMiddleware is the middleware for a command.
type CommandMiddleware func(ctx *Context, next NextFunc)

// Command is a command
type Command struct {
	middleware     []CommandMiddleware
	name        string
	aliases     []string
	description string
	argCodecs   []*ArgumentCodec
	enabled     bool
}

// Usage returns how to use the command
func (cmd *Command) Usage() string {
	var usageBldr strings.Builder
	for i, codec := range cmd.argCodecs {
		usageBldr.WriteString(codec.DisplayName())
		if i < len(cmd.argCodecs)-1 {
			usageBldr.WriteString(" ")
		}
	}

	return usageBldr.String()
}
