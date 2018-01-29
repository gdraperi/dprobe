package main

import (
	"bytes"
	"fmt"
	"go/build"
	"log"
	"net"
	"net/rpc"
	"os"
	"path/filepath"
	"reflect"
	"runtime"
	"time"
)

func do_server() int ***REMOVED***
	g_config.read()
	if g_config.ForceDebugOutput != "" ***REMOVED***
		// forcefully enable debugging and redirect logging into the
		// specified file
		*g_debug = true
		f, err := os.Create(g_config.ForceDebugOutput)
		if err != nil ***REMOVED***
			panic(err)
		***REMOVED***
		log.SetOutput(f)
	***REMOVED***

	addr := *g_addr
	if *g_sock == "unix" ***REMOVED***
		addr = get_socket_filename()
		if file_exists(addr) ***REMOVED***
			log.Printf("unix socket: '%s' already exists\n", addr)
			return 1
		***REMOVED***
	***REMOVED***
	g_daemon = new_daemon(*g_sock, addr)
	if *g_sock == "unix" ***REMOVED***
		// cleanup unix socket file
		defer os.Remove(addr)
	***REMOVED***

	rpc.Register(new(RPC))

	g_daemon.loop()
	return 0
***REMOVED***

//-------------------------------------------------------------------------
// daemon
//-------------------------------------------------------------------------

type daemon struct ***REMOVED***
	listener     net.Listener
	cmd_in       chan int
	autocomplete *auto_complete_context
	pkgcache     package_cache
	declcache    *decl_cache
	context      package_lookup_context
***REMOVED***

func new_daemon(network, address string) *daemon ***REMOVED***
	var err error

	d := new(daemon)
	d.listener, err = net.Listen(network, address)
	if err != nil ***REMOVED***
		panic(err)
	***REMOVED***

	d.cmd_in = make(chan int, 1)
	d.pkgcache = new_package_cache()
	d.declcache = new_decl_cache(&d.context)
	d.autocomplete = new_auto_complete_context(d.pkgcache, d.declcache)
	return d
***REMOVED***

func (this *daemon) drop_cache() ***REMOVED***
	this.pkgcache = new_package_cache()
	this.declcache = new_decl_cache(&this.context)
	this.autocomplete = new_auto_complete_context(this.pkgcache, this.declcache)
***REMOVED***

const (
	daemon_close = iota
)

func (this *daemon) loop() ***REMOVED***
	conn_in := make(chan net.Conn)
	go func() ***REMOVED***
		for ***REMOVED***
			c, err := this.listener.Accept()
			if err != nil ***REMOVED***
				panic(err)
			***REMOVED***
			conn_in <- c
		***REMOVED***
	***REMOVED***()

	timeout := time.Duration(g_config.CloseTimeout) * time.Second
	countdown := time.NewTimer(timeout)

	for ***REMOVED***
		// handle connections or server CMDs (currently one CMD)
		select ***REMOVED***
		case c := <-conn_in:
			rpc.ServeConn(c)
			countdown.Reset(timeout)
			runtime.GC()
		case cmd := <-this.cmd_in:
			switch cmd ***REMOVED***
			case daemon_close:
				return
			***REMOVED***
		case <-countdown.C:
			return
		***REMOVED***
	***REMOVED***
***REMOVED***

func (this *daemon) close() ***REMOVED***
	this.cmd_in <- daemon_close
***REMOVED***

var g_daemon *daemon

//-------------------------------------------------------------------------
// server_* functions
//
// Corresponding client_* functions are autogenerated by goremote.
//-------------------------------------------------------------------------

func server_auto_complete(file []byte, filename string, cursor int, context_packed go_build_context) (c []candidate, d int) ***REMOVED***
	context := unpack_build_context(&context_packed)
	defer func() ***REMOVED***
		if err := recover(); err != nil ***REMOVED***
			print_backtrace(err)
			c = []candidate***REMOVED***
				***REMOVED***"PANIC", "PANIC", decl_invalid, "panic"***REMOVED***,
			***REMOVED***

			// drop cache
			g_daemon.drop_cache()
		***REMOVED***
	***REMOVED***()
	// TODO: Probably we don't care about comparing all the fields, checking GOROOT and GOPATH
	// should be enough.
	if !reflect.DeepEqual(g_daemon.context.Context, context.Context) ***REMOVED***
		g_daemon.context = context
		g_daemon.drop_cache()
	***REMOVED***
	switch g_config.PackageLookupMode ***REMOVED***
	case "bzl":
		// when package lookup mode is bzl, we set GOPATH to "" explicitly and
		// BzlProjectRoot becomes valid (or empty)
		var err error
		g_daemon.context.GOPATH = ""
		g_daemon.context.BzlProjectRoot, err = find_bzl_project_root(g_config.LibPath, filename)
		if *g_debug && err != nil ***REMOVED***
			log.Printf("Bzl project root not found: %s", err)
		***REMOVED***
	case "gb":
		// when package lookup mode is gb, we set GOPATH to "" explicitly and
		// GBProjectRoot becomes valid (or empty)
		var err error
		g_daemon.context.GOPATH = ""
		g_daemon.context.GBProjectRoot, err = find_gb_project_root(filename)
		if *g_debug && err != nil ***REMOVED***
			log.Printf("Gb project root not found: %s", err)
		***REMOVED***
	case "go":
		// get current package path for GO15VENDOREXPERIMENT hack
		g_daemon.context.CurrentPackagePath = ""
		pkg, err := g_daemon.context.ImportDir(filepath.Dir(filename), build.FindOnly)
		if err == nil ***REMOVED***
			if *g_debug ***REMOVED***
				log.Printf("Go project path: %s", pkg.ImportPath)
			***REMOVED***
			g_daemon.context.CurrentPackagePath = pkg.ImportPath
		***REMOVED*** else if *g_debug ***REMOVED***
			log.Printf("Go project path not found: %s", err)
		***REMOVED***
	***REMOVED***
	if *g_debug ***REMOVED***
		var buf bytes.Buffer
		log.Printf("Got autocompletion request for '%s'\n", filename)
		log.Printf("Cursor at: %d\n", cursor)
		if cursor > len(file) || cursor < 0 ***REMOVED***
			log.Println("ERROR! Cursor is outside of the boundaries of the buffer, " +
				"this is most likely a text editor plugin bug. Text editor is responsible " +
				"for passing the correct cursor position to gocode.")
		***REMOVED*** else ***REMOVED***
			buf.WriteString("-------------------------------------------------------\n")
			buf.Write(file[:cursor])
			buf.WriteString("#")
			buf.Write(file[cursor:])
			log.Print(buf.String())
			log.Println("-------------------------------------------------------")
		***REMOVED***
	***REMOVED***
	candidates, d := g_daemon.autocomplete.apropos(file, filename, cursor)
	if *g_debug ***REMOVED***
		log.Printf("Offset: %d\n", d)
		log.Printf("Number of candidates found: %d\n", len(candidates))
		log.Printf("Candidates are:\n")
		for _, c := range candidates ***REMOVED***
			abbr := fmt.Sprintf("%s %s %s", c.Class, c.Name, c.Type)
			if c.Class == decl_func ***REMOVED***
				abbr = fmt.Sprintf("%s %s%s", c.Class, c.Name, c.Type[len("func"):])
			***REMOVED***
			log.Printf("  %s\n", abbr)
		***REMOVED***
		log.Println("=======================================================")
	***REMOVED***
	return candidates, d
***REMOVED***

func server_close(notused int) int ***REMOVED***
	g_daemon.close()
	return 0
***REMOVED***

func server_status(notused int) string ***REMOVED***
	return g_daemon.autocomplete.status()
***REMOVED***

func server_drop_cache(notused int) int ***REMOVED***
	// drop cache
	g_daemon.drop_cache()
	return 0
***REMOVED***

func server_set(key, value string) string ***REMOVED***
	if key == "\x00" ***REMOVED***
		return g_config.list()
	***REMOVED*** else if value == "\x00" ***REMOVED***
		return g_config.list_option(key)
	***REMOVED***
	// drop cache on settings changes
	g_daemon.drop_cache()
	return g_config.set_option(key, value)
***REMOVED***

func server_options(notused int) string ***REMOVED***
	return g_config.options()
***REMOVED***
