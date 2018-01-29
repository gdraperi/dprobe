package main

import (
	"fmt"
	"go/ast"
	"go/build"
	"go/parser"
	"go/token"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"
)

//-------------------------------------------------------------------------
// []package_import
//-------------------------------------------------------------------------

type package_import struct ***REMOVED***
	alias   string
	abspath string
	path    string
***REMOVED***

// Parses import declarations until the first non-import declaration and fills
// `packages` array with import information.
func collect_package_imports(filename string, decls []ast.Decl, context *package_lookup_context) []package_import ***REMOVED***
	pi := make([]package_import, 0, 16)
	for _, decl := range decls ***REMOVED***
		if gd, ok := decl.(*ast.GenDecl); ok && gd.Tok == token.IMPORT ***REMOVED***
			for _, spec := range gd.Specs ***REMOVED***
				imp := spec.(*ast.ImportSpec)
				path, alias := path_and_alias(imp)
				abspath, ok := abs_path_for_package(filename, path, context)
				if ok && alias != "_" ***REMOVED***
					pi = append(pi, package_import***REMOVED***alias, abspath, path***REMOVED***)
				***REMOVED***
			***REMOVED***
		***REMOVED*** else ***REMOVED***
			break
		***REMOVED***
	***REMOVED***
	return pi
***REMOVED***

//-------------------------------------------------------------------------
// decl_file_cache
//
// Contains cache for top-level declarations of a file as well as its
// contents, AST and import information.
//-------------------------------------------------------------------------

type decl_file_cache struct ***REMOVED***
	name  string // file name
	mtime int64  // last modification time

	decls     map[string]*decl // top-level declarations
	error     error            // last error
	packages  []package_import // import information
	filescope *scope

	fset    *token.FileSet
	context *package_lookup_context
***REMOVED***

func new_decl_file_cache(name string, context *package_lookup_context) *decl_file_cache ***REMOVED***
	return &decl_file_cache***REMOVED***
		name:    name,
		context: context,
	***REMOVED***
***REMOVED***

func (f *decl_file_cache) update() ***REMOVED***
	stat, err := os.Stat(f.name)
	if err != nil ***REMOVED***
		f.decls = nil
		f.error = err
		f.fset = nil
		return
	***REMOVED***

	statmtime := stat.ModTime().UnixNano()
	if f.mtime == statmtime ***REMOVED***
		return
	***REMOVED***

	f.mtime = statmtime
	f.read_file()
***REMOVED***

func (f *decl_file_cache) read_file() ***REMOVED***
	var data []byte
	data, f.error = file_reader.read_file(f.name)
	if f.error != nil ***REMOVED***
		return
	***REMOVED***
	data, _ = filter_out_shebang(data)

	f.process_data(data)
***REMOVED***

func (f *decl_file_cache) process_data(data []byte) ***REMOVED***
	var file *ast.File
	f.fset = token.NewFileSet()
	file, f.error = parser.ParseFile(f.fset, "", data, 0)
	f.filescope = new_scope(nil)
	for _, d := range file.Decls ***REMOVED***
		anonymify_ast(d, 0, f.filescope)
	***REMOVED***
	f.packages = collect_package_imports(f.name, file.Decls, f.context)
	f.decls = make(map[string]*decl, len(file.Decls))
	for _, decl := range file.Decls ***REMOVED***
		append_to_top_decls(f.decls, decl, f.filescope)
	***REMOVED***
***REMOVED***

func append_to_top_decls(decls map[string]*decl, decl ast.Decl, scope *scope) ***REMOVED***
	foreach_decl(decl, func(data *foreach_decl_struct) ***REMOVED***
		class := ast_decl_class(data.decl)
		for i, name := range data.names ***REMOVED***
			typ, v, vi := data.type_value_index(i)

			d := new_decl_full(name.Name, class, ast_decl_flags(data.decl), typ, v, vi, scope)
			if d == nil ***REMOVED***
				return
			***REMOVED***

			methodof := method_of(decl)
			if methodof != "" ***REMOVED***
				decl, ok := decls[methodof]
				if ok ***REMOVED***
					decl.add_child(d)
				***REMOVED*** else ***REMOVED***
					decl = new_decl(methodof, decl_methods_stub, scope)
					decls[methodof] = decl
					decl.add_child(d)
				***REMOVED***
			***REMOVED*** else ***REMOVED***
				decl, ok := decls[d.name]
				if ok ***REMOVED***
					decl.expand_or_replace(d)
				***REMOVED*** else ***REMOVED***
					decls[d.name] = d
				***REMOVED***
			***REMOVED***
		***REMOVED***
	***REMOVED***)
***REMOVED***

func abs_path_for_package(filename, p string, context *package_lookup_context) (string, bool) ***REMOVED***
	dir, _ := filepath.Split(filename)
	if len(p) == 0 ***REMOVED***
		return "", false
	***REMOVED***
	if p[0] == '.' ***REMOVED***
		return fmt.Sprintf("%s.a", filepath.Join(dir, p)), true
	***REMOVED***
	pkg, ok := find_go_dag_package(p, dir)
	if ok ***REMOVED***
		return pkg, true
	***REMOVED***
	return find_global_file(p, context)
***REMOVED***

func path_and_alias(imp *ast.ImportSpec) (string, string) ***REMOVED***
	path := ""
	if imp.Path != nil && len(imp.Path.Value) > 0 ***REMOVED***
		path = string(imp.Path.Value)
		path = path[1 : len(path)-1]
	***REMOVED***
	alias := ""
	if imp.Name != nil ***REMOVED***
		alias = imp.Name.Name
	***REMOVED***
	return path, alias
***REMOVED***

func find_go_dag_package(imp, filedir string) (string, bool) ***REMOVED***
	// Support godag directory structure
	dir, pkg := filepath.Split(imp)
	godag_pkg := filepath.Join(filedir, "..", dir, "_obj", pkg+".a")
	if file_exists(godag_pkg) ***REMOVED***
		return godag_pkg, true
	***REMOVED***
	return "", false
***REMOVED***

// autobuild compares the mod time of the source files of the package, and if any of them is newer
// than the package object file will rebuild it.
func autobuild(p *build.Package) error ***REMOVED***
	if p.Dir == "" ***REMOVED***
		return fmt.Errorf("no files to build")
	***REMOVED***
	ps, err := os.Stat(p.PkgObj)
	if err != nil ***REMOVED***
		// Assume package file does not exist and build for the first time.
		return build_package(p)
	***REMOVED***
	pt := ps.ModTime()
	fs, err := readdir_lstat(p.Dir)
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	for _, f := range fs ***REMOVED***
		if f.IsDir() ***REMOVED***
			continue
		***REMOVED***
		if f.ModTime().After(pt) ***REMOVED***
			// Source file is newer than package file; rebuild.
			return build_package(p)
		***REMOVED***
	***REMOVED***
	return nil
***REMOVED***

// build_package builds the package by calling `go install package/import`. If everything compiles
// correctly, the newly compiled package should then be in the usual place in the `$GOPATH/pkg`
// directory, and gocode will pick it up from there.
func build_package(p *build.Package) error ***REMOVED***
	if *g_debug ***REMOVED***
		log.Printf("-------------------")
		log.Printf("rebuilding package %s", p.Name)
		log.Printf("package import: %s", p.ImportPath)
		log.Printf("package object: %s", p.PkgObj)
		log.Printf("package source dir: %s", p.Dir)
		log.Printf("package source files: %v", p.GoFiles)
		log.Printf("GOPATH: %v", g_daemon.context.GOPATH)
		log.Printf("GOROOT: %v", g_daemon.context.GOROOT)
	***REMOVED***
	env := os.Environ()
	for i, v := range env ***REMOVED***
		if strings.HasPrefix(v, "GOPATH=") ***REMOVED***
			env[i] = "GOPATH=" + g_daemon.context.GOPATH
		***REMOVED*** else if strings.HasPrefix(v, "GOROOT=") ***REMOVED***
			env[i] = "GOROOT=" + g_daemon.context.GOROOT
		***REMOVED***
	***REMOVED***

	cmd := exec.Command("go", "install", p.ImportPath)
	cmd.Env = env

	// TODO: Should read STDERR rather than STDOUT.
	out, err := cmd.CombinedOutput()
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	if *g_debug ***REMOVED***
		log.Printf("build out: %s\n", string(out))
	***REMOVED***
	return nil
***REMOVED***

// executes autobuild function if autobuild option is enabled, logs error and
// ignores it
func try_autobuild(p *build.Package) ***REMOVED***
	if g_config.Autobuild ***REMOVED***
		err := autobuild(p)
		if err != nil && *g_debug ***REMOVED***
			log.Printf("Autobuild error: %s\n", err)
		***REMOVED***
	***REMOVED***
***REMOVED***

func log_found_package_maybe(imp, pkgpath string) ***REMOVED***
	if *g_debug ***REMOVED***
		log.Printf("Found %q at %q\n", imp, pkgpath)
	***REMOVED***
***REMOVED***

func log_build_context(context *package_lookup_context) ***REMOVED***
	log.Printf(" GOROOT: %s\n", context.GOROOT)
	log.Printf(" GOPATH: %s\n", context.GOPATH)
	log.Printf(" GOOS: %s\n", context.GOOS)
	log.Printf(" GOARCH: %s\n", context.GOARCH)
	log.Printf(" BzlProjectRoot: %q\n", context.BzlProjectRoot)
	log.Printf(" GBProjectRoot: %q\n", context.GBProjectRoot)
	log.Printf(" lib-path: %q\n", g_config.LibPath)
***REMOVED***

// find_global_file returns the file path of the compiled package corresponding to the specified
// import, and a boolean stating whether such path is valid.
// TODO: Return only one value, possibly empty string if not found.
func find_global_file(imp string, context *package_lookup_context) (string, bool) ***REMOVED***
	// gocode synthetically generates the builtin package
	// "unsafe", since the "unsafe.a" package doesn't really exist.
	// Thus, when the user request for the package "unsafe" we
	// would return synthetic global file that would be used
	// just as a key name to find this synthetic package
	if imp == "unsafe" ***REMOVED***
		return "unsafe", true
	***REMOVED***

	pkgfile := fmt.Sprintf("%s.a", imp)

	// if lib-path is defined, use it
	if g_config.LibPath != "" ***REMOVED***
		for _, p := range filepath.SplitList(g_config.LibPath) ***REMOVED***
			pkg_path := filepath.Join(p, pkgfile)
			if file_exists(pkg_path) ***REMOVED***
				log_found_package_maybe(imp, pkg_path)
				return pkg_path, true
			***REMOVED***
			// Also check the relevant pkg/OS_ARCH dir for the libpath, if provided.
			pkgdir := fmt.Sprintf("%s_%s", context.GOOS, context.GOARCH)
			pkg_path = filepath.Join(p, "pkg", pkgdir, pkgfile)
			if file_exists(pkg_path) ***REMOVED***
				log_found_package_maybe(imp, pkg_path)
				return pkg_path, true
			***REMOVED***
		***REMOVED***
	***REMOVED***

	// gb-specific lookup mode, only if the root dir was found
	if g_config.PackageLookupMode == "gb" && context.GBProjectRoot != "" ***REMOVED***
		root := context.GBProjectRoot
		pkgdir := filepath.Join(root, "pkg", context.GOOS+"-"+context.GOARCH)
		if !is_dir(pkgdir) ***REMOVED***
			pkgdir = filepath.Join(root, "pkg", context.GOOS+"-"+context.GOARCH+"-race")
		***REMOVED***
		pkg_path := filepath.Join(pkgdir, pkgfile)
		if file_exists(pkg_path) ***REMOVED***
			log_found_package_maybe(imp, pkg_path)
			return pkg_path, true
		***REMOVED***
	***REMOVED***

	// bzl-specific lookup mode, only if the root dir was found
	if g_config.PackageLookupMode == "bzl" && context.BzlProjectRoot != "" ***REMOVED***
		var root, impath string
		if strings.HasPrefix(imp, g_config.CustomPkgPrefix+"/") ***REMOVED***
			root = filepath.Join(context.BzlProjectRoot, "bazel-bin")
			impath = imp[len(g_config.CustomPkgPrefix)+1:]
		***REMOVED*** else if g_config.CustomVendorDir != "" ***REMOVED***
			// Try custom vendor dir.
			root = filepath.Join(context.BzlProjectRoot, "bazel-bin", g_config.CustomVendorDir)
			impath = imp
		***REMOVED***

		if root != "" && impath != "" ***REMOVED***
			// There might be more than one ".a" files in the pkg path with bazel.
			// But the best practice is to keep one go_library build target in each
			// pakcage directory so that it follows the standard Go package
			// structure. Thus here we assume there is at most one ".a" file existing
			// in the pkg path.
			if d, err := os.Open(filepath.Join(root, impath)); err == nil ***REMOVED***
				defer d.Close()

				if fis, err := d.Readdir(-1); err == nil ***REMOVED***
					for _, fi := range fis ***REMOVED***
						if !fi.IsDir() && filepath.Ext(fi.Name()) == ".a" ***REMOVED***
							pkg_path := filepath.Join(root, impath, fi.Name())
							log_found_package_maybe(imp, pkg_path)
							return pkg_path, true
						***REMOVED***
					***REMOVED***
				***REMOVED***
			***REMOVED***
		***REMOVED***
	***REMOVED***

	if context.CurrentPackagePath != "" ***REMOVED***
		// Try vendor path first, see GO15VENDOREXPERIMENT.
		// We don't check this environment variable however, seems like there is
		// almost no harm in doing so (well.. if you experiment with vendoring,
		// gocode will fail after enabling/disabling the flag, and you'll be
		// forced to get rid of vendor binaries). But asking users to set this
		// env var is up will bring more trouble. Because we also need to pass
		// it from client to server, make sure their editors set it, etc.
		// So, whatever, let's just pretend it's always on.
		package_path := context.CurrentPackagePath
		for ***REMOVED***
			limp := filepath.Join(package_path, "vendor", imp)
			if p, err := context.Import(limp, "", build.AllowBinary|build.FindOnly); err == nil ***REMOVED***
				try_autobuild(p)
				if file_exists(p.PkgObj) ***REMOVED***
					log_found_package_maybe(imp, p.PkgObj)
					return p.PkgObj, true
				***REMOVED***
			***REMOVED***
			if package_path == "" ***REMOVED***
				break
			***REMOVED***
			next_path := filepath.Dir(package_path)
			// let's protect ourselves from inf recursion here
			if next_path == package_path ***REMOVED***
				break
			***REMOVED***
			package_path = next_path
		***REMOVED***
	***REMOVED***

	if p, err := context.Import(imp, "", build.AllowBinary|build.FindOnly); err == nil ***REMOVED***
		try_autobuild(p)
		if file_exists(p.PkgObj) ***REMOVED***
			log_found_package_maybe(imp, p.PkgObj)
			return p.PkgObj, true
		***REMOVED***
	***REMOVED***

	if *g_debug ***REMOVED***
		log.Printf("Import path %q was not resolved\n", imp)
		log.Println("Gocode's build context is:")
		log_build_context(context)
	***REMOVED***
	return "", false
***REMOVED***

func package_name(file *ast.File) string ***REMOVED***
	if file.Name != nil ***REMOVED***
		return file.Name.Name
	***REMOVED***
	return ""
***REMOVED***

//-------------------------------------------------------------------------
// decl_cache
//
// Thread-safe collection of DeclFileCache entities.
//-------------------------------------------------------------------------

type package_lookup_context struct ***REMOVED***
	build.Context
	BzlProjectRoot     string
	GBProjectRoot      string
	CurrentPackagePath string
***REMOVED***

// gopath returns the list of Go path directories.
func (ctxt *package_lookup_context) gopath() []string ***REMOVED***
	var all []string
	for _, p := range filepath.SplitList(ctxt.GOPATH) ***REMOVED***
		if p == "" || p == ctxt.GOROOT ***REMOVED***
			// Empty paths are uninteresting.
			// If the path is the GOROOT, ignore it.
			// People sometimes set GOPATH=$GOROOT.
			// Do not get confused by this common mistake.
			continue
		***REMOVED***
		if strings.HasPrefix(p, "~") ***REMOVED***
			// Path segments starting with ~ on Unix are almost always
			// users who have incorrectly quoted ~ while setting GOPATH,
			// preventing it from expanding to $HOME.
			// The situation is made more confusing by the fact that
			// bash allows quoted ~ in $PATH (most shells do not).
			// Do not get confused by this, and do not try to use the path.
			// It does not exist, and printing errors about it confuses
			// those users even more, because they think "sure ~ exists!".
			// The go command diagnoses this situation and prints a
			// useful error.
			// On Windows, ~ is used in short names, such as c:\progra~1
			// for c:\program files.
			continue
		***REMOVED***
		all = append(all, p)
	***REMOVED***
	return all
***REMOVED***

func (ctxt *package_lookup_context) pkg_dirs() (string, []string) ***REMOVED***
	pkgdir := fmt.Sprintf("%s_%s", ctxt.GOOS, ctxt.GOARCH)

	var currentPackagePath string
	var all []string
	if ctxt.GOROOT != "" ***REMOVED***
		dir := filepath.Join(ctxt.GOROOT, "pkg", pkgdir)
		if is_dir(dir) ***REMOVED***
			all = append(all, dir)
		***REMOVED***
	***REMOVED***

	switch g_config.PackageLookupMode ***REMOVED***
	case "go":
		currentPackagePath = ctxt.CurrentPackagePath
		for _, p := range ctxt.gopath() ***REMOVED***
			dir := filepath.Join(p, "pkg", pkgdir)
			if is_dir(dir) ***REMOVED***
				all = append(all, dir)
			***REMOVED***
			dir = filepath.Join(dir, currentPackagePath, "vendor")
			if is_dir(dir) ***REMOVED***
				all = append(all, dir)
			***REMOVED***
		***REMOVED***
	case "gb":
		if ctxt.GBProjectRoot != "" ***REMOVED***
			pkgdir := fmt.Sprintf("%s-%s", ctxt.GOOS, ctxt.GOARCH)
			if !is_dir(pkgdir) ***REMOVED***
				pkgdir = fmt.Sprintf("%s-%s-race", ctxt.GOOS, ctxt.GOARCH)
			***REMOVED***
			dir := filepath.Join(ctxt.GBProjectRoot, "pkg", pkgdir)
			if is_dir(dir) ***REMOVED***
				all = append(all, dir)
			***REMOVED***
		***REMOVED***
	case "bzl":
		// TODO: Support bazel mode
	***REMOVED***
	return currentPackagePath, all
***REMOVED***

type decl_cache struct ***REMOVED***
	cache   map[string]*decl_file_cache
	context *package_lookup_context
	sync.Mutex
***REMOVED***

func new_decl_cache(context *package_lookup_context) *decl_cache ***REMOVED***
	return &decl_cache***REMOVED***
		cache:   make(map[string]*decl_file_cache),
		context: context,
	***REMOVED***
***REMOVED***

func (c *decl_cache) get(filename string) *decl_file_cache ***REMOVED***
	c.Lock()
	defer c.Unlock()

	f, ok := c.cache[filename]
	if !ok ***REMOVED***
		f = new_decl_file_cache(filename, c.context)
		c.cache[filename] = f
	***REMOVED***
	return f
***REMOVED***

func (c *decl_cache) get_and_update(filename string) *decl_file_cache ***REMOVED***
	f := c.get(filename)
	f.update()
	return f
***REMOVED***
