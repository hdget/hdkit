package cmd

import (
	"github.com/spf13/cobra"
)

// newCmd represents the new command
var newCmd = &cobra.Command{
	Use:   "new",
	Short: "new project",
}

func init() {
	newCmd.AddCommand(newDaprCmd)
	newCmd.AddCommand(newGokitCmd)
}
