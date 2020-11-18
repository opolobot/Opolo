package main

import (
	"github.com/bwmarrin/discordgo"
	"github.com/zorbyte/whiskey/cmd"
	"github.com/zorbyte/whiskey/cmd/core"
	"github.com/zorbyte/whiskey/cmd/fun"
	"github.com/zorbyte/whiskey/cmd/mod"
	"github.com/zorbyte/whiskey/hdlr"
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
