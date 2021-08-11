package grpc

import (
	"github.com/dave/jennifer/jen"
	"github.com/hdget/hdkit/g"
	"github.com/hdget/hdkit/generator"
)

type GrpcClientFile struct {
	*generator.BaseGenerator
	Meta       *generator.Meta
	StructName string
	PbDir      string
}

const (
	ClientFileName    = "client.go"
	ClientStructName  = "Client"
)

var (
	ClientStructComments = []string{
		"Client collects all of the handlers that compose a service.",
	}
)

func NewClientFile(meta *generator.Meta) (generator.Generator, error) {
	baseGenerator, err := generator.NewBaseGenerator(meta.Dirs[g.Grpc], ClientFileName, true)
	if err != nil {
		return nil, err
	}

	return &GrpcClientFile{
		BaseGenerator: baseGenerator,
		Meta:          meta,
		StructName:    ClientStructName,
		PbDir:         meta.Dirs[g.Pb],
	}, nil
}

func (f *GrpcClientFile) GetGenCodeFuncs() []func() {
	return []func(){
		f.genImports,
		f.genClientStruct,
		f.genNewClientFunc,
	}
}

func (f *GrpcClientFile) genImports() {
	f.JenFile.ImportName(f.PbDir, "pb")
	f.JenFile.ImportName(g.ImportPaths[g.Errors], "errors")
	f.JenFile.ImportName(g.ImportPaths[g.StdGrpc], "grpc")
	f.JenFile.ImportName(g.ImportPaths[g.HdSdk], "hdsdk")
	f.JenFile.ImportName(g.ImportPaths[g.KitEndpoint], "endpoint")
}

// genClientStruct collects all of the handlers that compose a service.
//type Client struct {
//	Search endpoint.Endpoint
//	Hello  endpoint.Endpoint
//}
func (f *GrpcClientFile) genClientStruct() {
	f.JenFile.PackageComment("Package grpc THIS FILE IS AUTO GENERATED PLEASE DO NOT EDIT!!")

	fields := make([]jen.Code, 0)
	for _, v := range f.Meta.SvcServerInterface.Methods {
		fields = append(fields, jen.Id(v.Name).Qual(g.ImportPaths[g.KitEndpoint], "Endpoint"))
	}

	f.Builder.AppendMultilineComment(ClientStructComments)
	f.Builder.NewLine()
	f.Builder.AppendStruct(
		f.StructName,
		fields...,
	)
}

// genNewClientFunc returns new client function
//func NewClient(conn *grpc.ClientConn, args ...string) (*Client, error) {
//   var name string
//   if len(args) > 0 {
//     name = args[0]
//   }
//
//  var ms types.MicroService
//  if name == "" {
//    ms = hdsdk.MicroService.My()
//  }else{
//    ms = hdsdk.MicroService.By(name)
//  }
//  if ms == nil {
//     return nil, fmt.Errorf("microservice:%s not found", name)
//  }
//
//	manager := ms.NewGrpcClientManager()
//	if manager == nil {
//		return nil, fmt.Errorf("new grpc client manager:%s failed", name)
//	}
//
//	return &Client{
//		Search: manager.CreateEndpoint(conn, &SearchAspect{}),
//		Hello:  manager.CreateEndpoint(conn, &HelloAspect{}),
//	}, nil
//}
func (f *GrpcClientFile) genNewClientFunc() {
	body := []jen.Code{
		jen.Var().Id("name").String(),
		jen.If(jen.Id("len").Call(jen.Id("args")).Op(">").Lit(0)).Block(
			jen.Id("name").Op("=").Id("args").Index(jen.Lit(0)),
		),
		jen.Line(),
		jen.Var().Id("ms").Qual(g.ImportPaths[g.HdSdkTypes], "MicroService"),
		jen.If(jen.Id("name").Op("==").Lit("")).Block(
			jen.Id("ms").Op("=").Qual(g.ImportPaths[g.HdSdk], "MicroService").Dot("My").Call(),
		).Else().Block(
			jen.Id("ms").Op("=").Qual(g.ImportPaths[g.HdSdk], "MicroService").Dot("By").Call(jen.Id("name")),
		),
		jen.Line(),
		jen.Id("manager").Op(":=").Qual("ms", "NewGrpcClientManager").Call(),
		jen.If(jen.Id("manager").Op("==").Nil()).Block(
			jen.Return(jen.Nil(), jen.Qual("fmt", "Errorf").Call(jen.Lit("new grpc client manager:%s failed"), jen.Id("name"))),
		),
		jen.Line(),
	}

	clients := jen.Dict{}
	for _, m := range f.Meta.SvcServerInterface.Methods {
		aspectStructName := m.Name + AspectSuffix
		clients[jen.Id(m.Name)] = jen.Qual("manager", "CreateEndpoint").Call(
			jen.Id("conn"), jen.Id("&" + aspectStructName + "{}"))
	}

	body = append(body,
		jen.Return(
			jen.Op("&").Id("Client").Values(clients),
			jen.Nil(),
		),
	)

	f.Builder.AppendFunction(
		"NewClient",
		nil,
		[]jen.Code{
			jen.Id("conn").Op("*").Qual(g.ImportPaths[g.StdGrpc], "ClientConn"),
			jen.Id("args").Op("...").String(),
		},
		[]jen.Code{
			jen.Op("*").Id("Client"),
			jen.Error(),
		},
		"",
		body...
	)
}
