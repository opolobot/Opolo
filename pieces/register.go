package pieces

import (
	"github.com/bwmarrin/discordgo"
	"github.com/opolobot/opolo/ocl"
	"github.com/opolobot/opolo/pieces/cmds/core"
	"github.com/opolobot/opolo/pieces/cmds/fun"
	"github.com/opolobot/opolo/pieces/events"
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
}
