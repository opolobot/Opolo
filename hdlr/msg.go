package hdlr

import (
	"github.com/bwmarrin/discordgo"
	"github.com/zorbyte/whiskey/cmd"
	"github.com/zorbyte/whiskey/util/msgcol"
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
