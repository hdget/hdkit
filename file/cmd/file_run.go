package cmd

import (
	"github.com/dave/jennifer/jen"
	"github.com/hdget/hdkit/g"
	"github.com/hdget/hdkit/generator"
)

type CmdRunServerFile struct {
	*generator.BaseGenerator
	Meta    *generator.Meta
	AppName string
	SvcDir  string
	GrpcDir string
	PbDir   string
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
	}, nil
}

func (f *CmdRunServerFile) PreGenerate() error {
	err := f.BaseGenerator.PreGenerate()
	if err != nil {
		return err
	}

	// package a should be aliased to "b"
	f.JenFile.ImportName(CobraImportPath, "cobra")
	return nil
}

func (f CmdRunServerFile) GetGenCodeFuncs() []func() {
	return []func(){
		f.genVar,
		f.genRunServerFunc,
	}
}

//var runCmd = &cobra.Command{
//	Use:   "run",
//	Short: "Run server",
//	Long:  "Run server",
//	Run: func(cmd *cobra.Command, args []string) {
//		runServer()
//	},
//	PreRun: func(cmd *cobra.Command, args []string) {
//      err := sdk.Initialize(g.Config)
//      if err != nil {
//          utils.Fatal("sdk initialize", "err", err)
//      }
//	},
//	PostRun: func(cmd *cobra.Command, args []string) {
//		sdk.Shutdown()
//	},
//}
func (f CmdRunServerFile) genVar() {
	found, _ := f.FindVar(VarRunCmd)
	if found == nil {
		f.Builder.Raw().Var().Id("runCmd").Op("=").Id("&").Qual(CobraImportPath, "Command").Values(
			jen.Dict{
				jen.Id("Use"):   jen.Lit("run"),
				jen.Id("Short"): jen.Lit("run server short description"),
				jen.Id("Long"):  jen.Lit("run server long description"),
				jen.Id("Run"): jen.Func().Params(
					jen.Id("cmd").Op("*").Qual("github.com/spf13/cobra", "Command"),
					jen.Id("args").Index().String(),
				).Block(
					jen.Id("runServer").Call(),
				),
				jen.Id("PreRun"): jen.Func().Params(
					jen.Id("cmd").Op("*").Qual("github.com/spf13/cobra", "Command"),
					jen.Id("args").Index().String(),
				).Block(
					jen.Err().Op(":=").Qual(SdkImportPath, "Initialize").Call(jen.Id("g").Dot("Config")),
					jen.If(jen.Err().Op("!=").Nil()).Block(
						jen.Qual(UtilsImportPath, "Fatal").Call(jen.Lit("sdk initialize"), jen.Lit("err"), jen.Err()),
					),
				),
				jen.Id("PostRun"): jen.Func().Params(
					jen.Id("cmd").Op("*").Qual("github.com/spf13/cobra", "Command"),
					jen.Id("args").Index().String(),
				).Block(
				// jen.Qual(SdkImportPath, "Shutdown").Call(),
				),
			},
		).Line()
	}
}

//func runServer() {
//	server := sdk.MicroService.My().CreateGrpcServer()
//	svc := service.NewSearchService()
//
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
			jen.Id("server").Op(":=").Qual(SdkImportPath, "MicroService").Dot("My").Call().Dot("CreateGrpcServer").Call(),
			jen.Id("svc").Op(":=").Qual(f.SvcDir, "New"+f.Meta.RawSvcName).Call(),
			jen.Line(),
			jen.Id("handlers").Op(":=").Qual(f.GrpcDir, "NewHandlers").Call(jen.Id("server"), jen.Id("svc")),
			jen.Qual(f.PbDir, "Register"+f.Meta.SvcServerInterfaceName).Call(jen.Id("server").Dot("GetServer").Call(), jen.Id("handlers")),
			jen.Line(),
			jen.Var().Id("group").Qual(ParallelImportPath, "Group"),
			jen.Qual("group", "Add").Call(
				jen.Func().Params().Id("error").Block(
					jen.Return(jen.Id("server").Dot("Run").Call()),
				),
				jen.Func().Params(jen.Id("err error")).Block(
					jen.Id("server").Dot("Close").Call(),
				),
			),
			jen.Id("group").Dot("Run").Call(),
		}

		f.Builder.AppendFunction(MethodRunServer, nil, nil, nil, "", body...)
	}
}
