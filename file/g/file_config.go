package g

import (
	"github.com/dave/jennifer/jen"
	"github.com/hdget/hdkit/g"
	"github.com/hdget/hdkit/generator"
)

type GConfigFile struct {
	*generator.BaseGenerator
	Meta             *generator.Meta
	ConfigStructName string // config struct name
	PbDir            string
}

const (
	GConfigFilename = "config.go"
)

func NewGConfigFile(meta *generator.Meta) (generator.Generator, error) {
	baseGenerator, err := generator.NewBaseGenerator(meta.Dirs[g.Global], GConfigFilename, false)
	if err != nil {
		return nil, err
	}

	return &GConfigFile{
		BaseGenerator:    baseGenerator,
		Meta:             meta,
		ConfigStructName: meta.RawSvcName + "Config",
		PbDir:            meta.Dirs[g.Pb],
	}, nil
}

func (f *GConfigFile) GetGenCodeFuncs() []func() {
	return []func(){
		f.genConfigDefines,
	}
}

//type XxxServiceConfig struct {
//	hdsdk.Config `mapstructure:",squash"`
//}
func (f *GConfigFile) genConfigDefines() {
	found, _ := f.FindStructure(f.ConfigStructName)
	if found == nil {
		f.Builder.AppendStruct(
			f.ConfigStructName,
			jen.Qual(g.ImportPaths[g.HdSdk], "Config").Tag(map[string]string{"mapstructure": ",squash"}),
		)

		// add `var Config *XxxServiceConfig
		f.Builder.Raw().Var().Id("Config").Op("*").Id(f.ConfigStructName)
	}
}
