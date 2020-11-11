package mod

import "github.com/TeamWhiskey/whiskey/cmds"

// Category for the moderation commands.
var Category *cmds.CommandCategory

func init() {
	Category = cmds.NewCommandCategory("Moderation", ":hammer:")
}