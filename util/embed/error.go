package embed

import "github.com/bwmarrin/discordgo"

const errorEmoji = ":rotating_light:"
const errorColour = 0xDD2E44

// Error is an error embed.
// TODO(@zorbyte): Reduce repetition.
func Error(title, description string) *discordgo.MessageEmbed {
	embedTitle := "Error"
	if title != "" {
		embedTitle = Subtitle(embedTitle, title)
	}

	embed := QuickEmbed(errorColour, errorEmoji, embedTitle, description)
	embed.Footer = &discordgo.MessageEmbedFooter{
		Text: "The error has been reported",
	}

	return embed
}
