package hdlrs

import (
	"log"

	"github.com/bwmarrin/discordgo"
	"github.com/TeamWhiskey/whiskey/utils"
)

// Ready event.
func Ready(s *discordgo.Session, m *discordgo.Ready) {
	log.Printf("Logged in as %v#%v (%v)", s.State.User.Username, s.State.User.Discriminator, s.State.User.ID)
	config := utils.GetConfig()
	s.UpdateStatus(1, config.Status)
}
