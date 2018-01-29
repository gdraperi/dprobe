// Copyright 2017 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package knownhosts implements a parser for the OpenSSH
// known_hosts host key database.
package knownhosts

import (
	"bufio"
	"bytes"
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha1"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"net"
	"os"
	"strings"

	"golang.org/x/crypto/ssh"
)

// See the sshd manpage
// (http://man.openbsd.org/sshd#SSH_KNOWN_HOSTS_FILE_FORMAT) for
// background.

type addr struct***REMOVED*** host, port string ***REMOVED***

func (a *addr) String() string ***REMOVED***
	h := a.host
	if strings.Contains(h, ":") ***REMOVED***
		h = "[" + h + "]"
	***REMOVED***
	return h + ":" + a.port
***REMOVED***

type matcher interface ***REMOVED***
	match([]addr) bool
***REMOVED***

type hostPattern struct ***REMOVED***
	negate bool
	addr   addr
***REMOVED***

func (p *hostPattern) String() string ***REMOVED***
	n := ""
	if p.negate ***REMOVED***
		n = "!"
	***REMOVED***

	return n + p.addr.String()
***REMOVED***

type hostPatterns []hostPattern

func (ps hostPatterns) match(addrs []addr) bool ***REMOVED***
	matched := false
	for _, p := range ps ***REMOVED***
		for _, a := range addrs ***REMOVED***
			m := p.match(a)
			if !m ***REMOVED***
				continue
			***REMOVED***
			if p.negate ***REMOVED***
				return false
			***REMOVED***
			matched = true
		***REMOVED***
	***REMOVED***
	return matched
***REMOVED***

// See
// https://android.googlesource.com/platform/external/openssh/+/ab28f5495c85297e7a597c1ba62e996416da7c7e/addrmatch.c
// The matching of * has no regard for separators, unlike filesystem globs
func wildcardMatch(pat []byte, str []byte) bool ***REMOVED***
	for ***REMOVED***
		if len(pat) == 0 ***REMOVED***
			return len(str) == 0
		***REMOVED***
		if len(str) == 0 ***REMOVED***
			return false
		***REMOVED***

		if pat[0] == '*' ***REMOVED***
			if len(pat) == 1 ***REMOVED***
				return true
			***REMOVED***

			for j := range str ***REMOVED***
				if wildcardMatch(pat[1:], str[j:]) ***REMOVED***
					return true
				***REMOVED***
			***REMOVED***
			return false
		***REMOVED***

		if pat[0] == '?' || pat[0] == str[0] ***REMOVED***
			pat = pat[1:]
			str = str[1:]
		***REMOVED*** else ***REMOVED***
			return false
		***REMOVED***
	***REMOVED***
***REMOVED***

func (p *hostPattern) match(a addr) bool ***REMOVED***
	return wildcardMatch([]byte(p.addr.host), []byte(a.host)) && p.addr.port == a.port
***REMOVED***

type keyDBLine struct ***REMOVED***
	cert     bool
	matcher  matcher
	knownKey KnownKey
***REMOVED***

func serialize(k ssh.PublicKey) string ***REMOVED***
	return k.Type() + " " + base64.StdEncoding.EncodeToString(k.Marshal())
***REMOVED***

func (l *keyDBLine) match(addrs []addr) bool ***REMOVED***
	return l.matcher.match(addrs)
***REMOVED***

type hostKeyDB struct ***REMOVED***
	// Serialized version of revoked keys
	revoked map[string]*KnownKey
	lines   []keyDBLine
***REMOVED***

func newHostKeyDB() *hostKeyDB ***REMOVED***
	db := &hostKeyDB***REMOVED***
		revoked: make(map[string]*KnownKey),
	***REMOVED***

	return db
***REMOVED***

func keyEq(a, b ssh.PublicKey) bool ***REMOVED***
	return bytes.Equal(a.Marshal(), b.Marshal())
***REMOVED***

// IsAuthorityForHost can be used as a callback in ssh.CertChecker
func (db *hostKeyDB) IsHostAuthority(remote ssh.PublicKey, address string) bool ***REMOVED***
	h, p, err := net.SplitHostPort(address)
	if err != nil ***REMOVED***
		return false
	***REMOVED***
	a := addr***REMOVED***host: h, port: p***REMOVED***

	for _, l := range db.lines ***REMOVED***
		if l.cert && keyEq(l.knownKey.Key, remote) && l.match([]addr***REMOVED***a***REMOVED***) ***REMOVED***
			return true
		***REMOVED***
	***REMOVED***
	return false
***REMOVED***

// IsRevoked can be used as a callback in ssh.CertChecker
func (db *hostKeyDB) IsRevoked(key *ssh.Certificate) bool ***REMOVED***
	_, ok := db.revoked[string(key.Marshal())]
	return ok
***REMOVED***

const markerCert = "@cert-authority"
const markerRevoked = "@revoked"

func nextWord(line []byte) (string, []byte) ***REMOVED***
	i := bytes.IndexAny(line, "\t ")
	if i == -1 ***REMOVED***
		return string(line), nil
	***REMOVED***

	return string(line[:i]), bytes.TrimSpace(line[i:])
***REMOVED***

func parseLine(line []byte) (marker, host string, key ssh.PublicKey, err error) ***REMOVED***
	if w, next := nextWord(line); w == markerCert || w == markerRevoked ***REMOVED***
		marker = w
		line = next
	***REMOVED***

	host, line = nextWord(line)
	if len(line) == 0 ***REMOVED***
		return "", "", nil, errors.New("knownhosts: missing host pattern")
	***REMOVED***

	// ignore the keytype as it's in the key blob anyway.
	_, line = nextWord(line)
	if len(line) == 0 ***REMOVED***
		return "", "", nil, errors.New("knownhosts: missing key type pattern")
	***REMOVED***

	keyBlob, _ := nextWord(line)

	keyBytes, err := base64.StdEncoding.DecodeString(keyBlob)
	if err != nil ***REMOVED***
		return "", "", nil, err
	***REMOVED***
	key, err = ssh.ParsePublicKey(keyBytes)
	if err != nil ***REMOVED***
		return "", "", nil, err
	***REMOVED***

	return marker, host, key, nil
***REMOVED***

func (db *hostKeyDB) parseLine(line []byte, filename string, linenum int) error ***REMOVED***
	marker, pattern, key, err := parseLine(line)
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	if marker == markerRevoked ***REMOVED***
		db.revoked[string(key.Marshal())] = &KnownKey***REMOVED***
			Key:      key,
			Filename: filename,
			Line:     linenum,
		***REMOVED***

		return nil
	***REMOVED***

	entry := keyDBLine***REMOVED***
		cert: marker == markerCert,
		knownKey: KnownKey***REMOVED***
			Filename: filename,
			Line:     linenum,
			Key:      key,
		***REMOVED***,
	***REMOVED***

	if pattern[0] == '|' ***REMOVED***
		entry.matcher, err = newHashedHost(pattern)
	***REMOVED*** else ***REMOVED***
		entry.matcher, err = newHostnameMatcher(pattern)
	***REMOVED***

	if err != nil ***REMOVED***
		return err
	***REMOVED***

	db.lines = append(db.lines, entry)
	return nil
***REMOVED***

func newHostnameMatcher(pattern string) (matcher, error) ***REMOVED***
	var hps hostPatterns
	for _, p := range strings.Split(pattern, ",") ***REMOVED***
		if len(p) == 0 ***REMOVED***
			continue
		***REMOVED***

		var a addr
		var negate bool
		if p[0] == '!' ***REMOVED***
			negate = true
			p = p[1:]
		***REMOVED***

		if len(p) == 0 ***REMOVED***
			return nil, errors.New("knownhosts: negation without following hostname")
		***REMOVED***

		var err error
		if p[0] == '[' ***REMOVED***
			a.host, a.port, err = net.SplitHostPort(p)
			if err != nil ***REMOVED***
				return nil, err
			***REMOVED***
		***REMOVED*** else ***REMOVED***
			a.host, a.port, err = net.SplitHostPort(p)
			if err != nil ***REMOVED***
				a.host = p
				a.port = "22"
			***REMOVED***
		***REMOVED***
		hps = append(hps, hostPattern***REMOVED***
			negate: negate,
			addr:   a,
		***REMOVED***)
	***REMOVED***
	return hps, nil
***REMOVED***

// KnownKey represents a key declared in a known_hosts file.
type KnownKey struct ***REMOVED***
	Key      ssh.PublicKey
	Filename string
	Line     int
***REMOVED***

func (k *KnownKey) String() string ***REMOVED***
	return fmt.Sprintf("%s:%d: %s", k.Filename, k.Line, serialize(k.Key))
***REMOVED***

// KeyError is returned if we did not find the key in the host key
// database, or there was a mismatch.  Typically, in batch
// applications, this should be interpreted as failure. Interactive
// applications can offer an interactive prompt to the user.
type KeyError struct ***REMOVED***
	// Want holds the accepted host keys. For each key algorithm,
	// there can be one hostkey.  If Want is empty, the host is
	// unknown. If Want is non-empty, there was a mismatch, which
	// can signify a MITM attack.
	Want []KnownKey
***REMOVED***

func (u *KeyError) Error() string ***REMOVED***
	if len(u.Want) == 0 ***REMOVED***
		return "knownhosts: key is unknown"
	***REMOVED***
	return "knownhosts: key mismatch"
***REMOVED***

// RevokedError is returned if we found a key that was revoked.
type RevokedError struct ***REMOVED***
	Revoked KnownKey
***REMOVED***

func (r *RevokedError) Error() string ***REMOVED***
	return "knownhosts: key is revoked"
***REMOVED***

// check checks a key against the host database. This should not be
// used for verifying certificates.
func (db *hostKeyDB) check(address string, remote net.Addr, remoteKey ssh.PublicKey) error ***REMOVED***
	if revoked := db.revoked[string(remoteKey.Marshal())]; revoked != nil ***REMOVED***
		return &RevokedError***REMOVED***Revoked: *revoked***REMOVED***
	***REMOVED***

	host, port, err := net.SplitHostPort(remote.String())
	if err != nil ***REMOVED***
		return fmt.Errorf("knownhosts: SplitHostPort(%s): %v", remote, err)
	***REMOVED***

	addrs := []addr***REMOVED***
		***REMOVED***host, port***REMOVED***,
	***REMOVED***

	if address != "" ***REMOVED***
		host, port, err := net.SplitHostPort(address)
		if err != nil ***REMOVED***
			return fmt.Errorf("knownhosts: SplitHostPort(%s): %v", address, err)
		***REMOVED***

		addrs = append(addrs, addr***REMOVED***host, port***REMOVED***)
	***REMOVED***

	return db.checkAddrs(addrs, remoteKey)
***REMOVED***

// checkAddrs checks if we can find the given public key for any of
// the given addresses.  If we only find an entry for the IP address,
// or only the hostname, then this still succeeds.
func (db *hostKeyDB) checkAddrs(addrs []addr, remoteKey ssh.PublicKey) error ***REMOVED***
	// TODO(hanwen): are these the right semantics? What if there
	// is just a key for the IP address, but not for the
	// hostname?

	// Algorithm => key.
	knownKeys := map[string]KnownKey***REMOVED******REMOVED***
	for _, l := range db.lines ***REMOVED***
		if l.match(addrs) ***REMOVED***
			typ := l.knownKey.Key.Type()
			if _, ok := knownKeys[typ]; !ok ***REMOVED***
				knownKeys[typ] = l.knownKey
			***REMOVED***
		***REMOVED***
	***REMOVED***

	keyErr := &KeyError***REMOVED******REMOVED***
	for _, v := range knownKeys ***REMOVED***
		keyErr.Want = append(keyErr.Want, v)
	***REMOVED***

	// Unknown remote host.
	if len(knownKeys) == 0 ***REMOVED***
		return keyErr
	***REMOVED***

	// If the remote host starts using a different, unknown key type, we
	// also interpret that as a mismatch.
	if known, ok := knownKeys[remoteKey.Type()]; !ok || !keyEq(known.Key, remoteKey) ***REMOVED***
		return keyErr
	***REMOVED***

	return nil
***REMOVED***

// The Read function parses file contents.
func (db *hostKeyDB) Read(r io.Reader, filename string) error ***REMOVED***
	scanner := bufio.NewScanner(r)

	lineNum := 0
	for scanner.Scan() ***REMOVED***
		lineNum++
		line := scanner.Bytes()
		line = bytes.TrimSpace(line)
		if len(line) == 0 || line[0] == '#' ***REMOVED***
			continue
		***REMOVED***

		if err := db.parseLine(line, filename, lineNum); err != nil ***REMOVED***
			return fmt.Errorf("knownhosts: %s:%d: %v", filename, lineNum, err)
		***REMOVED***
	***REMOVED***
	return scanner.Err()
***REMOVED***

// New creates a host key callback from the given OpenSSH host key
// files. The returned callback is for use in
// ssh.ClientConfig.HostKeyCallback.
func New(files ...string) (ssh.HostKeyCallback, error) ***REMOVED***
	db := newHostKeyDB()
	for _, fn := range files ***REMOVED***
		f, err := os.Open(fn)
		if err != nil ***REMOVED***
			return nil, err
		***REMOVED***
		defer f.Close()
		if err := db.Read(f, fn); err != nil ***REMOVED***
			return nil, err
		***REMOVED***
	***REMOVED***

	var certChecker ssh.CertChecker
	certChecker.IsHostAuthority = db.IsHostAuthority
	certChecker.IsRevoked = db.IsRevoked
	certChecker.HostKeyFallback = db.check

	return certChecker.CheckHostKey, nil
***REMOVED***

// Normalize normalizes an address into the form used in known_hosts
func Normalize(address string) string ***REMOVED***
	host, port, err := net.SplitHostPort(address)
	if err != nil ***REMOVED***
		host = address
		port = "22"
	***REMOVED***
	entry := host
	if port != "22" ***REMOVED***
		entry = "[" + entry + "]:" + port
	***REMOVED*** else if strings.Contains(host, ":") && !strings.HasPrefix(host, "[") ***REMOVED***
		entry = "[" + entry + "]"
	***REMOVED***
	return entry
***REMOVED***

// Line returns a line to add append to the known_hosts files.
func Line(addresses []string, key ssh.PublicKey) string ***REMOVED***
	var trimmed []string
	for _, a := range addresses ***REMOVED***
		trimmed = append(trimmed, Normalize(a))
	***REMOVED***

	return strings.Join(trimmed, ",") + " " + serialize(key)
***REMOVED***

// HashHostname hashes the given hostname. The hostname is not
// normalized before hashing.
func HashHostname(hostname string) string ***REMOVED***
	// TODO(hanwen): check if we can safely normalize this always.
	salt := make([]byte, sha1.Size)

	_, err := rand.Read(salt)
	if err != nil ***REMOVED***
		panic(fmt.Sprintf("crypto/rand failure %v", err))
	***REMOVED***

	hash := hashHost(hostname, salt)
	return encodeHash(sha1HashType, salt, hash)
***REMOVED***

func decodeHash(encoded string) (hashType string, salt, hash []byte, err error) ***REMOVED***
	if len(encoded) == 0 || encoded[0] != '|' ***REMOVED***
		err = errors.New("knownhosts: hashed host must start with '|'")
		return
	***REMOVED***
	components := strings.Split(encoded, "|")
	if len(components) != 4 ***REMOVED***
		err = fmt.Errorf("knownhosts: got %d components, want 3", len(components))
		return
	***REMOVED***

	hashType = components[1]
	if salt, err = base64.StdEncoding.DecodeString(components[2]); err != nil ***REMOVED***
		return
	***REMOVED***
	if hash, err = base64.StdEncoding.DecodeString(components[3]); err != nil ***REMOVED***
		return
	***REMOVED***
	return
***REMOVED***

func encodeHash(typ string, salt []byte, hash []byte) string ***REMOVED***
	return strings.Join([]string***REMOVED***"",
		typ,
		base64.StdEncoding.EncodeToString(salt),
		base64.StdEncoding.EncodeToString(hash),
	***REMOVED***, "|")
***REMOVED***

// See https://android.googlesource.com/platform/external/openssh/+/ab28f5495c85297e7a597c1ba62e996416da7c7e/hostfile.c#120
func hashHost(hostname string, salt []byte) []byte ***REMOVED***
	mac := hmac.New(sha1.New, salt)
	mac.Write([]byte(hostname))
	return mac.Sum(nil)
***REMOVED***

type hashedHost struct ***REMOVED***
	salt []byte
	hash []byte
***REMOVED***

const sha1HashType = "1"

func newHashedHost(encoded string) (*hashedHost, error) ***REMOVED***
	typ, salt, hash, err := decodeHash(encoded)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	// The type field seems for future algorithm agility, but it's
	// actually hardcoded in openssh currently, see
	// https://android.googlesource.com/platform/external/openssh/+/ab28f5495c85297e7a597c1ba62e996416da7c7e/hostfile.c#120
	if typ != sha1HashType ***REMOVED***
		return nil, fmt.Errorf("knownhosts: got hash type %s, must be '1'", typ)
	***REMOVED***

	return &hashedHost***REMOVED***salt: salt, hash: hash***REMOVED***, nil
***REMOVED***

func (h *hashedHost) match(addrs []addr) bool ***REMOVED***
	for _, a := range addrs ***REMOVED***
		if bytes.Equal(hashHost(Normalize(a.String()), h.salt), h.hash) ***REMOVED***
			return true
		***REMOVED***
	***REMOVED***
	return false
***REMOVED***
