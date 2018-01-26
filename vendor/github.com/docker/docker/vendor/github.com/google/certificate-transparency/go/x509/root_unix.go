// Copyright 2011 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build dragonfly freebsd linux openbsd netbsd solaris

package x509

import "io/ioutil"

// Possible certificate files; stop after finding one.
var certFiles = []string***REMOVED***
	"/etc/ssl/certs/ca-certificates.crt",     // Debian/Ubuntu/Gentoo etc.
	"/etc/pki/tls/certs/ca-bundle.crt",       // Fedora/RHEL
	"/etc/ssl/ca-bundle.pem",                 // OpenSUSE
	"/etc/ssl/cert.pem",                      // OpenBSD
	"/usr/local/share/certs/ca-root-nss.crt", // FreeBSD/DragonFly
***REMOVED***

func (c *Certificate) systemVerify(opts *VerifyOptions) (chains [][]*Certificate, err error) ***REMOVED***
	return nil, nil
***REMOVED***

func initSystemRoots() ***REMOVED***
	roots := NewCertPool()
	for _, file := range certFiles ***REMOVED***
		data, err := ioutil.ReadFile(file)
		if err == nil ***REMOVED***
			roots.AppendCertsFromPEM(data)
			systemRoots = roots
			return
		***REMOVED***
	***REMOVED***

	// All of the files failed to load. systemRoots will be nil which will
	// trigger a specific error at verification time.
***REMOVED***
