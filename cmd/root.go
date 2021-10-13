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

// rootCmd is the root command of kit
var rootCmd = &cobra.Command{
	Use:   "hdkit",
	Short: "hd kit",
	Run: func(cmd *cobra.Command, args []string) {
		if err := cmd.Help(); err != nil {
			os.Exit(1)
		}
	},
}

// Execute runs the root command
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

//nolint:errcheck
func init() {
	rootCmd.AddCommand(newCmd)
	rootCmd.AddCommand(generateCmd)
	rootCmd.AddCommand(compileProtoCmd)

	rootCmd.PersistentFlags().StringVarP(&cliProtoPath, "proto-path", "p", "", "Specify protobuf filepath")
	rootCmd.PersistentFlags().StringVarP(&cliFramework, "framework", "t", "dapr", "Specify framework type, e,g: dapr, gokit")
}

func getRootDir(name string) string {
	return strings.ToLower(name)
}
