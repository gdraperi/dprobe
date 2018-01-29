// Copyright 2017 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package autocert_test

import (
	"crypto/tls"
	"fmt"
	"log"
	"net/http"

	"golang.org/x/crypto/acme/autocert"
)

func ExampleNewListener() ***REMOVED***
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) ***REMOVED***
		fmt.Fprintf(w, "Hello, TLS user! Your config: %+v", r.TLS)
	***REMOVED***)
	log.Fatal(http.Serve(autocert.NewListener("example.com"), mux))
***REMOVED***

func ExampleManager() ***REMOVED***
	m := &autocert.Manager***REMOVED***
		Cache:      autocert.DirCache("secret-dir"),
		Prompt:     autocert.AcceptTOS,
		HostPolicy: autocert.HostWhitelist("example.org"),
	***REMOVED***
	go http.ListenAndServe(":http", m.HTTPHandler(nil))
	s := &http.Server***REMOVED***
		Addr:      ":https",
		TLSConfig: &tls.Config***REMOVED***GetCertificate: m.GetCertificate***REMOVED***,
	***REMOVED***
	s.ListenAndServeTLS("", "")
***REMOVED***
