package lib

import (
	"log"
	"time"

	"github.com/bwmarrin/discordgo"
)

// Ctx is the message context used for command execution
type Ctx struct {
	W *Whiskey
	S *discordgo.Session
	M *discordgo.Message
	C *Cmd

	// String used to call the command
	CmdCallKey string
	Args       []string

	StartTime time.Time

	// TODO(@zorbyte): Create a db, consider using gorm.
	DB interface{}
}

// Send sends a message to the channel the msg came from.
func (ctx *Ctx) Send(content string) (*discordgo.Message, error) {
	return ctx.S.ChannelMessageSend(ctx.M.ChannelID, content)
}

// SendError reports an error to the err channel and to the user
func (ctx *Ctx) SendError(err error) {
	cmdName := (func() string {
		if ctx.C.Name != "" {
			return ctx.C.Name
		}

		return "N/A"
	})()

	errTxt := ":rotating_light: An error occurred while handling the command `" + cmdName + "`:\n```" + err.Error() + "```"
	ctx.S.ChannelMessageSend(ctx.M.ChannelID, errTxt+"\n\nThe error has been reported")
	if ctx.W.Config.LogChannel != "" {
		ctx.S.ChannelMessageSend(ctx.W.Config.LogChannel, errTxt)
	}

	log.Println("An error occurred while handling the command "+cmdName+":", err)
}
