package cmdgen

import "github.com/hdget/hdkit/generator"

type FileFactory interface {
	Create() error // generate something based on already created one
}

type NewFileFunc = func(*generator.Meta) (generator.Generator, error)

