package generator

import "github.com/dave/jennifer/jen"

// CodeBuilder wraps a jen statement
type CodeBuilder struct {
	raw *jen.Statement
}

// NewCodeBuilder returns a partial generator
func NewCodeBuilder(st *jen.Statement) *CodeBuilder {
	if st != nil {
		return &CodeBuilder{
			raw: st,
		}
	}
	return &CodeBuilder{
		raw: &jen.Statement{},
	}
}

func (cg *CodeBuilder) AppendMultilineComment(c []string) {
	for i, v := range c {
		if i != len(c)-1 {
			cg.raw.Comment(v).Line()
			continue
		}
		cg.raw.Comment(v)
	}
}

// Raw returns the jen statement.
func (cg *CodeBuilder) Raw() *jen.Statement {
	return cg.raw
}

// String returns the source code string
func (cg *CodeBuilder) String() string {
	return cg.raw.GoString()
}

func (cg *CodeBuilder) AppendInterface(name string, methods []jen.Code) {
	cg.raw.Type().Id(name).Interface(methods...).Line()
}

func (cg *CodeBuilder) AppendStruct(name string, fields ...jen.Code) {
	cg.raw.Type().Id(name).Struct(fields...).Line()
}

// NewLine insert a new line in code.
func (cg *CodeBuilder) NewLine() {
	cg.raw.Line()
}

// AppendFunction new a function
// name is the function name
// receiver is (* T), if no receiver, use nil
// parameters is function parameters inside func()
// if specified oneResponse, then function return type is `oneResponse`
// if oneResponse is empty, it will use results to set function return stuff
func (cg *CodeBuilder) AppendFunction(name string, receiver *jen.Statement,
	parameters []jen.Code, results []jen.Code, oneResponse string, body ...jen.Code) {
	cg.raw.Func()
	if receiver != nil {
		cg.raw.Params(receiver)
	}
	if name != "" {
		cg.raw.Id(name)
	}
	cg.raw.Params(parameters...)

	// if we passed a returnType, then we only return one type,
	// in other words, returns should be nil in this situation
	if oneResponse != "" {
		cg.raw.Id(oneResponse)
	} else if len(results) > 0 {
		cg.raw.Params(results...)
	}
	cg.raw.Block(body...)
}
