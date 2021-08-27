package service

import (
	"github.com/dave/jennifer/jen"
	"github.com/hdget/hdkit/g"
	"github.com/hdget/hdkit/generator"
	"github.com/hdget/hdkit/parser"
	"github.com/hdget/hdkit/utils"
)

type InvocationHandlersFile struct {
	*generator.BaseGenerator
	Meta  *generator.Meta
	PbDir string
}

const (
	InvocationStruct          = "Invocation"
	NewInvocation             = "NewInvocation"
	InvocationHandlerFilename = "invocation_handlers.go"
)

func NewInvocationHandlersFile(meta *generator.Meta) (generator.Generator, error) {
	baseGenerator, err := generator.NewBaseGenerator(g.GetDir(meta.RootDir, g.Service), InvocationHandlerFilename, false)
	if err != nil {
		return nil, err
	}

	return &InvocationHandlersFile{
		BaseGenerator: baseGenerator,
		Meta:          meta,
		PbDir:         g.GetDir(meta.RootDir, g.Pb),
	}, nil
}

func (f *InvocationHandlersFile) GetGenCodeFuncs() []func() {
	return []func(){
		f.genImports,
		f.genHandlersStruct,
		f.genNewHandlersFunc,
		f.genInvocationHandlers,
	}
}

func (f *InvocationHandlersFile) genImports() {
	f.JenFile.ImportName(g.ImportPaths[g.DaprCommon], "common")
}

func (f *InvocationHandlersFile) genHandlersStruct() {
	found, _ := f.FindStructure(InvocationStruct)
	if found == nil {
		f.Builder.Raw().Type().Id(InvocationStruct).Struct()
		f.Builder.NewLine()
		f.Builder.NewLine()
	}
}

func (f *InvocationHandlersFile) genNewHandlersFunc() {
	found, _ := f.FindMethod(NewInvocation)
	if found == nil {
		body := []jen.Code{
			jen.Return(jen.Op("&").Id(InvocationStruct).Block()),
		}

		f.Builder.AppendFunction(
			NewInvocation,
			nil,
			nil,
			nil,
			"*"+InvocationStruct,
			body...,
		)
		f.Builder.NewLine()
		f.Builder.NewLine()
	}
}

func (f *InvocationHandlersFile) genInvocationHandlers() {
	for _, method := range f.Meta.SvcServerInterface.Methods {
		f.genInvocationHandler(method)
	}
}

// genInvocationHandler generate invocation handlers
//func echoHandler(ctx context.Context, in *common.InvocationEvent) (*common.Content, error) {
//	log.Printf("echo - ContentType:%s, Verb:%s, QueryString:%s, %+v", in.ContentType, in.Verb, in.QueryString, string(in.Data))
//	// do something with the invocation here
//	out = &common.Content{
//		Data:        in.Data,
//		ContentType: in.ContentType,
//		DataTypeURL: in.DataTypeURL,
//	}
//	return
//}
func (f *InvocationHandlersFile) genInvocationHandler(method parser.Method) {
	found, _ := f.FindMethod(utils.ToLowerSnakeCase(method.Name) + "Handler")
	if found == nil {
		body := []jen.Code{
			jen.Return(
				jen.Op("&").Qual(g.ImportPaths[g.DaprCommon], "Content").Values(
					jen.Dict{
						jen.Id("Data"):        jen.Id("event").Dot("Data"),
						jen.Id("ContentType"): jen.Id("event").Dot("ContentType"),
						jen.Id("DataTypeURL"): jen.Id("event").Dot("DataTypeURL"),
					},
				),
				jen.Nil(),
			),
		}

		f.Builder.AppendFunction(
			utils.ToLowerFirstCamelCase(method.Name)+"Handler",
			jen.Id("h").Id(InvocationStruct),
			[]jen.Code{
				jen.Id("ctx").Qual("context", "Context"),
				jen.Id("event").Op("*").Qual(g.ImportPaths[g.DaprCommon], "InvocationEvent"),
			},
			[]jen.Code{
				jen.Op("*").Qual(g.ImportPaths[g.DaprCommon], "Content"),
				jen.Id("error"),
			},
			"",
			body...,
		)
		f.Builder.NewLine()
		f.Builder.NewLine()
	}
}
