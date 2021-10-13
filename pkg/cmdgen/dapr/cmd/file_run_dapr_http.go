package cmd

import (
	"github.com/dave/jennifer/jen"
	"github.com/hdget/hdkit/g"
	"github.com/hdget/hdkit/generator"
)

type CmdRunDaprHttpServerFile struct {
	*generator.BaseGenerator
	Meta      *generator.Meta
	AppName   string
	SvcDir    string
	HttpDir   string
	PbDir     string
	GlobalDir string
}

const (
	CmdRunGrpcHttpServerFilename = "run_dapr_http.go"
	VarRunGrpcHttpServerCmd      = "runDaprHttpServerCmd"
	MethodRunGrpcHttpServer      = "runDaprHttpServer"
)

func NewCmdRunHttpServerFile(meta *generator.Meta) (generator.Generator, error) {
	baseGenerator, err := generator.NewBaseGenerator(g.GetDir(meta.RootDir, g.Cmd), CmdRunGrpcHttpServerFilename, false)
	if err != nil {
		return nil, err
	}

	return &CmdRunDaprHttpServerFile{
		BaseGenerator: baseGenerator,
		Meta:          meta,
		AppName:       meta.RootDir,
		SvcDir:        g.GetDir(meta.RootDir, g.Service),
		HttpDir:       g.GetDir(meta.RootDir, g.Http),
		PbDir:         g.GetDir(meta.RootDir, g.Pb),
		GlobalDir:     g.GetDir(meta.RootDir, g.Global),
	}, nil
}

func (f CmdRunDaprHttpServerFile) GetGenCodeFuncs() []func() {
	return []func(){
		f.genImports,
		f.genVar,
		f.genRunServerFunc,
	}
}

func (f *CmdRunDaprHttpServerFile) genImports() {
	f.JenFile.ImportName(f.GlobalDir, "g")
	f.JenFile.ImportName(f.SvcDir, "service")
	f.JenFile.ImportName(g.ImportPaths[g.Errors], "errors")
	f.JenFile.ImportAlias(g.ImportPaths[g.DaprHttp], "daprd")
	f.JenFile.ImportName(g.ImportPaths[g.HdSdk], "hdsdk")
	f.JenFile.ImportName(g.ImportPaths[g.HdUtils], "utils")
	f.JenFile.ImportName(g.ImportPaths[g.Cobra], "cobra")
}

//var runHttpServerCmd = &cobra.Command{
//	Use:   "run",
//	Short: "Run http server",
//	Long:  "Run http server",
//	Run: func(cmd *cobra.Command, args []string) {
//		runServer()
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
func (f CmdRunDaprHttpServerFile) genVar() {
	found, _ := f.FindVar(VarRunGrpcHttpServerCmd)
	if found == nil {
		f.Builder.Raw().Var().Id(VarRunGrpcHttpServerCmd).Op("=").Id("&").Qual(g.ImportPaths[g.Cobra], "Command").Values(
			jen.Dict{
				jen.Id("Use"):   jen.Lit("dapr_http"),
				jen.Id("Short"): jen.Lit("run dapr http server"),
				jen.Id("Long"):  jen.Lit("run dapr http server"),
				jen.Id("Run"): jen.Func().Params(
					jen.Id("cmd").Op("*").Qual(g.ImportPaths[g.Cobra], "Command"),
					jen.Id("args").Index().String(),
				).Block(
					jen.Id(MethodRunGrpcHttpServer).Call(),
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
//  server := daprd.NewService(cliHttpAddress)
//  if server == nil {
//    hdsdk.Logger.Fatal("new http service", "error", err)
//  }
//
//  svc := service.NewUserCenter()
//
//  for url, handler := service.GetInvocationHandlers() {
//    if err := s.AddServiceInvocationHandler(url, handler); err != nil {
//      hdsdk.Logger.Fatal("adding invocation handler", "error", err)
//    }
//  }
//
//  for url, handler := service.GetBindingHandlers() {
//    if err := s.AddBindingInvocationHandler(url, handler); err != nil {
//      hdsdk.Logger.Fatal("adding binding handler", "error", err)
//    }
//  }
//
//  for sub, handler := service.GetEventHandlers() {
//    if err := s.AddTopicEventHandler(sub, handler); err != nil {
//      hdsdk.Logger.Fatal("adding event handler", "error", err)
//    }
//  }
//
// if err := s.Start(); err != nil && err != http.ErrServerClosed {
//     hdsdk.Logger.Fatal("start http service", "error", err)
// }
//
//  hdsdk.Logger.Debug("start http service", "address", address)
//}
func (f CmdRunDaprHttpServerFile) genRunServerFunc() {
	found, _ := f.FindMethod(MethodRunGrpcHttpServer)
	if found == nil {
		body := []jen.Code{
			jen.Id("server").Op(":=").Qual(g.ImportPaths[g.DaprHttp], "NewService").Call(jen.Id(VarAddress)),
			jen.If(jen.Id("server").Op("==").Nil()).Block(
				jen.Qual(g.ImportPaths[g.HdSdk], "Logger").Dot("Fatal").Call(jen.Lit("new http server"), jen.Lit("error"), jen.Lit("error new http server")),
			),
			jen.Line(),
			jen.Id("svc").Op(":=").Qual(f.SvcDir, "New"+f.Meta.RawSvcName).Call(),
			jen.Line(),
			jen.For(
				jen.List(jen.Id("url"), jen.Id("handler")).Op(":=").Range().Qual("svc", "GetInvocationHandlers").Call(),
			).Block(
				jen.If(
					jen.Err().Op(":=").Qual("server", "AddServiceInvocationHandler").Call(jen.Id("url"), jen.Id("handler")),
					jen.Err().Op("!=").Nil(),
				).Block(
					jen.Qual(g.ImportPaths[g.HdSdk], "Logger").Dot("Fatal").Call(jen.Lit("adding invocation handler"), jen.Lit("error"), jen.Id("err")),
				),
			),
			jen.Line(),
			jen.For(
				jen.List(jen.Id("url"), jen.Id("handler")).Op(":=").Range().Qual("svc", "GetBindingHandlers").Call(),
			).Block(
				jen.If(
					jen.Err().Op(":=").Qual("server", "AddBindingInvocationHandler").Call(jen.Id("url"), jen.Id("handler")),
					jen.Err().Op("!=").Nil(),
				).Block(
					jen.Qual(g.ImportPaths[g.HdSdk], "Logger").Dot("Fatal").Call(jen.Lit("adding binding handler"), jen.Lit("error"), jen.Id("err")),
				),
			),
			jen.Line(),
			jen.For(
				jen.List(jen.Id("_"), jen.Id("event")).Op(":=").Range().Qual("svc", "GetEvents").Call(),
			).Block(
				jen.If(
					jen.Err().Op(":=").Qual("server", "AddTopicEventHandler").Call(jen.Qual("event", "Sub"), jen.Qual("event", "Handler")),
					jen.Err().Op("!=").Nil(),
				).Block(
					jen.Qual(g.ImportPaths[g.HdSdk], "Logger").Dot("Fatal").Call(jen.Lit("adding event handler"), jen.Lit("error"), jen.Id("err")),
				),
			),
			jen.Line(),
			jen.Qual(g.ImportPaths[g.HdSdk], "Logger").Dot("Debug").Call(jen.Lit("start http service"), jen.Lit("address"), jen.Id(VarAddress)),
			jen.If(
				jen.Err().Op(":=").Qual("server", "Start").Call(),
				jen.Err().Op("!=").Nil().Op("&&").Err().Op("!=").Qual("net/http", "ErrServerClosed"),
			).Block(
				jen.Qual(g.ImportPaths[g.HdSdk], "Logger").Dot("Fatal").Call(jen.Lit("start http service"), jen.Lit("error"), jen.Id("err")),
			),
			jen.Line(),
		}

		f.Builder.AppendFunction(
			MethodRunGrpcHttpServer,
			nil,
			nil,
			nil,
			"",
			body...,
		)
	}
}
