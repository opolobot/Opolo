package handlers

import "github.com/zorbyte/whiskey/lib"

// RegisterHandlers registers event handlers.
func RegisterHandlers(w *lib.Whiskey) {
	w.S.AddHandler(MsgCreate(w))
	w.S.AddHandlerOnce(Ready(w))
}
