package hdlr

import (
	"log"

	"github.com/bwmarrin/discordgo"
	"github.com/TeamWhiskey/whiskey/util"
)

// Ready event.
func Ready(s *discordgo.Session, m *discordgo.Ready) {
	log.Printf("Logged in as %v#%v (%v)", s.State.User.Username, s.State.User.Discriminator, s.State.User.ID)
	config := util.GetConfig()
	s.UpdateStatus(1, config.Status)
}
