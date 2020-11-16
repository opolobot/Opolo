package fun

import (
	"fmt"
	"math/rand"
	"strings"
	"time"

	"github.com/zorbyte/whiskey/cmd"
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
	cmd := cmd.New()
	cmd.Name("8ball")
	cmd.Aliases("8")
	cmd.Description("Ask 8ball a burning yes/no question")
	cmd.Use(eightball)

	Category.AddCommand(cmd.Command())
}

func eightball(ctx *cmd.Context, next cmd.NextFunc) {
	if len(ctx.RawArgs) == 0 {
		ctx.Send(":8ball: Cannot reply, did not receive a fucking question")
		return
	}

	var emoji string
	var message string

	rand.Seed(time.Now().UnixNano())
	hazy := rand.Intn(6)

	if hazy == 5 { // 1/6 chance of giving a hazy response
		i := rand.Intn(len(hazyResponses))

		message = hazyResponses[i]
		emoji = ":8ball:"
	} else {
		rand.Seed(genSeed(ctx))

		i := rand.Intn(len(responses))

		message = responses[i]

		if i > len(positiveResponses) {
			emoji = ":red_circle:"
		} else {
			emoji = ":green_circle:"
		}
	}

	ctx.Send(fmt.Sprintf("%s %s", emoji, message))
}

func genSeed(ctx *cmd.Context) int64 {
	joinedArgs := strings.Join(ctx.RawArgs, " ")
	slice := []rune(joinedArgs + ctx.Msg.Author.ID)
	seed := sum(slice)
	return seed
}

func sum(array []rune) int64 {
	var result int64
	for _, v := range array {
		result += int64(v)
	}
	return result
}
