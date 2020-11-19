package events

import (
	"github.com/bwmarrin/discordgo"
	"github.com/opolobot/opolo/ocl"
	"github.com/opolobot/opolo/ocl/msgcol"
)

// MessageCreate event.
func MessageCreate(session *discordgo.Session, msg *discordgo.MessageCreate) {
	// Ignore messages by bots.
	if msg.Author.Bot {
		return
	}

	// This NextFunc instantiates the middleware chain.
	next := ocl.Dispatch(session, msg.Message)
	if next == nil {
		msgcol.GetCollectionManager().Dispatch(msg.Message)
		return
	}

	next()
}
