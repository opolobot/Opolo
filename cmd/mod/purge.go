package mod

import (
	"fmt"
	"math"
	"time"

	"github.com/TeamWhiskey/whiskey/arg"
	"github.com/TeamWhiskey/whiskey/arg/parsers"
	"github.com/TeamWhiskey/whiskey/cmd"
	"github.com/TeamWhiskey/whiskey/util/mdlw"
	"github.com/bwmarrin/discordgo"
)

func init() {
	cmd := cmd.New()
	cmd.Name("purge")
	cmd.Description("Purges the desired amount of messages from the channel")
	cmd.Use(mdlw.PermCheck(discordgo.PermissionManageMessages, "ManageMessages"))
	cmd.Use(checkValidAmnt)
	cmd.Use(purge)
	cmd.Arg(&arg.Argument{
		Name:        "amnt",
		Constraints: "0<a<=500",
		Parser:      parsers.ParseInt,
		Required:    true,
	})

	Category.AddCommand(cmd.Command())
}

func checkValidAmnt(ctx *cmd.Context, next cmd.NextFunc) {
	amnt := ctx.Args["amnt"].(int)

	if amnt > 500 {
		ctx.Send("Amount can not be greater than 500")
		return
	}

	next()
}

func purge(ctx *cmd.Context, next cmd.NextFunc) {
	amnt := ctx.Args["amnt"].(int)

	res, err := ctx.Prompt(fmt.Sprintf("Are you sure you wish to delete %v msgs?", amnt))
	if err != nil {
		next(err)
		return
	}

	if res == cmd.PromptTimeout {
		ctx.Send("Timed out.", 10*time.Second)
		next()
		return
	} else if res == cmd.PromptDeny {
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
