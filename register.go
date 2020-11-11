package main

import (
	"github.com/bwmarrin/discordgo"
	"github.com/TeamWhiskey/whiskey/cmds"
	"github.com/TeamWhiskey/whiskey/cmds/core"
	"github.com/TeamWhiskey/whiskey/cmds/fun"
	"github.com/TeamWhiskey/whiskey/cmds/mod"
	"github.com/TeamWhiskey/whiskey/hdlrs"
)

func registerHandlers(session *discordgo.Session) {
	session.AddHandler(hdlrs.Ready)
	session.AddHandler(hdlrs.MessageCreate)
}

func registerCommandCategories() {
	cmdUI := cmds.GetCommandUI()
	cmdUI.AddCategory(core.Category)
	cmdUI.AddCategory(fun.Category)
	cmdUI.AddCategory(mod.Category)
	cmdUI.Build()
}
