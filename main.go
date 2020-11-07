package main

import (
	"fmt"
	"log"
	"regexp"
	"strings"

	"github.com/bwmarrin/discordgo"
)

func init() {
	log.Printf("Loading config file @ %v\n", ConfFileName)
	var err error
	config, err = FetchConf()
	if err != nil {
		log.Fatal("Failed to load config: ", err)
	}
}

var (
	stringSplitter = regexp.MustCompile(" +")

	config *Config
	cmds   = make(map[string]func(s *discordgo.Session, m *discordgo.MessageCreate, args ...string) (string, error))
)

func main() {
	w := NewWhiskey()

	cmds["test"] = func(s *discordgo.Session, m *discordgo.MessageCreate, args ...string) (string, error) {
		if strings.Join(args, " ") == "throw an error" {
			return "", fmt.Errorf("Arguments to throw an error were satisfied")
		}
		return "it worked, args: " + strings.Join(args, ", "), nil
	}

	w.Start()
}

func ready(s *discordgo.Session, m *discordgo.Ready) {
	log.Printf("Logged in as %v#%v (%v)", s.State.User.Username, s.State.User.Discriminator, s.State.User.ID)
	s.UpdateStatus(1, config.Status)
}

func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	// Ignore messages by the bot itself.
	if m.Author.ID == s.State.User.ID || !strings.HasPrefix(m.Content, config.Prefix) {
		return
	}

	cmdStr := strings.TrimSpace(m.Content)[len(config.Prefix):]

	cmdSegs := stringSplitter.Split(cmdStr, -1)
	cmdName := cmdSegs[0]
	cmdArgs := cmdSegs[1:]
	log.Printf("cmdName (cmdSegs[0]) -> %v, cmdSegs -> %v, cmdSegs[0:] -> %v, cmdSegs[1:] -> %v\n", cmdName, cmdSegs, cmdSegs[0:], cmdArgs)

	if cmdRunner, ok := cmds[cmdName]; ok {
		strToSend, err := cmdRunner(s, m, cmdArgs...)
		if err != nil {
			go sendError(s, m.ChannelID, err, cmdName)
			return
		}

		if strToSend != "" {
			s.ChannelMessageSend(m.ChannelID, strToSend)
		}
	}
}

func sendError(s *discordgo.Session, chanID string, err error, runningCmd ...string) {
	cmdName := (func() string {
		if len(runningCmd) > 0 {
			return runningCmd[0]
		}

		return "N/A"
	})()

	errTxt := ":rotating_light: An error occurred while handling the command `" + cmdName + "`:\n```" + err.Error() + "```"
	s.ChannelMessageSend(chanID, errTxt+"\n\nThe error has been reported")
	s.ChannelMessageSend(config.ErrChannel, errTxt)
	log.Println("An error occurred while handling the command "+cmdName+":", err)
}
