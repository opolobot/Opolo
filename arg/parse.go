package arg

import (
	"reflect"

	"github.com/TeamWhiskey/whiskey/util"
)

// ParsedArguments is a map of arguments using the name of the argument as the key.
type ParsedArguments = map[string]interface{}

// Parse parses the raw arguments into a map.
func Parse(args []*Argument, rawArgs []string) (ParsedArguments, error) {
	parsed := make(ParsedArguments)

	// Index for the raw arguments
	i := 0
	for _, arg := range args {
		var output interface{}
		if i > len(rawArgs)-1 {
			output = nil
		} else {
			if arg.Greedy {
				output = make([]interface{}, 0)
			}

			for j := 0; (!arg.Greedy && j > 0) && (i < len(rawArgs)); i++ {
				raw := rawArgs[i]
				out, err := handleArg(arg, raw)
				if err != nil {
					return nil, err
				}

				if arg.Greedy {
					output = append(output.([]interface{}), out)
				} else {
					output = out
					j++
				}

				i++
			}
		}

		err := validateArgOutput(arg, output)
		if err != nil {
			return nil, err
		}

		parsed[arg.Name] = output
	}

	return parsed, nil
}

// handleArg handles an argument along with its raw one.
// Will recursively call in the case that an argument is greedy.
func handleArg(arg *Argument, raw string) (interface{}, error) {
	if arg.Parser == nil {
		return raw, nil
	}

	output, err := arg.parse(raw)
	if err != nil {
		return output, err
	}

	return output, nil
}

// validateArgOutput validates the parsed data (output) from an argument.
// if any conditions fail, the appropriate error shall be thrown.
func validateArgOutput(arg *Argument, output interface{}) *ParsingError {
	if arg.Required {
		isZero := util.IsZero(reflect.ValueOf(output))

		if slice, isSlice := output.([]interface{}); arg.Greedy && (isZero || isSlice && len(slice) == 0) {
			return NewParsingError(arg, InsufficientArguments)
		}

		if isZero || output == nil {
			return NewParsingError(arg, RequiredArgumentMissing)
		}
	}

	return nil
}
