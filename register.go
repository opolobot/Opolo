package main

import (
	"github.com/bwmarrin/discordgo"
	"github.com/opolobot/opolo/ocl"
	"github.com/opolobot/opolo/pieces/cmds/core"
	"github.com/opolobot/opolo/pieces/cmds/fun"
	"github.com/opolobot/opolo/pieces/events"
)

func registerHandlers(session *discordgo.Session) {
	session.AddHandler(events.Ready)
	session.AddHandler(events.MessageCreate)
}

func registerCommandCategories() {
	reg := ocl.GetRegistry()
	reg.AddCategory(fun.Category)
	reg.AddCategory(core.Category)
	reg.Populate()
}
