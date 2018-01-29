package main

import "debug/elf"
import "text/scanner"
import "bytes"
import "errors"
import "io"
import "fmt"
import "strconv"
import "go/ast"
import "go/token"
import "strings"

var builtin_type_names = []*ast.Ident***REMOVED***
	nil,
	ast.NewIdent("int8"),
	ast.NewIdent("int16"),
	ast.NewIdent("int32"),
	ast.NewIdent("int64"),
	ast.NewIdent("uint8"),
	ast.NewIdent("uint16"),
	ast.NewIdent("uint32"),
	ast.NewIdent("uint64"),
	ast.NewIdent("float32"),
	ast.NewIdent("float64"),
	ast.NewIdent("int"),
	ast.NewIdent("uint"),
	ast.NewIdent("uintptr"),
	nil,
	ast.NewIdent("bool"),
	ast.NewIdent("string"),
	ast.NewIdent("complex64"),
	ast.NewIdent("complex128"),
	ast.NewIdent("error"),
	ast.NewIdent("byte"),
	ast.NewIdent("rune"),
***REMOVED***

const (
	smallest_builtin_code = -21
)

func read_import_data(import_path string) ([]byte, error) ***REMOVED***
	// TODO: find file location
	filename := import_path + ".gox"

	f, err := elf.Open(filename)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	defer f.Close()

	sec := f.Section(".go_export")
	if sec == nil ***REMOVED***
		return nil, errors.New("missing .go_export section in the file: " + filename)
	***REMOVED***

	return sec.Data()
***REMOVED***

func parse_import_data(data []byte) ***REMOVED***
	buf := bytes.NewBuffer(data)
	var p import_data_parser
	p.init(buf)

	// magic
	p.expect_ident("v1")
	p.expect(';')

	// package ident
	p.expect_ident("package")
	pkgid := p.expect(scanner.Ident)
	p.expect(';')

	println("package ident: " + pkgid)

	// package path
	p.expect_ident("pkgpath")
	pkgpath := p.expect(scanner.Ident)
	p.expect(';')

	println("package path: " + pkgpath)

	// package priority
	p.expect_ident("priority")
	priority := p.expect(scanner.Int)
	p.expect(';')

	println("package priority: " + priority)

	// import init functions
	for p.toktype == scanner.Ident && p.token() == "import" ***REMOVED***
		p.expect_ident("import")
		pkgname := p.expect(scanner.Ident)
		pkgpath := p.expect(scanner.Ident)
		importpath := p.expect(scanner.String)
		p.expect(';')
		println("import " + pkgname + " " + pkgpath + " " + importpath)
	***REMOVED***

	if p.toktype == scanner.Ident && p.token() == "init" ***REMOVED***
		p.expect_ident("init")
		for p.toktype != ';' ***REMOVED***
			pkgname := p.expect(scanner.Ident)
			initname := p.expect(scanner.Ident)
			prio := p.expect(scanner.Int)
			println("init " + pkgname + " " + initname + " " + fmt.Sprint(prio))
		***REMOVED***
		p.expect(';')
	***REMOVED***

loop:
	for ***REMOVED***
		switch tok := p.expect(scanner.Ident); tok ***REMOVED***
		case "const":
			p.read_const()
		case "type":
			p.read_type_decl()
		case "var":
			p.read_var()
		case "func":
			p.read_func()
		case "checksum":
			p.read_checksum()
			break loop
		default:
			panic(errors.New("unexpected identifier token: '" + tok + "'"))
		***REMOVED***
	***REMOVED***
***REMOVED***

//----------------------------------------------------------------------------
// import data parser
//----------------------------------------------------------------------------

type import_data_type struct ***REMOVED***
	name  string
	type_ ast.Expr
***REMOVED***

type import_data_parser struct ***REMOVED***
	scanner   scanner.Scanner
	toktype   rune
	typetable []*import_data_type
***REMOVED***

func (this *import_data_parser) init(reader io.Reader) ***REMOVED***
	this.scanner.Mode = scanner.ScanIdents | scanner.ScanInts | scanner.ScanStrings | scanner.ScanFloats
	this.scanner.Init(reader)
	this.next()

	// len == 1 here, because 0 is an invalid type index
	this.typetable = make([]*import_data_type, 1, 50)
***REMOVED***

func (this *import_data_parser) next() ***REMOVED***
	this.toktype = this.scanner.Scan()
***REMOVED***

func (this *import_data_parser) token() string ***REMOVED***
	return this.scanner.TokenText()
***REMOVED***

// internal, use expect(scanner.Ident) instead
func (this *import_data_parser) read_ident() string ***REMOVED***
	id := ""
	prev := rune(0)

loop:
	for ***REMOVED***
		switch this.toktype ***REMOVED***
		case scanner.Ident:
			if prev == scanner.Ident ***REMOVED***
				break loop
			***REMOVED***

			prev = this.toktype
			id += this.token()
			this.next()
		case '.', '?', '$':
			prev = this.toktype
			id += string(this.toktype)
			this.next()
		default:
			break loop
		***REMOVED***
	***REMOVED***

	if id == "" ***REMOVED***
		this.errorf("identifier expected, got %s", scanner.TokenString(this.toktype))
	***REMOVED***
	return id
***REMOVED***

func (this *import_data_parser) read_int() string ***REMOVED***
	val := ""
	if this.toktype == '-' ***REMOVED***
		this.next()
		val += "-"
	***REMOVED***
	if this.toktype != scanner.Int ***REMOVED***
		this.errorf("expected: %s, got: %s", scanner.TokenString(scanner.Int), scanner.TokenString(this.toktype))
	***REMOVED***

	val += this.token()
	this.next()
	return val
***REMOVED***

func (this *import_data_parser) errorf(format string, args ...interface***REMOVED******REMOVED***) ***REMOVED***
	panic(errors.New(fmt.Sprintf(format, args...)))
***REMOVED***

// makes sure that the current token is 'x', returns it and reads the next one
func (this *import_data_parser) expect(x rune) string ***REMOVED***
	if x == scanner.Ident ***REMOVED***
		// special case, in gccgo import data identifier is not exactly a scanner.Ident
		return this.read_ident()
	***REMOVED***

	if x == scanner.Int ***REMOVED***
		// another special case, handle negative ints as well
		return this.read_int()
	***REMOVED***

	if this.toktype != x ***REMOVED***
		this.errorf("expected: %s, got: %s", scanner.TokenString(x), scanner.TokenString(this.toktype))
	***REMOVED***

	tok := this.token()
	this.next()
	return tok
***REMOVED***

// makes sure that the following set of tokens matches 'special', reads the next one
func (this *import_data_parser) expect_special(special string) ***REMOVED***
	i := 0
	for i < len(special) ***REMOVED***
		if this.toktype != rune(special[i]) ***REMOVED***
			break
		***REMOVED***

		this.next()
		i++
	***REMOVED***

	if i < len(special) ***REMOVED***
		this.errorf("expected: \"%s\", got something else", special)
	***REMOVED***
***REMOVED***

// makes sure that the current token is scanner.Ident and is equals to 'ident', reads the next one
func (this *import_data_parser) expect_ident(ident string) ***REMOVED***
	tok := this.expect(scanner.Ident)
	if tok != ident ***REMOVED***
		this.errorf("expected identifier: \"%s\", got: \"%s\"", ident, tok)
	***REMOVED***
***REMOVED***

func (this *import_data_parser) read_type() ast.Expr ***REMOVED***
	type_, name := this.read_type_full()
	if name != "" ***REMOVED***
		return ast.NewIdent(name)
	***REMOVED***
	return type_
***REMOVED***

func (this *import_data_parser) read_type_full() (ast.Expr, string) ***REMOVED***
	this.expect('<')
	this.expect_ident("type")

	numstr := this.expect(scanner.Int)
	num, err := strconv.ParseInt(numstr, 10, 32)
	if err != nil ***REMOVED***
		panic(err)
	***REMOVED***

	if this.toktype == '>' ***REMOVED***
		// was already declared previously
		this.next()
		if num < 0 ***REMOVED***
			if num < smallest_builtin_code ***REMOVED***
				this.errorf("out of range built-in type code")
			***REMOVED***
			return builtin_type_names[-num], ""
		***REMOVED*** else ***REMOVED***
			// lookup type table
			type_ := this.typetable[num]
			return type_.type_, type_.name
		***REMOVED***
	***REMOVED***

	this.typetable = append(this.typetable, &import_data_type***REMOVED******REMOVED***)
	var type_ = this.typetable[len(this.typetable)-1]

	switch this.toktype ***REMOVED***
	case scanner.String:
		// named type
		s := this.expect(scanner.String)
		type_.name = s[1 : len(s)-1] // remove ""
		fallthrough
	default:
		// unnamed type
		switch this.toktype ***REMOVED***
		case scanner.Ident:
			switch tok := this.token(); tok ***REMOVED***
			case "struct":
				type_.type_ = this.read_struct_type()
			case "interface":
				type_.type_ = this.read_interface_type()
			case "map":
				type_.type_ = this.read_map_type()
			case "chan":
				type_.type_ = this.read_chan_type()
			default:
				this.errorf("unknown type class token: \"%s\"", tok)
			***REMOVED***
		case '[':
			type_.type_ = this.read_array_or_slice_type()
		case '*':
			this.next()
			if this.token() == "any" ***REMOVED***
				this.next()
				type_.type_ = &ast.StarExpr***REMOVED***X: ast.NewIdent("any")***REMOVED***
			***REMOVED*** else ***REMOVED***
				type_.type_ = &ast.StarExpr***REMOVED***X: this.read_type()***REMOVED***
			***REMOVED***
		case '(':
			type_.type_ = this.read_func_type()
		case '<':
			type_.type_ = this.read_type()
		***REMOVED***
	***REMOVED***

	for this.toktype != '>' ***REMOVED***
		// must be a method or many methods
		this.expect_ident("func")
		this.read_method()
	***REMOVED***

	this.expect('>')
	return type_.type_, type_.name
***REMOVED***

func (this *import_data_parser) read_map_type() ast.Expr ***REMOVED***
	this.expect_ident("map")
	this.expect('[')
	key := this.read_type()
	this.expect(']')
	val := this.read_type()
	return &ast.MapType***REMOVED***Key: key, Value: val***REMOVED***
***REMOVED***

func (this *import_data_parser) read_chan_type() ast.Expr ***REMOVED***
	dir := ast.SEND | ast.RECV
	this.expect_ident("chan")
	switch this.toktype ***REMOVED***
	case '-':
		// chan -< <type>
		this.expect_special("-<")
		dir = ast.SEND
	case '<':
		// slight ambiguity here
		if this.scanner.Peek() == '-' ***REMOVED***
			// chan <- <type>
			this.expect_special("<-")
			dir = ast.RECV
		***REMOVED***
		// chan <type>
	default:
		this.errorf("unexpected token: \"%s\"", this.token())
	***REMOVED***

	return &ast.ChanType***REMOVED***Dir: dir, Value: this.read_type()***REMOVED***
***REMOVED***

func (this *import_data_parser) read_field() *ast.Field ***REMOVED***
	var tag string
	name := this.expect(scanner.Ident)
	type_ := this.read_type()
	if this.toktype == scanner.String ***REMOVED***
		tag = this.expect(scanner.String)
	***REMOVED***

	return &ast.Field***REMOVED***
		Names: []*ast.Ident***REMOVED***ast.NewIdent(name)***REMOVED***,
		Type:  type_,
		Tag:   &ast.BasicLit***REMOVED***Kind: token.STRING, Value: tag***REMOVED***,
	***REMOVED***
***REMOVED***

func (this *import_data_parser) read_struct_type() ast.Expr ***REMOVED***
	var fields []*ast.Field
	read_field := func() ***REMOVED***
		field := this.read_field()
		fields = append(fields, field)
	***REMOVED***

	this.expect_ident("struct")
	this.expect('***REMOVED***')
	for this.toktype != '***REMOVED***' ***REMOVED***
		read_field()
		this.expect(';')
	***REMOVED***
	this.expect('***REMOVED***')
	return &ast.StructType***REMOVED***Fields: &ast.FieldList***REMOVED***List: fields***REMOVED******REMOVED***
***REMOVED***

func (this *import_data_parser) read_parameter() *ast.Field ***REMOVED***
	name := this.expect(scanner.Ident)

	var type_ ast.Expr
	if this.toktype == '.' ***REMOVED***
		this.expect_special("...")
		type_ = &ast.Ellipsis***REMOVED***Elt: this.read_type()***REMOVED***
	***REMOVED*** else ***REMOVED***
		type_ = this.read_type()
	***REMOVED***

	var tag string
	if this.toktype == scanner.String ***REMOVED***
		tag = this.expect(scanner.String)
	***REMOVED***

	return &ast.Field***REMOVED***
		Names: []*ast.Ident***REMOVED***ast.NewIdent(name)***REMOVED***,
		Type:  type_,
		Tag:   &ast.BasicLit***REMOVED***Kind: token.STRING, Value: tag***REMOVED***,
	***REMOVED***
***REMOVED***

func (this *import_data_parser) read_parameters() *ast.FieldList ***REMOVED***
	var fields []*ast.Field
	read_parameter := func() ***REMOVED***
		parameter := this.read_parameter()
		fields = append(fields, parameter)
	***REMOVED***

	this.expect('(')
	if this.toktype != ')' ***REMOVED***
		read_parameter()
		for this.toktype == ',' ***REMOVED***
			this.next() // skip ','
			read_parameter()
		***REMOVED***
	***REMOVED***
	this.expect(')')

	if fields == nil ***REMOVED***
		return nil
	***REMOVED***
	return &ast.FieldList***REMOVED***List: fields***REMOVED***
***REMOVED***

func (this *import_data_parser) read_func_type() *ast.FuncType ***REMOVED***
	var params, results *ast.FieldList

	params = this.read_parameters()
	switch this.toktype ***REMOVED***
	case '<':
		field := &ast.Field***REMOVED***Type: this.read_type()***REMOVED***
		results = &ast.FieldList***REMOVED***List: []*ast.Field***REMOVED***field***REMOVED******REMOVED***
	case '(':
		results = this.read_parameters()
	***REMOVED***

	return &ast.FuncType***REMOVED***Params: params, Results: results***REMOVED***
***REMOVED***

func (this *import_data_parser) read_method_or_embed_spec() *ast.Field ***REMOVED***
	var type_ ast.Expr
	name := this.expect(scanner.Ident)
	if name == "?" ***REMOVED***
		// TODO: ast.SelectorExpr conversion here possibly
		type_ = this.read_type()
	***REMOVED*** else ***REMOVED***
		type_ = this.read_func_type()
	***REMOVED***
	return &ast.Field***REMOVED***
		Names: []*ast.Ident***REMOVED***ast.NewIdent(name)***REMOVED***,
		Type:  type_,
	***REMOVED***
***REMOVED***

func (this *import_data_parser) read_interface_type() ast.Expr ***REMOVED***
	var methods []*ast.Field
	read_method := func() ***REMOVED***
		method := this.read_method_or_embed_spec()
		methods = append(methods, method)
	***REMOVED***

	this.expect_ident("interface")
	this.expect('***REMOVED***')
	for this.toktype != '***REMOVED***' ***REMOVED***
		read_method()
		this.expect(';')
	***REMOVED***
	this.expect('***REMOVED***')
	return &ast.InterfaceType***REMOVED***Methods: &ast.FieldList***REMOVED***List: methods***REMOVED******REMOVED***
***REMOVED***

func (this *import_data_parser) read_method() ***REMOVED***
	var buf1, buf2 bytes.Buffer
	recv := this.read_parameters()
	name := this.expect(scanner.Ident)
	type_ := this.read_func_type()
	this.expect(';')
	pretty_print_type_expr(&buf1, recv.List[0].Type)
	pretty_print_type_expr(&buf2, type_)
	println("func (" + buf1.String() + ") " + name + buf2.String()[4:])
***REMOVED***

func (this *import_data_parser) read_array_or_slice_type() ast.Expr ***REMOVED***
	var length ast.Expr

	this.expect('[')
	if this.toktype == scanner.Int ***REMOVED***
		// array type
		length = &ast.BasicLit***REMOVED***Kind: token.INT, Value: this.expect(scanner.Int)***REMOVED***
	***REMOVED***
	this.expect(']')
	return &ast.ArrayType***REMOVED***
		Len: length,
		Elt: this.read_type(),
	***REMOVED***
***REMOVED***

func (this *import_data_parser) read_const() ***REMOVED***
	var buf bytes.Buffer

	// const keyword was already consumed
	c := "const " + this.expect(scanner.Ident)
	if this.toktype != '=' ***REMOVED***
		// parse type
		type_ := this.read_type()
		pretty_print_type_expr(&buf, type_)
		c += " " + buf.String()
	***REMOVED***

	this.expect('=')

	// parse expr
	this.next()
	this.expect(';')
	println(c)
***REMOVED***

func (this *import_data_parser) read_checksum() ***REMOVED***
	// checksum keyword was already consumed
	for this.toktype != ';' ***REMOVED***
		this.next()
	***REMOVED***
	this.expect(';')
***REMOVED***

func (this *import_data_parser) read_type_decl() ***REMOVED***
	var buf bytes.Buffer
	// type keyword was already consumed
	type_, name := this.read_type_full()
	this.expect(';')
	pretty_print_type_expr(&buf, type_)
	println("type " + name + " " + buf.String())
***REMOVED***

func (this *import_data_parser) read_var() ***REMOVED***
	var buf bytes.Buffer
	// var keyword was already consumed
	name := this.expect(scanner.Ident)
	type_ := this.read_type()
	this.expect(';')
	pretty_print_type_expr(&buf, type_)
	println("var " + name + " " + buf.String())
***REMOVED***

func (this *import_data_parser) read_func() ***REMOVED***
	var buf bytes.Buffer
	// func keyword was already consumed
	name := this.expect(scanner.Ident)
	type_ := this.read_func_type()
	this.expect(';')
	pretty_print_type_expr(&buf, type_)
	println("func " + name + buf.String()[4:])
***REMOVED***

//-------------------------------------------------------------------------
// Pretty printing
//-------------------------------------------------------------------------

func get_array_len(e ast.Expr) string ***REMOVED***
	switch t := e.(type) ***REMOVED***
	case *ast.BasicLit:
		return string(t.Value)
	case *ast.Ellipsis:
		return "..."
	***REMOVED***
	return ""
***REMOVED***

func pretty_print_type_expr(out io.Writer, e ast.Expr) ***REMOVED***
	switch t := e.(type) ***REMOVED***
	case *ast.StarExpr:
		fmt.Fprintf(out, "*")
		pretty_print_type_expr(out, t.X)
	case *ast.Ident:
		if strings.HasPrefix(t.Name, "$") ***REMOVED***
			// beautify anonymous types
			switch t.Name[1] ***REMOVED***
			case 's':
				fmt.Fprintf(out, "struct")
			case 'i':
				fmt.Fprintf(out, "interface")
			***REMOVED***
		***REMOVED*** else ***REMOVED***
			fmt.Fprintf(out, t.Name)
		***REMOVED***
	case *ast.ArrayType:
		al := ""
		if t.Len != nil ***REMOVED***
			println(t.Len)
			al = get_array_len(t.Len)
		***REMOVED***
		if al != "" ***REMOVED***
			fmt.Fprintf(out, "[%s]", al)
		***REMOVED*** else ***REMOVED***
			fmt.Fprintf(out, "[]")
		***REMOVED***
		pretty_print_type_expr(out, t.Elt)
	case *ast.SelectorExpr:
		pretty_print_type_expr(out, t.X)
		fmt.Fprintf(out, ".%s", t.Sel.Name)
	case *ast.FuncType:
		fmt.Fprintf(out, "func(")
		pretty_print_func_field_list(out, t.Params)
		fmt.Fprintf(out, ")")

		buf := bytes.NewBuffer(make([]byte, 0, 256))
		nresults := pretty_print_func_field_list(buf, t.Results)
		if nresults > 0 ***REMOVED***
			results := buf.String()
			if strings.Index(results, ",") != -1 ***REMOVED***
				results = "(" + results + ")"
			***REMOVED***
			fmt.Fprintf(out, " %s", results)
		***REMOVED***
	case *ast.MapType:
		fmt.Fprintf(out, "map[")
		pretty_print_type_expr(out, t.Key)
		fmt.Fprintf(out, "]")
		pretty_print_type_expr(out, t.Value)
	case *ast.InterfaceType:
		fmt.Fprintf(out, "interface***REMOVED******REMOVED***")
	case *ast.Ellipsis:
		fmt.Fprintf(out, "...")
		pretty_print_type_expr(out, t.Elt)
	case *ast.StructType:
		fmt.Fprintf(out, "struct")
	case *ast.ChanType:
		switch t.Dir ***REMOVED***
		case ast.RECV:
			fmt.Fprintf(out, "<-chan ")
		case ast.SEND:
			fmt.Fprintf(out, "chan<- ")
		case ast.SEND | ast.RECV:
			fmt.Fprintf(out, "chan ")
		***REMOVED***
		pretty_print_type_expr(out, t.Value)
	case *ast.ParenExpr:
		fmt.Fprintf(out, "(")
		pretty_print_type_expr(out, t.X)
		fmt.Fprintf(out, ")")
	case *ast.BadExpr:
		// TODO: probably I should check that in a separate function
		// and simply discard declarations with BadExpr as a part of their
		// type
	default:
		// should never happen
		panic("unknown type")
	***REMOVED***
***REMOVED***

func pretty_print_func_field_list(out io.Writer, f *ast.FieldList) int ***REMOVED***
	count := 0
	if f == nil ***REMOVED***
		return count
	***REMOVED***
	for i, field := range f.List ***REMOVED***
		// names
		if field.Names != nil ***REMOVED***
			hasNonblank := false
			for j, name := range field.Names ***REMOVED***
				if name.Name != "?" ***REMOVED***
					hasNonblank = true
					fmt.Fprintf(out, "%s", name.Name)
					if j != len(field.Names)-1 ***REMOVED***
						fmt.Fprintf(out, ", ")
					***REMOVED***
				***REMOVED***
				count++
			***REMOVED***
			if hasNonblank ***REMOVED***
				fmt.Fprintf(out, " ")
			***REMOVED***
		***REMOVED*** else ***REMOVED***
			count++
		***REMOVED***

		// type
		pretty_print_type_expr(out, field.Type)

		// ,
		if i != len(f.List)-1 ***REMOVED***
			fmt.Fprintf(out, ", ")
		***REMOVED***
	***REMOVED***
	return count
***REMOVED***

func main() ***REMOVED***
	data, err := read_import_data("io")
	if err != nil ***REMOVED***
		panic(err)
	***REMOVED***
	parse_import_data(data)
***REMOVED***
