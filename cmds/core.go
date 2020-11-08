package cmds

import (
	"fmt"
	"strings"
	"time"

	"github.com/zorbyte/whiskey/lib"
)

var coreCmds []*lib.Cmd

const coreCmdsEmoji string = ":tumbler_glass:"

func init() {
	coreCmds = append(
		coreCmds,
		&lib.Cmd{
			Runner:      ping,
			Aliases:     []string{"p"},
			Description: "Tests the bot latency",
		},
		&lib.Cmd{
			Runner:      help,
			Usage:       "[cmd]",
			Aliases:     []string{"h"},
			Description: "List of the bot's commands and how to use them",
		},
	)
}

func ping(ctx *lib.Ctx) (string, error) {
	startTime := time.Now()
	m, err := ctx.Send(":ping_pong: Ping?")
	if err != nil {
		return "", err
	}

	ctx.S.ChannelMessageEdit(m.ChannelID, m.ID, fmt.Sprintf(":ping_pong: Pong! %v", time.Since(startTime)))
	return "", nil
}

func help(ctx *lib.Ctx) (string, error) {
	// TODO: Support command lookup.
	var helpStrBldr strings.Builder
	helpStrBldr.WriteString("**Whiskey (development) help**\n\n")
	for category, cmds := range helpOutputs {
		helpStrBldr.WriteString(category)
		helpStrBldr.WriteString("\n\n")
		for _, cmdHelpStr := range *cmds {
			helpStrBldr.WriteString("> **`")
			helpStrBldr.WriteString(ctx.W.Config.Prefix)
			helpStrBldr.WriteString(cmdHelpStr)
			helpStrBldr.WriteString("\n\n")
		}
	}

	return helpStrBldr.String(), nil
}
