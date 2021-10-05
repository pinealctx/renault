package share

import "path"

const (
	RenaultPath              = ".renault"
	RenaultProjectConfigPath = "project.yaml"
)

var (
	PWD string
)

func RenaultAbsolutePath() string {
	return path.Join(PWD, RenaultPath)
}

func ConfigAbsoluteFile() string {
	return path.Join(PWD, RenaultPath, RenaultProjectConfigPath)
}

func ProjectAbsolutePath(p string) string {
	return path.Join(PWD, p)
}
