// Copyright 2011 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package ssh

// Session tests.

import (
	"bytes"
	crypto_rand "crypto/rand"
	"errors"
	"io"
	"io/ioutil"
	"math/rand"
	"net"
	"testing"

	"golang.org/x/crypto/ssh/terminal"
)

type serverType func(Channel, <-chan *Request, *testing.T)

// dial constructs a new test server and returns a *ClientConn.
func dial(handler serverType, t *testing.T) *Client ***REMOVED***
	c1, c2, err := netPipe()
	if err != nil ***REMOVED***
		t.Fatalf("netPipe: %v", err)
	***REMOVED***

	go func() ***REMOVED***
		defer c1.Close()
		conf := ServerConfig***REMOVED***
			NoClientAuth: true,
		***REMOVED***
		conf.AddHostKey(testSigners["rsa"])

		_, chans, reqs, err := NewServerConn(c1, &conf)
		if err != nil ***REMOVED***
			t.Fatalf("Unable to handshake: %v", err)
		***REMOVED***
		go DiscardRequests(reqs)

		for newCh := range chans ***REMOVED***
			if newCh.ChannelType() != "session" ***REMOVED***
				newCh.Reject(UnknownChannelType, "unknown channel type")
				continue
			***REMOVED***

			ch, inReqs, err := newCh.Accept()
			if err != nil ***REMOVED***
				t.Errorf("Accept: %v", err)
				continue
			***REMOVED***
			go func() ***REMOVED***
				handler(ch, inReqs, t)
			***REMOVED***()
		***REMOVED***
	***REMOVED***()

	config := &ClientConfig***REMOVED***
		User:            "testuser",
		HostKeyCallback: InsecureIgnoreHostKey(),
	***REMOVED***

	conn, chans, reqs, err := NewClientConn(c2, "", config)
	if err != nil ***REMOVED***
		t.Fatalf("unable to dial remote side: %v", err)
	***REMOVED***

	return NewClient(conn, chans, reqs)
***REMOVED***

// Test a simple string is returned to session.Stdout.
func TestSessionShell(t *testing.T) ***REMOVED***
	conn := dial(shellHandler, t)
	defer conn.Close()
	session, err := conn.NewSession()
	if err != nil ***REMOVED***
		t.Fatalf("Unable to request new session: %v", err)
	***REMOVED***
	defer session.Close()
	stdout := new(bytes.Buffer)
	session.Stdout = stdout
	if err := session.Shell(); err != nil ***REMOVED***
		t.Fatalf("Unable to execute command: %s", err)
	***REMOVED***
	if err := session.Wait(); err != nil ***REMOVED***
		t.Fatalf("Remote command did not exit cleanly: %v", err)
	***REMOVED***
	actual := stdout.String()
	if actual != "golang" ***REMOVED***
		t.Fatalf("Remote shell did not return expected string: expected=golang, actual=%s", actual)
	***REMOVED***
***REMOVED***

// TODO(dfc) add support for Std***REMOVED***in,err***REMOVED***Pipe when the Server supports it.

// Test a simple string is returned via StdoutPipe.
func TestSessionStdoutPipe(t *testing.T) ***REMOVED***
	conn := dial(shellHandler, t)
	defer conn.Close()
	session, err := conn.NewSession()
	if err != nil ***REMOVED***
		t.Fatalf("Unable to request new session: %v", err)
	***REMOVED***
	defer session.Close()
	stdout, err := session.StdoutPipe()
	if err != nil ***REMOVED***
		t.Fatalf("Unable to request StdoutPipe(): %v", err)
	***REMOVED***
	var buf bytes.Buffer
	if err := session.Shell(); err != nil ***REMOVED***
		t.Fatalf("Unable to execute command: %v", err)
	***REMOVED***
	done := make(chan bool, 1)
	go func() ***REMOVED***
		if _, err := io.Copy(&buf, stdout); err != nil ***REMOVED***
			t.Errorf("Copy of stdout failed: %v", err)
		***REMOVED***
		done <- true
	***REMOVED***()
	if err := session.Wait(); err != nil ***REMOVED***
		t.Fatalf("Remote command did not exit cleanly: %v", err)
	***REMOVED***
	<-done
	actual := buf.String()
	if actual != "golang" ***REMOVED***
		t.Fatalf("Remote shell did not return expected string: expected=golang, actual=%s", actual)
	***REMOVED***
***REMOVED***

// Test that a simple string is returned via the Output helper,
// and that stderr is discarded.
func TestSessionOutput(t *testing.T) ***REMOVED***
	conn := dial(fixedOutputHandler, t)
	defer conn.Close()
	session, err := conn.NewSession()
	if err != nil ***REMOVED***
		t.Fatalf("Unable to request new session: %v", err)
	***REMOVED***
	defer session.Close()

	buf, err := session.Output("") // cmd is ignored by fixedOutputHandler
	if err != nil ***REMOVED***
		t.Error("Remote command did not exit cleanly:", err)
	***REMOVED***
	w := "this-is-stdout."
	g := string(buf)
	if g != w ***REMOVED***
		t.Error("Remote command did not return expected string:")
		t.Logf("want %q", w)
		t.Logf("got  %q", g)
	***REMOVED***
***REMOVED***

// Test that both stdout and stderr are returned
// via the CombinedOutput helper.
func TestSessionCombinedOutput(t *testing.T) ***REMOVED***
	conn := dial(fixedOutputHandler, t)
	defer conn.Close()
	session, err := conn.NewSession()
	if err != nil ***REMOVED***
		t.Fatalf("Unable to request new session: %v", err)
	***REMOVED***
	defer session.Close()

	buf, err := session.CombinedOutput("") // cmd is ignored by fixedOutputHandler
	if err != nil ***REMOVED***
		t.Error("Remote command did not exit cleanly:", err)
	***REMOVED***
	const stdout = "this-is-stdout."
	const stderr = "this-is-stderr."
	g := string(buf)
	if g != stdout+stderr && g != stderr+stdout ***REMOVED***
		t.Error("Remote command did not return expected string:")
		t.Logf("want %q, or %q", stdout+stderr, stderr+stdout)
		t.Logf("got  %q", g)
	***REMOVED***
***REMOVED***

// Test non-0 exit status is returned correctly.
func TestExitStatusNonZero(t *testing.T) ***REMOVED***
	conn := dial(exitStatusNonZeroHandler, t)
	defer conn.Close()
	session, err := conn.NewSession()
	if err != nil ***REMOVED***
		t.Fatalf("Unable to request new session: %v", err)
	***REMOVED***
	defer session.Close()
	if err := session.Shell(); err != nil ***REMOVED***
		t.Fatalf("Unable to execute command: %v", err)
	***REMOVED***
	err = session.Wait()
	if err == nil ***REMOVED***
		t.Fatalf("expected command to fail but it didn't")
	***REMOVED***
	e, ok := err.(*ExitError)
	if !ok ***REMOVED***
		t.Fatalf("expected *ExitError but got %T", err)
	***REMOVED***
	if e.ExitStatus() != 15 ***REMOVED***
		t.Fatalf("expected command to exit with 15 but got %v", e.ExitStatus())
	***REMOVED***
***REMOVED***

// Test 0 exit status is returned correctly.
func TestExitStatusZero(t *testing.T) ***REMOVED***
	conn := dial(exitStatusZeroHandler, t)
	defer conn.Close()
	session, err := conn.NewSession()
	if err != nil ***REMOVED***
		t.Fatalf("Unable to request new session: %v", err)
	***REMOVED***
	defer session.Close()

	if err := session.Shell(); err != nil ***REMOVED***
		t.Fatalf("Unable to execute command: %v", err)
	***REMOVED***
	err = session.Wait()
	if err != nil ***REMOVED***
		t.Fatalf("expected nil but got %v", err)
	***REMOVED***
***REMOVED***

// Test exit signal and status are both returned correctly.
func TestExitSignalAndStatus(t *testing.T) ***REMOVED***
	conn := dial(exitSignalAndStatusHandler, t)
	defer conn.Close()
	session, err := conn.NewSession()
	if err != nil ***REMOVED***
		t.Fatalf("Unable to request new session: %v", err)
	***REMOVED***
	defer session.Close()
	if err := session.Shell(); err != nil ***REMOVED***
		t.Fatalf("Unable to execute command: %v", err)
	***REMOVED***
	err = session.Wait()
	if err == nil ***REMOVED***
		t.Fatalf("expected command to fail but it didn't")
	***REMOVED***
	e, ok := err.(*ExitError)
	if !ok ***REMOVED***
		t.Fatalf("expected *ExitError but got %T", err)
	***REMOVED***
	if e.Signal() != "TERM" || e.ExitStatus() != 15 ***REMOVED***
		t.Fatalf("expected command to exit with signal TERM and status 15 but got signal %s and status %v", e.Signal(), e.ExitStatus())
	***REMOVED***
***REMOVED***

// Test exit signal and status are both returned correctly.
func TestKnownExitSignalOnly(t *testing.T) ***REMOVED***
	conn := dial(exitSignalHandler, t)
	defer conn.Close()
	session, err := conn.NewSession()
	if err != nil ***REMOVED***
		t.Fatalf("Unable to request new session: %v", err)
	***REMOVED***
	defer session.Close()
	if err := session.Shell(); err != nil ***REMOVED***
		t.Fatalf("Unable to execute command: %v", err)
	***REMOVED***
	err = session.Wait()
	if err == nil ***REMOVED***
		t.Fatalf("expected command to fail but it didn't")
	***REMOVED***
	e, ok := err.(*ExitError)
	if !ok ***REMOVED***
		t.Fatalf("expected *ExitError but got %T", err)
	***REMOVED***
	if e.Signal() != "TERM" || e.ExitStatus() != 143 ***REMOVED***
		t.Fatalf("expected command to exit with signal TERM and status 143 but got signal %s and status %v", e.Signal(), e.ExitStatus())
	***REMOVED***
***REMOVED***

// Test exit signal and status are both returned correctly.
func TestUnknownExitSignal(t *testing.T) ***REMOVED***
	conn := dial(exitSignalUnknownHandler, t)
	defer conn.Close()
	session, err := conn.NewSession()
	if err != nil ***REMOVED***
		t.Fatalf("Unable to request new session: %v", err)
	***REMOVED***
	defer session.Close()
	if err := session.Shell(); err != nil ***REMOVED***
		t.Fatalf("Unable to execute command: %v", err)
	***REMOVED***
	err = session.Wait()
	if err == nil ***REMOVED***
		t.Fatalf("expected command to fail but it didn't")
	***REMOVED***
	e, ok := err.(*ExitError)
	if !ok ***REMOVED***
		t.Fatalf("expected *ExitError but got %T", err)
	***REMOVED***
	if e.Signal() != "SYS" || e.ExitStatus() != 128 ***REMOVED***
		t.Fatalf("expected command to exit with signal SYS and status 128 but got signal %s and status %v", e.Signal(), e.ExitStatus())
	***REMOVED***
***REMOVED***

func TestExitWithoutStatusOrSignal(t *testing.T) ***REMOVED***
	conn := dial(exitWithoutSignalOrStatus, t)
	defer conn.Close()
	session, err := conn.NewSession()
	if err != nil ***REMOVED***
		t.Fatalf("Unable to request new session: %v", err)
	***REMOVED***
	defer session.Close()
	if err := session.Shell(); err != nil ***REMOVED***
		t.Fatalf("Unable to execute command: %v", err)
	***REMOVED***
	err = session.Wait()
	if err == nil ***REMOVED***
		t.Fatalf("expected command to fail but it didn't")
	***REMOVED***
	if _, ok := err.(*ExitMissingError); !ok ***REMOVED***
		t.Fatalf("got %T want *ExitMissingError", err)
	***REMOVED***
***REMOVED***

// windowTestBytes is the number of bytes that we'll send to the SSH server.
const windowTestBytes = 16000 * 200

// TestServerWindow writes random data to the server. The server is expected to echo
// the same data back, which is compared against the original.
func TestServerWindow(t *testing.T) ***REMOVED***
	origBuf := bytes.NewBuffer(make([]byte, 0, windowTestBytes))
	io.CopyN(origBuf, crypto_rand.Reader, windowTestBytes)
	origBytes := origBuf.Bytes()

	conn := dial(echoHandler, t)
	defer conn.Close()
	session, err := conn.NewSession()
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
	defer session.Close()
	result := make(chan []byte)

	go func() ***REMOVED***
		defer close(result)
		echoedBuf := bytes.NewBuffer(make([]byte, 0, windowTestBytes))
		serverStdout, err := session.StdoutPipe()
		if err != nil ***REMOVED***
			t.Errorf("StdoutPipe failed: %v", err)
			return
		***REMOVED***
		n, err := copyNRandomly("stdout", echoedBuf, serverStdout, windowTestBytes)
		if err != nil && err != io.EOF ***REMOVED***
			t.Errorf("Read only %d bytes from server, expected %d: %v", n, windowTestBytes, err)
		***REMOVED***
		result <- echoedBuf.Bytes()
	***REMOVED***()

	serverStdin, err := session.StdinPipe()
	if err != nil ***REMOVED***
		t.Fatalf("StdinPipe failed: %v", err)
	***REMOVED***
	written, err := copyNRandomly("stdin", serverStdin, origBuf, windowTestBytes)
	if err != nil ***REMOVED***
		t.Fatalf("failed to copy origBuf to serverStdin: %v", err)
	***REMOVED***
	if written != windowTestBytes ***REMOVED***
		t.Fatalf("Wrote only %d of %d bytes to server", written, windowTestBytes)
	***REMOVED***

	echoedBytes := <-result

	if !bytes.Equal(origBytes, echoedBytes) ***REMOVED***
		t.Fatalf("Echoed buffer differed from original, orig %d, echoed %d", len(origBytes), len(echoedBytes))
	***REMOVED***
***REMOVED***

// Verify the client can handle a keepalive packet from the server.
func TestClientHandlesKeepalives(t *testing.T) ***REMOVED***
	conn := dial(channelKeepaliveSender, t)
	defer conn.Close()
	session, err := conn.NewSession()
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
	defer session.Close()
	if err := session.Shell(); err != nil ***REMOVED***
		t.Fatalf("Unable to execute command: %v", err)
	***REMOVED***
	err = session.Wait()
	if err != nil ***REMOVED***
		t.Fatalf("expected nil but got: %v", err)
	***REMOVED***
***REMOVED***

type exitStatusMsg struct ***REMOVED***
	Status uint32
***REMOVED***

type exitSignalMsg struct ***REMOVED***
	Signal     string
	CoreDumped bool
	Errmsg     string
	Lang       string
***REMOVED***

func handleTerminalRequests(in <-chan *Request) ***REMOVED***
	for req := range in ***REMOVED***
		ok := false
		switch req.Type ***REMOVED***
		case "shell":
			ok = true
			if len(req.Payload) > 0 ***REMOVED***
				// We don't accept any commands, only the default shell.
				ok = false
			***REMOVED***
		case "env":
			ok = true
		***REMOVED***
		req.Reply(ok, nil)
	***REMOVED***
***REMOVED***

func newServerShell(ch Channel, in <-chan *Request, prompt string) *terminal.Terminal ***REMOVED***
	term := terminal.NewTerminal(ch, prompt)
	go handleTerminalRequests(in)
	return term
***REMOVED***

func exitStatusZeroHandler(ch Channel, in <-chan *Request, t *testing.T) ***REMOVED***
	defer ch.Close()
	// this string is returned to stdout
	shell := newServerShell(ch, in, "> ")
	readLine(shell, t)
	sendStatus(0, ch, t)
***REMOVED***

func exitStatusNonZeroHandler(ch Channel, in <-chan *Request, t *testing.T) ***REMOVED***
	defer ch.Close()
	shell := newServerShell(ch, in, "> ")
	readLine(shell, t)
	sendStatus(15, ch, t)
***REMOVED***

func exitSignalAndStatusHandler(ch Channel, in <-chan *Request, t *testing.T) ***REMOVED***
	defer ch.Close()
	shell := newServerShell(ch, in, "> ")
	readLine(shell, t)
	sendStatus(15, ch, t)
	sendSignal("TERM", ch, t)
***REMOVED***

func exitSignalHandler(ch Channel, in <-chan *Request, t *testing.T) ***REMOVED***
	defer ch.Close()
	shell := newServerShell(ch, in, "> ")
	readLine(shell, t)
	sendSignal("TERM", ch, t)
***REMOVED***

func exitSignalUnknownHandler(ch Channel, in <-chan *Request, t *testing.T) ***REMOVED***
	defer ch.Close()
	shell := newServerShell(ch, in, "> ")
	readLine(shell, t)
	sendSignal("SYS", ch, t)
***REMOVED***

func exitWithoutSignalOrStatus(ch Channel, in <-chan *Request, t *testing.T) ***REMOVED***
	defer ch.Close()
	shell := newServerShell(ch, in, "> ")
	readLine(shell, t)
***REMOVED***

func shellHandler(ch Channel, in <-chan *Request, t *testing.T) ***REMOVED***
	defer ch.Close()
	// this string is returned to stdout
	shell := newServerShell(ch, in, "golang")
	readLine(shell, t)
	sendStatus(0, ch, t)
***REMOVED***

// Ignores the command, writes fixed strings to stderr and stdout.
// Strings are "this-is-stdout." and "this-is-stderr.".
func fixedOutputHandler(ch Channel, in <-chan *Request, t *testing.T) ***REMOVED***
	defer ch.Close()
	_, err := ch.Read(nil)

	req, ok := <-in
	if !ok ***REMOVED***
		t.Fatalf("error: expected channel request, got: %#v", err)
		return
	***REMOVED***

	// ignore request, always send some text
	req.Reply(true, nil)

	_, err = io.WriteString(ch, "this-is-stdout.")
	if err != nil ***REMOVED***
		t.Fatalf("error writing on server: %v", err)
	***REMOVED***
	_, err = io.WriteString(ch.Stderr(), "this-is-stderr.")
	if err != nil ***REMOVED***
		t.Fatalf("error writing on server: %v", err)
	***REMOVED***
	sendStatus(0, ch, t)
***REMOVED***

func readLine(shell *terminal.Terminal, t *testing.T) ***REMOVED***
	if _, err := shell.ReadLine(); err != nil && err != io.EOF ***REMOVED***
		t.Errorf("unable to read line: %v", err)
	***REMOVED***
***REMOVED***

func sendStatus(status uint32, ch Channel, t *testing.T) ***REMOVED***
	msg := exitStatusMsg***REMOVED***
		Status: status,
	***REMOVED***
	if _, err := ch.SendRequest("exit-status", false, Marshal(&msg)); err != nil ***REMOVED***
		t.Errorf("unable to send status: %v", err)
	***REMOVED***
***REMOVED***

func sendSignal(signal string, ch Channel, t *testing.T) ***REMOVED***
	sig := exitSignalMsg***REMOVED***
		Signal:     signal,
		CoreDumped: false,
		Errmsg:     "Process terminated",
		Lang:       "en-GB-oed",
	***REMOVED***
	if _, err := ch.SendRequest("exit-signal", false, Marshal(&sig)); err != nil ***REMOVED***
		t.Errorf("unable to send signal: %v", err)
	***REMOVED***
***REMOVED***

func discardHandler(ch Channel, t *testing.T) ***REMOVED***
	defer ch.Close()
	io.Copy(ioutil.Discard, ch)
***REMOVED***

func echoHandler(ch Channel, in <-chan *Request, t *testing.T) ***REMOVED***
	defer ch.Close()
	if n, err := copyNRandomly("echohandler", ch, ch, windowTestBytes); err != nil ***REMOVED***
		t.Errorf("short write, wrote %d, expected %d: %v ", n, windowTestBytes, err)
	***REMOVED***
***REMOVED***

// copyNRandomly copies n bytes from src to dst. It uses a variable, and random,
// buffer size to exercise more code paths.
func copyNRandomly(title string, dst io.Writer, src io.Reader, n int) (int, error) ***REMOVED***
	var (
		buf       = make([]byte, 32*1024)
		written   int
		remaining = n
	)
	for remaining > 0 ***REMOVED***
		l := rand.Intn(1 << 15)
		if remaining < l ***REMOVED***
			l = remaining
		***REMOVED***
		nr, er := src.Read(buf[:l])
		nw, ew := dst.Write(buf[:nr])
		remaining -= nw
		written += nw
		if ew != nil ***REMOVED***
			return written, ew
		***REMOVED***
		if nr != nw ***REMOVED***
			return written, io.ErrShortWrite
		***REMOVED***
		if er != nil && er != io.EOF ***REMOVED***
			return written, er
		***REMOVED***
	***REMOVED***
	return written, nil
***REMOVED***

func channelKeepaliveSender(ch Channel, in <-chan *Request, t *testing.T) ***REMOVED***
	defer ch.Close()
	shell := newServerShell(ch, in, "> ")
	readLine(shell, t)
	if _, err := ch.SendRequest("keepalive@openssh.com", true, nil); err != nil ***REMOVED***
		t.Errorf("unable to send channel keepalive request: %v", err)
	***REMOVED***
	sendStatus(0, ch, t)
***REMOVED***

func TestClientWriteEOF(t *testing.T) ***REMOVED***
	conn := dial(simpleEchoHandler, t)
	defer conn.Close()

	session, err := conn.NewSession()
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
	defer session.Close()
	stdin, err := session.StdinPipe()
	if err != nil ***REMOVED***
		t.Fatalf("StdinPipe failed: %v", err)
	***REMOVED***
	stdout, err := session.StdoutPipe()
	if err != nil ***REMOVED***
		t.Fatalf("StdoutPipe failed: %v", err)
	***REMOVED***

	data := []byte(`0000`)
	_, err = stdin.Write(data)
	if err != nil ***REMOVED***
		t.Fatalf("Write failed: %v", err)
	***REMOVED***
	stdin.Close()

	res, err := ioutil.ReadAll(stdout)
	if err != nil ***REMOVED***
		t.Fatalf("Read failed: %v", err)
	***REMOVED***

	if !bytes.Equal(data, res) ***REMOVED***
		t.Fatalf("Read differed from write, wrote: %v, read: %v", data, res)
	***REMOVED***
***REMOVED***

func simpleEchoHandler(ch Channel, in <-chan *Request, t *testing.T) ***REMOVED***
	defer ch.Close()
	data, err := ioutil.ReadAll(ch)
	if err != nil ***REMOVED***
		t.Errorf("handler read error: %v", err)
	***REMOVED***
	_, err = ch.Write(data)
	if err != nil ***REMOVED***
		t.Errorf("handler write error: %v", err)
	***REMOVED***
***REMOVED***

func TestSessionID(t *testing.T) ***REMOVED***
	c1, c2, err := netPipe()
	if err != nil ***REMOVED***
		t.Fatalf("netPipe: %v", err)
	***REMOVED***
	defer c1.Close()
	defer c2.Close()

	serverID := make(chan []byte, 1)
	clientID := make(chan []byte, 1)

	serverConf := &ServerConfig***REMOVED***
		NoClientAuth: true,
	***REMOVED***
	serverConf.AddHostKey(testSigners["ecdsa"])
	clientConf := &ClientConfig***REMOVED***
		HostKeyCallback: InsecureIgnoreHostKey(),
		User:            "user",
	***REMOVED***

	go func() ***REMOVED***
		conn, chans, reqs, err := NewServerConn(c1, serverConf)
		if err != nil ***REMOVED***
			t.Fatalf("server handshake: %v", err)
		***REMOVED***
		serverID <- conn.SessionID()
		go DiscardRequests(reqs)
		for ch := range chans ***REMOVED***
			ch.Reject(Prohibited, "")
		***REMOVED***
	***REMOVED***()

	go func() ***REMOVED***
		conn, chans, reqs, err := NewClientConn(c2, "", clientConf)
		if err != nil ***REMOVED***
			t.Fatalf("client handshake: %v", err)
		***REMOVED***
		clientID <- conn.SessionID()
		go DiscardRequests(reqs)
		for ch := range chans ***REMOVED***
			ch.Reject(Prohibited, "")
		***REMOVED***
	***REMOVED***()

	s := <-serverID
	c := <-clientID
	if bytes.Compare(s, c) != 0 ***REMOVED***
		t.Errorf("server session ID (%x) != client session ID (%x)", s, c)
	***REMOVED*** else if len(s) == 0 ***REMOVED***
		t.Errorf("client and server SessionID were empty.")
	***REMOVED***
***REMOVED***

type noReadConn struct ***REMOVED***
	readSeen bool
	net.Conn
***REMOVED***

func (c *noReadConn) Close() error ***REMOVED***
	return nil
***REMOVED***

func (c *noReadConn) Read(b []byte) (int, error) ***REMOVED***
	c.readSeen = true
	return 0, errors.New("noReadConn error")
***REMOVED***

func TestInvalidServerConfiguration(t *testing.T) ***REMOVED***
	c1, c2, err := netPipe()
	if err != nil ***REMOVED***
		t.Fatalf("netPipe: %v", err)
	***REMOVED***
	defer c1.Close()
	defer c2.Close()

	serveConn := noReadConn***REMOVED***Conn: c1***REMOVED***
	serverConf := &ServerConfig***REMOVED******REMOVED***

	NewServerConn(&serveConn, serverConf)
	if serveConn.readSeen ***REMOVED***
		t.Fatalf("NewServerConn attempted to Read() from Conn while configuration is missing host key")
	***REMOVED***

	serverConf.AddHostKey(testSigners["ecdsa"])

	NewServerConn(&serveConn, serverConf)
	if serveConn.readSeen ***REMOVED***
		t.Fatalf("NewServerConn attempted to Read() from Conn while configuration is missing authentication method")
	***REMOVED***
***REMOVED***

func TestHostKeyAlgorithms(t *testing.T) ***REMOVED***
	serverConf := &ServerConfig***REMOVED***
		NoClientAuth: true,
	***REMOVED***
	serverConf.AddHostKey(testSigners["rsa"])
	serverConf.AddHostKey(testSigners["ecdsa"])

	connect := func(clientConf *ClientConfig, want string) ***REMOVED***
		var alg string
		clientConf.HostKeyCallback = func(h string, a net.Addr, key PublicKey) error ***REMOVED***
			alg = key.Type()
			return nil
		***REMOVED***
		c1, c2, err := netPipe()
		if err != nil ***REMOVED***
			t.Fatalf("netPipe: %v", err)
		***REMOVED***
		defer c1.Close()
		defer c2.Close()

		go NewServerConn(c1, serverConf)
		_, _, _, err = NewClientConn(c2, "", clientConf)
		if err != nil ***REMOVED***
			t.Fatalf("NewClientConn: %v", err)
		***REMOVED***
		if alg != want ***REMOVED***
			t.Errorf("selected key algorithm %s, want %s", alg, want)
		***REMOVED***
	***REMOVED***

	// By default, we get the preferred algorithm, which is ECDSA 256.

	clientConf := &ClientConfig***REMOVED***
		HostKeyCallback: InsecureIgnoreHostKey(),
	***REMOVED***
	connect(clientConf, KeyAlgoECDSA256)

	// Client asks for RSA explicitly.
	clientConf.HostKeyAlgorithms = []string***REMOVED***KeyAlgoRSA***REMOVED***
	connect(clientConf, KeyAlgoRSA)

	c1, c2, err := netPipe()
	if err != nil ***REMOVED***
		t.Fatalf("netPipe: %v", err)
	***REMOVED***
	defer c1.Close()
	defer c2.Close()

	go NewServerConn(c1, serverConf)
	clientConf.HostKeyAlgorithms = []string***REMOVED***"nonexistent-hostkey-algo"***REMOVED***
	_, _, _, err = NewClientConn(c2, "", clientConf)
	if err == nil ***REMOVED***
		t.Fatal("succeeded connecting with unknown hostkey algorithm")
	***REMOVED***
***REMOVED***
