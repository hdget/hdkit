package cmd

import (
	"bytes"
	"fmt"
	"github.com/hdget/hdkit/data"
	"github.com/hdget/hdkit/g"
	"github.com/spf13/cobra"
	iofs "io/fs"
	"os"
	"os/exec"
	"path"
	"runtime"
	"strings"
)

// newCmd represents the new command
var newCmd = &cobra.Command{
	Use:   "new",
	Short: "new project",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) < 1 {
			fmt.Println("You must provide a name for the project")
			return
		}

		rootDir := getRootDir(args[0])
		err := newProject(rootDir)
		if err != nil {
			fmt.Println("create project, err:", err)
			return
		}

		fmt.Printf("Successfully create project: %s\n\n", args[0])
		if cliPbPath == "" {
			switch runtime.GOOS {
			case "windows":
				fmt.Println(data.MsgWinSetup)
			default:
				fmt.Println(data.MsgLinuxSetup)
			}
		}
	},
}

func newProject(rootDir string) error {
	exists, err := g.GetFs().Exists(rootDir)
	if err != nil {
		return err
	}
	if exists {
		return fmt.Errorf("project %s already exist", rootDir)
	}

	// create project dirs
	for _, dir := range g.GetProjectDirs(rootDir) {
		err := g.GetFs().MakeDir(dir)
		if err != nil {
			return err
		}
	}

	err = createScriptFile(g.GetDir(rootDir, g.Binary))
	if err != nil {
		return err
	}

	err = copy3rdProtoFiles(rootDir, g.GetDir(rootDir, g.Proto))
	if err != nil {
		return err
	}

	err = createGoModuleFile(rootDir)
	if err != nil {
		return err
	}

	if cliPbPath != "" {
		return processProtoFiles(rootDir, cliPbPath)
	}
	return nil
}

// create script files under `bin` dir
func createScriptFile(binaryDir string) error {
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
func copy3rdProtoFiles(rootDir, protoDir string) error {
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
		newDir := path.Join(rootDir, subdir)
		if err := g.GetFs().MakeDir(newDir); err != nil {
			return err
		}
	}

	// 拷贝proto文件
	for _, f := range protofiles {
		destPath := path.Join(rootDir, f)
		bs, _ := data.ProtoFs.ReadFile(f)
		if err := g.GetFs().WriteFile(destPath, bs, true); err != nil {
			return err
		}
	}

	return nil
}

// create go.mod file
func createGoModuleFile(projectName string) error {
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
