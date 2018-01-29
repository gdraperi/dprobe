// Copyright 2012 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package agent

import (
	"bytes"
	"crypto/rand"
	"errors"
	"net"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"testing"
	"time"

	"golang.org/x/crypto/ssh"
)

// startOpenSSHAgent executes ssh-agent, and returns an Agent interface to it.
func startOpenSSHAgent(t *testing.T) (client Agent, socket string, cleanup func()) ***REMOVED***
	if testing.Short() ***REMOVED***
		// ssh-agent is not always available, and the key
		// types supported vary by platform.
		t.Skip("skipping test due to -short")
	***REMOVED***

	bin, err := exec.LookPath("ssh-agent")
	if err != nil ***REMOVED***
		t.Skip("could not find ssh-agent")
	***REMOVED***

	cmd := exec.Command(bin, "-s")
	out, err := cmd.Output()
	if err != nil ***REMOVED***
		t.Fatalf("cmd.Output: %v", err)
	***REMOVED***

	/* Output looks like:

		   SSH_AUTH_SOCK=/tmp/ssh-P65gpcqArqvH/agent.15541; export SSH_AUTH_SOCK;
	           SSH_AGENT_PID=15542; export SSH_AGENT_PID;
	           echo Agent pid 15542;
	*/
	fields := bytes.Split(out, []byte(";"))
	line := bytes.SplitN(fields[0], []byte("="), 2)
	line[0] = bytes.TrimLeft(line[0], "\n")
	if string(line[0]) != "SSH_AUTH_SOCK" ***REMOVED***
		t.Fatalf("could not find key SSH_AUTH_SOCK in %q", fields[0])
	***REMOVED***
	socket = string(line[1])

	line = bytes.SplitN(fields[2], []byte("="), 2)
	line[0] = bytes.TrimLeft(line[0], "\n")
	if string(line[0]) != "SSH_AGENT_PID" ***REMOVED***
		t.Fatalf("could not find key SSH_AGENT_PID in %q", fields[2])
	***REMOVED***
	pidStr := line[1]
	pid, err := strconv.Atoi(string(pidStr))
	if err != nil ***REMOVED***
		t.Fatalf("Atoi(%q): %v", pidStr, err)
	***REMOVED***

	conn, err := net.Dial("unix", string(socket))
	if err != nil ***REMOVED***
		t.Fatalf("net.Dial: %v", err)
	***REMOVED***

	ac := NewClient(conn)
	return ac, socket, func() ***REMOVED***
		proc, _ := os.FindProcess(pid)
		if proc != nil ***REMOVED***
			proc.Kill()
		***REMOVED***
		conn.Close()
		os.RemoveAll(filepath.Dir(socket))
	***REMOVED***
***REMOVED***

// startKeyringAgent uses Keyring to simulate a ssh-agent Server and returns a client.
func startKeyringAgent(t *testing.T) (client Agent, cleanup func()) ***REMOVED***
	c1, c2, err := netPipe()
	if err != nil ***REMOVED***
		t.Fatalf("netPipe: %v", err)
	***REMOVED***
	go ServeAgent(NewKeyring(), c2)

	return NewClient(c1), func() ***REMOVED***
		c1.Close()
		c2.Close()
	***REMOVED***
***REMOVED***

func testOpenSSHAgent(t *testing.T, key interface***REMOVED******REMOVED***, cert *ssh.Certificate, lifetimeSecs uint32) ***REMOVED***
	agent, _, cleanup := startOpenSSHAgent(t)
	defer cleanup()

	testAgentInterface(t, agent, key, cert, lifetimeSecs)
***REMOVED***

func testKeyringAgent(t *testing.T, key interface***REMOVED******REMOVED***, cert *ssh.Certificate, lifetimeSecs uint32) ***REMOVED***
	agent, cleanup := startKeyringAgent(t)
	defer cleanup()

	testAgentInterface(t, agent, key, cert, lifetimeSecs)
***REMOVED***

func testAgentInterface(t *testing.T, agent Agent, key interface***REMOVED******REMOVED***, cert *ssh.Certificate, lifetimeSecs uint32) ***REMOVED***
	signer, err := ssh.NewSignerFromKey(key)
	if err != nil ***REMOVED***
		t.Fatalf("NewSignerFromKey(%T): %v", key, err)
	***REMOVED***
	// The agent should start up empty.
	if keys, err := agent.List(); err != nil ***REMOVED***
		t.Fatalf("RequestIdentities: %v", err)
	***REMOVED*** else if len(keys) > 0 ***REMOVED***
		t.Fatalf("got %d keys, want 0: %v", len(keys), keys)
	***REMOVED***

	// Attempt to insert the key, with certificate if specified.
	var pubKey ssh.PublicKey
	if cert != nil ***REMOVED***
		err = agent.Add(AddedKey***REMOVED***
			PrivateKey:   key,
			Certificate:  cert,
			Comment:      "comment",
			LifetimeSecs: lifetimeSecs,
		***REMOVED***)
		pubKey = cert
	***REMOVED*** else ***REMOVED***
		err = agent.Add(AddedKey***REMOVED***PrivateKey: key, Comment: "comment", LifetimeSecs: lifetimeSecs***REMOVED***)
		pubKey = signer.PublicKey()
	***REMOVED***
	if err != nil ***REMOVED***
		t.Fatalf("insert(%T): %v", key, err)
	***REMOVED***

	// Did the key get inserted successfully?
	if keys, err := agent.List(); err != nil ***REMOVED***
		t.Fatalf("List: %v", err)
	***REMOVED*** else if len(keys) != 1 ***REMOVED***
		t.Fatalf("got %v, want 1 key", keys)
	***REMOVED*** else if keys[0].Comment != "comment" ***REMOVED***
		t.Fatalf("key comment: got %v, want %v", keys[0].Comment, "comment")
	***REMOVED*** else if !bytes.Equal(keys[0].Blob, pubKey.Marshal()) ***REMOVED***
		t.Fatalf("key mismatch")
	***REMOVED***

	// Can the agent make a valid signature?
	data := []byte("hello")
	sig, err := agent.Sign(pubKey, data)
	if err != nil ***REMOVED***
		t.Fatalf("Sign(%s): %v", pubKey.Type(), err)
	***REMOVED***

	if err := pubKey.Verify(data, sig); err != nil ***REMOVED***
		t.Fatalf("Verify(%s): %v", pubKey.Type(), err)
	***REMOVED***

	// If the key has a lifetime, is it removed when it should be?
	if lifetimeSecs > 0 ***REMOVED***
		time.Sleep(time.Second*time.Duration(lifetimeSecs) + 100*time.Millisecond)
		keys, err := agent.List()
		if err != nil ***REMOVED***
			t.Fatalf("List: %v", err)
		***REMOVED***
		if len(keys) > 0 ***REMOVED***
			t.Fatalf("key not expired")
		***REMOVED***
	***REMOVED***

***REMOVED***

func TestAgent(t *testing.T) ***REMOVED***
	for _, keyType := range []string***REMOVED***"rsa", "dsa", "ecdsa", "ed25519"***REMOVED*** ***REMOVED***
		testOpenSSHAgent(t, testPrivateKeys[keyType], nil, 0)
		testKeyringAgent(t, testPrivateKeys[keyType], nil, 0)
	***REMOVED***
***REMOVED***

func TestCert(t *testing.T) ***REMOVED***
	cert := &ssh.Certificate***REMOVED***
		Key:         testPublicKeys["rsa"],
		ValidBefore: ssh.CertTimeInfinity,
		CertType:    ssh.UserCert,
	***REMOVED***
	cert.SignCert(rand.Reader, testSigners["ecdsa"])

	testOpenSSHAgent(t, testPrivateKeys["rsa"], cert, 0)
	testKeyringAgent(t, testPrivateKeys["rsa"], cert, 0)
***REMOVED***

// netPipe is analogous to net.Pipe, but it uses a real net.Conn, and
// therefore is buffered (net.Pipe deadlocks if both sides start with
// a write.)
func netPipe() (net.Conn, net.Conn, error) ***REMOVED***
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil ***REMOVED***
		listener, err = net.Listen("tcp", "[::1]:0")
		if err != nil ***REMOVED***
			return nil, nil, err
		***REMOVED***
	***REMOVED***
	defer listener.Close()
	c1, err := net.Dial("tcp", listener.Addr().String())
	if err != nil ***REMOVED***
		return nil, nil, err
	***REMOVED***

	c2, err := listener.Accept()
	if err != nil ***REMOVED***
		c1.Close()
		return nil, nil, err
	***REMOVED***

	return c1, c2, nil
***REMOVED***

func TestAuth(t *testing.T) ***REMOVED***
	agent, _, cleanup := startOpenSSHAgent(t)
	defer cleanup()

	a, b, err := netPipe()
	if err != nil ***REMOVED***
		t.Fatalf("netPipe: %v", err)
	***REMOVED***

	defer a.Close()
	defer b.Close()

	if err := agent.Add(AddedKey***REMOVED***PrivateKey: testPrivateKeys["rsa"], Comment: "comment"***REMOVED***); err != nil ***REMOVED***
		t.Errorf("Add: %v", err)
	***REMOVED***

	serverConf := ssh.ServerConfig***REMOVED******REMOVED***
	serverConf.AddHostKey(testSigners["rsa"])
	serverConf.PublicKeyCallback = func(c ssh.ConnMetadata, key ssh.PublicKey) (*ssh.Permissions, error) ***REMOVED***
		if bytes.Equal(key.Marshal(), testPublicKeys["rsa"].Marshal()) ***REMOVED***
			return nil, nil
		***REMOVED***

		return nil, errors.New("pubkey rejected")
	***REMOVED***

	go func() ***REMOVED***
		conn, _, _, err := ssh.NewServerConn(a, &serverConf)
		if err != nil ***REMOVED***
			t.Fatalf("Server: %v", err)
		***REMOVED***
		conn.Close()
	***REMOVED***()

	conf := ssh.ClientConfig***REMOVED***
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	***REMOVED***
	conf.Auth = append(conf.Auth, ssh.PublicKeysCallback(agent.Signers))
	conn, _, _, err := ssh.NewClientConn(b, "", &conf)
	if err != nil ***REMOVED***
		t.Fatalf("NewClientConn: %v", err)
	***REMOVED***
	conn.Close()
***REMOVED***

func TestLockOpenSSHAgent(t *testing.T) ***REMOVED***
	agent, _, cleanup := startOpenSSHAgent(t)
	defer cleanup()
	testLockAgent(agent, t)
***REMOVED***

func TestLockKeyringAgent(t *testing.T) ***REMOVED***
	agent, cleanup := startKeyringAgent(t)
	defer cleanup()
	testLockAgent(agent, t)
***REMOVED***

func testLockAgent(agent Agent, t *testing.T) ***REMOVED***
	if err := agent.Add(AddedKey***REMOVED***PrivateKey: testPrivateKeys["rsa"], Comment: "comment 1"***REMOVED***); err != nil ***REMOVED***
		t.Errorf("Add: %v", err)
	***REMOVED***
	if err := agent.Add(AddedKey***REMOVED***PrivateKey: testPrivateKeys["dsa"], Comment: "comment dsa"***REMOVED***); err != nil ***REMOVED***
		t.Errorf("Add: %v", err)
	***REMOVED***
	if keys, err := agent.List(); err != nil ***REMOVED***
		t.Errorf("List: %v", err)
	***REMOVED*** else if len(keys) != 2 ***REMOVED***
		t.Errorf("Want 2 keys, got %v", keys)
	***REMOVED***

	passphrase := []byte("secret")
	if err := agent.Lock(passphrase); err != nil ***REMOVED***
		t.Errorf("Lock: %v", err)
	***REMOVED***

	if keys, err := agent.List(); err != nil ***REMOVED***
		t.Errorf("List: %v", err)
	***REMOVED*** else if len(keys) != 0 ***REMOVED***
		t.Errorf("Want 0 keys, got %v", keys)
	***REMOVED***

	signer, _ := ssh.NewSignerFromKey(testPrivateKeys["rsa"])
	if _, err := agent.Sign(signer.PublicKey(), []byte("hello")); err == nil ***REMOVED***
		t.Fatalf("Sign did not fail")
	***REMOVED***

	if err := agent.Remove(signer.PublicKey()); err == nil ***REMOVED***
		t.Fatalf("Remove did not fail")
	***REMOVED***

	if err := agent.RemoveAll(); err == nil ***REMOVED***
		t.Fatalf("RemoveAll did not fail")
	***REMOVED***

	if err := agent.Unlock(nil); err == nil ***REMOVED***
		t.Errorf("Unlock with wrong passphrase succeeded")
	***REMOVED***
	if err := agent.Unlock(passphrase); err != nil ***REMOVED***
		t.Errorf("Unlock: %v", err)
	***REMOVED***

	if err := agent.Remove(signer.PublicKey()); err != nil ***REMOVED***
		t.Fatalf("Remove: %v", err)
	***REMOVED***

	if keys, err := agent.List(); err != nil ***REMOVED***
		t.Errorf("List: %v", err)
	***REMOVED*** else if len(keys) != 1 ***REMOVED***
		t.Errorf("Want 1 keys, got %v", keys)
	***REMOVED***
***REMOVED***

func testOpenSSHAgentLifetime(t *testing.T) ***REMOVED***
	agent, _, cleanup := startOpenSSHAgent(t)
	defer cleanup()
	testAgentLifetime(t, agent)
***REMOVED***

func testKeyringAgentLifetime(t *testing.T) ***REMOVED***
	agent, cleanup := startKeyringAgent(t)
	defer cleanup()
	testAgentLifetime(t, agent)
***REMOVED***

func testAgentLifetime(t *testing.T, agent Agent) ***REMOVED***
	for _, keyType := range []string***REMOVED***"rsa", "dsa", "ecdsa"***REMOVED*** ***REMOVED***
		// Add private keys to the agent.
		err := agent.Add(AddedKey***REMOVED***
			PrivateKey:   testPrivateKeys[keyType],
			Comment:      "comment",
			LifetimeSecs: 1,
		***REMOVED***)
		if err != nil ***REMOVED***
			t.Fatalf("add: %v", err)
		***REMOVED***
		// Add certs to the agent.
		cert := &ssh.Certificate***REMOVED***
			Key:         testPublicKeys[keyType],
			ValidBefore: ssh.CertTimeInfinity,
			CertType:    ssh.UserCert,
		***REMOVED***
		cert.SignCert(rand.Reader, testSigners[keyType])
		err = agent.Add(AddedKey***REMOVED***
			PrivateKey:   testPrivateKeys[keyType],
			Certificate:  cert,
			Comment:      "comment",
			LifetimeSecs: 1,
		***REMOVED***)
		if err != nil ***REMOVED***
			t.Fatalf("add: %v", err)
		***REMOVED***
	***REMOVED***
	time.Sleep(1100 * time.Millisecond)
	if keys, err := agent.List(); err != nil ***REMOVED***
		t.Errorf("List: %v", err)
	***REMOVED*** else if len(keys) != 0 ***REMOVED***
		t.Errorf("Want 0 keys, got %v", len(keys))
	***REMOVED***
***REMOVED***
