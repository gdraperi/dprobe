package main

import (
	"fmt"
	"net"
	"net/http"
	"os"
	"path/filepath"
)

func main() ***REMOVED***
	p, err := filepath.Abs(filepath.Join("run", "docker", "plugins"))
	if err != nil ***REMOVED***
		panic(err)
	***REMOVED***
	if err := os.MkdirAll(p, 0755); err != nil ***REMOVED***
		panic(err)
	***REMOVED***
	l, err := net.Listen("unix", filepath.Join(p, "basic.sock"))
	if err != nil ***REMOVED***
		panic(err)
	***REMOVED***

	mux := http.NewServeMux()
	server := http.Server***REMOVED***
		Addr:    l.Addr().String(),
		Handler: http.NewServeMux(),
	***REMOVED***
	mux.HandleFunc("/Plugin.Activate", func(w http.ResponseWriter, r *http.Request) ***REMOVED***
		w.Header().Set("Content-Type", "application/vnd.docker.plugins.v1.1+json")
		fmt.Println(w, `***REMOVED***"Implements": ["dummy"]***REMOVED***`)
	***REMOVED***)
	server.Serve(l)
***REMOVED***
