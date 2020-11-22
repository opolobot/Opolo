package fun

import (
	"fmt"
	"math/rand"
	"strings"
	"time"

	"github.com/opolobot/Opolo/ocl"
	"github.com/opolobot/Opolo/ocl/args"
	"github.com/opolobot/Opolo/pieces/parsers"
)

var positiveResponses = []string{
	"It is certain",
	"It is decidedly so",
	"Without a doubt",
	"Yes - definitely",
	"You may rely on it",
	"As I see it, yes",
	"Most likely",
	"Outlook good",
	"Yes",
	"Signs point to yes",
}

var negativeResponses = []string{
	"Don't count on it",
	"My reply is no",
	"My sources say no",
	"Outlook not so good",
	"Very doubtful",
}

var responses = append(positiveResponses, negativeResponses...)

var hazyResponses = [...]string{
	"Reply hazy, try again",
	"Ask again later",
	"Better not tell you now",
	"Cannot predict now",
	"Concentrate and ask again",
}

func init() {
	cmd := ocl.New()
	cmd.Name("8ball")
	cmd.Aliases("8")
	cmd.Description("Ask 8ball a burning yes/no question")
	cmd.Args(args.New("<...question>", &parsers.String{}))
	cmd.Use(eightball)

	Category.Add(cmd)
}

func eightball(ctx *ocl.Context, _ ocl.Next) {
	var emoji, message string

	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	hazy := r.Intn(6)

	// 1/6 chance of giving a hazy response
	if hazy == 5 {
		i := r.Intn(len(hazyResponses))
		message = hazyResponses[i]
		emoji = ":8ball:"
	} else {
		r.Seed(genSeed(ctx))
		i := r.Intn(len(responses))
		message = responses[i]

		if i >= len(positiveResponses) {
			emoji = ":red_circle:"
		} else {
			emoji = ":green_circle:"
		}
	}

	ctx.Send(fmt.Sprintf("%s %s", emoji, message))
}

func genSeed(ctx *ocl.Context) int64 {
	joinedArgs := strings.Join(ctx.RawArgs, " ")
	joinedRunesAndAuthor := []rune(joinedArgs + ctx.Msg.Author.ID)
	seed := runesToSeed(joinedRunesAndAuthor)
	return seed
}

func runesToSeed(runes []rune) (result int64) {
	for i, v := range runes {
		result += int64(v) * int64(i)
	}

	return
}
