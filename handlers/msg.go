package handlers

import (
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

		// Ignore messages by the bot itself.
		if m.Author.ID == s.State.User.ID || !strings.HasPrefix(m.Content, w.Config.Prefix) {
			return
		}

		cmdStr := strings.TrimSpace(m.Content)[len(w.Config.Prefix):]

		cmdSegs := stringSplitter.Split(cmdStr, -1)
		cmdCallKey := cmdSegs[0]
		cmdArgs := cmdSegs[1:]
		log.Printf("cmdName -> %v, cmdArgs -> %v\n", cmdCallKey, cmdArgs)

		if cmd := w.FindCmd(cmdCallKey); cmd != nil {
			c := &lib.Ctx{
				W: w,
				S: s,
				M: m.Message,
				C: cmd,

				CmdCallKey: cmdCallKey,
				Args:       cmdArgs,

				StartTime: startTime,
			}

			strToSend, err := cmd.Runner(c)
			if err != nil {
				go c.SendError(err)
				return
			}

			log.Printf("Time for command execution: %v", time.Since(startTime))
			if strToSend != "" {
				c.Send(strToSend)
			}
		}
	}
}
