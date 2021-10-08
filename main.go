package main

import (
	"fmt"
	"github.com/pinealctx/renault/cmd/project"
	"github.com/pinealctx/renault/cmd/workspace"
	"github.com/pinealctx/renault/pkg/share"
	"github.com/urfave/cli/v2"
	"os"
)

func main() {
	var app = cli.NewApp()
	app.Name = "Renault"
	app.Usage = "A useful tools."
	app.Version = "0.0.2"
	app.Before = beforeAction
	app.After = afterAction
	app.Flags = []cli.Flag{
		&cli.StringFlag{
			Name:    "workspace",
			Aliases: []string{"w"},
			Usage:   "Specify the workspace.",
		},
	}
	app.Commands = cli.Commands{
		project.Command,
		workspace.Command,
	}
	var err = app.Run(os.Args)
	if err != nil {
		fmt.Println(err)
	}
}

func beforeAction(c *cli.Context) error {
	var w = c.String("workspace")
	if w != "" {
		share.PWD = w
		return nil
	}
	var pwd, err = os.Getwd()
	if err != nil {
		return fmt.Errorf("getwd error: %+v", err)
	}
	share.PWD = pwd
	return nil
}

func afterAction(*cli.Context) error {
	return nil
}
