package workspace

import (
	"fmt"
	"github.com/pinealctx/renault/pkg/paths"
	"github.com/pinealctx/renault/pkg/regex"
	"github.com/pinealctx/renault/pkg/share"
	"github.com/urfave/cli/v2"
	"io/ioutil"
	"os"
	"path"
	"strings"
)

const (
	GitConfigPath = ".git/config"
)

var initCommand = &cli.Command{
	Name:    "init",
	Aliases: []string{"i"},
	Usage:   "Initialize the workspace.",
	Action:  initWorkspace,
}

func initWorkspace(*cli.Context) error {
	var renaultPath = share.RenaultAbsolutePath()
	var exist, err = paths.Exists(renaultPath)
	if err != nil {
		return fmt.Errorf("check renault path exists error: %+v", err)
	}
	if exist {
		fmt.Println("Workspace already exists.")
		return nil
	}
	if err = os.Mkdir(renaultPath, 0755); err != nil {
		return fmt.Errorf("make renault dir error: %+v", err)
	}
	dirs, err := os.ReadDir(share.PWD)
	if err != nil {
		return fmt.Errorf("readDir error: %+v", err)
	}
	var projects = make([]Project, 0, len(dirs))
	for _, dir := range dirs {
		if strings.HasPrefix(dir.Name(), ".") {
			continue
		}
		var url = getGitURL(share.ProjectAbsolutePath(dir.Name()))
		if url != "" {
			projects = append(projects, Project{
				Name: getGitName(url),
				URL:  url,
			})
		}
	}
	if err = saveProjects(projects); err != nil {
		return fmt.Errorf("saveProjects error: %+v", err)
	}
	fmt.Println("Initialization of workspace completed.")
	return nil
}

type Project struct {
	Name string `yaml:"name"`
	URL  string `yaml:"url"`
}

func getGitURL(dir string) string {
	gitConfigPath := path.Join(dir, GitConfigPath)
	exist, err := paths.Exists(gitConfigPath)
	if err != nil || !exist {
		return ""
	}
	config, err := ioutil.ReadFile(gitConfigPath)
	if err != nil {
		return ""
	}
	uri, ok := regex.MatchStr("url = (.*?)\n", string(config), 1)
	if !ok {
		return ""
	}
	return uri
}

func getGitName(url string) string {
	var list = strings.Split(url, "/")
	var l = len(list)
	if l == 0 {
		return ""
	}
	return strings.ReplaceAll(list[l-1], ".git", "")
}
