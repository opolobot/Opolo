package cmds

import "fmt"

const (
	requiredArgumentMissing = iota
	insufficientArguments
	invalidArgument
)

type parsingError struct {
	ErrorType      int
	wrapped        error
	argDisplayName string
}

func (err *parsingError) errorStrToFormat() string {
	switch err.ErrorType {
	case requiredArgumentMissing:
		return "An invalid argument for %v was supplied"
	case insufficientArguments:
		return "You are missing some arguments in the %v for this command"
	case invalidArgument:
		return "You supplied an argument that was invalid or malformed for %v"
	default:
		return "An error occurred while processing the command's arguments (fault argument: %v), please try again later"
	}
}

func (err *parsingError) Error() string {
	return fmt.Sprintf(err.errorStrToFormat(), "`"+err.argDisplayName+"`")
}

func (err *parsingError) Unwrap() error {
	return err.wrapped
}

// NewParsingError creates an argument parsing error thats user friendly.
func NewParsingError(argDisplayName string, errorType int, wrapped ...error) error {
	err := &parsingError{
		argDisplayName: argDisplayName,
		ErrorType:      errorType,
		wrapped:        nil,
	}

	if len(wrapped) > 0 {
		err.wrapped = wrapped[0]
	}

	return err
}
