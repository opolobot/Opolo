package fun

import (
	"bytes"
	"image"
	"image/gif"
	"image/jpeg"
	"image/png"
	"io"
	"io/ioutil"
	"net/http"

	"github.com/bwmarrin/discordgo"
	"github.com/fogleman/gg"
	"github.com/opolobot/Opolo/ocl"
	"github.com/opolobot/Opolo/ocl/args"
	"github.com/opolobot/Opolo/pieces/parsers"
	"mvdan.cc/xurls/v2"
)

func init() {
	cmd := ocl.New()
	cmd.Name("meme")
	cmd.Aliases("text")
	cmd.Description("add top and bottom text to an image")
	cmd.Args(
		args.New("[image-url]", &parsers.String{}),
		args.New("[bottom=]", &parsers.String{}),
		args.New("[top=]", &parsers.String{}),
	)
	cmd.Use(meme)

	Category.Add(cmd)
}

func meme(ctx *ocl.Context, next ocl.Next) {
	img, t := getAttachment(ctx.Msg, ctx.Session, next)

	if img == nil {
		ctx.Send("no attachment found")
		return
	}

	r := img.Bounds()
	width := r.Dx()
	height := r.Dy()
	fontSize := float64(((height + width) / 2) / 8)
	stroke := fontSize / 16

	dc := gg.NewContext(width, height)
	dc.DrawImage(img, 0, 0)
	dc.LoadFontFace("assets/fonts/impact.ttf", fontSize)

	bottomtext := ctx.Args["bottom"].(string)
	toptext := ctx.Args["top"].(string)

	if bottomtext != "" {
		x := float64(width / 2)
		y := float64(height) - fontSize

		drawTextWithStroke(dc, bottomtext, x, y, fontSize, stroke, true)
	}
	if toptext != "" {
		x := float64(width / 2)
		y := fontSize

		drawTextWithStroke(dc, toptext, x, y, fontSize, stroke, false)
	}

	var buf bytes.Buffer
	w := io.MultiWriter(&buf)
	dc.EncodePNG(w)

	ctx.SendComplex(&discordgo.MessageSend{
		Files: []*discordgo.File{
			{
				Name:        "hello." + t,
				ContentType: "image/" + t,
				Reader:      bytes.NewReader(buf.Bytes()),
			},
		},
	})
}

func drawTextWithStroke(dc *gg.Context, text string, x, y, fontSize, stroke float64, down bool) {
	dc.SetRGB(0, 0, 0)

	// literally shamelessly stolen from https://github.com/fogleman/gg/blob/master/examples/meme.go
	for dy := -stroke; dy <= stroke; dy++ {
		for dx := -stroke; dx <= stroke; dx++ {
			if dx*dx+dy*dy >= stroke*stroke {
				// give it rounded corners
				continue
			}
			x := x + float64(dx)
			y := y + float64(dy)
			drawText(dc, text, x, y, fontSize, down)
		}
	}

	dc.SetRGB(1, 1, 1)
	drawText(dc, text, x, y, fontSize, down)
}

func drawText(dc *gg.Context, text string, x, y, fontSize float64, up bool) {
	texts := dc.WordWrap(text, float64(dc.Width()))
	if up {
		y -= float64(len(texts)-1) * dc.FontHeight()
	}
	for i, t := range texts {
		newy := y + float64(i)*dc.FontHeight()
		dc.DrawStringAnchored(t, x, newy, 0.5, 0.5)
	}
}

func getAttachment(m *discordgo.Message, s *discordgo.Session, next ocl.Next) (img image.Image, t string) {
	url := ""
	if len(m.Attachments) == 0 {
		mr := m.MessageReference
		if mr != nil {
			msg, err := s.ChannelMessage(mr.ChannelID, mr.MessageID)
			if err != nil {
				next(err)
			}
			img, t = getAttachment(msg, s, next)
		}
		rxStrict := xurls.Strict()
		url = rxStrict.FindString(m.Content)
		if url == "" {
			return
		}
	} else {
		url = m.Attachments[0].URL
	}

	resp, err := http.Get(url)
	if err != nil {
		next(err)
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		next(err)
	}

	contentType := http.DetectContentType(body)

	r := bytes.NewReader(body)

	t = "png"
	switch contentType {
	case "image/jpeg":
		img, err = jpeg.Decode(r)
		t = "jpeg"
	case "image/png":
		img, err = png.Decode(r)
	case "image/gif":
		img, err = gif.Decode(r)
		t = "gif"
	}

	if err != nil {
		next(err)
	}

	return
}
