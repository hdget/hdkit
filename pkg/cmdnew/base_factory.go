package cmdnew

import (
	"bytes"
	"fmt"
	"github.com/hdget/hdkit/data"
	"github.com/hdget/hdkit/g"
	"github.com/hdget/hdkit/utils"
	iofs "io/fs"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"runtime"
	"strings"
)

type baseProject struct {
	name    string
	rootDir string
	dirs    []string
}

func start(realFactory ProjectFactory) error {
	// create project dirs
	err := realFactory.createProjectDirs()
	if err != nil {
		return err
	}

	err = realFactory.copyScriptFiles()
	if err != nil {
		return err
	}

	err = realFactory.copySettingFiles()
	if err != nil {
		return err
	}

	err = realFactory.copy3rdProtoFiles()
	if err != nil {
		return err
	}

	err = realFactory.createGoModuleFile()
	if err != nil {
		return err
	}

	return nil
}

func (factory baseProject) createProjectDirs() error {
	// create project dirs
	for _, dir := range factory.dirs {
		err := g.GetFs().MakeDir(dir)
		if err != nil {
			return err
		}
	}
	return nil
}

// create script files under `bin` dir
func (factory baseProject) copyScriptFiles() error {
	files, err := iofs.ReadDir(data.ScriptFs, path.Join("script", factory.name))
	if err != nil {
		return err
	}

	fileSuffix := ".sh"
	if runtime.GOOS == "windows" {
		fileSuffix = ".bat"
	}

	binaryDir := g.GetDir(factory.rootDir, g.Binary)
	for _, f := range files {
		if strings.HasSuffix(f.Name(), fileSuffix) {
			bs, _ := data.ScriptFs.ReadFile(path.Join("script", factory.name, f.Name()))
			err := g.GetFs().WriteFile(path.Join(binaryDir, f.Name()), bs, true)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

// create setting files under `settings` dir
func (factory baseProject) copySettingFiles() error {
	sourceSettingDir := path.Join("setting", factory.name)
	targetSettingDir := g.GetDir(factory.rootDir, g.Setting)
	dirs, files, err := utils.TraverseDirFiles(data.SettingFs, sourceSettingDir)
	if err != nil {
		return err
	}

	// 创建对应子目录
	for _, subdir := range dirs {
		relativePath, _ := filepath.Rel(sourceSettingDir, subdir)
		if relativePath != "" {
			newDir := path.Join(targetSettingDir, relativePath)
			if err := g.GetFs().MakeDir(newDir); err != nil {
				return err
			}
		}
	}

	for _, f := range files {
		bs, _ := data.SettingFs.ReadFile(f)

		relativePath, _ := filepath.Rel(sourceSettingDir, f)
		destPath := path.Join(targetSettingDir, relativePath)

		// app.local.toml ==> app, local, toml
		// try change app.local.toml to <app>.local.toml
		tokens := strings.Split(filepath.Base(destPath), ".")
		if len(tokens) > 2 && tokens[0] == "app" {
			destPath = path.Join(filepath.Dir(destPath), fmt.Sprintf("%s.%s.%s", factory.rootDir, tokens[1], tokens[2]))
		}
		err := g.GetFs().WriteFile(destPath, bs, true)
		if err != nil {
			return err
		}
	}

	return nil
}

// create 3rd party proto files under `proto` dir
func (factory baseProject) copy3rdProtoFiles() error {
	dirs, files, err := utils.TraverseDirFiles(data.ProtoFs, "proto")
	if err != nil {
		return err
	}

	// 创建对应子目录
	for _, subdir := range dirs {
		newDir := path.Join(factory.rootDir, subdir)
		if err := g.GetFs().MakeDir(newDir); err != nil {
			return err
		}
	}

	// 拷贝proto文件
	for _, f := range files {
		if strings.HasSuffix(f, ".proto") {
			destPath := path.Join(factory.rootDir, f)
			bs, _ := data.ProtoFs.ReadFile(f)
			if err := g.GetFs().WriteFile(destPath, bs, true); err != nil {
				return err
			}
		}
	}
	return nil
}

// create go.mod file
func (factory baseProject) createGoModuleFile() error {
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

// GetProjectDirs get default project dir
func getProjectCommonDirs(rootDir string) []string {
	return []string{
		g.GetDir(rootDir, g.Binary),
		g.GetDir(rootDir, g.Proto),
		g.GetDir(rootDir, g.Global),
		g.GetDir(rootDir, g.Service),
		g.GetDir(rootDir, g.Pb),
		g.GetDir(rootDir, g.Cmd),
		g.GetDir(rootDir, g.Setting),
	}
}
