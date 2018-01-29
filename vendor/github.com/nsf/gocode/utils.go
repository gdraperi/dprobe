package main

import (
	"bytes"
	"fmt"
	"go/build"
	"io/ioutil"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
	"unicode/utf8"
)

// our own readdir, which skips the files it cannot lstat
func readdir_lstat(name string) ([]os.FileInfo, error) ***REMOVED***
	f, err := os.Open(name)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	defer f.Close()

	names, err := f.Readdirnames(-1)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	out := make([]os.FileInfo, 0, len(names))
	for _, lname := range names ***REMOVED***
		s, err := os.Lstat(filepath.Join(name, lname))
		if err != nil ***REMOVED***
			continue
		***REMOVED***
		out = append(out, s)
	***REMOVED***
	return out, nil
***REMOVED***

// our other readdir function, only opens and reads
func readdir(dirname string) []os.FileInfo ***REMOVED***
	f, err := os.Open(dirname)
	if err != nil ***REMOVED***
		return nil
	***REMOVED***
	fi, err := f.Readdir(-1)
	f.Close()
	if err != nil ***REMOVED***
		panic(err)
	***REMOVED***
	return fi
***REMOVED***

// returns truncated 'data' and amount of bytes skipped (for cursor pos adjustment)
func filter_out_shebang(data []byte) ([]byte, int) ***REMOVED***
	if len(data) > 2 && data[0] == '#' && data[1] == '!' ***REMOVED***
		newline := bytes.Index(data, []byte("\n"))
		if newline != -1 && len(data) > newline+1 ***REMOVED***
			return data[newline+1:], newline + 1
		***REMOVED***
	***REMOVED***
	return data, 0
***REMOVED***

func file_exists(filename string) bool ***REMOVED***
	_, err := os.Stat(filename)
	if err != nil ***REMOVED***
		return false
	***REMOVED***
	return true
***REMOVED***

func is_dir(path string) bool ***REMOVED***
	fi, err := os.Stat(path)
	return err == nil && fi.IsDir()
***REMOVED***

func char_to_byte_offset(s []byte, offset_c int) (offset_b int) ***REMOVED***
	for offset_b = 0; offset_c > 0 && offset_b < len(s); offset_b++ ***REMOVED***
		if utf8.RuneStart(s[offset_b]) ***REMOVED***
			offset_c--
		***REMOVED***
	***REMOVED***
	return offset_b
***REMOVED***

func xdg_home_dir() string ***REMOVED***
	xdghome := os.Getenv("XDG_CONFIG_HOME")
	if xdghome == "" ***REMOVED***
		xdghome = filepath.Join(os.Getenv("HOME"), ".config")
	***REMOVED***
	return xdghome
***REMOVED***

func has_prefix(s, prefix string, ignorecase bool) bool ***REMOVED***
	if ignorecase ***REMOVED***
		s = strings.ToLower(s)
		prefix = strings.ToLower(prefix)
	***REMOVED***
	return strings.HasPrefix(s, prefix)
***REMOVED***

func find_bzl_project_root(libpath, path string) (string, error) ***REMOVED***
	if libpath == "" ***REMOVED***
		return "", fmt.Errorf("could not find project root, libpath is empty")
	***REMOVED***

	pathMap := map[string]struct***REMOVED******REMOVED******REMOVED******REMOVED***
	for _, lp := range strings.Split(libpath, ":") ***REMOVED***
		lp := strings.TrimSpace(lp)
		pathMap[filepath.Clean(lp)] = struct***REMOVED******REMOVED******REMOVED******REMOVED***
	***REMOVED***

	path = filepath.Dir(path)
	if path == "" ***REMOVED***
		return "", fmt.Errorf("project root is blank")
	***REMOVED***

	start := path
	for path != "/" ***REMOVED***
		if _, ok := pathMap[filepath.Clean(path)]; ok ***REMOVED***
			return path, nil
		***REMOVED***
		path = filepath.Dir(path)
	***REMOVED***
	return "", fmt.Errorf("could not find project root in %q or its parents", start)
***REMOVED***

// Code taken directly from `gb`, I hope author doesn't mind.
func find_gb_project_root(path string) (string, error) ***REMOVED***
	path = filepath.Dir(path)
	if path == "" ***REMOVED***
		return "", fmt.Errorf("project root is blank")
	***REMOVED***
	start := path
	for path != "/" ***REMOVED***
		root := filepath.Join(path, "src")
		if _, err := os.Stat(root); err != nil ***REMOVED***
			if os.IsNotExist(err) ***REMOVED***
				path = filepath.Dir(path)
				continue
			***REMOVED***
			return "", err
		***REMOVED***
		path, err := filepath.EvalSymlinks(path)
		if err != nil ***REMOVED***
			return "", err
		***REMOVED***
		return path, nil
	***REMOVED***
	return "", fmt.Errorf("could not find project root in %q or its parents", start)
***REMOVED***

// vendorlessImportPath returns the devendorized version of the provided import path.
// e.g. "foo/bar/vendor/a/b" => "a/b"
func vendorlessImportPath(ipath string, currentPackagePath string) (string, bool) ***REMOVED***
	split := strings.Split(ipath, "vendor/")
	// no vendor in path
	if len(split) == 1 ***REMOVED***
		return ipath, true
	***REMOVED***
	// this import path does not belong to the current package
	if currentPackagePath != "" && !strings.Contains(currentPackagePath, split[0]) ***REMOVED***
		return "", false
	***REMOVED***
	// Devendorize for use in import statement.
	if i := strings.LastIndex(ipath, "/vendor/"); i >= 0 ***REMOVED***
		return ipath[i+len("/vendor/"):], true
	***REMOVED***
	if strings.HasPrefix(ipath, "vendor/") ***REMOVED***
		return ipath[len("vendor/"):], true
	***REMOVED***
	return ipath, true
***REMOVED***

//-------------------------------------------------------------------------
// print_backtrace
//
// a nicer backtrace printer than the default one
//-------------------------------------------------------------------------

var g_backtrace_mutex sync.Mutex

func print_backtrace(err interface***REMOVED******REMOVED***) ***REMOVED***
	g_backtrace_mutex.Lock()
	defer g_backtrace_mutex.Unlock()
	fmt.Printf("panic: %v\n", err)
	i := 2
	for ***REMOVED***
		pc, file, line, ok := runtime.Caller(i)
		if !ok ***REMOVED***
			break
		***REMOVED***
		f := runtime.FuncForPC(pc)
		fmt.Printf("%d(%s): %s:%d\n", i-1, f.Name(), file, line)
		i++
	***REMOVED***
	fmt.Println("")
***REMOVED***

//-------------------------------------------------------------------------
// File reader goroutine
//
// It's a bad idea to block multiple goroutines on file I/O. Creates many
// threads which fight for HDD. Therefore only single goroutine should read HDD
// at the same time.
//-------------------------------------------------------------------------

type file_read_request struct ***REMOVED***
	filename string
	out      chan file_read_response
***REMOVED***

type file_read_response struct ***REMOVED***
	data  []byte
	error error
***REMOVED***

type file_reader_type struct ***REMOVED***
	in chan file_read_request
***REMOVED***

func new_file_reader() *file_reader_type ***REMOVED***
	this := new(file_reader_type)
	this.in = make(chan file_read_request)
	go func() ***REMOVED***
		var rsp file_read_response
		for ***REMOVED***
			req := <-this.in
			rsp.data, rsp.error = ioutil.ReadFile(req.filename)
			req.out <- rsp
		***REMOVED***
	***REMOVED***()
	return this
***REMOVED***

func (this *file_reader_type) read_file(filename string) ([]byte, error) ***REMOVED***
	req := file_read_request***REMOVED***
		filename,
		make(chan file_read_response),
	***REMOVED***
	this.in <- req
	rsp := <-req.out
	return rsp.data, rsp.error
***REMOVED***

var file_reader = new_file_reader()

//-------------------------------------------------------------------------
// copy of the build.Context without func fields
//-------------------------------------------------------------------------

type go_build_context struct ***REMOVED***
	GOARCH        string
	GOOS          string
	GOROOT        string
	GOPATH        string
	CgoEnabled    bool
	UseAllFiles   bool
	Compiler      string
	BuildTags     []string
	ReleaseTags   []string
	InstallSuffix string
***REMOVED***

func pack_build_context(ctx *build.Context) go_build_context ***REMOVED***
	return go_build_context***REMOVED***
		GOARCH:        ctx.GOARCH,
		GOOS:          ctx.GOOS,
		GOROOT:        ctx.GOROOT,
		GOPATH:        ctx.GOPATH,
		CgoEnabled:    ctx.CgoEnabled,
		UseAllFiles:   ctx.UseAllFiles,
		Compiler:      ctx.Compiler,
		BuildTags:     ctx.BuildTags,
		ReleaseTags:   ctx.ReleaseTags,
		InstallSuffix: ctx.InstallSuffix,
	***REMOVED***
***REMOVED***

func unpack_build_context(ctx *go_build_context) package_lookup_context ***REMOVED***
	return package_lookup_context***REMOVED***
		Context: build.Context***REMOVED***
			GOARCH:        ctx.GOARCH,
			GOOS:          ctx.GOOS,
			GOROOT:        ctx.GOROOT,
			GOPATH:        ctx.GOPATH,
			CgoEnabled:    ctx.CgoEnabled,
			UseAllFiles:   ctx.UseAllFiles,
			Compiler:      ctx.Compiler,
			BuildTags:     ctx.BuildTags,
			ReleaseTags:   ctx.ReleaseTags,
			InstallSuffix: ctx.InstallSuffix,
		***REMOVED***,
	***REMOVED***
***REMOVED***
