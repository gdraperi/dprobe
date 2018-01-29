// Copyright 2017 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package knownhosts

import (
	"bytes"
	"fmt"
	"net"
	"reflect"
	"testing"

	"golang.org/x/crypto/ssh"
)

const edKeyStr = "ssh-ed25519 AAAAC3NzaC1lZDI1NTE5AAAAIGBAarftlLeoyf+v+nVchEZII/vna2PCV8FaX4vsF5BX"
const alternateEdKeyStr = "ssh-ed25519 AAAAC3NzaC1lZDI1NTE5AAAAIIXffBYeYL+WVzVru8npl5JHt2cjlr4ornFTWzoij9sx"
const ecKeyStr = "ecdsa-sha2-nistp256 AAAAE2VjZHNhLXNoYTItbmlzdHAyNTYAAAAIbmlzdHAyNTYAAABBBNLCu01+wpXe3xB5olXCN4SqU2rQu0qjSRKJO4Bg+JRCPU+ENcgdA5srTU8xYDz/GEa4dzK5ldPw4J/gZgSXCMs="

var ecKey, alternateEdKey, edKey ssh.PublicKey
var testAddr = &net.TCPAddr***REMOVED***
	IP:   net.IP***REMOVED***198, 41, 30, 196***REMOVED***,
	Port: 22,
***REMOVED***

var testAddr6 = &net.TCPAddr***REMOVED***
	IP: net.IP***REMOVED***198, 41, 30, 196,
		1, 2, 3, 4,
		1, 2, 3, 4,
		1, 2, 3, 4,
	***REMOVED***,
	Port: 22,
***REMOVED***

func init() ***REMOVED***
	var err error
	ecKey, _, _, _, err = ssh.ParseAuthorizedKey([]byte(ecKeyStr))
	if err != nil ***REMOVED***
		panic(err)
	***REMOVED***
	edKey, _, _, _, err = ssh.ParseAuthorizedKey([]byte(edKeyStr))
	if err != nil ***REMOVED***
		panic(err)
	***REMOVED***
	alternateEdKey, _, _, _, err = ssh.ParseAuthorizedKey([]byte(alternateEdKeyStr))
	if err != nil ***REMOVED***
		panic(err)
	***REMOVED***
***REMOVED***

func testDB(t *testing.T, s string) *hostKeyDB ***REMOVED***
	db := newHostKeyDB()
	if err := db.Read(bytes.NewBufferString(s), "testdb"); err != nil ***REMOVED***
		t.Fatalf("Read: %v", err)
	***REMOVED***

	return db
***REMOVED***

func TestRevoked(t *testing.T) ***REMOVED***
	db := testDB(t, "\n\n@revoked * "+edKeyStr+"\n")
	want := &RevokedError***REMOVED***
		Revoked: KnownKey***REMOVED***
			Key:      edKey,
			Filename: "testdb",
			Line:     3,
		***REMOVED***,
	***REMOVED***
	if err := db.check("", &net.TCPAddr***REMOVED***
		Port: 42,
	***REMOVED***, edKey); err == nil ***REMOVED***
		t.Fatal("no error for revoked key")
	***REMOVED*** else if !reflect.DeepEqual(want, err) ***REMOVED***
		t.Fatalf("got %#v, want %#v", want, err)
	***REMOVED***
***REMOVED***

func TestHostAuthority(t *testing.T) ***REMOVED***
	for _, m := range []struct ***REMOVED***
		authorityFor string
		address      string

		good bool
	***REMOVED******REMOVED***
		***REMOVED***authorityFor: "localhost", address: "localhost:22", good: true***REMOVED***,
		***REMOVED***authorityFor: "localhost", address: "localhost", good: false***REMOVED***,
		***REMOVED***authorityFor: "localhost", address: "localhost:1234", good: false***REMOVED***,
		***REMOVED***authorityFor: "[localhost]:1234", address: "localhost:1234", good: true***REMOVED***,
		***REMOVED***authorityFor: "[localhost]:1234", address: "localhost:22", good: false***REMOVED***,
		***REMOVED***authorityFor: "[localhost]:1234", address: "localhost", good: false***REMOVED***,
	***REMOVED*** ***REMOVED***
		db := testDB(t, `@cert-authority `+m.authorityFor+` `+edKeyStr)
		if ok := db.IsHostAuthority(db.lines[0].knownKey.Key, m.address); ok != m.good ***REMOVED***
			t.Errorf("IsHostAuthority: authority %s, address %s, wanted good = %v, got good = %v",
				m.authorityFor, m.address, m.good, ok)
		***REMOVED***
	***REMOVED***
***REMOVED***

func TestBracket(t *testing.T) ***REMOVED***
	db := testDB(t, `[git.eclipse.org]:29418,[198.41.30.196]:29418 `+edKeyStr)

	if err := db.check("git.eclipse.org:29418", &net.TCPAddr***REMOVED***
		IP:   net.IP***REMOVED***198, 41, 30, 196***REMOVED***,
		Port: 29418,
	***REMOVED***, edKey); err != nil ***REMOVED***
		t.Errorf("got error %v, want none", err)
	***REMOVED***

	if err := db.check("git.eclipse.org:29419", &net.TCPAddr***REMOVED***
		Port: 42,
	***REMOVED***, edKey); err == nil ***REMOVED***
		t.Fatalf("no error for unknown address")
	***REMOVED*** else if ke, ok := err.(*KeyError); !ok ***REMOVED***
		t.Fatalf("got type %T, want *KeyError", err)
	***REMOVED*** else if len(ke.Want) > 0 ***REMOVED***
		t.Fatalf("got Want %v, want []", ke.Want)
	***REMOVED***
***REMOVED***

func TestNewKeyType(t *testing.T) ***REMOVED***
	str := fmt.Sprintf("%s %s", testAddr, edKeyStr)
	db := testDB(t, str)
	if err := db.check("", testAddr, ecKey); err == nil ***REMOVED***
		t.Fatalf("no error for unknown address")
	***REMOVED*** else if ke, ok := err.(*KeyError); !ok ***REMOVED***
		t.Fatalf("got type %T, want *KeyError", err)
	***REMOVED*** else if len(ke.Want) == 0 ***REMOVED***
		t.Fatalf("got empty KeyError.Want")
	***REMOVED***
***REMOVED***

func TestSameKeyType(t *testing.T) ***REMOVED***
	str := fmt.Sprintf("%s %s", testAddr, edKeyStr)
	db := testDB(t, str)
	if err := db.check("", testAddr, alternateEdKey); err == nil ***REMOVED***
		t.Fatalf("no error for unknown address")
	***REMOVED*** else if ke, ok := err.(*KeyError); !ok ***REMOVED***
		t.Fatalf("got type %T, want *KeyError", err)
	***REMOVED*** else if len(ke.Want) == 0 ***REMOVED***
		t.Fatalf("got empty KeyError.Want")
	***REMOVED*** else if got, want := ke.Want[0].Key.Marshal(), edKey.Marshal(); !bytes.Equal(got, want) ***REMOVED***
		t.Fatalf("got key %q, want %q", got, want)
	***REMOVED***
***REMOVED***

func TestIPAddress(t *testing.T) ***REMOVED***
	str := fmt.Sprintf("%s %s", testAddr, edKeyStr)
	db := testDB(t, str)
	if err := db.check("", testAddr, edKey); err != nil ***REMOVED***
		t.Errorf("got error %q, want none", err)
	***REMOVED***
***REMOVED***

func TestIPv6Address(t *testing.T) ***REMOVED***
	str := fmt.Sprintf("%s %s", testAddr6, edKeyStr)
	db := testDB(t, str)

	if err := db.check("", testAddr6, edKey); err != nil ***REMOVED***
		t.Errorf("got error %q, want none", err)
	***REMOVED***
***REMOVED***

func TestBasic(t *testing.T) ***REMOVED***
	str := fmt.Sprintf("#comment\n\nserver.org,%s %s\notherhost %s", testAddr, edKeyStr, ecKeyStr)
	db := testDB(t, str)
	if err := db.check("server.org:22", testAddr, edKey); err != nil ***REMOVED***
		t.Errorf("got error %q, want none", err)
	***REMOVED***

	want := KnownKey***REMOVED***
		Key:      edKey,
		Filename: "testdb",
		Line:     3,
	***REMOVED***
	if err := db.check("server.org:22", testAddr, ecKey); err == nil ***REMOVED***
		t.Errorf("succeeded, want KeyError")
	***REMOVED*** else if ke, ok := err.(*KeyError); !ok ***REMOVED***
		t.Errorf("got %T, want *KeyError", err)
	***REMOVED*** else if len(ke.Want) != 1 ***REMOVED***
		t.Errorf("got %v, want 1 entry", ke)
	***REMOVED*** else if !reflect.DeepEqual(ke.Want[0], want) ***REMOVED***
		t.Errorf("got %v, want %v", ke.Want[0], want)
	***REMOVED***
***REMOVED***

func TestNegate(t *testing.T) ***REMOVED***
	str := fmt.Sprintf("%s,!server.org %s", testAddr, edKeyStr)
	db := testDB(t, str)
	if err := db.check("server.org:22", testAddr, ecKey); err == nil ***REMOVED***
		t.Errorf("succeeded")
	***REMOVED*** else if ke, ok := err.(*KeyError); !ok ***REMOVED***
		t.Errorf("got error type %T, want *KeyError", err)
	***REMOVED*** else if len(ke.Want) != 0 ***REMOVED***
		t.Errorf("got expected keys %d (first of type %s), want []", len(ke.Want), ke.Want[0].Key.Type())
	***REMOVED***
***REMOVED***

func TestWildcard(t *testing.T) ***REMOVED***
	str := fmt.Sprintf("server*.domain %s", edKeyStr)
	db := testDB(t, str)

	want := &KeyError***REMOVED***
		Want: []KnownKey***REMOVED******REMOVED***
			Filename: "testdb",
			Line:     1,
			Key:      edKey,
		***REMOVED******REMOVED***,
	***REMOVED***

	got := db.check("server.domain:22", &net.TCPAddr***REMOVED******REMOVED***, ecKey)
	if !reflect.DeepEqual(got, want) ***REMOVED***
		t.Errorf("got %s, want %s", got, want)
	***REMOVED***
***REMOVED***

func TestLine(t *testing.T) ***REMOVED***
	for in, want := range map[string]string***REMOVED***
		"server.org":                             "server.org " + edKeyStr,
		"server.org:22":                          "server.org " + edKeyStr,
		"server.org:23":                          "[server.org]:23 " + edKeyStr,
		"[c629:1ec4:102:304:102:304:102:304]:22": "[c629:1ec4:102:304:102:304:102:304] " + edKeyStr,
		"[c629:1ec4:102:304:102:304:102:304]:23": "[c629:1ec4:102:304:102:304:102:304]:23 " + edKeyStr,
	***REMOVED*** ***REMOVED***
		if got := Line([]string***REMOVED***in***REMOVED***, edKey); got != want ***REMOVED***
			t.Errorf("Line(%q) = %q, want %q", in, got, want)
		***REMOVED***
	***REMOVED***
***REMOVED***

func TestWildcardMatch(t *testing.T) ***REMOVED***
	for _, c := range []struct ***REMOVED***
		pat, str string
		want     bool
	***REMOVED******REMOVED***
		***REMOVED***"a?b", "abb", true***REMOVED***,
		***REMOVED***"ab", "abc", false***REMOVED***,
		***REMOVED***"abc", "ab", false***REMOVED***,
		***REMOVED***"a*b", "axxxb", true***REMOVED***,
		***REMOVED***"a*b", "axbxb", true***REMOVED***,
		***REMOVED***"a*b", "axbxbc", false***REMOVED***,
		***REMOVED***"a*?", "axbxc", true***REMOVED***,
		***REMOVED***"a*b*", "axxbxxxxxx", true***REMOVED***,
		***REMOVED***"a*b*c", "axxbxxxxxxc", true***REMOVED***,
		***REMOVED***"a*b*?", "axxbxxxxxxc", true***REMOVED***,
		***REMOVED***"a*b*z", "axxbxxbxxxz", true***REMOVED***,
		***REMOVED***"a*b*z", "axxbxxzxxxz", true***REMOVED***,
		***REMOVED***"a*b*z", "axxbxxzxxx", false***REMOVED***,
	***REMOVED*** ***REMOVED***
		got := wildcardMatch([]byte(c.pat), []byte(c.str))
		if got != c.want ***REMOVED***
			t.Errorf("wildcardMatch(%q, %q) = %v, want %v", c.pat, c.str, got, c.want)
		***REMOVED***

	***REMOVED***
***REMOVED***

// TODO(hanwen): test coverage for certificates.

const testHostname = "hostname"

// generated with keygen -H -f
const encodedTestHostnameHash = "|1|IHXZvQMvTcZTUU29+2vXFgx8Frs=|UGccIWfRVDwilMBnA3WJoRAC75Y="

func TestHostHash(t *testing.T) ***REMOVED***
	testHostHash(t, testHostname, encodedTestHostnameHash)
***REMOVED***

func TestHashList(t *testing.T) ***REMOVED***
	encoded := HashHostname(testHostname)
	testHostHash(t, testHostname, encoded)
***REMOVED***

func testHostHash(t *testing.T, hostname, encoded string) ***REMOVED***
	typ, salt, hash, err := decodeHash(encoded)
	if err != nil ***REMOVED***
		t.Fatalf("decodeHash: %v", err)
	***REMOVED***

	if got := encodeHash(typ, salt, hash); got != encoded ***REMOVED***
		t.Errorf("got encoding %s want %s", got, encoded)
	***REMOVED***

	if typ != sha1HashType ***REMOVED***
		t.Fatalf("got hash type %q, want %q", typ, sha1HashType)
	***REMOVED***

	got := hashHost(hostname, salt)
	if !bytes.Equal(got, hash) ***REMOVED***
		t.Errorf("got hash %x want %x", got, hash)
	***REMOVED***
***REMOVED***

func TestNormalize(t *testing.T) ***REMOVED***
	for in, want := range map[string]string***REMOVED***
		"127.0.0.1:22":             "127.0.0.1",
		"[127.0.0.1]:22":           "127.0.0.1",
		"[127.0.0.1]:23":           "[127.0.0.1]:23",
		"127.0.0.1:23":             "[127.0.0.1]:23",
		"[a.b.c]:22":               "a.b.c",
		"[abcd:abcd:abcd:abcd]":    "[abcd:abcd:abcd:abcd]",
		"[abcd:abcd:abcd:abcd]:22": "[abcd:abcd:abcd:abcd]",
		"[abcd:abcd:abcd:abcd]:23": "[abcd:abcd:abcd:abcd]:23",
	***REMOVED*** ***REMOVED***
		got := Normalize(in)
		if got != want ***REMOVED***
			t.Errorf("Normalize(%q) = %q, want %q", in, got, want)
		***REMOVED***
	***REMOVED***
***REMOVED***

func TestHashedHostkeyCheck(t *testing.T) ***REMOVED***
	str := fmt.Sprintf("%s %s", HashHostname(testHostname), edKeyStr)
	db := testDB(t, str)
	if err := db.check(testHostname+":22", testAddr, edKey); err != nil ***REMOVED***
		t.Errorf("check(%s): %v", testHostname, err)
	***REMOVED***
	want := &KeyError***REMOVED***
		Want: []KnownKey***REMOVED******REMOVED***
			Filename: "testdb",
			Line:     1,
			Key:      edKey,
		***REMOVED******REMOVED***,
	***REMOVED***
	if got := db.check(testHostname+":22", testAddr, alternateEdKey); !reflect.DeepEqual(got, want) ***REMOVED***
		t.Errorf("got error %v, want %v", got, want)
	***REMOVED***
***REMOVED***
