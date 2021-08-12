package cmd

import (
	"bytes"
	"fmt"
	"github.com/hdget/hdkit/data"
	"github.com/hdget/hdkit/g"
	"github.com/spf13/cobra"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"runtime"
)

var genPbCmd = &cobra.Command{
	Use:   "pb",
	Short: "generate pb files",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			fmt.Println("You must provide project dir")
			return
		}

		rootDir := getRootDir(args[0])
		// if we don't specify proto filepath, try check `<project>/proto` dir
		if cliPbPath == "" {
			protoDir := g.GetDir(rootDir, g.Proto)
			matches, err := filepath.Glob(filepath.Join(protoDir, "*.proto"))
			if err != nil {
				fmt.Printf("Check proto files in: %s, err:%v\n", protoDir, err)
				return
			}
			if len(matches) == 0 {
				fmt.Printf("no proto files in: %s\n", protoDir)
				return
			}

			for _, fp := range matches {
				err := processProtoFiles(rootDir, fp)
				if err != nil {
					fmt.Printf("Error process proto file:%s, err:%v", fp, err)
					os.Exit(1)
				}
			}
		} else {
			err := processProtoFiles(rootDir, cliPbPath)
			if err != nil {
				fmt.Println("Error generate pb files:", err)
				os.Exit(1)
			}
		}
	},
}

func processProtoFiles(rootDir, protoFilePath string) error {
	exists, err := g.GetFs().Exists(protoFilePath)
	if err != nil {
		return err
	}
	if !exists {
		return fmt.Errorf("protobuf filepath: %s not found", protoFilePath)
	}

	// copy proto file to `proto` directory
	protoDir := g.GetDir(rootDir, g.Proto)
	isDir, _ := g.GetFs().IsDir(protoFilePath)
	if isDir {
		matched, err := filepath.Glob(path.Join(protoFilePath, "*.proto"))
		if err != nil {
			return fmt.Errorf("search *.proto under path: %s", protoFilePath)
		}
		if len(matched) == 0 {
			return fmt.Errorf("no *.proto files under path: %s", protoFilePath)
		}

		for _, m := range matched {
			if err := copyFile(m, protoDir); err != nil {
				return err
			}
		}
	} else {
		if err := copyFile(protoFilePath, protoDir); err != nil {
			return err
		}
	}

	if !checkProtoc() {
		fmt.Println(data.MsgInstallProtoc)
		return nil
	}

	// if we found protoc, try to comple .proto files
	err = compileProto(g.GetDir(rootDir, g.Binary), g.GetDir(rootDir, g.Pb))
	if err != nil {
		return err
	}
	return nil
}

func copyFile(srcPath, destDir string) error {
	filename := filepath.Base(srcPath)
	destPath := path.Join(destDir, filename)
	if filepath.Dir(destPath) != srcPath {
		err := g.GetFs().Copy(srcPath, destPath)
		if err != nil {
			return err
		}
	}
	return nil
}

// check if we install protoc
func checkProtoc() bool {
	foundProtoc := false
	if p := exec.Command("protoc"); p.Run() == nil {
		foundProtoc = true
	}

	return foundProtoc
}

// create go.mod file
func compileProto(binDir, pbPath string) error {
	exist, _ := g.GetFs().Exists(pbPath)
	if !exist {
		return fmt.Errorf("proto file: %s doesn't exist", pbPath)
	}

	var cmds []string
	switch runtime.GOOS {
	case "windows":
		cmds = []string{
			"cmd", "/c", "gen_grpc.bat",
		}
	default:
		cmds = []string{
			"sh", "-c", "gen_gprc.sh",
		}
	}

	var stderr bytes.Buffer
	cmd := exec.Command(cmds[0], cmds[1:]...)
	cmd.Dir = binDir
	cmd.Stderr = &stderr
	stdout, err := cmd.Output()
	if err != nil {
		return fmt.Errorf("compile proto file: %s, err:%v", pbPath, err)
	}
	fmt.Println(string(stdout))
	return nil
}
