package main

import (
	"math/rand"
	"time"

	"github.com/zorbyte/whiskey/cmdsOld"
	"github.com/zorbyte/whiskey/handlers"
	"github.com/zorbyte/whiskey/lib"
)

func main() {
	// Seed math/rand
	rand.Seed(time.Now().UnixNano())

	w := lib.NewWhiskey()
	handlers.RegisterHandlers(w)
	cmdsOld.RegisterCmds(w)
	w.Start()
}
