package core

import (
	"sort"
	"strings"

	"github.com/TeamWhiskey/whiskey/arg"
	"github.com/TeamWhiskey/whiskey/cmd"
	"github.com/TeamWhiskey/whiskey/util"
)

func init() {
	cmd := cmd.New()
	cmd.Name("help")
	cmd.Aliases("h")
	cmd.Description("List of the bot's commands and how to use them")
	cmd.Use(help)
	cmd.Arg(&arg.Argument{
		Name: "cmd",
	})

	Category.AddCommand(cmd.Command())
}

func help(ctx *cmd.Context, next cmd.NextFunc) {
	var helpStrBldr strings.Builder
	helpStrBldr.WriteString("**Whiskey (development) help**\n\n")

	// TODO(@zorbyte): Dynamic prefix support.
	prefix := util.GetConfig().Prefix

	// w!help [cmd] <- Args[]
	// get help info on a specific command
	if len(ctx.Args) > 0 {
		callKeyArg := ctx.Args["cmd"]
		if callKeyArg != nil {
			callKey := callKeyArg.(string)

			reg := cmd.GetRegistry()
			cmnd, err := reg.LookupCommand(callKey)
			if err != nil {
				next(err)
				return
			}

			if cmnd != nil {
				buildLookupHelp(&helpStrBldr, prefix, cmnd)
			} else {
				buildRegularHelp(&helpStrBldr, prefix)
			}
		}
	} else {
		buildRegularHelp(&helpStrBldr, prefix)
	}

	ctx.Send(helpStrBldr.String())
	next()
}

func buildLookupHelp(helpStrBldr *strings.Builder, prefix string, cmnd *cmd.Command) {
	helpStrBldr.WriteString(cmnd.Category.DisplayName())
	helpStrBldr.WriteString("\n\n")
	finaliseCmdHelpStr(cmnd.Help(), prefix, helpStrBldr)
}

func buildRegularHelp(helpStrBldr *strings.Builder, prefix string) {
	reg := cmd.GetRegistry()
	var keys []string
	for key := range reg.Categories {
		keys = append(keys, key)
	}

	sort.Strings(keys)
	for _, key := range keys {
		cat := reg.Categories[key]
		helpStrBldr.WriteString(cat.DisplayName())
		helpStrBldr.WriteString("\n\n")

		var cmdKeys []string
		var cmdKeysToIdx map[string]int
		for idx, cmnd := range cat.Commands {
			cmdKeys = append(cmdKeys, cmnd.Name)
			cmdKeysToIdx[cmnd.Name] = idx
		}

		sort.Strings(cmdKeys)
		for _, cmdKey := range cmdKeys {
			idx := cmdKeysToIdx[cmdKey]
			cmnd := cat.Commands[idx]
			finaliseCmdHelpStr(cmnd.Help(), prefix, helpStrBldr)
		}
	}
}

func finaliseCmdHelpStr(cmdHelpStr string, prefix string, strBldr *strings.Builder) {
	strBldr.WriteString("> **`")
	strBldr.WriteString(prefix)
	strBldr.WriteString(cmdHelpStr)
	strBldr.WriteString("\n\n")
}
