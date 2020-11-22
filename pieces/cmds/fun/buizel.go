package fun

import (
	"io/ioutil"
	"net/http"

	"github.com/opolobot/Opolo/ocl"
)

func init() {
	cmd := ocl.New()
	cmd.Name("buizel")
	cmd.Aliases("bui")
	cmd.Description("Grab a random gif of a buizel")
	cmd.Use(buizel)

	Category.Add(cmd)
}

func buizel(ctx *ocl.Context, next ocl.Next) {
	resp, err := http.Get("https://b.piapiac.org/api")
	if err != nil {
		next(err)
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		next(err)
	}

	img := string(body)
	ctx.Send(img)
}
