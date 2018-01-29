// Copyright 2012 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package agent implements the ssh-agent protocol, and provides both
// a client and a server. The client can talk to a standard ssh-agent
// that uses UNIX sockets, and one could implement an alternative
// ssh-agent process using the sample server.
//
// References:
//  [PROTOCOL.agent]:    http://cvsweb.openbsd.org/cgi-bin/cvsweb/src/usr.bin/ssh/PROTOCOL.agent?rev=HEAD
package agent // import "golang.org/x/crypto/ssh/agent"

import (
	"bytes"
	"crypto/dsa"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rsa"
	"encoding/base64"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"math/big"
	"sync"

	"golang.org/x/crypto/ed25519"
	"golang.org/x/crypto/ssh"
)

// Agent represents the capabilities of an ssh-agent.
type Agent interface ***REMOVED***
	// List returns the identities known to the agent.
	List() ([]*Key, error)

	// Sign has the agent sign the data using a protocol 2 key as defined
	// in [PROTOCOL.agent] section 2.6.2.
	Sign(key ssh.PublicKey, data []byte) (*ssh.Signature, error)

	// Add adds a private key to the agent.
	Add(key AddedKey) error

	// Remove removes all identities with the given public key.
	Remove(key ssh.PublicKey) error

	// RemoveAll removes all identities.
	RemoveAll() error

	// Lock locks the agent. Sign and Remove will fail, and List will empty an empty list.
	Lock(passphrase []byte) error

	// Unlock undoes the effect of Lock
	Unlock(passphrase []byte) error

	// Signers returns signers for all the known keys.
	Signers() ([]ssh.Signer, error)
***REMOVED***

// ConstraintExtension describes an optional constraint defined by users.
type ConstraintExtension struct ***REMOVED***
	// ExtensionName consist of a UTF-8 string suffixed by the
	// implementation domain following the naming scheme defined
	// in Section 4.2 of [RFC4251], e.g.  "foo@example.com".
	ExtensionName string
	// ExtensionDetails contains the actual content of the extended
	// constraint.
	ExtensionDetails []byte
***REMOVED***

// AddedKey describes an SSH key to be added to an Agent.
type AddedKey struct ***REMOVED***
	// PrivateKey must be a *rsa.PrivateKey, *dsa.PrivateKey or
	// *ecdsa.PrivateKey, which will be inserted into the agent.
	PrivateKey interface***REMOVED******REMOVED***
	// Certificate, if not nil, is communicated to the agent and will be
	// stored with the key.
	Certificate *ssh.Certificate
	// Comment is an optional, free-form string.
	Comment string
	// LifetimeSecs, if not zero, is the number of seconds that the
	// agent will store the key for.
	LifetimeSecs uint32
	// ConfirmBeforeUse, if true, requests that the agent confirm with the
	// user before each use of this key.
	ConfirmBeforeUse bool
	// ConstraintExtensions are the experimental or private-use constraints
	// defined by users.
	ConstraintExtensions []ConstraintExtension
***REMOVED***

// See [PROTOCOL.agent], section 3.
const (
	agentRequestV1Identities   = 1
	agentRemoveAllV1Identities = 9

	// 3.2 Requests from client to agent for protocol 2 key operations
	agentAddIdentity         = 17
	agentRemoveIdentity      = 18
	agentRemoveAllIdentities = 19
	agentAddIDConstrained    = 25

	// 3.3 Key-type independent requests from client to agent
	agentAddSmartcardKey            = 20
	agentRemoveSmartcardKey         = 21
	agentLock                       = 22
	agentUnlock                     = 23
	agentAddSmartcardKeyConstrained = 26

	// 3.7 Key constraint identifiers
	agentConstrainLifetime  = 1
	agentConstrainConfirm   = 2
	agentConstrainExtension = 3
)

// maxAgentResponseBytes is the maximum agent reply size that is accepted. This
// is a sanity check, not a limit in the spec.
const maxAgentResponseBytes = 16 << 20

// Agent messages:
// These structures mirror the wire format of the corresponding ssh agent
// messages found in [PROTOCOL.agent].

// 3.4 Generic replies from agent to client
const agentFailure = 5

type failureAgentMsg struct***REMOVED******REMOVED***

const agentSuccess = 6

type successAgentMsg struct***REMOVED******REMOVED***

// See [PROTOCOL.agent], section 2.5.2.
const agentRequestIdentities = 11

type requestIdentitiesAgentMsg struct***REMOVED******REMOVED***

// See [PROTOCOL.agent], section 2.5.2.
const agentIdentitiesAnswer = 12

type identitiesAnswerAgentMsg struct ***REMOVED***
	NumKeys uint32 `sshtype:"12"`
	Keys    []byte `ssh:"rest"`
***REMOVED***

// See [PROTOCOL.agent], section 2.6.2.
const agentSignRequest = 13

type signRequestAgentMsg struct ***REMOVED***
	KeyBlob []byte `sshtype:"13"`
	Data    []byte
	Flags   uint32
***REMOVED***

// See [PROTOCOL.agent], section 2.6.2.

// 3.6 Replies from agent to client for protocol 2 key operations
const agentSignResponse = 14

type signResponseAgentMsg struct ***REMOVED***
	SigBlob []byte `sshtype:"14"`
***REMOVED***

type publicKey struct ***REMOVED***
	Format string
	Rest   []byte `ssh:"rest"`
***REMOVED***

// 3.7 Key constraint identifiers
type constrainLifetimeAgentMsg struct ***REMOVED***
	LifetimeSecs uint32 `sshtype:"1"`
***REMOVED***

type constrainExtensionAgentMsg struct ***REMOVED***
	ExtensionName    string `sshtype:"3"`
	ExtensionDetails []byte

	// Rest is a field used for parsing, not part of message
	Rest []byte `ssh:"rest"`
***REMOVED***

// Key represents a protocol 2 public key as defined in
// [PROTOCOL.agent], section 2.5.2.
type Key struct ***REMOVED***
	Format  string
	Blob    []byte
	Comment string
***REMOVED***

func clientErr(err error) error ***REMOVED***
	return fmt.Errorf("agent: client error: %v", err)
***REMOVED***

// String returns the storage form of an agent key with the format, base64
// encoded serialized key, and the comment if it is not empty.
func (k *Key) String() string ***REMOVED***
	s := string(k.Format) + " " + base64.StdEncoding.EncodeToString(k.Blob)

	if k.Comment != "" ***REMOVED***
		s += " " + k.Comment
	***REMOVED***

	return s
***REMOVED***

// Type returns the public key type.
func (k *Key) Type() string ***REMOVED***
	return k.Format
***REMOVED***

// Marshal returns key blob to satisfy the ssh.PublicKey interface.
func (k *Key) Marshal() []byte ***REMOVED***
	return k.Blob
***REMOVED***

// Verify satisfies the ssh.PublicKey interface.
func (k *Key) Verify(data []byte, sig *ssh.Signature) error ***REMOVED***
	pubKey, err := ssh.ParsePublicKey(k.Blob)
	if err != nil ***REMOVED***
		return fmt.Errorf("agent: bad public key: %v", err)
	***REMOVED***
	return pubKey.Verify(data, sig)
***REMOVED***

type wireKey struct ***REMOVED***
	Format string
	Rest   []byte `ssh:"rest"`
***REMOVED***

func parseKey(in []byte) (out *Key, rest []byte, err error) ***REMOVED***
	var record struct ***REMOVED***
		Blob    []byte
		Comment string
		Rest    []byte `ssh:"rest"`
	***REMOVED***

	if err := ssh.Unmarshal(in, &record); err != nil ***REMOVED***
		return nil, nil, err
	***REMOVED***

	var wk wireKey
	if err := ssh.Unmarshal(record.Blob, &wk); err != nil ***REMOVED***
		return nil, nil, err
	***REMOVED***

	return &Key***REMOVED***
		Format:  wk.Format,
		Blob:    record.Blob,
		Comment: record.Comment,
	***REMOVED***, record.Rest, nil
***REMOVED***

// client is a client for an ssh-agent process.
type client struct ***REMOVED***
	// conn is typically a *net.UnixConn
	conn io.ReadWriter
	// mu is used to prevent concurrent access to the agent
	mu sync.Mutex
***REMOVED***

// NewClient returns an Agent that talks to an ssh-agent process over
// the given connection.
func NewClient(rw io.ReadWriter) Agent ***REMOVED***
	return &client***REMOVED***conn: rw***REMOVED***
***REMOVED***

// call sends an RPC to the agent. On success, the reply is
// unmarshaled into reply and replyType is set to the first byte of
// the reply, which contains the type of the message.
func (c *client) call(req []byte) (reply interface***REMOVED******REMOVED***, err error) ***REMOVED***
	c.mu.Lock()
	defer c.mu.Unlock()

	msg := make([]byte, 4+len(req))
	binary.BigEndian.PutUint32(msg, uint32(len(req)))
	copy(msg[4:], req)
	if _, err = c.conn.Write(msg); err != nil ***REMOVED***
		return nil, clientErr(err)
	***REMOVED***

	var respSizeBuf [4]byte
	if _, err = io.ReadFull(c.conn, respSizeBuf[:]); err != nil ***REMOVED***
		return nil, clientErr(err)
	***REMOVED***
	respSize := binary.BigEndian.Uint32(respSizeBuf[:])
	if respSize > maxAgentResponseBytes ***REMOVED***
		return nil, clientErr(err)
	***REMOVED***

	buf := make([]byte, respSize)
	if _, err = io.ReadFull(c.conn, buf); err != nil ***REMOVED***
		return nil, clientErr(err)
	***REMOVED***
	reply, err = unmarshal(buf)
	if err != nil ***REMOVED***
		return nil, clientErr(err)
	***REMOVED***
	return reply, err
***REMOVED***

func (c *client) simpleCall(req []byte) error ***REMOVED***
	resp, err := c.call(req)
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	if _, ok := resp.(*successAgentMsg); ok ***REMOVED***
		return nil
	***REMOVED***
	return errors.New("agent: failure")
***REMOVED***

func (c *client) RemoveAll() error ***REMOVED***
	return c.simpleCall([]byte***REMOVED***agentRemoveAllIdentities***REMOVED***)
***REMOVED***

func (c *client) Remove(key ssh.PublicKey) error ***REMOVED***
	req := ssh.Marshal(&agentRemoveIdentityMsg***REMOVED***
		KeyBlob: key.Marshal(),
	***REMOVED***)
	return c.simpleCall(req)
***REMOVED***

func (c *client) Lock(passphrase []byte) error ***REMOVED***
	req := ssh.Marshal(&agentLockMsg***REMOVED***
		Passphrase: passphrase,
	***REMOVED***)
	return c.simpleCall(req)
***REMOVED***

func (c *client) Unlock(passphrase []byte) error ***REMOVED***
	req := ssh.Marshal(&agentUnlockMsg***REMOVED***
		Passphrase: passphrase,
	***REMOVED***)
	return c.simpleCall(req)
***REMOVED***

// List returns the identities known to the agent.
func (c *client) List() ([]*Key, error) ***REMOVED***
	// see [PROTOCOL.agent] section 2.5.2.
	req := []byte***REMOVED***agentRequestIdentities***REMOVED***

	msg, err := c.call(req)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	switch msg := msg.(type) ***REMOVED***
	case *identitiesAnswerAgentMsg:
		if msg.NumKeys > maxAgentResponseBytes/8 ***REMOVED***
			return nil, errors.New("agent: too many keys in agent reply")
		***REMOVED***
		keys := make([]*Key, msg.NumKeys)
		data := msg.Keys
		for i := uint32(0); i < msg.NumKeys; i++ ***REMOVED***
			var key *Key
			var err error
			if key, data, err = parseKey(data); err != nil ***REMOVED***
				return nil, err
			***REMOVED***
			keys[i] = key
		***REMOVED***
		return keys, nil
	case *failureAgentMsg:
		return nil, errors.New("agent: failed to list keys")
	***REMOVED***
	panic("unreachable")
***REMOVED***

// Sign has the agent sign the data using a protocol 2 key as defined
// in [PROTOCOL.agent] section 2.6.2.
func (c *client) Sign(key ssh.PublicKey, data []byte) (*ssh.Signature, error) ***REMOVED***
	req := ssh.Marshal(signRequestAgentMsg***REMOVED***
		KeyBlob: key.Marshal(),
		Data:    data,
	***REMOVED***)

	msg, err := c.call(req)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	switch msg := msg.(type) ***REMOVED***
	case *signResponseAgentMsg:
		var sig ssh.Signature
		if err := ssh.Unmarshal(msg.SigBlob, &sig); err != nil ***REMOVED***
			return nil, err
		***REMOVED***

		return &sig, nil
	case *failureAgentMsg:
		return nil, errors.New("agent: failed to sign challenge")
	***REMOVED***
	panic("unreachable")
***REMOVED***

// unmarshal parses an agent message in packet, returning the parsed
// form and the message type of packet.
func unmarshal(packet []byte) (interface***REMOVED******REMOVED***, error) ***REMOVED***
	if len(packet) < 1 ***REMOVED***
		return nil, errors.New("agent: empty packet")
	***REMOVED***
	var msg interface***REMOVED******REMOVED***
	switch packet[0] ***REMOVED***
	case agentFailure:
		return new(failureAgentMsg), nil
	case agentSuccess:
		return new(successAgentMsg), nil
	case agentIdentitiesAnswer:
		msg = new(identitiesAnswerAgentMsg)
	case agentSignResponse:
		msg = new(signResponseAgentMsg)
	case agentV1IdentitiesAnswer:
		msg = new(agentV1IdentityMsg)
	default:
		return nil, fmt.Errorf("agent: unknown type tag %d", packet[0])
	***REMOVED***
	if err := ssh.Unmarshal(packet, msg); err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	return msg, nil
***REMOVED***

type rsaKeyMsg struct ***REMOVED***
	Type        string `sshtype:"17|25"`
	N           *big.Int
	E           *big.Int
	D           *big.Int
	Iqmp        *big.Int // IQMP = Inverse Q Mod P
	P           *big.Int
	Q           *big.Int
	Comments    string
	Constraints []byte `ssh:"rest"`
***REMOVED***

type dsaKeyMsg struct ***REMOVED***
	Type        string `sshtype:"17|25"`
	P           *big.Int
	Q           *big.Int
	G           *big.Int
	Y           *big.Int
	X           *big.Int
	Comments    string
	Constraints []byte `ssh:"rest"`
***REMOVED***

type ecdsaKeyMsg struct ***REMOVED***
	Type        string `sshtype:"17|25"`
	Curve       string
	KeyBytes    []byte
	D           *big.Int
	Comments    string
	Constraints []byte `ssh:"rest"`
***REMOVED***

type ed25519KeyMsg struct ***REMOVED***
	Type        string `sshtype:"17|25"`
	Pub         []byte
	Priv        []byte
	Comments    string
	Constraints []byte `ssh:"rest"`
***REMOVED***

// Insert adds a private key to the agent.
func (c *client) insertKey(s interface***REMOVED******REMOVED***, comment string, constraints []byte) error ***REMOVED***
	var req []byte
	switch k := s.(type) ***REMOVED***
	case *rsa.PrivateKey:
		if len(k.Primes) != 2 ***REMOVED***
			return fmt.Errorf("agent: unsupported RSA key with %d primes", len(k.Primes))
		***REMOVED***
		k.Precompute()
		req = ssh.Marshal(rsaKeyMsg***REMOVED***
			Type:        ssh.KeyAlgoRSA,
			N:           k.N,
			E:           big.NewInt(int64(k.E)),
			D:           k.D,
			Iqmp:        k.Precomputed.Qinv,
			P:           k.Primes[0],
			Q:           k.Primes[1],
			Comments:    comment,
			Constraints: constraints,
		***REMOVED***)
	case *dsa.PrivateKey:
		req = ssh.Marshal(dsaKeyMsg***REMOVED***
			Type:        ssh.KeyAlgoDSA,
			P:           k.P,
			Q:           k.Q,
			G:           k.G,
			Y:           k.Y,
			X:           k.X,
			Comments:    comment,
			Constraints: constraints,
		***REMOVED***)
	case *ecdsa.PrivateKey:
		nistID := fmt.Sprintf("nistp%d", k.Params().BitSize)
		req = ssh.Marshal(ecdsaKeyMsg***REMOVED***
			Type:        "ecdsa-sha2-" + nistID,
			Curve:       nistID,
			KeyBytes:    elliptic.Marshal(k.Curve, k.X, k.Y),
			D:           k.D,
			Comments:    comment,
			Constraints: constraints,
		***REMOVED***)
	case *ed25519.PrivateKey:
		req = ssh.Marshal(ed25519KeyMsg***REMOVED***
			Type:        ssh.KeyAlgoED25519,
			Pub:         []byte(*k)[32:],
			Priv:        []byte(*k),
			Comments:    comment,
			Constraints: constraints,
		***REMOVED***)
	default:
		return fmt.Errorf("agent: unsupported key type %T", s)
	***REMOVED***

	// if constraints are present then the message type needs to be changed.
	if len(constraints) != 0 ***REMOVED***
		req[0] = agentAddIDConstrained
	***REMOVED***

	resp, err := c.call(req)
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	if _, ok := resp.(*successAgentMsg); ok ***REMOVED***
		return nil
	***REMOVED***
	return errors.New("agent: failure")
***REMOVED***

type rsaCertMsg struct ***REMOVED***
	Type        string `sshtype:"17|25"`
	CertBytes   []byte
	D           *big.Int
	Iqmp        *big.Int // IQMP = Inverse Q Mod P
	P           *big.Int
	Q           *big.Int
	Comments    string
	Constraints []byte `ssh:"rest"`
***REMOVED***

type dsaCertMsg struct ***REMOVED***
	Type        string `sshtype:"17|25"`
	CertBytes   []byte
	X           *big.Int
	Comments    string
	Constraints []byte `ssh:"rest"`
***REMOVED***

type ecdsaCertMsg struct ***REMOVED***
	Type        string `sshtype:"17|25"`
	CertBytes   []byte
	D           *big.Int
	Comments    string
	Constraints []byte `ssh:"rest"`
***REMOVED***

type ed25519CertMsg struct ***REMOVED***
	Type        string `sshtype:"17|25"`
	CertBytes   []byte
	Pub         []byte
	Priv        []byte
	Comments    string
	Constraints []byte `ssh:"rest"`
***REMOVED***

// Add adds a private key to the agent. If a certificate is given,
// that certificate is added instead as public key.
func (c *client) Add(key AddedKey) error ***REMOVED***
	var constraints []byte

	if secs := key.LifetimeSecs; secs != 0 ***REMOVED***
		constraints = append(constraints, ssh.Marshal(constrainLifetimeAgentMsg***REMOVED***secs***REMOVED***)...)
	***REMOVED***

	if key.ConfirmBeforeUse ***REMOVED***
		constraints = append(constraints, agentConstrainConfirm)
	***REMOVED***

	cert := key.Certificate
	if cert == nil ***REMOVED***
		return c.insertKey(key.PrivateKey, key.Comment, constraints)
	***REMOVED***
	return c.insertCert(key.PrivateKey, cert, key.Comment, constraints)
***REMOVED***

func (c *client) insertCert(s interface***REMOVED******REMOVED***, cert *ssh.Certificate, comment string, constraints []byte) error ***REMOVED***
	var req []byte
	switch k := s.(type) ***REMOVED***
	case *rsa.PrivateKey:
		if len(k.Primes) != 2 ***REMOVED***
			return fmt.Errorf("agent: unsupported RSA key with %d primes", len(k.Primes))
		***REMOVED***
		k.Precompute()
		req = ssh.Marshal(rsaCertMsg***REMOVED***
			Type:        cert.Type(),
			CertBytes:   cert.Marshal(),
			D:           k.D,
			Iqmp:        k.Precomputed.Qinv,
			P:           k.Primes[0],
			Q:           k.Primes[1],
			Comments:    comment,
			Constraints: constraints,
		***REMOVED***)
	case *dsa.PrivateKey:
		req = ssh.Marshal(dsaCertMsg***REMOVED***
			Type:        cert.Type(),
			CertBytes:   cert.Marshal(),
			X:           k.X,
			Comments:    comment,
			Constraints: constraints,
		***REMOVED***)
	case *ecdsa.PrivateKey:
		req = ssh.Marshal(ecdsaCertMsg***REMOVED***
			Type:        cert.Type(),
			CertBytes:   cert.Marshal(),
			D:           k.D,
			Comments:    comment,
			Constraints: constraints,
		***REMOVED***)
	case *ed25519.PrivateKey:
		req = ssh.Marshal(ed25519CertMsg***REMOVED***
			Type:        cert.Type(),
			CertBytes:   cert.Marshal(),
			Pub:         []byte(*k)[32:],
			Priv:        []byte(*k),
			Comments:    comment,
			Constraints: constraints,
		***REMOVED***)
	default:
		return fmt.Errorf("agent: unsupported key type %T", s)
	***REMOVED***

	// if constraints are present then the message type needs to be changed.
	if len(constraints) != 0 ***REMOVED***
		req[0] = agentAddIDConstrained
	***REMOVED***

	signer, err := ssh.NewSignerFromKey(s)
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	if bytes.Compare(cert.Key.Marshal(), signer.PublicKey().Marshal()) != 0 ***REMOVED***
		return errors.New("agent: signer and cert have different public key")
	***REMOVED***

	resp, err := c.call(req)
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	if _, ok := resp.(*successAgentMsg); ok ***REMOVED***
		return nil
	***REMOVED***
	return errors.New("agent: failure")
***REMOVED***

// Signers provides a callback for client authentication.
func (c *client) Signers() ([]ssh.Signer, error) ***REMOVED***
	keys, err := c.List()
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	var result []ssh.Signer
	for _, k := range keys ***REMOVED***
		result = append(result, &agentKeyringSigner***REMOVED***c, k***REMOVED***)
	***REMOVED***
	return result, nil
***REMOVED***

type agentKeyringSigner struct ***REMOVED***
	agent *client
	pub   ssh.PublicKey
***REMOVED***

func (s *agentKeyringSigner) PublicKey() ssh.PublicKey ***REMOVED***
	return s.pub
***REMOVED***

func (s *agentKeyringSigner) Sign(rand io.Reader, data []byte) (*ssh.Signature, error) ***REMOVED***
	// The agent has its own entropy source, so the rand argument is ignored.
	return s.agent.Sign(s.pub, data)
***REMOVED***
