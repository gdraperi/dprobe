package main

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/token"
	"io"
	"reflect"
	"strings"
	"sync"
)

// decl.class
type decl_class int16

const (
	decl_invalid = decl_class(-1 + iota)

	// these are in a sorted order
	decl_const
	decl_func
	decl_import
	decl_package
	decl_type
	decl_var

	// this one serves as a temporary type for those methods that were
	// declared before their actual owner
	decl_methods_stub
)

func (this decl_class) String() string ***REMOVED***
	switch this ***REMOVED***
	case decl_invalid:
		return "PANIC"
	case decl_const:
		return "const"
	case decl_func:
		return "func"
	case decl_import:
		return "import"
	case decl_package:
		return "package"
	case decl_type:
		return "type"
	case decl_var:
		return "var"
	case decl_methods_stub:
		return "IF YOU SEE THIS, REPORT A BUG" // :D
	***REMOVED***
	panic("unreachable")
***REMOVED***

// decl.flags
type decl_flags int16

const (
	decl_foreign decl_flags = 1 << iota // imported from another package

	// means that the decl is a part of the range statement
	// its type is inferred in a special way
	decl_rangevar

	// decl of decl_type class is a type alias
	decl_alias

	// for preventing infinite recursions and loops in type inference code
	decl_visited
)

//-------------------------------------------------------------------------
// decl
//
// The most important data structure of the whole gocode project. It
// describes a single declaration and its children.
//-------------------------------------------------------------------------

type decl struct ***REMOVED***
	// Name starts with '$' if the declaration describes an anonymous type.
	// '$s_%d' for anonymous struct types
	// '$i_%d' for anonymous interface types
	name  string
	typ   ast.Expr
	class decl_class
	flags decl_flags

	// functions for interface type, fields+methods for struct type
	children map[string]*decl

	// embedded types
	embedded []ast.Expr

	// if the type is unknown at AST building time, I'm using these
	value ast.Expr

	// if it's a multiassignment and the Value is a CallExpr, it is being set
	// to an index into the return value tuple, otherwise it's a -1
	value_index int

	// scope where this Decl was declared in (not its visibilty scope!)
	// Decl uses it for type inference
	scope *scope
***REMOVED***

func ast_decl_type(d ast.Decl) ast.Expr ***REMOVED***
	switch t := d.(type) ***REMOVED***
	case *ast.GenDecl:
		switch t.Tok ***REMOVED***
		case token.CONST, token.VAR:
			c := t.Specs[0].(*ast.ValueSpec)
			return c.Type
		case token.TYPE:
			t := t.Specs[0].(*ast.TypeSpec)
			return t.Type
		***REMOVED***
	case *ast.FuncDecl:
		return t.Type
	***REMOVED***
	panic("unreachable")
***REMOVED***

func ast_decl_flags(d ast.Decl) decl_flags ***REMOVED***
	switch t := d.(type) ***REMOVED***
	case *ast.GenDecl:
		switch t.Tok ***REMOVED***
		case token.TYPE:
			if isAliasTypeSpec(t.Specs[0].(*ast.TypeSpec)) ***REMOVED***
				return decl_alias
			***REMOVED***
		***REMOVED***
	***REMOVED***
	return 0
***REMOVED***

func ast_decl_class(d ast.Decl) decl_class ***REMOVED***
	switch t := d.(type) ***REMOVED***
	case *ast.GenDecl:
		switch t.Tok ***REMOVED***
		case token.VAR:
			return decl_var
		case token.CONST:
			return decl_const
		case token.TYPE:
			return decl_type
		***REMOVED***
	case *ast.FuncDecl:
		return decl_func
	***REMOVED***
	panic("unreachable")
***REMOVED***

func ast_decl_convertable(d ast.Decl) bool ***REMOVED***
	switch t := d.(type) ***REMOVED***
	case *ast.GenDecl:
		switch t.Tok ***REMOVED***
		case token.VAR, token.CONST, token.TYPE:
			return true
		***REMOVED***
	case *ast.FuncDecl:
		return true
	***REMOVED***
	return false
***REMOVED***

func ast_field_list_to_decls(f *ast.FieldList, class decl_class, flags decl_flags, scope *scope, add_anonymous bool) map[string]*decl ***REMOVED***
	count := 0
	for _, field := range f.List ***REMOVED***
		count += len(field.Names)
	***REMOVED***

	decls := make(map[string]*decl, count)
	for _, field := range f.List ***REMOVED***
		for _, name := range field.Names ***REMOVED***
			if flags&decl_foreign != 0 && !ast.IsExported(name.Name) ***REMOVED***
				continue
			***REMOVED***
			d := &decl***REMOVED***
				name:        name.Name,
				typ:         field.Type,
				class:       class,
				flags:       flags,
				scope:       scope,
				value_index: -1,
			***REMOVED***
			decls[d.name] = d
		***REMOVED***

		// add anonymous field as a child (type embedding)
		if class == decl_var && field.Names == nil && add_anonymous ***REMOVED***
			tp := get_type_path(field.Type)
			if flags&decl_foreign != 0 && !ast.IsExported(tp.name) ***REMOVED***
				continue
			***REMOVED***
			d := &decl***REMOVED***
				name:        tp.name,
				typ:         field.Type,
				class:       class,
				flags:       flags,
				scope:       scope,
				value_index: -1,
			***REMOVED***
			decls[d.name] = d
		***REMOVED***
	***REMOVED***
	return decls
***REMOVED***

func ast_field_list_to_embedded(f *ast.FieldList) []ast.Expr ***REMOVED***
	count := 0
	for _, field := range f.List ***REMOVED***
		if field.Names == nil || field.Names[0].Name == "?" ***REMOVED***
			count++
		***REMOVED***
	***REMOVED***

	if count == 0 ***REMOVED***
		return nil
	***REMOVED***

	embedded := make([]ast.Expr, count)
	i := 0
	for _, field := range f.List ***REMOVED***
		if field.Names == nil || field.Names[0].Name == "?" ***REMOVED***
			embedded[i] = field.Type
			i++
		***REMOVED***
	***REMOVED***

	return embedded
***REMOVED***

func ast_type_to_embedded(ty ast.Expr) []ast.Expr ***REMOVED***
	switch t := ty.(type) ***REMOVED***
	case *ast.StructType:
		return ast_field_list_to_embedded(t.Fields)
	case *ast.InterfaceType:
		return ast_field_list_to_embedded(t.Methods)
	***REMOVED***
	return nil
***REMOVED***

func ast_type_to_children(ty ast.Expr, flags decl_flags, scope *scope) map[string]*decl ***REMOVED***
	switch t := ty.(type) ***REMOVED***
	case *ast.StructType:
		return ast_field_list_to_decls(t.Fields, decl_var, flags, scope, true)
	case *ast.InterfaceType:
		return ast_field_list_to_decls(t.Methods, decl_func, flags, scope, false)
	***REMOVED***
	return nil
***REMOVED***

//-------------------------------------------------------------------------
// anonymous_id_gen
//
// ID generator for anonymous types (thread-safe)
//-------------------------------------------------------------------------

type anonymous_id_gen struct ***REMOVED***
	sync.Mutex
	i int
***REMOVED***

func (a *anonymous_id_gen) gen() (id int) ***REMOVED***
	a.Lock()
	defer a.Unlock()
	id = a.i
	a.i++
	return
***REMOVED***

var g_anon_gen anonymous_id_gen

//-------------------------------------------------------------------------

func check_for_anon_type(t ast.Expr, flags decl_flags, s *scope) ast.Expr ***REMOVED***
	if t == nil ***REMOVED***
		return nil
	***REMOVED***
	var name string

	switch t.(type) ***REMOVED***
	case *ast.StructType:
		name = fmt.Sprintf("$s_%d", g_anon_gen.gen())
	case *ast.InterfaceType:
		name = fmt.Sprintf("$i_%d", g_anon_gen.gen())
	***REMOVED***

	if name != "" ***REMOVED***
		anonymify_ast(t, flags, s)
		d := new_decl_full(name, decl_type, flags, t, nil, -1, s)
		s.add_named_decl(d)
		return ast.NewIdent(name)
	***REMOVED***
	return t
***REMOVED***

//-------------------------------------------------------------------------

func new_decl_full(name string, class decl_class, flags decl_flags, typ, v ast.Expr, vi int, s *scope) *decl ***REMOVED***
	if name == "_" ***REMOVED***
		return nil
	***REMOVED***
	d := new(decl)
	d.name = name
	d.class = class
	d.flags = flags
	d.typ = typ
	d.value = v
	d.value_index = vi
	d.scope = s
	d.children = ast_type_to_children(d.typ, flags, s)
	d.embedded = ast_type_to_embedded(d.typ)
	return d
***REMOVED***

func new_decl(name string, class decl_class, scope *scope) *decl ***REMOVED***
	decl := new(decl)
	decl.name = name
	decl.class = class
	decl.value_index = -1
	decl.scope = scope
	return decl
***REMOVED***

func new_decl_var(name string, typ ast.Expr, value ast.Expr, vindex int, scope *scope) *decl ***REMOVED***
	if name == "_" ***REMOVED***
		return nil
	***REMOVED***
	decl := new(decl)
	decl.name = name
	decl.class = decl_var
	decl.typ = typ
	decl.value = value
	decl.value_index = vindex
	decl.scope = scope
	return decl
***REMOVED***

func method_of(d ast.Decl) string ***REMOVED***
	if t, ok := d.(*ast.FuncDecl); ok ***REMOVED***
		if t.Recv != nil && len(t.Recv.List) != 0 ***REMOVED***
			switch t := t.Recv.List[0].Type.(type) ***REMOVED***
			case *ast.StarExpr:
				if se, ok := t.X.(*ast.SelectorExpr); ok ***REMOVED***
					return se.Sel.Name
				***REMOVED***
				if ident, ok := t.X.(*ast.Ident); ok ***REMOVED***
					return ident.Name
				***REMOVED***
				return ""
			case *ast.Ident:
				return t.Name
			default:
				return ""
			***REMOVED***
		***REMOVED***
	***REMOVED***
	return ""
***REMOVED***

func (other *decl) deep_copy() *decl ***REMOVED***
	d := new(decl)
	d.name = other.name
	d.class = other.class
	d.flags = other.flags
	d.typ = other.typ
	d.value = other.value
	d.value_index = other.value_index
	d.children = make(map[string]*decl, len(other.children))
	for key, value := range other.children ***REMOVED***
		d.children[key] = value
	***REMOVED***
	if other.embedded != nil ***REMOVED***
		d.embedded = make([]ast.Expr, len(other.embedded))
		copy(d.embedded, other.embedded)
	***REMOVED***
	d.scope = other.scope
	return d
***REMOVED***

func (d *decl) is_rangevar() bool ***REMOVED***
	return d.flags&decl_rangevar != 0
***REMOVED***

func (d *decl) is_alias() bool ***REMOVED***
	return d.flags&decl_alias != 0
***REMOVED***

func (d *decl) is_visited() bool ***REMOVED***
	return d.flags&decl_visited != 0
***REMOVED***

func (d *decl) set_visited() ***REMOVED***
	d.flags |= decl_visited
***REMOVED***

func (d *decl) clear_visited() ***REMOVED***
	d.flags &^= decl_visited
***REMOVED***

func (d *decl) expand_or_replace(other *decl) ***REMOVED***
	// expand only if it's a methods stub, otherwise simply keep it as is
	if d.class != decl_methods_stub && other.class != decl_methods_stub ***REMOVED***
		return
	***REMOVED***

	if d.class == decl_methods_stub ***REMOVED***
		d.typ = other.typ
		d.class = other.class
		d.flags = other.flags
	***REMOVED***

	if other.children != nil ***REMOVED***
		for _, c := range other.children ***REMOVED***
			d.add_child(c)
		***REMOVED***
	***REMOVED***

	if other.embedded != nil ***REMOVED***
		d.embedded = other.embedded
		d.scope = other.scope
	***REMOVED***
***REMOVED***

func (d *decl) matches() bool ***REMOVED***
	if strings.HasPrefix(d.name, "$") || d.class == decl_methods_stub ***REMOVED***
		return false
	***REMOVED***
	return true
***REMOVED***

func (d *decl) pretty_print_type(out io.Writer, canonical_aliases map[string]string) ***REMOVED***
	switch d.class ***REMOVED***
	case decl_type:
		switch d.typ.(type) ***REMOVED***
		case *ast.StructType:
			// TODO: not used due to anonymify?
			fmt.Fprintf(out, "struct")
		case *ast.InterfaceType:
			// TODO: not used due to anonymify?
			fmt.Fprintf(out, "interface")
		default:
			if d.typ != nil ***REMOVED***
				pretty_print_type_expr(out, d.typ, canonical_aliases)
			***REMOVED***
		***REMOVED***
	case decl_var:
		if d.typ != nil ***REMOVED***
			pretty_print_type_expr(out, d.typ, canonical_aliases)
		***REMOVED***
	case decl_func:
		pretty_print_type_expr(out, d.typ, canonical_aliases)
	***REMOVED***
***REMOVED***

func (d *decl) add_child(cd *decl) ***REMOVED***
	if d.children == nil ***REMOVED***
		d.children = make(map[string]*decl)
	***REMOVED***
	d.children[cd.name] = cd
***REMOVED***

func check_for_builtin_funcs(typ *ast.Ident, c *ast.CallExpr, scope *scope) (ast.Expr, *scope) ***REMOVED***
	if strings.HasPrefix(typ.Name, "func(") ***REMOVED***
		if t, ok := c.Fun.(*ast.Ident); ok ***REMOVED***
			switch t.Name ***REMOVED***
			case "new":
				if len(c.Args) > 0 ***REMOVED***
					e := new(ast.StarExpr)
					e.X = c.Args[0]
					return e, scope
				***REMOVED***
			case "make":
				if len(c.Args) > 0 ***REMOVED***
					return c.Args[0], scope
				***REMOVED***
			case "append":
				if len(c.Args) > 0 ***REMOVED***
					t, scope, _ := infer_type(c.Args[0], scope, -1)
					return t, scope
				***REMOVED***
			case "complex":
				// TODO: fix it
				return ast.NewIdent("complex"), g_universe_scope
			case "closed":
				return ast.NewIdent("bool"), g_universe_scope
			case "cap":
				return ast.NewIdent("int"), g_universe_scope
			case "copy":
				return ast.NewIdent("int"), g_universe_scope
			case "len":
				return ast.NewIdent("int"), g_universe_scope
			***REMOVED***
			// TODO:
			// func recover() interface***REMOVED******REMOVED***
			// func imag(c ComplexType) FloatType
			// func real(c ComplexType) FloatType
		***REMOVED***
	***REMOVED***
	return nil, nil
***REMOVED***

func func_return_type(f *ast.FuncType, index int) ast.Expr ***REMOVED***
	if f.Results == nil ***REMOVED***
		return nil
	***REMOVED***

	if index == -1 ***REMOVED***
		return f.Results.List[0].Type
	***REMOVED***

	i := 0
	var field *ast.Field
	for _, field = range f.Results.List ***REMOVED***
		n := 1
		if field.Names != nil ***REMOVED***
			n = len(field.Names)
		***REMOVED***
		if i <= index && index < i+n ***REMOVED***
			return field.Type
		***REMOVED***
		i += n
	***REMOVED***
	return nil
***REMOVED***

type type_path struct ***REMOVED***
	pkg  string
	name string
***REMOVED***

func (tp *type_path) is_nil() bool ***REMOVED***
	return tp.pkg == "" && tp.name == ""
***REMOVED***

// converts type expressions like:
// ast.Expr
// *ast.Expr
// $ast$go/ast.Expr
// to a path that can be used to lookup a type related Decl
func get_type_path(e ast.Expr) (r type_path) ***REMOVED***
	if e == nil ***REMOVED***
		return type_path***REMOVED***"", ""***REMOVED***
	***REMOVED***

	switch t := e.(type) ***REMOVED***
	case *ast.Ident:
		r.name = t.Name
	case *ast.StarExpr:
		r = get_type_path(t.X)
	case *ast.SelectorExpr:
		if ident, ok := t.X.(*ast.Ident); ok ***REMOVED***
			r.pkg = ident.Name
		***REMOVED***
		r.name = t.Sel.Name
	***REMOVED***
	return
***REMOVED***

func lookup_path(tp type_path, scope *scope) *decl ***REMOVED***
	if tp.is_nil() ***REMOVED***
		return nil
	***REMOVED***
	var decl *decl
	if tp.pkg != "" ***REMOVED***
		decl = scope.lookup(tp.pkg)
		// return nil early if the package wasn't found but it's part
		// of the type specification
		if decl == nil ***REMOVED***
			return nil
		***REMOVED***
	***REMOVED***

	if decl != nil ***REMOVED***
		if tp.name != "" ***REMOVED***
			return decl.find_child(tp.name)
		***REMOVED*** else ***REMOVED***
			return decl
		***REMOVED***
	***REMOVED***

	return scope.lookup(tp.name)
***REMOVED***

func lookup_pkg(tp type_path, scope *scope) string ***REMOVED***
	if tp.is_nil() ***REMOVED***
		return ""
	***REMOVED***
	if tp.pkg == "" ***REMOVED***
		return ""
	***REMOVED***
	decl := scope.lookup(tp.pkg)
	if decl == nil ***REMOVED***
		return ""
	***REMOVED***
	return decl.name
***REMOVED***

func type_to_decl(t ast.Expr, scope *scope) *decl ***REMOVED***
	tp := get_type_path(t)
	d := lookup_path(tp, scope)
	if d != nil && d.class == decl_var ***REMOVED***
		// weird variable declaration pointing to itself
		return nil
	***REMOVED***
	return d
***REMOVED***

func expr_to_decl(e ast.Expr, scope *scope) *decl ***REMOVED***
	t, scope, _ := infer_type(e, scope, -1)
	return type_to_decl(t, scope)
***REMOVED***

//-------------------------------------------------------------------------
// Type inference
//-------------------------------------------------------------------------

type type_predicate func(ast.Expr) bool

func advance_to_type(pred type_predicate, v ast.Expr, scope *scope) (ast.Expr, *scope) ***REMOVED***
	if pred(v) ***REMOVED***
		return v, scope
	***REMOVED***

	decl := type_to_decl(v, scope)
	if decl == nil ***REMOVED***
		return nil, nil
	***REMOVED***

	if decl.is_visited() ***REMOVED***
		return nil, nil
	***REMOVED***
	decl.set_visited()
	defer decl.clear_visited()

	return advance_to_type(pred, decl.typ, decl.scope)
***REMOVED***

func advance_to_struct_or_interface(decl *decl) *decl ***REMOVED***
	if decl.is_visited() ***REMOVED***
		return nil
	***REMOVED***
	decl.set_visited()
	defer decl.clear_visited()

	if struct_interface_predicate(decl.typ) ***REMOVED***
		return decl
	***REMOVED***

	decl = type_to_decl(decl.typ, decl.scope)
	if decl == nil ***REMOVED***
		return nil
	***REMOVED***
	return advance_to_struct_or_interface(decl)
***REMOVED***

func struct_interface_predicate(v ast.Expr) bool ***REMOVED***
	switch v.(type) ***REMOVED***
	case *ast.StructType, *ast.InterfaceType:
		return true
	***REMOVED***
	return false
***REMOVED***

func chan_predicate(v ast.Expr) bool ***REMOVED***
	_, ok := v.(*ast.ChanType)
	return ok
***REMOVED***

func index_predicate(v ast.Expr) bool ***REMOVED***
	switch v.(type) ***REMOVED***
	case *ast.ArrayType, *ast.MapType, *ast.Ellipsis:
		return true
	***REMOVED***
	return false
***REMOVED***

func star_predicate(v ast.Expr) bool ***REMOVED***
	_, ok := v.(*ast.StarExpr)
	return ok
***REMOVED***

func func_predicate(v ast.Expr) bool ***REMOVED***
	_, ok := v.(*ast.FuncType)
	return ok
***REMOVED***

func range_predicate(v ast.Expr) bool ***REMOVED***
	switch t := v.(type) ***REMOVED***
	case *ast.Ident:
		if t.Name == "string" ***REMOVED***
			return true
		***REMOVED***
	case *ast.ArrayType, *ast.MapType, *ast.ChanType, *ast.Ellipsis:
		return true
	***REMOVED***
	return false
***REMOVED***

type anonymous_typer struct ***REMOVED***
	flags decl_flags
	scope *scope
***REMOVED***

func (a *anonymous_typer) Visit(node ast.Node) ast.Visitor ***REMOVED***
	switch t := node.(type) ***REMOVED***
	case *ast.CompositeLit:
		t.Type = check_for_anon_type(t.Type, a.flags, a.scope)
	case *ast.MapType:
		t.Key = check_for_anon_type(t.Key, a.flags, a.scope)
		t.Value = check_for_anon_type(t.Value, a.flags, a.scope)
	case *ast.ArrayType:
		t.Elt = check_for_anon_type(t.Elt, a.flags, a.scope)
	case *ast.Ellipsis:
		t.Elt = check_for_anon_type(t.Elt, a.flags, a.scope)
	case *ast.ChanType:
		t.Value = check_for_anon_type(t.Value, a.flags, a.scope)
	case *ast.Field:
		t.Type = check_for_anon_type(t.Type, a.flags, a.scope)
	case *ast.CallExpr:
		t.Fun = check_for_anon_type(t.Fun, a.flags, a.scope)
	case *ast.ParenExpr:
		t.X = check_for_anon_type(t.X, a.flags, a.scope)
	case *ast.StarExpr:
		t.X = check_for_anon_type(t.X, a.flags, a.scope)
	case *ast.GenDecl:
		switch t.Tok ***REMOVED***
		case token.VAR:
			for _, s := range t.Specs ***REMOVED***
				vs := s.(*ast.ValueSpec)
				vs.Type = check_for_anon_type(vs.Type, a.flags, a.scope)
			***REMOVED***
		***REMOVED***
	***REMOVED***
	return a
***REMOVED***

func anonymify_ast(node ast.Node, flags decl_flags, scope *scope) ***REMOVED***
	v := anonymous_typer***REMOVED***flags, scope***REMOVED***
	ast.Walk(&v, node)
***REMOVED***

// RETURNS:
// 	- type expression which represents a full name of a type
//	- bool whether a type expression is actually a type (used internally)
//	- scope in which type makes sense
func infer_type(v ast.Expr, scope *scope, index int) (ast.Expr, *scope, bool) ***REMOVED***
	switch t := v.(type) ***REMOVED***
	case *ast.CompositeLit:
		return t.Type, scope, true
	case *ast.Ident:
		if d := scope.lookup(t.Name); d != nil ***REMOVED***
			if d.class == decl_package ***REMOVED***
				return ast.NewIdent(t.Name), scope, false
			***REMOVED***
			typ, scope := d.infer_type()
			return typ, scope, d.class == decl_type
		***REMOVED***
	case *ast.UnaryExpr:
		switch t.Op ***REMOVED***
		case token.AND:
			// &a makes sense only with values, don't even check for type
			it, s, _ := infer_type(t.X, scope, -1)
			if it == nil ***REMOVED***
				break
			***REMOVED***

			e := new(ast.StarExpr)
			e.X = it
			return e, s, false
		case token.ARROW:
			// <-a makes sense only with values
			it, s, _ := infer_type(t.X, scope, -1)
			if it == nil ***REMOVED***
				break
			***REMOVED***
			switch index ***REMOVED***
			case -1, 0:
				it, s = advance_to_type(chan_predicate, it, s)
				return it.(*ast.ChanType).Value, s, false
			case 1:
				// technically it's a value, but in case of index == 1
				// it is always the last infer operation
				return ast.NewIdent("bool"), g_universe_scope, false
			***REMOVED***
		case token.ADD, token.NOT, token.SUB, token.XOR:
			it, s, _ := infer_type(t.X, scope, -1)
			if it == nil ***REMOVED***
				break
			***REMOVED***
			return it, s, false
		***REMOVED***
	case *ast.BinaryExpr:
		switch t.Op ***REMOVED***
		case token.EQL, token.NEQ, token.LSS, token.LEQ,
			token.GTR, token.GEQ, token.LOR, token.LAND:
			// logic operations, the result is a bool, always
			return ast.NewIdent("bool"), g_universe_scope, false
		case token.ADD, token.SUB, token.MUL, token.QUO, token.OR,
			token.XOR, token.REM, token.AND, token.AND_NOT:
			// try X, then Y, they should be the same anyway
			it, s, _ := infer_type(t.X, scope, -1)
			if it == nil ***REMOVED***
				it, s, _ = infer_type(t.Y, scope, -1)
				if it == nil ***REMOVED***
					break
				***REMOVED***
			***REMOVED***
			return it, s, false
		case token.SHL, token.SHR:
			// try only X for shifts, Y is always uint
			it, s, _ := infer_type(t.X, scope, -1)
			if it == nil ***REMOVED***
				break
			***REMOVED***
			return it, s, false
		***REMOVED***
	case *ast.IndexExpr:
		// something[another] always returns a value and it works on a value too
		it, s, _ := infer_type(t.X, scope, -1)
		if it == nil ***REMOVED***
			break
		***REMOVED***
		it, s = advance_to_type(index_predicate, it, s)
		switch t := it.(type) ***REMOVED***
		case *ast.ArrayType:
			return t.Elt, s, false
		case *ast.Ellipsis:
			return t.Elt, s, false
		case *ast.MapType:
			switch index ***REMOVED***
			case -1, 0:
				return t.Value, s, false
			case 1:
				return ast.NewIdent("bool"), g_universe_scope, false
			***REMOVED***
		***REMOVED***
	case *ast.SliceExpr:
		// something[start : end] always returns a value
		it, s, _ := infer_type(t.X, scope, -1)
		if it == nil ***REMOVED***
			break
		***REMOVED***
		it, s = advance_to_type(index_predicate, it, s)
		switch t := it.(type) ***REMOVED***
		case *ast.ArrayType:
			e := new(ast.ArrayType)
			e.Elt = t.Elt
			return e, s, false
		***REMOVED***
	case *ast.StarExpr:
		it, s, is_type := infer_type(t.X, scope, -1)
		if it == nil ***REMOVED***
			break
		***REMOVED***
		if is_type ***REMOVED***
			// if it's a type, add * modifier, make it a 'pointer of' type
			e := new(ast.StarExpr)
			e.X = it
			return e, s, true
		***REMOVED*** else ***REMOVED***
			it, s := advance_to_type(star_predicate, it, s)
			if se, ok := it.(*ast.StarExpr); ok ***REMOVED***
				return se.X, s, false
			***REMOVED***
		***REMOVED***
	case *ast.CallExpr:
		// this is a function call or a type cast:
		// myFunc(1,2,3) or int16(myvar)
		it, s, is_type := infer_type(t.Fun, scope, -1)
		if it == nil ***REMOVED***
			break
		***REMOVED***

		if is_type ***REMOVED***
			// a type cast
			return it, scope, false
		***REMOVED*** else ***REMOVED***
			// it must be a function call or a built-in function
			// first check for built-in
			if ct, ok := it.(*ast.Ident); ok ***REMOVED***
				ty, s := check_for_builtin_funcs(ct, t, scope)
				if ty != nil ***REMOVED***
					return ty, s, false
				***REMOVED***
			***REMOVED***

			// then check for an ordinary function call
			it, scope = advance_to_type(func_predicate, it, s)
			if ct, ok := it.(*ast.FuncType); ok ***REMOVED***
				return func_return_type(ct, index), s, false
			***REMOVED***
		***REMOVED***
	case *ast.ParenExpr:
		it, s, is_type := infer_type(t.X, scope, -1)
		if it == nil ***REMOVED***
			break
		***REMOVED***
		return it, s, is_type
	case *ast.SelectorExpr:
		it, s, _ := infer_type(t.X, scope, -1)
		if it == nil ***REMOVED***
			break
		***REMOVED***

		if d := type_to_decl(it, s); d != nil ***REMOVED***
			c := d.find_child_and_in_embedded(t.Sel.Name)
			if c != nil ***REMOVED***
				if c.class == decl_type ***REMOVED***
					return t, scope, true
				***REMOVED*** else ***REMOVED***
					typ, s := c.infer_type()
					return typ, s, false
				***REMOVED***
			***REMOVED***
		***REMOVED***
	case *ast.FuncLit:
		// it's a value, but I think most likely we don't even care, cause we can only
		// call it, and CallExpr uses the type itself to figure out
		return t.Type, scope, false
	case *ast.TypeAssertExpr:
		if t.Type == nil ***REMOVED***
			return infer_type(t.X, scope, -1)
		***REMOVED***
		switch index ***REMOVED***
		case -1, 0:
			// converting a value to a different type, but return thing is a value
			it, _, _ := infer_type(t.Type, scope, -1)
			return it, scope, false
		case 1:
			return ast.NewIdent("bool"), g_universe_scope, false
		***REMOVED***
	case *ast.ArrayType, *ast.MapType, *ast.ChanType, *ast.Ellipsis,
		*ast.FuncType, *ast.StructType, *ast.InterfaceType:
		return t, scope, true
	default:
		_ = reflect.TypeOf(v)
		//fmt.Println(ty)
	***REMOVED***
	return nil, nil, false
***REMOVED***

// Uses Value, ValueIndex and Scope to infer the type of this
// declaration. Returns the type itself and the scope where this type
// makes sense.
func (d *decl) infer_type() (ast.Expr, *scope) ***REMOVED***
	// special case for range vars
	if d.is_rangevar() ***REMOVED***
		var scope *scope
		d.typ, scope = infer_range_type(d.value, d.scope, d.value_index)
		return d.typ, scope
	***REMOVED***

	switch d.class ***REMOVED***
	case decl_package:
		// package is handled specially in inferType
		return nil, nil
	case decl_type:
		return ast.NewIdent(d.name), d.scope
	***REMOVED***

	// shortcut
	if d.typ != nil && d.value == nil ***REMOVED***
		return d.typ, d.scope
	***REMOVED***

	// prevent loops
	if d.is_visited() ***REMOVED***
		return nil, nil
	***REMOVED***
	d.set_visited()
	defer d.clear_visited()

	var scope *scope
	d.typ, scope, _ = infer_type(d.value, d.scope, d.value_index)
	return d.typ, scope
***REMOVED***

func (d *decl) type_dealias() *decl ***REMOVED***
	if d.is_visited() ***REMOVED***
		return nil
	***REMOVED***
	d.set_visited()
	defer d.clear_visited()

	dd := type_to_decl(d.typ, d.scope)
	if dd != nil && dd.is_alias() ***REMOVED***
		return dd.type_dealias()
	***REMOVED***
	return dd
***REMOVED***

func (d *decl) find_child(name string) *decl ***REMOVED***
	// type aliases don't really have any children on their own, but they
	// point to a different type, let's try to find one
	if d.is_alias() ***REMOVED***
		dd := d.type_dealias()
		if dd != nil ***REMOVED***
			return dd.find_child(name)
		***REMOVED***

		// note that type alias can also point to a type literal, something like
		// type A = struct ***REMOVED*** A int ***REMOVED***
		// in this case we rely on "advance_to_struct_or_interface" below
	***REMOVED***

	if d.children != nil ***REMOVED***
		if c, ok := d.children[name]; ok ***REMOVED***
			return c
		***REMOVED***
	***REMOVED***

	decl := advance_to_struct_or_interface(d)
	if decl != nil && decl != d ***REMOVED***
		if d.is_visited() ***REMOVED***
			return nil
		***REMOVED***
		d.set_visited()
		defer d.clear_visited()

		return decl.find_child(name)
	***REMOVED***
	return nil
***REMOVED***

func (d *decl) find_child_and_in_embedded(name string) *decl ***REMOVED***
	if d == nil ***REMOVED***
		return nil
	***REMOVED***

	c := d.find_child(name)
	if c == nil ***REMOVED***
		for _, e := range d.embedded ***REMOVED***
			typedecl := type_to_decl(e, d.scope)
			c = typedecl.find_child_and_in_embedded(name)
			if c != nil ***REMOVED***
				break
			***REMOVED***
		***REMOVED***
	***REMOVED***
	return c
***REMOVED***

// Special type inference for range statements.
// [int], [int] := range [string]
// [int], [value] := range [slice or array]
// [key], [value] := range [map]
// [value], [nil] := range [chan]
func infer_range_type(e ast.Expr, sc *scope, valueindex int) (ast.Expr, *scope) ***REMOVED***
	t, s, _ := infer_type(e, sc, -1)
	t, s = advance_to_type(range_predicate, t, s)
	if t != nil ***REMOVED***
		var t1, t2 ast.Expr
		var s1, s2 *scope
		s1 = s
		s2 = s

		switch t := t.(type) ***REMOVED***
		case *ast.Ident:
			// string
			if t.Name == "string" ***REMOVED***
				t1 = ast.NewIdent("int")
				t2 = ast.NewIdent("rune")
				s1 = g_universe_scope
				s2 = g_universe_scope
			***REMOVED*** else ***REMOVED***
				t1, t2 = nil, nil
			***REMOVED***
		case *ast.ArrayType:
			t1 = ast.NewIdent("int")
			s1 = g_universe_scope
			t2 = t.Elt
		case *ast.Ellipsis:
			t1 = ast.NewIdent("int")
			s1 = g_universe_scope
			t2 = t.Elt
		case *ast.MapType:
			t1 = t.Key
			t2 = t.Value
		case *ast.ChanType:
			t1 = t.Value
			t2 = nil
		default:
			t1, t2 = nil, nil
		***REMOVED***

		switch valueindex ***REMOVED***
		case 0:
			return t1, s1
		case 1:
			return t2, s2
		***REMOVED***
	***REMOVED***
	return nil, nil
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

func pretty_print_type_expr(out io.Writer, e ast.Expr, canonical_aliases map[string]string) ***REMOVED***
	switch t := e.(type) ***REMOVED***
	case *ast.StarExpr:
		fmt.Fprintf(out, "*")
		pretty_print_type_expr(out, t.X, canonical_aliases)
	case *ast.Ident:
		if strings.HasPrefix(t.Name, "$") ***REMOVED***
			// beautify anonymous types
			switch t.Name[1] ***REMOVED***
			case 's':
				fmt.Fprintf(out, "struct")
			case 'i':
				// ok, in most cases anonymous interface is an
				// empty interface, I'll just pretend that
				// it's always true
				fmt.Fprintf(out, "interface***REMOVED******REMOVED***")
			***REMOVED***
		***REMOVED*** else if !*g_debug && strings.HasPrefix(t.Name, "!") ***REMOVED***
			// these are full package names for disambiguating and pretty
			// printing packages within packages, e.g.
			// !go/ast!ast vs. !github.com/nsf/my/ast!ast
			// another ugly hack, if people are punished in hell for ugly hacks
			// I'm screwed...
			emarkIdx := strings.LastIndex(t.Name, "!")
			path := t.Name[1:emarkIdx]
			alias := canonical_aliases[path]
			if alias == "" ***REMOVED***
				alias = t.Name[emarkIdx+1:]
			***REMOVED***
			fmt.Fprintf(out, alias)
		***REMOVED*** else ***REMOVED***
			fmt.Fprintf(out, t.Name)
		***REMOVED***
	case *ast.ArrayType:
		al := ""
		if t.Len != nil ***REMOVED***
			al = get_array_len(t.Len)
		***REMOVED***
		if al != "" ***REMOVED***
			fmt.Fprintf(out, "[%s]", al)
		***REMOVED*** else ***REMOVED***
			fmt.Fprintf(out, "[]")
		***REMOVED***
		pretty_print_type_expr(out, t.Elt, canonical_aliases)
	case *ast.SelectorExpr:
		pretty_print_type_expr(out, t.X, canonical_aliases)
		fmt.Fprintf(out, ".%s", t.Sel.Name)
	case *ast.FuncType:
		fmt.Fprintf(out, "func(")
		pretty_print_func_field_list(out, t.Params, canonical_aliases)
		fmt.Fprintf(out, ")")

		buf := bytes.NewBuffer(make([]byte, 0, 256))
		nresults := pretty_print_func_field_list(buf, t.Results, canonical_aliases)
		if nresults > 0 ***REMOVED***
			results := buf.String()
			if strings.IndexAny(results, ", ") != -1 ***REMOVED***
				results = "(" + results + ")"
			***REMOVED***
			fmt.Fprintf(out, " %s", results)
		***REMOVED***
	case *ast.MapType:
		fmt.Fprintf(out, "map[")
		pretty_print_type_expr(out, t.Key, canonical_aliases)
		fmt.Fprintf(out, "]")
		pretty_print_type_expr(out, t.Value, canonical_aliases)
	case *ast.InterfaceType:
		fmt.Fprintf(out, "interface***REMOVED******REMOVED***")
	case *ast.Ellipsis:
		fmt.Fprintf(out, "...")
		pretty_print_type_expr(out, t.Elt, canonical_aliases)
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
		pretty_print_type_expr(out, t.Value, canonical_aliases)
	case *ast.ParenExpr:
		fmt.Fprintf(out, "(")
		pretty_print_type_expr(out, t.X, canonical_aliases)
		fmt.Fprintf(out, ")")
	case *ast.BadExpr:
		// TODO: probably I should check that in a separate function
		// and simply discard declarations with BadExpr as a part of their
		// type
	default:
		// the element has some weird type, just ignore it
	***REMOVED***
***REMOVED***

func pretty_print_func_field_list(out io.Writer, f *ast.FieldList, canonical_aliases map[string]string) int ***REMOVED***
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
		pretty_print_type_expr(out, field.Type, canonical_aliases)

		// ,
		if i != len(f.List)-1 ***REMOVED***
			fmt.Fprintf(out, ", ")
		***REMOVED***
	***REMOVED***
	return count
***REMOVED***

func ast_decl_names(d ast.Decl) []*ast.Ident ***REMOVED***
	var names []*ast.Ident

	switch t := d.(type) ***REMOVED***
	case *ast.GenDecl:
		switch t.Tok ***REMOVED***
		case token.CONST:
			c := t.Specs[0].(*ast.ValueSpec)
			names = make([]*ast.Ident, len(c.Names))
			for i, name := range c.Names ***REMOVED***
				names[i] = name
			***REMOVED***
		case token.TYPE:
			t := t.Specs[0].(*ast.TypeSpec)
			names = make([]*ast.Ident, 1)
			names[0] = t.Name
		case token.VAR:
			v := t.Specs[0].(*ast.ValueSpec)
			names = make([]*ast.Ident, len(v.Names))
			for i, name := range v.Names ***REMOVED***
				names[i] = name
			***REMOVED***
		***REMOVED***
	case *ast.FuncDecl:
		names = make([]*ast.Ident, 1)
		names[0] = t.Name
	***REMOVED***

	return names
***REMOVED***

func ast_decl_values(d ast.Decl) []ast.Expr ***REMOVED***
	// TODO: CONST values here too
	switch t := d.(type) ***REMOVED***
	case *ast.GenDecl:
		switch t.Tok ***REMOVED***
		case token.VAR:
			v := t.Specs[0].(*ast.ValueSpec)
			if v.Values != nil ***REMOVED***
				return v.Values
			***REMOVED***
		***REMOVED***
	***REMOVED***
	return nil
***REMOVED***

func ast_decl_split(d ast.Decl) []ast.Decl ***REMOVED***
	var decls []ast.Decl
	if t, ok := d.(*ast.GenDecl); ok ***REMOVED***
		decls = make([]ast.Decl, len(t.Specs))
		for i, s := range t.Specs ***REMOVED***
			decl := new(ast.GenDecl)
			*decl = *t
			decl.Specs = make([]ast.Spec, 1)
			decl.Specs[0] = s
			decls[i] = decl
		***REMOVED***
	***REMOVED*** else ***REMOVED***
		decls = make([]ast.Decl, 1)
		decls[0] = d
	***REMOVED***
	return decls
***REMOVED***

//-------------------------------------------------------------------------
// decl_pack
//-------------------------------------------------------------------------

type decl_pack struct ***REMOVED***
	names  []*ast.Ident
	typ    ast.Expr
	values []ast.Expr
***REMOVED***

type foreach_decl_struct struct ***REMOVED***
	decl_pack
	decl ast.Decl
***REMOVED***

func (f *decl_pack) value(i int) ast.Expr ***REMOVED***
	if f.values == nil ***REMOVED***
		return nil
	***REMOVED***
	if len(f.values) > 1 ***REMOVED***
		return f.values[i]
	***REMOVED***
	return f.values[0]
***REMOVED***

func (f *decl_pack) value_index(i int) (v ast.Expr, vi int) ***REMOVED***
	// default: nil value
	v = nil
	vi = -1

	if f.values != nil ***REMOVED***
		// A = B, if there is only one name, the value is solo too
		if len(f.names) == 1 ***REMOVED***
			return f.values[0], -1
		***REMOVED***

		if len(f.values) > 1 ***REMOVED***
			// in case if there are multiple values, it's a usual
			// multiassignment
			if i >= len(f.values) ***REMOVED***
				i = len(f.values) - 1
			***REMOVED***
			v = f.values[i]
		***REMOVED*** else ***REMOVED***
			// in case if there is one value, but many names, it's
			// a tuple unpack.. use index here
			v = f.values[0]
			vi = i
		***REMOVED***
	***REMOVED***
	return
***REMOVED***

func (f *decl_pack) type_value_index(i int) (ast.Expr, ast.Expr, int) ***REMOVED***
	if f.typ != nil ***REMOVED***
		// If there is a type, we don't care about value, just return the type
		// and zero value.
		return f.typ, nil, -1
	***REMOVED***

	// And otherwise we simply return nil type and a valid value for later inferring.
	v, vi := f.value_index(i)
	return nil, v, vi
***REMOVED***

type foreach_decl_func func(data *foreach_decl_struct)

func foreach_decl(decl ast.Decl, do foreach_decl_func) ***REMOVED***
	decls := ast_decl_split(decl)
	var data foreach_decl_struct
	for _, decl := range decls ***REMOVED***
		if !ast_decl_convertable(decl) ***REMOVED***
			continue
		***REMOVED***
		data.names = ast_decl_names(decl)
		data.typ = ast_decl_type(decl)
		data.values = ast_decl_values(decl)
		data.decl = decl

		do(&data)
	***REMOVED***
***REMOVED***

//-------------------------------------------------------------------------
// Built-in declarations
//-------------------------------------------------------------------------

var g_universe_scope = new_scope(nil)

func init() ***REMOVED***
	builtin := ast.NewIdent("built-in")

	add_type := func(name string) ***REMOVED***
		d := new_decl(name, decl_type, g_universe_scope)
		d.typ = builtin
		g_universe_scope.add_named_decl(d)
	***REMOVED***
	add_type("bool")
	add_type("byte")
	add_type("complex64")
	add_type("complex128")
	add_type("float32")
	add_type("float64")
	add_type("int8")
	add_type("int16")
	add_type("int32")
	add_type("int64")
	add_type("string")
	add_type("uint8")
	add_type("uint16")
	add_type("uint32")
	add_type("uint64")
	add_type("int")
	add_type("uint")
	add_type("uintptr")
	add_type("rune")

	add_const := func(name string) ***REMOVED***
		d := new_decl(name, decl_const, g_universe_scope)
		d.typ = builtin
		g_universe_scope.add_named_decl(d)
	***REMOVED***
	add_const("true")
	add_const("false")
	add_const("iota")
	add_const("nil")

	add_func := func(name, typ string) ***REMOVED***
		d := new_decl(name, decl_func, g_universe_scope)
		d.typ = ast.NewIdent(typ)
		g_universe_scope.add_named_decl(d)
	***REMOVED***
	add_func("append", "func([]type, ...type) []type")
	add_func("cap", "func(container) int")
	add_func("close", "func(channel)")
	add_func("complex", "func(real, imag) complex")
	add_func("copy", "func(dst, src)")
	add_func("delete", "func(map[typeA]typeB, typeA)")
	add_func("imag", "func(complex)")
	add_func("len", "func(container) int")
	add_func("make", "func(type, len[, cap]) type")
	add_func("new", "func(type) *type")
	add_func("panic", "func(interface***REMOVED******REMOVED***)")
	add_func("print", "func(...interface***REMOVED******REMOVED***)")
	add_func("println", "func(...interface***REMOVED******REMOVED***)")
	add_func("real", "func(complex)")
	add_func("recover", "func() interface***REMOVED******REMOVED***")

	// built-in error interface
	d := new_decl("error", decl_type, g_universe_scope)
	d.typ = &ast.InterfaceType***REMOVED******REMOVED***
	d.children = make(map[string]*decl)
	d.children["Error"] = new_decl("Error", decl_func, g_universe_scope)
	d.children["Error"].typ = &ast.FuncType***REMOVED***
		Results: &ast.FieldList***REMOVED***
			List: []*ast.Field***REMOVED***
				***REMOVED***
					Type: ast.NewIdent("string"),
				***REMOVED***,
			***REMOVED***,
		***REMOVED***,
	***REMOVED***
	g_universe_scope.add_named_decl(d)
***REMOVED***
