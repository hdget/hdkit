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
	HandlersFileName    = "handlers.go"
	HandlersStructName  = "Handlers"
	HandlerStructSuffix = "Handler"
	EndpointSuffix      = "Endpoint"
)

var (
	HandlersStructComments = []string{
		"Handlers collects all of the handlers that compose a service.",
	}
)

func NewHandlersFile(meta *generator.Meta) (generator.Generator, error) {
	baseGenerator, err := generator.NewBaseGenerator(meta.Dirs[g.Grpc], HandlersFileName, true)
	if err != nil {
		return nil, err
	}

	return &GrpcHandlersFile{
		BaseGenerator: baseGenerator,
		Meta:          meta,
		StructName:    HandlersStructName,
		PbDir:         meta.Dirs[g.Pb],
	}, nil
}

func (f *GrpcHandlersFile) PreGenerate() error {
	err := f.BaseGenerator.PreGenerate()
	if err != nil {
		return err
	}

	// package a should be aliased to "b"
	f.JenFile.ImportAlias("github.com/hdget/sdk/types", "sdktypes")
	f.JenFile.ImportAlias("github.com/go-kit/kit/transport/grpc", "kitgrpc")
	return nil
}

func (f *GrpcHandlersFile) GetGenCodeFuncs() []func() {
	return []func(){
		f.genHandlersStructure,
		f.genNewHandlersFunction,
		f.genGrpcServeMethods,
	}
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
// func NewHandlers(sever MsGrpcServer, svc Service) *Handlers {
//    return &Handlers{
//   			SearchHandler: sever.CreateHandler(svc, NewServerHandler()),
//				HelloHandler: sever.CreateHandler(svc, NewHelloHandler()),
//	  }
// }
func (f *GrpcHandlersFile) genNewHandlersFunction() {
	handlers := jen.Dict{}
	for _, m := range f.Meta.SvcServerInterface.Methods {
		handlerName := m.Name + HandlerStructSuffix
		endpointName := m.Name + EndpointSuffix
		handlers[jen.Id(handlerName)] = jen.Qual("server", "CreateHandler").Call(
			jen.Id("svc"), jen.Id("&"+endpointName+"{}"))
	}

	body := jen.Return(jen.Op("&").Id("Handlers").Values(handlers))

	f.Builder.Raw().
		Func().
		Id("NewHandlers").
		Params(
			jen.Id("server").Qual("github.com/hdget/sdk/types", "MsGrpcServer"),
			jen.Id("svc").Qual(f.PbDir, f.Meta.SvcServerInterface.Name),
		).
		Op("*").Id("Handlers").
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
		jen.Id("h").Op("*").Id("Handlers"),
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
