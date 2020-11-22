package fun

import (
	"fmt"
	"math/rand"

	"github.com/opolobot/Opolo/ocl"
	"github.com/opolobot/Opolo/ocl/embeds"
)

func init() {
	cmd := ocl.New()
	cmd.Name("flip")
	cmd.Aliases("coin", "coinflip")
	cmd.Description(embeds.Subtitle("Flips a coin", "heads or tails"))
	cmd.Use(flip)

	Category.Add(cmd)
}

func flip(ctx *ocl.Context, next ocl.Next) {
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
		next(fmt.Errorf("failed to choose between heads or tails"))
		return
	}

	ctx.SendEmbed(embeds.Info(fmt.Sprintf("It's __%v__!", headsOrTails), "coin", ""))
}
