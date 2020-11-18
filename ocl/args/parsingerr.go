package args

import "fmt"

const (
	// RequiredArgumentMissing is used when a required argument is missing.
	RequiredArgumentMissing = iota

	// InsufficientArguments is used when at least one greedy argument is required.
	InsufficientArguments

	// InvalidArgument is used when an invalid input is used for an argument.
	InvalidArgument
)

// ParsingError is an error in argument transformation.
type ParsingError struct {
	ErrorType      int
	argDisplayName string
}

func (err *ParsingError) errorStrToFormat() string {
	switch err.ErrorType {
	case RequiredArgumentMissing:
		return "An invalid argument for %v was supplied"
	case InsufficientArguments:
		return "You are missing some arguments in the %v for this command"
	case InvalidArgument:
		return "You supplied an argument that was invalid or malformed for %v"
	default:
		return "An error occurred while processing the command's arguments (faulty argument: %v), please try again later"
	}
}

func (err *ParsingError) Error() string {
	return fmt.Sprintf(err.errorStrToFormat(), err.argDisplayName)
}

// UIError is an error string suitable for sending to the user.
func (err *ParsingError) UIError() string {
	return fmt.Sprintf(err.errorStrToFormat(), "`"+err.argDisplayName+"`")
}

// NewParsingError creates an argument parsing error that's user friendly.
func NewParsingError(arg *Argument, errorType int) *ParsingError {
	return &ParsingError{errorType, arg.ID}
}
