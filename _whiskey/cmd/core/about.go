package core

import (
	"fmt"

	"github.com/zorbyte/whiskey/cmd"
)

const teamWhiskeyGithub string = "https://github.com/zorbyte/Whiskey"

func init() {
	cmd := cmd.New()
	cmd.Name("about")
	cmd.Description("Tells you about Whiskey")
	cmd.Use(about)

	Category.AddCommand(cmd.Command())
}

func about(ctx *cmd.Context, next cmd.NextFunc) {
	ctx.Send(fmt.Sprintf(
		"**Whiskey; Just a beverage\n\nour github repo~ :octopus: %v**"+
			"\n\n**the ~~dipsos~~ devs**\n\t\\~ zorbyte (https://github.com/zorbyte)\n\t\\~ itjk (https://github.com/itjk",
		teamWhiskeyGithub,
	))

	next()
}
