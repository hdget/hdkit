package parser

import (
	"fmt"
	"github.com/hdget/hdkit/g"
	"github.com/pkg/errors"
	"path"
	"path/filepath"
	"strings"
)

// ParsePbFiles Parse file if file doesn't exist, read FileContent from the file
// @return parsed service interface
func ParsePbFiles(pbDir string) (*Interface, error) {
	matches, err := filepath.Glob(path.Join(pbDir, "*.go"))
	if err != nil {
		return nil, err
	}

	for _, fp := range matches {
		svcInterface, _ := parsePbFile(fp)
		if svcInterface != nil {
			return svcInterface, nil
		}
	}

	return nil, errors.Wrapf(g.ErrServiceNotFound, "parse pb files in: %s", pbDir)
}

func parsePbFile(pbFilePath string) (*Interface, error) {
	fileContent, err := g.GetFs().ReadFile(pbFilePath)
	if err != nil {
		return nil, err
	}

	parsedFile, err := NewFileParser().Parse([]byte(fileContent))
	if err != nil {
		return nil, err
	}

	svcInterface, err := findServiceInterface(parsedFile)
	if err != nil {
		return nil, err
	}

	return svcInterface, nil
}

func findServiceInterface(pbFile *File) (*Interface, error) {
	var found *Interface
	for i, intf := range pbFile.Interfaces {
		if strings.HasSuffix(intf.Name, "Server") {
			found = &pbFile.Interfaces[i]
			break
		}
	}

	if found == nil {
		return nil, fmt.Errorf("no service defined in proto file")
	}

	return found, nil
}
