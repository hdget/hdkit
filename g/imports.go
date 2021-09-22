package g

type ImportType int

const (
	_ ImportType = iota
	HdSdk
	HdSdkTypes
	HdUtils
	HdWs
	HdParallel
	KitEndpoint
	KitGrpc
	KitHttp
	DaprGrpc
	DaprHttp
	DaprCommon
	Cobra
	Errors
	StdGrpc
	StdHttp
	StdJson
)

var (
	ImportPaths = map[ImportType]string{
		HdSdk:       "github.com/hdget/hdsdk",
		HdSdkTypes:  "github.com/hdget/hdsdk/types",
		HdUtils:     "github.com/hdget/hdsdk/utils",
		HdWs:        "github.com/hdget/hdsdk/lib/ws",
		HdParallel:  "github.com/hdget/hdsdk/utils/parallel",
		KitEndpoint: "github.com/go-kit/kit/endpoint",
		KitGrpc:     "github.com/go-kit/kit/transport/grpc",
		KitHttp:     "github.com/go-kit/kit/transport/http",
		Cobra:       "github.com/spf13/cobra",
		Errors:      "github.com/pkg/errors",
		StdGrpc:     "google.golang.org/grpc",
		StdHttp:     "net/http",
		StdJson:     "encoding/json",
		DaprGrpc:    "github.com/dapr/go-sdk/service/grpc",
		DaprHttp:    "github.com/dapr/go-sdk/service/http",
		DaprCommon:  "github.com/dapr/go-sdk/service/common",
	}
)
