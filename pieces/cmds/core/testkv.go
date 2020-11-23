package core

import (
	"fmt"

	"github.com/opolobot/Opolo/ocl"
	"github.com/opolobot/Opolo/ocl/args"
	"github.com/opolobot/Opolo/pieces/parsers"
)

func init() {
	cmd := ocl.New()
	cmd.Name("testkv")
	cmd.Args(args.New("[...nom_nom_greedy]", &parsers.String{}), args.New("<test=thing>", &parsers.String{}))
	cmd.Use(testkv)

	Category.Add(cmd)
}

func testkv(ctx *ocl.Context, _ ocl.Next) {
	ctx.Send(fmt.Sprint(ctx.Args))
}