package cmd

import (
	"fmt"
	"github.com/hdget/hdkit/pkg/cmdgen/dapr"
	"github.com/hdget/hdkit/pkg/cmdgen/gokit"
	"github.com/spf13/cobra"
	"os"
)

type GenerateFactory interface {
}

var generateCmd = &cobra.Command{
	Use:   "gen",
	Short: "generate service codes",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			fmt.Println("You must provide project root dir")
			return
		}

		rootDir := getRootDir(args[0])
		switch cliFramework {
		case "gokit":
			genGokitCodes(rootDir)
		default:
			genDaprCodes(rootDir)
		}
	},
}

func genDaprCodes(rootDir string) {
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
}

func genGokitCodes(rootDir string) {
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
}
