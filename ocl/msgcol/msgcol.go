package msgcol

import "github.com/bwmarrin/discordgo"

// MessageCollector collects messages.
type MessageCollector struct {
	ID        uint16
	ChannelID string
	Msgs      chan *discordgo.Message
}
