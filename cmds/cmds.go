package cmds

import (
	"log"

	"github.com/zorbyte/whiskey/lib"
)

// RegisterCmds registers all commands from each category
func RegisterCmds(w *lib.Whiskey) {
	log.Println("Registering random commands")
	for _, cmd := range randomCmds {
		cmd.Category = "random"
		log.Printf("Registering command %v\n", cmd.Name)
		w.RegCmd(cmd)
	}
}
