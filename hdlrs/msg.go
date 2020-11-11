package hdlrs

import (
	"github.com/bwmarrin/discordgo"
	"github.com/zorbyte/whiskey/cmds"
	"github.com/zorbyte/whiskey/utils/msgcol"
)

// MessageCreate event.
func MessageCreate(session *discordgo.Session, msg *discordgo.MessageCreate) {
	// Ignore messages by the bot itself.
	if msg.Author.ID == session.State.User.ID {
		return
	}

	cmdUI := cmds.GetCommandUI()
	next := cmdUI.Dispatch(session, msg.Message)
	if next == nil {
		msgcol.GetCollectionManager().Dispatch(msg.Message)
		return
	}

	next()
}
