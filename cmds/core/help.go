package core

import (
	"sort"
	"strings"

	"github.com/TeamWhiskey/whiskey/args"
	"github.com/TeamWhiskey/whiskey/cmds"
	"github.com/TeamWhiskey/whiskey/utils"
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
	var keys []string
	for key := range cmdUI.NameAndDisplayNames {
		keys = append(keys, key)
	}

	sort.Strings(keys)
	for _, key := range keys {
		displayName := cmdUI.NameAndDisplayNames[key]
		category := cmdUI.HelpOutputs[displayName]
		helpStrBldr.WriteString(displayName)
		helpStrBldr.WriteString("\n\n")

		var cmdKeys []string
		for key := range *category {
			cmdKeys = append(cmdKeys, key)
		}

		sort.Strings(cmdKeys)
		for _, cmdKey := range cmdKeys {
			cmdHelpStr := (*category)[cmdKey]
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
