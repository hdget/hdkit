package cmd

import (
	"github.com/dave/jennifer/jen"
	"github.com/hdget/hdkit/g"
	"github.com/hdget/hdkit/generator"
)

type CmdRunHttpServerFile struct {
	*generator.BaseGenerator
	Meta      *generator.Meta
	AppName   string
	SvcDir    string
	HttpDir   string
	PbDir     string
	GlobalDir string
}

const (
	CmdRunHttpServerFilename = "run_http.go"
	VarRunHttpServerCmd      = "runHttpServerCmd"
	MethodRunHttpServer      = "runHttpServer"
)

func NewCmdRunHttpServerFile(meta *generator.Meta) (generator.Generator, error) {
	baseGenerator, err := generator.NewBaseGenerator(g.GetDir(meta.RootDir, g.Cmd), CmdRunHttpServerFilename, false)
	if err != nil {
		return nil, err
	}

	return &CmdRunHttpServerFile{
		BaseGenerator: baseGenerator,
		Meta:          meta,
		AppName:       meta.RootDir,
		SvcDir:        g.GetDir(meta.RootDir, g.Service),
		HttpDir:       g.GetDir(meta.RootDir, g.Http),
		PbDir:         g.GetDir(meta.RootDir, g.Pb),
		GlobalDir:     g.GetDir(meta.RootDir, g.Global),
	}, nil
}

func (f CmdRunHttpServerFile) GetGenCodeFuncs() []func() {
	return []func(){
		f.genImports,
		f.genVar,
		f.genRunServerFunc,
	}
}

func (f *CmdRunHttpServerFile) genImports() {
	f.JenFile.ImportName(f.GlobalDir, "g")
	f.JenFile.ImportName(f.PbDir, "pb")
	f.JenFile.ImportName(f.SvcDir, "service")
	f.JenFile.ImportAlias(f.HttpDir, "genhttp")
	f.JenFile.ImportName(g.ImportPaths[g.Errors], "errors")
	f.JenFile.ImportName(g.ImportPaths[g.StdGrpc], "grpc")
	f.JenFile.ImportName(g.ImportPaths[g.HdSdk], "hdsdk")
	f.JenFile.ImportName(g.ImportPaths[g.HdUtils], "utils")
	f.JenFile.ImportName(g.ImportPaths[g.KitEndpoint], "endpoint")
	f.JenFile.ImportName(g.ImportPaths[g.Cobra], "cobra")
	f.JenFile.ImportName(g.ImportPaths[g.HdParallel], "parallel")
}

//var runCmd = &cobra.Command{
//	Use:   "run",
//	Short: "Run server",
//	Long:  "Run server",
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
func (f CmdRunHttpServerFile) genVar() {
	found, _ := f.FindVar(VarRunHttpServerCmd)
	if found == nil {
		f.Builder.Raw().Var().Id(VarRunHttpServerCmd).Op("=").Id("&").Qual(g.ImportPaths[g.Cobra], "Command").Values(
			jen.Dict{
				jen.Id("Use"):   jen.Lit("http"),
				jen.Id("Short"): jen.Lit("run server short description"),
				jen.Id("Long"):  jen.Lit("run server long description"),
				jen.Id("Run"): jen.Func().Params(
					jen.Id("cmd").Op("*").Qual(g.ImportPaths[g.Cobra], "Command"),
					jen.Id("args").Index().String(),
				).Block(
					jen.Id(MethodRunHttpServer).Call(),
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

//func runServer() {
//  ms := hdsdk.MicroService.My()
//  if ms == nil {
//    hdsdk.Logger.Fatal("get microservice instance", "err", "empty microservice instance")
//  }
//
//  manager := ms.NewGrpcServerManager()
//  if manager == nil {
//      hdsdk.Logger.Fatal("create grpc server", "err", "empty server")
//  }
//
//	svc := service.NewSearchService()
//	handlers := grpc.NewHandlers(manager, svc)
//	pb.RegisterSearchServiceServer(server.GetServer(), handlers)
//
//  svc := service.NewPartner()
//  handlers := genhttp.NewHandlers(manager, svc)
//  err := manager.RunServer(handlers)
//  if err != nil {
//      hdsdk.Logger.Fatal("http server exited", "err", err)
//  }
//}
func (f CmdRunHttpServerFile) genRunServerFunc() {
	found, _ := f.FindMethod(MethodRunHttpServer)
	if found == nil {
		body := []jen.Code{
			jen.Id("ms").Op(":=").Qual(g.ImportPaths[g.HdSdk], "MicroService").Dot("My").Call(),
			jen.If(jen.Id("ms").Op("==").Nil()).Block(
				jen.Qual(g.ImportPaths[g.HdSdk], "Logger").Dot("Fatal").Call(jen.Lit("get microservice instance"), jen.Lit("err"), jen.Lit("empty microservice instance")),
			),
			jen.Line(),
			jen.Id("manager").Op(":=").Qual("ms", "NewHttpServerManager").Call(),
			jen.If(jen.Id("manager").Op("==").Nil()).Block(
				jen.Qual(g.ImportPaths[g.HdSdk], "Logger").Dot("Fatal").Call(jen.Lit("new http server manager"), jen.Lit("err"), jen.Lit("create failed")),
			),
			jen.Line(),
			jen.Id("svc").Op(":=").Qual(f.SvcDir, "New"+f.Meta.RawSvcName).Call(),
			jen.Id("handlers").Op(":=").Qual(f.HttpDir, "NewHandlers").Call(jen.Id("manager"), jen.Id("svc")),

			jen.Line(),
			jen.Var().Id("group").Qual(g.ImportPaths[g.HdParallel], "Group"),
			jen.Qual("group", "Add").Call(
				jen.Func().Params().Id("error").Block(
					jen.Return(jen.Id("manager").Dot("RunServer").Call(jen.Id("handlers"))),
				),
				jen.Func().Params(jen.Id("err error")).Block(
					jen.Id("manager").Dot("Close").Call(),
				),
			),
			jen.Err().Op(":=").Id("group").Dot("Run").Call(),
			jen.If(jen.Err().Op("!=").Nil()).Block(
				jen.Qual(g.ImportPaths[g.HdSdk], "Logger").Dot("Fatal").Call(jen.Lit("http server exited"), jen.Lit("err"), jen.Id("err")),
			),
		}

		f.Builder.AppendFunction(MethodRunHttpServer, nil, nil, nil, "", body...)
	}
}
