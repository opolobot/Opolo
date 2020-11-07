package main

import "github.com/bwmarrin/discordgo"

// MCtx is the message context used for command execution
type MCtx struct {
	W *Whiskey
	M discordgo.MessageCreate

	// String used to call the command
	CmdKey string
	Args   []string

	// TODO(@zorbyte): Create a db, consider using gorm.
	DB interface{}
}
