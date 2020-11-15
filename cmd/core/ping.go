package core

import (
	"fmt"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/zorbyte/whiskey/cmd"
)

func init() {
	cmd := cmd.New()
	cmd.Name("ping")
	cmd.Aliases("p")
	cmd.Description("Pings the bot")
	cmd.Use(ping)

	Category.AddCommand(cmd.Command())
}

func ping(ctx *cmd.Context, next cmd.NextFunc) {
	executionTime := time.Since(ctx.StartTime)
	sentTime := time.Now()
	m, err := ctx.Send(":ping_pong: Ping?")
	if err != nil {
		next(err)
	}

	msgSentTime, err := discordgo.SnowflakeTimestamp(m.ID)
	if err != nil {
		next(err)
	}

	latency := msgSentTime.Sub(sentTime).Round(time.Millisecond)
	if latency < 0 {
		latency *= -1
	}

	heartbeatTime := ctx.Session.HeartbeatLatency().Round(time.Millisecond).Seconds()

	ctx.Edit(fmt.Sprintf(
		"***:ping_pong:  ~Pong!***\n"+
			"\n> __**latency**__        **~**   :arrows_counterclockwise: %v"+
			"\n> __**exec. time**__   **~**   :stopwatch: %v"+
			"\n> __**heartbeat**__    **~**   :heartbeat: %.1fs",
		latency,
		executionTime,
		heartbeatTime,
	))

	next()
}
