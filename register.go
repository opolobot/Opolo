package main

import (
	"github.com/bwmarrin/discordgo"
	"github.com/TeamWhiskey/whiskey/cmd"
	"github.com/TeamWhiskey/whiskey/cmd/core"
	"github.com/TeamWhiskey/whiskey/cmd/fun"
	"github.com/TeamWhiskey/whiskey/cmd/mod"
	"github.com/TeamWhiskey/whiskey/hdlr"
)

func registerHandlers(session *discordgo.Session) {
	session.AddHandler(hdlrs.Ready)
	session.AddHandler(hdlrs.MessageCreate)
}

func registerCommandCategories() {
	cmdUI := cmd.GetCommandUI()
	cmdUI.AddCategory(core.Category)
	cmdUI.AddCategory(fun.Category)
	cmdUI.AddCategory(mod.Category)
	cmdUI.Build()
}
