package http

import (
	"fmt"
	"github.com/dave/jennifer/jen"
	"github.com/hdget/hdkit/g"
	"github.com/hdget/hdkit/generator"
	"github.com/hdget/hdkit/parser"
	"github.com/hdget/hdkit/utils"
	"strings"
)

type HttpAspectMethodFile struct {
	*generator.BaseGenerator
	Meta       *generator.Meta
	Method     parser.Method // service interface's method name
	StructName string
	PbDir      string
}

const (
	AspectSuffix = "Aspect"
)

func NewHttpAspectMethodFile(method parser.Method, meta *generator.Meta) (generator.Generator, error) {
	filename := fmt.Sprintf("http_%s_%s.go", strings.ToLower(AspectSuffix), strings.ToLower(method.Name))
	baseGenerator, err := generator.NewBaseGenerator(g.GetDir(meta.RootDir, g.Http), filename, true)
	if err != nil {
		return nil, err
	}

	return &HttpAspectMethodFile{
		BaseGenerator: baseGenerator,
		Meta:          meta,
		Method:        method,
		StructName:    utils.ToCamelCase(method.Name) + AspectSuffix,
		PbDir:         g.GetDir(meta.RootDir, g.Pb),
	}, nil
}

func (f *HttpAspectMethodFile) GetGenCodeFuncs() []func() {
	return []func(){
		f.genImports,
		f.genAspectStruct,
		f.genGetMethodNameFunc,
		f.genMakeEndpointFunc,
		f.genServerDecodeRequestFunc,
		f.genServerEncodeResponseFunc,
	}
}

func (f *HttpAspectMethodFile) genImports() {
	f.JenFile.ImportName(f.PbDir, "pb")
	f.JenFile.ImportName(g.ImportPaths[g.KitEndpoint], "endpoint")
	f.JenFile.ImportName(g.ImportPaths[g.Errors], "errors")
}

func (f *HttpAspectMethodFile) genAspectStruct() {
	f.Builder.Raw().Type().Id(f.StructName).Struct().Line()
}

//func (ap SearchAspect) GetMethodName() string  {
//	return "Hello"
//}
func (f *HttpAspectMethodFile) genGetMethodNameFunc() {
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
func (f *HttpAspectMethodFile) genMakeEndpointFunc() {
	cg := generator.NewCodeBuilder(nil)
	body := f.TypeAssert("s", "svc", f.PbDir, f.Meta.SvcServerInterfaceName, "make endpoint: type assert error")
	body = append(body, f.TypeAssert("req", "request", f.PbDir, f.Deference(f.Method.Parameters[1].Type), "make endpoint: type assert error")...)
	body = append(body, jen.Return(jen.Qual("s", f.Method.Name).Call(jen.Id("ctx"), jen.Op("&").Id("req"))))

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
//func decodeHTTPSumRequest(_ context.Context, request *http.Request) (interface{}, error) {
//	var req pb.SumRequest
//	err := json.NewDecoder(request.Body).Decode(&req)
//	return req, err
//}
func (f *HttpAspectMethodFile) genServerDecodeRequestFunc() {
	body := []jen.Code{
		jen.Var().Id("req").Qual(f.PbDir, f.Deference(f.Method.Parameters[1].Type)),
		jen.Err().Op(":=").Qual(g.ImportPaths[g.StdJson], "NewDecoder").Call(
			jen.Qual("request", "Body")).Dot("Decode").Call(jen.Op("&").Id("req")),
		jen.Return(jen.Id("req"), jen.Err()),
	}

	f.Builder.AppendFunction(
		"ServerDecodeRequest",
		jen.Id("ap").Op("*").Id(f.StructName),
		[]jen.Code{
			jen.Id("ctx").Qual("context", "Context"),
			jen.Id("request").Op("*").Qual(g.ImportPaths[g.StdHttp], "Request"),
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
//func encodeHTTPGenericResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
//	w.Header().Set("Content-Type", "application/json; charset=utf-8")
//	return json.NewEncoder(w).Encode(response)
//}
func (f *HttpAspectMethodFile) genServerEncodeResponseFunc() {
	body := []jen.Code{
		jen.Qual("w", "Header").Call().Dot("Set").Parens(
			jen.List(
				jen.Lit("Content-Type"), jen.Lit("application/json; charset=utf-8"),
			),
		),
		jen.Return(
			jen.Qual(g.ImportPaths[g.StdJson], "NewEncoder").Call(jen.Id("w")).Dot("Encode").Call(jen.Id("response")),
		),
	}

	f.Builder.AppendFunction(
		"ServerEncodeResponse",
		jen.Id("ap").Op("*").Id(f.StructName),
		[]jen.Code{
			jen.Id("ctx").Qual("context", "Context"),
			jen.Id("w").Qual(g.ImportPaths[g.StdHttp], "ResponseWriter"),
			jen.Id("response").Interface(),
		},
		[]jen.Code{
			jen.Error(),
		},
		"",
		body...,
	)
	f.Builder.NewLine()
}
