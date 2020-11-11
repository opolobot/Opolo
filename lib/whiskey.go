package lib

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"runtime/debug"
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

	// We need information about guilds (which includes their channels) and msgs.
	s.Identify.Intents = discordgo.MakeIntent(discordgo.IntentsGuilds | discordgo.IntentsGuildMessages)

	whiskey := &Whiskey{
		S:      s,
		Config: config,

		Collectors: make(map[string][]*MsgCollector),

		Cmds:    map[string]*Cmd{},
		aliases: map[string]string{},

		startTime: startTime,
	}

	return whiskey
}

// SendError sends an error to the error channel
func (whiskey *Whiskey) SendError(err error) {
	log.Println(err)
	log.Println(debug.Stack())
	if whiskey.Config.LogChannel != "" {
		errStr := fmt.Sprintf("**:warning: An error occurred in-flight**\n```%v\n%v```", err.Error(), string(debug.Stack()))
		whiskey.S.ChannelMessageSend(whiskey.Config.LogChannel, errStr)
	}
}

// Start starts the bot and establishes a ws conn.
func (whiskey *Whiskey) Start() {
	// Measure bootstrap speed.
	log.Printf("Bootstap time was %v", time.Since(whiskey.startTime))

	// Open the websocket and begin listening.
	err := whiskey.S.Open()
	if err != nil {
		log.Fatal("Failed to establish ws connection:", err)
	}

	// Wait here until CTRL-C or other term signal is received.
	log.Println("Opened WS connection. Use CTRL-C to exit")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc

	// Cleanly close down the Discord session.
	whiskey.S.Close()
}
