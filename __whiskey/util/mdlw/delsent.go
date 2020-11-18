package mdlw

import (
	"github.com/bwmarrin/discordgo"
	"github.com/zorbyte/whiskey/cmd"
	"github.com/zorbyte/whiskey/util"
)

// DeleteSent deletes the message sent by the author if possible.
func DeleteSent(ctx *cmd.Context, next cmd.NextFunc) {
	hasPerm, err := util.HasPermission(ctx.Session.State, ctx.Msg.ChannelID, discordgo.PermissionManageMessages)
	if err != nil {
		next(err)
		return
	}

	if hasPerm {
		err = ctx.Delete(ctx.Msg)
		if err != nil {
			next(err)
		}

		return
	}

	next()
}
