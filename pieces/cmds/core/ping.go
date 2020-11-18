package core

import (
	"fmt"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/opolobot/opolo/ocl"
	"github.com/opolobot/opolo/ocl/embeds"
)

func init() {
	cmd := &ocl.Command{
		Name:        "ping",
		Aliases:     []string{"p"},
		Description: "Calculates the bot latency and execution time",
		Stack:       []ocl.Middleware{ping},
	}

	Category.AddCommand(cmd)
}

func ping(ctx *ocl.Context, next ocl.Next) {
	executionTime := time.Since(ctx.StartTime)
	sentTime := time.Now()
	m, err := ctx.SendEmbed(embeds.Info("Ping?", "ping_pong", ""))
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

	pongEmbed := embeds.Info("Pong", "ping_pong", "")
	pongEmbed.Fields = append(
		pongEmbed.Fields,
		&discordgo.MessageEmbedField{
			Name:   ":arrows_counterclockwise:  latency",
			Value:  fmt.Sprintf("%v", latency),
			Inline: true,
		},
		&discordgo.MessageEmbedField{
			Name:   ":stopwatch:  exec. time",
			Value:  fmt.Sprintf("%v", executionTime),
			Inline: true,
		},
	)

	ctx.EditEmbed(pongEmbed)

	next()
}
