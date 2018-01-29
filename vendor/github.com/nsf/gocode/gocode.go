package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	_ "net/http/pprof"
	"os"
	"path/filepath"
)

var (
	g_is_server = flag.Bool("s", false, "run a server instead of a client")
	g_format    = flag.String("f", "nice", "output format (vim | emacs | nice | csv | csv-with-package | json)")
	g_input     = flag.String("in", "", "use this file instead of stdin input")
	g_sock      = create_sock_flag("sock", "socket type (unix | tcp)")
	g_addr      = flag.String("addr", "127.0.0.1:37373", "address for tcp socket")
	g_debug     = flag.Bool("debug", false, "enable server-side debug mode")
	g_profile   = flag.Int("profile", 0, "port on which to expose profiling information for pprof; 0 to disable profiling")
)

func get_socket_filename() string ***REMOVED***
	user := os.Getenv("USER")
	if user == "" ***REMOVED***
		user = "all"
	***REMOVED***
	return filepath.Join(os.TempDir(), fmt.Sprintf("gocode-daemon.%s", user))
***REMOVED***

func show_usage() ***REMOVED***
	fmt.Fprintf(os.Stderr,
		"Usage: %s [-s] [-f=<format>] [-in=<path>] [-sock=<type>] [-addr=<addr>]\n"+
			"       <command> [<args>]\n\n",
		os.Args[0])
	fmt.Fprintf(os.Stderr,
		"Flags:\n")
	flag.PrintDefaults()
	fmt.Fprintf(os.Stderr,
		"\nCommands:\n"+
			"  autocomplete [<path>] <offset>     main autocompletion command\n"+
			"  close                              close the gocode daemon\n"+
			"  drop-cache                         drop gocode daemon's cache\n"+
			"  options                            list config options (extended)\n"+
			"  set [<name> [<value>]]             list or set config options\n"+
			"  status                             gocode daemon status report\n"+
			"")
***REMOVED***

func main() ***REMOVED***
	flag.Usage = show_usage
	flag.Parse()

	var retval int
	if *g_is_server ***REMOVED***
		go func() ***REMOVED***
			if *g_profile <= 0 ***REMOVED***
				return
			***REMOVED***
			addr := fmt.Sprintf("localhost:%d", *g_profile)
			// Use the following commands to profile the binary:
			// go tool pprof http://localhost:6060/debug/pprof/profile   # 30-second CPU profile
			// go tool pprof http://localhost:6060/debug/pprof/heap      # heap profile
			// go tool pprof http://localhost:6060/debug/pprof/block     # goroutine blocking profile
			// See http://blog.golang.org/profiling-go-programs for more info.
			log.Printf("enabling  profiler on %s", addr)
			log.Print(http.ListenAndServe(addr, nil))
		***REMOVED***()
		retval = do_server()
	***REMOVED*** else ***REMOVED***
		retval = do_client()
	***REMOVED***
	os.Exit(retval)
***REMOVED***
