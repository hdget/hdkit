package cmd

import (
	"fmt"
	"github.com/hdget/hdkit/data"
	"github.com/hdget/hdkit/file"
	"github.com/hdget/hdkit/g"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"os"
)

var (
	SupportTransports = []string{"http", "grpc"}
)

var genServiceCmd = &cobra.Command{
	Use:   "service",
	Short: "generate service",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			fmt.Println("You must provide project root dir")
			return
		}

		rootDir := getRootDir(args[0])
		m, err := file.NewServiceFactory(rootDir)
		if err != nil {
			fmt.Printf("Error new service: %v\n\n", err)
			if errors.Is(err, g.ErrServiceNotFound) {
				if !checkProtoc() {
					fmt.Println(data.MsgInstallProtoc)
				}
				fmt.Println(data.MsgWinSetup)
			}
			return
		}

		err = m.Create()
		if err != nil {
			fmt.Printf("Error generating service: %v\n", err)
			os.Exit(1)
		}
	},
}

//nolint:errcheck
func init() {
	// default generate grpc transport
	genServiceCmd.Flags().StringArrayP("transports", "t", []string{"grpc"}, "The transport you want your service to be initiated with")
	genServiceCmd.Flags().StringP("pb_path", "p", "", "Specify path to store pb dir")
	genServiceCmd.Flags().StringP("pb_import_path", "i", "", "Specify path to import pb")
	genServiceCmd.Flags().StringArrayP("methods", "m", []string{}, "Specify methods to be generated")

	viper.BindPFlag("g_s_transports", genServiceCmd.Flags().Lookup("transports"))
	viper.BindPFlag("g_s_methods", genServiceCmd.Flags().Lookup("methods"))
	viper.BindPFlag("g_s_pb_path", genServiceCmd.Flags().Lookup("pb_path"))
	viper.BindPFlag("g_s_pb_import_path", genServiceCmd.Flags().Lookup("pb_import_path"))
}
