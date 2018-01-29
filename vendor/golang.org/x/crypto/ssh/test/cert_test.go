// Copyright 2014 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build darwin dragonfly freebsd linux netbsd openbsd

package test

import (
	"bytes"
	"crypto/rand"
	"testing"

	"golang.org/x/crypto/ssh"
)

// Test both logging in with a cert, and also that the certificate presented by an OpenSSH host can be validated correctly
func TestCertLogin(t *testing.T) ***REMOVED***
	s := newServer(t)
	defer s.Shutdown()

	// Use a key different from the default.
	clientKey := testSigners["dsa"]
	caAuthKey := testSigners["ecdsa"]
	cert := &ssh.Certificate***REMOVED***
		Key:             clientKey.PublicKey(),
		ValidPrincipals: []string***REMOVED***username()***REMOVED***,
		CertType:        ssh.UserCert,
		ValidBefore:     ssh.CertTimeInfinity,
	***REMOVED***
	if err := cert.SignCert(rand.Reader, caAuthKey); err != nil ***REMOVED***
		t.Fatalf("SetSignature: %v", err)
	***REMOVED***

	certSigner, err := ssh.NewCertSigner(cert, clientKey)
	if err != nil ***REMOVED***
		t.Fatalf("NewCertSigner: %v", err)
	***REMOVED***

	conf := &ssh.ClientConfig***REMOVED***
		User: username(),
		HostKeyCallback: (&ssh.CertChecker***REMOVED***
			IsHostAuthority: func(pk ssh.PublicKey, addr string) bool ***REMOVED***
				return bytes.Equal(pk.Marshal(), testPublicKeys["ca"].Marshal())
			***REMOVED***,
		***REMOVED***).CheckHostKey,
	***REMOVED***
	conf.Auth = append(conf.Auth, ssh.PublicKeys(certSigner))

	for _, test := range []struct ***REMOVED***
		addr    string
		succeed bool
	***REMOVED******REMOVED***
		***REMOVED***addr: "host.example.com:22", succeed: true***REMOVED***,
		***REMOVED***addr: "host.example.com:10000", succeed: true***REMOVED***, // non-standard port must be OK
		***REMOVED***addr: "host.example.com", succeed: false***REMOVED***,      // port must be specified
		***REMOVED***addr: "host.ex4mple.com:22", succeed: false***REMOVED***,   // wrong host
	***REMOVED*** ***REMOVED***
		client, err := s.TryDialWithAddr(conf, test.addr)

		// Always close client if opened successfully
		if err == nil ***REMOVED***
			client.Close()
		***REMOVED***

		// Now evaluate whether the test failed or passed
		if test.succeed ***REMOVED***
			if err != nil ***REMOVED***
				t.Fatalf("TryDialWithAddr: %v", err)
			***REMOVED***
		***REMOVED*** else ***REMOVED***
			if err == nil ***REMOVED***
				t.Fatalf("TryDialWithAddr, unexpected success")
			***REMOVED***
		***REMOVED***
	***REMOVED***
***REMOVED***
