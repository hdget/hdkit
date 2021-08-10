package cmd

import (
	"bytes"
	"fmt"
	"github.com/hdget/hdkit/data"
	"github.com/hdget/hdkit/g"
	"github.com/spf13/cobra"
	"os"
	"os/exec"
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
				err := processProtoFile(rootDir, fp)
				if err != nil {
					fmt.Printf("Error process proto file:%s, err:%v", fp, err)
					os.Exit(1)
				}
			}
		} else {
			err := processProtoFile(rootDir, cliPbPath)
			if err != nil {
				fmt.Println("Error generate pb files:", err)
				os.Exit(1)
			}
		}
	},
}

func processProtoFile(rootDir, protoFilePath string) error {
	exists, err := g.GetFs().Exists(protoFilePath)
	if err != nil {
		return err
	}
	if !exists {
		return fmt.Errorf("protobuf file: %s not found", protoFilePath)
	}

	// copy proto file to `proto` directory
	protoDir := g.GetDir(rootDir, g.Proto)
	filename := filepath.Base(protoFilePath)
	destProtoPath := filepath.Join(protoDir, filename)
	if filepath.Dir(protoFilePath) != protoDir {
		err = g.GetFs().Copy(protoFilePath, destProtoPath)
		if err != nil {
			return err
		}
	}

	if !checkProtoc() {
		fmt.Println(data.MsgInstallProtoc)
		return nil
	}

	// if we found protoc, try to comple .proto files
	err = compileProto(g.GetDir(rootDir, g.Binary), destProtoPath)
	if err != nil {
		return err
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
func compileProto(binDir, protoPath string) error {
	exist, _ := g.GetFs().Exists(protoPath)
	if !exist {
		return fmt.Errorf("proto file: %s doesn't exist", protoPath)
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
		return fmt.Errorf("compile proto file: %s, err:%v", protoPath, err)
	}
	fmt.Println(string(stdout))
	return nil
}
