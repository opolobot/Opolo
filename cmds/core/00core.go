package core

import "github.com/TeamWhiskey/whiskey/cmds"

// Category for the core commands.
var Category *cmds.CommandCategory

func init() {
	Category = cmds.NewCommandCategory("Whiskey core", ":tumbler_glass:")
}
