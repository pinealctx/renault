package workspace

import (
	"fmt"
	"github.com/codeskyblue/go-sh"
	"github.com/panjf2000/ants/v2"
	"github.com/pinealctx/renault/pkg/gits"
	"github.com/pinealctx/renault/pkg/paths"
	"github.com/pinealctx/renault/pkg/share"
	"github.com/urfave/cli/v2"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
	"strings"
	"sync"
	"time"
)

const (
	poolSize = 5
)

var syncCommand = &cli.Command{
	Name:    "sync",
	Aliases: []string{"a"},
	Usage:   "Sync the workspace all projects.",
	Action:  syncWorkspace,
}

func syncWorkspace(*cli.Context) error {
	var renaultPath = share.RenaultAbsolutePath()
	var exist, err = paths.Exists(renaultPath)
	if err != nil {
		return fmt.Errorf("check renault path exists error: %+v", err)
	}
	if !exist {
		fmt.Println("Workspace don't initialize.")
		return nil
	}
	projects, err := loadProjects()
	if err != nil {
		return fmt.Errorf("loadProjects error: %+v", err)
	}
	dirs, err := os.ReadDir(share.PWD)
	if err != nil {
		return fmt.Errorf("readDir error: %+v", err)
	}
	var changed bool
	var hash = make(map[string]struct{})
	for _, p := range projects {
		hash[p.Name] = struct{}{}
	}
	for _, dir := range dirs {
		if strings.HasPrefix(dir.Name(), ".") {
			continue
		}
		var _, ok = hash[dir.Name()]
		if !ok {
			var url = getGitURL(share.ProjectAbsolutePath(dir.Name()))
			if url != "" {
				projects = append(projects, Project{
					Name: getGitName(url),
					URL:  url,
				})
				changed = true
			}
		}
	}
	var wg sync.WaitGroup
	pool, err := ants.NewPoolWithFunc(poolSize, func(i interface{}) {
		wg.Done()
		syncGitProject(i.(*Project))
	})
	if err != nil {
		return fmt.Errorf("newPoolWithFunc error: %+v", err)
	}
	defer pool.Release()
	for _, p := range projects {
		wg.Add(1)
		if err = pool.Invoke(&p); err != nil {
			return fmt.Errorf("invoke task error: %+v", err)
		}
	}
	wg.Wait()
	if changed {
		if err = saveProjects(projects); err != nil {
			return fmt.Errorf("saveProjects error: %+v", err)
		}
	}
	fmt.Println("Workspace synchronization completed.")
	return nil
}

func syncGitProject(p *Project) {
	var s = sh.NewSession()
	s.SetTimeout(time.Second * 15)
	defer s.Kill(os.Kill)

	var latest, err = cloneProject(s, p)
	if err != nil {
		fmt.Printf("[%s] clone project error: %+v\n", p.Name, err)
		return
	}
	if latest {
		return
	}
	if err = pullProject(s, p); err != nil {
		fmt.Printf("[%s] pull project error: %+v\n", p.Name, err)
	}
}

func cloneProject(s *sh.Session, p *Project) (bool, error) {
	var pp = share.ProjectAbsolutePath(p.Name)
	var exists, err = paths.Exists(pp)
	if err != nil {
		return false, fmt.Errorf("cloneProject exists error: %+v", err)
	}
	if exists {
		var url = getGitURL(pp)
		if url != p.URL {
			fmt.Printf("[%s] [Warning] The git project url does not match the configured url.\n", p.Name)
		}
		return false, nil
	}
	s.SetDir(share.PWD)
	out, err := s.Command("git", "clone", p.URL, p.Name).CombinedOutput()
	if err != nil {
		return false, fmt.Errorf("git clone error: %+v", err)
	}
	fmt.Printf("[%s] git clone success: %s", p.Name, out)
	return true, nil
}

func pullProject(s *sh.Session, p *Project) error {
	var pp = share.ProjectAbsolutePath(p.Name)
	s.SetDir(pp)
	var status, err = statusProject(s, p)
	if err != nil {
		return fmt.Errorf("status project error: %+v", err)
	}

	defer func() {
		fmt.Printf("[%s] git status: %s\n", p.Name, status.Fmt())
	}()

	if !status.CanPull(true) {
		return nil
	}
	_, err = s.Command("git", "pull").CombinedOutput()
	if err != nil {
		return fmt.Errorf("git pull project error: %+v", err)
	}
	fmt.Printf("[%s] git pull success.", p.Name)
	status, err = statusProject(s, p)
	if err != nil {
		return fmt.Errorf("status project agin error: %+v", err)
	}
	status.SetNewPull()
	return nil
}

func statusProject(s *sh.Session, p *Project) (*gits.Status, error) {
	var pp = share.ProjectAbsolutePath(p.Name)
	s.SetDir(pp)
	var output, err = s.Command("git", "fetch").CombinedOutput()
	if err != nil {
		return nil, fmt.Errorf("git fetch error: %+v\n%s", err, output)
	}
	output, err = s.Command("git", "status", "--porcelain=v2", "--branch").CombinedOutput()
	if err != nil {
		return nil, fmt.Errorf("git status error: %+v", err)
	}
	var status = gits.NewStatus(pp)
	if err = status.Parse(output); err != nil {
		return nil, fmt.Errorf("parse git status error: %+v", err)
	}
	output, err = s.Command("git", "describe", "--tags").CombinedOutput()
	if err == nil {
		status.SetTag(strings.TrimRight(string(output), "\n"))
	}
	return status, nil
}

func loadProjects() ([]Project, error) {
	var buff, err = ioutil.ReadFile(share.ConfigAbsoluteFile())
	if err != nil {
		return nil, fmt.Errorf("loadProjects readFile error: %+v", err)
	}
	var projects []Project
	if err = yaml.Unmarshal(buff, &projects); err != nil {
		return nil, fmt.Errorf("loadProjects unmarshal error: %+v", err)
	}
	return projects, nil
}

func saveProjects(projects []Project) error {
	var buf, err = yaml.Marshal(projects)
	if err != nil {
		return fmt.Errorf("saveProjects marshal error: %+v", err)
	}
	if err = ioutil.WriteFile(share.ConfigAbsoluteFile(), buf, 0755); err != nil {
		return fmt.Errorf("saveProjects writeFile error: %+v", err)
	}
	return nil
}
