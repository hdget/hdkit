// Package generator stores global variables
package g

import (
	"path"
)

type DirType int

const (
	_ DirType = iota
	Binary
	Proto
	Global
	Service
	Pb
	Grpc
	Http
	Cmd
	Setting
)

var (
	dirTokens = map[DirType][]string{
		Binary:  []string{"bin"},
		Proto:   []string{"proto"},
		Global:  []string{"g"}, // global directory which stores config, error
		Service: []string{"pkg", "service"},
		Pb:      []string{"autogen", "pb"},
		Grpc:    []string{"autogen", "grpc"},
		Http:    []string{"autogen", "http"},
		Cmd:     []string{"cmd"},
		Setting: []string{"setting"},
	}
)

// GetDir get default project dir
func GetDir(rootDir string, dirType DirType) string {
	return buildDir(rootDir, dirTokens[dirType])
}

// assembly the relative Dir for generating files
func buildDir(rootDir string, dirTokens []string, args ...string) string {
	dirParts := append([]string{rootDir}, dirTokens...)
	dirParts = append(dirParts, args...)
	return path.Join(dirParts...)
}
