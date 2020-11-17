package main

import (
	"rand"
	"time"

	"github.com/opolobot/opolo/services/cli"
)

func main() {
	// Seed math/rand
	rand.Seed(time.Now().UnixNano())

	cli.StartCLI()
}
