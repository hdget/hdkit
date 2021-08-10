package cmd

import (
	"github.com/spf13/cobra"
	"os"
)

var generateCmd = &cobra.Command{
	Use:     "gen",
	Short:   "generate service/middleware/client",
	Aliases: []string{"g"},
	Run: func(cmd *cobra.Command, args []string) {
		if err := cmd.Help(); err != nil {
			os.Exit(1)
		}
	},
}

func init() {
	generateCmd.AddCommand(genPbCmd)
	generateCmd.AddCommand(genClientCmd)
	generateCmd.AddCommand(genServiceCmd)
}
