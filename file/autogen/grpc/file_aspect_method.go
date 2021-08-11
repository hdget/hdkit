package grpc

import (
	"fmt"
	"github.com/dave/jennifer/jen"
	"github.com/hdget/hdkit/g"
	"github.com/hdget/hdkit/generator"
	"github.com/hdget/hdkit/parser"
	"github.com/hdget/hdkit/utils"
	"strings"
)

type AspectMethodFile struct {
	*generator.BaseGenerator
	Meta       *generator.Meta
	Method     parser.Method // service interface's method name
	StructName string
	PbDir      string
}

const (
	AspectSuffix = "Aspect"
)

func NewEndpointMethodFile(method parser.Method, meta *generator.Meta) (generator.Generator, error) {
	filename := fmt.Sprintf("%s_%s.go", strings.ToLower(AspectSuffix), strings.ToLower(method.Name))
	baseGenerator, err := generator.NewBaseGenerator(meta.Dirs[g.Grpc], filename, true)
	if err != nil {
		return nil, err
	}

	return &AspectMethodFile{
		BaseGenerator: baseGenerator,
		Meta:          meta,
		Method:        method,
		StructName:    utils.ToCamelCase(method.Name) + AspectSuffix,
		PbDir:         meta.Dirs[g.Pb],
	}, nil
}

func (f *AspectMethodFile) GetGenCodeFuncs() []func() {
	return []func(){
		f.genImports,
		f.genAspectStruct,
		f.genGetMethodNameFunc,
		f.genGetServiceNameFunc,
		f.genGetGrpcReplyTypeFunc,
		f.genMakeEndpointFunc,
		f.genServerDecodeRequestFunc,
		f.genServerEncodeResponseFunc,
		f.genClientEncodeRequestFunc,
		f.genClientDecodeResponseFunc,
	}
}

func (f *AspectMethodFile) genImports() {
	f.JenFile.ImportName(f.PbDir, "pb")
	f.JenFile.ImportName(g.ImportPaths[g.KitEndpoint], "endpoint")
	f.JenFile.ImportName(g.ImportPaths[g.Errors], "errors")
}

func (f *AspectMethodFile) genAspectStruct() {
	f.Builder.Raw().Type().Id(f.StructName).Struct().Line()
}


//func (ap SearchAspect) GetMethodName() string  {
//	return "Hello"
//}
func (f *AspectMethodFile) genGetMethodNameFunc() {
	f.Builder.AppendFunction(
		"GetMethodName",
		jen.Id("ap").Op("*").Id(f.StructName),
		nil,
		nil,
		"string",
		jen.Return(jen.Lit(f.Method.Name)),
	)
	f.Builder.NewLine()
}

//func (ap HelloAspect) GetServiceName() string  {
//	return "Hello"
//}
func (f *AspectMethodFile) genGetServiceNameFunc() {
	f.Builder.AppendFunction(
		"GetServiceName",
		jen.Id("ap").Op("*").Id(f.StructName),
		nil,
		nil,
		"string",
		jen.Return(jen.Lit("pb." + f.Meta.RawSvcName)),
	)
	f.Builder.NewLine()
}

//func (ap HelloAspect) GetServiceName() string  {
//	return "Hello"
//}
func (f *AspectMethodFile) genGetGrpcReplyTypeFunc() {
	f.Builder.AppendFunction(
		"GetGrpcReplyType",
		jen.Id("ap").Op("*").Id(f.StructName),
		nil,
		nil,
		"interface{}",
		jen.Return(
			jen.Qual(f.PbDir, f.Deference(f.Method.Results[0].Type)).Block(),
		),
	)
	f.Builder.NewLine()
}

// genMakeEndpointFunc generate MakeEndpoint function
//
//func MakeSearchEndpoint(svc interface{}) endpoint.Endpoint {
//	return func(ctx context.Context, request interface{}) (interface{}, error) {
//      s, ok := svc.(*SearchServiceImpl)
//      if !ok {
//         return nil, errors.New("invalid service")
//      }
//      req, ok := request.(*pb.ServiceRequest)
//      if !ok {
//         return nil, errors.New("invalid service request")
//      }
//      return s.Search(ctx, req)
//	}
//}
func (f *AspectMethodFile) genMakeEndpointFunc() {
	cg := generator.NewCodeBuilder(nil)
	body := f.TypeAssert("s", "svc", f.PbDir, f.Meta.SvcServerInterfaceName, "make endpoint: type assert error")
	body = append(body, f.TypeAssert("req", "request", f.PbDir, f.Method.Parameters[1].Type, "make endpoint: type assert error")...)
	body = append(body, jen.Return(jen.Qual("s", f.Method.Name).Call(jen.Id("ctx"), jen.Id("req"))))

	cg.AppendFunction(
		"",
		nil,
		[]jen.Code{
			jen.Id("ctx").Qual("context", "Context"),
			jen.Id("request").Interface(),
		},
		[]jen.Code{
			jen.Interface(),
			jen.Error(),
		},
		"",
		body...,
	)

	f.Builder.Raw().Commentf("MakeEndpoint returns an endpoint that invokes %s on the service.", f.Method.Name)
	f.Builder.NewLine()
	f.Builder.AppendFunction(
		"MakeEndpoint",
		jen.Id("ap").Op("*").Id(f.StructName),
		[]jen.Code{
			jen.Id("svc").Interface(),
		},
		[]jen.Code{
			jen.Qual(g.ImportPaths[g.KitEndpoint], "Endpoint"),
		},
		"",
		jen.Return(cg.Raw()),
	)
	f.Builder.NewLine()
}

// genServerDecodeRequestFunc generate ServerDecodeRequest function
//func (ep *SearchAspect) ServerDecodeRequest(ctx context.Context, request interface{}) (interface{}, error) {
//	return request.(*pb.SearchRequest), nil
//}
func (f *AspectMethodFile) genServerDecodeRequestFunc() {
	body := f.TypeAssert("req", "request", f.PbDir, f.Method.Parameters[1].Type, "server decode request: type assert error")
	body = append(body, jen.Return(jen.Id("req"), jen.Nil()))

	f.Builder.AppendFunction(
		"ServerDecodeRequest",
		jen.Id("ap").Op("*").Id(f.StructName),
		[]jen.Code{
			jen.Id("ctx").Qual("context", "Context"),
			jen.Id("request").Interface(),
		},
		[]jen.Code{
			jen.Interface(),
			jen.Error(),
		},
		"",
		body...,
	)
	f.Builder.NewLine()
}

// genServerEncodeResponseFunc
//func (s *SearchAspect) ServerEncodeResponse(ctx context.Context, response interface{}) (interface{}, error) {
//	return response.(*pb.SearchResponse), nil
//}
func (f *AspectMethodFile) genServerEncodeResponseFunc() {
	body := f.TypeAssert("resp", "response", f.PbDir, f.Method.Results[0].Type, "server encode response: type assert error")
	body = append(body, jen.Return(jen.Id("resp"), jen.Nil()))

	f.Builder.AppendFunction(
		"ServerEncodeResponse",
		jen.Id("ap").Op("*").Id(f.StructName),
		[]jen.Code{
			jen.Id("ctx").Qual("context", "Context"),
			jen.Id("response").Interface(),
		},
		[]jen.Code{
			jen.Interface(),
			jen.Error(),
		},
		"",
		body...,
	)
	f.Builder.NewLine()
}

// genClientEncodeRequestFunc
//func (ap HelloAspect) ClientEncodeRequest(ctx context.Context, request interface{}) (interface{}, error) {
//	return request.(*pb.HelloRequest), nil
//}
func (f *AspectMethodFile) genClientEncodeRequestFunc() {
	body := f.TypeAssert("req", "request", f.PbDir, f.Method.Parameters[1].Type, "client encode request: type assert error")
	body = append(body, jen.Return(jen.Id("req"), jen.Nil()))

	f.Builder.AppendFunction(
		"ClientEncodeRequest",
		jen.Id("ap").Op("*").Id(f.StructName),
		[]jen.Code{
			jen.Id("ctx").Qual("context", "Context"),
			jen.Id("request").Interface(),
		},
		[]jen.Code{
			jen.Interface(),
			jen.Error(),
		},
		"",
		body...,
	)
	f.Builder.NewLine()
}

// genClientEncodeRequestFunc
//func (ap HelloAspect) ClientEncodeRequest(ctx context.Context, request interface{}) (interface{}, error) {
//	return request.(*pb.HelloRequest), nil
//}
func (f *AspectMethodFile) genClientDecodeResponseFunc() {
	body := f.TypeAssert("resp", "response", f.PbDir, f.Method.Results[0].Type, "client decode response: type assert error")
	body = append(body, jen.Return(jen.Id("resp"), jen.Nil()))

	f.Builder.AppendFunction(
		"ClientDecodeResponse",
		jen.Id("ap").Op("*").Id(f.StructName),
		[]jen.Code{
			jen.Id("ctx").Qual("context", "Context"),
			jen.Id("response").Interface(),
		},
		[]jen.Code{
			jen.Interface(),
			jen.Error(),
		},
		"",
		body...,
	)
	f.Builder.NewLine()
}
