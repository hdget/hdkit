package cmd

import (
	"fmt"
	"github.com/hdget/hdkit/pkg/cmdnew"
	"github.com/hdget/hdkit/pkg/common"
	"github.com/spf13/cobra"
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
		switch cliFramework {
		case "gokit":
			newGokitProject(rootDir)
		default:
			newDaprProject(rootDir)
		}
	},
}

func newDaprProject(rootDir string) {
	factory, err := cmdnew.NewDaprProject(rootDir)
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

	fmt.Printf("Successfully create DAPR project: %s\n\n", rootDir)
}

func newGokitProject(rootDir string) {
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

	fmt.Printf("Successfully create GoKit project: %s\n\n", rootDir)
}
