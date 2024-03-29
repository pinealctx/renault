package workspace

import "github.com/urfave/cli/v2"

var Command = &cli.Command{
	Name:    "workspace",
	Aliases: []string{"w"},
	Usage:   "Commands related to the management workspace.",
	Subcommands: []*cli.Command{
		initCommand,
		syncCommand,
		addCommand,
	},
}
