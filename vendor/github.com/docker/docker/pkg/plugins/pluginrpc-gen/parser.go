package main

import (
	"errors"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"path"
	"reflect"
	"strings"
)

var errBadReturn = errors.New("found return arg with no name: all args must be named")

type errUnexpectedType struct ***REMOVED***
	expected string
	actual   interface***REMOVED******REMOVED***
***REMOVED***

func (e errUnexpectedType) Error() string ***REMOVED***
	return fmt.Sprintf("got wrong type expecting %s, got: %v", e.expected, reflect.TypeOf(e.actual))
***REMOVED***

// ParsedPkg holds information about a package that has been parsed,
// its name and the list of functions.
type ParsedPkg struct ***REMOVED***
	Name      string
	Functions []function
	Imports   []importSpec
***REMOVED***

type function struct ***REMOVED***
	Name    string
	Args    []arg
	Returns []arg
	Doc     string
***REMOVED***

type arg struct ***REMOVED***
	Name            string
	ArgType         string
	PackageSelector string
***REMOVED***

func (a *arg) String() string ***REMOVED***
	return a.Name + " " + a.ArgType
***REMOVED***

type importSpec struct ***REMOVED***
	Name string
	Path string
***REMOVED***

func (s *importSpec) String() string ***REMOVED***
	var ss string
	if len(s.Name) != 0 ***REMOVED***
		ss += s.Name
	***REMOVED***
	ss += s.Path
	return ss
***REMOVED***

// Parse parses the given file for an interface definition with the given name.
func Parse(filePath string, objName string) (*ParsedPkg, error) ***REMOVED***
	fs := token.NewFileSet()
	pkg, err := parser.ParseFile(fs, filePath, nil, parser.AllErrors)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	p := &ParsedPkg***REMOVED******REMOVED***
	p.Name = pkg.Name.Name
	obj, exists := pkg.Scope.Objects[objName]
	if !exists ***REMOVED***
		return nil, fmt.Errorf("could not find object %s in %s", objName, filePath)
	***REMOVED***
	if obj.Kind != ast.Typ ***REMOVED***
		return nil, fmt.Errorf("exected type, got %s", obj.Kind)
	***REMOVED***
	spec, ok := obj.Decl.(*ast.TypeSpec)
	if !ok ***REMOVED***
		return nil, errUnexpectedType***REMOVED***"*ast.TypeSpec", obj.Decl***REMOVED***
	***REMOVED***
	iface, ok := spec.Type.(*ast.InterfaceType)
	if !ok ***REMOVED***
		return nil, errUnexpectedType***REMOVED***"*ast.InterfaceType", spec.Type***REMOVED***
	***REMOVED***

	p.Functions, err = parseInterface(iface)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	// figure out what imports will be needed
	imports := make(map[string]importSpec)
	for _, f := range p.Functions ***REMOVED***
		args := append(f.Args, f.Returns...)
		for _, arg := range args ***REMOVED***
			if len(arg.PackageSelector) == 0 ***REMOVED***
				continue
			***REMOVED***

			for _, i := range pkg.Imports ***REMOVED***
				if i.Name != nil ***REMOVED***
					if i.Name.Name != arg.PackageSelector ***REMOVED***
						continue
					***REMOVED***
					imports[i.Path.Value] = importSpec***REMOVED***Name: arg.PackageSelector, Path: i.Path.Value***REMOVED***
					break
				***REMOVED***

				_, name := path.Split(i.Path.Value)
				splitName := strings.Split(name, "-")
				if len(splitName) > 1 ***REMOVED***
					name = splitName[len(splitName)-1]
				***REMOVED***
				// import paths have quotes already added in, so need to remove them for name comparison
				name = strings.TrimPrefix(name, `"`)
				name = strings.TrimSuffix(name, `"`)
				if name == arg.PackageSelector ***REMOVED***
					imports[i.Path.Value] = importSpec***REMOVED***Path: i.Path.Value***REMOVED***
					break
				***REMOVED***
			***REMOVED***
		***REMOVED***
	***REMOVED***

	for _, spec := range imports ***REMOVED***
		p.Imports = append(p.Imports, spec)
	***REMOVED***

	return p, nil
***REMOVED***

func parseInterface(iface *ast.InterfaceType) ([]function, error) ***REMOVED***
	var functions []function
	for _, field := range iface.Methods.List ***REMOVED***
		switch f := field.Type.(type) ***REMOVED***
		case *ast.FuncType:
			method, err := parseFunc(field)
			if err != nil ***REMOVED***
				return nil, err
			***REMOVED***
			if method == nil ***REMOVED***
				continue
			***REMOVED***
			functions = append(functions, *method)
		case *ast.Ident:
			spec, ok := f.Obj.Decl.(*ast.TypeSpec)
			if !ok ***REMOVED***
				return nil, errUnexpectedType***REMOVED***"*ast.TypeSpec", f.Obj.Decl***REMOVED***
			***REMOVED***
			iface, ok := spec.Type.(*ast.InterfaceType)
			if !ok ***REMOVED***
				return nil, errUnexpectedType***REMOVED***"*ast.TypeSpec", spec.Type***REMOVED***
			***REMOVED***
			funcs, err := parseInterface(iface)
			if err != nil ***REMOVED***
				fmt.Println(err)
				continue
			***REMOVED***
			functions = append(functions, funcs...)
		default:
			return nil, errUnexpectedType***REMOVED***"*astFuncType or *ast.Ident", f***REMOVED***
		***REMOVED***
	***REMOVED***
	return functions, nil
***REMOVED***

func parseFunc(field *ast.Field) (*function, error) ***REMOVED***
	f := field.Type.(*ast.FuncType)
	method := &function***REMOVED***Name: field.Names[0].Name***REMOVED***
	if _, exists := skipFuncs[method.Name]; exists ***REMOVED***
		fmt.Println("skipping:", method.Name)
		return nil, nil
	***REMOVED***
	if f.Params != nil ***REMOVED***
		args, err := parseArgs(f.Params.List)
		if err != nil ***REMOVED***
			return nil, err
		***REMOVED***
		method.Args = args
	***REMOVED***
	if f.Results != nil ***REMOVED***
		returns, err := parseArgs(f.Results.List)
		if err != nil ***REMOVED***
			return nil, fmt.Errorf("error parsing function returns for %q: %v", method.Name, err)
		***REMOVED***
		method.Returns = returns
	***REMOVED***
	return method, nil
***REMOVED***

func parseArgs(fields []*ast.Field) ([]arg, error) ***REMOVED***
	var args []arg
	for _, f := range fields ***REMOVED***
		if len(f.Names) == 0 ***REMOVED***
			return nil, errBadReturn
		***REMOVED***
		for _, name := range f.Names ***REMOVED***
			p, err := parseExpr(f.Type)
			if err != nil ***REMOVED***
				return nil, err
			***REMOVED***
			args = append(args, arg***REMOVED***name.Name, p.value, p.pkg***REMOVED***)
		***REMOVED***
	***REMOVED***
	return args, nil
***REMOVED***

type parsedExpr struct ***REMOVED***
	value string
	pkg   string
***REMOVED***

func parseExpr(e ast.Expr) (parsedExpr, error) ***REMOVED***
	var parsed parsedExpr
	switch i := e.(type) ***REMOVED***
	case *ast.Ident:
		parsed.value += i.Name
	case *ast.StarExpr:
		p, err := parseExpr(i.X)
		if err != nil ***REMOVED***
			return parsed, err
		***REMOVED***
		parsed.value += "*"
		parsed.value += p.value
		parsed.pkg = p.pkg
	case *ast.SelectorExpr:
		p, err := parseExpr(i.X)
		if err != nil ***REMOVED***
			return parsed, err
		***REMOVED***
		parsed.pkg = p.value
		parsed.value += p.value + "."
		parsed.value += i.Sel.Name
	case *ast.MapType:
		parsed.value += "map["
		p, err := parseExpr(i.Key)
		if err != nil ***REMOVED***
			return parsed, err
		***REMOVED***
		parsed.value += p.value
		parsed.value += "]"
		p, err = parseExpr(i.Value)
		if err != nil ***REMOVED***
			return parsed, err
		***REMOVED***
		parsed.value += p.value
		parsed.pkg = p.pkg
	case *ast.ArrayType:
		parsed.value += "[]"
		p, err := parseExpr(i.Elt)
		if err != nil ***REMOVED***
			return parsed, err
		***REMOVED***
		parsed.value += p.value
		parsed.pkg = p.pkg
	default:
		return parsed, errUnexpectedType***REMOVED***"*ast.Ident or *ast.StarExpr", i***REMOVED***
	***REMOVED***
	return parsed, nil
***REMOVED***
