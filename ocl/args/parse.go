package args

import (
	"reflect"
	"strings"

	"github.com/opolobot/Opolo/utils"
)

// ParsedArguments is a map of arguments using the name of the argument as the key.
type ParsedArguments = map[string]interface{}

// Parse parses the raw arguments into a map.
func Parse(args []*Argument, rawArgs []string) (ParsedArguments, error) {
	parsed := make(ParsedArguments)
	m := mapKeyedValues(rawArgs)
	// Index for the raw arguments
	i := 0
	for _, arg := range args {
		var output interface{}
		if i > len(rawArgs)-1 {
			output = arg.parser.ZeroVal()
		} else {
			if arg.greedy {
				output = make([]interface{}, 0)
			}

			if arg.keyed {
				out, err := handleArg(arg, m[arg.name])
				if err != nil {
					return nil, err
				}
				output = out
			}

			for j := 0; i < len(rawArgs) && !arg.keyed; i++ {
				raw := rawArgs[i]
				out, err := handleArg(arg, raw)
				if err != nil {
					return nil, err
				}

				if arg.greedy {
					// If the output of a greedy arg is a zero value, do not append it
					// and assume that the arguments have ended.
					if out == arg.parser.ZeroVal() {
						i--
						break
					}

					output = append(output.([]interface{}), out)
				} else {
					output = out
				}

				if !arg.greedy && j == 0 {
					break
				}

				j++
			}
		}

		err := validateArgOutput(arg, output)
		if err != nil {
			return nil, err
		}

		parsed[arg.name] = output
	}

	return parsed, nil
}

// handleArg handles an argument along with its raw one.
func handleArg(arg *Argument, raw string) (interface{}, error) {
	if arg.parser == nil {
		return raw, nil
	}

	output, err := arg.parser.Parse(raw)
	if err != nil {
		return output, err
	}

	return output, nil
}

// validateArgOutput validates the parsed data (output) from an argument.
// if any conditions fail, the appropriate error shall be thrown.
func validateArgOutput(arg *Argument, output interface{}) error {
	if arg.required {
		isZero := output == nil || utils.IsZero(reflect.ValueOf(output))

		if slice, isSlice := output.([]interface{}); arg.greedy && (isZero || isSlice && len(slice) == 0) {
			return NewParsingError(arg, InsufficientArguments)
		}

		if isZero || output == nil {
			return NewParsingError(arg, RequiredArgumentMissing)
		}
	}

	return nil
}

func mapKeyedValues(rawArgs []string) map[string]string {
	name := ""
	m := make(map[string]string)

	for _, arg := range rawArgs {
		if strings.Contains(arg, "=") {
			keyVal := strings.Split(arg, "=")
			name = keyVal[0]
			m[name] = keyVal[1]
		} else if name != "" {
			m[name] += " " + arg
		}
	}

	return m
}
