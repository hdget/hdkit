package cmdnew

import (
	"fmt"
	"github.com/hdget/hdkit/g"
)

type daprProject struct {
	baseProject
}

func NewDaprProject(rootDir string) (ProjectFactory, error) {
	exists, err := g.GetFs().Exists(rootDir)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, fmt.Errorf("project %s already exist", rootDir)
	}

	return &daprProject{
		baseProject: baseProject{
			name:    "dapr",
			rootDir: rootDir,
			dirs:    getProjectCommonDirs(rootDir),
		},
	}, nil
}

func (factory daprProject) Create() error {
	return start(factory)
}
