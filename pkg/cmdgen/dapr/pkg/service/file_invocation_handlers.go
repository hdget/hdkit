package service

import (
	"github.com/hdget/hdkit/g"
	"github.com/hdget/hdkit/generator"
)

type InvocationHandlersFile struct {
	*generator.BaseGenerator
	Meta  *generator.Meta
	PbDir string
}

func NewInvocationHandlersFile(meta *generator.Meta) (generator.Generator, error) {
	baseGenerator, err := generator.NewBaseGenerator(g.GetDir(meta.RootDir, g.Service), ServiceFilename, false)
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
	return []func(){}
}
