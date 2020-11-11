package core

import (
	"strings"

	"github.com/zorbyte/whiskey/args"
	"github.com/zorbyte/whiskey/cmds"
	"github.com/zorbyte/whiskey/utils"
)

var config *utils.Config
var cmdUI *cmds.CommandUI

func init() {
	config = utils.GetConfig()
	cmdUI = cmds.GetCommandUI()

	cmdBldr := cmds.NewCommandBuilder()
	cmdBldr.Name("help")
	cmdBldr.Aliases("h")
	cmdBldr.Description("List of the bot's commands and how to use them")
	cmdBldr.Use(help)
	cmdBldr.Args(&args.ArgumentCodec{
		Name: "cmd",
	})

	Category.AddCommand(cmdBldr.Build())
}

func help(ctx *cmds.Context, next cmds.NextFunc) {
	var helpStrBldr strings.Builder
	helpStrBldr.WriteString("**Whiskey (development) help**\n\n")

	// w!help [cmd] <- Args[]
	// get help info on a specific command
	if ctx.ArgsCtx.Amnt > 0 {
		cmdCallKeyArg := ctx.ArgsCtx.Args["cmd"]
		if cmdCallKeyArg != nil {
			cmdCallKey := cmdCallKeyArg.(string)
			cmd, err := cmdUI.LookupCommand(cmdCallKey)
			if err != nil {
				next(err)
			}

			if cmd != nil {
				buildLookupHelp(&helpStrBldr, cmd)
			} else {
				buildRegularHelp(&helpStrBldr)
			}
		}
	} else {
		buildRegularHelp(&helpStrBldr)
	}

	ctx.Send(helpStrBldr.String())
	next()
}

func buildLookupHelp(helpStrBldr *strings.Builder, cmd *cmds.Command) {
	cmdCat, ok := cmdUI.HelpOutputs[cmd.Category.DisplayName()]
	if ok {
		helpStrBldr.WriteString(cmd.Category.DisplayName())
		helpStrBldr.WriteString("\n\n")
		cmdHelpStr, ok := (*cmdCat)[cmd.Name]
		if ok {
			finaliseCmdHelpStr(cmdHelpStr, config.Prefix, helpStrBldr)
		}
	}
}

func buildRegularHelp(helpStrBldr *strings.Builder) {
	for category, cmds := range cmdUI.HelpOutputs {
		helpStrBldr.WriteString(category)
		helpStrBldr.WriteString("\n\n")
		for _, cmdHelpStr := range *cmds {
			finaliseCmdHelpStr(cmdHelpStr, config.Prefix, helpStrBldr)
		}
	}
}

func finaliseCmdHelpStr(cmdHelpStr string, prefix string, strBldr *strings.Builder) {
	strBldr.WriteString("> **`")
	strBldr.WriteString(prefix)
	strBldr.WriteString(cmdHelpStr)
	strBldr.WriteString("\n\n")
}
