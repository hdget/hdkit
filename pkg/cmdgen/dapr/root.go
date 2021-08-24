package dapr

import (
	"github.com/hdget/hdkit/generator"
	"github.com/hdget/hdkit/pkg/cmdgen"
	"github.com/hdget/hdkit/pkg/cmdgen/dapr/cmd"
	"github.com/hdget/hdkit/pkg/cmdgen/dapr/g"
	"github.com/hdget/hdkit/pkg/cmdgen/dapr/service"
)

type daprFileFactory struct {
	rootDir string
	meta *generator.Meta
}


// NewDaprFileFactory returns a initialized and ready generator.
func NewDaprFileFactory(rootDir string) (cmdgen.FileFactory, error) {
	meta, err := generator.NewMeta(rootDir)
	if err != nil {
		return nil, err
	}

	return &daprFileFactory{
		rootDir: rootDir,
		meta: meta,
	}, nil
}

// Create create files
func (factory *daprFileFactory) Create() error {
	// generate all individual files
	for _, newFunc := range factory.getNewFileFuncs() {
		g, err := newFunc(factory.meta)
		if err != nil {
			return err
		}

		if g != nil {
			err := g.Generate(g)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func (factory *daprFileFactory) getNewFileFuncs() []cmdgen.NewFileFunc {
	return []cmdgen.NewFileFunc{
		service.NewServiceFile,      // service/service.go
		g.NewGConfigFile,            // g/config.go
		cmd.NewCmdRootFile,          // cmd/root.go
		cmd.NewCmdRunFile,           // cmd/run.go
		cmd.NewCmdRunGrpcServerFile, // cmd/run_grpc.go
		cmd.NewCmdRunHttpServerFile, // cmd/run_http.go
		NewMainFile,            // main.go
	}
}

