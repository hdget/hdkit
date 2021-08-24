package cmd

import (
	"fmt"
	"github.com/hdget/hdkit/g"
	"github.com/hdget/hdkit/pkg"
	"github.com/spf13/cobra"
	"os"
	"path/filepath"
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
		protoFilePaths := make([]string, 0)
		if cliProtoPath == "" {
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

			protoFilePaths = append(protoFilePaths, matches...)
		} else {
			protoFilePaths = append(protoFilePaths, cliProtoPath)
		}

		for _, filePath := range protoFilePaths {
			err := pkg.ProcessProtoFiles(rootDir, filePath)
			if err != nil {
				fmt.Printf("Error process proto file:%s, err:%v", filePath, err)
				os.Exit(1)
			}
		}
	},
}
