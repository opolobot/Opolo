package cmds

import (
	"fmt"
	"math/rand"

	"github.com/zorbyte/whiskey/lib"
)

var funCmds *cmdCategory

func init() {
	funCmds = &cmdCategory{
		Emoji:       ":tada:",
		Name:        "fun",
		DisplayName: "Fun",
	}

	funCmds.Cmds = append(
		funCmds.Cmds,
		&lib.Cmd{
			Runner:      flip,
			Aliases:     []string{"f"},
			Description: "Flips a coin",
		},
	)
}

// -- flip --

// flip flips a coin with results being either heads or tails
func flip(ctx *lib.Ctx) (string, error) {
	headsOrTails := (func() string {
		switch rand.Intn(2) {
		case 0:
			return "heads"
		case 1:
			return "tails"
		default:
			return ""
		}
	})()

	if headsOrTails == "" {
		return "", fmt.Errorf("Failed to choose between heads or tails")
	}

	return fmt.Sprintf(":coin: **It's __%v__!**", headsOrTails), nil
}
