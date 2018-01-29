// Copyright 2012 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build !windows

package test

// Session functional tests.

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"strings"
	"testing"

	"golang.org/x/crypto/ssh"
)

func TestRunCommandSuccess(t *testing.T) ***REMOVED***
	server := newServer(t)
	defer server.Shutdown()
	conn := server.Dial(clientConfig())
	defer conn.Close()

	session, err := conn.NewSession()
	if err != nil ***REMOVED***
		t.Fatalf("session failed: %v", err)
	***REMOVED***
	defer session.Close()
	err = session.Run("true")
	if err != nil ***REMOVED***
		t.Fatalf("session failed: %v", err)
	***REMOVED***
***REMOVED***

func TestHostKeyCheck(t *testing.T) ***REMOVED***
	server := newServer(t)
	defer server.Shutdown()

	conf := clientConfig()
	hostDB := hostKeyDB()
	conf.HostKeyCallback = hostDB.Check

	// change the keys.
	hostDB.keys[ssh.KeyAlgoRSA][25]++
	hostDB.keys[ssh.KeyAlgoDSA][25]++
	hostDB.keys[ssh.KeyAlgoECDSA256][25]++

	conn, err := server.TryDial(conf)
	if err == nil ***REMOVED***
		conn.Close()
		t.Fatalf("dial should have failed.")
	***REMOVED*** else if !strings.Contains(err.Error(), "host key mismatch") ***REMOVED***
		t.Fatalf("'host key mismatch' not found in %v", err)
	***REMOVED***
***REMOVED***

func TestRunCommandStdin(t *testing.T) ***REMOVED***
	server := newServer(t)
	defer server.Shutdown()
	conn := server.Dial(clientConfig())
	defer conn.Close()

	session, err := conn.NewSession()
	if err != nil ***REMOVED***
		t.Fatalf("session failed: %v", err)
	***REMOVED***
	defer session.Close()

	r, w := io.Pipe()
	defer r.Close()
	defer w.Close()
	session.Stdin = r

	err = session.Run("true")
	if err != nil ***REMOVED***
		t.Fatalf("session failed: %v", err)
	***REMOVED***
***REMOVED***

func TestRunCommandStdinError(t *testing.T) ***REMOVED***
	server := newServer(t)
	defer server.Shutdown()
	conn := server.Dial(clientConfig())
	defer conn.Close()

	session, err := conn.NewSession()
	if err != nil ***REMOVED***
		t.Fatalf("session failed: %v", err)
	***REMOVED***
	defer session.Close()

	r, w := io.Pipe()
	defer r.Close()
	session.Stdin = r
	pipeErr := errors.New("closing write end of pipe")
	w.CloseWithError(pipeErr)

	err = session.Run("true")
	if err != pipeErr ***REMOVED***
		t.Fatalf("expected %v, found %v", pipeErr, err)
	***REMOVED***
***REMOVED***

func TestRunCommandFailed(t *testing.T) ***REMOVED***
	server := newServer(t)
	defer server.Shutdown()
	conn := server.Dial(clientConfig())
	defer conn.Close()

	session, err := conn.NewSession()
	if err != nil ***REMOVED***
		t.Fatalf("session failed: %v", err)
	***REMOVED***
	defer session.Close()
	err = session.Run(`bash -c "kill -9 $$"`)
	if err == nil ***REMOVED***
		t.Fatalf("session succeeded: %v", err)
	***REMOVED***
***REMOVED***

func TestRunCommandWeClosed(t *testing.T) ***REMOVED***
	server := newServer(t)
	defer server.Shutdown()
	conn := server.Dial(clientConfig())
	defer conn.Close()

	session, err := conn.NewSession()
	if err != nil ***REMOVED***
		t.Fatalf("session failed: %v", err)
	***REMOVED***
	err = session.Shell()
	if err != nil ***REMOVED***
		t.Fatalf("shell failed: %v", err)
	***REMOVED***
	err = session.Close()
	if err != nil ***REMOVED***
		t.Fatalf("shell failed: %v", err)
	***REMOVED***
***REMOVED***

func TestFuncLargeRead(t *testing.T) ***REMOVED***
	server := newServer(t)
	defer server.Shutdown()
	conn := server.Dial(clientConfig())
	defer conn.Close()

	session, err := conn.NewSession()
	if err != nil ***REMOVED***
		t.Fatalf("unable to create new session: %s", err)
	***REMOVED***

	stdout, err := session.StdoutPipe()
	if err != nil ***REMOVED***
		t.Fatalf("unable to acquire stdout pipe: %s", err)
	***REMOVED***

	err = session.Start("dd if=/dev/urandom bs=2048 count=1024")
	if err != nil ***REMOVED***
		t.Fatalf("unable to execute remote command: %s", err)
	***REMOVED***

	buf := new(bytes.Buffer)
	n, err := io.Copy(buf, stdout)
	if err != nil ***REMOVED***
		t.Fatalf("error reading from remote stdout: %s", err)
	***REMOVED***

	if n != 2048*1024 ***REMOVED***
		t.Fatalf("Expected %d bytes but read only %d from remote command", 2048, n)
	***REMOVED***
***REMOVED***

func TestKeyChange(t *testing.T) ***REMOVED***
	server := newServer(t)
	defer server.Shutdown()
	conf := clientConfig()
	hostDB := hostKeyDB()
	conf.HostKeyCallback = hostDB.Check
	conf.RekeyThreshold = 1024
	conn := server.Dial(conf)
	defer conn.Close()

	for i := 0; i < 4; i++ ***REMOVED***
		session, err := conn.NewSession()
		if err != nil ***REMOVED***
			t.Fatalf("unable to create new session: %s", err)
		***REMOVED***

		stdout, err := session.StdoutPipe()
		if err != nil ***REMOVED***
			t.Fatalf("unable to acquire stdout pipe: %s", err)
		***REMOVED***

		err = session.Start("dd if=/dev/urandom bs=1024 count=1")
		if err != nil ***REMOVED***
			t.Fatalf("unable to execute remote command: %s", err)
		***REMOVED***
		buf := new(bytes.Buffer)
		n, err := io.Copy(buf, stdout)
		if err != nil ***REMOVED***
			t.Fatalf("error reading from remote stdout: %s", err)
		***REMOVED***

		want := int64(1024)
		if n != want ***REMOVED***
			t.Fatalf("Expected %d bytes but read only %d from remote command", want, n)
		***REMOVED***
	***REMOVED***

	if changes := hostDB.checkCount; changes < 4 ***REMOVED***
		t.Errorf("got %d key changes, want 4", changes)
	***REMOVED***
***REMOVED***

func TestInvalidTerminalMode(t *testing.T) ***REMOVED***
	server := newServer(t)
	defer server.Shutdown()
	conn := server.Dial(clientConfig())
	defer conn.Close()

	session, err := conn.NewSession()
	if err != nil ***REMOVED***
		t.Fatalf("session failed: %v", err)
	***REMOVED***
	defer session.Close()

	if err = session.RequestPty("vt100", 80, 40, ssh.TerminalModes***REMOVED***255: 1984***REMOVED***); err == nil ***REMOVED***
		t.Fatalf("req-pty failed: successful request with invalid mode")
	***REMOVED***
***REMOVED***

func TestValidTerminalMode(t *testing.T) ***REMOVED***
	server := newServer(t)
	defer server.Shutdown()
	conn := server.Dial(clientConfig())
	defer conn.Close()

	session, err := conn.NewSession()
	if err != nil ***REMOVED***
		t.Fatalf("session failed: %v", err)
	***REMOVED***
	defer session.Close()

	stdout, err := session.StdoutPipe()
	if err != nil ***REMOVED***
		t.Fatalf("unable to acquire stdout pipe: %s", err)
	***REMOVED***

	stdin, err := session.StdinPipe()
	if err != nil ***REMOVED***
		t.Fatalf("unable to acquire stdin pipe: %s", err)
	***REMOVED***

	tm := ssh.TerminalModes***REMOVED***ssh.ECHO: 0***REMOVED***
	if err = session.RequestPty("xterm", 80, 40, tm); err != nil ***REMOVED***
		t.Fatalf("req-pty failed: %s", err)
	***REMOVED***

	err = session.Shell()
	if err != nil ***REMOVED***
		t.Fatalf("session failed: %s", err)
	***REMOVED***

	stdin.Write([]byte("stty -a && exit\n"))

	var buf bytes.Buffer
	if _, err := io.Copy(&buf, stdout); err != nil ***REMOVED***
		t.Fatalf("reading failed: %s", err)
	***REMOVED***

	if sttyOutput := buf.String(); !strings.Contains(sttyOutput, "-echo ") ***REMOVED***
		t.Fatalf("terminal mode failure: expected -echo in stty output, got %s", sttyOutput)
	***REMOVED***
***REMOVED***

func TestWindowChange(t *testing.T) ***REMOVED***
	server := newServer(t)
	defer server.Shutdown()
	conn := server.Dial(clientConfig())
	defer conn.Close()

	session, err := conn.NewSession()
	if err != nil ***REMOVED***
		t.Fatalf("session failed: %v", err)
	***REMOVED***
	defer session.Close()

	stdout, err := session.StdoutPipe()
	if err != nil ***REMOVED***
		t.Fatalf("unable to acquire stdout pipe: %s", err)
	***REMOVED***

	stdin, err := session.StdinPipe()
	if err != nil ***REMOVED***
		t.Fatalf("unable to acquire stdin pipe: %s", err)
	***REMOVED***

	tm := ssh.TerminalModes***REMOVED***ssh.ECHO: 0***REMOVED***
	if err = session.RequestPty("xterm", 80, 40, tm); err != nil ***REMOVED***
		t.Fatalf("req-pty failed: %s", err)
	***REMOVED***

	if err := session.WindowChange(100, 100); err != nil ***REMOVED***
		t.Fatalf("window-change failed: %s", err)
	***REMOVED***

	err = session.Shell()
	if err != nil ***REMOVED***
		t.Fatalf("session failed: %s", err)
	***REMOVED***

	stdin.Write([]byte("stty size && exit\n"))

	var buf bytes.Buffer
	if _, err := io.Copy(&buf, stdout); err != nil ***REMOVED***
		t.Fatalf("reading failed: %s", err)
	***REMOVED***

	if sttyOutput := buf.String(); !strings.Contains(sttyOutput, "100 100") ***REMOVED***
		t.Fatalf("terminal WindowChange failure: expected \"100 100\" stty output, got %s", sttyOutput)
	***REMOVED***
***REMOVED***

func testOneCipher(t *testing.T, cipher string, cipherOrder []string) ***REMOVED***
	server := newServer(t)
	defer server.Shutdown()
	conf := clientConfig()
	conf.Ciphers = []string***REMOVED***cipher***REMOVED***
	// Don't fail if sshd doesn't have the cipher.
	conf.Ciphers = append(conf.Ciphers, cipherOrder...)
	conn, err := server.TryDial(conf)
	if err != nil ***REMOVED***
		t.Fatalf("TryDial: %v", err)
	***REMOVED***
	defer conn.Close()

	numBytes := 4096

	// Exercise sending data to the server
	if _, _, err := conn.Conn.SendRequest("drop-me", false, make([]byte, numBytes)); err != nil ***REMOVED***
		t.Fatalf("SendRequest: %v", err)
	***REMOVED***

	// Exercise receiving data from the server
	session, err := conn.NewSession()
	if err != nil ***REMOVED***
		t.Fatalf("NewSession: %v", err)
	***REMOVED***

	out, err := session.Output(fmt.Sprintf("dd if=/dev/zero of=/dev/stdout bs=%d count=1", numBytes))
	if err != nil ***REMOVED***
		t.Fatalf("Output: %v", err)
	***REMOVED***

	if len(out) != numBytes ***REMOVED***
		t.Fatalf("got %d bytes, want %d bytes", len(out), numBytes)
	***REMOVED***
***REMOVED***

var deprecatedCiphers = []string***REMOVED***
	"aes128-cbc", "3des-cbc",
	"arcfour128", "arcfour256",
***REMOVED***

func TestCiphers(t *testing.T) ***REMOVED***
	var config ssh.Config
	config.SetDefaults()
	cipherOrder := append(config.Ciphers, deprecatedCiphers...)

	for _, ciph := range cipherOrder ***REMOVED***
		t.Run(ciph, func(t *testing.T) ***REMOVED***
			testOneCipher(t, ciph, cipherOrder)
		***REMOVED***)
	***REMOVED***
***REMOVED***

func TestMACs(t *testing.T) ***REMOVED***
	var config ssh.Config
	config.SetDefaults()
	macOrder := config.MACs

	for _, mac := range macOrder ***REMOVED***
		server := newServer(t)
		defer server.Shutdown()
		conf := clientConfig()
		conf.MACs = []string***REMOVED***mac***REMOVED***
		// Don't fail if sshd doesn't have the MAC.
		conf.MACs = append(conf.MACs, macOrder...)
		if conn, err := server.TryDial(conf); err == nil ***REMOVED***
			conn.Close()
		***REMOVED*** else ***REMOVED***
			t.Fatalf("failed for MAC %q", mac)
		***REMOVED***
	***REMOVED***
***REMOVED***

func TestKeyExchanges(t *testing.T) ***REMOVED***
	var config ssh.Config
	config.SetDefaults()
	kexOrder := config.KeyExchanges
	for _, kex := range kexOrder ***REMOVED***
		server := newServer(t)
		defer server.Shutdown()
		conf := clientConfig()
		// Don't fail if sshd doesn't have the kex.
		conf.KeyExchanges = append([]string***REMOVED***kex***REMOVED***, kexOrder...)
		conn, err := server.TryDial(conf)
		if err == nil ***REMOVED***
			conn.Close()
		***REMOVED*** else ***REMOVED***
			t.Errorf("failed for kex %q", kex)
		***REMOVED***
	***REMOVED***
***REMOVED***

func TestClientAuthAlgorithms(t *testing.T) ***REMOVED***
	for _, key := range []string***REMOVED***
		"rsa",
		"dsa",
		"ecdsa",
		"ed25519",
	***REMOVED*** ***REMOVED***
		server := newServer(t)
		conf := clientConfig()
		conf.SetDefaults()
		conf.Auth = []ssh.AuthMethod***REMOVED***
			ssh.PublicKeys(testSigners[key]),
		***REMOVED***

		conn, err := server.TryDial(conf)
		if err == nil ***REMOVED***
			conn.Close()
		***REMOVED*** else ***REMOVED***
			t.Errorf("failed for key %q", key)
		***REMOVED***

		server.Shutdown()
	***REMOVED***
***REMOVED***
