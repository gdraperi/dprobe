// Copyright 2012 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package agent

import (
	"crypto/dsa"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rsa"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"log"
	"math/big"

	"golang.org/x/crypto/ed25519"
	"golang.org/x/crypto/ssh"
)

// Server wraps an Agent and uses it to implement the agent side of
// the SSH-agent, wire protocol.
type server struct ***REMOVED***
	agent Agent
***REMOVED***

func (s *server) processRequestBytes(reqData []byte) []byte ***REMOVED***
	rep, err := s.processRequest(reqData)
	if err != nil ***REMOVED***
		if err != errLocked ***REMOVED***
			// TODO(hanwen): provide better logging interface?
			log.Printf("agent %d: %v", reqData[0], err)
		***REMOVED***
		return []byte***REMOVED***agentFailure***REMOVED***
	***REMOVED***

	if err == nil && rep == nil ***REMOVED***
		return []byte***REMOVED***agentSuccess***REMOVED***
	***REMOVED***

	return ssh.Marshal(rep)
***REMOVED***

func marshalKey(k *Key) []byte ***REMOVED***
	var record struct ***REMOVED***
		Blob    []byte
		Comment string
	***REMOVED***
	record.Blob = k.Marshal()
	record.Comment = k.Comment

	return ssh.Marshal(&record)
***REMOVED***

// See [PROTOCOL.agent], section 2.5.1.
const agentV1IdentitiesAnswer = 2

type agentV1IdentityMsg struct ***REMOVED***
	Numkeys uint32 `sshtype:"2"`
***REMOVED***

type agentRemoveIdentityMsg struct ***REMOVED***
	KeyBlob []byte `sshtype:"18"`
***REMOVED***

type agentLockMsg struct ***REMOVED***
	Passphrase []byte `sshtype:"22"`
***REMOVED***

type agentUnlockMsg struct ***REMOVED***
	Passphrase []byte `sshtype:"23"`
***REMOVED***

func (s *server) processRequest(data []byte) (interface***REMOVED******REMOVED***, error) ***REMOVED***
	switch data[0] ***REMOVED***
	case agentRequestV1Identities:
		return &agentV1IdentityMsg***REMOVED***0***REMOVED***, nil

	case agentRemoveAllV1Identities:
		return nil, nil

	case agentRemoveIdentity:
		var req agentRemoveIdentityMsg
		if err := ssh.Unmarshal(data, &req); err != nil ***REMOVED***
			return nil, err
		***REMOVED***

		var wk wireKey
		if err := ssh.Unmarshal(req.KeyBlob, &wk); err != nil ***REMOVED***
			return nil, err
		***REMOVED***

		return nil, s.agent.Remove(&Key***REMOVED***Format: wk.Format, Blob: req.KeyBlob***REMOVED***)

	case agentRemoveAllIdentities:
		return nil, s.agent.RemoveAll()

	case agentLock:
		var req agentLockMsg
		if err := ssh.Unmarshal(data, &req); err != nil ***REMOVED***
			return nil, err
		***REMOVED***

		return nil, s.agent.Lock(req.Passphrase)

	case agentUnlock:
		var req agentUnlockMsg
		if err := ssh.Unmarshal(data, &req); err != nil ***REMOVED***
			return nil, err
		***REMOVED***
		return nil, s.agent.Unlock(req.Passphrase)

	case agentSignRequest:
		var req signRequestAgentMsg
		if err := ssh.Unmarshal(data, &req); err != nil ***REMOVED***
			return nil, err
		***REMOVED***

		var wk wireKey
		if err := ssh.Unmarshal(req.KeyBlob, &wk); err != nil ***REMOVED***
			return nil, err
		***REMOVED***

		k := &Key***REMOVED***
			Format: wk.Format,
			Blob:   req.KeyBlob,
		***REMOVED***

		sig, err := s.agent.Sign(k, req.Data) //  TODO(hanwen): flags.
		if err != nil ***REMOVED***
			return nil, err
		***REMOVED***
		return &signResponseAgentMsg***REMOVED***SigBlob: ssh.Marshal(sig)***REMOVED***, nil

	case agentRequestIdentities:
		keys, err := s.agent.List()
		if err != nil ***REMOVED***
			return nil, err
		***REMOVED***

		rep := identitiesAnswerAgentMsg***REMOVED***
			NumKeys: uint32(len(keys)),
		***REMOVED***
		for _, k := range keys ***REMOVED***
			rep.Keys = append(rep.Keys, marshalKey(k)...)
		***REMOVED***
		return rep, nil

	case agentAddIDConstrained, agentAddIdentity:
		return nil, s.insertIdentity(data)
	***REMOVED***

	return nil, fmt.Errorf("unknown opcode %d", data[0])
***REMOVED***

func parseConstraints(constraints []byte) (lifetimeSecs uint32, confirmBeforeUse bool, extensions []ConstraintExtension, err error) ***REMOVED***
	for len(constraints) != 0 ***REMOVED***
		switch constraints[0] ***REMOVED***
		case agentConstrainLifetime:
			lifetimeSecs = binary.BigEndian.Uint32(constraints[1:5])
			constraints = constraints[5:]
		case agentConstrainConfirm:
			confirmBeforeUse = true
			constraints = constraints[1:]
		case agentConstrainExtension:
			var msg constrainExtensionAgentMsg
			if err = ssh.Unmarshal(constraints, &msg); err != nil ***REMOVED***
				return 0, false, nil, err
			***REMOVED***
			extensions = append(extensions, ConstraintExtension***REMOVED***
				ExtensionName:    msg.ExtensionName,
				ExtensionDetails: msg.ExtensionDetails,
			***REMOVED***)
			constraints = msg.Rest
		default:
			return 0, false, nil, fmt.Errorf("unknown constraint type: %d", constraints[0])
		***REMOVED***
	***REMOVED***
	return
***REMOVED***

func setConstraints(key *AddedKey, constraintBytes []byte) error ***REMOVED***
	lifetimeSecs, confirmBeforeUse, constraintExtensions, err := parseConstraints(constraintBytes)
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	key.LifetimeSecs = lifetimeSecs
	key.ConfirmBeforeUse = confirmBeforeUse
	key.ConstraintExtensions = constraintExtensions
	return nil
***REMOVED***

func parseRSAKey(req []byte) (*AddedKey, error) ***REMOVED***
	var k rsaKeyMsg
	if err := ssh.Unmarshal(req, &k); err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	if k.E.BitLen() > 30 ***REMOVED***
		return nil, errors.New("agent: RSA public exponent too large")
	***REMOVED***
	priv := &rsa.PrivateKey***REMOVED***
		PublicKey: rsa.PublicKey***REMOVED***
			E: int(k.E.Int64()),
			N: k.N,
		***REMOVED***,
		D:      k.D,
		Primes: []*big.Int***REMOVED***k.P, k.Q***REMOVED***,
	***REMOVED***
	priv.Precompute()

	addedKey := &AddedKey***REMOVED***PrivateKey: priv, Comment: k.Comments***REMOVED***
	if err := setConstraints(addedKey, k.Constraints); err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	return addedKey, nil
***REMOVED***

func parseEd25519Key(req []byte) (*AddedKey, error) ***REMOVED***
	var k ed25519KeyMsg
	if err := ssh.Unmarshal(req, &k); err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	priv := ed25519.PrivateKey(k.Priv)

	addedKey := &AddedKey***REMOVED***PrivateKey: &priv, Comment: k.Comments***REMOVED***
	if err := setConstraints(addedKey, k.Constraints); err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	return addedKey, nil
***REMOVED***

func parseDSAKey(req []byte) (*AddedKey, error) ***REMOVED***
	var k dsaKeyMsg
	if err := ssh.Unmarshal(req, &k); err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	priv := &dsa.PrivateKey***REMOVED***
		PublicKey: dsa.PublicKey***REMOVED***
			Parameters: dsa.Parameters***REMOVED***
				P: k.P,
				Q: k.Q,
				G: k.G,
			***REMOVED***,
			Y: k.Y,
		***REMOVED***,
		X: k.X,
	***REMOVED***

	addedKey := &AddedKey***REMOVED***PrivateKey: priv, Comment: k.Comments***REMOVED***
	if err := setConstraints(addedKey, k.Constraints); err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	return addedKey, nil
***REMOVED***

func unmarshalECDSA(curveName string, keyBytes []byte, privScalar *big.Int) (priv *ecdsa.PrivateKey, err error) ***REMOVED***
	priv = &ecdsa.PrivateKey***REMOVED***
		D: privScalar,
	***REMOVED***

	switch curveName ***REMOVED***
	case "nistp256":
		priv.Curve = elliptic.P256()
	case "nistp384":
		priv.Curve = elliptic.P384()
	case "nistp521":
		priv.Curve = elliptic.P521()
	default:
		return nil, fmt.Errorf("agent: unknown curve %q", curveName)
	***REMOVED***

	priv.X, priv.Y = elliptic.Unmarshal(priv.Curve, keyBytes)
	if priv.X == nil || priv.Y == nil ***REMOVED***
		return nil, errors.New("agent: point not on curve")
	***REMOVED***

	return priv, nil
***REMOVED***

func parseEd25519Cert(req []byte) (*AddedKey, error) ***REMOVED***
	var k ed25519CertMsg
	if err := ssh.Unmarshal(req, &k); err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	pubKey, err := ssh.ParsePublicKey(k.CertBytes)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	priv := ed25519.PrivateKey(k.Priv)
	cert, ok := pubKey.(*ssh.Certificate)
	if !ok ***REMOVED***
		return nil, errors.New("agent: bad ED25519 certificate")
	***REMOVED***

	addedKey := &AddedKey***REMOVED***PrivateKey: &priv, Certificate: cert, Comment: k.Comments***REMOVED***
	if err := setConstraints(addedKey, k.Constraints); err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	return addedKey, nil
***REMOVED***

func parseECDSAKey(req []byte) (*AddedKey, error) ***REMOVED***
	var k ecdsaKeyMsg
	if err := ssh.Unmarshal(req, &k); err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	priv, err := unmarshalECDSA(k.Curve, k.KeyBytes, k.D)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	addedKey := &AddedKey***REMOVED***PrivateKey: priv, Comment: k.Comments***REMOVED***
	if err := setConstraints(addedKey, k.Constraints); err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	return addedKey, nil
***REMOVED***

func parseRSACert(req []byte) (*AddedKey, error) ***REMOVED***
	var k rsaCertMsg
	if err := ssh.Unmarshal(req, &k); err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	pubKey, err := ssh.ParsePublicKey(k.CertBytes)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	cert, ok := pubKey.(*ssh.Certificate)
	if !ok ***REMOVED***
		return nil, errors.New("agent: bad RSA certificate")
	***REMOVED***

	// An RSA publickey as marshaled by rsaPublicKey.Marshal() in keys.go
	var rsaPub struct ***REMOVED***
		Name string
		E    *big.Int
		N    *big.Int
	***REMOVED***
	if err := ssh.Unmarshal(cert.Key.Marshal(), &rsaPub); err != nil ***REMOVED***
		return nil, fmt.Errorf("agent: Unmarshal failed to parse public key: %v", err)
	***REMOVED***

	if rsaPub.E.BitLen() > 30 ***REMOVED***
		return nil, errors.New("agent: RSA public exponent too large")
	***REMOVED***

	priv := rsa.PrivateKey***REMOVED***
		PublicKey: rsa.PublicKey***REMOVED***
			E: int(rsaPub.E.Int64()),
			N: rsaPub.N,
		***REMOVED***,
		D:      k.D,
		Primes: []*big.Int***REMOVED***k.Q, k.P***REMOVED***,
	***REMOVED***
	priv.Precompute()

	addedKey := &AddedKey***REMOVED***PrivateKey: &priv, Certificate: cert, Comment: k.Comments***REMOVED***
	if err := setConstraints(addedKey, k.Constraints); err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	return addedKey, nil
***REMOVED***

func parseDSACert(req []byte) (*AddedKey, error) ***REMOVED***
	var k dsaCertMsg
	if err := ssh.Unmarshal(req, &k); err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	pubKey, err := ssh.ParsePublicKey(k.CertBytes)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	cert, ok := pubKey.(*ssh.Certificate)
	if !ok ***REMOVED***
		return nil, errors.New("agent: bad DSA certificate")
	***REMOVED***

	// A DSA publickey as marshaled by dsaPublicKey.Marshal() in keys.go
	var w struct ***REMOVED***
		Name       string
		P, Q, G, Y *big.Int
	***REMOVED***
	if err := ssh.Unmarshal(cert.Key.Marshal(), &w); err != nil ***REMOVED***
		return nil, fmt.Errorf("agent: Unmarshal failed to parse public key: %v", err)
	***REMOVED***

	priv := &dsa.PrivateKey***REMOVED***
		PublicKey: dsa.PublicKey***REMOVED***
			Parameters: dsa.Parameters***REMOVED***
				P: w.P,
				Q: w.Q,
				G: w.G,
			***REMOVED***,
			Y: w.Y,
		***REMOVED***,
		X: k.X,
	***REMOVED***

	addedKey := &AddedKey***REMOVED***PrivateKey: priv, Certificate: cert, Comment: k.Comments***REMOVED***
	if err := setConstraints(addedKey, k.Constraints); err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	return addedKey, nil
***REMOVED***

func parseECDSACert(req []byte) (*AddedKey, error) ***REMOVED***
	var k ecdsaCertMsg
	if err := ssh.Unmarshal(req, &k); err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	pubKey, err := ssh.ParsePublicKey(k.CertBytes)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	cert, ok := pubKey.(*ssh.Certificate)
	if !ok ***REMOVED***
		return nil, errors.New("agent: bad ECDSA certificate")
	***REMOVED***

	// An ECDSA publickey as marshaled by ecdsaPublicKey.Marshal() in keys.go
	var ecdsaPub struct ***REMOVED***
		Name string
		ID   string
		Key  []byte
	***REMOVED***
	if err := ssh.Unmarshal(cert.Key.Marshal(), &ecdsaPub); err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	priv, err := unmarshalECDSA(ecdsaPub.ID, ecdsaPub.Key, k.D)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	addedKey := &AddedKey***REMOVED***PrivateKey: priv, Certificate: cert, Comment: k.Comments***REMOVED***
	if err := setConstraints(addedKey, k.Constraints); err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	return addedKey, nil
***REMOVED***

func (s *server) insertIdentity(req []byte) error ***REMOVED***
	var record struct ***REMOVED***
		Type string `sshtype:"17|25"`
		Rest []byte `ssh:"rest"`
	***REMOVED***

	if err := ssh.Unmarshal(req, &record); err != nil ***REMOVED***
		return err
	***REMOVED***

	var addedKey *AddedKey
	var err error

	switch record.Type ***REMOVED***
	case ssh.KeyAlgoRSA:
		addedKey, err = parseRSAKey(req)
	case ssh.KeyAlgoDSA:
		addedKey, err = parseDSAKey(req)
	case ssh.KeyAlgoECDSA256, ssh.KeyAlgoECDSA384, ssh.KeyAlgoECDSA521:
		addedKey, err = parseECDSAKey(req)
	case ssh.KeyAlgoED25519:
		addedKey, err = parseEd25519Key(req)
	case ssh.CertAlgoRSAv01:
		addedKey, err = parseRSACert(req)
	case ssh.CertAlgoDSAv01:
		addedKey, err = parseDSACert(req)
	case ssh.CertAlgoECDSA256v01, ssh.CertAlgoECDSA384v01, ssh.CertAlgoECDSA521v01:
		addedKey, err = parseECDSACert(req)
	case ssh.CertAlgoED25519v01:
		addedKey, err = parseEd25519Cert(req)
	default:
		return fmt.Errorf("agent: not implemented: %q", record.Type)
	***REMOVED***

	if err != nil ***REMOVED***
		return err
	***REMOVED***
	return s.agent.Add(*addedKey)
***REMOVED***

// ServeAgent serves the agent protocol on the given connection. It
// returns when an I/O error occurs.
func ServeAgent(agent Agent, c io.ReadWriter) error ***REMOVED***
	s := &server***REMOVED***agent***REMOVED***

	var length [4]byte
	for ***REMOVED***
		if _, err := io.ReadFull(c, length[:]); err != nil ***REMOVED***
			return err
		***REMOVED***
		l := binary.BigEndian.Uint32(length[:])
		if l > maxAgentResponseBytes ***REMOVED***
			// We also cap requests.
			return fmt.Errorf("agent: request too large: %d", l)
		***REMOVED***

		req := make([]byte, l)
		if _, err := io.ReadFull(c, req); err != nil ***REMOVED***
			return err
		***REMOVED***

		repData := s.processRequestBytes(req)
		if len(repData) > maxAgentResponseBytes ***REMOVED***
			return fmt.Errorf("agent: reply too large: %d bytes", len(repData))
		***REMOVED***

		binary.BigEndian.PutUint32(length[:], uint32(len(repData)))
		if _, err := c.Write(length[:]); err != nil ***REMOVED***
			return err
		***REMOVED***
		if _, err := c.Write(repData); err != nil ***REMOVED***
			return err
		***REMOVED***
	***REMOVED***
***REMOVED***
