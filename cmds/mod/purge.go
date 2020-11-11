package mod

import (
	"fmt"
	"math"
	"strings"
	"time"

	"github.com/zorbyte/whiskey/args"
	"github.com/zorbyte/whiskey/args/parsers"
	"github.com/zorbyte/whiskey/cmds"
)

func init() {
	cmdBldr := cmds.NewCommandBuilder()
	cmdBldr.Name("purge")
	cmdBldr.Description("Purges the desired amount of messages from the channel")
	cmdBldr.Use(purge)
	cmdBldr.Args(&args.ArgumentCodec{
		Name:      "amnt",
		ExtraInfo: "<=500",
		Parser:    parsers.ParseIntArg,
		Required:  true,
	})

	Category.AddCommand(cmdBldr.Build())
}

func purge(ctx *cmds.Context, next cmds.NextFunc) {
	amnt := ctx.ArgsCtx.Args["amnt"].(int)

	if amnt > 500 {
		ctx.Send("Amount can not be greater than 500")
		next()
		return
	}

	_, err := ctx.Send(fmt.Sprintf("Are you sure you wish to delete %v msgs? This action will cancel in 10 seconds. [y/N]", amnt))
	if err != nil {
		next(err)
		return
	}

	msgs, err := ctx.Collect(10 * time.Second)
	if err != nil {
		next(err)
		return
	}

	accepted := false
	denied := false
	for msg := range msgs {
		switch strings.ToLower(msg.Content) {
		case "y":
			accepted = true
			break
		case "n":
			denied = true
			break
		}
	}

	err = ctx.Delete()
	if err != nil {
		next(err)
		return
	}

	if !accepted {
		if denied {
			ctx.Send("Cancelling purge.")
			next()
			return
		}

		ctx.Send("Timed out.")
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
