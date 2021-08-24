package http

import (
	"github.com/dave/jennifer/jen"
	"github.com/hdget/hdkit/g"
	"github.com/hdget/hdkit/generator"
	"github.com/hdget/hdkit/utils"
	"path"
)

type HttpHandlersFile struct {
	*generator.BaseGenerator
	Meta       *generator.Meta
	StructName string
	PbDir      string
}

const (
	HandlersFileName    = "http_handlers.go"
	HandlersStructName  = "HttpHandlers"
	HandlerStructSuffix = "Handler"
)

var (
	HandlersStructComments = []string{
		"Handlers collects all of the handlers that compose a service.",
	}
)

func NewHttpHandlersFile(meta *generator.Meta) (generator.Generator, error) {
	baseGenerator, err := generator.NewBaseGenerator(g.GetDir(meta.RootDir, g.Http), HandlersFileName, true)
	if err != nil {
		return nil, err
	}

	return &HttpHandlersFile{
		BaseGenerator: baseGenerator,
		Meta:          meta,
		StructName:    HandlersStructName,
		PbDir:         g.GetDir(meta.RootDir, g.Pb),
	}, nil
}

func (f *HttpHandlersFile) GetGenCodeFuncs() []func() {
	return []func(){
		f.genImports,
		f.genNewHandlersFunction,
	}
}

func (f *HttpHandlersFile) genImports() {
	f.JenFile.ImportName(f.PbDir, "pb")
	f.JenFile.ImportName(g.ImportPaths[g.Errors], "errors")
	f.JenFile.ImportName(g.ImportPaths[g.HdSdkTypes], "types")
	f.JenFile.ImportAlias(g.ImportPaths[g.KitGrpc], "kithttp")
}

// genNewHandlersFunction returns new handlers function
//
// func NewHandlers(manager HttpServerManager, ap HttpAspect) map[string]http.Handler {
//    return map[string]{
//   			"Search": manager.CreateHandler(svc, &SearchAspect{}),
//				"Hello": sever.CreateHandler(svc, &HelloAspect{}),
//	  }
// }
func (f *HttpHandlersFile) genNewHandlersFunction() {
	handlers := jen.Dict{}
	for _, m := range f.Meta.SvcServerInterface.Methods {
		url := "/" + path.Join(utils.ToSnakeCase(f.Meta.RawSvcName), utils.ToSnakeCase(m.Name))
		aspectName := m.Name + AspectSuffix
		handlers[jen.Lit(url)] = jen.Qual("manager", "CreateHandler").Call(
			jen.Id("svc"), jen.Id("&"+aspectName+"{}"))
	}

	body := jen.Return(jen.Map(jen.String()).Qual(g.ImportPaths[g.StdHttp], "Handler").Values(handlers))

	f.Builder.AppendFunction(
		"NewHandlers",
		nil,
		[]jen.Code{
			jen.Id("manager").Qual(g.ImportPaths[g.HdSdkTypes], "HttpServerManager"),
			jen.Id("svc").Qual(f.PbDir, f.Meta.SvcServerInterface.Name),
		},
		[]jen.Code{
			jen.Map(jen.String()).Qual(g.ImportPaths[g.StdHttp], "Handler"),
		},
		"",
		body,
	)
	f.Builder.NewLine()
}
