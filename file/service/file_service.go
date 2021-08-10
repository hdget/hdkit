package service

import (
	"fmt"
	"github.com/dave/jennifer/jen"
	"github.com/hdget/hdkit/g"
	"github.com/hdget/hdkit/generator"
	"github.com/hdget/hdkit/parser"
	"github.com/hdget/hdkit/utils"
)

type ServiceFile struct {
	*generator.BaseGenerator
	Meta  *generator.Meta
	PbDir string
}

const (
	ServiceFilename = "service.go"
)

func NewServiceFile(meta *generator.Meta) (generator.Generator, error) {
	baseGenerator, err := generator.NewBaseGenerator(meta.Dirs[g.Service], ServiceFilename, false)
	if err != nil {
		return nil, err
	}

	return &ServiceFile{
		BaseGenerator: baseGenerator,
		Meta:          meta,
		PbDir:         meta.Dirs[g.Pb],
	}, nil
}

func (f *ServiceFile) GetGenCodeFuncs() []func() {
	return []func(){
		f.genServiceStruct,
		f.genNewServiceFunction,
		f.genServiceMethods,
	}
}

func (f *ServiceFile) genServiceStruct() {
	found, _ := f.FindStructure(f.Meta.SvcStructName)
	if found == nil {
		f.Builder.AppendStruct(f.Meta.SvcStructName)
		f.genVarSvcStructImplService()
	}
}

// genVarSvcStructImplService add fake var definition which will give error prompts
// if xxxServiceImpl struct not implement all methods defined in XxxServiceServer interface
//
// var _ pb.XxxServiceServer = (*XxxServiceImpl)(nil)
func (f *ServiceFile) genVarSvcStructImplService() {
	f.Builder.Raw().Var().Id("_").Qual(f.PbDir, f.Meta.SvcServerInterface.Name).
		Op("=").Parens(jen.Op("*").Id(f.Meta.SvcStructName)).Call(jen.Nil()).Line()
}

func (f ServiceFile) genNewServiceFunction() {
	funcName := fmt.Sprintf("New%s", utils.ToCamelCase(f.Meta.RawSvcName))

	found, _ := f.FindMethod(funcName)
	if found == nil {
		f.Builder.Raw().Commentf(
			"%s returns a naive, stateless implementation of %s.",
			funcName,
			f.Meta.SvcServerInterfaceName,
		).Line()
		body := []jen.Code{
			jen.Return(jen.Id(fmt.Sprintf("&%s{}", f.Meta.SvcStructName))),
		}
		f.Builder.AppendFunction(
			funcName,
			nil,
			[]jen.Code{},
			[]jen.Code{jen.Qual(f.PbDir, f.Meta.SvcServerInterfaceName)},
			"",
			body...,
		)
		f.Builder.NewLine()
		f.Builder.NewLine()
	}
}

func (f ServiceFile) genServiceMethods() {
	existMethodNames := make([]string, 0)
	for _, existMethod := range f.ParsedFile.Methods {
		existMethodNames = append(existMethodNames, existMethod.Name)
	}

	tbdMethods := make([]parser.Method, 0)
	for _, pbMethod := range f.Meta.SvcServerInterface.Methods {
		if !utils.StringSliceContains(existMethodNames, pbMethod.Name) {
			tbdMethods = append(tbdMethods, pbMethod)
		}
	}

	for _, m := range tbdMethods {
		f.genServiceMethod(m)
	}
}

// genServiceMethod generate method function as below
//func (s SearchServiceImpl) Hello(ctx context.Context, request *pb.HelloRequest) (*pb.HelloResponse, error) {
//	return &pb.HelloResponse{
//		Response: "hello world",
//	}, nil
//}
func (f ServiceFile) genServiceMethod(method parser.Method) {
	f.Builder.AppendFunction(
		method.Name,
		jen.Id("impl").Op("*").Id(f.Meta.SvcStructName),
		[]jen.Code{
			jen.Id("ctx").Qual("context", "Context"),
			utils.GetValidParameterCode("request", f.PbDir, method.Parameters[1].Type),
		},
		[]jen.Code{
			jen.List(utils.GetValidParameterCode("", f.PbDir, method.Results[0].Type), jen.Id("error")),
		},
		"",
		jen.Return(jen.Nil(), jen.Nil()),
	)
	f.Builder.NewLine()
}
