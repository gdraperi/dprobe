package main

import (
	"encoding/binary"
	"fmt"
	"go/ast"
	"go/token"
	"strconv"
	"strings"
	"unicode"
	"unicode/utf8"
)

//-------------------------------------------------------------------------
// gc_bin_parser
//
// The following part of the code may contain portions of the code from the Go
// standard library, which tells me to retain their copyright notice:
//
// Copyright (c) 2012 The Go Authors. All rights reserved.
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

type gc_bin_parser struct ***REMOVED***
	data    []byte
	buf     []byte // for reading strings
	version int    // export format version

	// object lists
	strList       []string   // in order of appearance
	pathList      []string   // in order of appearance
	pkgList       []string   // in order of appearance
	typList       []ast.Expr // in order of appearance
	callback      func(pkg string, decl ast.Decl)
	pfc           *package_file_cache
	trackAllTypes bool

	// position encoding
	posInfoFormat bool
	prevFile      string
	prevLine      int

	// debugging support
	debugFormat bool
	read        int // bytes read

***REMOVED***

func (p *gc_bin_parser) init(data []byte, pfc *package_file_cache) ***REMOVED***
	p.data = data
	p.version = -1            // unknown version
	p.strList = []string***REMOVED***""***REMOVED***  // empty string is mapped to 0
	p.pathList = []string***REMOVED***""***REMOVED*** // empty string is mapped to 0
	p.pfc = pfc
***REMOVED***

func (p *gc_bin_parser) parse_export(callback func(string, ast.Decl)) ***REMOVED***
	p.callback = callback

	// read version info
	var versionstr string
	if b := p.rawByte(); b == 'c' || b == 'd' ***REMOVED***
		// Go1.7 encoding; first byte encodes low-level
		// encoding format (compact vs debug).
		// For backward-compatibility only (avoid problems with
		// old installed packages). Newly compiled packages use
		// the extensible format string.
		// TODO(gri) Remove this support eventually; after Go1.8.
		if b == 'd' ***REMOVED***
			p.debugFormat = true
		***REMOVED***
		p.trackAllTypes = p.rawByte() == 'a'
		p.posInfoFormat = p.int() != 0
		versionstr = p.string()
		if versionstr == "v1" ***REMOVED***
			p.version = 0
		***REMOVED***
	***REMOVED*** else ***REMOVED***
		// Go1.8 extensible encoding
		// read version string and extract version number (ignore anything after the version number)
		versionstr = p.rawStringln(b)
		if s := strings.SplitN(versionstr, " ", 3); len(s) >= 2 && s[0] == "version" ***REMOVED***
			if v, err := strconv.Atoi(s[1]); err == nil && v > 0 ***REMOVED***
				p.version = v
			***REMOVED***
		***REMOVED***
	***REMOVED***

	// read version specific flags - extend as necessary
	switch p.version ***REMOVED***
	// case 6:
	// 	...
	//	fallthrough
	case 5, 4, 3, 2, 1:
		p.debugFormat = p.rawStringln(p.rawByte()) == "debug"
		p.trackAllTypes = p.int() != 0
		p.posInfoFormat = p.int() != 0
	case 0:
		// Go1.7 encoding format - nothing to do here
	default:
		panic(fmt.Errorf("unknown export format version %d (%q)", p.version, versionstr))
	***REMOVED***

	// --- generic export data ---

	// populate typList with predeclared "known" types
	p.typList = append(p.typList, predeclared...)

	// read package data
	pkgName := p.pkg()
	p.pfc.defalias = pkgName[strings.LastIndex(pkgName, "!")+1:]

	// read objects of phase 1 only (see cmd/compiler/internal/gc/bexport.go)
	objcount := 0
	for ***REMOVED***
		tag := p.tagOrIndex()
		if tag == endTag ***REMOVED***
			break
		***REMOVED***
		p.obj(tag)
		objcount++
	***REMOVED***

	// self-verification
	if count := p.int(); count != objcount ***REMOVED***
		panic(fmt.Sprintf("got %d objects; want %d", objcount, count))
	***REMOVED***
***REMOVED***

func (p *gc_bin_parser) pkg() string ***REMOVED***
	// if the package was seen before, i is its index (>= 0)
	i := p.tagOrIndex()
	if i >= 0 ***REMOVED***
		return p.pkgList[i]
	***REMOVED***

	// otherwise, i is the package tag (< 0)
	if i != packageTag ***REMOVED***
		panic(fmt.Sprintf("unexpected package tag %d version %d", i, p.version))
	***REMOVED***

	// read package data
	name := p.string()
	var path string
	if p.version >= 5 ***REMOVED***
		path = p.path()
	***REMOVED*** else ***REMOVED***
		path = p.string()
	***REMOVED***

	// we should never see an empty package name
	if name == "" ***REMOVED***
		panic("empty package name in import")
	***REMOVED***

	// an empty path denotes the package we are currently importing;
	// it must be the first package we see
	if (path == "") != (len(p.pkgList) == 0) ***REMOVED***
		panic(fmt.Sprintf("package path %q for pkg index %d", path, len(p.pkgList)))
	***REMOVED***

	var fullName string
	if path != "" ***REMOVED***
		fullName = "!" + path + "!" + name
		p.pfc.add_package_to_scope(fullName, path)
	***REMOVED*** else ***REMOVED***
		fullName = "!" + p.pfc.name + "!" + name
	***REMOVED***

	// if the package was imported before, use that one; otherwise create a new one
	p.pkgList = append(p.pkgList, fullName)
	return p.pkgList[len(p.pkgList)-1]
***REMOVED***

func (p *gc_bin_parser) obj(tag int) ***REMOVED***
	switch tag ***REMOVED***
	case constTag:
		p.pos()
		pkg, name := p.qualifiedName()
		typ := p.typ("")
		p.skipValue() // ignore const value, gocode's not interested
		p.callback(pkg, &ast.GenDecl***REMOVED***
			Tok: token.CONST,
			Specs: []ast.Spec***REMOVED***
				&ast.ValueSpec***REMOVED***
					Names:  []*ast.Ident***REMOVED***ast.NewIdent(name)***REMOVED***,
					Type:   typ,
					Values: []ast.Expr***REMOVED***&ast.BasicLit***REMOVED***Kind: token.INT, Value: "0"***REMOVED******REMOVED***,
				***REMOVED***,
			***REMOVED***,
		***REMOVED***)

	case aliasTag:
		// TODO(gri) verify type alias hookup is correct
		p.pos()
		pkg, name := p.qualifiedName()
		typ := p.typ("")
		p.callback(pkg, &ast.GenDecl***REMOVED***
			Tok:   token.TYPE,
			Specs: []ast.Spec***REMOVED***typeAliasSpec(name, typ)***REMOVED***,
		***REMOVED***)

	case typeTag:
		_ = p.typ("")

	case varTag:
		p.pos()
		pkg, name := p.qualifiedName()
		typ := p.typ("")
		p.callback(pkg, &ast.GenDecl***REMOVED***
			Tok: token.VAR,
			Specs: []ast.Spec***REMOVED***
				&ast.ValueSpec***REMOVED***
					Names: []*ast.Ident***REMOVED***ast.NewIdent(name)***REMOVED***,
					Type:  typ,
				***REMOVED***,
			***REMOVED***,
		***REMOVED***)

	case funcTag:
		p.pos()
		pkg, name := p.qualifiedName()
		params := p.paramList()
		results := p.paramList()
		p.callback(pkg, &ast.FuncDecl***REMOVED***
			Name: ast.NewIdent(name),
			Type: &ast.FuncType***REMOVED***Params: params, Results: results***REMOVED***,
		***REMOVED***)

	default:
		panic(fmt.Sprintf("unexpected object tag %d", tag))
	***REMOVED***
***REMOVED***

const deltaNewFile = -64 // see cmd/compile/internal/gc/bexport.go

func (p *gc_bin_parser) pos() ***REMOVED***
	if !p.posInfoFormat ***REMOVED***
		return
	***REMOVED***

	file := p.prevFile
	line := p.prevLine
	delta := p.int()
	line += delta
	if p.version >= 5 ***REMOVED***
		if delta == deltaNewFile ***REMOVED***
			if n := p.int(); n >= 0 ***REMOVED***
				// file changed
				file = p.path()
				line = n
			***REMOVED***
		***REMOVED***
	***REMOVED*** else ***REMOVED***
		if delta == 0 ***REMOVED***
			if n := p.int(); n >= 0 ***REMOVED***
				// file changed
				file = p.prevFile[:n] + p.string()
				line = p.int()
			***REMOVED***
		***REMOVED***
	***REMOVED***
	p.prevFile = file
	p.prevLine = line

	// TODO(gri) register new position
***REMOVED***

func (p *gc_bin_parser) qualifiedName() (pkg string, name string) ***REMOVED***
	name = p.string()
	pkg = p.pkg()
	return pkg, name
***REMOVED***

func (p *gc_bin_parser) reserveMaybe() int ***REMOVED***
	if p.trackAllTypes ***REMOVED***
		p.typList = append(p.typList, nil)
		return len(p.typList) - 1
	***REMOVED*** else ***REMOVED***
		return -1
	***REMOVED***
***REMOVED***

func (p *gc_bin_parser) recordMaybe(idx int, t ast.Expr) ast.Expr ***REMOVED***
	if idx == -1 ***REMOVED***
		return t
	***REMOVED***
	p.typList[idx] = t
	return t
***REMOVED***

func (p *gc_bin_parser) record(t ast.Expr) ***REMOVED***
	p.typList = append(p.typList, t)
***REMOVED***

// parent is the package which declared the type; parent == nil means
// the package currently imported. The parent package is needed for
// exported struct fields and interface methods which don't contain
// explicit package information in the export data.
func (p *gc_bin_parser) typ(parent string) ast.Expr ***REMOVED***
	// if the type was seen before, i is its index (>= 0)
	i := p.tagOrIndex()
	if i >= 0 ***REMOVED***
		return p.typList[i]
	***REMOVED***

	// otherwise, i is the type tag (< 0)
	switch i ***REMOVED***
	case namedTag:
		// read type object
		p.pos()
		parent, name := p.qualifiedName()
		tdecl := &ast.GenDecl***REMOVED***
			Tok: token.TYPE,
			Specs: []ast.Spec***REMOVED***
				&ast.TypeSpec***REMOVED***
					Name: ast.NewIdent(name),
				***REMOVED***,
			***REMOVED***,
		***REMOVED***

		// record it right away (underlying type can contain refs to t)
		t := &ast.SelectorExpr***REMOVED***X: ast.NewIdent(parent), Sel: ast.NewIdent(name)***REMOVED***
		p.record(t)

		// parse underlying type
		t0 := p.typ(parent)
		tdecl.Specs[0].(*ast.TypeSpec).Type = t0

		p.callback(parent, tdecl)

		// interfaces have no methods
		if _, ok := t0.(*ast.InterfaceType); ok ***REMOVED***
			return t
		***REMOVED***

		// read associated methods
		for i := p.int(); i > 0; i-- ***REMOVED***
			// TODO(gri) replace this with something closer to fieldName
			p.pos()
			name := p.string()
			if !exported(name) ***REMOVED***
				p.pkg()
			***REMOVED***

			recv := p.paramList()
			params := p.paramList()
			results := p.paramList()
			p.int() // go:nointerface pragma - discarded

			strip_method_receiver(recv)
			p.callback(parent, &ast.FuncDecl***REMOVED***
				Recv: recv,
				Name: ast.NewIdent(name),
				Type: &ast.FuncType***REMOVED***Params: params, Results: results***REMOVED***,
			***REMOVED***)
		***REMOVED***
		return t
	case arrayTag:
		i := p.reserveMaybe()
		n := p.int64()
		elt := p.typ(parent)
		return p.recordMaybe(i, &ast.ArrayType***REMOVED***
			Len: &ast.BasicLit***REMOVED***Kind: token.INT, Value: fmt.Sprint(n)***REMOVED***,
			Elt: elt,
		***REMOVED***)

	case sliceTag:
		i := p.reserveMaybe()
		elt := p.typ(parent)
		return p.recordMaybe(i, &ast.ArrayType***REMOVED***Len: nil, Elt: elt***REMOVED***)

	case dddTag:
		i := p.reserveMaybe()
		elt := p.typ(parent)
		return p.recordMaybe(i, &ast.Ellipsis***REMOVED***Elt: elt***REMOVED***)

	case structTag:
		i := p.reserveMaybe()
		return p.recordMaybe(i, p.structType(parent))

	case pointerTag:
		i := p.reserveMaybe()
		elt := p.typ(parent)
		return p.recordMaybe(i, &ast.StarExpr***REMOVED***X: elt***REMOVED***)

	case signatureTag:
		i := p.reserveMaybe()
		params := p.paramList()
		results := p.paramList()
		return p.recordMaybe(i, &ast.FuncType***REMOVED***Params: params, Results: results***REMOVED***)

	case interfaceTag:
		i := p.reserveMaybe()
		var embeddeds []*ast.SelectorExpr
		for n := p.int(); n > 0; n-- ***REMOVED***
			p.pos()
			if named, ok := p.typ(parent).(*ast.SelectorExpr); ok ***REMOVED***
				embeddeds = append(embeddeds, named)
			***REMOVED***
		***REMOVED***
		methods := p.methodList(parent)
		for _, field := range embeddeds ***REMOVED***
			methods = append(methods, &ast.Field***REMOVED***Type: field***REMOVED***)
		***REMOVED***
		return p.recordMaybe(i, &ast.InterfaceType***REMOVED***Methods: &ast.FieldList***REMOVED***List: methods***REMOVED******REMOVED***)

	case mapTag:
		i := p.reserveMaybe()
		key := p.typ(parent)
		val := p.typ(parent)
		return p.recordMaybe(i, &ast.MapType***REMOVED***Key: key, Value: val***REMOVED***)

	case chanTag:
		i := p.reserveMaybe()
		dir := ast.SEND | ast.RECV
		switch d := p.int(); d ***REMOVED***
		case 1:
			dir = ast.RECV
		case 2:
			dir = ast.SEND
		case 3:
			// already set
		default:
			panic(fmt.Sprintf("unexpected channel dir %d", d))
		***REMOVED***
		elt := p.typ(parent)
		return p.recordMaybe(i, &ast.ChanType***REMOVED***Dir: dir, Value: elt***REMOVED***)

	default:
		panic(fmt.Sprintf("unexpected type tag %d", i))
	***REMOVED***
***REMOVED***

func (p *gc_bin_parser) structType(parent string) *ast.StructType ***REMOVED***
	var fields []*ast.Field
	if n := p.int(); n > 0 ***REMOVED***
		fields = make([]*ast.Field, n)
		for i := range fields ***REMOVED***
			fields[i], _ = p.field(parent) // (*ast.Field, tag), not interested in tags
		***REMOVED***
	***REMOVED***
	return &ast.StructType***REMOVED***Fields: &ast.FieldList***REMOVED***List: fields***REMOVED******REMOVED***
***REMOVED***

func (p *gc_bin_parser) field(parent string) (*ast.Field, string) ***REMOVED***
	p.pos()
	_, name, _ := p.fieldName(parent)
	typ := p.typ(parent)
	tag := p.string()

	var names []*ast.Ident
	if name != "" ***REMOVED***
		names = []*ast.Ident***REMOVED***ast.NewIdent(name)***REMOVED***
	***REMOVED***
	return &ast.Field***REMOVED***
		Names: names,
		Type:  typ,
	***REMOVED***, tag
***REMOVED***

func (p *gc_bin_parser) methodList(parent string) (methods []*ast.Field) ***REMOVED***
	if n := p.int(); n > 0 ***REMOVED***
		methods = make([]*ast.Field, n)
		for i := range methods ***REMOVED***
			methods[i] = p.method(parent)
		***REMOVED***
	***REMOVED***
	return
***REMOVED***

func (p *gc_bin_parser) method(parent string) *ast.Field ***REMOVED***
	p.pos()
	_, name, _ := p.fieldName(parent)
	params := p.paramList()
	results := p.paramList()
	return &ast.Field***REMOVED***
		Names: []*ast.Ident***REMOVED***ast.NewIdent(name)***REMOVED***,
		Type:  &ast.FuncType***REMOVED***Params: params, Results: results***REMOVED***,
	***REMOVED***
***REMOVED***

func (p *gc_bin_parser) fieldName(parent string) (string, string, bool) ***REMOVED***
	name := p.string()
	pkg := parent
	if p.version == 0 && name == "_" ***REMOVED***
		// version 0 didn't export a package for _ fields
		return pkg, name, false
	***REMOVED***
	var alias bool
	switch name ***REMOVED***
	case "":
		// 1) field name matches base type name and is exported: nothing to do
	case "?":
		// 2) field name matches base type name and is not exported: need package
		name = ""
		pkg = p.pkg()
	case "@":
		// 3) field name doesn't match type name (alias)
		name = p.string()
		alias = true
		fallthrough
	default:
		if !exported(name) ***REMOVED***
			pkg = p.pkg()
		***REMOVED***
	***REMOVED***
	return pkg, name, alias
***REMOVED***

func (p *gc_bin_parser) paramList() *ast.FieldList ***REMOVED***
	n := p.int()
	if n == 0 ***REMOVED***
		return nil
	***REMOVED***
	// negative length indicates unnamed parameters
	named := true
	if n < 0 ***REMOVED***
		n = -n
		named = false
	***REMOVED***
	// n > 0
	flds := make([]*ast.Field, n)
	for i := range flds ***REMOVED***
		flds[i] = p.param(named)
	***REMOVED***
	return &ast.FieldList***REMOVED***List: flds***REMOVED***
***REMOVED***

func (p *gc_bin_parser) param(named bool) *ast.Field ***REMOVED***
	t := p.typ("")

	name := "?"
	if named ***REMOVED***
		name = p.string()
		if name == "" ***REMOVED***
			panic("expected named parameter")
		***REMOVED***
		if name != "_" ***REMOVED***
			p.pkg()
		***REMOVED***
		if i := strings.Index(name, "·"); i > 0 ***REMOVED***
			name = name[:i] // cut off gc-specific parameter numbering
		***REMOVED***
	***REMOVED***

	// read and discard compiler-specific info
	p.string()

	return &ast.Field***REMOVED***
		Names: []*ast.Ident***REMOVED***ast.NewIdent(name)***REMOVED***,
		Type:  t,
	***REMOVED***
***REMOVED***

func exported(name string) bool ***REMOVED***
	ch, _ := utf8.DecodeRuneInString(name)
	return unicode.IsUpper(ch)
***REMOVED***

func (p *gc_bin_parser) skipValue() ***REMOVED***
	switch tag := p.tagOrIndex(); tag ***REMOVED***
	case falseTag, trueTag:
	case int64Tag:
		p.int64()
	case floatTag:
		p.float()
	case complexTag:
		p.float()
		p.float()
	case stringTag:
		p.string()
	default:
		panic(fmt.Sprintf("unexpected value tag %d", tag))
	***REMOVED***
***REMOVED***

func (p *gc_bin_parser) float() ***REMOVED***
	sign := p.int()
	if sign == 0 ***REMOVED***
		return
	***REMOVED***

	p.int()    // exp
	p.string() // mant
***REMOVED***

// ----------------------------------------------------------------------------
// Low-level decoders

func (p *gc_bin_parser) tagOrIndex() int ***REMOVED***
	if p.debugFormat ***REMOVED***
		p.marker('t')
	***REMOVED***

	return int(p.rawInt64())
***REMOVED***

func (p *gc_bin_parser) int() int ***REMOVED***
	x := p.int64()
	if int64(int(x)) != x ***REMOVED***
		panic("exported integer too large")
	***REMOVED***
	return int(x)
***REMOVED***

func (p *gc_bin_parser) int64() int64 ***REMOVED***
	if p.debugFormat ***REMOVED***
		p.marker('i')
	***REMOVED***

	return p.rawInt64()
***REMOVED***

func (p *gc_bin_parser) path() string ***REMOVED***
	if p.debugFormat ***REMOVED***
		p.marker('p')
	***REMOVED***
	// if the path was seen before, i is its index (>= 0)
	// (the empty string is at index 0)
	i := p.rawInt64()
	if i >= 0 ***REMOVED***
		return p.pathList[i]
	***REMOVED***
	// otherwise, i is the negative path length (< 0)
	a := make([]string, -i)
	for n := range a ***REMOVED***
		a[n] = p.string()
	***REMOVED***
	s := strings.Join(a, "/")
	p.pathList = append(p.pathList, s)
	return s
***REMOVED***

func (p *gc_bin_parser) string() string ***REMOVED***
	if p.debugFormat ***REMOVED***
		p.marker('s')
	***REMOVED***
	// if the string was seen before, i is its index (>= 0)
	// (the empty string is at index 0)
	i := p.rawInt64()
	if i >= 0 ***REMOVED***
		return p.strList[i]
	***REMOVED***
	// otherwise, i is the negative string length (< 0)
	if n := int(-i); n <= cap(p.buf) ***REMOVED***
		p.buf = p.buf[:n]
	***REMOVED*** else ***REMOVED***
		p.buf = make([]byte, n)
	***REMOVED***
	for i := range p.buf ***REMOVED***
		p.buf[i] = p.rawByte()
	***REMOVED***
	s := string(p.buf)
	p.strList = append(p.strList, s)
	return s
***REMOVED***

func (p *gc_bin_parser) marker(want byte) ***REMOVED***
	if got := p.rawByte(); got != want ***REMOVED***
		panic(fmt.Sprintf("incorrect marker: got %c; want %c (pos = %d)", got, want, p.read))
	***REMOVED***

	pos := p.read
	if n := int(p.rawInt64()); n != pos ***REMOVED***
		panic(fmt.Sprintf("incorrect position: got %d; want %d", n, pos))
	***REMOVED***
***REMOVED***

// rawInt64 should only be used by low-level decoders.
func (p *gc_bin_parser) rawInt64() int64 ***REMOVED***
	i, err := binary.ReadVarint(p)
	if err != nil ***REMOVED***
		panic(fmt.Sprintf("read error: %v", err))
	***REMOVED***
	return i
***REMOVED***

// rawStringln should only be used to read the initial version string.
func (p *gc_bin_parser) rawStringln(b byte) string ***REMOVED***
	p.buf = p.buf[:0]
	for b != '\n' ***REMOVED***
		p.buf = append(p.buf, b)
		b = p.rawByte()
	***REMOVED***
	return string(p.buf)
***REMOVED***

// needed for binary.ReadVarint in rawInt64
func (p *gc_bin_parser) ReadByte() (byte, error) ***REMOVED***
	return p.rawByte(), nil
***REMOVED***

// byte is the bottleneck interface for reading p.data.
// It unescapes '|' 'S' to '$' and '|' '|' to '|'.
// rawByte should only be used by low-level decoders.
func (p *gc_bin_parser) rawByte() byte ***REMOVED***
	b := p.data[0]
	r := 1
	if b == '|' ***REMOVED***
		b = p.data[1]
		r = 2
		switch b ***REMOVED***
		case 'S':
			b = '$'
		case '|':
			// nothing to do
		default:
			panic("unexpected escape sequence in export data")
		***REMOVED***
	***REMOVED***
	p.data = p.data[r:]
	p.read += r
	return b

***REMOVED***

// ----------------------------------------------------------------------------
// Export format

// Tags. Must be < 0.
const (
	// Objects
	packageTag = -(iota + 1)
	constTag
	typeTag
	varTag
	funcTag
	endTag

	// Types
	namedTag
	arrayTag
	sliceTag
	dddTag
	structTag
	pointerTag
	signatureTag
	interfaceTag
	mapTag
	chanTag

	// Values
	falseTag
	trueTag
	int64Tag
	floatTag
	fractionTag // not used by gc
	complexTag
	stringTag
	nilTag     // only used by gc (appears in exported inlined function bodies)
	unknownTag // not used by gc (only appears in packages with errors)

	// Type aliases
	aliasTag
)

var predeclared = []ast.Expr***REMOVED***
	// basic types
	ast.NewIdent("bool"),
	ast.NewIdent("int"),
	ast.NewIdent("int8"),
	ast.NewIdent("int16"),
	ast.NewIdent("int32"),
	ast.NewIdent("int64"),
	ast.NewIdent("uint"),
	ast.NewIdent("uint8"),
	ast.NewIdent("uint16"),
	ast.NewIdent("uint32"),
	ast.NewIdent("uint64"),
	ast.NewIdent("uintptr"),
	ast.NewIdent("float32"),
	ast.NewIdent("float64"),
	ast.NewIdent("complex64"),
	ast.NewIdent("complex128"),
	ast.NewIdent("string"),

	// basic type aliases
	ast.NewIdent("byte"),
	ast.NewIdent("rune"),

	// error
	ast.NewIdent("error"),

	// TODO(nsf): don't think those are used in just package type info,
	// maybe for consts, but we are not interested in that
	// untyped types
	ast.NewIdent(">_<"), // TODO: types.Typ[types.UntypedBool],
	ast.NewIdent(">_<"), // TODO: types.Typ[types.UntypedInt],
	ast.NewIdent(">_<"), // TODO: types.Typ[types.UntypedRune],
	ast.NewIdent(">_<"), // TODO: types.Typ[types.UntypedFloat],
	ast.NewIdent(">_<"), // TODO: types.Typ[types.UntypedComplex],
	ast.NewIdent(">_<"), // TODO: types.Typ[types.UntypedString],
	ast.NewIdent(">_<"), // TODO: types.Typ[types.UntypedNil],

	// package unsafe
	&ast.SelectorExpr***REMOVED***X: ast.NewIdent("unsafe"), Sel: ast.NewIdent("Pointer")***REMOVED***,

	// invalid type
	ast.NewIdent(">_<"), // TODO: types.Typ[types.Invalid], // only appears in packages with errors

	// used internally by gc; never used by this package or in .a files
	ast.NewIdent("any"),
***REMOVED***
