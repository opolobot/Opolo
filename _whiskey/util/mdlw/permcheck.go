package mdlw

import (
	"fmt"

	"github.com/zorbyte/whiskey/cmd"
	"github.com/zorbyte/whiskey/util"
)

// PermCheck is a middleware for checking permissions.
func PermCheck(permMask int, name string) cmd.Middleware {
	return func(ctx *cmd.Context, next cmd.NextFunc) {
		hasPerm, err := util.HasPermission(ctx.Session.State, ctx.Msg.ChannelID, permMask)
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
