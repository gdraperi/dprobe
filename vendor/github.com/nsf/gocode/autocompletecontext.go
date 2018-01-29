package main

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"sort"
	"strings"
	"time"
)

//-------------------------------------------------------------------------
// out_buffers
//
// Temporary structure for writing autocomplete response.
//-------------------------------------------------------------------------

// fields must be exported for RPC
type candidate struct ***REMOVED***
	Name    string
	Type    string
	Class   decl_class
	Package string
***REMOVED***

type out_buffers struct ***REMOVED***
	tmpbuf            *bytes.Buffer
	candidates        []candidate
	canonical_aliases map[string]string
	ctx               *auto_complete_context
	tmpns             map[string]bool
	ignorecase        bool
***REMOVED***

func new_out_buffers(ctx *auto_complete_context) *out_buffers ***REMOVED***
	b := new(out_buffers)
	b.tmpbuf = bytes.NewBuffer(make([]byte, 0, 1024))
	b.candidates = make([]candidate, 0, 64)
	b.ctx = ctx
	b.canonical_aliases = make(map[string]string)
	for _, imp := range b.ctx.current.packages ***REMOVED***
		b.canonical_aliases[imp.abspath] = imp.alias
	***REMOVED***
	return b
***REMOVED***

func (b *out_buffers) Len() int ***REMOVED***
	return len(b.candidates)
***REMOVED***

func (b *out_buffers) Less(i, j int) bool ***REMOVED***
	x := b.candidates[i]
	y := b.candidates[j]
	if x.Class == y.Class ***REMOVED***
		return x.Name < y.Name
	***REMOVED***
	return x.Class < y.Class
***REMOVED***

func (b *out_buffers) Swap(i, j int) ***REMOVED***
	b.candidates[i], b.candidates[j] = b.candidates[j], b.candidates[i]
***REMOVED***

func (b *out_buffers) append_decl(p, name, pkg string, decl *decl, class decl_class) ***REMOVED***
	c1 := !g_config.ProposeBuiltins && decl.scope == g_universe_scope && decl.name != "Error"
	c2 := class != decl_invalid && decl.class != class
	c3 := class == decl_invalid && !has_prefix(name, p, b.ignorecase)
	c4 := !decl.matches()
	c5 := !check_type_expr(decl.typ)

	if c1 || c2 || c3 || c4 || c5 ***REMOVED***
		return
	***REMOVED***

	decl.pretty_print_type(b.tmpbuf, b.canonical_aliases)
	b.candidates = append(b.candidates, candidate***REMOVED***
		Name:    name,
		Type:    b.tmpbuf.String(),
		Class:   decl.class,
		Package: pkg,
	***REMOVED***)
	b.tmpbuf.Reset()
***REMOVED***

func (b *out_buffers) append_embedded(p string, decl *decl, pkg string, class decl_class) ***REMOVED***
	if decl.embedded == nil ***REMOVED***
		return
	***REMOVED***

	first_level := false
	if b.tmpns == nil ***REMOVED***
		// first level, create tmp namespace
		b.tmpns = make(map[string]bool)
		first_level = true

		// add all children of the current decl to the namespace
		for _, c := range decl.children ***REMOVED***
			b.tmpns[c.name] = true
		***REMOVED***
	***REMOVED***

	for _, emb := range decl.embedded ***REMOVED***
		typedecl := type_to_decl(emb, decl.scope)
		if typedecl == nil ***REMOVED***
			continue
		***REMOVED***

		// could be type alias
		if typedecl.is_alias() ***REMOVED***
			typedecl = typedecl.type_dealias()
		***REMOVED***

		// prevent infinite recursion here
		if typedecl.is_visited() ***REMOVED***
			continue
		***REMOVED***
		typedecl.set_visited()
		defer typedecl.clear_visited()

		for _, c := range typedecl.children ***REMOVED***
			if _, has := b.tmpns[c.name]; has ***REMOVED***
				continue
			***REMOVED***
			b.append_decl(p, c.name, pkg, c, class)
			b.tmpns[c.name] = true
		***REMOVED***
		b.append_embedded(p, typedecl, pkg, class)
	***REMOVED***

	if first_level ***REMOVED***
		// remove tmp namespace
		b.tmpns = nil
	***REMOVED***
***REMOVED***

//-------------------------------------------------------------------------
// auto_complete_context
//
// Context that holds cache structures for autocompletion needs. It
// includes cache for packages and for main package files.
//-------------------------------------------------------------------------

type auto_complete_context struct ***REMOVED***
	current *auto_complete_file // currently edited file
	others  []*decl_file_cache  // other files of the current package
	pkg     *scope

	pcache    package_cache // packages cache
	declcache *decl_cache   // top-level declarations cache
***REMOVED***

func new_auto_complete_context(pcache package_cache, declcache *decl_cache) *auto_complete_context ***REMOVED***
	c := new(auto_complete_context)
	c.current = new_auto_complete_file("", declcache.context)
	c.pcache = pcache
	c.declcache = declcache
	return c
***REMOVED***

func (c *auto_complete_context) update_caches() ***REMOVED***
	// temporary map for packages that we need to check for a cache expiration
	// map is used as a set of unique items to prevent double checks
	ps := make(map[string]*package_file_cache)

	// collect import information from all of the files
	c.pcache.append_packages(ps, c.current.packages)
	c.others = get_other_package_files(c.current.name, c.current.package_name, c.declcache)
	for _, other := range c.others ***REMOVED***
		c.pcache.append_packages(ps, other.packages)
	***REMOVED***

	update_packages(ps)

	// fix imports for all files
	fixup_packages(c.current.filescope, c.current.packages, c.pcache)
	for _, f := range c.others ***REMOVED***
		fixup_packages(f.filescope, f.packages, c.pcache)
	***REMOVED***

	// At this point we have collected all top level declarations, now we need to
	// merge them in the common package block.
	c.merge_decls()
***REMOVED***

func (c *auto_complete_context) merge_decls() ***REMOVED***
	c.pkg = new_scope(g_universe_scope)
	merge_decls(c.current.filescope, c.pkg, c.current.decls)
	merge_decls_from_packages(c.pkg, c.current.packages, c.pcache)
	for _, f := range c.others ***REMOVED***
		merge_decls(f.filescope, c.pkg, f.decls)
		merge_decls_from_packages(c.pkg, f.packages, c.pcache)
	***REMOVED***

	// special pass for type aliases which also have methods, while this is
	// valid code, it shouldn't happen a lot in practice, so, whatever
	// let's move all type alias methods to their first non-alias type down in
	// the chain
	propagate_type_alias_methods(c.pkg)
***REMOVED***

func (c *auto_complete_context) make_decl_set(scope *scope) map[string]*decl ***REMOVED***
	set := make(map[string]*decl, len(c.pkg.entities)*2)
	make_decl_set_recursive(set, scope)
	return set
***REMOVED***

func (c *auto_complete_context) get_candidates_from_set(set map[string]*decl, partial string, class decl_class, b *out_buffers) ***REMOVED***
	for key, value := range set ***REMOVED***
		if value == nil ***REMOVED***
			continue
		***REMOVED***
		value.infer_type()
		pkgname := ""
		if pkg, ok := c.pcache[value.name]; ok ***REMOVED***
			pkgname = pkg.import_name
		***REMOVED***
		b.append_decl(partial, key, pkgname, value, class)
	***REMOVED***
***REMOVED***

func (c *auto_complete_context) get_candidates_from_decl_alias(cc cursor_context, class decl_class, b *out_buffers) ***REMOVED***
	if cc.decl.is_visited() ***REMOVED***
		return
	***REMOVED***

	cc.decl = cc.decl.type_dealias()
	if cc.decl == nil ***REMOVED***
		return
	***REMOVED***

	cc.decl.set_visited()
	defer cc.decl.clear_visited()

	c.get_candidates_from_decl(cc, class, b)
	return
***REMOVED***

func (c *auto_complete_context) decl_package_import_path(decl *decl) string ***REMOVED***
	if decl == nil || decl.scope == nil ***REMOVED***
		return ""
	***REMOVED***
	if pkg, ok := c.pcache[decl.scope.pkgname]; ok ***REMOVED***
		return pkg.import_name
	***REMOVED***
	return ""
***REMOVED***

func (c *auto_complete_context) get_candidates_from_decl(cc cursor_context, class decl_class, b *out_buffers) ***REMOVED***
	if cc.decl.is_alias() ***REMOVED***
		c.get_candidates_from_decl_alias(cc, class, b)
		return
	***REMOVED***

	// propose all children of a subject declaration and
	for _, decl := range cc.decl.children ***REMOVED***
		if cc.decl.class == decl_package && !ast.IsExported(decl.name) ***REMOVED***
			continue
		***REMOVED***
		if cc.struct_field ***REMOVED***
			// if we're autocompleting struct field init, skip all methods
			if _, ok := decl.typ.(*ast.FuncType); ok ***REMOVED***
				continue
			***REMOVED***
		***REMOVED***
		b.append_decl(cc.partial, decl.name, c.decl_package_import_path(decl), decl, class)
	***REMOVED***
	// propose all children of an underlying struct/interface type
	adecl := advance_to_struct_or_interface(cc.decl)
	if adecl != nil && adecl != cc.decl ***REMOVED***
		for _, decl := range adecl.children ***REMOVED***
			if decl.class == decl_var ***REMOVED***
				b.append_decl(cc.partial, decl.name, c.decl_package_import_path(decl), decl, class)
			***REMOVED***
		***REMOVED***
	***REMOVED***
	// propose all children of its embedded types
	b.append_embedded(cc.partial, cc.decl, c.decl_package_import_path(cc.decl), class)
***REMOVED***

func (c *auto_complete_context) get_import_candidates(partial string, b *out_buffers) ***REMOVED***
	currentPackagePath, pkgdirs := g_daemon.context.pkg_dirs()
	resultSet := map[string]struct***REMOVED******REMOVED******REMOVED******REMOVED***
	for _, pkgdir := range pkgdirs ***REMOVED***
		// convert srcpath to pkgpath and get candidates
		get_import_candidates_dir(pkgdir, filepath.FromSlash(partial), b.ignorecase, currentPackagePath, resultSet)
	***REMOVED***
	for k := range resultSet ***REMOVED***
		b.candidates = append(b.candidates, candidate***REMOVED***Name: k, Class: decl_import***REMOVED***)
	***REMOVED***
***REMOVED***

func get_import_candidates_dir(root, partial string, ignorecase bool, currentPackagePath string, r map[string]struct***REMOVED******REMOVED***) ***REMOVED***
	var fpath string
	var match bool
	if strings.HasSuffix(partial, "/") ***REMOVED***
		fpath = filepath.Join(root, partial)
	***REMOVED*** else ***REMOVED***
		fpath = filepath.Join(root, filepath.Dir(partial))
		match = true
	***REMOVED***
	fi := readdir(fpath)
	for i := range fi ***REMOVED***
		name := fi[i].Name()
		rel, err := filepath.Rel(root, filepath.Join(fpath, name))
		if err != nil ***REMOVED***
			panic(err)
		***REMOVED***
		if match && !has_prefix(rel, partial, ignorecase) ***REMOVED***
			continue
		***REMOVED*** else if fi[i].IsDir() ***REMOVED***
			get_import_candidates_dir(root, rel+string(filepath.Separator), ignorecase, currentPackagePath, r)
		***REMOVED*** else ***REMOVED***
			ext := filepath.Ext(name)
			if ext != ".a" ***REMOVED***
				continue
			***REMOVED*** else ***REMOVED***
				rel = rel[0 : len(rel)-2]
			***REMOVED***
			if ipath, ok := vendorlessImportPath(filepath.ToSlash(rel), currentPackagePath); ok ***REMOVED***
				r[ipath] = struct***REMOVED******REMOVED******REMOVED******REMOVED***
			***REMOVED***
		***REMOVED***
	***REMOVED***
***REMOVED***

// returns three slices of the same length containing:
// 1. apropos names
// 2. apropos types (pretty-printed)
// 3. apropos classes
// and length of the part that should be replaced (if any)
func (c *auto_complete_context) apropos(file []byte, filename string, cursor int) ([]candidate, int) ***REMOVED***
	c.current.cursor = cursor
	c.current.name = filename

	// Update caches and parse the current file.
	// This process is quite complicated, because I was trying to design it in a
	// concurrent fashion. Apparently I'm not really good at that. Hopefully
	// will be better in future.

	// Ugly hack, but it actually may help in some cases. Insert a
	// semicolon right at the cursor location.
	filesemi := make([]byte, len(file)+1)
	copy(filesemi, file[:cursor])
	filesemi[cursor] = ';'
	copy(filesemi[cursor+1:], file[cursor:])

	// Does full processing of the currently edited file (top-level declarations plus
	// active function).
	c.current.process_data(filesemi)

	// Updates cache of other files and packages. See the function for details of
	// the process. At the end merges all the top-level declarations into the package
	// block.
	c.update_caches()

	// And we're ready to Go. ;)

	b := new_out_buffers(c)
	if g_config.IgnoreCase ***REMOVED***
		if *g_debug ***REMOVED***
			log.Printf("ignoring case sensitivity")
		***REMOVED***
		b.ignorecase = true
	***REMOVED***

	cc, ok := c.deduce_cursor_context(file, cursor)
	partial := len(cc.partial)
	if !g_config.Partials ***REMOVED***
		if *g_debug ***REMOVED***
			log.Printf("not performing partial prefix matching")
		***REMOVED***
		cc.partial = ""
	***REMOVED***
	if !ok ***REMOVED***
		var d *decl
		if ident, ok := cc.expr.(*ast.Ident); ok && g_config.UnimportedPackages ***REMOVED***
			p := resolveKnownPackageIdent(ident.Name, c.current.name, c.current.context)
			if p != nil ***REMOVED***
				c.pcache[p.name] = p
				d = p.main
			***REMOVED***
		***REMOVED***
		if d == nil ***REMOVED***
			return nil, 0
		***REMOVED***
		cc.decl = d
	***REMOVED***

	class := decl_invalid
	if g_config.ClassFiltering ***REMOVED***
		switch cc.partial ***REMOVED***
		case "const":
			class = decl_const
		case "var":
			class = decl_var
		case "type":
			class = decl_type
		case "func":
			class = decl_func
		case "package":
			class = decl_package
		***REMOVED***
	***REMOVED***

	if cc.decl_import ***REMOVED***
		c.get_import_candidates(cc.partial, b)
		if cc.partial != "" && len(b.candidates) == 0 ***REMOVED***
			// as a fallback, try case insensitive approach
			b.ignorecase = true
			c.get_import_candidates(cc.partial, b)
		***REMOVED***
	***REMOVED*** else if cc.decl == nil ***REMOVED***
		// In case if no declaraion is a subject of completion, propose all:
		set := c.make_decl_set(c.current.scope)
		c.get_candidates_from_set(set, cc.partial, class, b)
		if cc.partial != "" && len(b.candidates) == 0 ***REMOVED***
			// as a fallback, try case insensitive approach
			b.ignorecase = true
			c.get_candidates_from_set(set, cc.partial, class, b)
		***REMOVED***
	***REMOVED*** else ***REMOVED***
		c.get_candidates_from_decl(cc, class, b)
		if cc.partial != "" && len(b.candidates) == 0 ***REMOVED***
			// as a fallback, try case insensitive approach
			b.ignorecase = true
			c.get_candidates_from_decl(cc, class, b)
		***REMOVED***
	***REMOVED***

	if len(b.candidates) == 0 ***REMOVED***
		return nil, 0
	***REMOVED***

	sort.Sort(b)
	return b.candidates, partial
***REMOVED***

func update_packages(ps map[string]*package_file_cache) ***REMOVED***
	// initiate package cache update
	done := make(chan bool)
	for _, p := range ps ***REMOVED***
		go func(p *package_file_cache) ***REMOVED***
			defer func() ***REMOVED***
				if err := recover(); err != nil ***REMOVED***
					print_backtrace(err)
					done <- false
				***REMOVED***
			***REMOVED***()
			p.update_cache()
			done <- true
		***REMOVED***(p)
	***REMOVED***

	// wait for its completion
	for _ = range ps ***REMOVED***
		if !<-done ***REMOVED***
			panic("One of the package cache updaters panicked")
		***REMOVED***
	***REMOVED***
***REMOVED***

func collect_type_alias_methods(d *decl) map[string]*decl ***REMOVED***
	if d == nil || d.is_visited() || !d.is_alias() ***REMOVED***
		return nil
	***REMOVED***
	d.set_visited()
	defer d.clear_visited()

	// add own methods
	m := map[string]*decl***REMOVED******REMOVED***
	for k, v := range d.children ***REMOVED***
		m[k] = v
	***REMOVED***

	// recurse into more aliases
	dd := type_to_decl(d.typ, d.scope)
	for k, v := range collect_type_alias_methods(dd) ***REMOVED***
		m[k] = v
	***REMOVED***

	return m
***REMOVED***

func propagate_type_alias_methods(s *scope) ***REMOVED***
	for _, e := range s.entities ***REMOVED***
		if !e.is_alias() ***REMOVED***
			continue
		***REMOVED***

		methods := collect_type_alias_methods(e)
		if len(methods) == 0 ***REMOVED***
			continue
		***REMOVED***

		dd := e.type_dealias()
		if dd == nil ***REMOVED***
			continue
		***REMOVED***

		decl := dd.deep_copy()
		for _, v := range methods ***REMOVED***
			decl.add_child(v)
		***REMOVED***
		s.entities[decl.name] = decl
	***REMOVED***
***REMOVED***

func merge_decls(filescope *scope, pkg *scope, decls map[string]*decl) ***REMOVED***
	for _, d := range decls ***REMOVED***
		pkg.merge_decl(d)
	***REMOVED***
	filescope.parent = pkg
***REMOVED***

func merge_decls_from_packages(pkgscope *scope, pkgs []package_import, pcache package_cache) ***REMOVED***
	for _, p := range pkgs ***REMOVED***
		path, alias := p.abspath, p.alias
		if alias != "." ***REMOVED***
			continue
		***REMOVED***
		p := pcache[path].main
		if p == nil ***REMOVED***
			continue
		***REMOVED***
		for _, d := range p.children ***REMOVED***
			if ast.IsExported(d.name) ***REMOVED***
				pkgscope.merge_decl(d)
			***REMOVED***
		***REMOVED***
	***REMOVED***
***REMOVED***

func fixup_packages(filescope *scope, pkgs []package_import, pcache package_cache) ***REMOVED***
	for _, p := range pkgs ***REMOVED***
		path, alias := p.abspath, p.alias
		if alias == "" ***REMOVED***
			alias = pcache[path].defalias
		***REMOVED***
		// skip packages that will be merged to the package scope
		if alias == "." ***REMOVED***
			continue
		***REMOVED***
		filescope.replace_decl(alias, pcache[path].main)
	***REMOVED***
***REMOVED***

func get_other_package_files(filename, packageName string, declcache *decl_cache) []*decl_file_cache ***REMOVED***
	others := find_other_package_files(filename, packageName)

	ret := make([]*decl_file_cache, len(others))
	done := make(chan *decl_file_cache)

	for _, nm := range others ***REMOVED***
		go func(name string) ***REMOVED***
			defer func() ***REMOVED***
				if err := recover(); err != nil ***REMOVED***
					print_backtrace(err)
					done <- nil
				***REMOVED***
			***REMOVED***()
			done <- declcache.get_and_update(name)
		***REMOVED***(nm)
	***REMOVED***

	for i := range others ***REMOVED***
		ret[i] = <-done
		if ret[i] == nil ***REMOVED***
			panic("One of the decl cache updaters panicked")
		***REMOVED***
	***REMOVED***

	return ret
***REMOVED***

func find_other_package_files(filename, package_name string) []string ***REMOVED***
	if filename == "" ***REMOVED***
		return nil
	***REMOVED***

	dir, file := filepath.Split(filename)
	files_in_dir, err := readdir_lstat(dir)
	if err != nil ***REMOVED***
		panic(err)
	***REMOVED***

	count := 0
	for _, stat := range files_in_dir ***REMOVED***
		ok, _ := filepath.Match("*.go", stat.Name())
		if !ok || stat.Name() == file ***REMOVED***
			continue
		***REMOVED***
		count++
	***REMOVED***

	out := make([]string, 0, count)
	for _, stat := range files_in_dir ***REMOVED***
		const non_regular = os.ModeDir | os.ModeSymlink |
			os.ModeDevice | os.ModeNamedPipe | os.ModeSocket

		ok, _ := filepath.Match("*.go", stat.Name())
		if !ok || stat.Name() == file || stat.Mode()&non_regular != 0 ***REMOVED***
			continue
		***REMOVED***

		abspath := filepath.Join(dir, stat.Name())
		if file_package_name(abspath) == package_name ***REMOVED***
			n := len(out)
			out = out[:n+1]
			out[n] = abspath
		***REMOVED***
	***REMOVED***

	return out
***REMOVED***

func file_package_name(filename string) string ***REMOVED***
	file, _ := parser.ParseFile(token.NewFileSet(), filename, nil, parser.PackageClauseOnly)
	return file.Name.Name
***REMOVED***

func make_decl_set_recursive(set map[string]*decl, scope *scope) ***REMOVED***
	for name, ent := range scope.entities ***REMOVED***
		if _, ok := set[name]; !ok ***REMOVED***
			set[name] = ent
		***REMOVED***
	***REMOVED***
	if scope.parent != nil ***REMOVED***
		make_decl_set_recursive(set, scope.parent)
	***REMOVED***
***REMOVED***

func check_func_field_list(f *ast.FieldList) bool ***REMOVED***
	if f == nil ***REMOVED***
		return true
	***REMOVED***

	for _, field := range f.List ***REMOVED***
		if !check_type_expr(field.Type) ***REMOVED***
			return false
		***REMOVED***
	***REMOVED***
	return true
***REMOVED***

// checks for a type expression correctness, it the type expression has
// ast.BadExpr somewhere, returns false, otherwise true
func check_type_expr(e ast.Expr) bool ***REMOVED***
	switch t := e.(type) ***REMOVED***
	case *ast.StarExpr:
		return check_type_expr(t.X)
	case *ast.ArrayType:
		return check_type_expr(t.Elt)
	case *ast.SelectorExpr:
		return check_type_expr(t.X)
	case *ast.FuncType:
		a := check_func_field_list(t.Params)
		b := check_func_field_list(t.Results)
		return a && b
	case *ast.MapType:
		a := check_type_expr(t.Key)
		b := check_type_expr(t.Value)
		return a && b
	case *ast.Ellipsis:
		return check_type_expr(t.Elt)
	case *ast.ChanType:
		return check_type_expr(t.Value)
	case *ast.BadExpr:
		return false
	default:
		return true
	***REMOVED***
***REMOVED***

//-------------------------------------------------------------------------
// Status output
//-------------------------------------------------------------------------

type decl_slice []*decl

func (s decl_slice) Less(i, j int) bool ***REMOVED***
	if s[i].class != s[j].class ***REMOVED***
		return s[i].name < s[j].name
	***REMOVED***
	return s[i].class < s[j].class
***REMOVED***
func (s decl_slice) Len() int      ***REMOVED*** return len(s) ***REMOVED***
func (s decl_slice) Swap(i, j int) ***REMOVED*** s[i], s[j] = s[j], s[i] ***REMOVED***

const (
	color_red          = "\033[0;31m"
	color_red_bold     = "\033[1;31m"
	color_green        = "\033[0;32m"
	color_green_bold   = "\033[1;32m"
	color_yellow       = "\033[0;33m"
	color_yellow_bold  = "\033[1;33m"
	color_blue         = "\033[0;34m"
	color_blue_bold    = "\033[1;34m"
	color_magenta      = "\033[0;35m"
	color_magenta_bold = "\033[1;35m"
	color_cyan         = "\033[0;36m"
	color_cyan_bold    = "\033[1;36m"
	color_white        = "\033[0;37m"
	color_white_bold   = "\033[1;37m"
	color_none         = "\033[0m"
)

var g_decl_class_to_color = [...]string***REMOVED***
	decl_const:        color_white_bold,
	decl_var:          color_magenta,
	decl_type:         color_cyan,
	decl_func:         color_green,
	decl_package:      color_red,
	decl_methods_stub: color_red,
***REMOVED***

var g_decl_class_to_string_status = [...]string***REMOVED***
	decl_const:        "  const",
	decl_var:          "    var",
	decl_type:         "   type",
	decl_func:         "   func",
	decl_package:      "package",
	decl_methods_stub: "   stub",
***REMOVED***

func (c *auto_complete_context) status() string ***REMOVED***

	buf := bytes.NewBuffer(make([]byte, 0, 4096))
	fmt.Fprintf(buf, "Server's GOMAXPROCS == %d\n", runtime.GOMAXPROCS(0))
	fmt.Fprintf(buf, "\nPackage cache contains %d entries\n", len(c.pcache))
	fmt.Fprintf(buf, "\nListing these entries:\n")
	for _, mod := range c.pcache ***REMOVED***
		fmt.Fprintf(buf, "\tname: %s (default alias: %s)\n", mod.name, mod.defalias)
		fmt.Fprintf(buf, "\timports %d declarations and %d packages\n", len(mod.main.children), len(mod.others))
		if mod.mtime == -1 ***REMOVED***
			fmt.Fprintf(buf, "\tthis package stays in cache forever (built-in package)\n")
		***REMOVED*** else ***REMOVED***
			mtime := time.Unix(0, mod.mtime)
			fmt.Fprintf(buf, "\tlast modification time: %s\n", mtime)
		***REMOVED***
		fmt.Fprintf(buf, "\n")
	***REMOVED***
	if c.current.name != "" ***REMOVED***
		fmt.Fprintf(buf, "Last edited file: %s (package: %s)\n", c.current.name, c.current.package_name)
		if len(c.others) > 0 ***REMOVED***
			fmt.Fprintf(buf, "\nOther files from the current package:\n")
		***REMOVED***
		for _, f := range c.others ***REMOVED***
			fmt.Fprintf(buf, "\t%s\n", f.name)
		***REMOVED***
		fmt.Fprintf(buf, "\nListing declarations from files:\n")

		const status_decls = "\t%s%s" + color_none + " " + color_yellow + "%s" + color_none + "\n"
		const status_decls_children = "\t%s%s" + color_none + " " + color_yellow + "%s" + color_none + " (%d)\n"

		fmt.Fprintf(buf, "\n%s:\n", c.current.name)
		ds := make(decl_slice, len(c.current.decls))
		i := 0
		for _, d := range c.current.decls ***REMOVED***
			ds[i] = d
			i++
		***REMOVED***
		sort.Sort(ds)
		for _, d := range ds ***REMOVED***
			if len(d.children) > 0 ***REMOVED***
				fmt.Fprintf(buf, status_decls_children,
					g_decl_class_to_color[d.class],
					g_decl_class_to_string_status[d.class],
					d.name, len(d.children))
			***REMOVED*** else ***REMOVED***
				fmt.Fprintf(buf, status_decls,
					g_decl_class_to_color[d.class],
					g_decl_class_to_string_status[d.class],
					d.name)
			***REMOVED***
		***REMOVED***

		for _, f := range c.others ***REMOVED***
			fmt.Fprintf(buf, "\n%s:\n", f.name)
			ds = make(decl_slice, len(f.decls))
			i = 0
			for _, d := range f.decls ***REMOVED***
				ds[i] = d
				i++
			***REMOVED***
			sort.Sort(ds)
			for _, d := range ds ***REMOVED***
				if len(d.children) > 0 ***REMOVED***
					fmt.Fprintf(buf, status_decls_children,
						g_decl_class_to_color[d.class],
						g_decl_class_to_string_status[d.class],
						d.name, len(d.children))
				***REMOVED*** else ***REMOVED***
					fmt.Fprintf(buf, status_decls,
						g_decl_class_to_color[d.class],
						g_decl_class_to_string_status[d.class],
						d.name)
				***REMOVED***
			***REMOVED***
		***REMOVED***
	***REMOVED***
	return buf.String()
***REMOVED***
