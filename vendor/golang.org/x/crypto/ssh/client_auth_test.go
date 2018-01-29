// Copyright 2011 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package ssh

import (
	"bytes"
	"crypto/rand"
	"errors"
	"fmt"
	"os"
	"strings"
	"testing"
)

type keyboardInteractive map[string]string

func (cr keyboardInteractive) Challenge(user string, instruction string, questions []string, echos []bool) ([]string, error) ***REMOVED***
	var answers []string
	for _, q := range questions ***REMOVED***
		answers = append(answers, cr[q])
	***REMOVED***
	return answers, nil
***REMOVED***

// reused internally by tests
var clientPassword = "tiger"

// tryAuth runs a handshake with a given config against an SSH server
// with config serverConfig
func tryAuth(t *testing.T, config *ClientConfig) error ***REMOVED***
	c1, c2, err := netPipe()
	if err != nil ***REMOVED***
		t.Fatalf("netPipe: %v", err)
	***REMOVED***
	defer c1.Close()
	defer c2.Close()

	certChecker := CertChecker***REMOVED***
		IsUserAuthority: func(k PublicKey) bool ***REMOVED***
			return bytes.Equal(k.Marshal(), testPublicKeys["ecdsa"].Marshal())
		***REMOVED***,
		UserKeyFallback: func(conn ConnMetadata, key PublicKey) (*Permissions, error) ***REMOVED***
			if conn.User() == "testuser" && bytes.Equal(key.Marshal(), testPublicKeys["rsa"].Marshal()) ***REMOVED***
				return nil, nil
			***REMOVED***

			return nil, fmt.Errorf("pubkey for %q not acceptable", conn.User())
		***REMOVED***,
		IsRevoked: func(c *Certificate) bool ***REMOVED***
			return c.Serial == 666
		***REMOVED***,
	***REMOVED***

	serverConfig := &ServerConfig***REMOVED***
		PasswordCallback: func(conn ConnMetadata, pass []byte) (*Permissions, error) ***REMOVED***
			if conn.User() == "testuser" && string(pass) == clientPassword ***REMOVED***
				return nil, nil
			***REMOVED***
			return nil, errors.New("password auth failed")
		***REMOVED***,
		PublicKeyCallback: certChecker.Authenticate,
		KeyboardInteractiveCallback: func(conn ConnMetadata, challenge KeyboardInteractiveChallenge) (*Permissions, error) ***REMOVED***
			ans, err := challenge("user",
				"instruction",
				[]string***REMOVED***"question1", "question2"***REMOVED***,
				[]bool***REMOVED***true, true***REMOVED***)
			if err != nil ***REMOVED***
				return nil, err
			***REMOVED***
			ok := conn.User() == "testuser" && ans[0] == "answer1" && ans[1] == "answer2"
			if ok ***REMOVED***
				challenge("user", "motd", nil, nil)
				return nil, nil
			***REMOVED***
			return nil, errors.New("keyboard-interactive failed")
		***REMOVED***,
	***REMOVED***
	serverConfig.AddHostKey(testSigners["rsa"])

	go newServer(c1, serverConfig)
	_, _, _, err = NewClientConn(c2, "", config)
	return err
***REMOVED***

func TestClientAuthPublicKey(t *testing.T) ***REMOVED***
	config := &ClientConfig***REMOVED***
		User: "testuser",
		Auth: []AuthMethod***REMOVED***
			PublicKeys(testSigners["rsa"]),
		***REMOVED***,
		HostKeyCallback: InsecureIgnoreHostKey(),
	***REMOVED***
	if err := tryAuth(t, config); err != nil ***REMOVED***
		t.Fatalf("unable to dial remote side: %s", err)
	***REMOVED***
***REMOVED***

func TestAuthMethodPassword(t *testing.T) ***REMOVED***
	config := &ClientConfig***REMOVED***
		User: "testuser",
		Auth: []AuthMethod***REMOVED***
			Password(clientPassword),
		***REMOVED***,
		HostKeyCallback: InsecureIgnoreHostKey(),
	***REMOVED***

	if err := tryAuth(t, config); err != nil ***REMOVED***
		t.Fatalf("unable to dial remote side: %s", err)
	***REMOVED***
***REMOVED***

func TestAuthMethodFallback(t *testing.T) ***REMOVED***
	var passwordCalled bool
	config := &ClientConfig***REMOVED***
		User: "testuser",
		Auth: []AuthMethod***REMOVED***
			PublicKeys(testSigners["rsa"]),
			PasswordCallback(
				func() (string, error) ***REMOVED***
					passwordCalled = true
					return "WRONG", nil
				***REMOVED***),
		***REMOVED***,
		HostKeyCallback: InsecureIgnoreHostKey(),
	***REMOVED***

	if err := tryAuth(t, config); err != nil ***REMOVED***
		t.Fatalf("unable to dial remote side: %s", err)
	***REMOVED***

	if passwordCalled ***REMOVED***
		t.Errorf("password auth tried before public-key auth.")
	***REMOVED***
***REMOVED***

func TestAuthMethodWrongPassword(t *testing.T) ***REMOVED***
	config := &ClientConfig***REMOVED***
		User: "testuser",
		Auth: []AuthMethod***REMOVED***
			Password("wrong"),
			PublicKeys(testSigners["rsa"]),
		***REMOVED***,
		HostKeyCallback: InsecureIgnoreHostKey(),
	***REMOVED***

	if err := tryAuth(t, config); err != nil ***REMOVED***
		t.Fatalf("unable to dial remote side: %s", err)
	***REMOVED***
***REMOVED***

func TestAuthMethodKeyboardInteractive(t *testing.T) ***REMOVED***
	answers := keyboardInteractive(map[string]string***REMOVED***
		"question1": "answer1",
		"question2": "answer2",
	***REMOVED***)
	config := &ClientConfig***REMOVED***
		User: "testuser",
		Auth: []AuthMethod***REMOVED***
			KeyboardInteractive(answers.Challenge),
		***REMOVED***,
		HostKeyCallback: InsecureIgnoreHostKey(),
	***REMOVED***

	if err := tryAuth(t, config); err != nil ***REMOVED***
		t.Fatalf("unable to dial remote side: %s", err)
	***REMOVED***
***REMOVED***

func TestAuthMethodWrongKeyboardInteractive(t *testing.T) ***REMOVED***
	answers := keyboardInteractive(map[string]string***REMOVED***
		"question1": "answer1",
		"question2": "WRONG",
	***REMOVED***)
	config := &ClientConfig***REMOVED***
		User: "testuser",
		Auth: []AuthMethod***REMOVED***
			KeyboardInteractive(answers.Challenge),
		***REMOVED***,
	***REMOVED***

	if err := tryAuth(t, config); err == nil ***REMOVED***
		t.Fatalf("wrong answers should not have authenticated with KeyboardInteractive")
	***REMOVED***
***REMOVED***

// the mock server will only authenticate ssh-rsa keys
func TestAuthMethodInvalidPublicKey(t *testing.T) ***REMOVED***
	config := &ClientConfig***REMOVED***
		User: "testuser",
		Auth: []AuthMethod***REMOVED***
			PublicKeys(testSigners["dsa"]),
		***REMOVED***,
	***REMOVED***

	if err := tryAuth(t, config); err == nil ***REMOVED***
		t.Fatalf("dsa private key should not have authenticated with rsa public key")
	***REMOVED***
***REMOVED***

// the client should authenticate with the second key
func TestAuthMethodRSAandDSA(t *testing.T) ***REMOVED***
	config := &ClientConfig***REMOVED***
		User: "testuser",
		Auth: []AuthMethod***REMOVED***
			PublicKeys(testSigners["dsa"], testSigners["rsa"]),
		***REMOVED***,
		HostKeyCallback: InsecureIgnoreHostKey(),
	***REMOVED***
	if err := tryAuth(t, config); err != nil ***REMOVED***
		t.Fatalf("client could not authenticate with rsa key: %v", err)
	***REMOVED***
***REMOVED***

func TestClientHMAC(t *testing.T) ***REMOVED***
	for _, mac := range supportedMACs ***REMOVED***
		config := &ClientConfig***REMOVED***
			User: "testuser",
			Auth: []AuthMethod***REMOVED***
				PublicKeys(testSigners["rsa"]),
			***REMOVED***,
			Config: Config***REMOVED***
				MACs: []string***REMOVED***mac***REMOVED***,
			***REMOVED***,
			HostKeyCallback: InsecureIgnoreHostKey(),
		***REMOVED***
		if err := tryAuth(t, config); err != nil ***REMOVED***
			t.Fatalf("client could not authenticate with mac algo %s: %v", mac, err)
		***REMOVED***
	***REMOVED***
***REMOVED***

// issue 4285.
func TestClientUnsupportedCipher(t *testing.T) ***REMOVED***
	config := &ClientConfig***REMOVED***
		User: "testuser",
		Auth: []AuthMethod***REMOVED***
			PublicKeys(),
		***REMOVED***,
		Config: Config***REMOVED***
			Ciphers: []string***REMOVED***"aes128-cbc"***REMOVED***, // not currently supported
		***REMOVED***,
	***REMOVED***
	if err := tryAuth(t, config); err == nil ***REMOVED***
		t.Errorf("expected no ciphers in common")
	***REMOVED***
***REMOVED***

func TestClientUnsupportedKex(t *testing.T) ***REMOVED***
	if os.Getenv("GO_BUILDER_NAME") != "" ***REMOVED***
		t.Skip("skipping known-flaky test on the Go build dashboard; see golang.org/issue/15198")
	***REMOVED***
	config := &ClientConfig***REMOVED***
		User: "testuser",
		Auth: []AuthMethod***REMOVED***
			PublicKeys(),
		***REMOVED***,
		Config: Config***REMOVED***
			KeyExchanges: []string***REMOVED***"diffie-hellman-group-exchange-sha256"***REMOVED***, // not currently supported
		***REMOVED***,
		HostKeyCallback: InsecureIgnoreHostKey(),
	***REMOVED***
	if err := tryAuth(t, config); err == nil || !strings.Contains(err.Error(), "common algorithm") ***REMOVED***
		t.Errorf("got %v, expected 'common algorithm'", err)
	***REMOVED***
***REMOVED***

func TestClientLoginCert(t *testing.T) ***REMOVED***
	cert := &Certificate***REMOVED***
		Key:         testPublicKeys["rsa"],
		ValidBefore: CertTimeInfinity,
		CertType:    UserCert,
	***REMOVED***
	cert.SignCert(rand.Reader, testSigners["ecdsa"])
	certSigner, err := NewCertSigner(cert, testSigners["rsa"])
	if err != nil ***REMOVED***
		t.Fatalf("NewCertSigner: %v", err)
	***REMOVED***

	clientConfig := &ClientConfig***REMOVED***
		User:            "user",
		HostKeyCallback: InsecureIgnoreHostKey(),
	***REMOVED***
	clientConfig.Auth = append(clientConfig.Auth, PublicKeys(certSigner))

	// should succeed
	if err := tryAuth(t, clientConfig); err != nil ***REMOVED***
		t.Errorf("cert login failed: %v", err)
	***REMOVED***

	// corrupted signature
	cert.Signature.Blob[0]++
	if err := tryAuth(t, clientConfig); err == nil ***REMOVED***
		t.Errorf("cert login passed with corrupted sig")
	***REMOVED***

	// revoked
	cert.Serial = 666
	cert.SignCert(rand.Reader, testSigners["ecdsa"])
	if err := tryAuth(t, clientConfig); err == nil ***REMOVED***
		t.Errorf("revoked cert login succeeded")
	***REMOVED***
	cert.Serial = 1

	// sign with wrong key
	cert.SignCert(rand.Reader, testSigners["dsa"])
	if err := tryAuth(t, clientConfig); err == nil ***REMOVED***
		t.Errorf("cert login passed with non-authoritative key")
	***REMOVED***

	// host cert
	cert.CertType = HostCert
	cert.SignCert(rand.Reader, testSigners["ecdsa"])
	if err := tryAuth(t, clientConfig); err == nil ***REMOVED***
		t.Errorf("cert login passed with wrong type")
	***REMOVED***
	cert.CertType = UserCert

	// principal specified
	cert.ValidPrincipals = []string***REMOVED***"user"***REMOVED***
	cert.SignCert(rand.Reader, testSigners["ecdsa"])
	if err := tryAuth(t, clientConfig); err != nil ***REMOVED***
		t.Errorf("cert login failed: %v", err)
	***REMOVED***

	// wrong principal specified
	cert.ValidPrincipals = []string***REMOVED***"fred"***REMOVED***
	cert.SignCert(rand.Reader, testSigners["ecdsa"])
	if err := tryAuth(t, clientConfig); err == nil ***REMOVED***
		t.Errorf("cert login passed with wrong principal")
	***REMOVED***
	cert.ValidPrincipals = nil

	// added critical option
	cert.CriticalOptions = map[string]string***REMOVED***"root-access": "yes"***REMOVED***
	cert.SignCert(rand.Reader, testSigners["ecdsa"])
	if err := tryAuth(t, clientConfig); err == nil ***REMOVED***
		t.Errorf("cert login passed with unrecognized critical option")
	***REMOVED***

	// allowed source address
	cert.CriticalOptions = map[string]string***REMOVED***"source-address": "127.0.0.42/24,::42/120"***REMOVED***
	cert.SignCert(rand.Reader, testSigners["ecdsa"])
	if err := tryAuth(t, clientConfig); err != nil ***REMOVED***
		t.Errorf("cert login with source-address failed: %v", err)
	***REMOVED***

	// disallowed source address
	cert.CriticalOptions = map[string]string***REMOVED***"source-address": "127.0.0.42,::42"***REMOVED***
	cert.SignCert(rand.Reader, testSigners["ecdsa"])
	if err := tryAuth(t, clientConfig); err == nil ***REMOVED***
		t.Errorf("cert login with source-address succeeded")
	***REMOVED***
***REMOVED***

func testPermissionsPassing(withPermissions bool, t *testing.T) ***REMOVED***
	serverConfig := &ServerConfig***REMOVED***
		PublicKeyCallback: func(conn ConnMetadata, key PublicKey) (*Permissions, error) ***REMOVED***
			if conn.User() == "nopermissions" ***REMOVED***
				return nil, nil
			***REMOVED***
			return &Permissions***REMOVED******REMOVED***, nil
		***REMOVED***,
	***REMOVED***
	serverConfig.AddHostKey(testSigners["rsa"])

	clientConfig := &ClientConfig***REMOVED***
		Auth: []AuthMethod***REMOVED***
			PublicKeys(testSigners["rsa"]),
		***REMOVED***,
		HostKeyCallback: InsecureIgnoreHostKey(),
	***REMOVED***
	if withPermissions ***REMOVED***
		clientConfig.User = "permissions"
	***REMOVED*** else ***REMOVED***
		clientConfig.User = "nopermissions"
	***REMOVED***

	c1, c2, err := netPipe()
	if err != nil ***REMOVED***
		t.Fatalf("netPipe: %v", err)
	***REMOVED***
	defer c1.Close()
	defer c2.Close()

	go NewClientConn(c2, "", clientConfig)
	serverConn, err := newServer(c1, serverConfig)
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
	if p := serverConn.Permissions; (p != nil) != withPermissions ***REMOVED***
		t.Fatalf("withPermissions is %t, but Permissions object is %#v", withPermissions, p)
	***REMOVED***
***REMOVED***

func TestPermissionsPassing(t *testing.T) ***REMOVED***
	testPermissionsPassing(true, t)
***REMOVED***

func TestNoPermissionsPassing(t *testing.T) ***REMOVED***
	testPermissionsPassing(false, t)
***REMOVED***

func TestRetryableAuth(t *testing.T) ***REMOVED***
	n := 0
	passwords := []string***REMOVED***"WRONG1", "WRONG2"***REMOVED***

	config := &ClientConfig***REMOVED***
		User: "testuser",
		Auth: []AuthMethod***REMOVED***
			RetryableAuthMethod(PasswordCallback(func() (string, error) ***REMOVED***
				p := passwords[n]
				n++
				return p, nil
			***REMOVED***), 2),
			PublicKeys(testSigners["rsa"]),
		***REMOVED***,
		HostKeyCallback: InsecureIgnoreHostKey(),
	***REMOVED***

	if err := tryAuth(t, config); err != nil ***REMOVED***
		t.Fatalf("unable to dial remote side: %s", err)
	***REMOVED***
	if n != 2 ***REMOVED***
		t.Fatalf("Did not try all passwords")
	***REMOVED***
***REMOVED***

func ExampleRetryableAuthMethod(t *testing.T) ***REMOVED***
	user := "testuser"
	NumberOfPrompts := 3

	// Normally this would be a callback that prompts the user to answer the
	// provided questions
	Cb := func(user, instruction string, questions []string, echos []bool) (answers []string, err error) ***REMOVED***
		return []string***REMOVED***"answer1", "answer2"***REMOVED***, nil
	***REMOVED***

	config := &ClientConfig***REMOVED***
		HostKeyCallback: InsecureIgnoreHostKey(),
		User:            user,
		Auth: []AuthMethod***REMOVED***
			RetryableAuthMethod(KeyboardInteractiveChallenge(Cb), NumberOfPrompts),
		***REMOVED***,
	***REMOVED***

	if err := tryAuth(t, config); err != nil ***REMOVED***
		t.Fatalf("unable to dial remote side: %s", err)
	***REMOVED***
***REMOVED***

// Test if username is received on server side when NoClientAuth is used
func TestClientAuthNone(t *testing.T) ***REMOVED***
	user := "testuser"
	serverConfig := &ServerConfig***REMOVED***
		NoClientAuth: true,
	***REMOVED***
	serverConfig.AddHostKey(testSigners["rsa"])

	clientConfig := &ClientConfig***REMOVED***
		User:            user,
		HostKeyCallback: InsecureIgnoreHostKey(),
	***REMOVED***

	c1, c2, err := netPipe()
	if err != nil ***REMOVED***
		t.Fatalf("netPipe: %v", err)
	***REMOVED***
	defer c1.Close()
	defer c2.Close()

	go NewClientConn(c2, "", clientConfig)
	serverConn, err := newServer(c1, serverConfig)
	if err != nil ***REMOVED***
		t.Fatalf("newServer: %v", err)
	***REMOVED***
	if serverConn.User() != user ***REMOVED***
		t.Fatalf("server: got %q, want %q", serverConn.User(), user)
	***REMOVED***
***REMOVED***

// Test if authentication attempts are limited on server when MaxAuthTries is set
func TestClientAuthMaxAuthTries(t *testing.T) ***REMOVED***
	user := "testuser"

	serverConfig := &ServerConfig***REMOVED***
		MaxAuthTries: 2,
		PasswordCallback: func(conn ConnMetadata, pass []byte) (*Permissions, error) ***REMOVED***
			if conn.User() == "testuser" && string(pass) == "right" ***REMOVED***
				return nil, nil
			***REMOVED***
			return nil, errors.New("password auth failed")
		***REMOVED***,
	***REMOVED***
	serverConfig.AddHostKey(testSigners["rsa"])

	expectedErr := fmt.Errorf("ssh: handshake failed: %v", &disconnectMsg***REMOVED***
		Reason:  2,
		Message: "too many authentication failures",
	***REMOVED***)

	for tries := 2; tries < 4; tries++ ***REMOVED***
		n := tries
		clientConfig := &ClientConfig***REMOVED***
			User: user,
			Auth: []AuthMethod***REMOVED***
				RetryableAuthMethod(PasswordCallback(func() (string, error) ***REMOVED***
					n--
					if n == 0 ***REMOVED***
						return "right", nil
					***REMOVED***
					return "wrong", nil
				***REMOVED***), tries),
			***REMOVED***,
			HostKeyCallback: InsecureIgnoreHostKey(),
		***REMOVED***

		c1, c2, err := netPipe()
		if err != nil ***REMOVED***
			t.Fatalf("netPipe: %v", err)
		***REMOVED***
		defer c1.Close()
		defer c2.Close()

		go newServer(c1, serverConfig)
		_, _, _, err = NewClientConn(c2, "", clientConfig)
		if tries > 2 ***REMOVED***
			if err == nil ***REMOVED***
				t.Fatalf("client: got no error, want %s", expectedErr)
			***REMOVED*** else if err.Error() != expectedErr.Error() ***REMOVED***
				t.Fatalf("client: got %s, want %s", err, expectedErr)
			***REMOVED***
		***REMOVED*** else ***REMOVED***
			if err != nil ***REMOVED***
				t.Fatalf("client: got %s, want no error", err)
			***REMOVED***
		***REMOVED***
	***REMOVED***
***REMOVED***

// Test if authentication attempts are correctly limited on server
// when more public keys are provided then MaxAuthTries
func TestClientAuthMaxAuthTriesPublicKey(t *testing.T) ***REMOVED***
	signers := []Signer***REMOVED******REMOVED***
	for i := 0; i < 6; i++ ***REMOVED***
		signers = append(signers, testSigners["dsa"])
	***REMOVED***

	validConfig := &ClientConfig***REMOVED***
		User: "testuser",
		Auth: []AuthMethod***REMOVED***
			PublicKeys(append([]Signer***REMOVED***testSigners["rsa"]***REMOVED***, signers...)...),
		***REMOVED***,
		HostKeyCallback: InsecureIgnoreHostKey(),
	***REMOVED***
	if err := tryAuth(t, validConfig); err != nil ***REMOVED***
		t.Fatalf("unable to dial remote side: %s", err)
	***REMOVED***

	expectedErr := fmt.Errorf("ssh: handshake failed: %v", &disconnectMsg***REMOVED***
		Reason:  2,
		Message: "too many authentication failures",
	***REMOVED***)
	invalidConfig := &ClientConfig***REMOVED***
		User: "testuser",
		Auth: []AuthMethod***REMOVED***
			PublicKeys(append(signers, testSigners["rsa"])...),
		***REMOVED***,
		HostKeyCallback: InsecureIgnoreHostKey(),
	***REMOVED***
	if err := tryAuth(t, invalidConfig); err == nil ***REMOVED***
		t.Fatalf("client: got no error, want %s", expectedErr)
	***REMOVED*** else if err.Error() != expectedErr.Error() ***REMOVED***
		t.Fatalf("client: got %s, want %s", err, expectedErr)
	***REMOVED***
***REMOVED***

// Test whether authentication errors are being properly logged if all
// authentication methods have been exhausted
func TestClientAuthErrorList(t *testing.T) ***REMOVED***
	publicKeyErr := errors.New("This is an error from PublicKeyCallback")

	clientConfig := &ClientConfig***REMOVED***
		Auth: []AuthMethod***REMOVED***
			PublicKeys(testSigners["rsa"]),
		***REMOVED***,
		HostKeyCallback: InsecureIgnoreHostKey(),
	***REMOVED***
	serverConfig := &ServerConfig***REMOVED***
		PublicKeyCallback: func(_ ConnMetadata, _ PublicKey) (*Permissions, error) ***REMOVED***
			return nil, publicKeyErr
		***REMOVED***,
	***REMOVED***
	serverConfig.AddHostKey(testSigners["rsa"])

	c1, c2, err := netPipe()
	if err != nil ***REMOVED***
		t.Fatalf("netPipe: %v", err)
	***REMOVED***
	defer c1.Close()
	defer c2.Close()

	go NewClientConn(c2, "", clientConfig)
	_, err = newServer(c1, serverConfig)
	if err == nil ***REMOVED***
		t.Fatal("newServer: got nil, expected errors")
	***REMOVED***

	authErrs, ok := err.(*ServerAuthError)
	if !ok ***REMOVED***
		t.Fatalf("errors: got %T, want *ssh.ServerAuthError", err)
	***REMOVED***
	for i, e := range authErrs.Errors ***REMOVED***
		switch i ***REMOVED***
		case 0:
			if e.Error() != "no auth passed yet" ***REMOVED***
				t.Fatalf("errors: got %v, want no auth passed yet", e.Error())
			***REMOVED***
		case 1:
			if e != publicKeyErr ***REMOVED***
				t.Fatalf("errors: got %v, want %v", e, publicKeyErr)
			***REMOVED***
		default:
			t.Fatalf("errors: got %v, expected 2 errors", authErrs.Errors)
		***REMOVED***
	***REMOVED***
***REMOVED***
