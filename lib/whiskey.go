package lib

import (
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/bwmarrin/discordgo"
)

// Whiskey is the mediator for the bot
type Whiskey struct {
	S      *discordgo.Session
	Config *Config

	Collectors map[string][]*MsgCollector

	Cmds    map[string]*Cmd
	aliases map[string]string

	startTime time.Time
}

// NewWhiskey creates a Whiskey instance
func NewWhiskey() *Whiskey {
	startTime := time.Now()
	log.Printf("Loading config file @ %v\n", configFileName)
	config, err := readConfig()
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

		Collectors: make(map[string][]*MsgCollector),

		Cmds:    map[string]*Cmd{},
		aliases: map[string]string{},

		startTime: startTime,
	}

	return w
}

// SendError sends an error to the error channel
func (w *Whiskey) SendError(errText string) {
	if w.Config.LogChannel != "" {
		w.S.ChannelMessageSend(w.Config.LogChannel, errText)
	}
}

// Start starts the bot and establishes a ws conn.
func (w *Whiskey) Start() {
	// Measure bootstrap speed.
	log.Printf("Bootstap time was %v", time.Since(w.startTime))

	// Open the websocket and begin listening.
	err := w.S.Open()
	if err != nil {
		log.Fatal("Failed to establish ws connection:", err)
	}

	// Wait here until CTRL-C or other term signal is received.
	log.Println("Opened WS connection. Use CTRL-C to exit")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc

	// Cleanly close down the Discord session.
	w.S.Close()
}
