package hdlr

import (
	"github.com/TeamWhiskey/whiskey/cmd"
	"github.com/TeamWhiskey/whiskey/util/msgcol"
	"github.com/bwmarrin/discordgo"
)

// MessageCreate event.
func MessageCreate(session *discordgo.Session, msg *discordgo.MessageCreate) {
	// Ignore messages by bots.
	if msg.Author.Bot {
		return
	}

	// This NextFunc instantiates the middleware chain.
	next := cmd.Dispatch(session, msg.Message)
	if next == nil {
		msgcol.GetCollectionManager().Dispatch(msg.Message)
		return
	}

	next()
}
