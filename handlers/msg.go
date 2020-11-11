package handlers

import (
	"fmt"
	"log"
	"regexp"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/zorbyte/whiskey/lib"
)

var stringSplitter = regexp.MustCompile(" +")

// MsgCreate handles msg create event and dispatches commands
func MsgCreate(w *lib.Whiskey) func(s *discordgo.Session, m *discordgo.MessageCreate) {
	return func(s *discordgo.Session, m *discordgo.MessageCreate) {
		startTime := time.Now()
		defer (func() {
			if r := recover(); r != nil {
				w.SendError(fmt.Errorf("Panic in message event!\n%v", r))
			}
		})()

		// Ignore messages by the bot itself.
		if m.Author.ID == s.State.User.ID {
			return
		}

		if !strings.HasPrefix(m.Content, w.Config.Prefix) {
			for chanID, cols := range w.Collectors {
				if chanID == m.ChannelID {
					for _, col := range cols {
						col.Msgs <- m.Message
						if len(col.Msgs) == cap(col.Msgs) {
							w.CancelCollector(chanID, col.ID)
						}
					}
				}
			}

			return
		}

		cmdStr := strings.TrimSpace(m.Content)[len(w.Config.Prefix):]

		cmdSegs := stringSplitter.Split(cmdStr, -1)
		cmdCallKey := cmdSegs[0]
		cmdArgs := cmdSegs[1:]
		log.Printf("cmdName -> %v, cmdArgs -> %v\n", cmdCallKey, cmdArgs)

		if cmd := w.FindCmd(cmdCallKey); cmd != nil {
			ctx := &lib.Ctx{
				W: w,
				S: s,
				M: m.Message,
				C: cmd,

				CmdCallKey: cmdCallKey,
				Args:       cmdArgs,

				StartTime: startTime,
			}

			// It's pretty nasty that this is being repeated, hopefully an
			// alternative can be found.
			defer (func() {
				if r := recover(); r != nil {
					ctx.SendError(fmt.Errorf("%v", r))
				}
			})()

			strToSend, err := cmd.Runner(ctx)
			if err != nil {
				ctx.SendError(err)
				return
			}

			log.Printf("Time for command execution: %v", time.Since(startTime))
			if strToSend != "" {
				ctx.Send(strToSend)
			}
		}
	}
}
