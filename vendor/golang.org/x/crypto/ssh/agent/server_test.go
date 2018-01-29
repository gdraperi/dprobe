// Copyright 2012 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package agent

import (
	"crypto"
	"crypto/rand"
	"fmt"
	pseudorand "math/rand"
	"reflect"
	"strings"
	"testing"

	"golang.org/x/crypto/ssh"
)

func TestServer(t *testing.T) ***REMOVED***
	c1, c2, err := netPipe()
	if err != nil ***REMOVED***
		t.Fatalf("netPipe: %v", err)
	***REMOVED***
	defer c1.Close()
	defer c2.Close()
	client := NewClient(c1)

	go ServeAgent(NewKeyring(), c2)

	testAgentInterface(t, client, testPrivateKeys["rsa"], nil, 0)
***REMOVED***

func TestLockServer(t *testing.T) ***REMOVED***
	testLockAgent(NewKeyring(), t)
***REMOVED***

func TestSetupForwardAgent(t *testing.T) ***REMOVED***
	a, b, err := netPipe()
	if err != nil ***REMOVED***
		t.Fatalf("netPipe: %v", err)
	***REMOVED***

	defer a.Close()
	defer b.Close()

	_, socket, cleanup := startOpenSSHAgent(t)
	defer cleanup()

	serverConf := ssh.ServerConfig***REMOVED***
		NoClientAuth: true,
	***REMOVED***
	serverConf.AddHostKey(testSigners["rsa"])
	incoming := make(chan *ssh.ServerConn, 1)
	go func() ***REMOVED***
		conn, _, _, err := ssh.NewServerConn(a, &serverConf)
		if err != nil ***REMOVED***
			t.Fatalf("Server: %v", err)
		***REMOVED***
		incoming <- conn
	***REMOVED***()

	conf := ssh.ClientConfig***REMOVED***
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	***REMOVED***
	conn, chans, reqs, err := ssh.NewClientConn(b, "", &conf)
	if err != nil ***REMOVED***
		t.Fatalf("NewClientConn: %v", err)
	***REMOVED***
	client := ssh.NewClient(conn, chans, reqs)

	if err := ForwardToRemote(client, socket); err != nil ***REMOVED***
		t.Fatalf("SetupForwardAgent: %v", err)
	***REMOVED***

	server := <-incoming
	ch, reqs, err := server.OpenChannel(channelType, nil)
	if err != nil ***REMOVED***
		t.Fatalf("OpenChannel(%q): %v", channelType, err)
	***REMOVED***
	go ssh.DiscardRequests(reqs)

	agentClient := NewClient(ch)
	testAgentInterface(t, agentClient, testPrivateKeys["rsa"], nil, 0)
	conn.Close()
***REMOVED***

func TestV1ProtocolMessages(t *testing.T) ***REMOVED***
	c1, c2, err := netPipe()
	if err != nil ***REMOVED***
		t.Fatalf("netPipe: %v", err)
	***REMOVED***
	defer c1.Close()
	defer c2.Close()
	c := NewClient(c1)

	go ServeAgent(NewKeyring(), c2)

	testV1ProtocolMessages(t, c.(*client))
***REMOVED***

func testV1ProtocolMessages(t *testing.T, c *client) ***REMOVED***
	reply, err := c.call([]byte***REMOVED***agentRequestV1Identities***REMOVED***)
	if err != nil ***REMOVED***
		t.Fatalf("v1 request all failed: %v", err)
	***REMOVED***
	if msg, ok := reply.(*agentV1IdentityMsg); !ok || msg.Numkeys != 0 ***REMOVED***
		t.Fatalf("invalid request all response: %#v", reply)
	***REMOVED***

	reply, err = c.call([]byte***REMOVED***agentRemoveAllV1Identities***REMOVED***)
	if err != nil ***REMOVED***
		t.Fatalf("v1 remove all failed: %v", err)
	***REMOVED***
	if _, ok := reply.(*successAgentMsg); !ok ***REMOVED***
		t.Fatalf("invalid remove all response: %#v", reply)
	***REMOVED***
***REMOVED***

func verifyKey(sshAgent Agent) error ***REMOVED***
	keys, err := sshAgent.List()
	if err != nil ***REMOVED***
		return fmt.Errorf("listing keys: %v", err)
	***REMOVED***

	if len(keys) != 1 ***REMOVED***
		return fmt.Errorf("bad number of keys found. expected 1, got %d", len(keys))
	***REMOVED***

	buf := make([]byte, 128)
	if _, err := rand.Read(buf); err != nil ***REMOVED***
		return fmt.Errorf("rand: %v", err)
	***REMOVED***

	sig, err := sshAgent.Sign(keys[0], buf)
	if err != nil ***REMOVED***
		return fmt.Errorf("sign: %v", err)
	***REMOVED***

	if err := keys[0].Verify(buf, sig); err != nil ***REMOVED***
		return fmt.Errorf("verify: %v", err)
	***REMOVED***
	return nil
***REMOVED***

func addKeyToAgent(key crypto.PrivateKey) error ***REMOVED***
	sshAgent := NewKeyring()
	if err := sshAgent.Add(AddedKey***REMOVED***PrivateKey: key***REMOVED***); err != nil ***REMOVED***
		return fmt.Errorf("add: %v", err)
	***REMOVED***
	return verifyKey(sshAgent)
***REMOVED***

func TestKeyTypes(t *testing.T) ***REMOVED***
	for k, v := range testPrivateKeys ***REMOVED***
		if err := addKeyToAgent(v); err != nil ***REMOVED***
			t.Errorf("error adding key type %s, %v", k, err)
		***REMOVED***
		if err := addCertToAgentSock(v, nil); err != nil ***REMOVED***
			t.Errorf("error adding key type %s, %v", k, err)
		***REMOVED***
	***REMOVED***
***REMOVED***

func addCertToAgentSock(key crypto.PrivateKey, cert *ssh.Certificate) error ***REMOVED***
	a, b, err := netPipe()
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	agentServer := NewKeyring()
	go ServeAgent(agentServer, a)

	agentClient := NewClient(b)
	if err := agentClient.Add(AddedKey***REMOVED***PrivateKey: key, Certificate: cert***REMOVED***); err != nil ***REMOVED***
		return fmt.Errorf("add: %v", err)
	***REMOVED***
	return verifyKey(agentClient)
***REMOVED***

func addCertToAgent(key crypto.PrivateKey, cert *ssh.Certificate) error ***REMOVED***
	sshAgent := NewKeyring()
	if err := sshAgent.Add(AddedKey***REMOVED***PrivateKey: key, Certificate: cert***REMOVED***); err != nil ***REMOVED***
		return fmt.Errorf("add: %v", err)
	***REMOVED***
	return verifyKey(sshAgent)
***REMOVED***

func TestCertTypes(t *testing.T) ***REMOVED***
	for keyType, key := range testPublicKeys ***REMOVED***
		cert := &ssh.Certificate***REMOVED***
			ValidPrincipals: []string***REMOVED***"gopher1"***REMOVED***,
			ValidAfter:      0,
			ValidBefore:     ssh.CertTimeInfinity,
			Key:             key,
			Serial:          1,
			CertType:        ssh.UserCert,
			SignatureKey:    testPublicKeys["rsa"],
			Permissions: ssh.Permissions***REMOVED***
				CriticalOptions: map[string]string***REMOVED******REMOVED***,
				Extensions:      map[string]string***REMOVED******REMOVED***,
			***REMOVED***,
		***REMOVED***
		if err := cert.SignCert(rand.Reader, testSigners["rsa"]); err != nil ***REMOVED***
			t.Fatalf("signcert: %v", err)
		***REMOVED***
		if err := addCertToAgent(testPrivateKeys[keyType], cert); err != nil ***REMOVED***
			t.Fatalf("%v", err)
		***REMOVED***
		if err := addCertToAgentSock(testPrivateKeys[keyType], cert); err != nil ***REMOVED***
			t.Fatalf("%v", err)
		***REMOVED***
	***REMOVED***
***REMOVED***

func TestParseConstraints(t *testing.T) ***REMOVED***
	// Test LifetimeSecs
	var msg = constrainLifetimeAgentMsg***REMOVED***pseudorand.Uint32()***REMOVED***
	lifetimeSecs, _, _, err := parseConstraints(ssh.Marshal(msg))
	if err != nil ***REMOVED***
		t.Fatalf("parseConstraints: %v", err)
	***REMOVED***
	if lifetimeSecs != msg.LifetimeSecs ***REMOVED***
		t.Errorf("got lifetime %v, want %v", lifetimeSecs, msg.LifetimeSecs)
	***REMOVED***

	// Test ConfirmBeforeUse
	_, confirmBeforeUse, _, err := parseConstraints([]byte***REMOVED***agentConstrainConfirm***REMOVED***)
	if err != nil ***REMOVED***
		t.Fatalf("%v", err)
	***REMOVED***
	if !confirmBeforeUse ***REMOVED***
		t.Error("got comfirmBeforeUse == false")
	***REMOVED***

	// Test ConstraintExtensions
	var data []byte
	var expect []ConstraintExtension
	for i := 0; i < 10; i++ ***REMOVED***
		var ext = ConstraintExtension***REMOVED***
			ExtensionName:    fmt.Sprintf("name%d", i),
			ExtensionDetails: []byte(fmt.Sprintf("details: %d", i)),
		***REMOVED***
		expect = append(expect, ext)
		data = append(data, agentConstrainExtension)
		data = append(data, ssh.Marshal(ext)...)
	***REMOVED***
	_, _, extensions, err := parseConstraints(data)
	if err != nil ***REMOVED***
		t.Fatalf("%v", err)
	***REMOVED***
	if !reflect.DeepEqual(expect, extensions) ***REMOVED***
		t.Errorf("got extension %v, want %v", extensions, expect)
	***REMOVED***

	// Test Unknown Constraint
	_, _, _, err = parseConstraints([]byte***REMOVED***128***REMOVED***)
	if err == nil || !strings.Contains(err.Error(), "unknown constraint") ***REMOVED***
		t.Errorf("unexpected error: %v", err)
	***REMOVED***
***REMOVED***
