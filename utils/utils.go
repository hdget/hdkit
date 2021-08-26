package utils

import (
	"bytes"
	"embed"
	"github.com/dave/jennifer/jen"
	"github.com/hdget/hdkit/g"
	"golang.org/x/tools/imports"
	iofs "io/fs"
	"os"
	"path"
	"regexp"
	"runtime"
	"strings"
	"unicode"
)

const (
	Cameling string = `[\p{L}\p{N}]+`
)

var (
	rxCameling = regexp.MustCompile(Cameling)
)

// ToCamelCase converts from underscore separated form to camel case form.
func ToCamelCase(str string) string {
	byteSrc := []byte(str)
	chunks := rxCameling.FindAll(byteSrc, -1)
	for idx, val := range chunks {
		chunks[idx] = bytes.Title(val)
	}
	return string(bytes.Join(chunks, nil))
}

// ToLowerFirstCamelCase returns the given string in camelcase formatted string
// but with the first letter being lowercase.
func ToLowerFirstCamelCase(s string) string {
	if s == "" {
		return s
	}
	if len(s) == 1 {
		return strings.ToLower(string(s[0]))
	}
	return strings.ToLower(string(s[0])) + ToCamelCase(s)[1:]
}

// ToUpperFirst returns the given string with the first letter being uppercase.
func ToUpperFirst(s string) string {
	if s == "" {
		return s
	}
	if len(s) == 1 {
		return strings.ToLower(string(s[0]))
	}
	return strings.ToUpper(string(s[0])) + s[1:]
}

// ToLowerSnakeCase the given string in snake-case format.
func ToLowerSnakeCase(s string) string {
	return strings.ToLower(ToSnakeCase(s))
}

// ToSnakeCase converts from camel case form to underscore separated form.
func ToSnakeCase(s string) string {
	s = ToCamelCase(s)
	runes := []rune(s)
	length := len(runes)
	var out []rune
	for i := 0; i < length; i++ {
		out = append(out, unicode.ToLower(runes[i]))
		if i+1 < length && (unicode.IsUpper(runes[i+1]) && unicode.IsLower(runes[i])) {
			out = append(out, '_')
		}
	}

	return string(out)
}

// GoImportsSource is used to format and optimize imports the
// given source.
func GoImportsSource(path string, s string) ([]byte, error) {
	return imports.Process(path, []byte(s), nil)
}

// GetCmdServiceImportPath returns the import path of the cmd service (used by cmd/main.go).
func GetCmdServiceImportPath(name string) (string, error) {
	return GetImportPath(name, "cmd_service_path_format")
}

// GetEndpointImportPath returns the import path of the service endpoints.
func GetEndpointImportPath(name string) (string, error) {
	return GetImportPath(name, "endpoint_path_format")
}

// GetGRPCTransportImportPath returns the import path of the service grpc transport.
func GetGRPCTransportImportPath(name string) (string, error) {
	return GetImportPath(name, "grpc_path_format")
}

// GetPbImportPath returns the import path of the generated service grpc pb.
func GetPbImportPath(name, pathByFlag string) (string, error) {
	if pathByFlag != "" {
		return pathByFlag, nil
	}
	return GetImportPath(name, "grpc_pb_path_format")
}

// GetHTTPTransportImportPath returns the import path of the service http transport.
func GetHTTPTransportImportPath(name string) (string, error) {
	return GetImportPath(name, "http_path_format")
}

// GetGoPath returns the gopath.
func GetGoPath() string {
	env := "HOME"
	if runtime.GOOS == "windows" {
		env = "USERPROFILE"
	} else if runtime.GOOS == "plan9" {
		env = "home"
	}
	if home := os.Getenv(env); home != "" {
		def := path.Join(home, "go")
		if path.Clean(def) == path.Clean(runtime.GOROOT()) {
			// Don't set the default GOPATH to GOROOT,
			// as that will trigger warnings from the go tool.
			return ""
		}
		return def
	}
	return ""
}

func GetImportPath(rootDir string, svcPath string) (string, error) {
	modName, err := GetModNameFromModFile(rootDir)
	if err != nil {
		return "", err
	}

	gosrc := GetGoPath() + "/src/"
	gosrc = strings.Replace(gosrc, "\\", "/", -1)
	pwd, err := os.Getwd()
	if err != nil {
		return "", err
	}

	pwd = strings.Replace(pwd, "\\", "/", -1)
	projectPath := strings.Replace(pwd, gosrc, "", 1)

	path := strings.Replace(svcPath, "\\", "/", -1)
	if modName != "" {
		modName = strings.Replace(modName, "\\", "/", -1)
		modNameArr := strings.Split(modName, "/")
		if len(modNameArr) <= 1 {
			projectPath = ""
		} else {
			projectPath = strings.Join(modNameArr[0:len(modNameArr)-1], "/")
		}
	}
	var importPath string
	// Change: here should not use os.Getwd() as projectPath
	// Desc:It can't pass go test, on windows, projectPath will be "c:/User/xxx/...", this will cause err certainly.
	// projectPath = ""
	if projectPath == "" {
		importPath = path
	} else {
		importPath = projectPath + "/" + path
	}
	return importPath, nil
}

func GetModNameFromModFile(name string) (string, error) {
	modFile := "go.mod"
	filePath := ToLowerSnakeCase(name) + "/" + modFile
	exists, _ := g.GetFs().Exists(filePath)
	var modFileInParentLevel bool
	if !exists {
		//if the service level has no go.mod file, it will check the parent level
		exists, err := g.GetFs().Exists(modFile)
		if !exists {
			return "", err
		}
		filePath = modFile
		modFileInParentLevel = true
	}

	data, err := g.GetFs().ReadFile(filePath)
	if err != nil {
		return "", err
	}

	modDataArr := strings.Split(string(data), "\n")
	if len(modDataArr) != 0 {
		modNameArr := strings.Split(modDataArr[0], " ")
		if len(modNameArr) < 2 { // go.mod file: module XXXX/XXXX/{projectName}
			return "", nil
		}
		if modFileInParentLevel {
			return modNameArr[1] + "/" + name, nil
		}
		return modNameArr[1], nil
	}
	return "", nil
}

func IsExist(path string) bool {
	_, err := os.Stat(path)
	return err == nil || os.IsExist(err)
}

// StringSliceContains 检查字符串slice中是否含有
func StringSliceContains(list []string, checkItem string) bool {
	if len(list) == 0 {
		return false
	}

	for _, item := range list {
		if item == checkItem {
			return true
		}
	}

	return false
}

// IsUpperCase check if the string is uppercase. Empty string is valid.
func IsUpperCase(str string) bool {
	if len(str) == 0 {
		return true
	}
	return str == strings.ToUpper(str)
}

func GetValidParameterCode(paramName, impPath, paramType string) jen.Code {
	// check if target is a point or not
	isParamTypeAPoint := false
	if paramType[0] == '*' {
		paramType = paramType[1:]
		isParamTypeAPoint = true
	}

	var ret jen.Code
	if paramName != "" {
		ret = jen.Id(paramName).Qual(impPath, paramType)
		if isParamTypeAPoint {
			ret = jen.Id(paramName).Op("*").Qual(impPath, paramType)
		}
	} else {
		ret = jen.Qual(impPath, paramType)
		if isParamTypeAPoint {
			ret = jen.Op("*").Qual(impPath, paramType)
		}
	}

	return ret
}

// TraverseDirFiles return relative dirs
func TraverseDirFiles(fs embed.FS, subdir string) ([]string, []string, error) {
	dirs := make([]string, 0)
	files := make([]string, 0)

	err := iofs.WalkDir(fs, subdir, func(path string, d iofs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if d.IsDir() {
			dirs = append(dirs, path)
		} else {
			files = append(files, path)
		}
		return nil
	})
	if err != nil {
		return nil, nil, err
	}

	return dirs, files, nil
}
