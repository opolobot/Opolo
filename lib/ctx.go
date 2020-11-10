package lib

import (
	"errors"
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

	lastMsg *discordgo.Message

	// String used to call the command
	CmdCallKey string
	Args       []string

	StartTime time.Time

	// TODO(@zorbyte): Create a db, consider using gorm.
	DB interface{}
}

// Send sends a message to the channel the msg came from.
func (ctx *Ctx) Send(content string) (*discordgo.Message, error) {
	m, err := ctx.S.ChannelMessageSend(ctx.M.ChannelID, content)
	ctx.lastMsg = m
	return m, err
}

// Edit edits the last msg sent by the context.
func (ctx *Ctx) Edit(content string) (*discordgo.Message, error) {
	if ctx.lastMsg == nil {
		return nil, errors.New("Tried to edit a message that does not exist")
	}

	m, err := ctx.S.ChannelMessageEdit(ctx.lastMsg.ChannelID, ctx.lastMsg.ID, content)
	ctx.lastMsg = m
	return m, err
}

// Delete deletes the last recently sent message.
func (ctx *Ctx) Delete() error {
	if ctx.lastMsg == nil {
		return errors.New("Tried to delete a message that does not exist")
	}

	return ctx.S.ChannelMessageDelete(ctx.lastMsg.ChannelID, ctx.lastMsg.ID)
}

// Collect collects messages.
func (ctx *Ctx) Collect(time time.Duration, amnt ...int) (chan *discordgo.Message, error) {
	amntEmpty := len(amnt) == 0
	if time == 0 && amntEmpty {
		return nil, errors.New("ctx.Collect can not collect indefinitely if no amount is supplied")
	}

	var amntToPurge int
	if amntEmpty {
		amntToPurge = 0
	} else {
		amntToPurge = amnt[0]
	}

	col := ctx.W.Collect(ctx.M.ChannelID, time, amntToPurge)
	return col.Msgs, nil
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
	ctx.W.SendError(errTxt)

	log.Println("An error occurred while handling the command "+cmdName+":", err)
}
