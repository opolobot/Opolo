package cmds

import (
	"fmt"
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
	categoryDisplayName := emoji + " __**" + category + "**__"
	curHelpOutputs := make(map[string]string)
	helpOutputs[categoryDisplayName] = &curHelpOutputs
	for _, cmd := range *categoryCmds {
		cmd.Category = categoryDisplayName

		// In the case that no name was provided.
		if cmd.Name == "" {
			cmd.Name = getCmdRunnerName(cmd.Runner)
		}

		// Prebuild help strings for performance (maybe?)
		// Note: The first **` is missing for dynamic prefix support, as
		// including it and styling the dynamic prefix would ruin the formatting
		usageStr := (func() string {
			if cmd.Usage != "" {
				return " " + cmd.Usage
			}

			return ""
		})()

		cmdHelpStr := fmt.Sprintf(
			"%v%v`**\n> **desc ~** %v\n> **aliases ~** %v",
			cmd.Name, usageStr,
			cmd.Description, strings.Join(cmd.Aliases, ", "),
		)
		curHelpOutputs[cmd.Name] = cmdHelpStr

		log.Printf("Registering command %v\n", cmd.Name)
		w.RegCmd(cmd)
	}
}

func getCmdRunnerName(i interface{}) string {
	nameWithPkg := runtime.FuncForPC(reflect.ValueOf(i).Pointer()).Name()
	return nameWithPkg[pkgNameLen:]
}
