package cmds

import (
	"errors"
	"log"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/zorbyte/whiskey/lib"
)

// Context is the message context used for command execution
type Context struct {
	Whiskey *lib.Whiskey
	Session *discordgo.Session
	Msg *discordgo.Message
	Cmd *Command

	lastMsg *discordgo.Message

	// String used to call the command
	CmdCallKey string
	Args       []string

	StartTime time.Time

	// TODO(@zorbyte): Create a db, consider using gorm.
	DB interface{}
}

// Send sends a message to the channel the msg came from.
func (ctx *Context) Send(content string) (*discordgo.Message, error) {
	m, err := ctx.Session.ChannelMessageSend(ctx.Msg.ChannelID, content)
	ctx.lastMsg = m
	return m, err
}

// Edit edits the last msg sent by the context.
func (ctx *Context) Edit(content string) (*discordgo.Message, error) {
	if ctx.lastMsg == nil {
		return nil, errors.New("Tried to edit a message that does not exist")
	}

	m, err := ctx.Session.ChannelMessageEdit(ctx.lastMsg.ChannelID, ctx.lastMsg.ID, content)
	ctx.lastMsg = m
	return m, err
}

// Delete deletes the last recently sent message.
func (ctx *Context) Delete() error {
	if ctx.lastMsg == nil {
		return errors.New("Tried to delete a message that does not exist")
	}

	return ctx.Session.ChannelMessageDelete(ctx.lastMsg.ChannelID, ctx.lastMsg.ID)
}

// Collect collects messages.
func (ctx *Context) Collect(time time.Duration, amnt ...int) (chan *discordgo.Message, error) {
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

	col := ctx.Whiskey.Collect(ctx.Msg.ChannelID, time, amntToPurge)
	return col.Msgs, nil
}

// SendError reports an error to the err channel and to the user
func (ctx *Context) SendError(err error) {
	cmdName := (func() string {
		if ctx.Cmd.name != "" {
			return ctx.Cmd.name
		}

		return "N/A"
	})()

	errTxt := ":rotating_light: An error occurred while handling the command `" + cmdName + "`:\n```" + err.Error() + "```"
	ctx.Session.ChannelMessageSend(ctx.Msg.ChannelID, errTxt+"\n\nThe error has been reported")
	ctx.Whiskey.SendError(errors.New(errTxt))

	log.Println("An error occurred while handling the command "+cmdName+":", err)
}
