package cmd

import (
	"github.com/dave/jennifer/jen"
	"github.com/hdget/hdkit/g"
	"github.com/hdget/hdkit/generator"
)

type CmdRunDaprGrpcServerFile struct {
	*generator.BaseGenerator
	Meta      *generator.Meta
	AppName   string
	SvcDir    string
	GlobalDir string
}

const (
	CmdRunDaprGrpcServerFilename = "run_dapr_grpc.go"
	VarRunDaprGrpcServerCmd      = "runDaprGrpcServerCmd"
	MethodRunDaprGprcServer      = "runDaprGrpcServer"
)

func NewCmdRunGrpcServerFile(meta *generator.Meta) (generator.Generator, error) {
	baseGenerator, err := generator.NewBaseGenerator(g.GetDir(meta.RootDir, g.Cmd), CmdRunDaprGrpcServerFilename, false)
	if err != nil {
		return nil, err
	}

	return &CmdRunDaprGrpcServerFile{
		BaseGenerator: baseGenerator,
		Meta:          meta,
		AppName:       meta.RootDir,
		SvcDir:        g.GetDir(meta.RootDir, g.Service),
		GlobalDir:     g.GetDir(meta.RootDir, g.Global),
	}, nil
}

func (f CmdRunDaprGrpcServerFile) GetGenCodeFuncs() []func() {
	return []func(){
		f.genImports,
		f.genVar,
		f.genRunServerFunc,
	}
}

func (f *CmdRunDaprGrpcServerFile) genImports() {
	f.JenFile.ImportName(f.GlobalDir, "g")
	f.JenFile.ImportName(f.SvcDir, "service")
	f.JenFile.ImportName(g.ImportPaths[g.Errors], "errors")
	f.JenFile.ImportAlias(g.ImportPaths[g.DaprGrpc], "daprd")
	f.JenFile.ImportName(g.ImportPaths[g.HdSdk], "hdsdk")
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
func (f CmdRunDaprGrpcServerFile) genVar() {
	found, _ := f.FindVar(VarRunDaprGrpcServerCmd)
	if found == nil {
		f.Builder.Raw().Var().Id(VarRunDaprGrpcServerCmd).Op("=").Id("&").Qual(g.ImportPaths[g.Cobra], "Command").Values(
			jen.Dict{
				jen.Id("Use"):   jen.Lit("dapr_grpc"),
				jen.Id("Short"): jen.Lit("run dapr grpc server"),
				jen.Id("Long"):  jen.Lit("run dapr grpc server"),
				jen.Id("Run"): jen.Func().Params(
					jen.Id("cmd").Op("*").Qual(g.ImportPaths[g.Cobra], "Command"),
					jen.Id("args").Index().String(),
				).Block(
					jen.Id(MethodRunDaprGprcServer).Call(),
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

//func runDaprGrpcServer() {
//  server, err := daprd.NewService(address)
//  if err != nil {
//    hdsdk.Logger.Fatal("new dapr service", "error", err)
//  }
//
//  svc := service.NewUserCenter()
//
//  for method, handler := svc.GetInvocationHandlers() {
//    if err := s.AddServiceInvocationHandler(method, handler); err != nil {
//      hdsdk.Logger.Fatal("adding invocation handler", "error", err)
//    }
//  }
//
//  for name, handler := svc.GetBindingHandlers() {
//    if err := s.AddBindingInvocationHandler(name, handler); err != nil {
//      hdsdk.Logger.Fatal("adding binding handler", "error", err)
//    }
//  }
//
//  for sub, handler := svc.GetEventHandlers() {
//    if err := s.AddTopicEventHandler(sub, handler); err != nil {
//      hdsdk.Logger.Fatal("adding event handler", "error", err)
//    }
//  }
//
//  if err := svc.Start(); err != nil {
//     hdsdk.Logger.Fatal("grpc service start", "error", err)
//  }
//
//  hdsdk.Logger.Debug("grpc service start", "address", address)
//}
func (f CmdRunDaprGrpcServerFile) genRunServerFunc() {
	found, _ := f.FindMethod(MethodRunDaprGprcServer)
	if found == nil {
		body := []jen.Code{
			jen.List(jen.Id("server"), jen.Err()).Op(":=").Qual(g.ImportPaths[g.DaprGrpc], "NewService").Call(jen.Id(VarAddress)),
			jen.If(jen.Err().Op("!=").Nil()).Block(
				jen.Qual(g.ImportPaths[g.HdSdk], "Logger").Dot("Fatal").Call(jen.Lit("new dapr service"), jen.Lit("error"), jen.Id("err")),
			),
			jen.Line(),
			jen.Id("svc").Op(":=").Qual(f.SvcDir, "New"+f.Meta.RawSvcName).Call(),
			jen.Line(),
			jen.For(
				jen.List(jen.Id("method"), jen.Id("handler")).Op(":=").Range().Qual("svc", "GetInvocationHandlers").Call(),
			).Block(
				jen.If(
					jen.Err().Op(":=").Qual("server", "AddServiceInvocationHandler").Call(jen.Id("method"), jen.Id("handler")),
					jen.Err().Op("!=").Nil(),
				).Block(
					jen.Qual(g.ImportPaths[g.HdSdk], "Logger").Dot("Fatal").Call(jen.Lit("adding invocation handler"), jen.Lit("error"), jen.Id("err")),
				),
			),
			jen.Line(),
			jen.For(
				jen.List(jen.Id("name"), jen.Id("handler")).Op(":=").Range().Qual("svc", "GetBindingHandlers").Call(),
			).Block(
				jen.If(
					jen.Err().Op(":=").Qual("server", "AddBindingInvocationHandler").Call(jen.Id("name"), jen.Id("handler")),
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
			jen.Qual(g.ImportPaths[g.HdSdk], "Logger").Dot("Debug").Call(jen.Lit("start grpc service"), jen.Lit("address"), jen.Id(VarAddress)),
			jen.If(
				jen.Err().Op(":=").Qual("server", "Start").Call(),
				jen.Err().Op("!=").Nil(),
			).Block(
				jen.Qual(g.ImportPaths[g.HdSdk], "Logger").Dot("Fatal").Call(jen.Lit("start grpc service"), jen.Lit("error"), jen.Id("err")),
			),
			jen.Line(),
		}

		f.Builder.AppendFunction(
			MethodRunDaprGprcServer,
			nil,
			nil,
			nil,
			"",
			body...,
		)
	}
}
