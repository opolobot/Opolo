package core

import (
	"fmt"
	"math"
	"sort"
	"strings"

	"github.com/zorbyte/whiskey/arg"
	"github.com/zorbyte/whiskey/cmd"
	"github.com/zorbyte/whiskey/util"
	"github.com/zorbyte/whiskey/util/embed"
)

const helpEmoji = ":grey_question:"
const helpColour = 0xCCD6DD

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

func help(ctx *cmd.Context, _ cmd.NextFunc) {
	prefix := util.GetConfig().Prefix
	embed := embed.QuickEmbed(helpColour, helpEmoji, "Whiskey help", fmt.Sprintf("**prefix:** `%v`", prefix))
	embed.Fields = append(embed.Fields)
	if len(ctx.RawArgs) > 0 {
		callKeyArg := ctx.Args["cmd"]
		if callKeyArg != nil {
			_ = callKeyArg.(string)
			// TODO
		}
	}
}

func individualHelp(cmd *cmd.Command) {}

func helpOld(ctx *cmd.Context, next cmd.NextFunc) {
	var helpStrBldr strings.Builder
	helpMenuName := "**Whiskey help**"
	helpStrBldr.WriteString(helpMenuName)

	// TODO(@zorbyte): Dynamic prefix support.
	prefix := util.GetConfig().Prefix

	dividerLen := float64(len(helpMenuName)) - 4.0

	writeGap(&helpStrBldr)
	writeDivider(&helpStrBldr, dividerLen)
	writeGap(&helpStrBldr)

	// w!help [cmd] <- RawArgs[0]
	// get help info on a specific command
	if len(ctx.RawArgs) > 0 {
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

	writeDivider(&helpStrBldr, dividerLen)
	writeGap(&helpStrBldr)
	helpStrBldr.WriteString(fmt.Sprintf("**prefix  ~  **`%v`", prefix))
	writeGap(&helpStrBldr)
	writeDivider(&helpStrBldr, dividerLen)
	helpStrBldr.WriteString(fmt.Sprintf("\n*Whiskey (%v) by zorbyte and itjk.*", util.Version()))

	ctx.Send(helpStrBldr.String())
	next()
}

func writeGap(helpStrBldr *strings.Builder) {
	helpStrBldr.WriteString("\n\n")
}

func writeDivider(helpStrBldr *strings.Builder, dividerLen float64) {
	divider := "~~" + strings.Repeat("-", int(math.Floor(dividerLen*1.2))) + "~~"
	helpStrBldr.WriteString(divider)
}

func buildLookupHelp(helpStrBldr *strings.Builder, prefix string, cmnd *cmd.Command) {
	helpStrBldr.WriteString(cmnd.Category.DisplayName())
	helpStrBldr.WriteString(" **~** ")
	helpStrBldr.WriteString(fmt.Sprintf(cmnd.Help(), prefix))
}

func buildRegularHelp(helpStrBldr *strings.Builder, prefix string) {
	reg := cmd.GetRegistry()
	var keys []string
	for key := range reg.Categories {
		keys = append(keys, key)
	}

	sort.Strings(keys)
	_ = sort.Reverse(sort.StringSlice(keys))
	for _, key := range keys {
		cat := reg.Categories[key]
		helpStrBldr.WriteString(cat.DisplayName())
		helpStrBldr.WriteString(" **~** ")

		var cmdKeys []string
		for _, cmnd := range cat.Commands {
			cmdKeys = append(cmdKeys, "`"+cmnd.Name+"`")
		}

		sort.Strings(cmdKeys)
		helpStrBldr.WriteString(strings.Join(cmdKeys, ", "))
		writeGap(helpStrBldr)
	}
}
