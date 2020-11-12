package utils

import (
	"fmt"
	"log"
	"runtime/debug"

	"github.com/bwmarrin/discordgo"
)

// SendError sends an error to the error channel
func SendError(session *discordgo.Session, err error) {
	config := GetConfig()
	log.Println(err)
	log.Println(string(debug.Stack()))
	if config.LogChannel != "" {
		errStr := fmt.Sprintf("**:warning: An error occurred in-flight**\n\n%v\n**Stacktrace**\n```%v```", err.Error(), string(debug.Stack()))
		session.ChannelMessageSend(config.LogChannel, errStr)
	}
}
