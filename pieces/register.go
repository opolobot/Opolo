package pieces

import (
	"github.com/bwmarrin/discordgo"
	"github.com/opolobot/Opolo/ocl"
	"github.com/opolobot/Opolo/pieces/cmds/core"
	"github.com/opolobot/Opolo/pieces/cmds/fun"
	"github.com/opolobot/Opolo/pieces/cmds/mod"
	"github.com/opolobot/Opolo/pieces/events"
)

// TODO(@zorbyte): Makes this more dynamic instead of registrering the individual pieces.

// RegisterHandlers registers event handlers.
func RegisterHandlers(session *discordgo.Session) {
	session.AddHandler(events.Ready)
	session.AddHandler(events.MessageCreate)
}

// RegisterCommandCategories registers command categories.
func RegisterCommandCategories() {
	reg := ocl.GetRegistry()
	defer reg.Populate()

	reg.AddCategory(core.Category)
	reg.AddCategory(fun.Category)
	reg.AddCategory(mod.Category)
}
