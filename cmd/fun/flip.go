package fun

import (
	"fmt"
	"math/rand"

	"github.com/TeamWhiskey/whiskey/cmd"
)

func init() {
	cmd := cmd.New()
	cmd.Name("flip")
	cmd.Description("Flips a coin! ~ heads or tails")
	cmd.Use(flip)

	Category.AddCommand(cmd.Command())
}

func flip(ctx *cmd.Context, next cmd.NextFunc) {
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
