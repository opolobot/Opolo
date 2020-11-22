package fun

import (
	"strings"

	"github.com/opolobot/Opolo/ocl"
	"github.com/opolobot/Opolo/ocl/args"
	"github.com/opolobot/Opolo/pieces/parsers"
	"github.com/rivo/uniseg"
)

const max = 41
const cow = `
       \   ^__^
        \  (oo)\_______
           (__)\       )\/\
                ||----w |
                ||     ||
`

func init() {
	cmd := ocl.New()
	cmd.Name("cowsay")
	cmd.Aliases("cow", "moo")
	cmd.Description("Generates an ASCII picture of a cow saying something provided by the user")
	cmd.Args(args.New("<...message>", &parsers.String{}))
	cmd.Use(cowsay)

	Category.Add(cmd)
}

func cowsay(ctx *ocl.Context, next ocl.Next) {
	message := strings.Join(ctx.RawArgs, " ")
	send := createBubble(message, max)
	send += cow
	ctx.Send("```\n" + send + "```")
}

func createBubble(text string, lineWidth int) (bubble string) {
	lines := strings.Split(wrap(text, lineWidth), "\n")

	var width int
	if len(lines) == 1 {
		width = wordlen(lines[0])
	} else {
		width = lineWidth
	}

	bubble += "  " + strings.Repeat("_", width) + "\n"

	for i, line := range lines {
		l, r := border(lines, i)

		bubble += l + " "
		bubble += line
		bubble += strings.Repeat(" ", 1+(width-wordlen(line))) + r + "\n"
	}

	bubble += "  " + strings.Repeat("-", width)

	return
}

func border(lines []string, i int) (string, string) {
	if len(lines) == 1 {
		return "<", ">"
	} else if i == 0 {
		return "/", "\\"
	} else if i == len(lines)-1 {
		return "\\", "/"
	}
	return "|", "|"
}

func wrap(text string, lineWidth int) (result string) {
	words := strings.Split(text, " ")
	words = splitLongWords(words, lineWidth)

	if len(words) == 0 {
		return
	}

	remaining := lineWidth

	for _, word := range words {
		wordlen := wordlen(word)
		if wordlen+1 > remaining {
			if uniseg.GraphemeClusterCount(result) > 0 {
				result += "\n"
			}
			result += word
			remaining = lineWidth - wordlen
		} else {
			if uniseg.GraphemeClusterCount(result) > 0 {
				result += " "
			}
			result += word
			remaining -= wordlen + 1
		}
	}

	return
}

func splitLongWords(words []string, lineWidth int) (result []string) {
	for _, word := range words {
		if wordlen(word) > lineWidth {
			split := splitWord(word, lineWidth)
			result = append(result, split...)
		} else {
			result = append(result, word)
		}
	}

	return
}

func splitWord(word string, lineWidth int) (parts []string) {
	var all string

	gr := uniseg.NewGraphemes(word)
	for gr.Next() {
		if wordlen(all) == lineWidth {
			parts = append(parts, all)
			all = ""
		}
		all += gr.Str()
	}

	if wordlen(all) > 0 {
		parts = append(parts, all)
	}

	return
}

func wordlen(word string) int {
	return uniseg.GraphemeClusterCount(word)
}
