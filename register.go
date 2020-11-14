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
	session.AddHandler(hdlr.Ready)
	session.AddHandler(hdlr.MessageCreate)
}

func registerCommandCategories() {
	reg := cmd.GetRegistry()
	reg.AddCategory(core.Category)
	reg.AddCategory(fun.Category)
	reg.AddCategory(mod.Category)
	reg.Populate()
}
