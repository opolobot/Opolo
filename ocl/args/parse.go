package args

import (
	"reflect"
	"strings"

	"github.com/opolobot/Opolo/pieces/parsers"
	"github.com/opolobot/Opolo/utils"
)

// ParsedArguments is a map of arguments using the name of the argument as the key.
type ParsedArguments = map[string]interface{}

// Parse parses the raw arguments into a map.
func Parse(args []*Argument, rawArgs []string) (ParsedArguments, error) {
	parsed := make(ParsedArguments)

	// Index for the raw arguments
	rawIdx := 0
	for argIdx, arg := range args {
		var output interface{}
		if rawIdx > len(rawArgs)-1 {
			output = arg.parser.ZeroVal()
		} else {
			if arg.greedy {
				output = make([]interface{}, 0)
			}

			for greedyIdx := 0; rawIdx < len(rawArgs); rawIdx++ {
				raw := rawArgs[rawIdx]
				out, err := handleArg(arg, raw)
				if err != nil {
					return nil, err
				}

				if arg.greedy {
					if out == arg.parser.ZeroVal() {
						break
					}

					if strings.Contains(raw, "=") && argIdx+1 < len(args) {
						if _, ok := args[argIdx+1].parser.(*parsers.KeyValue); ok {
							break
						}
					}

					output = append(output.([]interface{}), out)
				} else {
					output = out
				}

				if !arg.greedy && greedyIdx == 0 {
					break
				}

				greedyIdx++
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
