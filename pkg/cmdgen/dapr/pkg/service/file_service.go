package service

import (
	"fmt"
	"github.com/dave/jennifer/jen"
	"github.com/hdget/hdkit/g"
	"github.com/hdget/hdkit/generator"
	"github.com/hdget/hdkit/utils"
)

type ServiceFile struct {
	*generator.BaseGenerator
	Meta  *generator.Meta
	PbDir string
}

const (
	DaprService     = "DaprService"
	ServiceFilename = "service.go"
)

func NewServiceFile(meta *generator.Meta) (generator.Generator, error) {
	baseGenerator, err := generator.NewBaseGenerator(g.GetDir(meta.RootDir, g.Service), ServiceFilename, false)
	if err != nil {
		return nil, err
	}

	return &ServiceFile{
		BaseGenerator: baseGenerator,
		Meta:          meta,
		PbDir:         g.GetDir(meta.RootDir, g.Pb),
	}, nil
}

func (f *ServiceFile) GetGenCodeFuncs() []func() {
	return []func(){
		f.genImports,
		f.genFuncTypes,
		f.genServiceInterface,
		f.genServiceStruct,
		f.genNewServiceFunction,
		f.genServiceMethods,
	}
}

func (f *ServiceFile) genImports() {
	f.JenFile.ImportName(g.ImportPaths[g.DaprCommon], "common")
}

// type InvocationHandler func echoHandler(ctx context.Context, in *common.InvocationEvent) (out *common.Content, err error)
func (f *ServiceFile) genFuncTypes() {
	found, _ := f.FindFuncType("InvocationHandler")
	if found == nil {
		f.Builder.Raw().Type().Id("InvocationHandler").Func().Params(
			jen.List(
				jen.Id("ctx").Qual("context", "Context"),
				jen.Id("event").Op("*").Qual(g.ImportPaths[g.DaprCommon], "InvocationEvent"),
			),
		).Parens(
			jen.List(
				jen.Id("out").Op("*").Qual(g.ImportPaths[g.DaprCommon], "Content"),
				jen.Id("err").Id("error"),
			),
		).Line()

		f.Builder.Raw().Type().Id("BindingHandler").Func().Params(
			jen.List(
				jen.Id("ctx").Qual("context", "Context"),
				jen.Id("event").Op("*").Qual(g.ImportPaths[g.DaprCommon], "BindingEvent"),
			),
		).Parens(
			jen.List(
				jen.Id("out").Index().Id("byte"),
				jen.Id("err").Id("error"),
			),
		).Line()

		f.Builder.Raw().Type().Id("EventHandler").Func().Params(
			jen.List(
				jen.Id("ctx").Qual("context", "Context"),
				jen.Id("event").Op("*").Qual(g.ImportPaths[g.DaprCommon], "TopicEvent"),
			),
		).Parens(
			jen.List(
				jen.Id("retry").Bool(),
				jen.Id("err").Id("error"),
			),
		).Line().Line()

		f.Builder.Raw().Type().Id("Event").Struct(
			jen.Id("Sub").Op("*").Qual(g.ImportPaths[g.DaprCommon], "Subscription"),
			jen.Id("Handler").Id("EventHandler"),
		).Line().Line()
	}
}

func (f *ServiceFile) genServiceInterface() {
	found, _ := f.FindInterface(DaprService)
	if found == nil {
		f.Builder.AppendInterface(
			DaprService,
			[]jen.Code{
				jen.Id("GetInvocationHandlers").Params().Map(jen.String()).Id("InvocationHandler"),
				jen.Id("GetBindingHandlers").Params().Map(jen.String()).Id("BindingHandler"),
				jen.Id("GetEvents").Params().Index().Id("Event"),
			},
		)
		//f.genVarSvcStructImplService()
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
	f.Builder.Raw().Var().Id("_").Id(DaprService).
		Op("=").Parens(jen.Op("*").Id(f.Meta.SvcStructName)).Call(jen.Nil()).Line().Line()
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
			[]jen.Code{jen.Id(DaprService)},
			"",
			body...,
		)
		f.Builder.NewLine()
		f.Builder.NewLine()
	}
}

func (f ServiceFile) genServiceMethods() {
	found, _ := f.FindMethod("GetInvocationHandlers")
	if found == nil {
		invocationHandlers := f.getHandlerValues()

		f.Builder.AppendFunction(
			"GetInvocationHandlers",
			jen.Id("impl").Op("*").Id(f.Meta.SvcStructName),
			nil,
			[]jen.Code{
				jen.Map(jen.String()).Id("InvocationHandler"),
			},
			"",
			[]jen.Code{
				jen.Id("invoke").Op(":=").Id(NewInvocation).Call(),
				jen.Return(jen.Map(jen.String()).Id("InvocationHandler").Values(invocationHandlers)),
			}...,
		)
		f.Builder.NewLine()
		f.Builder.NewLine()

		f.Builder.AppendFunction(
			"GetBindingHandlers",
			jen.Id("impl").Op("*").Id(f.Meta.SvcStructName),
			nil,
			[]jen.Code{
				jen.Map(jen.String()).Id("BindingHandler"),
			},
			"",
			jen.Return(jen.Map(jen.String()).Id("BindingHandler").Values()),
		)
		f.Builder.NewLine()
		f.Builder.NewLine()

		f.Builder.AppendFunction(
			"GetEvents",
			jen.Id("impl").Op("*").Id(f.Meta.SvcStructName),
			nil,
			[]jen.Code{
				jen.Index().Id("Event"),
			},
			"",
			jen.Return(jen.Index().Id("Event").Block()),
		)
		f.Builder.NewLine()
		f.Builder.NewLine()
	}

}

// genServiceMethod generate method function as below
//func (impl *userCenterServiceImpl) GetInvocationHandlers() map[string]InvocationHandler {
//	return map[string]InvocationHandler{
//		"xxx": nil,
//	}
//}
func (f ServiceFile) getHandlerValues() jen.Dict {
	return jen.DictFunc(func(d jen.Dict) {
		for _, method := range f.Meta.SvcServerInterface.Methods {
			d[jen.Lit(method.Name)] = jen.Id("invoke").Dot(method.Name + "Handler")
		}
	})
}
