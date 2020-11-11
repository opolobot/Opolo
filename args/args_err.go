package args

import "fmt"

const (
	requiredArgumentMissing = iota
	insufficientArguments
	invalidArgument
)

// ParsingError is an error in argument parsing.
type ParsingError struct {
	ErrorType      int
	wrapped        error
	argDisplayName string
}

func (err *ParsingError) errorStrToFormat() string {
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

func (err *ParsingError) Error() string {
	return fmt.Sprintf(err.errorStrToFormat(), "`"+err.argDisplayName+"`")
}

func (err *ParsingError) Unwrap() error {
	return err.wrapped
}

// NewParsingError creates an argument parsing error thats user friendly.
func NewParsingError(argDisplayName string, errorType int, wrapped ...error) error {
	err := &ParsingError{
		argDisplayName: argDisplayName,
		ErrorType:      errorType,
		wrapped:        nil,
	}

	if len(wrapped) > 0 {
		err.wrapped = wrapped[0]
	}

	return err
}
