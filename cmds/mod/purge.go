package mod

import (
	"fmt"
	"math"
	"time"

	"github.com/TeamWhiskey/whiskey/args"
	"github.com/TeamWhiskey/whiskey/args/parsers"
	"github.com/TeamWhiskey/whiskey/cmds"
	"github.com/TeamWhiskey/whiskey/utils/mdlwr"
	"github.com/bwmarrin/discordgo"
)

func init() {
	cmdBldr := cmds.NewCommandBuilder()
	cmdBldr.Name("purge")
	cmdBldr.Description("Purges the desired amount of messages from the channel")
	cmdBldr.Use(mdlwr.PermCheck(discordgo.PermissionManageMessages, "ManageMessages"))
	cmdBldr.Use(checkValidAmnt)
	cmdBldr.Use(purge)
	cmdBldr.Args(&args.ArgumentCodec{
		Name:      "amnt",
		ExtraInfo: "<=500",
		Parser:    parsers.ParseIntArg,
		Required:  true,
	})

	Category.AddCommand(cmdBldr.Build())
}

func checkValidAmnt(ctx *cmds.Context, next cmds.NextFunc) {
	amnt := ctx.ArgsCtx.Args["amnt"].(int)

	if amnt > 500 {
		ctx.Send("Amount can not be greater than 500")
		return
	}

	next()
}

func purge(ctx *cmds.Context, next cmds.NextFunc) {
	amnt := ctx.ArgsCtx.Args["amnt"].(int)

	res, err := ctx.Prompt(fmt.Sprintf("Are you sure you wish to delete %v msgs?", amnt))
	if err != nil {
		next(err)
		return
	}

	if res == cmds.PromptTimeout {
		ctx.Send("Timed out.", 10*time.Second)
		next()
		return
	} else if res == cmds.PromptDeny {
		ctx.Send("Cancelled.", 10*time.Second)
		next()
		return
	}

	quotient := int(math.Floor(float64(amnt) / 100))
	iters := quotient
	rem := amnt % 100
	if rem > 0 {
		iters++
	}

	curMsgsIdx := 0
	for i := 0; i < iters; i++ {
		amntToDel := 0
		if i <= quotient && quotient != 0 {
			if err != nil {
				next(err)
				return
			}

			amntToDel = curMsgsIdx * 100
			curMsgsIdx++
		} else {
			amntToDel = int(rem)
		}

		msgs, err := ctx.Session.ChannelMessages(ctx.Msg.ChannelID, amntToDel, ctx.Msg.ID, "", "")
		if err != nil {
			next(err)
			return
		}

		var msgIDs []string
		for _, msg := range msgs {
			msgIDs = append(msgIDs, msg.ID)
		}

		go ctx.Session.ChannelMessagesBulkDelete(ctx.Msg.ChannelID, msgIDs[:amntToDel])
	}

	next()
}
