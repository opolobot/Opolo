package handlers

import (
	"log"

	"github.com/zorbyte/whiskey/lib"
)

// RegisterHandlers registers event handlers.
func RegisterHandlers(w *lib.Whiskey) {
	log.Println("Registering event handlers")
	w.S.AddHandler(MsgCreate(w))
	w.S.AddHandlerOnce(Ready(w))
}
