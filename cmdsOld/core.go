package cmdsOld

import (
	"fmt"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/zorbyte/whiskey/lib"
)

var coreCmds *cmdCategory

const teamWhiskeyGithub string = "https://github.com/TeamWhiskey"

func init() {
	coreCmds = &cmdCategory{
		Emoji:       ":tumbler_glass:",
		Name:        "core",
		DisplayName: "Whiskey core",
	}

	coreCmds.Cmds = append(
		coreCmds.Cmds,
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
		&lib.Cmd{
			Runner:      about,
			Description: "Tell's you about Whiskey and TeamWhiskey",
		},
	)
}

// -- ping --

func ping(ctx *lib.Ctx) (string, error) {
	executionTime := time.Since(ctx.StartTime)
	startTime := time.Now()
	m, err := ctx.Send(":ping_pong: Ping?")
	if err != nil {
		return "", err
	}

	messageSentTime, err := discordgo.SnowflakeTimestamp(m.ID)
	if err != nil {
		return "", err
	}

	ctx.Edit(fmt.Sprintf(
		"***:ping_pong:  ~Pong!***\n"+
			"\n> __**latency**__        **~**   :arrows_counterclockwise: %v"+
			"\n> __**exec. time**__   **~**   :stopwatch: %v",
		messageSentTime.Sub(startTime).Round(time.Millisecond),
		executionTime.Round(time.Microsecond),
	))

	return "", nil
}

// -- help --

func help(ctx *lib.Ctx) (string, error) {
	var helpStrBldr strings.Builder
	helpStrBldr.WriteString("**Whiskey (development) help**\n\n")

	// w!help [cmd] <- Args[0]
	// get help info on a specific command
	if len(ctx.Args) > 0 {
		cmdCallKey := ctx.Args[0]
		cmd := ctx.W.FindCmd(cmdCallKey)
		if cmd != nil {
			cmdCat, ok := helpOutputs[cmd.Category]
			if ok {
				helpStrBldr.WriteString(cmd.Category)
				helpStrBldr.WriteString("\n\n")
				cmdHelpStr, ok := (*cmdCat)[cmd.Name]
				if ok {
					finaliseCmdHelpStr(cmdHelpStr, ctx.W.Config.Prefix, &helpStrBldr)
					return helpStrBldr.String(), nil
				}
			}
		}
	}

	for category, cmds := range helpOutputs {
		helpStrBldr.WriteString(category)
		helpStrBldr.WriteString("\n\n")
		for _, cmdHelpStr := range *cmds {
			finaliseCmdHelpStr(cmdHelpStr, ctx.W.Config.Prefix, &helpStrBldr)
		}
	}

	return helpStrBldr.String(), nil
}

func finaliseCmdHelpStr(cmdHelpStr string, prefix string, strBldr *strings.Builder) {
	strBldr.WriteString("> **`")
	strBldr.WriteString(prefix)
	strBldr.WriteString(cmdHelpStr)
	strBldr.WriteString("\n\n")
}

// -- about --

func about(ctx *lib.Ctx) (string, error) {
	return fmt.Sprintf(
		"**Whiskey is a bot by :tumbler_glass: TeamWhiskey\n\nfind us on github ~ :octopus: %v**"+
			"\n\n**the ~~dipsos~~ team**\n\t\\~ zorbyte (Founder)\n\t\\~ MountainWhale\n\t\\~ FardinDaDev",
		teamWhiskeyGithub,
	), nil
}
