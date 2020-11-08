package lib

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/bwmarrin/discordgo"
)

// CmdRunner is the func to run for the cmd
type CmdRunner func(ctx *Ctx) (string, error)

// Cmd is information about the command as well as its runner
type Cmd struct {
	Runner      CmdRunner
	Name        string
	Usage       string
	Aliases     []string
	Description string
	Category    string
}

// Whiskey is the mediator for the bot
type Whiskey struct {
	S      *discordgo.Session
	Config *Config

	cmds    map[string]*Cmd
	aliases map[string]string
}

// NewWhiskey creates a Whiskey instance
func NewWhiskey() *Whiskey {
	log.Printf("Loading config file @ %v\n", ConfFileName)
	config, err := FetchConf()
	if err != nil {
		log.Fatal("Failed to load config: ", err)
	}

	s, err := discordgo.New("Bot " + config.Token)
	if err != nil {
		log.Fatal("Failed to create discord session:", err)
	}

	// We need information about guilds (which includes their channels),
	// messages and voice states.
	s.Identify.Intents = discordgo.MakeIntent(discordgo.IntentsGuilds | discordgo.IntentsGuildMessages | discordgo.IntentsGuildVoiceStates)

	w := &Whiskey{
		S:      s,
		Config: config,

		cmds:    map[string]*Cmd{},
		aliases: map[string]string{},
	}

	return w
}

// RegCmd registers a command to whiskey.
func (w *Whiskey) RegCmd(cmd *Cmd) {
	_, ok := w.cmds[cmd.Name]
	if ok {
		log.Fatal("Duplicate command key " + cmd.Name + " was registered")
	}

	w.cmds[cmd.Name] = cmd

	for _, alias := range cmd.Aliases {
		_, ok = w.aliases[alias]
		if ok {
			log.Fatal("Duplicate command alias " + alias + " was registered")
		}

		w.aliases[alias] = cmd.Name
	}
}

// FindCmd finds the cmd with either the name or the alias
func (w *Whiskey) FindCmd(name string) *Cmd {
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
