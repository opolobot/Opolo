package ocl

import (
	"errors"
	"log"
	"runtime/debug"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/opolobot/opolo/ocl/embeds"
	"github.com/opolobot/opolo/ocl/msgcol"
)

const (
	// PromptAccept for acceptance of a prompt.
	PromptAccept = iota

	// PromptDeny for denial of a prompt.
	PromptDeny

	// PromptTimeout is when the prompt times out.
	PromptTimeout
)

// Context is the message context used for command execution
type Context struct {
	Session *discordgo.Session
	Msg     *discordgo.Message
	Cmd     *Command

	lastMsg *discordgo.Message

	// String used to call the command
	CallKey string
	RawArgs []string
	Args    map[string]interface{}

	StartTime time.Time
}

// ---- Send ----

// SendComplex sends a message of any type to the channel the msg came from.
func (ctx *Context) SendComplex(data *discordgo.MessageSend, deleteTime ...time.Duration) (*discordgo.Message, error) {
	m, err := ctx.Session.ChannelMessageSendComplex(ctx.Msg.ChannelID, data)
	ctx.lastMsg = m

	if len(deleteTime) > 0 {
		go (func() {
			time.Sleep(deleteTime[0])
			ctx.Delete(m)
		})()
	}

	return m, err
}

// SendEmbed sends an embed.
func (ctx *Context) SendEmbed(embed *discordgo.MessageEmbed, deleteTime ...time.Duration) (*discordgo.Message, error) {
	return ctx.SendComplex(&discordgo.MessageSend{
		Embed: embed,
	}, deleteTime...)
}

// Send sends a message to the channel the msg came from.
func (ctx *Context) Send(content string, deleteTime ...time.Duration) (*discordgo.Message, error) {
	return ctx.SendComplex(&discordgo.MessageSend{
		Content: content,
	}, deleteTime...)
}

// SendError reports an error to the err channel and to the user
func (ctx *Context) SendError(err error) {
	cmdName := (func() string {
		if ctx.Cmd != nil && ctx.Cmd.Name != "" {
			return ctx.Cmd.Name
		}

		return "N/A"
	})()

	errEmbed := embed.Error("Failed to run command `"+cmdName+"`", "```yaml\n"+err.Error()+"```")
	ctx.SendEmbed(errEmbed)

	config := util.GetConfig()
	if config.LogChannel != "" {
		stacktrace := string(debug.Stack())
		log.Print("An error occurred while handling the command " + cmdName + ":\n" + err.Error() + "\n" + stacktrace)
		ctx.Session.ChannelMessageSendEmbed(config.LogChannel, errEmbed)
		ctx.Session.ChannelMessageSend(config.LogChannel, "**Stacktrace**\n```go\n" + stacktrace + "```")
	}
}

// ---- Edit ----

// EditComplex edits the last message sent by the bot in the channel with complex content.
func (ctx *Context) EditComplex(data *discordgo.MessageEdit) (*discordgo.Message, error) {
	if ctx.lastMsg == nil {
		return nil, errors.New("Tried to edit a message that does not exist")
	}

	data.ID = ctx.lastMsg.ID
	data.Channel = ctx.lastMsg.ChannelID

	m, err := ctx.Session.ChannelMessageEditComplex(data)
	ctx.lastMsg = m
	return m, err
}

// EditEmbed edits a message replacing the previous
func (ctx *Context) EditEmbed(embed *discordgo.MessageEmbed) (*discordgo.Message, error) {
	return ctx.EditComplex(&discordgo.MessageEdit{
		Embed: embed,
	})
}

// Edit edits the last msg sent by the context.
func (ctx *Context) Edit(content string) (*discordgo.Message, error) {
	return ctx.EditComplex(&discordgo.MessageEdit{
		Content: &content,
	})
}

// ---- Delete ----

// Delete deletes the last recently sent message or the supplied one.
func (ctx *Context) Delete(msg ...*discordgo.Message) error {
	msgSupplied := len(msg) > 0
	if msgSupplied && ctx.lastMsg == nil {
		return nil
	}

	var ID string
	var chanID string
	if msgSupplied {
		ID = msg[0].ID
		chanID = msg[0].ChannelID
	} else {
		ID = ctx.lastMsg.ID
		chanID = ctx.lastMsg.ChannelID
	}

	return ctx.Session.ChannelMessageDelete(chanID, ID)
}

// ---- Message Collectors ----

// Prompt sends a prompt to accept or deny within 10 seconds.
func (ctx *Context) Prompt(prompt string) (int, error) {
	_, err := ctx.Send(prompt + " This action will cancel in 10 seconds. [y/N]")
	if err != nil {
		return PromptDeny, err
	}

	msgs, err := ctx.Collect(10 * time.Second)
	if err != nil {
		return PromptDeny, err
	}

	accepted := false
	denied := false
	for msg := range msgs {
		switch strings.ToLower(msg.Content) {
		case "y":
			accepted = true
			break
		case "n":
			denied = true
			break
		}
	}

	if err = ctx.Delete(); err != nil {
		return PromptDeny, err
	}

	if !accepted {
		if denied {
			ctx.Send("Cancelling purge.")

			return PromptDeny, nil
		}

		ctx.Send("Timed out.")
		return PromptTimeout, nil
	}

	return PromptAccept, nil
}

// CleanUp deletes the last recently sent message.
func (ctx *Context) CleanUp() {
	if ctx.lastMsg != nil {
		time.Sleep(10 * time.Second)
		ctx.Session.ChannelMessageDelete(ctx.lastMsg.ChannelID, ctx.lastMsg.ID)
	}
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

	colMnger := msgcol.GetCollectionManager()
	col := colMnger.NewCollector(ctx.Msg.ChannelID, time, amntToPurge)
	return col.Msgs, nil
}
