package embeds

import "github.com/bwmarrin/discordgo"

// QuickEmbed is a means to create an embed with the Whiskey Standard format.
func QuickEmbed(colour int, emoji, title, description string) *discordgo.MessageEmbed {
	var formattedEmoji string
	if emoji != "" {
		formattedEmoji = ":" + emoji + ":   "
	}

	embed := &discordgo.MessageEmbed{
		Title: formattedEmoji + title,

		// Ha! en-GB shall prevail!
		Color: colour,
	}

	if description != "" {
		embed.Description = description
	}

	return embed
}

// Subtitle is a string concatenation utility for formatting embed titles.
func Subtitle(mainTitle, subtitle string) string {
	return mainTitle + " ~ " + subtitle
}
