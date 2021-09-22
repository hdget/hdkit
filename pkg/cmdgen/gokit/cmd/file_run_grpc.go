package cmd

import (
	"github.com/dave/jennifer/jen"
	"github.com/hdget/hdkit/g"
	"github.com/hdget/hdkit/generator"
)

type CmdRunGrpcServerFile struct {
	*generator.BaseGenerator
	Meta      *generator.Meta
	AppName   string
	SvcDir    string
	GrpcDir   string
	PbDir     string
	GlobalDir string
}

const (
	CmdRunGrpcServerFilename = "run_grpc.go"
	VarRunGrpcServerCmd      = "runGrpcServerCmd"
	MethodRunGprcServer      = "runGrpcServer"
)

func NewCmdRunGrpcServerFile(meta *generator.Meta) (generator.Generator, error) {
	baseGenerator, err := generator.NewBaseGenerator(g.GetDir(meta.RootDir, g.Cmd), CmdRunGrpcServerFilename, false)
	if err != nil {
		return nil, err
	}

	return &CmdRunGrpcServerFile{
		BaseGenerator: baseGenerator,
		Meta:          meta,
		AppName:       meta.RootDir,
		SvcDir:        g.GetDir(meta.RootDir, g.Service),
		GrpcDir:       g.GetDir(meta.RootDir, g.Grpc),
		PbDir:         g.GetDir(meta.RootDir, g.Pb),
		GlobalDir:     g.GetDir(meta.RootDir, g.Global),
	}, nil
}

func (f CmdRunGrpcServerFile) GetGenCodeFuncs() []func() {
	return []func(){
		f.genImports,
		f.genVar,
		f.genRunServerFunc,
	}
}

func (f *CmdRunGrpcServerFile) genImports() {
	f.JenFile.ImportName(f.GlobalDir, "g")
	f.JenFile.ImportName(f.PbDir, "pb")
	f.JenFile.ImportName(f.SvcDir, "service")
	f.JenFile.ImportAlias(f.GrpcDir, "gengrpc")
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
func (f CmdRunGrpcServerFile) genVar() {
	found, _ := f.FindVar(VarRunGrpcServerCmd)
	if found == nil {
		f.Builder.Raw().Var().Id(VarRunGrpcServerCmd).Op("=").Id("&").Qual(g.ImportPaths[g.Cobra], "Command").Values(
			jen.Dict{
				jen.Id("Use"):   jen.Lit("grpc"),
				jen.Id("Short"): jen.Lit("run server short description"),
				jen.Id("Long"):  jen.Lit("run server long description"),
				jen.Id("Run"): jen.Func().Params(
					jen.Id("cmd").Op("*").Qual(g.ImportPaths[g.Cobra], "Command"),
					jen.Id("args").Index().String(),
				).Block(
					jen.Id(MethodRunGprcServer).Call(),
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
//	var group parallel.Group
//  group.Add(
//    func() error {
//	    return server.Run()
//    },
//    func(err error) {
//	    server.Close()
//    },
// )
// group.Run()
//}
func (f CmdRunGrpcServerFile) genRunServerFunc() {
	found, _ := f.FindMethod(MethodRunGprcServer)
	if found == nil {
		body := []jen.Code{
			jen.Id("ms").Op(":=").Qual(g.ImportPaths[g.HdSdk], "MicroService").Dot("My").Call(),
			jen.If(jen.Id("ms").Op("==").Nil()).Block(
				jen.Qual(g.ImportPaths[g.HdSdk], "Logger").Dot("Fatal").Call(jen.Lit("get microservice instance"), jen.Lit("err"), jen.Lit("empty microservice instance")),
			),
			jen.Line(),
			jen.Id("manager").Op(":=").Qual("ms", "NewGrpcServerManager").Call(),
			jen.If(jen.Id("manager").Op("==").Nil()).Block(
				jen.Qual(g.ImportPaths[g.HdSdk], "Logger").Dot("Fatal").Call(jen.Lit("new grpc server manager"), jen.Lit("err"), jen.Lit("create failed")),
			),
			jen.Line(),
			jen.Id("svc").Op(":=").Qual(f.SvcDir, "New"+f.Meta.RawSvcName).Call(),
			jen.Id("handlers").Op(":=").Qual(f.GrpcDir, "NewHandlers").Call(jen.Id("manager"), jen.Id("svc")),
			jen.Qual(f.PbDir, "Register"+f.Meta.SvcServerInterfaceName).Call(jen.Id("manager").Dot("GetServer").Call(), jen.Id("handlers")),
			jen.Line(),
			jen.Var().Id("group").Qual(g.ImportPaths[g.HdParallel], "Group"),
			jen.Qual("group", "Add").Call(
				jen.Func().Params().Id("error").Block(
					jen.Return(jen.Id("manager").Dot("RunServer").Call()),
				),
				jen.Func().Params(jen.Id("err error")).Block(
					jen.Id("manager").Dot("Close").Call(),
				),
			),
			jen.Err().Op(":=").Id("group").Dot("Run").Call(),
			jen.If(jen.Err().Op("!=").Nil()).Block(
				jen.Qual(g.ImportPaths[g.HdSdk], "Logger").Dot("Fatal").Call(jen.Lit("grpc server exited"), jen.Lit("err"), jen.Id("err")),
			),
		}

		f.Builder.AppendFunction(MethodRunGprcServer, nil, nil, nil, "", body...)
	}
}
