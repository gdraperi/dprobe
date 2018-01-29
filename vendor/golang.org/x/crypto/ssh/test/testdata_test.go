// Copyright 2014 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// IMPLEMENTATION NOTE: To avoid a package loop, this file is in three places:
// ssh/, ssh/agent, and ssh/test/. It should be kept in sync across all three
// instances.

package test

import (
	"crypto/rand"
	"fmt"

	"golang.org/x/crypto/ssh"
	"golang.org/x/crypto/ssh/testdata"
)

var (
	testPrivateKeys map[string]interface***REMOVED******REMOVED***
	testSigners     map[string]ssh.Signer
	testPublicKeys  map[string]ssh.PublicKey
)

func init() ***REMOVED***
	var err error

	n := len(testdata.PEMBytes)
	testPrivateKeys = make(map[string]interface***REMOVED******REMOVED***, n)
	testSigners = make(map[string]ssh.Signer, n)
	testPublicKeys = make(map[string]ssh.PublicKey, n)
	for t, k := range testdata.PEMBytes ***REMOVED***
		testPrivateKeys[t], err = ssh.ParseRawPrivateKey(k)
		if err != nil ***REMOVED***
			panic(fmt.Sprintf("Unable to parse test key %s: %v", t, err))
		***REMOVED***
		testSigners[t], err = ssh.NewSignerFromKey(testPrivateKeys[t])
		if err != nil ***REMOVED***
			panic(fmt.Sprintf("Unable to create signer for test key %s: %v", t, err))
		***REMOVED***
		testPublicKeys[t] = testSigners[t].PublicKey()
	***REMOVED***

	// Create a cert and sign it for use in tests.
	testCert := &ssh.Certificate***REMOVED***
		Nonce:           []byte***REMOVED******REMOVED***,                       // To pass reflect.DeepEqual after marshal & parse, this must be non-nil
		ValidPrincipals: []string***REMOVED***"gopher1", "gopher2"***REMOVED***, // increases test coverage
		ValidAfter:      0,                              // unix epoch
		ValidBefore:     ssh.CertTimeInfinity,           // The end of currently representable time.
		Reserved:        []byte***REMOVED******REMOVED***,                       // To pass reflect.DeepEqual after marshal & parse, this must be non-nil
		Key:             testPublicKeys["ecdsa"],
		SignatureKey:    testPublicKeys["rsa"],
		Permissions: ssh.Permissions***REMOVED***
			CriticalOptions: map[string]string***REMOVED******REMOVED***,
			Extensions:      map[string]string***REMOVED******REMOVED***,
		***REMOVED***,
	***REMOVED***
	testCert.SignCert(rand.Reader, testSigners["rsa"])
	testPrivateKeys["cert"] = testPrivateKeys["ecdsa"]
	testSigners["cert"], err = ssh.NewCertSigner(testCert, testSigners["ecdsa"])
	if err != nil ***REMOVED***
		panic(fmt.Sprintf("Unable to create certificate signer: %v", err))
	***REMOVED***
***REMOVED***
