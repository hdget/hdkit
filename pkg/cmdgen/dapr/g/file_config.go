package g

import (
	"github.com/hdget/hdkit/g"
	"github.com/hdget/hdkit/generator"
)

type GConfigFile struct {
	*generator.BaseGenerator
	Meta             *generator.Meta
	ConfigStructName string // config struct name
	ConfDir          string
}

const (
	GConfigFilename = "config.go"
)

func NewGConfigFile(meta *generator.Meta) (generator.Generator, error) {
	baseGenerator, err := generator.NewBaseGenerator(g.GetDir(meta.RootDir, g.Global), GConfigFilename, false)
	if err != nil {
		return nil, err
	}

	return &GConfigFile{
		BaseGenerator:    baseGenerator,
		Meta:             meta,
		ConfigStructName: meta.RawSvcName + "Config",
		ConfDir:          g.GetDir(meta.RootDir, g.Config),
	}, nil
}

func (f *GConfigFile) GetGenCodeFuncs() []func() {
	return []func(){
		f.genImports,
		f.genConfigDefines,
	}
}

func (f *GConfigFile) genImports() {
	f.JenFile.ImportName(g.ImportPaths[g.HdSdk], "hdsdk")
	f.JenFile.ImportName(f.ConfDir, "conf")
}

//type XxxServiceConfig struct {
//	hdsdk.Config `mapstructure:",squash"`
//}
func (f *GConfigFile) genConfigDefines() {
	found, _ := f.FindVar("Config")
	if found == nil {
		// add `var Config *XxxServiceConfig
		f.Builder.Raw().Var().Id("Config").Op("*").Qual(f.ConfDir, f.ConfigStructName)
	}
}
