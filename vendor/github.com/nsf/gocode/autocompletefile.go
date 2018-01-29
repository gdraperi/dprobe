package main

import (
	"bytes"
	"go/ast"
	"go/parser"
	"go/scanner"
	"go/token"
	"log"
)

func parse_decl_list(fset *token.FileSet, data []byte) ([]ast.Decl, error) ***REMOVED***
	var buf bytes.Buffer
	buf.WriteString("package p;")
	buf.Write(data)
	file, err := parser.ParseFile(fset, "", buf.Bytes(), parser.AllErrors)
	if err != nil ***REMOVED***
		return file.Decls, err
	***REMOVED***
	return file.Decls, nil
***REMOVED***

func log_parse_error(intro string, err error) ***REMOVED***
	if el, ok := err.(scanner.ErrorList); ok ***REMOVED***
		log.Printf("%s:", intro)
		for _, er := range el ***REMOVED***
			log.Printf(" %s", er)
		***REMOVED***
	***REMOVED*** else ***REMOVED***
		log.Printf("%s: %s", intro, err)
	***REMOVED***
***REMOVED***

//-------------------------------------------------------------------------
// auto_complete_file
//-------------------------------------------------------------------------

type auto_complete_file struct ***REMOVED***
	name         string
	package_name string

	decls     map[string]*decl
	packages  []package_import
	filescope *scope
	scope     *scope

	cursor  int // for current file buffer only
	fset    *token.FileSet
	context *package_lookup_context
***REMOVED***

func new_auto_complete_file(name string, context *package_lookup_context) *auto_complete_file ***REMOVED***
	p := new(auto_complete_file)
	p.name = name
	p.cursor = -1
	p.fset = token.NewFileSet()
	p.context = context
	return p
***REMOVED***

func (f *auto_complete_file) offset(p token.Pos) int ***REMOVED***
	const fixlen = len("package p;")
	return f.fset.Position(p).Offset - fixlen
***REMOVED***

// this one is used for current file buffer exclusively
func (f *auto_complete_file) process_data(data []byte) ***REMOVED***
	cur, filedata, block := rip_off_decl(data, f.cursor)
	file, err := parser.ParseFile(f.fset, "", filedata, parser.AllErrors)
	if err != nil && *g_debug ***REMOVED***
		log_parse_error("Error parsing input file (outer block)", err)
	***REMOVED***
	f.package_name = package_name(file)

	f.decls = make(map[string]*decl)
	f.packages = collect_package_imports(f.name, file.Decls, f.context)
	f.filescope = new_scope(nil)
	f.scope = f.filescope

	for _, d := range file.Decls ***REMOVED***
		anonymify_ast(d, 0, f.filescope)
	***REMOVED***

	// process all top-level declarations
	for _, decl := range file.Decls ***REMOVED***
		append_to_top_decls(f.decls, decl, f.scope)
	***REMOVED***
	if block != nil ***REMOVED***
		// process local function as top-level declaration
		decls, err := parse_decl_list(f.fset, block)
		if err != nil && *g_debug ***REMOVED***
			log_parse_error("Error parsing input file (inner block)", err)
		***REMOVED***

		for _, d := range decls ***REMOVED***
			anonymify_ast(d, 0, f.filescope)
		***REMOVED***

		for _, decl := range decls ***REMOVED***
			append_to_top_decls(f.decls, decl, f.scope)
		***REMOVED***

		// process function internals
		f.cursor = cur
		for _, decl := range decls ***REMOVED***
			f.process_decl_locals(decl)
		***REMOVED***
	***REMOVED***

***REMOVED***

func (f *auto_complete_file) process_decl_locals(decl ast.Decl) ***REMOVED***
	switch t := decl.(type) ***REMOVED***
	case *ast.FuncDecl:
		if f.cursor_in(t.Body) ***REMOVED***
			s := f.scope
			f.scope = new_scope(f.scope)

			f.process_field_list(t.Recv, s)
			f.process_field_list(t.Type.Params, s)
			f.process_field_list(t.Type.Results, s)
			f.process_block_stmt(t.Body)
		***REMOVED***
	default:
		v := new(func_lit_visitor)
		v.ctx = f
		ast.Walk(v, decl)
	***REMOVED***
***REMOVED***

func (f *auto_complete_file) process_decl(decl ast.Decl) ***REMOVED***
	if t, ok := decl.(*ast.GenDecl); ok && f.offset(t.TokPos) > f.cursor ***REMOVED***
		return
	***REMOVED***
	prevscope := f.scope
	foreach_decl(decl, func(data *foreach_decl_struct) ***REMOVED***
		class := ast_decl_class(data.decl)
		if class != decl_type ***REMOVED***
			f.scope, prevscope = advance_scope(f.scope)
		***REMOVED***
		for i, name := range data.names ***REMOVED***
			typ, v, vi := data.type_value_index(i)

			d := new_decl_full(name.Name, class, ast_decl_flags(data.decl), typ, v, vi, prevscope)
			if d == nil ***REMOVED***
				return
			***REMOVED***

			f.scope.add_named_decl(d)
		***REMOVED***
	***REMOVED***)
***REMOVED***

func (f *auto_complete_file) process_block_stmt(block *ast.BlockStmt) ***REMOVED***
	if block != nil && f.cursor_in(block) ***REMOVED***
		f.scope, _ = advance_scope(f.scope)

		for _, stmt := range block.List ***REMOVED***
			f.process_stmt(stmt)
		***REMOVED***

		// hack to process all func literals
		v := new(func_lit_visitor)
		v.ctx = f
		ast.Walk(v, block)
	***REMOVED***
***REMOVED***

type func_lit_visitor struct ***REMOVED***
	ctx *auto_complete_file
***REMOVED***

func (v *func_lit_visitor) Visit(node ast.Node) ast.Visitor ***REMOVED***
	if t, ok := node.(*ast.FuncLit); ok && v.ctx.cursor_in(t.Body) ***REMOVED***
		s := v.ctx.scope
		v.ctx.scope = new_scope(v.ctx.scope)

		v.ctx.process_field_list(t.Type.Params, s)
		v.ctx.process_field_list(t.Type.Results, s)
		v.ctx.process_block_stmt(t.Body)

		return nil
	***REMOVED***
	return v
***REMOVED***

func (f *auto_complete_file) process_stmt(stmt ast.Stmt) ***REMOVED***
	switch t := stmt.(type) ***REMOVED***
	case *ast.DeclStmt:
		f.process_decl(t.Decl)
	case *ast.AssignStmt:
		f.process_assign_stmt(t)
	case *ast.IfStmt:
		if f.cursor_in_if_head(t) ***REMOVED***
			f.process_stmt(t.Init)
		***REMOVED*** else if f.cursor_in_if_stmt(t) ***REMOVED***
			f.scope, _ = advance_scope(f.scope)
			f.process_stmt(t.Init)
			f.process_block_stmt(t.Body)
			f.process_stmt(t.Else)
		***REMOVED***
	case *ast.BlockStmt:
		f.process_block_stmt(t)
	case *ast.RangeStmt:
		f.process_range_stmt(t)
	case *ast.ForStmt:
		if f.cursor_in_for_head(t) ***REMOVED***
			f.process_stmt(t.Init)
		***REMOVED*** else if f.cursor_in(t.Body) ***REMOVED***
			f.scope, _ = advance_scope(f.scope)

			f.process_stmt(t.Init)
			f.process_block_stmt(t.Body)
		***REMOVED***
	case *ast.SwitchStmt:
		f.process_switch_stmt(t)
	case *ast.TypeSwitchStmt:
		f.process_type_switch_stmt(t)
	case *ast.SelectStmt:
		f.process_select_stmt(t)
	case *ast.LabeledStmt:
		f.process_stmt(t.Stmt)
	***REMOVED***
***REMOVED***

func (f *auto_complete_file) process_select_stmt(a *ast.SelectStmt) ***REMOVED***
	if !f.cursor_in(a.Body) ***REMOVED***
		return
	***REMOVED***
	var prevscope *scope
	f.scope, prevscope = advance_scope(f.scope)

	var last_cursor_after *ast.CommClause
	for _, s := range a.Body.List ***REMOVED***
		if cc := s.(*ast.CommClause); f.cursor > f.offset(cc.Colon) ***REMOVED***
			last_cursor_after = cc
		***REMOVED***
	***REMOVED***

	if last_cursor_after != nil ***REMOVED***
		if last_cursor_after.Comm != nil ***REMOVED***
			//if lastCursorAfter.Lhs != nil && lastCursorAfter.Tok == token.DEFINE ***REMOVED***
			if astmt, ok := last_cursor_after.Comm.(*ast.AssignStmt); ok && astmt.Tok == token.DEFINE ***REMOVED***
				vname := astmt.Lhs[0].(*ast.Ident).Name
				v := new_decl_var(vname, nil, astmt.Rhs[0], -1, prevscope)
				if v != nil ***REMOVED***
					f.scope.add_named_decl(v)
				***REMOVED***
			***REMOVED***
		***REMOVED***
		for _, s := range last_cursor_after.Body ***REMOVED***
			f.process_stmt(s)
		***REMOVED***
	***REMOVED***
***REMOVED***

func (f *auto_complete_file) process_type_switch_stmt(a *ast.TypeSwitchStmt) ***REMOVED***
	if !f.cursor_in(a.Body) ***REMOVED***
		return
	***REMOVED***
	var prevscope *scope
	f.scope, prevscope = advance_scope(f.scope)

	f.process_stmt(a.Init)
	// type var
	var tv *decl
	if a, ok := a.Assign.(*ast.AssignStmt); ok ***REMOVED***
		lhs := a.Lhs
		rhs := a.Rhs
		if lhs != nil && len(lhs) == 1 ***REMOVED***
			tvname := lhs[0].(*ast.Ident).Name
			tv = new_decl_var(tvname, nil, rhs[0], -1, prevscope)
		***REMOVED***
	***REMOVED***

	var last_cursor_after *ast.CaseClause
	for _, s := range a.Body.List ***REMOVED***
		if cc := s.(*ast.CaseClause); f.cursor > f.offset(cc.Colon) ***REMOVED***
			last_cursor_after = cc
		***REMOVED***
	***REMOVED***

	if last_cursor_after != nil ***REMOVED***
		if tv != nil ***REMOVED***
			if last_cursor_after.List != nil && len(last_cursor_after.List) == 1 ***REMOVED***
				tv.typ = last_cursor_after.List[0]
				tv.value = nil
			***REMOVED***
			f.scope.add_named_decl(tv)
		***REMOVED***
		for _, s := range last_cursor_after.Body ***REMOVED***
			f.process_stmt(s)
		***REMOVED***
	***REMOVED***
***REMOVED***

func (f *auto_complete_file) process_switch_stmt(a *ast.SwitchStmt) ***REMOVED***
	if !f.cursor_in(a.Body) ***REMOVED***
		return
	***REMOVED***
	f.scope, _ = advance_scope(f.scope)

	f.process_stmt(a.Init)
	var last_cursor_after *ast.CaseClause
	for _, s := range a.Body.List ***REMOVED***
		if cc := s.(*ast.CaseClause); f.cursor > f.offset(cc.Colon) ***REMOVED***
			last_cursor_after = cc
		***REMOVED***
	***REMOVED***
	if last_cursor_after != nil ***REMOVED***
		for _, s := range last_cursor_after.Body ***REMOVED***
			f.process_stmt(s)
		***REMOVED***
	***REMOVED***
***REMOVED***

func (f *auto_complete_file) process_range_stmt(a *ast.RangeStmt) ***REMOVED***
	if !f.cursor_in(a.Body) ***REMOVED***
		return
	***REMOVED***
	var prevscope *scope
	f.scope, prevscope = advance_scope(f.scope)

	if a.Tok == token.DEFINE ***REMOVED***
		if t, ok := a.Key.(*ast.Ident); ok ***REMOVED***
			d := new_decl_var(t.Name, nil, a.X, 0, prevscope)
			if d != nil ***REMOVED***
				d.flags |= decl_rangevar
				f.scope.add_named_decl(d)
			***REMOVED***
		***REMOVED***

		if a.Value != nil ***REMOVED***
			if t, ok := a.Value.(*ast.Ident); ok ***REMOVED***
				d := new_decl_var(t.Name, nil, a.X, 1, prevscope)
				if d != nil ***REMOVED***
					d.flags |= decl_rangevar
					f.scope.add_named_decl(d)
				***REMOVED***
			***REMOVED***
		***REMOVED***
	***REMOVED***

	f.process_block_stmt(a.Body)
***REMOVED***

func (f *auto_complete_file) process_assign_stmt(a *ast.AssignStmt) ***REMOVED***
	if a.Tok != token.DEFINE || f.offset(a.TokPos) > f.cursor ***REMOVED***
		return
	***REMOVED***

	names := make([]*ast.Ident, len(a.Lhs))
	for i, name := range a.Lhs ***REMOVED***
		id, ok := name.(*ast.Ident)
		if !ok ***REMOVED***
			// something is wrong, just ignore the whole stmt
			return
		***REMOVED***
		names[i] = id
	***REMOVED***

	var prevscope *scope
	f.scope, prevscope = advance_scope(f.scope)

	pack := decl_pack***REMOVED***names, nil, a.Rhs***REMOVED***
	for i, name := range pack.names ***REMOVED***
		typ, v, vi := pack.type_value_index(i)
		d := new_decl_var(name.Name, typ, v, vi, prevscope)
		if d == nil ***REMOVED***
			continue
		***REMOVED***

		f.scope.add_named_decl(d)
	***REMOVED***
***REMOVED***

func (f *auto_complete_file) process_field_list(field_list *ast.FieldList, s *scope) ***REMOVED***
	if field_list != nil ***REMOVED***
		decls := ast_field_list_to_decls(field_list, decl_var, 0, s, false)
		for _, d := range decls ***REMOVED***
			f.scope.add_named_decl(d)
		***REMOVED***
	***REMOVED***
***REMOVED***

func (f *auto_complete_file) cursor_in_if_head(s *ast.IfStmt) bool ***REMOVED***
	if f.cursor > f.offset(s.If) && f.cursor <= f.offset(s.Body.Lbrace) ***REMOVED***
		return true
	***REMOVED***
	return false
***REMOVED***

func (f *auto_complete_file) cursor_in_if_stmt(s *ast.IfStmt) bool ***REMOVED***
	if f.cursor > f.offset(s.If) ***REMOVED***
		// magic -10 comes from auto_complete_file.offset method, see
		// len() expr in there
		if f.offset(s.End()) == -10 || f.cursor < f.offset(s.End()) ***REMOVED***
			return true
		***REMOVED***
	***REMOVED***
	return false
***REMOVED***

func (f *auto_complete_file) cursor_in_for_head(s *ast.ForStmt) bool ***REMOVED***
	if f.cursor > f.offset(s.For) && f.cursor <= f.offset(s.Body.Lbrace) ***REMOVED***
		return true
	***REMOVED***
	return false
***REMOVED***

func (f *auto_complete_file) cursor_in(block *ast.BlockStmt) bool ***REMOVED***
	if f.cursor == -1 || block == nil ***REMOVED***
		return false
	***REMOVED***

	if f.cursor > f.offset(block.Lbrace) && f.cursor <= f.offset(block.Rbrace) ***REMOVED***
		return true
	***REMOVED***
	return false
***REMOVED***
