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

type EndpointMethodFile struct {
	*generator.BaseGenerator
	Meta       *generator.Meta
	Method     parser.Method // service interface's method name
	StructName string
	PbDir      string
}

func NewEndpointMethodFile(method parser.Method, meta *generator.Meta) (generator.Generator, error) {
	filename := fmt.Sprintf("endpoint_%s.go", strings.ToLower(method.Name))
	baseGenerator, err := generator.NewBaseGenerator(meta.Dirs[g.Grpc], filename, true)
	if err != nil {
		return nil, err
	}

	return &EndpointMethodFile{
		BaseGenerator: baseGenerator,
		Meta:          meta,
		Method:        method,
		StructName:    utils.ToCamelCase(method.Name) + "Endpoint",
		PbDir:         meta.Dirs[g.Pb],
	}, nil
}

func (f *EndpointMethodFile) GetGenCodeFuncs() []func() {
	return []func(){
		f.genEndpointStruct,
		f.genGetNameFunc,
		f.genMakeEndpointFunction,
		f.genServerDecodeRequest,
		f.genServerEncodeResponse,
	}
}

func (f *EndpointMethodFile) genEndpointStruct() {
	f.Builder.Raw().Type().Id(f.StructName).Struct().Line()
}

// genMakeEndpointFunction generate MakeEndpoint function
//
//func MakeHelloEndpoint(svc interface{}) endpoint.Endpoint {
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
func (f *EndpointMethodFile) genMakeEndpointFunction() {
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
		jen.Id("ep").Op("*").Id(f.StructName),
		[]jen.Code{
			jen.Id("svc").Interface(),
		},
		[]jen.Code{
			jen.Qual("github.com/go-kit/kit/endpoint", "Endpoint"),
		},
		"",
		jen.Return(cg.Raw()),
	)
	f.Builder.NewLine()
}

// genServerDecodeRequest generate ServerDecodeRequest function
//func (ep *SearchEndpoint) ServerDecodeRequest(ctx context.Context, request interface{}) (interface{}, error) {
//	return request.(*pb.SearchRequest), nil
//}
func (f *EndpointMethodFile) genServerDecodeRequest() {
	body := f.TypeAssert("req", "request", f.PbDir, f.Method.Parameters[1].Type, "decode request: type assert error")
	body = append(body, jen.Return(jen.Id("req"), jen.Nil()))

	f.Builder.AppendFunction(
		"ServerDecodeRequest",
		jen.Id("ep").Op("*").Id(f.StructName),
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

// genServerEncodeResponse
//func (s *SearchEndpoint) ServerEncodeResponse(ctx context.Context, response interface{}) (interface{}, error) {
//	return response.(*pb.SearchResponse), nil
//}
func (f *EndpointMethodFile) genServerEncodeResponse() {
	body := f.TypeAssert("resp", "response", f.PbDir, f.Method.Results[0].Type, "encode response: type assert error")
	body = append(body, jen.Return(jen.Id("resp"), jen.Nil()))

	f.Builder.AppendFunction(
		"ServerEncodeResponse",
		jen.Id("ep").Op("*").Id(f.StructName),
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

//func (h HelloHandler) GetName() string  {
//	return "hello"
//}
func (f *EndpointMethodFile) genGetNameFunc() {
	f.Builder.AppendFunction(
		"GetName",
		jen.Id("ep").Op("*").Id(f.Method.Name+"Endpoint"),
		nil,
		nil,
		"string",
		jen.Return(jen.Lit(f.Method.Name)),
	)
	f.Builder.NewLine()
}
