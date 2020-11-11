package cmds

import "github.com/zorbyte/whiskey/lib"

// ArgumentsContext is information about the arguments currently in use for this command.
type ArgumentsContext struct {
	Amnt    int
	RawArgs []string
	Args    map[string]interface{}
}

// NewArgumentsContext processes the raw arguments into the Args map and creates an arguments context instance.
func NewArgumentsContext(whiskey *lib.Whiskey, codecs []*ArgumentCodec, rawArgs []string) (*ArgumentsContext, error) {
	args := make(map[string]interface{})
	amntProcessed := 0
	finished := false
	for _, codec := range codecs {
		finished = amntProcessed > len(rawArgs)-1
		if finished {
			if codec.required {
				if codec.greedy {
					return nil, NewParsingError(codec.DisplayName(), insufficientArguments)
				}

				return nil, NewParsingError(codec.DisplayName(), requiredArgumentMissing)
			}
		}

		rawArg := rawArgs[amntProcessed]
		parsed, err := codec.parser(rawArg)
		if err != nil {
			if parsingError, ok := err.(*parsingError); ok {
				return nil, parsingError
			}

			whiskey.SendError(err)
			return nil, NewParsingError(codec.DisplayName(), -1, err)
		}

		args[codec.name] = parsed
		amntProcessed++
	}

	argsCtx := &ArgumentsContext{
		Amnt:    amntProcessed,
		RawArgs: rawArgs,
		Args:    args,
	}

	return argsCtx, nil
}
