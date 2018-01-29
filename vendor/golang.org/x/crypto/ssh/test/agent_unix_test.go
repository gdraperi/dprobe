// Copyright 2014 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build darwin dragonfly freebsd linux netbsd openbsd

package test

import (
	"bytes"
	"testing"

	"golang.org/x/crypto/ssh"
	"golang.org/x/crypto/ssh/agent"
)

func TestAgentForward(t *testing.T) ***REMOVED***
	server := newServer(t)
	defer server.Shutdown()
	conn := server.Dial(clientConfig())
	defer conn.Close()

	keyring := agent.NewKeyring()
	if err := keyring.Add(agent.AddedKey***REMOVED***PrivateKey: testPrivateKeys["dsa"]***REMOVED***); err != nil ***REMOVED***
		t.Fatalf("Error adding key: %s", err)
	***REMOVED***
	if err := keyring.Add(agent.AddedKey***REMOVED***
		PrivateKey:       testPrivateKeys["dsa"],
		ConfirmBeforeUse: true,
		LifetimeSecs:     3600,
	***REMOVED***); err != nil ***REMOVED***
		t.Fatalf("Error adding key with constraints: %s", err)
	***REMOVED***
	pub := testPublicKeys["dsa"]

	sess, err := conn.NewSession()
	if err != nil ***REMOVED***
		t.Fatalf("NewSession: %v", err)
	***REMOVED***
	if err := agent.RequestAgentForwarding(sess); err != nil ***REMOVED***
		t.Fatalf("RequestAgentForwarding: %v", err)
	***REMOVED***

	if err := agent.ForwardToAgent(conn, keyring); err != nil ***REMOVED***
		t.Fatalf("SetupForwardKeyring: %v", err)
	***REMOVED***
	out, err := sess.CombinedOutput("ssh-add -L")
	if err != nil ***REMOVED***
		t.Fatalf("running ssh-add: %v, out %s", err, out)
	***REMOVED***
	key, _, _, _, err := ssh.ParseAuthorizedKey(out)
	if err != nil ***REMOVED***
		t.Fatalf("ParseAuthorizedKey(%q): %v", out, err)
	***REMOVED***

	if !bytes.Equal(key.Marshal(), pub.Marshal()) ***REMOVED***
		t.Fatalf("got key %s, want %s", ssh.MarshalAuthorizedKey(key), ssh.MarshalAuthorizedKey(pub))
	***REMOVED***
***REMOVED***
