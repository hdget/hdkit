package cmd

import (
	"fmt"
	"github.com/hdget/hdkit/pkg/cmdnew"
	"github.com/hdget/hdkit/pkg/common"
	"github.com/spf13/cobra"
)

// newGokitCmd generate gokit based code boilerplate
var newGokitCmd = &cobra.Command{
	Use:   "gokit",
	Short: "create GoKit based project",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) < 1 {
			fmt.Println("You must provide a name for the project")
			return
		}

		rootDir := getRootDir(args[0])
		factory, err := cmdnew.NewGokitProject(rootDir)
		if err != nil {
			fmt.Println("new project factory, err:", err)
			return
		}

		err = factory.Create()
		if err != nil {
			fmt.Println("create project, err:", err)
			return
		}

		if cliProtoPath != "" {
			err = common.ProcessProtoFiles(rootDir, cliProtoPath)
			if err != nil {
				fmt.Println("create project, err:", err)
				return
			}
		} else {
			common.RemindGeneratePb()
		}

		fmt.Printf("Successfully create GoKit project: %s\n\n", args[0])
	},
}
