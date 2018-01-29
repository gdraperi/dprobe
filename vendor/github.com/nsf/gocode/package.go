package main

import (
	"bytes"
	"fmt"
	"go/ast"
	"os"
	"strings"
)

type package_parser interface ***REMOVED***
	parse_export(callback func(pkg string, decl ast.Decl))
***REMOVED***

//-------------------------------------------------------------------------
// package_file_cache
//
// Structure that represents a cache for an imported pacakge. In other words
// these are the contents of an archive (*.a) file.
//-------------------------------------------------------------------------

type package_file_cache struct ***REMOVED***
	name        string // file name
	import_name string
	mtime       int64
	defalias    string

	scope  *scope
	main   *decl // package declaration
	others map[string]*decl
***REMOVED***

func new_package_file_cache(absname, name string) *package_file_cache ***REMOVED***
	m := new(package_file_cache)
	m.name = absname
	m.import_name = name
	m.mtime = 0
	m.defalias = ""
	return m
***REMOVED***

// Creates a cache that stays in cache forever. Useful for built-in packages.
func new_package_file_cache_forever(name, defalias string) *package_file_cache ***REMOVED***
	m := new(package_file_cache)
	m.name = name
	m.mtime = -1
	m.defalias = defalias
	return m
***REMOVED***

func (m *package_file_cache) find_file() string ***REMOVED***
	if file_exists(m.name) ***REMOVED***
		return m.name
	***REMOVED***

	n := len(m.name)
	filename := m.name[:n-1] + "6"
	if file_exists(filename) ***REMOVED***
		return filename
	***REMOVED***

	filename = m.name[:n-1] + "8"
	if file_exists(filename) ***REMOVED***
		return filename
	***REMOVED***

	filename = m.name[:n-1] + "5"
	if file_exists(filename) ***REMOVED***
		return filename
	***REMOVED***
	return m.name
***REMOVED***

func (m *package_file_cache) update_cache() ***REMOVED***
	if m.mtime == -1 ***REMOVED***
		return
	***REMOVED***
	fname := m.find_file()
	stat, err := os.Stat(fname)
	if err != nil ***REMOVED***
		return
	***REMOVED***

	statmtime := stat.ModTime().UnixNano()
	if m.mtime != statmtime ***REMOVED***
		m.mtime = statmtime

		data, err := file_reader.read_file(fname)
		if err != nil ***REMOVED***
			return
		***REMOVED***
		m.process_package_data(data)
	***REMOVED***
***REMOVED***

func (m *package_file_cache) process_package_data(data []byte) ***REMOVED***
	m.scope = new_named_scope(g_universe_scope, m.name)

	// find import section
	i := bytes.Index(data, []byte***REMOVED***'\n', '$', '$'***REMOVED***)
	if i == -1 ***REMOVED***
		panic(fmt.Sprintf("Can't find the import section in the package file %s", m.name))
	***REMOVED***
	data = data[i+len("\n$$"):]

	// main package
	m.main = new_decl(m.name, decl_package, nil)
	// create map for other packages
	m.others = make(map[string]*decl)

	var pp package_parser
	if data[0] == 'B' ***REMOVED***
		// binary format, skip 'B\n'
		data = data[2:]
		var p gc_bin_parser
		p.init(data, m)
		pp = &p
	***REMOVED*** else ***REMOVED***
		// textual format, find the beginning of the package clause
		i = bytes.Index(data, []byte***REMOVED***'p', 'a', 'c', 'k', 'a', 'g', 'e'***REMOVED***)
		if i == -1 ***REMOVED***
			panic("Can't find the package clause")
		***REMOVED***
		data = data[i:]

		var p gc_parser
		p.init(data, m)
		pp = &p
	***REMOVED***

	prefix := "!" + m.name + "!"
	pp.parse_export(func(pkg string, decl ast.Decl) ***REMOVED***
		anonymify_ast(decl, decl_foreign, m.scope)
		if pkg == "" || strings.HasPrefix(pkg, prefix) ***REMOVED***
			// main package
			add_ast_decl_to_package(m.main, decl, m.scope)
		***REMOVED*** else ***REMOVED***
			// others
			if _, ok := m.others[pkg]; !ok ***REMOVED***
				m.others[pkg] = new_decl(pkg, decl_package, nil)
			***REMOVED***
			add_ast_decl_to_package(m.others[pkg], decl, m.scope)
		***REMOVED***
	***REMOVED***)

	// hack, add ourselves to the package scope
	mainName := "!" + m.name + "!" + m.defalias
	m.add_package_to_scope(mainName, m.name)

	// replace dummy package decls in package scope to actual packages
	for key := range m.scope.entities ***REMOVED***
		if !strings.HasPrefix(key, "!") ***REMOVED***
			continue
		***REMOVED***
		pkg, ok := m.others[key]
		if !ok && key == mainName ***REMOVED***
			pkg = m.main
		***REMOVED***
		m.scope.replace_decl(key, pkg)
	***REMOVED***
***REMOVED***

func (m *package_file_cache) add_package_to_scope(alias, realname string) ***REMOVED***
	d := new_decl(realname, decl_package, nil)
	m.scope.add_decl(alias, d)
***REMOVED***

func add_ast_decl_to_package(pkg *decl, decl ast.Decl, scope *scope) ***REMOVED***
	foreach_decl(decl, func(data *foreach_decl_struct) ***REMOVED***
		class := ast_decl_class(data.decl)
		for i, name := range data.names ***REMOVED***
			typ, v, vi := data.type_value_index(i)

			d := new_decl_full(name.Name, class, decl_foreign|ast_decl_flags(data.decl), typ, v, vi, scope)
			if d == nil ***REMOVED***
				return
			***REMOVED***

			if !name.IsExported() && d.class != decl_type ***REMOVED***
				return
			***REMOVED***

			methodof := method_of(data.decl)
			if methodof != "" ***REMOVED***
				decl := pkg.find_child(methodof)
				if decl != nil ***REMOVED***
					decl.add_child(d)
				***REMOVED*** else ***REMOVED***
					decl = new_decl(methodof, decl_methods_stub, scope)
					decl.add_child(d)
					pkg.add_child(decl)
				***REMOVED***
			***REMOVED*** else ***REMOVED***
				decl := pkg.find_child(d.name)
				if decl != nil ***REMOVED***
					decl.expand_or_replace(d)
				***REMOVED*** else ***REMOVED***
					pkg.add_child(d)
				***REMOVED***
			***REMOVED***
		***REMOVED***
	***REMOVED***)
***REMOVED***

//-------------------------------------------------------------------------
// package_cache
//-------------------------------------------------------------------------

type package_cache map[string]*package_file_cache

func new_package_cache() package_cache ***REMOVED***
	m := make(package_cache)

	// add built-in "unsafe" package
	m.add_builtin_unsafe_package()

	return m
***REMOVED***

// Function fills 'ps' set with packages from 'packages' import information.
// In case if package is not in the cache, it creates one and adds one to the cache.
func (c package_cache) append_packages(ps map[string]*package_file_cache, pkgs []package_import) ***REMOVED***
	for _, m := range pkgs ***REMOVED***
		if _, ok := ps[m.abspath]; ok ***REMOVED***
			continue
		***REMOVED***

		if mod, ok := c[m.abspath]; ok ***REMOVED***
			ps[m.abspath] = mod
		***REMOVED*** else ***REMOVED***
			mod = new_package_file_cache(m.abspath, m.path)
			ps[m.abspath] = mod
			c[m.abspath] = mod
		***REMOVED***
	***REMOVED***
***REMOVED***

var g_builtin_unsafe_package = []byte(`
import
$$
package unsafe
	type @"".Pointer uintptr
	func @"".Offsetof (? any) uintptr
	func @"".Sizeof (? any) uintptr
	func @"".Alignof (? any) uintptr

$$
`)

func (c package_cache) add_builtin_unsafe_package() ***REMOVED***
	pkg := new_package_file_cache_forever("unsafe", "unsafe")
	pkg.process_package_data(g_builtin_unsafe_package)
	c["unsafe"] = pkg
***REMOVED***
