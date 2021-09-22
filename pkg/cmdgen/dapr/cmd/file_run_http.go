package cmd

import (
	"github.com/dave/jennifer/jen"
	"github.com/hdget/hdkit/g"
	"github.com/hdget/hdkit/generator"
)

type CmdRunNormalHttpServerFile struct {
	*generator.BaseGenerator
	Meta      *generator.Meta
	AppName   string
	SvcDir    string
	GlobalDir string
}

const (
	CmdRunNormalHttpServerFilename = "run_http.go"
	VarRunNormalHttpServerCmd      = "runHttpServerCmd"
	MethodRunNormalHttpServer      = "runHttpServer"
)

func NewCmdRunNormalHttpServerFile(meta *generator.Meta) (generator.Generator, error) {
	baseGenerator, err := generator.NewBaseGenerator(g.GetDir(meta.RootDir, g.Cmd), CmdRunNormalHttpServerFilename, false)
	if err != nil {
		return nil, err
	}

	return &CmdRunNormalHttpServerFile{
		BaseGenerator: baseGenerator,
		Meta:          meta,
		AppName:       meta.RootDir,
		SvcDir:        g.GetDir(meta.RootDir, g.Service),
		GlobalDir:     g.GetDir(meta.RootDir, g.Global),
	}, nil
}

func (f CmdRunNormalHttpServerFile) GetGenCodeFuncs() []func() {
	return []func(){
		f.genImports,
		f.genVar,
		f.genRunServerFunc,
	}
}

func (f *CmdRunNormalHttpServerFile) genImports() {
	f.JenFile.ImportName(f.GlobalDir, "g")
	f.JenFile.ImportName(f.SvcDir, "service")
	f.JenFile.ImportName(g.ImportPaths[g.Errors], "errors")
	f.JenFile.ImportName(g.ImportPaths[g.HdSdk], "hdsdk")
	f.JenFile.ImportName(g.ImportPaths[g.HdWs], "ws")
	f.JenFile.ImportName(g.ImportPaths[g.HdUtils], "utils")
	f.JenFile.ImportName(g.ImportPaths[g.Cobra], "cobra")
}

//var runCmd = &cobra.Command{
//	Use:   "run",
//	Short: "Run server",
//	Long:  "Run server",
//	Run: func(cmd *cobra.Command, args []string) {
//		runServer(args[0])
//	},
//	PreRun: func(cmd *cobra.Command, args []string) {
//      err := hdsdk.Initialize(g.Config)
//      if err != nil {
//          utils.LogFatal("hdsdk initialize", "err", err)
//      }
//	},
//	PostRun: func(cmd *cobra.Command, args []string) {
//		hdsdk.Shutdown()
//	},
//}
func (f CmdRunNormalHttpServerFile) genVar() {
	found, _ := f.FindVar(VarRunNormalHttpServerCmd)
	if found == nil {
		f.Builder.Raw().Var().Id(VarRunNormalHttpServerCmd).Op("=").Id("&").Qual(g.ImportPaths[g.Cobra], "Command").Values(
			jen.Dict{
				jen.Id("Use"):   jen.Lit("http"),
				jen.Id("Short"): jen.Lit("run normal http server"),
				jen.Id("Long"):  jen.Lit("run normal http server"),
				jen.Id("Run"): jen.Func().Params(
					jen.Id("cmd").Op("*").Qual(g.ImportPaths[g.Cobra], "Command"),
					jen.Id("args").Index().String(),
				).Block(
					jen.Id(MethodRunNormalHttpServer).Call(),
				),
				jen.Id("PreRun"): jen.Func().Params(
					jen.Id("cmd").Op("*").Qual(g.ImportPaths[g.Cobra], "Command"),
					jen.Id("args").Index().String(),
				).Block(
					jen.Err().Op(":=").Qual(g.ImportPaths[g.HdSdk], "Initialize").Call(jen.Qual(f.GlobalDir, "Config")),
					jen.If(jen.Err().Op("!=").Nil()).Block(
						jen.Qual(g.ImportPaths[g.HdUtils], "LogFatal").Call(jen.Lit("sdk initialize"), jen.Lit("err"), jen.Err()),
					),
				),
				jen.Id("PostRun"): jen.Func().Params(
					jen.Id("cmd").Op("*").Qual(g.ImportPaths[g.Cobra], "Command"),
					jen.Id("args").Index().String(),
				).Block(
				// jen.Qual(HdSdkImportPath, "Shutdown").Call(),
				),
			},
		).Line()
	}
}

//func runHttpServer() {
// ws.SetReleaseMode()
// srv := ws.NewHttpServer(hdsdk.Logger, cliAddress)
// srv.SetupRoutes(oidc.NewOidcProvider().GetRoutes())
// srv.Run()
//
//}
func (f CmdRunNormalHttpServerFile) genRunServerFunc() {
	found, _ := f.FindMethod(MethodRunNormalHttpServer)
	if found == nil {
		body := []jen.Code{
			jen.Qual(g.ImportPaths[g.HdWs], "SetReleaseMode").Call(),
			jen.Id("svc").Op(":=").Qual(g.ImportPaths[g.HdWs], "NewHttpServer").Call(jen.Qual(g.ImportPaths[g.HdSdk], "Logger"), jen.Id("cliAddress")),
			jen.Id("svc").Dot("SetupRoutes").Call(jen.Nil()),
			jen.Id("svc").Dot("Run").Call(),
		}

		f.Builder.AppendFunction(
			MethodRunNormalHttpServer,
			nil,
			nil,
			nil,
			"",
			body...,
		)
	}
}
