package grpc

import (
	"github.com/dave/jennifer/jen"
	"github.com/hdget/hdkit/g"
	"github.com/hdget/hdkit/generator"
	"github.com/hdget/hdkit/parser"
	"github.com/hdget/hdkit/utils"
)

type GrpcHandlersFile struct {
	*generator.BaseGenerator
	Meta       *generator.Meta
	StructName string
	PbDir      string
}

const (
	HandlersFileName    = "grpc_handlers.go"
	HandlersStructName  = "GrpcHandlers"
	HandlerStructSuffix = "Handler"
)

var (
	HandlersStructComments = []string{
		"Handlers collects all of the handlers that compose a service.",
	}
)

func NewGrpcHandlersFile(meta *generator.Meta) (generator.Generator, error) {
	baseGenerator, err := generator.NewBaseGenerator(g.GetDir(meta.RootDir, g.Grpc), HandlersFileName, true)
	if err != nil {
		return nil, err
	}

	return &GrpcHandlersFile{
		BaseGenerator: baseGenerator,
		Meta:          meta,
		StructName:    HandlersStructName,
		PbDir:         g.GetDir(meta.RootDir, g.Pb),
	}, nil
}

func (f *GrpcHandlersFile) GetGenCodeFuncs() []func() {
	return []func(){
		f.genImports,
		f.genHandlersStructure,
		f.genNewHandlersFunction,
		f.genGrpcServeMethods,
	}
}

func (f *GrpcHandlersFile) genImports() {
	f.JenFile.ImportName(f.PbDir, "pb")
	f.JenFile.ImportName(g.ImportPaths[g.Errors], "errors")
	f.JenFile.ImportName(g.ImportPaths[g.HdSdkTypes], "types")
	f.JenFile.ImportAlias(g.ImportPaths[g.KitGrpc], "kitgrpc")
}

// genHandlersStructure collects all of the handlers that compose a service.
// type Handlers struct {
//		SearchHandler kitgrpc.Handler
//		HelloHandler  kitgrpc.Handler
//		WorldHandler  kitgrpc.Handler
// }
func (f *GrpcHandlersFile) genHandlersStructure() {
	f.JenFile.PackageComment("Package grpc THIS FILE IS AUTO GENERATED PLEASE DO NOT EDIT!!")

	fields := make([]jen.Code, 0)
	for _, v := range f.Meta.SvcServerInterface.Methods {
		fields = append(fields, jen.Id(v.Name+HandlerStructSuffix).Qual("github.com/go-kit/kit/transport/grpc", "Handler"))
	}

	f.Builder.AppendMultilineComment(HandlersStructComments)
	f.Builder.NewLine()
	f.Builder.AppendStruct(
		HandlersStructName,
		fields...,
	)
}

// genNewHandlersFunction returns new handlers function
//
// func NewHandlers(manager GrpcServerManager, svc Service) *Handlers {
//    return &Handlers{
//   			SearchHandler: sever.CreateHandler(svc, NewServerHandler()),
//				HelloHandler: sever.CreateHandler(svc, NewHelloHandler()),
//	  }
// }
func (f *GrpcHandlersFile) genNewHandlersFunction() {
	handlers := jen.Dict{}
	for _, m := range f.Meta.SvcServerInterface.Methods {
		handlerName := m.Name + HandlerStructSuffix
		aspectName := m.Name + AspectSuffix
		handlers[jen.Id(handlerName)] = jen.Qual("manager", "CreateHandler").Call(
			jen.Id("svc"), jen.Id("&"+aspectName+"{}"))
	}

	body := jen.Return(jen.Op("&").Id(HandlersStructName).Values(handlers))

	f.Builder.Raw().
		Func().
		Id("NewHandlers").
		Params(
			jen.Id("manager").Qual(g.ImportPaths[g.HdSdkTypes], "GrpcServerManager"),
			jen.Id("svc").Qual(f.PbDir, f.Meta.SvcServerInterface.Name),
		).
		Op("*").Id(HandlersStructName).
		Block(body)
	f.Builder.NewLine()
}

func (f *GrpcHandlersFile) genGrpcServeMethods() {
	for _, method := range f.Meta.SvcServerInterface.Methods {
		f.genGrpcServeMethod(method)
	}
}

// genEndpointMethod generate endpoint corresponding method
//
//func (h *Handlers) Hello(ctx context.Context, request *pb.HelloRequest) (*pb.HelloResponse, error) {
//	_, resp, err := h.HelloHandler.ServeGRPC(ctx, request)
//	return resp.(*pb.HelloResponse), err
//}
func (f *GrpcHandlersFile) genGrpcServeMethod(method parser.Method) {
	body := []jen.Code{
		jen.List(jen.Id("_"), jen.Id("resp"), jen.Id("err")).Op(":=").
			Id("h").Dot(method.Name+"Handler").Dot("ServeGRPC").Call(jen.Id("ctx"), jen.Id("request")),
		jen.If(jen.Err().Op("!=").Nil()).Block(
			jen.Return(jen.Nil(), jen.Err()),
		),
	}

	typeAsserts := f.TypeAssert("response", "resp", f.PbDir, method.Results[0].Type, "server grpc: type assert error")
	body = append(body, typeAsserts...)
	body = append(body, []jen.Code{
		jen.Return(jen.Id("response"), jen.Nil()),
	}...)

	f.Builder.AppendFunction(
		method.Name,
		jen.Id("h").Op("*").Id(HandlersStructName),
		[]jen.Code{
			jen.Id("ctx").Qual("context", "Context"),
			utils.GetValidParameterCode("request", f.PbDir, method.Parameters[1].Type),
		},
		[]jen.Code{
			jen.List(utils.GetValidParameterCode("", f.PbDir, method.Results[0].Type), jen.Error()),
		},
		"",
		body...,
	)
	f.Builder.NewLine()
}
