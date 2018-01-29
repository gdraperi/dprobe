// Copyright 2014 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package ssh

import (
	"bytes"
	"crypto/dsa"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"fmt"
	"reflect"
	"strings"
	"testing"

	"golang.org/x/crypto/ed25519"
	"golang.org/x/crypto/ssh/testdata"
)

func rawKey(pub PublicKey) interface***REMOVED******REMOVED*** ***REMOVED***
	switch k := pub.(type) ***REMOVED***
	case *rsaPublicKey:
		return (*rsa.PublicKey)(k)
	case *dsaPublicKey:
		return (*dsa.PublicKey)(k)
	case *ecdsaPublicKey:
		return (*ecdsa.PublicKey)(k)
	case ed25519PublicKey:
		return (ed25519.PublicKey)(k)
	case *Certificate:
		return k
	***REMOVED***
	panic("unknown key type")
***REMOVED***

func TestKeyMarshalParse(t *testing.T) ***REMOVED***
	for _, priv := range testSigners ***REMOVED***
		pub := priv.PublicKey()
		roundtrip, err := ParsePublicKey(pub.Marshal())
		if err != nil ***REMOVED***
			t.Errorf("ParsePublicKey(%T): %v", pub, err)
		***REMOVED***

		k1 := rawKey(pub)
		k2 := rawKey(roundtrip)

		if !reflect.DeepEqual(k1, k2) ***REMOVED***
			t.Errorf("got %#v in roundtrip, want %#v", k2, k1)
		***REMOVED***
	***REMOVED***
***REMOVED***

func TestUnsupportedCurves(t *testing.T) ***REMOVED***
	raw, err := ecdsa.GenerateKey(elliptic.P224(), rand.Reader)
	if err != nil ***REMOVED***
		t.Fatalf("GenerateKey: %v", err)
	***REMOVED***

	if _, err = NewSignerFromKey(raw); err == nil || !strings.Contains(err.Error(), "only P-256") ***REMOVED***
		t.Fatalf("NewPrivateKey should not succeed with P-224, got: %v", err)
	***REMOVED***

	if _, err = NewPublicKey(&raw.PublicKey); err == nil || !strings.Contains(err.Error(), "only P-256") ***REMOVED***
		t.Fatalf("NewPublicKey should not succeed with P-224, got: %v", err)
	***REMOVED***
***REMOVED***

func TestNewPublicKey(t *testing.T) ***REMOVED***
	for _, k := range testSigners ***REMOVED***
		raw := rawKey(k.PublicKey())
		// Skip certificates, as NewPublicKey does not support them.
		if _, ok := raw.(*Certificate); ok ***REMOVED***
			continue
		***REMOVED***
		pub, err := NewPublicKey(raw)
		if err != nil ***REMOVED***
			t.Errorf("NewPublicKey(%#v): %v", raw, err)
		***REMOVED***
		if !reflect.DeepEqual(k.PublicKey(), pub) ***REMOVED***
			t.Errorf("NewPublicKey(%#v) = %#v, want %#v", raw, pub, k.PublicKey())
		***REMOVED***
	***REMOVED***
***REMOVED***

func TestKeySignVerify(t *testing.T) ***REMOVED***
	for _, priv := range testSigners ***REMOVED***
		pub := priv.PublicKey()

		data := []byte("sign me")
		sig, err := priv.Sign(rand.Reader, data)
		if err != nil ***REMOVED***
			t.Fatalf("Sign(%T): %v", priv, err)
		***REMOVED***

		if err := pub.Verify(data, sig); err != nil ***REMOVED***
			t.Errorf("publicKey.Verify(%T): %v", priv, err)
		***REMOVED***
		sig.Blob[5]++
		if err := pub.Verify(data, sig); err == nil ***REMOVED***
			t.Errorf("publicKey.Verify on broken sig did not fail")
		***REMOVED***
	***REMOVED***
***REMOVED***

func TestParseRSAPrivateKey(t *testing.T) ***REMOVED***
	key := testPrivateKeys["rsa"]

	rsa, ok := key.(*rsa.PrivateKey)
	if !ok ***REMOVED***
		t.Fatalf("got %T, want *rsa.PrivateKey", rsa)
	***REMOVED***

	if err := rsa.Validate(); err != nil ***REMOVED***
		t.Errorf("Validate: %v", err)
	***REMOVED***
***REMOVED***

func TestParseECPrivateKey(t *testing.T) ***REMOVED***
	key := testPrivateKeys["ecdsa"]

	ecKey, ok := key.(*ecdsa.PrivateKey)
	if !ok ***REMOVED***
		t.Fatalf("got %T, want *ecdsa.PrivateKey", ecKey)
	***REMOVED***

	if !validateECPublicKey(ecKey.Curve, ecKey.X, ecKey.Y) ***REMOVED***
		t.Fatalf("public key does not validate.")
	***REMOVED***
***REMOVED***

// See Issue https://github.com/golang/go/issues/6650.
func TestParseEncryptedPrivateKeysFails(t *testing.T) ***REMOVED***
	const wantSubstring = "encrypted"
	for i, tt := range testdata.PEMEncryptedKeys ***REMOVED***
		_, err := ParsePrivateKey(tt.PEMBytes)
		if err == nil ***REMOVED***
			t.Errorf("#%d key %s: ParsePrivateKey successfully parsed, expected an error", i, tt.Name)
			continue
		***REMOVED***

		if !strings.Contains(err.Error(), wantSubstring) ***REMOVED***
			t.Errorf("#%d key %s: got error %q, want substring %q", i, tt.Name, err, wantSubstring)
		***REMOVED***
	***REMOVED***
***REMOVED***

// Parse encrypted private keys with passphrase
func TestParseEncryptedPrivateKeysWithPassphrase(t *testing.T) ***REMOVED***
	data := []byte("sign me")
	for _, tt := range testdata.PEMEncryptedKeys ***REMOVED***
		s, err := ParsePrivateKeyWithPassphrase(tt.PEMBytes, []byte(tt.EncryptionKey))
		if err != nil ***REMOVED***
			t.Fatalf("ParsePrivateKeyWithPassphrase returned error: %s", err)
			continue
		***REMOVED***
		sig, err := s.Sign(rand.Reader, data)
		if err != nil ***REMOVED***
			t.Fatalf("dsa.Sign: %v", err)
		***REMOVED***
		if err := s.PublicKey().Verify(data, sig); err != nil ***REMOVED***
			t.Errorf("Verify failed: %v", err)
		***REMOVED***
	***REMOVED***

	tt := testdata.PEMEncryptedKeys[0]
	_, err := ParsePrivateKeyWithPassphrase(tt.PEMBytes, []byte("incorrect"))
	if err != x509.IncorrectPasswordError ***REMOVED***
		t.Fatalf("got %v want IncorrectPasswordError", err)
	***REMOVED***
***REMOVED***

func TestParseDSA(t *testing.T) ***REMOVED***
	// We actually exercise the ParsePrivateKey codepath here, as opposed to
	// using the ParseRawPrivateKey+NewSignerFromKey path that testdata_test.go
	// uses.
	s, err := ParsePrivateKey(testdata.PEMBytes["dsa"])
	if err != nil ***REMOVED***
		t.Fatalf("ParsePrivateKey returned error: %s", err)
	***REMOVED***

	data := []byte("sign me")
	sig, err := s.Sign(rand.Reader, data)
	if err != nil ***REMOVED***
		t.Fatalf("dsa.Sign: %v", err)
	***REMOVED***

	if err := s.PublicKey().Verify(data, sig); err != nil ***REMOVED***
		t.Errorf("Verify failed: %v", err)
	***REMOVED***
***REMOVED***

// Tests for authorized_keys parsing.

// getTestKey returns a public key, and its base64 encoding.
func getTestKey() (PublicKey, string) ***REMOVED***
	k := testPublicKeys["rsa"]

	b := &bytes.Buffer***REMOVED******REMOVED***
	e := base64.NewEncoder(base64.StdEncoding, b)
	e.Write(k.Marshal())
	e.Close()

	return k, b.String()
***REMOVED***

func TestMarshalParsePublicKey(t *testing.T) ***REMOVED***
	pub, pubSerialized := getTestKey()
	line := fmt.Sprintf("%s %s user@host", pub.Type(), pubSerialized)

	authKeys := MarshalAuthorizedKey(pub)
	actualFields := strings.Fields(string(authKeys))
	if len(actualFields) == 0 ***REMOVED***
		t.Fatalf("failed authKeys: %v", authKeys)
	***REMOVED***

	// drop the comment
	expectedFields := strings.Fields(line)[0:2]

	if !reflect.DeepEqual(actualFields, expectedFields) ***REMOVED***
		t.Errorf("got %v, expected %v", actualFields, expectedFields)
	***REMOVED***

	actPub, _, _, _, err := ParseAuthorizedKey([]byte(line))
	if err != nil ***REMOVED***
		t.Fatalf("cannot parse %v: %v", line, err)
	***REMOVED***
	if !reflect.DeepEqual(actPub, pub) ***REMOVED***
		t.Errorf("got %v, expected %v", actPub, pub)
	***REMOVED***
***REMOVED***

type authResult struct ***REMOVED***
	pubKey   PublicKey
	options  []string
	comments string
	rest     string
	ok       bool
***REMOVED***

func testAuthorizedKeys(t *testing.T, authKeys []byte, expected []authResult) ***REMOVED***
	rest := authKeys
	var values []authResult
	for len(rest) > 0 ***REMOVED***
		var r authResult
		var err error
		r.pubKey, r.comments, r.options, rest, err = ParseAuthorizedKey(rest)
		r.ok = (err == nil)
		t.Log(err)
		r.rest = string(rest)
		values = append(values, r)
	***REMOVED***

	if !reflect.DeepEqual(values, expected) ***REMOVED***
		t.Errorf("got %#v, expected %#v", values, expected)
	***REMOVED***
***REMOVED***

func TestAuthorizedKeyBasic(t *testing.T) ***REMOVED***
	pub, pubSerialized := getTestKey()
	line := "ssh-rsa " + pubSerialized + " user@host"
	testAuthorizedKeys(t, []byte(line),
		[]authResult***REMOVED***
			***REMOVED***pub, nil, "user@host", "", true***REMOVED***,
		***REMOVED***)
***REMOVED***

func TestAuth(t *testing.T) ***REMOVED***
	pub, pubSerialized := getTestKey()
	authWithOptions := []string***REMOVED***
		`# comments to ignore before any keys...`,
		``,
		`env="HOME=/home/root",no-port-forwarding ssh-rsa ` + pubSerialized + ` user@host`,
		`# comments to ignore, along with a blank line`,
		``,
		`env="HOME=/home/root2" ssh-rsa ` + pubSerialized + ` user2@host2`,
		``,
		`# more comments, plus a invalid entry`,
		`ssh-rsa data-that-will-not-parse user@host3`,
	***REMOVED***
	for _, eol := range []string***REMOVED***"\n", "\r\n"***REMOVED*** ***REMOVED***
		authOptions := strings.Join(authWithOptions, eol)
		rest2 := strings.Join(authWithOptions[3:], eol)
		rest3 := strings.Join(authWithOptions[6:], eol)
		testAuthorizedKeys(t, []byte(authOptions), []authResult***REMOVED***
			***REMOVED***pub, []string***REMOVED***`env="HOME=/home/root"`, "no-port-forwarding"***REMOVED***, "user@host", rest2, true***REMOVED***,
			***REMOVED***pub, []string***REMOVED***`env="HOME=/home/root2"`***REMOVED***, "user2@host2", rest3, true***REMOVED***,
			***REMOVED***nil, nil, "", "", false***REMOVED***,
		***REMOVED***)
	***REMOVED***
***REMOVED***

func TestAuthWithQuotedSpaceInEnv(t *testing.T) ***REMOVED***
	pub, pubSerialized := getTestKey()
	authWithQuotedSpaceInEnv := []byte(`env="HOME=/home/root dir",no-port-forwarding ssh-rsa ` + pubSerialized + ` user@host`)
	testAuthorizedKeys(t, []byte(authWithQuotedSpaceInEnv), []authResult***REMOVED***
		***REMOVED***pub, []string***REMOVED***`env="HOME=/home/root dir"`, "no-port-forwarding"***REMOVED***, "user@host", "", true***REMOVED***,
	***REMOVED***)
***REMOVED***

func TestAuthWithQuotedCommaInEnv(t *testing.T) ***REMOVED***
	pub, pubSerialized := getTestKey()
	authWithQuotedCommaInEnv := []byte(`env="HOME=/home/root,dir",no-port-forwarding ssh-rsa ` + pubSerialized + `   user@host`)
	testAuthorizedKeys(t, []byte(authWithQuotedCommaInEnv), []authResult***REMOVED***
		***REMOVED***pub, []string***REMOVED***`env="HOME=/home/root,dir"`, "no-port-forwarding"***REMOVED***, "user@host", "", true***REMOVED***,
	***REMOVED***)
***REMOVED***

func TestAuthWithQuotedQuoteInEnv(t *testing.T) ***REMOVED***
	pub, pubSerialized := getTestKey()
	authWithQuotedQuoteInEnv := []byte(`env="HOME=/home/\"root dir",no-port-forwarding` + "\t" + `ssh-rsa` + "\t" + pubSerialized + `   user@host`)
	authWithDoubleQuotedQuote := []byte(`no-port-forwarding,env="HOME=/home/ \"root dir\"" ssh-rsa ` + pubSerialized + "\t" + `user@host`)
	testAuthorizedKeys(t, []byte(authWithQuotedQuoteInEnv), []authResult***REMOVED***
		***REMOVED***pub, []string***REMOVED***`env="HOME=/home/\"root dir"`, "no-port-forwarding"***REMOVED***, "user@host", "", true***REMOVED***,
	***REMOVED***)

	testAuthorizedKeys(t, []byte(authWithDoubleQuotedQuote), []authResult***REMOVED***
		***REMOVED***pub, []string***REMOVED***"no-port-forwarding", `env="HOME=/home/ \"root dir\""`***REMOVED***, "user@host", "", true***REMOVED***,
	***REMOVED***)
***REMOVED***

func TestAuthWithInvalidSpace(t *testing.T) ***REMOVED***
	_, pubSerialized := getTestKey()
	authWithInvalidSpace := []byte(`env="HOME=/home/root dir", no-port-forwarding ssh-rsa ` + pubSerialized + ` user@host
#more to follow but still no valid keys`)
	testAuthorizedKeys(t, []byte(authWithInvalidSpace), []authResult***REMOVED***
		***REMOVED***nil, nil, "", "", false***REMOVED***,
	***REMOVED***)
***REMOVED***

func TestAuthWithMissingQuote(t *testing.T) ***REMOVED***
	pub, pubSerialized := getTestKey()
	authWithMissingQuote := []byte(`env="HOME=/home/root,no-port-forwarding ssh-rsa ` + pubSerialized + ` user@host
env="HOME=/home/root",shared-control ssh-rsa ` + pubSerialized + ` user@host`)

	testAuthorizedKeys(t, []byte(authWithMissingQuote), []authResult***REMOVED***
		***REMOVED***pub, []string***REMOVED***`env="HOME=/home/root"`, `shared-control`***REMOVED***, "user@host", "", true***REMOVED***,
	***REMOVED***)
***REMOVED***

func TestInvalidEntry(t *testing.T) ***REMOVED***
	authInvalid := []byte(`ssh-rsa`)
	_, _, _, _, err := ParseAuthorizedKey(authInvalid)
	if err == nil ***REMOVED***
		t.Errorf("got valid entry for %q", authInvalid)
	***REMOVED***
***REMOVED***

var knownHostsParseTests = []struct ***REMOVED***
	input string
	err   string

	marker  string
	comment string
	hosts   []string
	rest    string
***REMOVED******REMOVED***
	***REMOVED***
		"",
		"EOF",

		"", "", nil, "",
	***REMOVED***,
	***REMOVED***
		"# Just a comment",
		"EOF",

		"", "", nil, "",
	***REMOVED***,
	***REMOVED***
		"   \t   ",
		"EOF",

		"", "", nil, "",
	***REMOVED***,
	***REMOVED***
		"localhost ssh-rsa ***REMOVED***RSAPUB***REMOVED***",
		"",

		"", "", []string***REMOVED***"localhost"***REMOVED***, "",
	***REMOVED***,
	***REMOVED***
		"localhost\tssh-rsa ***REMOVED***RSAPUB***REMOVED***",
		"",

		"", "", []string***REMOVED***"localhost"***REMOVED***, "",
	***REMOVED***,
	***REMOVED***
		"localhost\tssh-rsa ***REMOVED***RSAPUB***REMOVED***\tcomment comment",
		"",

		"", "comment comment", []string***REMOVED***"localhost"***REMOVED***, "",
	***REMOVED***,
	***REMOVED***
		"localhost\tssh-rsa ***REMOVED***RSAPUB***REMOVED***\tcomment comment\n",
		"",

		"", "comment comment", []string***REMOVED***"localhost"***REMOVED***, "",
	***REMOVED***,
	***REMOVED***
		"localhost\tssh-rsa ***REMOVED***RSAPUB***REMOVED***\tcomment comment\r\n",
		"",

		"", "comment comment", []string***REMOVED***"localhost"***REMOVED***, "",
	***REMOVED***,
	***REMOVED***
		"localhost\tssh-rsa ***REMOVED***RSAPUB***REMOVED***\tcomment comment\r\nnext line",
		"",

		"", "comment comment", []string***REMOVED***"localhost"***REMOVED***, "next line",
	***REMOVED***,
	***REMOVED***
		"localhost,[host2:123]\tssh-rsa ***REMOVED***RSAPUB***REMOVED***\tcomment comment",
		"",

		"", "comment comment", []string***REMOVED***"localhost", "[host2:123]"***REMOVED***, "",
	***REMOVED***,
	***REMOVED***
		"@marker \tlocalhost,[host2:123]\tssh-rsa ***REMOVED***RSAPUB***REMOVED***",
		"",

		"marker", "", []string***REMOVED***"localhost", "[host2:123]"***REMOVED***, "",
	***REMOVED***,
	***REMOVED***
		"@marker \tlocalhost,[host2:123]\tssh-rsa aabbccdd",
		"short read",

		"", "", nil, "",
	***REMOVED***,
***REMOVED***

func TestKnownHostsParsing(t *testing.T) ***REMOVED***
	rsaPub, rsaPubSerialized := getTestKey()

	for i, test := range knownHostsParseTests ***REMOVED***
		var expectedKey PublicKey
		const rsaKeyToken = "***REMOVED***RSAPUB***REMOVED***"

		input := test.input
		if strings.Contains(input, rsaKeyToken) ***REMOVED***
			expectedKey = rsaPub
			input = strings.Replace(test.input, rsaKeyToken, rsaPubSerialized, -1)
		***REMOVED***

		marker, hosts, pubKey, comment, rest, err := ParseKnownHosts([]byte(input))
		if err != nil ***REMOVED***
			if len(test.err) == 0 ***REMOVED***
				t.Errorf("#%d: unexpectedly failed with %q", i, err)
			***REMOVED*** else if !strings.Contains(err.Error(), test.err) ***REMOVED***
				t.Errorf("#%d: expected error containing %q, but got %q", i, test.err, err)
			***REMOVED***
			continue
		***REMOVED*** else if len(test.err) != 0 ***REMOVED***
			t.Errorf("#%d: succeeded but expected error including %q", i, test.err)
			continue
		***REMOVED***

		if !reflect.DeepEqual(expectedKey, pubKey) ***REMOVED***
			t.Errorf("#%d: expected key %#v, but got %#v", i, expectedKey, pubKey)
		***REMOVED***

		if marker != test.marker ***REMOVED***
			t.Errorf("#%d: expected marker %q, but got %q", i, test.marker, marker)
		***REMOVED***

		if comment != test.comment ***REMOVED***
			t.Errorf("#%d: expected comment %q, but got %q", i, test.comment, comment)
		***REMOVED***

		if !reflect.DeepEqual(test.hosts, hosts) ***REMOVED***
			t.Errorf("#%d: expected hosts %#v, but got %#v", i, test.hosts, hosts)
		***REMOVED***

		if rest := string(rest); rest != test.rest ***REMOVED***
			t.Errorf("#%d: expected remaining input to be %q, but got %q", i, test.rest, rest)
		***REMOVED***
	***REMOVED***
***REMOVED***

func TestFingerprintLegacyMD5(t *testing.T) ***REMOVED***
	pub, _ := getTestKey()
	fingerprint := FingerprintLegacyMD5(pub)
	want := "fb:61:6d:1a:e3:f0:95:45:3c:a0:79:be:4a:93:63:66" // ssh-keygen -lf -E md5 rsa
	if fingerprint != want ***REMOVED***
		t.Errorf("got fingerprint %q want %q", fingerprint, want)
	***REMOVED***
***REMOVED***

func TestFingerprintSHA256(t *testing.T) ***REMOVED***
	pub, _ := getTestKey()
	fingerprint := FingerprintSHA256(pub)
	want := "SHA256:Anr3LjZK8YVpjrxu79myrW9Hrb/wpcMNpVvTq/RcBm8" // ssh-keygen -lf rsa
	if fingerprint != want ***REMOVED***
		t.Errorf("got fingerprint %q want %q", fingerprint, want)
	***REMOVED***
***REMOVED***
