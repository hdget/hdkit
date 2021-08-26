package cmd

import (
	"github.com/spf13/cobra"
	"os"
)

type GenerateFactory interface {
}

var (
	cliFramework int
)

var generateCmd = &cobra.Command{
	Use:   "gen",
	Short: "generate service codes",
	Run: func(cmd *cobra.Command, args []string) {
		if err := cmd.Help(); err != nil {
			os.Exit(1)
		}
	},
}

func init() {
	generateCmd.AddCommand(genPbCmd)
	generateCmd.AddCommand(genDaprCodesCmd)
	generateCmd.AddCommand(genGokitCodesCmd)

	generateCmd.PersistentFlags().IntVarP(&cliFramework, "framework", "t", int(FRAMEWORK_DAPR), "Specify framework type")
}
