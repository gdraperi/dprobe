// Copyright 2012 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build darwin dragonfly freebsd linux netbsd openbsd plan9

package test

// functional test harness for unix.

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"os"
	"os/exec"
	"os/user"
	"path/filepath"
	"testing"
	"text/template"

	"golang.org/x/crypto/ssh"
	"golang.org/x/crypto/ssh/testdata"
)

const sshdConfig = `
Protocol 2
Banner ***REMOVED******REMOVED***.Dir***REMOVED******REMOVED***/banner
HostKey ***REMOVED******REMOVED***.Dir***REMOVED******REMOVED***/id_rsa
HostKey ***REMOVED******REMOVED***.Dir***REMOVED******REMOVED***/id_dsa
HostKey ***REMOVED******REMOVED***.Dir***REMOVED******REMOVED***/id_ecdsa
HostCertificate ***REMOVED******REMOVED***.Dir***REMOVED******REMOVED***/id_rsa-cert.pub
Pidfile ***REMOVED******REMOVED***.Dir***REMOVED******REMOVED***/sshd.pid
#UsePrivilegeSeparation no
KeyRegenerationInterval 3600
ServerKeyBits 768
SyslogFacility AUTH
LogLevel DEBUG2
LoginGraceTime 120
PermitRootLogin no
StrictModes no
RSAAuthentication yes
PubkeyAuthentication yes
AuthorizedKeysFile	***REMOVED******REMOVED***.Dir***REMOVED******REMOVED***/authorized_keys
TrustedUserCAKeys ***REMOVED******REMOVED***.Dir***REMOVED******REMOVED***/id_ecdsa.pub
IgnoreRhosts yes
RhostsRSAAuthentication no
HostbasedAuthentication no
PubkeyAcceptedKeyTypes=*
`

var configTmpl = template.Must(template.New("").Parse(sshdConfig))

type server struct ***REMOVED***
	t          *testing.T
	cleanup    func() // executed during Shutdown
	configfile string
	cmd        *exec.Cmd
	output     bytes.Buffer // holds stderr from sshd process

	// Client half of the network connection.
	clientConn net.Conn
***REMOVED***

func username() string ***REMOVED***
	var username string
	if user, err := user.Current(); err == nil ***REMOVED***
		username = user.Username
	***REMOVED*** else ***REMOVED***
		// user.Current() currently requires cgo. If an error is
		// returned attempt to get the username from the environment.
		log.Printf("user.Current: %v; falling back on $USER", err)
		username = os.Getenv("USER")
	***REMOVED***
	if username == "" ***REMOVED***
		panic("Unable to get username")
	***REMOVED***
	return username
***REMOVED***

type storedHostKey struct ***REMOVED***
	// keys map from an algorithm string to binary key data.
	keys map[string][]byte

	// checkCount counts the Check calls. Used for testing
	// rekeying.
	checkCount int
***REMOVED***

func (k *storedHostKey) Add(key ssh.PublicKey) ***REMOVED***
	if k.keys == nil ***REMOVED***
		k.keys = map[string][]byte***REMOVED******REMOVED***
	***REMOVED***
	k.keys[key.Type()] = key.Marshal()
***REMOVED***

func (k *storedHostKey) Check(addr string, remote net.Addr, key ssh.PublicKey) error ***REMOVED***
	k.checkCount++
	algo := key.Type()

	if k.keys == nil || bytes.Compare(key.Marshal(), k.keys[algo]) != 0 ***REMOVED***
		return fmt.Errorf("host key mismatch. Got %q, want %q", key, k.keys[algo])
	***REMOVED***
	return nil
***REMOVED***

func hostKeyDB() *storedHostKey ***REMOVED***
	keyChecker := &storedHostKey***REMOVED******REMOVED***
	keyChecker.Add(testPublicKeys["ecdsa"])
	keyChecker.Add(testPublicKeys["rsa"])
	keyChecker.Add(testPublicKeys["dsa"])
	return keyChecker
***REMOVED***

func clientConfig() *ssh.ClientConfig ***REMOVED***
	config := &ssh.ClientConfig***REMOVED***
		User: username(),
		Auth: []ssh.AuthMethod***REMOVED***
			ssh.PublicKeys(testSigners["user"]),
		***REMOVED***,
		HostKeyCallback: hostKeyDB().Check,
		HostKeyAlgorithms: []string***REMOVED*** // by default, don't allow certs as this affects the hostKeyDB checker
			ssh.KeyAlgoECDSA256, ssh.KeyAlgoECDSA384, ssh.KeyAlgoECDSA521,
			ssh.KeyAlgoRSA, ssh.KeyAlgoDSA,
			ssh.KeyAlgoED25519,
		***REMOVED***,
	***REMOVED***
	return config
***REMOVED***

// unixConnection creates two halves of a connected net.UnixConn.  It
// is used for connecting the Go SSH client with sshd without opening
// ports.
func unixConnection() (*net.UnixConn, *net.UnixConn, error) ***REMOVED***
	dir, err := ioutil.TempDir("", "unixConnection")
	if err != nil ***REMOVED***
		return nil, nil, err
	***REMOVED***
	defer os.Remove(dir)

	addr := filepath.Join(dir, "ssh")
	listener, err := net.Listen("unix", addr)
	if err != nil ***REMOVED***
		return nil, nil, err
	***REMOVED***
	defer listener.Close()
	c1, err := net.Dial("unix", addr)
	if err != nil ***REMOVED***
		return nil, nil, err
	***REMOVED***

	c2, err := listener.Accept()
	if err != nil ***REMOVED***
		c1.Close()
		return nil, nil, err
	***REMOVED***

	return c1.(*net.UnixConn), c2.(*net.UnixConn), nil
***REMOVED***

func (s *server) TryDial(config *ssh.ClientConfig) (*ssh.Client, error) ***REMOVED***
	return s.TryDialWithAddr(config, "")
***REMOVED***

// addr is the user specified host:port. While we don't actually dial it,
// we need to know this for host key matching
func (s *server) TryDialWithAddr(config *ssh.ClientConfig, addr string) (*ssh.Client, error) ***REMOVED***
	sshd, err := exec.LookPath("sshd")
	if err != nil ***REMOVED***
		s.t.Skipf("skipping test: %v", err)
	***REMOVED***

	c1, c2, err := unixConnection()
	if err != nil ***REMOVED***
		s.t.Fatalf("unixConnection: %v", err)
	***REMOVED***

	s.cmd = exec.Command(sshd, "-f", s.configfile, "-i", "-e")
	f, err := c2.File()
	if err != nil ***REMOVED***
		s.t.Fatalf("UnixConn.File: %v", err)
	***REMOVED***
	defer f.Close()
	s.cmd.Stdin = f
	s.cmd.Stdout = f
	s.cmd.Stderr = &s.output
	if err := s.cmd.Start(); err != nil ***REMOVED***
		s.t.Fail()
		s.Shutdown()
		s.t.Fatalf("s.cmd.Start: %v", err)
	***REMOVED***
	s.clientConn = c1
	conn, chans, reqs, err := ssh.NewClientConn(c1, addr, config)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	return ssh.NewClient(conn, chans, reqs), nil
***REMOVED***

func (s *server) Dial(config *ssh.ClientConfig) *ssh.Client ***REMOVED***
	conn, err := s.TryDial(config)
	if err != nil ***REMOVED***
		s.t.Fail()
		s.Shutdown()
		s.t.Fatalf("ssh.Client: %v", err)
	***REMOVED***
	return conn
***REMOVED***

func (s *server) Shutdown() ***REMOVED***
	if s.cmd != nil && s.cmd.Process != nil ***REMOVED***
		// Don't check for errors; if it fails it's most
		// likely "os: process already finished", and we don't
		// care about that. Use os.Interrupt, so child
		// processes are killed too.
		s.cmd.Process.Signal(os.Interrupt)
		s.cmd.Wait()
	***REMOVED***
	if s.t.Failed() ***REMOVED***
		// log any output from sshd process
		s.t.Logf("sshd: %s", s.output.String())
	***REMOVED***
	s.cleanup()
***REMOVED***

func writeFile(path string, contents []byte) ***REMOVED***
	f, err := os.OpenFile(path, os.O_WRONLY|os.O_TRUNC|os.O_CREATE, 0600)
	if err != nil ***REMOVED***
		panic(err)
	***REMOVED***
	defer f.Close()
	if _, err := f.Write(contents); err != nil ***REMOVED***
		panic(err)
	***REMOVED***
***REMOVED***

// newServer returns a new mock ssh server.
func newServer(t *testing.T) *server ***REMOVED***
	if testing.Short() ***REMOVED***
		t.Skip("skipping test due to -short")
	***REMOVED***
	dir, err := ioutil.TempDir("", "sshtest")
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
	f, err := os.Create(filepath.Join(dir, "sshd_config"))
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
	err = configTmpl.Execute(f, map[string]string***REMOVED***
		"Dir": dir,
	***REMOVED***)
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
	f.Close()

	writeFile(filepath.Join(dir, "banner"), []byte("Server Banner"))

	for k, v := range testdata.PEMBytes ***REMOVED***
		filename := "id_" + k
		writeFile(filepath.Join(dir, filename), v)
		writeFile(filepath.Join(dir, filename+".pub"), ssh.MarshalAuthorizedKey(testPublicKeys[k]))
	***REMOVED***

	for k, v := range testdata.SSHCertificates ***REMOVED***
		filename := "id_" + k + "-cert.pub"
		writeFile(filepath.Join(dir, filename), v)
	***REMOVED***

	var authkeys bytes.Buffer
	for k := range testdata.PEMBytes ***REMOVED***
		authkeys.Write(ssh.MarshalAuthorizedKey(testPublicKeys[k]))
	***REMOVED***
	writeFile(filepath.Join(dir, "authorized_keys"), authkeys.Bytes())

	return &server***REMOVED***
		t:          t,
		configfile: f.Name(),
		cleanup: func() ***REMOVED***
			if err := os.RemoveAll(dir); err != nil ***REMOVED***
				t.Error(err)
			***REMOVED***
		***REMOVED***,
	***REMOVED***
***REMOVED***

func newTempSocket(t *testing.T) (string, func()) ***REMOVED***
	dir, err := ioutil.TempDir("", "socket")
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
	deferFunc := func() ***REMOVED*** os.RemoveAll(dir) ***REMOVED***
	addr := filepath.Join(dir, "sock")
	return addr, deferFunc
***REMOVED***
