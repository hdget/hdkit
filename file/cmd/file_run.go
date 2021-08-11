package cmd

import (
	"github.com/dave/jennifer/jen"
	"github.com/hdget/hdkit/g"
	"github.com/hdget/hdkit/generator"
)

type CmdRunServerFile struct {
	*generator.BaseGenerator
	Meta      *generator.Meta
	AppName   string
	SvcDir    string
	GrpcDir   string
	PbDir     string
	GlobalDir string
}

const (
	CmdRunServerFilename = "run.go"
	VarRunCmd            = "runCmd"
	MethodRunServer      = "runServer"
)

func NewCmdRunServerFile(meta *generator.Meta) (generator.Generator, error) {
	baseGenerator, err := generator.NewBaseGenerator(meta.Dirs[g.Cmd], CmdRunServerFilename, false)
	if err != nil {
		return nil, err
	}

	return &CmdRunServerFile{
		BaseGenerator: baseGenerator,
		Meta:          meta,
		AppName:       meta.RootDir,
		SvcDir:        meta.Dirs[g.Service],
		GrpcDir:       meta.Dirs[g.Grpc],
		PbDir:         meta.Dirs[g.Pb],
		GlobalDir:     meta.Dirs[g.Global],
	}, nil
}

func (f CmdRunServerFile) GetGenCodeFuncs() []func() {
	return []func(){
		f.genImports,
		f.genVar,
		f.genRunServerFunc,
	}
}

func (f *CmdRunServerFile) genImports() {
	f.JenFile.ImportName(f.GlobalDir, "g")
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
//          utils.Fatal("hdsdk initialize", "err", err)
//      }
//	},
//	PostRun: func(cmd *cobra.Command, args []string) {
//		hdsdk.Shutdown()
//	},
//}
func (f CmdRunServerFile) genVar() {
	found, _ := f.FindVar(VarRunCmd)
	if found == nil {
		f.Builder.Raw().Var().Id("runCmd").Op("=").Id("&").Qual(g.ImportPaths[g.Cobra], "Command").Values(
			jen.Dict{
				jen.Id("Use"):   jen.Lit("run"),
				jen.Id("Short"): jen.Lit("run server short description"),
				jen.Id("Long"):  jen.Lit("run server long description"),
				jen.Id("Run"): jen.Func().Params(
					jen.Id("cmd").Op("*").Qual(g.ImportPaths[g.Cobra], "Command"),
					jen.Id("args").Index().String(),
				).Block(
					jen.Id("runServer").Call(),
				),
				jen.Id("PreRun"): jen.Func().Params(
					jen.Id("cmd").Op("*").Qual(g.ImportPaths[g.Cobra], "Command"),
					jen.Id("args").Index().String(),
				).Block(
					jen.Err().Op(":=").Qual(g.ImportPaths[g.HdSdk], "Initialize").Call(jen.Qual(f.GlobalDir, "Config")),
					jen.If(jen.Err().Op("!=").Nil()).Block(
						jen.Qual(g.ImportPaths[g.HdUtils], "Fatal").Call(jen.Lit("sdk initialize"), jen.Lit("err"), jen.Err()),
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
//  server := ms.CreateGrpcServer()
//  if server == nil {
//      hdsdk.Logger.Fatal("create grpc server", "err", "empty server")
//  }
//
//	svc := service.NewSearchService()
//	handlers := grpc.NewHandlers(server, svc)
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
func (f CmdRunServerFile) genRunServerFunc() {
	found, _ := f.FindMethod(MethodRunServer)
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

		f.Builder.AppendFunction(MethodRunServer, nil, nil, nil, "", body...)
	}
}
