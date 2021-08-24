package generator

import (
	"github.com/hdget/hdkit/g"
	"github.com/hdget/hdkit/parser"
	"github.com/hdget/hdkit/utils"
	"strings"
)

type Meta struct {
	RawSvcName             string            // service defined in .proto file
	SvcStructName          string            // service struct
	SvcServerInterfaceName string            // service interface name
	SvcServerInterface     *parser.Interface // service server interface found in .pb.go files
	RootDir                string
	ImportPaths            map[g.ImportType]string
}

func NewMeta(rootDir string) (*Meta, error) {
	svcInterface, err := parser.ParsePbFiles(g.GetDir(rootDir, g.Pb))
	if err != nil {
		return nil, err
	}

	// Create service struct name for pb's `ServiceSever` interface
	rawSvcName := svcInterface.Name
	beginPosition := strings.LastIndex(rawSvcName, "Server")
	if beginPosition > 0 {
		rawSvcName = svcInterface.Name[:beginPosition]
	}

	// get service struct name
	// if `service` defined in protobuf file has suffix `XxxService`, do nothing
	// if `service` defined in protobuf file doesn't has suffix `Service`, append it
	var svcStructName string
	svcStructSuffix := "Impl"
	if !strings.HasSuffix(strings.ToLower(rawSvcName), "service") {
		svcStructSuffix = "Service" + svcStructSuffix
	}
	svcStructName = utils.ToLowerFirstCamelCase(rawSvcName) + svcStructSuffix

	return &Meta{
		RawSvcName:             rawSvcName,
		SvcServerInterfaceName: svcInterface.Name,
		SvcStructName:          svcStructName,
		SvcServerInterface:     svcInterface,
		RootDir:                rootDir,
	}, nil
}
