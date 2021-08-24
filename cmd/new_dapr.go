package cmd

import (
	"fmt"
	"github.com/hdget/hdkit/pkg"
	"github.com/hdget/hdkit/pkg/cmdnew/dapr"
	"github.com/spf13/cobra"
)

// newDaprCmd generate dapr based code boilerplate
var newDaprCmd = &cobra.Command{
	Use:     "dapr",
	Short:   "create DAPR based project",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) < 1 {
			fmt.Println("You must provide a name for the project")
			return
		}

		rootDir := getRootDir(args[0])
		factory, err := dapr.NewDaprProject(rootDir)
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
			err = pkg.ProcessProtoFiles(rootDir, cliProtoPath)
			if err != nil {
				fmt.Println("create project, err:", err)
				return
			}
		} else {
			pkg.RemindGeneratePb()
		}

		fmt.Printf("Successfully create DAPR project: %s\n\n", args[0])
	},
}


