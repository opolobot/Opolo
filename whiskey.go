package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/bwmarrin/discordgo"
)

// Cmd is a cmd
type Cmd func(s *discordgo.Session, m *discordgo.MessageCreate, args ...string) (string, error)

type cmdData struct {
	name        string
	usage       string
	aliases     []string
	description string
}

// CmdRegistration is the struct used to register a cmd.
type CmdRegistration struct {
	cmdData
	cmd Cmd
}

// Whiskey is the mediator for the bot
type Whiskey struct {
	S *discordgo.Session

	// May seem counterintuitive to not use structures, however this works too and a refactor will likely be in
	// the works for that anyway.
	cmds       map[string]Cmd
	cmdDescips map[string]string
	cmdUsages  map[string]string
	aliases    map[string]string
}

// NewWhiskey creates a Whiskey instance
func NewWhiskey() *Whiskey {
	s, err := discordgo.New("Bot " + config.Token)
	if err != nil {
		log.Fatal("Failed to create discord session:", err)
	}

	s.AddHandler(messageCreate)
	s.AddHandler(ready)

	// We need information about guilds (which includes their channels),
	// messages and voice states.
	s.Identify.Intents = discordgo.MakeIntent(discordgo.IntentsGuilds | discordgo.IntentsGuildMessages | discordgo.IntentsGuildVoiceStates)

	w := &Whiskey{
		S: s,
	}

	return w
}

// RegCmd registers a command to whiskey.
func (w *Whiskey) RegCmd(cmd *CmdRegistration) {
	_, ok := w.cmds[cmd.name]
	if ok {
		log.Fatal("Duplicate command key " + cmd.name + " was registered")
	}

	w.cmds[cmd.name] = cmd.cmd

	for _, alias := range cmd.aliases {
		_, ok = w.aliases[alias]
		if ok {
			log.Fatal("Duplicate command alias " + alias + " was registered")
		}

		w.aliases[alias] = cmd.name
	}
}

// FindCmd finds the cmd with either the name or the alias
func (w *Whiskey) FindCmd(name string) Cmd {
	cmd, ok := w.cmds[name]
	if !ok {
		name, ok = w.aliases[name]
		if ok {
			return w.FindCmd(name)
		}

		return nil
	}

	return cmd
}

// Start starts the bot and establishes a ws conn.
func (w *Whiskey) Start() {
	// Open the websocket and begin listening.
	err := w.S.Open()
	if err != nil {
		log.Fatal("Failed to establish ws connection:", err)
	}

	// Wait here until CTRL-C or other term signal is received.
	log.Println("Whiskey is now running. Press CTRL-C to exit")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc

	// Cleanly close down the Discord session.
	log.Println("Shutting down...")
	w.S.Close()
}
