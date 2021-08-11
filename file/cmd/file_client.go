package cmd

import (
	"github.com/dave/jennifer/jen"
	"github.com/hdget/hdkit/g"
	"github.com/hdget/hdkit/generator"
)

type CmdRunClientFile struct {
	*generator.BaseGenerator
	Meta      *generator.Meta
	AppName   string
	SvcDir    string
	GrpcDir   string
	PbDir     string
	GlobalDir string
}

const (
	CmdRunClientFilename = "client.go"
	VarClientCmd         = "clientCmd"
	MethodRunClient      = "runClient"
)

func NewCmdRunClientFile(meta *generator.Meta) (generator.Generator, error) {
	baseGenerator, err := generator.NewBaseGenerator(meta.Dirs[g.Cmd], CmdRunClientFilename, false)
	if err != nil {
		return nil, err
	}

	return &CmdRunClientFile{
		BaseGenerator: baseGenerator,
		Meta:          meta,
		AppName:       meta.RootDir,
		GrpcDir:       meta.Dirs[g.Grpc],
		PbDir:         meta.Dirs[g.Pb],
		GlobalDir:     meta.Dirs[g.Global],
	}, nil
}

func (f CmdRunClientFile) GetGenCodeFuncs() []func() {
	return []func(){
		f.genImports,
		f.genVar,
		f.genRunClientFunc,
	}
}

func (f *CmdRunClientFile) genImports() {
	f.JenFile.ImportName(f.GlobalDir, "g")
	f.JenFile.ImportAlias(f.GrpcDir, "gengrpc")
}

//var runCmd = &cobra.Command{
//	Use:   "client",
//	Short: "run client short desc",
//	Long:  "run client long desc",
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
//		// hdsdk.Shutdown()
//	},
//}
func (f CmdRunClientFile) genVar() {
	found, _ := f.FindVar(VarRunCmd)
	if found == nil {
		f.Builder.Raw().Var().Id(VarClientCmd).Op("=").Id("&").Qual(g.ImportPaths[g.Cobra], "Command").Values(
			jen.Dict{
				jen.Id("Use"):   jen.Lit("client"),
				jen.Id("Short"): jen.Lit("run client short description"),
				jen.Id("Long"):  jen.Lit("run client long description"),
				jen.Id("Run"): jen.Func().Params(
					jen.Id("cmd").Op("*").Qual(g.ImportPaths[g.Cobra], "Command"),
					jen.Id("args").Index().String(),
				).Block(
					jen.Err().Op(":=").Id(MethodRunClient).Call(),
					jen.If(jen.Err().Op("!=").Nil()).Block(
						jen.Qual(g.ImportPaths[g.HdUtils], "Fatal").Call(jen.Lit("run client"), jen.Lit("err"), jen.Err()),
					),
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

//func runClient() error {
//  ms := hdsdk.MicroService.My()
//  if ms == nil {
//    return errors.New("microservice not found")
//  }
//
//  manager := ms.NewGrpcServerManager()
//  if manager == nil {
//    return errors.New("empty grpc server manager")
//  }
//
//  conn, err := grpc.DialContext(context.Background(), "0.0.0.0:12345", grpc.WithInsecure())
//  if err != nil {
//    return err
//  }
//  defer conn.Close()
//
//  client, err := gengrpc.NewClient(conn)
//  if err != nil {
//     return err
//  }
//
//  result, err := client.WhoIs(context.Background(), &pb.EmptyRequest{})
//  if err != nil {
//    return err
//  }
//  hdsdk.Logger.Debug(result)
func (f CmdRunClientFile) genRunClientFunc() {
	found, _ := f.FindMethod(MethodRunServer)
	if found == nil {
		body := []jen.Code{
			jen.Id("ms").Op(":=").Qual(g.ImportPaths[g.HdSdk], "MicroService").Dot("My").Call(),
			jen.If(jen.Id("ms").Op("==").Nil()).Block(
				jen.Return(jen.Qual(g.ImportPaths[g.Errors], "New").Call(jen.Lit("microservice not found"))),
			),
			jen.Line(),
			jen.Id("manager").Op(":=").Qual("ms", "NewGrpcClientManager").Call(),
			jen.If(jen.Id("manager").Op("==").Nil()).Block(
				jen.Return(jen.Qual(g.ImportPaths[g.Errors], "New").Call(jen.Lit("empty grpc server manager"))),
			),
			jen.Line(),
			jen.List(jen.Id("conn"), jen.Err()).Op(":=").Qual(g.ImportPaths[g.StdGrpc], "DialContext").Call(
				jen.Qual("context", "Background").Call(),
				jen.Lit("0.0.0.0:12345"),
				jen.Qual(g.ImportPaths[g.StdGrpc], "WithInsecure").Call(),
			),
			jen.If(jen.Err().Op("!=").Nil()).Block(
				jen.Return(jen.Err()),
			),
			jen.Defer().Id("conn").Dot("Close").Call(),
			jen.Line(),
			jen.List(jen.Id("client"), jen.Err()).Op(":=").Qual("gengrpc", "NewClient").Call(
				jen.Id("conn"),
			),
			jen.If(jen.Err().Op("!=").Nil()).Block(
				jen.Return(jen.Err()),
			),
			jen.Line(),
			jen.List(jen.Id("response"), jen.Err()).Op(":=").Id("client").Dot("WhoIs").Call(
				jen.Qual("context", "Background").Call(),
				jen.Nil(),
			),
			jen.If(jen.Err().Op("!=").Nil()).Block(
				jen.Return(jen.Err()),
			),
			jen.Qual("fmt", "Println").Call(jen.Id("response")),
			jen.Return(jen.Nil()),
		}

		f.Builder.AppendFunction(MethodRunClient, nil, nil, []jen.Code{jen.Error()}, "", body...)
	}
}
