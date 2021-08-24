package gokit

import (
	"github.com/hdget/hdkit/generator"
	"github.com/hdget/hdkit/pkg/cmdgen"
	"github.com/hdget/hdkit/pkg/cmdgen/gokit/autogen/grpc"
	"github.com/hdget/hdkit/pkg/cmdgen/gokit/autogen/http"
	"github.com/hdget/hdkit/pkg/cmdgen/gokit/cmd"
	"github.com/hdget/hdkit/pkg/cmdgen/gokit/g"
	"github.com/hdget/hdkit/pkg/cmdgen/gokit/service"
)

type gokitFileFactory struct {
	rootDir string
	meta *generator.Meta
}


// NewGokitCodeFactory returns a initialized and ready generator.
func NewGokitFileFactory(rootDir string) (cmdgen.FileFactory, error) {
	meta, err := generator.NewMeta(rootDir)
	if err != nil {
		return nil, err
	}

	return &gokitFileFactory{
		rootDir: rootDir,
		meta: meta,
	}, nil
}

// Create create files
func (factory *gokitFileFactory) Create() error {
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

	// generate method based files
	// e,g: endpoint_<method>.go
	for _, method := range factory.meta.SvcServerInterface.Methods {
		g, err := grpc.NewGrpcAspectMethodFile(method, factory.meta)
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

	for _, method := range factory.meta.SvcServerInterface.Methods {
		g, err := http.NewHttpAspectMethodFile(method, factory.meta)
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

func (factory *gokitFileFactory) getNewFileFuncs() []cmdgen.NewFileFunc {
	return []cmdgen.NewFileFunc{
		service.NewServiceFile,      // service/service.go
		grpc.NewGrpcHandlersFile,    // autogen/grpc/handlers.go
		grpc.NewClientFile,          // autogen/grpc/client.go
		http.NewHttpHandlersFile,    // autogen/http/handlers.go
		g.NewGConfigFile,            // g/config.go
		cmd.NewCmdRootFile,          // cmd/root.go
		cmd.NewCmdRunFile,           // cmd/run.go
		cmd.NewCmdRunGrpcServerFile, // cmd/run_grpc.go
		cmd.NewCmdRunHttpServerFile, // cmd/run_http.go
		cmd.NewCmdRunClientFile,     // cmd/client.go
		NewMainFile,            // main.go
	}
}

