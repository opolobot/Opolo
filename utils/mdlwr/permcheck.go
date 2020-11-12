package mdlwr

import (
	"fmt"

	"github.com/TeamWhiskey/whiskey/cmds"
	"github.com/TeamWhiskey/whiskey/utils"
)

// PermCheck is a middleware for checking permissions.
func PermCheck(permMask int, name string) cmds.CommandMiddleware {
	return func(ctx *cmds.Context, next cmds.NextFunc) {
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
