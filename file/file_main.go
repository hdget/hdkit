package file

import (
	"github.com/dave/jennifer/jen"
	"github.com/hdget/hdkit/g"
	"github.com/hdget/hdkit/generator"
)

type MainFile struct {
	*generator.BaseGenerator
	Meta   *generator.Meta
	CmdDir string
}

const (
	MainPackageName = "main"
	MainFilename    = "main.go"
)

func NewMainFile(meta *generator.Meta) (generator.Generator, error) {
	baseGenerator, err := generator.NewBaseGenerator(meta.RootDir, MainFilename, true, MainPackageName)
	if err != nil {
		return nil, err
	}

	return &MainFile{
		BaseGenerator: baseGenerator,
		Meta:          meta,
		CmdDir:        meta.Dirs[g.Cmd],
	}, nil
}

func (f MainFile) GetGenCodeFuncs() []func() {
	return []func(){
		f.genImports,
		f.genMain,
	}
}

func (f *MainFile) genImports() {
	f.JenFile.ImportName(f.CmdDir, "cmd")
}

//func main() {
//	cmd.Execute()
//}
func (f MainFile) genMain() {
	f.Builder.AppendFunction(
		"main",
		nil,
		nil,
		nil,
		"",
		jen.Qual(f.CmdDir, "Execute").Call(),
	)
	f.Builder.NewLine()
}
