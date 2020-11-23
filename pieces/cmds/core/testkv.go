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
	cmd.Args(args.New("<test=thing>", &parsers.String{}))
	cmd.Use(testkv)

	Category.Add(cmd)
}

func testkv(ctx *ocl.Context, _ ocl.Next) {
	fmt.Print(ctx.Args)
	test := ctx.Args["test"].(string)
	ctx.Send(test)
}