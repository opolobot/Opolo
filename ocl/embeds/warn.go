package embeds

import "github.com/bwmarrin/discordgo"

const warningEmoji = "warning"
const warningColour = 0xFFCC4D


// Warn creates a warning embed.
func Warn(title, description string) *discordgo.MessageEmbed {
	embedTitle := "Warning"
	if title != "" {
		embedTitle = Subtitle(embedTitle, title)
	}

	return QuickEmbed(warningColour, warningEmoji, embedTitle, description)
}
