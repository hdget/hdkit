package cmd

import (
	_ "embed"
	"github.com/spf13/cobra"
	"os"
	"strings"
)

var (
	// pb path specified in cli mode
	cliProtoPath string
	// new/gen framework type
	cliFramework string
)

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
	RootCmd.AddCommand(compileProtoCmd)

	RootCmd.PersistentFlags().StringVarP(&cliProtoPath, "proto-path", "p", "", "Specify protobuf filepath")
	generateCmd.PersistentFlags().StringVarP(&cliFramework, "framework", "t", "dapr", "Specify framework type, e,g: dapr, gokit")
}

func getRootDir(name string) string {
	return strings.ToLower(name)
}
