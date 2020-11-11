package fun

import "github.com/zorbyte/whiskey/cmds"

// Category for the fun commands.
var Category *cmds.CommandCategory

func init() {
	Category = cmds.NewCommandCategory("Fun", ":tada:")
}