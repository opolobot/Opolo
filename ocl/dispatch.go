package ocl

import (
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/opolobot/opolo/common"
	"github.com/opolobot/opolo/ocl/args"
	"github.com/opolobot/opolo/ocl/embeds"
)

var stringSplitter = regexp.MustCompile(" +")

// Dispatch dispatches a command.
func Dispatch(session *discordgo.Session, msg *discordgo.Message) (next Next) {
	startTime := time.Now()

	prefix := common.GetConfig().Prefix
	if !strings.HasPrefix(msg.Content, prefix) {
		return
	}

	callKey, rawArgs := parseContent(prefix, msg.Content)
	ctx := &Context{
		Session: session,
		Msg:     msg,

		CallKey: callKey,
		RawArgs: rawArgs,

		StartTime: startTime,
	}

	reg := GetRegistry()
	cmd, err := reg.LookupCommand(callKey)
	if err != nil {
		ctx.SendError(err)
	}

	if cmd != nil {
		ctx.Cmd = cmd
		success := parseArgs(ctx)
		if !success {
			return
		}

		idx := -1
		next = func(err ...error) {
			if idx < 0 {
				defer handleInFlightPanic(ctx)
			}

			if len(err) > 0 {
				ctx.SendError(err[0])
				return
			} else if idx++; idx <= len(cmd.Stack)-1 {
				cmd.Stack[idx](ctx, next)
			}
		}

		return
	}

	didYouMean(ctx)

	return
}

func parseContent(prefix, content string) (callKey string, rawArgs []string) {
	excludingPrefix := strings.TrimSpace(content)[len(prefix):]
	segments := stringSplitter.Split(excludingPrefix, -1)
	callKey = strings.ToLower(segments[0])
	rawArgs = segments[1:]

	return
}

func parseArgs(ctx *Context) (success bool) {
	var err error
	ctx.Args, err = args.Parse(ctx.Cmd.Arguments, ctx.RawArgs)
	if success = err == nil; !success {
		if pErr, ok := err.(*args.ParsingError); ok {
			ctx.Send(pErr.UIError())
		} else {
			ctx.SendError(err)
		}
	}

	return
}

func handleInFlightPanic(ctx *Context) {
	if r := recover(); r != nil {
		ctx.SendError(fmt.Errorf("%v", r))
	}
}

func didYouMean(ctx *Context) {
	reg := GetRegistry()
	prefix := common.GetConfig().Prefix

	closest, distance := reg.FindClosestCmdMatch(ctx.CallKey)
	if distance <= 2 && distance != 0 {
		ctx.SendEmbed(embeds.Info(fmt.Sprintf("Did you mean `%v`?", prefix+closest), "question", ""))
	}
}
