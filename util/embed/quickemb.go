package embed

import "github.com/bwmarrin/discordgo"

// WhiskeyColour is the colour for whiskey embeds.
const WhiskeyColour = 0xB6801F

// QuickEmbed is a means to create an embed with the Whiskey Standard format.
func QuickEmbed(colour int, emoji, title, description string) *discordgo.MessageEmbed {
	embed := &discordgo.MessageEmbed{
		Title: emoji + " " + title,

		// Ha!
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
