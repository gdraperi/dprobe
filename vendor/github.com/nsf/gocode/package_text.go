package main

import (
	"bytes"
	"errors"
	"fmt"
	"go/ast"
	"go/token"
	"strconv"
	"text/scanner"
)

//-------------------------------------------------------------------------
// gc_parser
//
// The following part of the code may contain portions of the code from the Go
// standard library, which tells me to retain their copyright notice:
//
// Copyright (c) 2009 The Go Authors. All rights reserved.
//
// Redistribution and use in source and binary forms, with or without
// modification, are permitted provided that the following conditions are
// met:
//
//    * Redistributions of source code must retain the above copyright
// notice, this list of conditions and the following disclaimer.
//    * Redistributions in binary form must reproduce the above
// copyright notice, this list of conditions and the following disclaimer
// in the documentation and/or other materials provided with the
// distribution.
//    * Neither the name of Google Inc. nor the names of its
// contributors may be used to endorse or promote products derived from
// this software without specific prior written permission.
//
// THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS
// "AS IS" AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT
// LIMITED TO, THE IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR
// A PARTICULAR PURPOSE ARE DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT
// OWNER OR CONTRIBUTORS BE LIABLE FOR ANY DIRECT, INDIRECT, INCIDENTAL,
// SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES (INCLUDING, BUT NOT
// LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES; LOSS OF USE,
// DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND ON ANY
// THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT
// (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE
// OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.
//-------------------------------------------------------------------------

type gc_parser struct ***REMOVED***
	scanner      scanner.Scanner
	tok          rune
	lit          string
	path_to_name map[string]string
	beautify     bool
	pfc          *package_file_cache
***REMOVED***

func (p *gc_parser) init(data []byte, pfc *package_file_cache) ***REMOVED***
	p.scanner.Init(bytes.NewReader(data))
	p.scanner.Error = func(_ *scanner.Scanner, msg string) ***REMOVED*** p.error(msg) ***REMOVED***
	p.scanner.Mode = scanner.ScanIdents | scanner.ScanInts | scanner.ScanStrings |
		scanner.ScanComments | scanner.ScanChars | scanner.SkipComments
	p.scanner.Whitespace = 1<<'\t' | 1<<' ' | 1<<'\r' | 1<<'\v' | 1<<'\f'
	p.scanner.Filename = "package.go"
	p.next()
	// and the built-in "unsafe" package to the path_to_name map
	p.path_to_name = map[string]string***REMOVED***"unsafe": "unsafe"***REMOVED***
	p.pfc = pfc
***REMOVED***

func (p *gc_parser) next() ***REMOVED***
	p.tok = p.scanner.Scan()
	switch p.tok ***REMOVED***
	case scanner.Ident, scanner.Int, scanner.String:
		p.lit = p.scanner.TokenText()
	default:
		p.lit = ""
	***REMOVED***
***REMOVED***

func (p *gc_parser) error(msg string) ***REMOVED***
	panic(errors.New(msg))
***REMOVED***

func (p *gc_parser) errorf(format string, args ...interface***REMOVED******REMOVED***) ***REMOVED***
	p.error(fmt.Sprintf(format, args...))
***REMOVED***

func (p *gc_parser) expect(tok rune) string ***REMOVED***
	lit := p.lit
	if p.tok != tok ***REMOVED***
		p.errorf("expected %s, got %s (%q)", scanner.TokenString(tok),
			scanner.TokenString(p.tok), lit)
	***REMOVED***
	p.next()
	return lit
***REMOVED***

func (p *gc_parser) expect_keyword(keyword string) ***REMOVED***
	lit := p.expect(scanner.Ident)
	if lit != keyword ***REMOVED***
		p.errorf("expected keyword: %s, got: %q", keyword, lit)
	***REMOVED***
***REMOVED***

func (p *gc_parser) expect_special(what string) ***REMOVED***
	i := 0
	for i < len(what) ***REMOVED***
		if p.tok != rune(what[i]) ***REMOVED***
			break
		***REMOVED***

		nc := p.scanner.Peek()
		if i != len(what)-1 && nc <= ' ' ***REMOVED***
			break
		***REMOVED***

		p.next()
		i++
	***REMOVED***

	if i < len(what) ***REMOVED***
		p.errorf("expected: %q, got: %q", what, what[0:i])
	***REMOVED***
***REMOVED***

// dotIdentifier = "?" | ( ident | '·' ) ***REMOVED*** ident | int | '·' ***REMOVED*** .
// we're doing lexer job here, kind of
func (p *gc_parser) parse_dot_ident() string ***REMOVED***
	if p.tok == '?' ***REMOVED***
		p.next()
		return "?"
	***REMOVED***

	ident := ""
	sep := 'x'
	i, j := 0, -1
	for (p.tok == scanner.Ident || p.tok == scanner.Int || p.tok == '·') && sep > ' ' ***REMOVED***
		ident += p.lit
		if p.tok == '·' ***REMOVED***
			ident += "·"
			j = i
			i++
		***REMOVED***
		i += len(p.lit)
		sep = p.scanner.Peek()
		p.next()
	***REMOVED***
	// middot = \xc2\xb7
	if j != -1 && i > j+1 ***REMOVED***
		c := ident[j+2]
		if c >= '0' && c <= '9' ***REMOVED***
			ident = ident[0:j]
		***REMOVED***
	***REMOVED***
	return ident
***REMOVED***

// ImportPath = string_lit .
// quoted name of the path, but we return it as an identifier, taking an alias
// from 'pathToAlias' map, it is filled by import statements
func (p *gc_parser) parse_package() *ast.Ident ***REMOVED***
	path, err := strconv.Unquote(p.expect(scanner.String))
	if err != nil ***REMOVED***
		panic(err)
	***REMOVED***

	return ast.NewIdent(path)
***REMOVED***

// ExportedName = "@" ImportPath "." dotIdentifier .
func (p *gc_parser) parse_exported_name() *ast.SelectorExpr ***REMOVED***
	p.expect('@')
	pkg := p.parse_package()
	if pkg.Name == "" ***REMOVED***
		pkg.Name = "!" + p.pfc.name + "!" + p.pfc.defalias
	***REMOVED*** else ***REMOVED***
		pkg.Name = p.path_to_name[pkg.Name]
	***REMOVED***
	p.expect('.')
	name := ast.NewIdent(p.parse_dot_ident())
	return &ast.SelectorExpr***REMOVED***X: pkg, Sel: name***REMOVED***
***REMOVED***

// Name = identifier | "?" | ExportedName .
func (p *gc_parser) parse_name() (string, ast.Expr) ***REMOVED***
	switch p.tok ***REMOVED***
	case scanner.Ident:
		name := p.lit
		p.next()
		return name, ast.NewIdent(name)
	case '?':
		p.next()
		return "?", ast.NewIdent("?")
	case '@':
		en := p.parse_exported_name()
		return en.Sel.Name, en
	***REMOVED***
	p.error("name expected")
	return "", nil
***REMOVED***

// Field = Name Type [ string_lit ] .
func (p *gc_parser) parse_field() *ast.Field ***REMOVED***
	var tag string
	name, _ := p.parse_name()
	typ := p.parse_type()
	if p.tok == scanner.String ***REMOVED***
		tag = p.expect(scanner.String)
	***REMOVED***

	var names []*ast.Ident
	if name != "?" ***REMOVED***
		names = []*ast.Ident***REMOVED***ast.NewIdent(name)***REMOVED***
	***REMOVED***

	return &ast.Field***REMOVED***
		Names: names,
		Type:  typ,
		Tag:   &ast.BasicLit***REMOVED***Kind: token.STRING, Value: tag***REMOVED***,
	***REMOVED***
***REMOVED***

// Parameter = ( identifier | "?" ) [ "..." ] Type [ string_lit ] .
func (p *gc_parser) parse_parameter() *ast.Field ***REMOVED***
	// name
	name, _ := p.parse_name()

	// type
	var typ ast.Expr
	if p.tok == '.' ***REMOVED***
		p.expect_special("...")
		typ = &ast.Ellipsis***REMOVED***Elt: p.parse_type()***REMOVED***
	***REMOVED*** else ***REMOVED***
		typ = p.parse_type()
	***REMOVED***

	var tag string
	if p.tok == scanner.String ***REMOVED***
		tag = p.expect(scanner.String)
	***REMOVED***

	return &ast.Field***REMOVED***
		Names: []*ast.Ident***REMOVED***ast.NewIdent(name)***REMOVED***,
		Type:  typ,
		Tag:   &ast.BasicLit***REMOVED***Kind: token.STRING, Value: tag***REMOVED***,
	***REMOVED***
***REMOVED***

// Parameters = "(" [ ParameterList ] ")" .
// ParameterList = ***REMOVED*** Parameter "," ***REMOVED*** Parameter .
func (p *gc_parser) parse_parameters() *ast.FieldList ***REMOVED***
	flds := []*ast.Field***REMOVED******REMOVED***
	parse_parameter := func() ***REMOVED***
		par := p.parse_parameter()
		flds = append(flds, par)
	***REMOVED***

	p.expect('(')
	if p.tok != ')' ***REMOVED***
		parse_parameter()
		for p.tok == ',' ***REMOVED***
			p.next()
			parse_parameter()
		***REMOVED***
	***REMOVED***
	p.expect(')')
	return &ast.FieldList***REMOVED***List: flds***REMOVED***
***REMOVED***

// Signature = Parameters [ Result ] .
// Result = Type | Parameters .
func (p *gc_parser) parse_signature() *ast.FuncType ***REMOVED***
	var params *ast.FieldList
	var results *ast.FieldList

	params = p.parse_parameters()
	switch p.tok ***REMOVED***
	case scanner.Ident, '[', '*', '<', '@':
		fld := &ast.Field***REMOVED***Type: p.parse_type()***REMOVED***
		results = &ast.FieldList***REMOVED***List: []*ast.Field***REMOVED***fld***REMOVED******REMOVED***
	case '(':
		results = p.parse_parameters()
	***REMOVED***
	return &ast.FuncType***REMOVED***Params: params, Results: results***REMOVED***
***REMOVED***

// MethodOrEmbedSpec = Name [ Signature ] .
func (p *gc_parser) parse_method_or_embed_spec() *ast.Field ***REMOVED***
	name, nameexpr := p.parse_name()
	if p.tok == '(' ***REMOVED***
		typ := p.parse_signature()
		return &ast.Field***REMOVED***
			Names: []*ast.Ident***REMOVED***ast.NewIdent(name)***REMOVED***,
			Type:  typ,
		***REMOVED***
	***REMOVED***

	return &ast.Field***REMOVED***
		Type: nameexpr,
	***REMOVED***
***REMOVED***

// int_lit = [ "-" | "+" ] ***REMOVED*** "0" ... "9" ***REMOVED*** .
func (p *gc_parser) parse_int() ***REMOVED***
	switch p.tok ***REMOVED***
	case '-', '+':
		p.next()
	***REMOVED***
	p.expect(scanner.Int)
***REMOVED***

// number = int_lit [ "p" int_lit ] .
func (p *gc_parser) parse_number() ***REMOVED***
	p.parse_int()
	if p.lit == "p" ***REMOVED***
		p.next()
		p.parse_int()
	***REMOVED***
***REMOVED***

//-------------------------------------------------------------------------------
// gc_parser.types
//-------------------------------------------------------------------------------

// InterfaceType = "interface" "***REMOVED***" [ MethodOrEmbedList ] "***REMOVED***" .
// MethodOrEmbedList = MethodOrEmbedSpec ***REMOVED*** ";" MethodOrEmbedSpec ***REMOVED*** .
func (p *gc_parser) parse_interface_type() ast.Expr ***REMOVED***
	var methods []*ast.Field
	parse_method := func() ***REMOVED***
		meth := p.parse_method_or_embed_spec()
		methods = append(methods, meth)
	***REMOVED***

	p.expect_keyword("interface")
	p.expect('***REMOVED***')
	if p.tok != '***REMOVED***' ***REMOVED***
		parse_method()
		for p.tok == ';' ***REMOVED***
			p.next()
			parse_method()
		***REMOVED***
	***REMOVED***
	p.expect('***REMOVED***')
	return &ast.InterfaceType***REMOVED***Methods: &ast.FieldList***REMOVED***List: methods***REMOVED******REMOVED***
***REMOVED***

// StructType = "struct" "***REMOVED***" [ FieldList ] "***REMOVED***" .
// FieldList = Field ***REMOVED*** ";" Field ***REMOVED*** .
func (p *gc_parser) parse_struct_type() ast.Expr ***REMOVED***
	var fields []*ast.Field
	parse_field := func() ***REMOVED***
		fld := p.parse_field()
		fields = append(fields, fld)
	***REMOVED***

	p.expect_keyword("struct")
	p.expect('***REMOVED***')
	if p.tok != '***REMOVED***' ***REMOVED***
		parse_field()
		for p.tok == ';' ***REMOVED***
			p.next()
			parse_field()
		***REMOVED***
	***REMOVED***
	p.expect('***REMOVED***')
	return &ast.StructType***REMOVED***Fields: &ast.FieldList***REMOVED***List: fields***REMOVED******REMOVED***
***REMOVED***

// MapType = "map" "[" Type "]" Type .
func (p *gc_parser) parse_map_type() ast.Expr ***REMOVED***
	p.expect_keyword("map")
	p.expect('[')
	key := p.parse_type()
	p.expect(']')
	elt := p.parse_type()
	return &ast.MapType***REMOVED***Key: key, Value: elt***REMOVED***
***REMOVED***

// ChanType = ( "chan" [ "<-" ] | "<-" "chan" ) Type .
func (p *gc_parser) parse_chan_type() ast.Expr ***REMOVED***
	dir := ast.SEND | ast.RECV
	if p.tok == scanner.Ident ***REMOVED***
		p.expect_keyword("chan")
		if p.tok == '<' ***REMOVED***
			p.expect_special("<-")
			dir = ast.SEND
		***REMOVED***
	***REMOVED*** else ***REMOVED***
		p.expect_special("<-")
		p.expect_keyword("chan")
		dir = ast.RECV
	***REMOVED***

	elt := p.parse_type()
	return &ast.ChanType***REMOVED***Dir: dir, Value: elt***REMOVED***
***REMOVED***

// ArrayOrSliceType = ArrayType | SliceType .
// ArrayType = "[" int_lit "]" Type .
// SliceType = "[" "]" Type .
func (p *gc_parser) parse_array_or_slice_type() ast.Expr ***REMOVED***
	p.expect('[')
	if p.tok == ']' ***REMOVED***
		// SliceType
		p.next() // skip ']'
		return &ast.ArrayType***REMOVED***Len: nil, Elt: p.parse_type()***REMOVED***
	***REMOVED***

	// ArrayType
	lit := p.expect(scanner.Int)
	p.expect(']')
	return &ast.ArrayType***REMOVED***
		Len: &ast.BasicLit***REMOVED***Kind: token.INT, Value: lit***REMOVED***,
		Elt: p.parse_type(),
	***REMOVED***
***REMOVED***

// Type =
//	BasicType | TypeName | ArrayType | SliceType | StructType |
//      PointerType | FuncType | InterfaceType | MapType | ChanType |
//      "(" Type ")" .
// BasicType = ident .
// TypeName = ExportedName .
// SliceType = "[" "]" Type .
// PointerType = "*" Type .
// FuncType = "func" Signature .
func (p *gc_parser) parse_type() ast.Expr ***REMOVED***
	switch p.tok ***REMOVED***
	case scanner.Ident:
		switch p.lit ***REMOVED***
		case "struct":
			return p.parse_struct_type()
		case "func":
			p.next()
			return p.parse_signature()
		case "interface":
			return p.parse_interface_type()
		case "map":
			return p.parse_map_type()
		case "chan":
			return p.parse_chan_type()
		default:
			lit := p.lit
			p.next()
			return ast.NewIdent(lit)
		***REMOVED***
	case '@':
		return p.parse_exported_name()
	case '[':
		return p.parse_array_or_slice_type()
	case '*':
		p.next()
		return &ast.StarExpr***REMOVED***X: p.parse_type()***REMOVED***
	case '<':
		return p.parse_chan_type()
	case '(':
		p.next()
		typ := p.parse_type()
		p.expect(')')
		return typ
	***REMOVED***
	p.errorf("unexpected token: %s", scanner.TokenString(p.tok))
	return nil
***REMOVED***

//-------------------------------------------------------------------------------
// gc_parser.declarations
//-------------------------------------------------------------------------------

// ImportDecl = "import" identifier string_lit .
func (p *gc_parser) parse_import_decl() ***REMOVED***
	p.expect_keyword("import")
	alias := p.expect(scanner.Ident)
	path := p.parse_package()
	fullName := "!" + path.Name + "!" + alias
	p.path_to_name[path.Name] = fullName
	p.pfc.add_package_to_scope(fullName, path.Name)
***REMOVED***

// ConstDecl   = "const" ExportedName [ Type ] "=" Literal .
// Literal     = bool_lit | int_lit | float_lit | complex_lit | string_lit .
// bool_lit    = "true" | "false" .
// complex_lit = "(" float_lit "+" float_lit ")" .
// rune_lit    = "(" int_lit "+" int_lit ")" .
// string_lit  = `"` ***REMOVED*** unicode_char ***REMOVED*** `"` .
func (p *gc_parser) parse_const_decl() (string, *ast.GenDecl) ***REMOVED***
	// TODO: do we really need actual const value? gocode doesn't use this
	p.expect_keyword("const")
	name := p.parse_exported_name()

	var typ ast.Expr
	if p.tok != '=' ***REMOVED***
		typ = p.parse_type()
	***REMOVED***

	p.expect('=')

	// skip the value
	switch p.tok ***REMOVED***
	case scanner.Ident:
		// must be bool, true or false
		p.next()
	case '-', '+', scanner.Int:
		// number
		p.parse_number()
	case '(':
		// complex_lit or rune_lit
		p.next() // skip '('
		if p.tok == scanner.Char ***REMOVED***
			p.next()
		***REMOVED*** else ***REMOVED***
			p.parse_number()
		***REMOVED***
		p.expect('+')
		p.parse_number()
		p.expect(')')
	case scanner.Char:
		p.next()
	case scanner.String:
		p.next()
	default:
		p.error("expected literal")
	***REMOVED***

	return name.X.(*ast.Ident).Name, &ast.GenDecl***REMOVED***
		Tok: token.CONST,
		Specs: []ast.Spec***REMOVED***
			&ast.ValueSpec***REMOVED***
				Names:  []*ast.Ident***REMOVED***name.Sel***REMOVED***,
				Type:   typ,
				Values: []ast.Expr***REMOVED***&ast.BasicLit***REMOVED***Kind: token.INT, Value: "0"***REMOVED******REMOVED***,
			***REMOVED***,
		***REMOVED***,
	***REMOVED***
***REMOVED***

// TypeDecl = "type" ExportedName Type .
func (p *gc_parser) parse_type_decl() (string, *ast.GenDecl) ***REMOVED***
	p.expect_keyword("type")
	name := p.parse_exported_name()
	typ := p.parse_type()
	return name.X.(*ast.Ident).Name, &ast.GenDecl***REMOVED***
		Tok: token.TYPE,
		Specs: []ast.Spec***REMOVED***
			&ast.TypeSpec***REMOVED***
				Name: name.Sel,
				Type: typ,
			***REMOVED***,
		***REMOVED***,
	***REMOVED***
***REMOVED***

// VarDecl = "var" ExportedName Type .
func (p *gc_parser) parse_var_decl() (string, *ast.GenDecl) ***REMOVED***
	p.expect_keyword("var")
	name := p.parse_exported_name()
	typ := p.parse_type()
	return name.X.(*ast.Ident).Name, &ast.GenDecl***REMOVED***
		Tok: token.VAR,
		Specs: []ast.Spec***REMOVED***
			&ast.ValueSpec***REMOVED***
				Names: []*ast.Ident***REMOVED***name.Sel***REMOVED***,
				Type:  typ,
			***REMOVED***,
		***REMOVED***,
	***REMOVED***
***REMOVED***

// FuncBody = "***REMOVED***" ... "***REMOVED***" .
func (p *gc_parser) parse_func_body() ***REMOVED***
	p.expect('***REMOVED***')
	for i := 1; i > 0; p.next() ***REMOVED***
		switch p.tok ***REMOVED***
		case '***REMOVED***':
			i++
		case '***REMOVED***':
			i--
		***REMOVED***
	***REMOVED***
***REMOVED***

// FuncDecl = "func" ExportedName Signature [ FuncBody ] .
func (p *gc_parser) parse_func_decl() (string, *ast.FuncDecl) ***REMOVED***
	// "func" was already consumed by lookahead
	name := p.parse_exported_name()
	typ := p.parse_signature()
	if p.tok == '***REMOVED***' ***REMOVED***
		p.parse_func_body()
	***REMOVED***
	return name.X.(*ast.Ident).Name, &ast.FuncDecl***REMOVED***
		Name: name.Sel,
		Type: typ,
	***REMOVED***
***REMOVED***

func strip_method_receiver(recv *ast.FieldList) string ***REMOVED***
	var sel *ast.SelectorExpr

	// find selector expression
	typ := recv.List[0].Type
	switch t := typ.(type) ***REMOVED***
	case *ast.StarExpr:
		sel = t.X.(*ast.SelectorExpr)
	case *ast.SelectorExpr:
		sel = t
	***REMOVED***

	// extract package path
	pkg := sel.X.(*ast.Ident).Name

	// write back stripped type
	switch t := typ.(type) ***REMOVED***
	case *ast.StarExpr:
		t.X = sel.Sel
	case *ast.SelectorExpr:
		recv.List[0].Type = sel.Sel
	***REMOVED***

	return pkg
***REMOVED***

// MethodDecl = "func" Receiver Name Signature .
// Receiver = "(" ( identifier | "?" ) [ "*" ] ExportedName ")" [ FuncBody ] .
func (p *gc_parser) parse_method_decl() (string, *ast.FuncDecl) ***REMOVED***
	recv := p.parse_parameters()
	pkg := strip_method_receiver(recv)
	name, _ := p.parse_name()
	typ := p.parse_signature()
	if p.tok == '***REMOVED***' ***REMOVED***
		p.parse_func_body()
	***REMOVED***
	return pkg, &ast.FuncDecl***REMOVED***
		Recv: recv,
		Name: ast.NewIdent(name),
		Type: typ,
	***REMOVED***
***REMOVED***

// Decl = [ ImportDecl | ConstDecl | TypeDecl | VarDecl | FuncDecl | MethodDecl ] "\n" .
func (p *gc_parser) parse_decl() (pkg string, decl ast.Decl) ***REMOVED***
	switch p.lit ***REMOVED***
	case "import":
		p.parse_import_decl()
	case "const":
		pkg, decl = p.parse_const_decl()
	case "type":
		pkg, decl = p.parse_type_decl()
	case "var":
		pkg, decl = p.parse_var_decl()
	case "func":
		p.next()
		if p.tok == '(' ***REMOVED***
			pkg, decl = p.parse_method_decl()
		***REMOVED*** else ***REMOVED***
			pkg, decl = p.parse_func_decl()
		***REMOVED***
	***REMOVED***
	p.expect('\n')
	return
***REMOVED***

// Export = PackageClause ***REMOVED*** Decl ***REMOVED*** "$$" .
// PackageClause = "package" identifier [ "safe" ] "\n" .
func (p *gc_parser) parse_export(callback func(string, ast.Decl)) ***REMOVED***
	p.expect_keyword("package")
	p.pfc.defalias = p.expect(scanner.Ident)
	if p.tok != '\n' ***REMOVED***
		p.expect_keyword("safe")
	***REMOVED***
	p.expect('\n')

	for p.tok != '$' && p.tok != scanner.EOF ***REMOVED***
		pkg, decl := p.parse_decl()
		if decl != nil ***REMOVED***
			callback(pkg, decl)
		***REMOVED***
	***REMOVED***
***REMOVED***
