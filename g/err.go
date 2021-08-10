package g

import (
	"github.com/pkg/errors"
)

var (
	ErrServiceNotFound            = errors.New("service not found")
	ErrProtocNotFound             = errors.New("protoc not found")
	ErrInvalidGeneratorParameters = errors.New("invalid base generator parameters")
	ErrInterfaceNotFound          = errors.New("interface not found")
	ErrMethodNotFound             = errors.New("method not found")
	ErrVarNotFound                = errors.New("var not found")
	ErrConstNotFound              = errors.New("const not found")
)
