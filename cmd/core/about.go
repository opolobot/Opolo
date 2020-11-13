package core

import (
	"fmt"

	"github.com/TeamWhiskey/whiskey/cmd"
)

const teamWhiskeyGithub string = "https://github.com/TeamWhiskey"

func init() {
	cmd := cmd.New()
	cmd.Name("about")
	cmd.Description("Tells you about Whiskey and TeamWhiskey")
	cmd.Use(about)

	Category.AddCommand(cmd.Command())
}

func about(ctx *cmd.Context, next cmd.NextFunc) {
	ctx.Send(fmt.Sprintf(
		"**Whiskey is a bot by :tumbler_glass: TeamWhiskey\n\nfind us on github ~ :octopus: %v**"+
			"\n\n**the ~~dipsos~~ team**\n\t\\~ zorbyte (Founder)\n\t\\~ MountainWhale\n\t\\~ FardinDaDev",
		teamWhiskeyGithub,
	))

	next()
}
