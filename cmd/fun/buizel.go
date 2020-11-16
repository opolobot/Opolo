package fun

import (
	"io/ioutil"
	"net/http"

	"github.com/zorbyte/whiskey/cmd"
)

func init() {
	cmd := cmd.New()
	cmd.Name("buizel")
	cmd.Aliases("bui")
	cmd.Description("Grab a random gif of a buizel")
	cmd.Use(bui)

	Category.AddCommand(cmd.Command())
}

func bui(ctx *cmd.Context, next cmd.NextFunc) {
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
