package cmd

import (
	"fmt"
	"github.com/hdget/hdkit/pkg/cmdgen/dapr"
	"github.com/spf13/cobra"
	"os"
)

var genDaprCodesCmd = &cobra.Command{
	Use:   "dapr",
	Short: "generate dapr codes",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			fmt.Println("You must provide project root dir")
			return
		}

		rootDir := getRootDir(args[0])

		factory, err := dapr.NewDaprFileFactory(rootDir)
		if err != nil {
			fmt.Printf("Error new dapr file factory: %v\n\n", err)
			os.Exit(1)
		}

		err = factory.Create()
		if err != nil {
			fmt.Printf("Error generating dapr codes: %v\n", err)
			os.Exit(1)
		}
	},
}
