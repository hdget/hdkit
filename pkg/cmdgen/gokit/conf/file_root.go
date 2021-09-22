package conf

import (
	"github.com/dave/jennifer/jen"
	"github.com/hdget/hdkit/g"
	"github.com/hdget/hdkit/generator"
)

type RootConfigFile struct {
	*generator.BaseGenerator
	Meta             *generator.Meta
	ConfigStructName string // config struct name
}

const (
	RootConfigFilename = "root.go"
)

func NewRootConfigFile(meta *generator.Meta) (generator.Generator, error) {
	baseGenerator, err := generator.NewBaseGenerator(g.GetDir(meta.RootDir, g.Config), RootConfigFilename, false)
	if err != nil {
		return nil, err
	}

	return &RootConfigFile{
		BaseGenerator:    baseGenerator,
		Meta:             meta,
		ConfigStructName: meta.RawSvcName + "Config",
	}, nil
}

func (f *RootConfigFile) GetGenCodeFuncs() []func() {
	return []func(){
		f.genImports,
		f.genConfigDefines,
	}
}

func (f *RootConfigFile) genImports() {
	f.JenFile.ImportName(g.ImportPaths[g.HdSdk], "hdsdk")
}

//type XxxServiceConfig struct {
//	hdsdk.Config `mapstructure:",squash"`
//}
func (f *RootConfigFile) genConfigDefines() {
	found, _ := f.FindStructure(f.ConfigStructName)
	if found == nil {
		f.Builder.AppendStruct(
			f.ConfigStructName,
			jen.Qual(g.ImportPaths[g.HdSdk], "Config").Tag(map[string]string{"mapstructure": ",squash"}),
		)
	}
}
