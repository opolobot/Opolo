package core

import (
	"fmt"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/opolobot/opolo/ocl"
	"github.com/opolobot/opolo/ocl/args"
	"github.com/opolobot/opolo/ocl/embeds"
	"github.com/opolobot/opolo/pieces/parsers"
	"github.com/opolobot/opolo/utils"
)

func init() {
	cmd := &ocl.Command{
		Name:        "help",
		Aliases:     []string{"h", "cmds", "commands"},
		Description: "Provides help for using the opolo.",
		Arguments:   []*args.Argument{args.Create("[cmd]", &parsers.String{})},
		Stack:       []ocl.Middleware{help},
	}

	Category.Add(cmd)
}

func help(ctx *ocl.Context, next ocl.Next) {
	callKey := ctx.Args["cmd"].(string)
	if callKey != "" {
		err := individualHelp(ctx, callKey)
		if err != nil {
			next(err)
		}

		return
	}

	regularHelp(ctx)
	next()
}

func individualHelp(ctx *ocl.Context, callKey string) error {
	cmd, err := ocl.GetRegistry().LookupCommand(callKey)
	if err != nil {
		return err
	}

	if cmd == nil {
		ctx.SendEmbed(embeds.Warn(fmt.Sprintf("Command `%v` not found", callKey), ""))
		return nil
	}

	// Adds embed.Fields space in front of the usage string if needed.
	usageStr := cmd.Usage()
	if usageStr != "" {
		usageStr = " " + usageStr
	}

	usage := fmt.Sprintf("`%v%v%v`\n", utils.StubPrefix(), cmd.Name, usageStr)

	embed := embeds.Info(embeds.Subtitle("Help", usage), "", cmd.Description)

	// If there are any aliases, add the field.
	if len(cmd.Aliases) > 0 {
		embed.Fields = append(embed.Fields, &discordgo.MessageEmbedField{
			Name:  ":scroll:  Aliases",
			Value: "`" + strings.Join(cmd.Aliases, "`, `") + "`",
		})
	}

	// Add the permissions field.
	embed.Fields = append(
		embed.Fields,
		&discordgo.MessageEmbedField{
			Name:  ":busts_in_silhouette:  Permission Level",
			Value: fmt.Sprint(cmd.Permission),
		},
		&discordgo.MessageEmbedField{
			Name:  ":books:  Category",
			Value: cmd.Category().Name,
		},
	)

	ctx.SendEmbed(embed)

	return nil
}

func regularHelp(ctx *ocl.Context) {
	prefix := utils.StubPrefix()

	embed := embeds.Info("Help", "", fmt.Sprintf("Your server prefix is `%v`", prefix))
	embed.Footer = &discordgo.MessageEmbedFooter{
		Text: fmt.Sprintf("Run %vhelp <command> for info about a command.", prefix),
	}

	var lastCategoryName string
	var field *discordgo.MessageEmbedField
	var counted int
	for _, cmd := range ocl.GetRegistry().Commands {
		cat := cmd.Category()
		if lastCategoryName != cat.Name {
			lastCategoryName = cat.Name
			field = &discordgo.MessageEmbedField{
				Name:  ":" + cat.Emoji + ":  " + cat.Name,
				Value: "",
			}

			counted = 0
			embed.Fields = append(embed.Fields, field)
		}

		if ctx.HasPermission(cmd.Permission) {
			field.Value += "`" + cmd.Name + "`"
			if counted < len(cat.Commands)-1 {
				field.Value += ", "
			}

			counted++
		}

		if counted == 0 {
			copy(embed.Fields[len(embed.Fields)-1:], embed.Fields[len(embed.Fields):])
			embed.Fields[len(embed.Fields)-1] = nil
			embed.Fields = embed.Fields[:len(embed.Fields)-1]
		}
	}

	ctx.SendEmbed(embed)

	return
}
