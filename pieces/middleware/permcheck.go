package middleware

import (
	"fmt"

	"github.com/opolobot/Opolo/ocl"
	"github.com/opolobot/Opolo/utils"
)

// PermCheck is a middleware for checking permissions.
func PermCheck(permMask int, name string) ocl.Middleware {
	return func(ctx *ocl.Context, next ocl.Next) {
		hasPerm, err := utils.HasPermission(ctx.Session.State, ctx.Msg.ChannelID, permMask)
		if err != nil {
			next(err)
			return
		}

		if hasPerm {
			next()
		} else {
			ctx.Send(fmt.Sprintf("Whiskey does not have the required permission `%v` for this command.", name))
		}
	}
}
