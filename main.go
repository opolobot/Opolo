package main

import (
	"log"
	"math/rand"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/TeamWhiskey/whiskey/util"
	"github.com/bwmarrin/discordgo"
)

func main() {
	// Seed math/rand
	rand.Seed(time.Now().UnixNano())

	config := util.GetConfig()

	session, err := discordgo.New("Bot " + config.Token)
	if err != nil {
		log.Fatal("Failed to create discord session:", err)
	}

	// We need information about guilds (which includes their channels) and msgs.
	session.Identify.Intents = discordgo.MakeIntent(discordgo.IntentsGuilds | discordgo.IntentsGuildMessages)

	startTime := time.Now()
	registerHandlers(session)
	registerCommandCategories()

	// Measure bootstrap speed.
	log.Printf("Bootstap time was %v", time.Since(startTime))

	// Open the websocket and begin listening.
	err = session.Open()
	if err != nil {
		log.Fatal("Failed to establish ws connection:", err)
	}

	// Wait here until CTRL-C or other term signal is received.
	log.Println("Opened WS connection. Use CTRL-C to exit")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc

	// Cleanly close down the Discord session.
	err = session.Close()
	if err != nil {
		log.Fatalf("Failed to close session")
	}
}
