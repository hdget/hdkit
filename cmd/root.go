package cmd

import (
	_ "embed"
	"github.com/spf13/cobra"
	"os"
	"strings"
)

type framework int

const (
	_ framework = iota
	FRAMEWORK_DAPR
	FRAMEWORK_GOKIT
)

// pb path specified in cli mode
var (
	cliProtoPath string
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

	RootCmd.PersistentFlags().StringVarP(&cliProtoPath, "proto_path", "p", "", "Specify protobuf filepath")
}

func getRootDir(name string) string {
	return strings.ToLower(name)
}
