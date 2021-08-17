package cmd

import (
	"github.com/dave/jennifer/jen"
	"github.com/hdget/hdkit/g"
	"github.com/hdget/hdkit/generator"
)

type CmdRunFile struct {
	*generator.BaseGenerator
	Meta      *generator.Meta
	AppName   string
	GlobalDir string
}

const (
	RunFilename = "run.go"
	VarRunCmd   = "runCmd"
)

func NewCmdRunFile(meta *generator.Meta) (generator.Generator, error) {
	baseGenerator, err := generator.NewBaseGenerator(meta.Dirs[g.Cmd], RunFilename, false)
	if err != nil {
		return nil, err
	}

	return &CmdRunFile{
		BaseGenerator: baseGenerator,
		Meta:          meta,
		AppName:       meta.RootDir,
		GlobalDir:     meta.Dirs[g.Global],
	}, nil
}

func (f CmdRunFile) GetGenCodeFuncs() []func() {
	return []func(){
		f.genVar,
		f.genInitFunc,
	}
}

// var(
//  env        string
//  configFile string
// )
//var rootCmd = &cobra.Command{
//	Use:   APP,
//	Short: "bd server",
//	Long:  `bd server serves for all kinds of API`,
//}
func (f CmdRunFile) genVar() {
	found, _ := f.FindVar(VarRunCmd)
	if found == nil {
		f.Builder.Raw().Var().Id(VarRunCmd).Op("=").Id("&").Qual(g.ImportPaths[g.Cobra], "Command").Values(
			jen.Dict{
				jen.Id("Use"):   jen.Lit("run"),
				jen.Id("Short"): jen.Lit("run short description"),
				jen.Id("Long"):  jen.Lit("run long description"),
			},
		).Line()
	}
}

//func init() {
//	cobra.OnInitialize(loadConfig)
//
//	rootCmd.PersistentFlags().StringP("env", "e", "", "running environment, e,g: [prod, sim, pre, test, dev, local]")
//	rootCmd.PersistentFlags().StringP("config", "c", "", "config file, default: config.toml")
//	rootCmd.AddCommand(runServerCmd)
//}
func (f CmdRunFile) genInitFunc() {
	found, _ := f.FindMethod("init")
	if found == nil {
		body := []jen.Code{
			jen.Id(VarRunCmd).Dot("AddCommand").Call(jen.Id(VarRunGrpcServerCmd)),
			jen.Id(VarRunCmd).Dot("AddCommand").Call(jen.Id(VarRunHttpServerCmd)),
			jen.Id(VarRunCmd).Dot("AddCommand").Call(jen.Id(VarRunClientCmd)),
		}

		f.Builder.AppendFunction(
			"init",
			nil,
			nil,
			nil,
			"",
			body...,
		)
		f.Builder.NewLine()
	}
}
