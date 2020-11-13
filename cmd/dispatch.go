package cmd

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/TeamWhiskey/whiskey/arg"
	"github.com/TeamWhiskey/whiskey/util"
	"github.com/bwmarrin/discordgo"
)

var stringSplitter = regexp.MustCompile(" +")

// Dispatch dispatches a command.
func Dispatch(session *discordgo.Session, msg *discordgo.Message) NextFunc {
	startTime := time.Now()

	prefix := util.GetConfig().Prefix
	if !strings.HasPrefix(msg.Content, prefix) {
		return nil
	}

	callKey, rawArgs := parseContent(msg.Content, prefix)

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
			return nil
		}

		var nextFunc NextFunc
		idx := -1
		nextFunc = func(err ...error) {
			if idx < 0 {
				defer handleInFlightPanic(ctx)
			}

			if len(err) > 0 {
				ctx.SendError(err[0])
				return
			} else if idx++; idx <= len(cmd.stack)-1 {
				cmd.stack[idx](ctx, nextFunc)
			}
		}

		return nextFunc
	}

	return nil
}

func parseContent(prefix, content string) (callKey string, rawArgs []string) {
	excludingPrefix := strings.TrimSpace(strings.ToLower(content))[len(prefix):]
	segments := stringSplitter.Split(excludingPrefix, -1)
	callKey = segments[0]
	rawArgs = segments[1:]

	return callKey, rawArgs
}

func parseArgs(ctx *Context) (success bool) {
	var err error
	ctx.Args, err = arg.Parse(ctx.Cmd.args, ctx.RawArgs)
	if success = err == nil; !success {
		if pErr, ok := err.(*arg.ParsingError); ok {
			err = errors.Unwrap(pErr)
			if err != nil {
				ctx.SendError(err)
			} else if pErr.ErrorType != arg.InternalParsingError {
				ctx.Send(pErr.UIError())
			} else {
				ctx.SendError(pErr)
			}
		}
	}

	return success
}

func handleInFlightPanic(ctx *Context) {
	if r := recover(); r != nil {
		ctx.SendError(fmt.Errorf("%v", r))
	}
}

func didYouMean(ctx *Context) {
	reg := GetRegistry()
	config := util.GetConfig()

	closest, distance := reg.FindClosestCmdMatch(ctx.CallKey)
	if distance <= 2 && distance != 0 {
		ctx.Send(fmt.Sprintf("**:question: ~ Did you mean `%v`?**", config.Prefix+closest))
	}
}
