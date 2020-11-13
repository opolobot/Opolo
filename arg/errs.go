package arg

import "fmt"

const (
	// InternalParsingError is used when the parser itself fails for a reason unrelated to the data type.
	InternalParsingError = iota - 1

	// RequiredArgumentMissing is used when a required argument is missing.
	RequiredArgumentMissing

	// InsufficientArguments is used when at least one greedy argument is required.
	InsufficientArguments

	// InvalidArgument is used when an invalid input is used for an argument.
	InvalidArgument
)

// ParsingError is an error in argument transformation.
type ParsingError struct {
	ErrorType      int
	wrapped        error
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
		return "An error occurred while processing the command's arguments (fault argument: %v), please try again later"
	}
}

func (err *ParsingError) Error() string {
	return fmt.Sprintf(err.errorStrToFormat(), err.argDisplayName)
}

// UIError is an error string suitable for sending to the user.
func (err *ParsingError) UIError() string {
	return fmt.Sprintf(err.errorStrToFormat(), "`"+err.argDisplayName+"`")
}

func (err *ParsingError) Unwrap() error {
	return err.wrapped
}

// NewParsingError creates an argument parsing error that's user friendly.
func NewParsingError(arg *Argument, errorType int, wrapped ...error) *ParsingError {
	err := &ParsingError{
		argDisplayName: arg.DisplayName(),
		ErrorType:      errorType,
		wrapped:        nil,
	}

	if len(wrapped) > 0 {
		err.wrapped = wrapped[0]
	}

	return err
}
