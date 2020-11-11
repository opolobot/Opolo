package args

import (
	"github.com/bwmarrin/discordgo"
	"github.com/TeamWhiskey/whiskey/utils"
)

// ArgumentsContext is information about the arguments currently in use for this command.
type ArgumentsContext struct {
	Amnt    int
	RawArgs []string
	Args    map[string]interface{}
}

// NewArgumentsContext processes the raw arguments into the Args map and creates an arguments context instance.
func NewArgumentsContext(session *discordgo.Session, codecs []*ArgumentCodec, rawArgs []string) (*ArgumentsContext, error) {
	args := make(map[string]interface{})
	amntProcessed := 0
	finished := false
	for _, codec := range codecs {
		finished = amntProcessed > len(rawArgs)-1
		if finished {
			if codec.Required {
				if codec.Greedy {
					return nil, NewParsingError(codec.DisplayName(), insufficientArguments)
				}

				return nil, NewParsingError(codec.DisplayName(), requiredArgumentMissing)
			}
		} else {
			rawArg := rawArgs[amntProcessed]

			var parsed interface{}
			if codec.Parser == nil {
				parsed = rawArg
			} else {
				var err error
				parsed, err = codec.Parser(rawArg)
				if err != nil {
					if parsingError, ok := err.(*ParsingError); ok {
						return nil, parsingError
					}
	
					utils.SendError(session, err)
					return nil, NewParsingError(codec.DisplayName(), -1, err)
				}
			}
	
			args[codec.Name] = parsed
			amntProcessed++
		}
	}

	argsCtx := &ArgumentsContext{
		Amnt:    amntProcessed,
		RawArgs: rawArgs,
		Args:    args,
	}

	return argsCtx, nil
}
