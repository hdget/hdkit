package dapr

import (
	"bytes"
	"fmt"
	"github.com/hdget/hdkit/data"
	"github.com/hdget/hdkit/g"
	"github.com/hdget/hdkit/pkg/cmdnew"
	iofs "io/fs"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"runtime"
	"strings"
)

type daprProject struct {
	rootDir string
}

func NewDaprProject(rootDir string) (cmdnew.ProjectFactory, error) {
	exists, err := g.GetFs().Exists(rootDir)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, fmt.Errorf("project %s already exist", rootDir)
	}

	return &daprProject{
		rootDir: rootDir,
	}, nil
}

func (factory daprProject) Create() error {
	// create project dirs
	for _, dir := range factory.getProjectDirs() {
		err := g.GetFs().MakeDir(dir)
		if err != nil {
			return err
		}
	}

	err := factory.createScriptFile()
	if err != nil {
		return err
	}

	err = factory.copy3rdProtoFiles()
	if err != nil {
		return err
	}

	err = factory.createGoModuleFile()
	if err != nil {
		return err
	}

	return nil
}


// GetProjectDirs get default project dir
func (factory daprProject) getProjectDirs() []string {
	return []string{
		g.GetDir(factory.rootDir, g.Binary),
		g.GetDir(factory.rootDir, g.Proto),
		g.GetDir(factory.rootDir, g.Global),
		g.GetDir(factory.rootDir, g.Service),
		g.GetDir(factory.rootDir, g.Pb),
		g.GetDir(factory.rootDir, g.Cmd),
	}
}

// create script files under `bin` dir
func (factory daprProject) createScriptFile() error {
	binaryDir := g.GetDir(factory.rootDir, g.Binary)

	files, err := iofs.ReadDir(data.ScriptFs, "script")
	if err != nil {
		return err
	}

	fileSuffix := ".sh"
	if runtime.GOOS == "windows" {
		fileSuffix = ".bat"
	}

	for _, f := range files {
		if strings.HasSuffix(f.Name(), fileSuffix) {
			bs, _ := data.ScriptFs.ReadFile(path.Join("script", f.Name()))
			err := g.GetFs().WriteFile(path.Join(binaryDir, f.Name()), bs, true)
			if err != nil {
				return err
			}
		}
	}

	return nil
}


// create 3rd party proto files under `proto` dir
func (factory daprProject) copy3rdProtoFiles() error {
	subdirs := make([]string, 0)
	protofiles := make([]string, 0)
	err := iofs.WalkDir(data.ProtoFs, "proto", func(path string, d iofs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if d.IsDir() {
			subdirs = append(subdirs, path)
		} else {
			if strings.HasSuffix(d.Name(), ".proto") {
				protofiles = append(protofiles, path)
			}
		}
		return nil
	})
	if err != nil {
		return err
	}

	// 创建对应子目录
	for _, subdir := range subdirs {
		newDir := path.Join(factory.rootDir, subdir)
		if err := g.GetFs().MakeDir(newDir); err != nil {
			return err
		}
	}

	// 拷贝proto文件
	for _, f := range protofiles {
		destPath := path.Join(factory.rootDir, f)
		bs, _ := data.ProtoFs.ReadFile(f)
		if err := g.GetFs().WriteFile(destPath, bs, true); err != nil {
			return err
		}
	}

	return nil
}


// create go.mod file
func (factory daprProject) createGoModuleFile() error {
	projectName := filepath.Base(factory.rootDir)

	exist, _ := g.GetFs().Exists(path.Join(projectName, "go.mod"))
	if exist {
		return nil
	}

	var cmds []string
	switch runtime.GOOS {
	case "windows":
		cmds = []string{
			"cmd", "/c", fmt.Sprintf("go mod init %s", strings.ToLower(projectName)),
		}
	default:
		cmds = []string{
			"sh", "-c", fmt.Sprintf("go mod init %s", strings.ToLower(projectName)),
		}
	}

	// obtain current working Dir
	baseDir, err := os.Getwd()
	if err != nil {
		return err
	}

	var stderr bytes.Buffer
	cmd := exec.Command(cmds[0], cmds[1:]...)
	cmd.Dir = path.Join([]string{baseDir, projectName}...)
	cmd.Stderr = &stderr
	_, err = cmd.Output()
	if err != nil {
		return fmt.Errorf("create go.mod: %s, err: %v", strings.Join(cmds, " "), err.Error()+" , "+stderr.String())
	}

	return nil
}
