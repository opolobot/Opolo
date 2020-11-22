package middleware

import (
	"github.com/bwmarrin/discordgo"
	"github.com/opolobot/Opolo/ocl"
	"github.com/opolobot/Opolo/utils"
)

// DeleteSent deletes the message sent by the author if possible.
func DeleteSent(ctx *ocl.Context, next ocl.Next) {
	hasPerm, err := utils.HasPermission(ctx.Session.State, ctx.Msg.ChannelID, discordgo.PermissionManageMessages)
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
