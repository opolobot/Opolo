package fun

import (
	"fmt"
	"math/rand"

	"github.com/zorbyte/whiskey/cmds"
)

func init() {
	cmdBldr := cmds.NewCommandBuilder()
	cmdBldr.Name("flip")
	cmdBldr.Description("Flips a coin! ~ heads or tails")
	cmdBldr.Use(flip)

	Category.AddCommand(cmdBldr.Build())
}

func flip(ctx *cmds.Context, next cmds.NextFunc) {
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
		next(fmt.Errorf("Failed to choose between heads or tails"))
	}

	ctx.Send(fmt.Sprintf(":coin: **It's __%v__!**", headsOrTails))
}
