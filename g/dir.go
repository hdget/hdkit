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
)

var (
	dirTokens = map[DirType][]string{
		Binary:  []string{"bin"},
		Proto:   []string{"proto"},
		Global:  []string{"g"}, // global directory which stores config, error
		Service: []string{"service"},
		Pb:      []string{"autogen", "pb"},
		Grpc:    []string{"autogen", "grpc"},
		Http:    []string{"autogen", "http"},
		Cmd:     []string{"cmd"},
	}
)

// GetDir get default project dir
func GetDir(rootDir string, dirType DirType) string {
	return buildDir(rootDir, dirTokens[dirType])
}
//
//// GetProjectDirs get default project dir
//func GetProjectDirs(rootDir string) map[DirType]string {
//	if len(gDirs) == 0 {
//		gDirs = map[DirType]string{
//			Binary:  buildDir(rootDir, dirTokens[Binary]),
//			Proto:   buildDir(rootDir, dirTokens[Proto]),
//			Global:  buildDir(rootDir, dirTokens[Global]),
//			Service: buildDir(rootDir, dirTokens[Service]),
//			Pb:      buildDir(rootDir, dirTokens[Pb]),
//			Grpc:    buildDir(rootDir, dirTokens[Grpc]),
//			Http:    buildDir(rootDir, dirTokens[Http]),
//			Cmd:     buildDir(rootDir, dirTokens[Cmd]),
//		}
//	}
//	return gDirs
//}

// assembly the relative Dir for generating files
func buildDir(rootDir string, dirTokens []string, args ...string) string {
	dirParts := append([]string{rootDir}, dirTokens...)
	dirParts = append(dirParts, args...)
	return path.Join(dirParts...)
}
