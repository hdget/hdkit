package generator

import (
	"bytes"
	"fmt"
	"github.com/dave/jennifer/jen"
	"github.com/hdget/hdkit/g"
	"github.com/hdget/hdkit/parser"
	"github.com/hdget/hdkit/utils"
	"github.com/pkg/errors"
	"go/ast"
	"go/format"
	ps "go/parser"
	"go/token"
	"path"
	"strconv"
	"strings"
)

type BaseGenerator struct {
	pkg       string // package name
	Dir       string // the directory save the output file
	Filename  string // Filename for output file
	overwrite bool   // if the file need to be overwrite or appended

	JenFile     *jen.File     // original go file
	ParsedFile  *parser.File  // go file -> parsed file
	fileContent *bytes.Buffer // source file content

	Builder *CodeBuilder

	GenFuncs []func() // sub generate functions

	isNewCreated bool // flag to indicate the file is new create or not
}

// NewBaseGenerator parse file if file doesn't exist, read FileContent from the file
// the first variadic is package name if specified
func NewBaseGenerator(dir, filename string, overwrite bool, args ...string) (*BaseGenerator, error) {
	var pkg string
	if len(args) > 0 {
		pkg = args[0]
	}
	// if specified package name then use NewFile to create jen.JenFile
	// or it will use the basename of the Dir as package name
	var jenFile *jen.File
	if pkg != "" {
		jenFile = jen.NewFile(pkg)
	} else {
		jenFile = jen.NewFilePath(dir)
	}

	// check if the Dir/Filename exists or not,
	// if not, create it
	filepath := path.Join(dir, filename)
	exists, err := g.GetFs().Exists(filepath)
	if err != nil {
		return nil, err
	}

	// if no file exists, it will create automatically
	isNewCreated := false
	if !exists {
		if err := g.GetFs().WriteFile(filepath, []byte(jenFile.GoString()), false); err != nil {
			return nil, err
		}
		isNewCreated = true
	}

	fileData, err := g.GetFs().ReadFile(filepath)
	if err != nil {
		return nil, err
	}

	var buffer = &bytes.Buffer{}
	if _, err = buffer.Write(fileData); err != nil {
		return nil, err
	}

	parsedFile, err := parser.NewFileParser().Parse(buffer.Bytes())
	if err != nil {
		return nil, err
	}

	return &BaseGenerator{
		pkg:          pkg,
		Dir:          dir,
		Filename:     filename,
		overwrite:    overwrite,
		JenFile:      jenFile,
		ParsedFile:   parsedFile,
		fileContent:  buffer,
		Builder:      NewCodeBuilder(jenFile.Empty()),
		isNewCreated: isNewCreated,
	}, nil
}

func (bg *BaseGenerator) PreGenerate() error {
	if bg.Dir == "" || bg.Filename == "" || bg.JenFile == nil ||
		bg.ParsedFile == nil || bg.Builder == nil {
		return g.ErrInvalidGeneratorParameters
	}

	err := g.GetFs().MakeDir(bg.Dir)
	if err != nil {
		return err
	}

	return nil
}

func (bg *BaseGenerator) PostGenerate() error {
	err := bg.save()
	if err != nil {
		return err
	}
	return nil
}

func (bg *BaseGenerator) Generate(concrete Generator) error {
	fmt.Println("Generating: ", path.Join(bg.Dir, bg.Filename))

	if err := concrete.PreGenerate(); err != nil {
		return err
	}

	for _, genFunc := range concrete.GetGenCodeFuncs() {
		genFunc()
	}

	if err := concrete.PostGenerate(); err != nil {
		return err
	}

	return nil
}

func (bg *BaseGenerator) FindInterface(interfaceName string) (*parser.Interface, error) {
	var found *parser.Interface
	for i, v := range bg.ParsedFile.Interfaces {
		if v.Name == interfaceName {
			found = &bg.ParsedFile.Interfaces[i]
			break
		}
	}
	if found == nil {
		return nil, errors.Wrap(g.ErrInterfaceNotFound, interfaceName)
	}
	return found, nil
}

func (bg *BaseGenerator) FindVar(varName string) (*parser.NamedTypeValue, error) {
	var found *parser.NamedTypeValue
	for i, v := range bg.ParsedFile.Vars {
		if (v.Type == "" && v.Name == varName) || (v.Type != "" && v.Type == varName) {
			found = &bg.ParsedFile.Vars[i]
			break
		}
	}
	if found == nil {
		return nil, errors.Wrap(g.ErrVarNotFound, varName)
	}
	return found, nil
}

func (bg *BaseGenerator) FindConst(constName string) (*parser.NamedTypeValue, error) {
	var found *parser.NamedTypeValue
	for i, v := range bg.ParsedFile.Constants {
		if (v.Type == "" && v.Name == constName) || (v.Type != "" && v.Type == constName) {
			found = &bg.ParsedFile.Constants[i]
			break
		}
	}
	if found == nil {
		return nil, errors.Wrap(g.ErrConstNotFound, constName)
	}
	return found, nil
}

func (bg *BaseGenerator) FindMethod(methodName string) (*parser.Method, error) {
	var found *parser.Method
	for i, v := range bg.ParsedFile.Methods {
		if v.Name == methodName {
			found = &bg.ParsedFile.Methods[i]
			break
		}
	}
	if found == nil {
		return nil, errors.Wrap(g.ErrMethodNotFound, methodName)
	}
	return found, nil
}

func (bg *BaseGenerator) FindStructure(structureName string) (*parser.Struct, error) {
	var found *parser.Struct
	for i, v := range bg.ParsedFile.Structures {
		if v.Name == structureName {
			found = &bg.ParsedFile.Structures[i]
			break
		}
	}
	if found == nil {
		return nil, errors.Wrap(g.ErrStructureNotFound, structureName)
	}
	return found, nil
}

// GenerateNameBySample is used to generate a variable name using a sample.
//
// The exclude parameter represents the names that it can not use.
//
// E.x  sample = "hello" this will return the name "h" if it is not in any NamedTypeValue name.
func (bg *BaseGenerator) GenerateNameBySample(sample string, exclude []parser.NamedTypeValue) string {
	sn := 1
	name := utils.ToLowerFirstCamelCase(sample)[:sn]
	for _, v := range exclude {
		if v.Name == name {
			sn++
			if sn > len(sample) {
				sample = sample[len(sample)-sn:]
			}
			name = utils.ToLowerFirstCamelCase(sample)[:sn]
		}
	}
	return name
}

// EnsureThatWeUseQualifierIfNeeded is used to see if we need to import a path of a given type.
func (bg *BaseGenerator) EnsureThatWeUseQualifierIfNeeded(tp string, imp []parser.NamedTypeValue) string {
	if bytes.HasPrefix([]byte(tp), []byte("...")) {
		return ""
	}
	if t := strings.Split(tp, "."); len(t) > 0 {
		s := t[0]
		for _, v := range imp {
			i, _ := strconv.Unquote(v.Type)
			if strings.HasSuffix(i, s) || v.Name == s {
				return i
			}
		}
		return ""
	}
	return ""
}

// AddImportsToFile adds missing imports to a file that we edit with the generator
func (bg *BaseGenerator) AddImportsToFile(imp []parser.NamedTypeValue, src string) (string, error) {
	// Create the AST by parsing src
	fset := token.NewFileSet()
	f, err := ps.ParseFile(fset, "", src, 0)
	if err != nil {
		return "", err
	}
	found := false
	// Add the imports
	for i := 0; i < len(f.Decls); i++ {
		d := f.Decls[i]
		switch d := d.(type) {
		case *ast.FuncDecl:
			// No action
		case *ast.GenDecl:
			// IMPORT Declarations
			if d.Tok == token.IMPORT {
				if d.Rparen == 0 || d.Lparen == 0 {
					d.Rparen = f.Package
					d.Lparen = f.Package
				}
				found = true
				// Add the new import
				for _, v := range imp {
					iSpec := &ast.ImportSpec{
						Name: &ast.Ident{Name: v.Name},
						Path: &ast.BasicLit{Value: v.Type},
					}
					d.Specs = append(d.Specs, iSpec)
				}
			}
		}
	}

	if !found {
		dd := ast.GenDecl{
			TokPos: f.Package + 1,
			Tok:    token.IMPORT,
			Specs:  []ast.Spec{},
			Lparen: f.Package,
			Rparen: f.Package,
		}
		lastPos := 0
		for _, v := range imp {
			lastPos += len(v.Name) + len(v.Type) + 1
			iSpec := &ast.ImportSpec{
				Name:   &ast.Ident{Name: v.Name},
				Path:   &ast.BasicLit{Value: v.Type},
				EndPos: token.Pos(lastPos),
			}
			dd.Specs = append(dd.Specs, iSpec)

		}
		f.Decls = append([]ast.Decl{&dd}, f.Decls...)
	}

	// Sort the imports
	ast.SortImports(fset, f)
	var buf bytes.Buffer
	if err := format.Node(&buf, fset, f); err != nil {
		return "", err
	}
	return buf.String(), nil
}

// GuessType guess what's the valid parameter type in method
// only process following cases:
// - `map[]Xxx or map[]*Xxx`
// - `[]Xxx or []*Xxx`
// - `Xxx or *Xxx`
// Above upper cased struct/interface will regards as comes from service package
func (bg *BaseGenerator) GuessType(rawType string) string {
	validParamType := rawType

	pkgName := path.Base(bg.Dir)

	// check map[] or slice []
	separator := "]"
	arrayTokens := strings.Split(rawType, separator)
	// map[]Xxx or []xxx
	if len(arrayTokens) > 1 {
		checkToken := arrayTokens[len(arrayTokens)-1]

		var point string
		if checkToken[0] == '*' {
			checkToken = checkToken[1:]
			point = "*"
		}

		if utils.IsUpperCase(checkToken[:1]) {
			validType := pkgName + "." + checkToken
			validParamType = strings.Join(arrayTokens[:len(arrayTokens)-1], "") + separator + point + validType
			return validParamType
		}
	}

	// if rawType is not something like "x.y" style
	// which means rawType is a whole word
	dotTokens := strings.Split(rawType, ".")
	if len(dotTokens) == 1 {
		checkToken := rawType
		var point string
		if checkToken[0] == '*' {
			checkToken = checkToken[1:]
			point = "*"
		}

		// If the type of the parameter is not `something.MyType` and it starts with an uppercase
		// than the type was defined inside the service package.
		if utils.IsUpperCase(checkToken[:1]) && checkToken[0] != '[' {
			validParamType = point + pkgName + "." + checkToken
			return validParamType
		}
	}
	return validParamType
}

// TypeAssert try to type assert
// to, ok = source.(target)
// if !ok {
//    return errors.New(errMsg)
// }
func (bg *BaseGenerator) TypeAssert(to, source, tgtImportPath, target, errMsg string) []jen.Code {
	// check if target is a point or not
	isTargetAPoint := false
	if target[0] == '*' {
		target = target[1:]
		isTargetAPoint = true
	}

	// if target is a point, then add Op("*")
	typeAssert := jen.List(jen.Id(to), jen.Id("ok")).Op(":=").Id(source).Assert(jen.Qual(tgtImportPath, target))
	if isTargetAPoint {
		typeAssert = jen.List(jen.Id(to), jen.Id("ok")).Op(":=").Id(source).Assert(jen.Op("*").Qual(tgtImportPath, target))
	}
	return []jen.Code{
		typeAssert,
		jen.If(jen.Op("!").Id("ok")).Block(
			jen.Return(jen.Nil(), jen.Qual("github.com/pkg/errors", "New").Call(jen.Lit(errMsg))),
		),
	}
}

func (bg *BaseGenerator) Deference(v string) string {
	if len(v) == 0 {
		return ""
	}

	if v[0] == '*' {
		v = v[1:]
	}
	return v
}

// Save forced will overwrite the file
func (bg *BaseGenerator) save() error {
	data := bg.JenFile.GoString()
	if !bg.overwrite && !bg.isNewCreated && bg.fileContent.Len() > 0 {
		bg.fileContent.WriteString("\n")
		bg.fileContent.WriteString(bg.Builder.String())
		data = bg.fileContent.String()
	}

	toWrite, err := utils.GoImportsSource(bg.Dir, data)
	if err != nil {
		return err
	}

	err = g.GetFs().WriteFile(path.Join(bg.Dir, bg.Filename), toWrite, true)
	if err != nil {
		return err
	}

	return nil
}
