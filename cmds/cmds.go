package cmds

import (
	"log"
	"reflect"
	"runtime"
	"strings"

	"github.com/zorbyte/whiskey/lib"
)

// Length of "github.com/zorbyte/whiskey/cmds."
const pkgNameLen uint16 = 32

var w *lib.Whiskey
var helpOutputs = make(map[string]*map[string]string)

// RegisterCmds registers all commands from each category
func RegisterCmds(whiskey *lib.Whiskey) {
	// Set the w variable so that we don't have to pass it repeteatedly
	w = whiskey
	log.Println("Registering commands")

	registerCategory("core", "Whiskey Core", coreCmdsEmoji, &coreCmds)
}

func registerCategory(name string, category string, emoji string, categoryCmds *[]*lib.Cmd) {
	log.Printf("Registering %v commands\n", name)

	// Set categories with styling for pre-made help components.
	categoryDisplayName := emoji + " __**" + category + "**__"
	curHelpOutputs := make(map[string]string)
	helpOutputs[categoryDisplayName] = &curHelpOutputs

	for _, cmd := range *categoryCmds {
		cmd.Category = categoryDisplayName

		// In the case that no name was provided.
		if cmd.Name == "" {
			cmd.Name = getCmdRunnerName(cmd.Runner)
		}

		curHelpOutputs[cmd.Name] = generateHelpStr(cmd)

		log.Printf("Registering command %v\n", cmd.Name)
		w.RegCmd(cmd)
	}
}

func getCmdRunnerName(i interface{}) string {
	nameWithPkg := runtime.FuncForPC(reflect.ValueOf(i).Pointer()).Name()
	return nameWithPkg[pkgNameLen:]
}

func generateHelpStr(cmd *lib.Cmd) string {
	// Prebuild help strings for performance (maybe?)
	var cmdHelpStrBldr strings.Builder

	// Name and usage.
	cmdHelpStrBldr.WriteString(cmd.Name)
	if cmd.Usage != "" {
		cmdHelpStrBldr.WriteString(" ")
		cmdHelpStrBldr.WriteString(cmd.Usage)
	}

	// The first **` is missing for dynamic prefix support, as
	// including it and styling the dynamic prefix would ruin the formatting
	cmdHelpStrBldr.WriteString("`**")

	// Description.
	if cmd.Description != "" {
		cmdHelpStrBldr.WriteString("\n> **desc ~** ")
		cmdHelpStrBldr.WriteString(cmd.Description)
	}

	// Aliases
	if len(cmd.Aliases) > 0 {
		cmdHelpStrBldr.WriteString("\n> **aliases ~** ")
		cmdHelpStrBldr.WriteString(strings.Join(cmd.Aliases, ", "))
	}

	return cmdHelpStrBldr.String()
}
