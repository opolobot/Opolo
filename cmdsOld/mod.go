package cmdsOld

import (
	"fmt"
	"math"
	"strconv"
	"strings"
	"time"

	"github.com/zorbyte/whiskey/lib"
)

var modCmds *cmdCategory

func init() {
	modCmds = &cmdCategory{
		Emoji:       ":hammer:",
		Name:        "moderation",
		DisplayName: "Moderation",
	}

	modCmds.Cmds = append(
		modCmds.Cmds,
		&lib.Cmd{
			Runner:      purge,
			Description: "Purges the desired amount of messages from the chat",
			Usage:       "<amount (<=500)>",
		},
	)
}

func purge(ctx *lib.Ctx) (string, error) {
	if ctx.Args[0] == "" {
		return "Insufficient arguments! Please supply an amount", nil
	}

	amnt, err := strconv.Atoi(ctx.Args[0])
	if err != nil {
		return fmt.Sprintf("Input %v is not a valid number", amnt), nil
	}

	if amnt > 500 {
		return "Amount can not be greater than 500", nil
	}

	_, err = ctx.Send(fmt.Sprintf("Are you sure you wish to delete %v msgs? This action will cancel in 10 seconds. [y/N]", amnt))
	if err != nil {
		return "", err
	}

	msgs, err := ctx.Collect(10 * time.Second)
	if err != nil {
		return "", err
	}

	go (func() {
		time.Sleep(5 * time.Second)
		ctx.Edit("")
	})()

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
		return "", err
	}

	if !accepted {
		if denied {
			return "Cancelling purge.", nil
		}

		return "Timed out.", nil
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
				return "", err
			}

			amntToDel = curMsgsIdx * 100
			curMsgsIdx++
		} else {
			amntToDel = int(rem)
		}

		msgs, err := ctx.S.ChannelMessages(ctx.M.ChannelID, amntToDel, ctx.M.ID, "", "")
		if err != nil {
			return "", err
		}

		var msgIDs []string
		for _, msg := range msgs {
			msgIDs = append(msgIDs, msg.ID)
		}

		go ctx.S.ChannelMessagesBulkDelete(ctx.M.ChannelID, msgIDs[:amntToDel])
	}

	return "", nil
}
