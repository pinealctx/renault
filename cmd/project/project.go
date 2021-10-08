package project

import "github.com/urfave/cli/v2"

var Command = &cli.Command{
	Name:    "project",
	Aliases: []string{"p"},
	Usage:   "Commands related to the management project.",
	Subcommands: []*cli.Command{
		initCommand,
	},
}
