package main

import (
	"github.com/hdget/hdkit/cmd"
	"github.com/spf13/viper"
	"path"
	"runtime"
)

var (
	defaultSettings = map[string]string{
		// client related
		"client_cmd_path_format": path.Join("%s", "cmd", "client"),
		// cmd related
		"cmd_service_path_format": path.Join("%s", "cmd", "service"),
		"cmd_path_format":         path.Join("%s", "cmd"),
		"cmd_base_file_name":      "service_gen.go",
		"cmd_svc_file_name":       "service.go",
		// endpoint related
		"endpoint_path_format":          path.Join("%s", "pkg", "endpoint"),
		"endpoint_base_file_name":       "endpoint_gen.go",
		"endpoint_file_name":            "endpoint.go",
		"endpoint_middleware_file_name": "middleware.go",
		// grpc related
		"grpc_client_path_format": path.Join("%s", "client", "grpc"),
		"grpc_path_format":        path.Join("%s", "pkg", "grpc"),
		"grpc_pb_path_format":     path.Join("%s", "pkg", "grpc", "pb"),
		"grpc_client_file_name":   "grpc.go",
		"grpc_pb_file_name":       "%s.proto",
		"grpc_base_file_name":     "handler_gen.go",
		"grpc_file_name":          "handler.go",
		// http related
		"http_file_name":          "handler.go",
		"http_base_file_name":     "handler_gen.go",
		"http_client_file_name":   "http.go",
		"http_path_format":        path.Join("%s", "pkg", "http"),
		"http_client_path_format": path.Join("%s", "client", "http"),
		// service related
		"service_file_name":            "service.go",
		"service_middleware_file_name": "middleware.go",
		"service_path_format":          path.Join("%s", "pkg", "service"),
		"service_struct_prefix":        "basic",
	}
)

func main() {
	setDefaults()
	viper.AutomaticEnv()
	cmd.Execute()
}

func setDefaults() {
	for k, v := range defaultSettings {
		viper.SetDefault(k, v)
	}
	if runtime.GOOS == "windows" {
		viper.SetDefault("grpc_compile_file_name", "compile.bat")
	} else {
		viper.SetDefault("grpc_compile_file_name", "compile.sh")
	}
}
