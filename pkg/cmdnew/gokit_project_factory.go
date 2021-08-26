package cmdnew

import (
	"fmt"
	"github.com/hdget/hdkit/g"
)

type gokitFactory struct {
	baseProject
}

func NewGokitProject(rootDir string) (ProjectFactory, error) {
	exists, err := g.GetFs().Exists(rootDir)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, fmt.Errorf("project %s already exist", rootDir)
	}

	return &gokitFactory{
		baseProject: baseProject{
			name:    "gokit",
			rootDir: rootDir,
			dirs:    getGokitProjectDirs(rootDir),
		},
	}, nil
}

func (factory gokitFactory) Create() error {
	return start(factory)
}

// GetProjectDirs get default project dir
func getGokitProjectDirs(rootDir string) []string {
	dirs := getProjectCommonDirs(rootDir)
	dirs = append(dirs,
		g.GetDir(rootDir, g.Grpc),
		g.GetDir(rootDir, g.Http),
	)
	return dirs
}
