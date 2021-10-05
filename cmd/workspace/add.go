package workspace

import "github.com/urfave/cli/v2"

var addCommand = &cli.Command{
	Name:    "add",
	Aliases: []string{"a"},
	Usage:   "Add project into the workspace.",
	Action:  addWorkspace,
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:  "git",
			Usage: "Specify the project git repo.",
		},
		&cli.StringFlag{
			Name:  "name",
			Usage: "Specify the project name.",
		},
	},
}

func addWorkspace(c *cli.Context) error {

	return nil
}
