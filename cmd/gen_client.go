package cmd

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// genClientCmd represents the client command
var genClientCmd = &cobra.Command{
	Use:     "client",
	Short:   "Create simple client lib",
	Aliases: []string{"c"},
	Run: func(cmd *cobra.Command, args []string) {
		//if len(args) == 0 {
		//	logrus.Error("You must provide a name for the service")
		//	return
		//}
		//
		//pbImportPath := viper.GetString("g_c_pb_import_path")
		//if viper.GetString("g_c_transport") == "grpc" {
		//	if pbImportPath == "" {
		//		logrus.Error("You must provide pb import path by --pb_import_path or -i, because transport is grpc")
		//		return
		//	}
		//}
		//g := generator.NewGenerateClient(
		//	args[0],
		//	viper.GetString("g_c_transport"),
		//	pbImportPath,
		//)
		//if err := g.Create(); err != nil {
		//	logrus.Error(err)
		//}
	},
}

//nolint:errcheck
func init() {
	genClientCmd.Flags().StringP("transport", "t", "http", "The transport you want your client to be initiated")
	genClientCmd.Flags().StringP("pb_import_path", "i", "", "Specify path to import pb")
	viper.BindPFlag("g_c_transport", genClientCmd.Flags().Lookup("transport"))
	viper.BindPFlag("g_c_pb_import_path", genClientCmd.Flags().Lookup("pb_import_path"))
}
