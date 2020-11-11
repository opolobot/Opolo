package mod

import "github.com/zorbyte/whiskey/cmds"

// Category for the moderation commands.
var Category *cmds.CommandCategory

func init() {
	Category = cmds.NewCommandCategory("Moderation", ":hammer:")
}