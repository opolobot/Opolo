package embeds

import "github.com/bwmarrin/discordgo"

const infoColour = 0x3669FA

// Info is a simple embed.
func Info(title, emoji, description string) *discordgo.MessageEmbed {
	embed := QuickEmbed(infoColour, emoji, title, description)

	return embed
}
