package cmd

import (
	"fmt"
	"github.com/hdget/hdkit/pkg/cmdgen/gokit"
	"github.com/spf13/cobra"
	"os"
)

var (
	SupportTransports = []string{"http", "grpc"}
)

var genGokitCodesCmd = &cobra.Command{
	Use:   "gokit",
	Short: "generate gokit codes",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			fmt.Println("You must provide project root dir")
			return
		}

		rootDir := getRootDir(args[0])

		factory, err := gokit.NewGokitFileFactory(rootDir)
		if err != nil {
			fmt.Printf("Error new gokit file factory: %v\n\n", err)
			os.Exit(1)
		}

		err = factory.Create()
		if err != nil {
			fmt.Printf("Error generating gokit codes: %v\n", err)
			os.Exit(1)
		}
	},
}
