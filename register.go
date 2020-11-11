package main

import (
	"github.com/bwmarrin/discordgo"
	"github.com/zorbyte/whiskey/cmds"
	"github.com/zorbyte/whiskey/cmds/core"
	"github.com/zorbyte/whiskey/cmds/fun"
	"github.com/zorbyte/whiskey/cmds/mod"
	"github.com/zorbyte/whiskey/hdlrs"
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
