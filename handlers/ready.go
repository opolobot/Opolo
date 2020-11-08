package handlers

import (
	"log"

	"github.com/bwmarrin/discordgo"
	"github.com/zorbyte/whiskey/lib"
)

// Ready event.
func Ready(w *lib.Whiskey) func(s *discordgo.Session, m *discordgo.Ready) {
	return func(s *discordgo.Session, m *discordgo.Ready) {
		log.Printf("Logged in as %v#%v (%v)", s.State.User.Username, s.State.User.Discriminator, s.State.User.ID)
		s.UpdateStatus(1, w.Config.Status)
	}
}
