package g

type ImportType int

const (
	_ ImportType = iota
	HdSdk
	HdSdkTypes
	HdUtils
	HdParallel
	KitEndpoint
	KitGrpc
	Cobra
	Errors
	StdGrpc
)

var (
	ImportPaths = map[ImportType]string{
		HdSdk:       "github.com/hdget/hdsdk",
		HdSdkTypes:  "github.com/hdget/hdsdk/types",
		HdUtils:     "github.com/hdget/hdsdk/utils",
		HdParallel:  "github.com/hdget/hdsdk/utils/parallel",
		KitEndpoint: "github.com/go-kit/kit/endpoint",
		KitGrpc:     "github.com/go-kit/kit/transport/grpc",
		Cobra:       "github.com/spf13/cobra",
		Errors:      "github.com/pkg/errors",
		StdGrpc:     "google.golang.org/grpc",
	}
)
