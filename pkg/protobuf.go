package pkg

import (
	"bytes"
	"fmt"
	"github.com/hdget/hdkit/data"
	"github.com/hdget/hdkit/g"
	"os/exec"
	"path"
	"path/filepath"
	"runtime"
)

func RemindGeneratePb() {
	switch runtime.GOOS {
	case "windows":
		fmt.Println(data.MsgWinSetup)
	default:
		fmt.Println(data.MsgLinuxSetup)
	}
}

func ProcessProtoFiles(rootDir, protoPath string) error {
	exists, err := g.GetFs().Exists(protoPath)
	if err != nil {
		return err
	}
	if !exists {
		return fmt.Errorf("protobuf filepath: %s not found", protoPath)
	}

	// copy proto file to `proto` directory
	protoDir := g.GetDir(rootDir, g.Proto)
	isDir, _ := g.GetFs().IsDir(protoPath)
	if isDir {
		matched, err := filepath.Glob(path.Join(protoPath, "*.proto"))
		if err != nil {
			return fmt.Errorf("search *.proto under path: %s", protoPath)
		}
		if len(matched) == 0 {
			return fmt.Errorf("no *.proto files under path: %s", protoPath)
		}

		for _, m := range matched {
			if err := copyFile(m, protoDir); err != nil {
				return err
			}
		}
	} else {
		if err := copyFile(protoPath, protoDir); err != nil {
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
	// copy proto files
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
