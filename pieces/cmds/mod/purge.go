package mod

import (
	"fmt"
	"math"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/opolobot/Opolo/ocl"
	"github.com/opolobot/Opolo/ocl/args"
	"github.com/opolobot/Opolo/pieces/middleware"
	"github.com/opolobot/Opolo/pieces/parsers"
)

func init() {
	cmd := ocl.New()
	cmd.Name("purge")
	cmd.Description("Purges the desired amount of messages from the channel")
	cmd.Args(args.New("<amnt>", &parsers.Int{}, "0<amnt<=500"))
	cmd.Use(middleware.PermCheck(discordgo.PermissionManageMessages, "ManageMessages"))
	cmd.Use(middleware.DeleteSent)
	cmd.Use(checkValidAmnt)
	cmd.Use(purge)

	Category.Add(cmd)
}

func checkValidAmnt(ctx *ocl.Context, next ocl.Next) {
	amnt := ctx.Args["amnt"].(int)

	if amnt > 500 {
		ctx.Send("Amount can not be greater than 500")
		return
	}

	next()
}

func purge(ctx *ocl.Context, next ocl.Next) {
	amnt := ctx.Args["amnt"].(int)

	res, err := ctx.Prompt(fmt.Sprintf("Are you sure you wish to delete %v msgs?", amnt))
	if err != nil {
		next(err)
		return
	}

	if res == ocl.PromptTimeout {
		ctx.Send("Timed out.", 10*time.Second)
		next()
		return
	} else if res == ocl.PromptDeny {
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
