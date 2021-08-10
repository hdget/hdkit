package cmd

import (
	_ "embed"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"os"
	"strings"
)

// pb path specified in cli mode
var cliPbPath string

// RootCmd is the root command of kit
var RootCmd = &cobra.Command{
	Use:   "hdkit",
	Short: "micro service kit",
	Run: func(cmd *cobra.Command, args []string) {
		if err := cmd.Help(); err != nil {
			os.Exit(1)
		}
	},
}

// Execute runs the root command
func Execute() {
	if err := RootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

//nolint:errcheck
func init() {
	RootCmd.AddCommand(newCmd)
	RootCmd.AddCommand(generateCmd)

	RootCmd.PersistentFlags().StringVarP(&cliPbPath, "pb_path", "p", "", "Specify protobuf filepath")

	RootCmd.PersistentFlags().BoolP("debug", "d", false, "If you want to see the debug logs.")
	RootCmd.PersistentFlags().BoolP("force", "f", false, "Force overide existing files without asking.")
	RootCmd.PersistentFlags().StringP("folder", "b", "", "If you want to specify the base folder of the project.")

	viper.BindPFlag("folder", RootCmd.PersistentFlags().Lookup("folder"))
	viper.BindPFlag("force", RootCmd.PersistentFlags().Lookup("force"))
	viper.BindPFlag("debug", RootCmd.PersistentFlags().Lookup("debug"))
}

func getRootDir(name string) string {
	return strings.ToLower(name)
}
