package hdlr

import (
	"time"

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

	// <3 from zorbyte
	if msg.Content == ";babylon_exists_don't_tell_me_otherwise" {
		easterEggSponsoredByZorbyte(session, msg)
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

// easterEggSponsoredByZorbyte is an easter egg triggered by a special message content string.
// Please note that any religious references are not a representation of the views of the developers
// and are rather a reference to a running joke amongst the developers.
func easterEggSponsoredByZorbyte(session *discordgo.Session, msg *discordgo.MessageCreate) {
	m, _ := session.ChannelMessageSend(msg.ChannelID, "Where are we talking about?")
	go (func() {
		time.Sleep(1 * time.Second)
		session.ChannelMessageEdit(m.ChannelID, m.ID, "Oh wait, it's gone:tm:")
	})()
}
