package core

import (
	"fmt"

	"github.com/zorbyte/whiskey/cmds"
)

const teamWhiskeyGithub string = "https://github.com/TeamWhiskey"

func init() {
	cmdBldr := cmds.NewCommandBuilder()
	cmdBldr.Name("about")
	cmdBldr.Description("Tells you about Whiskey and TeamWhiskey")
	cmdBldr.Use(about)

	Category.AddCommand(cmdBldr.Build())
}

func about(ctx *cmds.Context, next cmds.NextFunc) {
	ctx.Send(fmt.Sprintf(
		"**Whiskey is a bot by :tumbler_glass: TeamWhiskey\n\nfind us on github ~ :octopus: %v**"+
			"\n\n**the ~~dipsos~~ team**\n\t\\~ zorbyte (Founder)\n\t\\~ MountainWhale\n\t\\~ FardinDaDev",
		teamWhiskeyGithub,
	))

	next()
}
