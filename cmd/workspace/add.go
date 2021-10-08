package workspace

import (
	"fmt"
	"github.com/pinealctx/renault/pkg/paths"
	"github.com/pinealctx/renault/pkg/share"
	"github.com/urfave/cli/v2"
	"os"
)

var addCommand = &cli.Command{
	Name:    "add",
	Aliases: []string{"a"},
	Usage:   "Add project into the workspace.",
	Action:  addWorkspace,
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:     "url",
			Usage:    "Specify the project git repo url.",
			Required: true,
		},
		&cli.StringFlag{
			Name:  "name",
			Usage: "Specify the project name.",
		},
	},
}

func addWorkspace(c *cli.Context) error {
	var renaultPath = share.RenaultAbsolutePath()
	var exist, err = paths.Exists(renaultPath)
	if err != nil {
		return fmt.Errorf("check renault path exists error: %+v", err)
	}
	if !exist {
		fmt.Println("Workspace don't initialize.")
		return nil
	}
	var url = c.String("url")
	var name = c.String("name")
	if url == "" {
		return fmt.Errorf("git repo url must be invalid")
	}
	if name == "" {
		name = getGitName(url)
	}
	dirs, err := os.ReadDir(share.PWD)
	if err != nil {
		return fmt.Errorf("readDir error: %+v", err)
	}
	for _, dir := range dirs {
		if dir.Name() == name {
			return fmt.Errorf("project name already exists")
		}
	}
	projects, err := loadProjects()
	if err != nil {
		return fmt.Errorf("loadProjects error: %+v", err)
	}
	for _, project := range projects {
		if project.Name == name {
			return fmt.Errorf("project name already exists")
		}
	}
	projects = append(projects, Project{
		Name: name,
		URL:  url,
	})
	if err = saveProjects(projects); err != nil {
		return fmt.Errorf("saveProjects error: %+v", err)
	}
	return nil
}
