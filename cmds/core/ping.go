package core

import (
	"fmt"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/TeamWhiskey/whiskey/cmds"
)

func init() {
	cmdBldr := cmds.NewCommandBuilder()
	cmdBldr.Name("ping")
	cmdBldr.Aliases("p")
	cmdBldr.Description("Pings the bot")
	cmdBldr.Use(ping)

	Category.AddCommand(cmdBldr.Build())
}

func ping(ctx *cmds.Context, next cmds.NextFunc) {
	executionTime := time.Since(ctx.StartTime)
	startTime := time.Now()
	m, err := ctx.Send(":ping_pong: Ping?")
	if err != nil {
		next(err)
	}

	messageSentTime, err := discordgo.SnowflakeTimestamp(m.ID)
	if err != nil {
		next(err)
	}

	ctx.Edit(fmt.Sprintf(
		"***:ping_pong:  ~Pong!***\n"+
			"\n> __**latency**__        **~**   :arrows_counterclockwise: %v"+
			"\n> __**exec. time**__   **~**   :stopwatch: %v"+
			"\n> __**heartbeat**__    **~**   :heartbeat: %v",
		startTime.Sub(messageSentTime).Round(time.Millisecond),
		executionTime.Round(time.Microsecond),
		ctx.Session.HeartbeatLatency().Round(time.Microsecond),
	))

	next()
}
