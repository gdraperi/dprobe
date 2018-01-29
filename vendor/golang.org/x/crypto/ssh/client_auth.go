// Copyright 2011 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package ssh

import (
	"bytes"
	"errors"
	"fmt"
	"io"
)

// clientAuthenticate authenticates with the remote server. See RFC 4252.
func (c *connection) clientAuthenticate(config *ClientConfig) error ***REMOVED***
	// initiate user auth session
	if err := c.transport.writePacket(Marshal(&serviceRequestMsg***REMOVED***serviceUserAuth***REMOVED***)); err != nil ***REMOVED***
		return err
	***REMOVED***
	packet, err := c.transport.readPacket()
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	var serviceAccept serviceAcceptMsg
	if err := Unmarshal(packet, &serviceAccept); err != nil ***REMOVED***
		return err
	***REMOVED***

	// during the authentication phase the client first attempts the "none" method
	// then any untried methods suggested by the server.
	tried := make(map[string]bool)
	var lastMethods []string

	sessionID := c.transport.getSessionID()
	for auth := AuthMethod(new(noneAuth)); auth != nil; ***REMOVED***
		ok, methods, err := auth.auth(sessionID, config.User, c.transport, config.Rand)
		if err != nil ***REMOVED***
			return err
		***REMOVED***
		if ok ***REMOVED***
			// success
			return nil
		***REMOVED***
		tried[auth.method()] = true
		if methods == nil ***REMOVED***
			methods = lastMethods
		***REMOVED***
		lastMethods = methods

		auth = nil

	findNext:
		for _, a := range config.Auth ***REMOVED***
			candidateMethod := a.method()
			if tried[candidateMethod] ***REMOVED***
				continue
			***REMOVED***
			for _, meth := range methods ***REMOVED***
				if meth == candidateMethod ***REMOVED***
					auth = a
					break findNext
				***REMOVED***
			***REMOVED***
		***REMOVED***
	***REMOVED***
	return fmt.Errorf("ssh: unable to authenticate, attempted methods %v, no supported methods remain", keys(tried))
***REMOVED***

func keys(m map[string]bool) []string ***REMOVED***
	s := make([]string, 0, len(m))

	for key := range m ***REMOVED***
		s = append(s, key)
	***REMOVED***
	return s
***REMOVED***

// An AuthMethod represents an instance of an RFC 4252 authentication method.
type AuthMethod interface ***REMOVED***
	// auth authenticates user over transport t.
	// Returns true if authentication is successful.
	// If authentication is not successful, a []string of alternative
	// method names is returned. If the slice is nil, it will be ignored
	// and the previous set of possible methods will be reused.
	auth(session []byte, user string, p packetConn, rand io.Reader) (bool, []string, error)

	// method returns the RFC 4252 method name.
	method() string
***REMOVED***

// "none" authentication, RFC 4252 section 5.2.
type noneAuth int

func (n *noneAuth) auth(session []byte, user string, c packetConn, rand io.Reader) (bool, []string, error) ***REMOVED***
	if err := c.writePacket(Marshal(&userAuthRequestMsg***REMOVED***
		User:    user,
		Service: serviceSSH,
		Method:  "none",
	***REMOVED***)); err != nil ***REMOVED***
		return false, nil, err
	***REMOVED***

	return handleAuthResponse(c)
***REMOVED***

func (n *noneAuth) method() string ***REMOVED***
	return "none"
***REMOVED***

// passwordCallback is an AuthMethod that fetches the password through
// a function call, e.g. by prompting the user.
type passwordCallback func() (password string, err error)

func (cb passwordCallback) auth(session []byte, user string, c packetConn, rand io.Reader) (bool, []string, error) ***REMOVED***
	type passwordAuthMsg struct ***REMOVED***
		User     string `sshtype:"50"`
		Service  string
		Method   string
		Reply    bool
		Password string
	***REMOVED***

	pw, err := cb()
	// REVIEW NOTE: is there a need to support skipping a password attempt?
	// The program may only find out that the user doesn't have a password
	// when prompting.
	if err != nil ***REMOVED***
		return false, nil, err
	***REMOVED***

	if err := c.writePacket(Marshal(&passwordAuthMsg***REMOVED***
		User:     user,
		Service:  serviceSSH,
		Method:   cb.method(),
		Reply:    false,
		Password: pw,
	***REMOVED***)); err != nil ***REMOVED***
		return false, nil, err
	***REMOVED***

	return handleAuthResponse(c)
***REMOVED***

func (cb passwordCallback) method() string ***REMOVED***
	return "password"
***REMOVED***

// Password returns an AuthMethod using the given password.
func Password(secret string) AuthMethod ***REMOVED***
	return passwordCallback(func() (string, error) ***REMOVED*** return secret, nil ***REMOVED***)
***REMOVED***

// PasswordCallback returns an AuthMethod that uses a callback for
// fetching a password.
func PasswordCallback(prompt func() (secret string, err error)) AuthMethod ***REMOVED***
	return passwordCallback(prompt)
***REMOVED***

type publickeyAuthMsg struct ***REMOVED***
	User    string `sshtype:"50"`
	Service string
	Method  string
	// HasSig indicates to the receiver packet that the auth request is signed and
	// should be used for authentication of the request.
	HasSig   bool
	Algoname string
	PubKey   []byte
	// Sig is tagged with "rest" so Marshal will exclude it during
	// validateKey
	Sig []byte `ssh:"rest"`
***REMOVED***

// publicKeyCallback is an AuthMethod that uses a set of key
// pairs for authentication.
type publicKeyCallback func() ([]Signer, error)

func (cb publicKeyCallback) method() string ***REMOVED***
	return "publickey"
***REMOVED***

func (cb publicKeyCallback) auth(session []byte, user string, c packetConn, rand io.Reader) (bool, []string, error) ***REMOVED***
	// Authentication is performed by sending an enquiry to test if a key is
	// acceptable to the remote. If the key is acceptable, the client will
	// attempt to authenticate with the valid key.  If not the client will repeat
	// the process with the remaining keys.

	signers, err := cb()
	if err != nil ***REMOVED***
		return false, nil, err
	***REMOVED***
	var methods []string
	for _, signer := range signers ***REMOVED***
		ok, err := validateKey(signer.PublicKey(), user, c)
		if err != nil ***REMOVED***
			return false, nil, err
		***REMOVED***
		if !ok ***REMOVED***
			continue
		***REMOVED***

		pub := signer.PublicKey()
		pubKey := pub.Marshal()
		sign, err := signer.Sign(rand, buildDataSignedForAuth(session, userAuthRequestMsg***REMOVED***
			User:    user,
			Service: serviceSSH,
			Method:  cb.method(),
		***REMOVED***, []byte(pub.Type()), pubKey))
		if err != nil ***REMOVED***
			return false, nil, err
		***REMOVED***

		// manually wrap the serialized signature in a string
		s := Marshal(sign)
		sig := make([]byte, stringLength(len(s)))
		marshalString(sig, s)
		msg := publickeyAuthMsg***REMOVED***
			User:     user,
			Service:  serviceSSH,
			Method:   cb.method(),
			HasSig:   true,
			Algoname: pub.Type(),
			PubKey:   pubKey,
			Sig:      sig,
		***REMOVED***
		p := Marshal(&msg)
		if err := c.writePacket(p); err != nil ***REMOVED***
			return false, nil, err
		***REMOVED***
		var success bool
		success, methods, err = handleAuthResponse(c)
		if err != nil ***REMOVED***
			return false, nil, err
		***REMOVED***

		// If authentication succeeds or the list of available methods does not
		// contain the "publickey" method, do not attempt to authenticate with any
		// other keys.  According to RFC 4252 Section 7, the latter can occur when
		// additional authentication methods are required.
		if success || !containsMethod(methods, cb.method()) ***REMOVED***
			return success, methods, err
		***REMOVED***
	***REMOVED***

	return false, methods, nil
***REMOVED***

func containsMethod(methods []string, method string) bool ***REMOVED***
	for _, m := range methods ***REMOVED***
		if m == method ***REMOVED***
			return true
		***REMOVED***
	***REMOVED***

	return false
***REMOVED***

// validateKey validates the key provided is acceptable to the server.
func validateKey(key PublicKey, user string, c packetConn) (bool, error) ***REMOVED***
	pubKey := key.Marshal()
	msg := publickeyAuthMsg***REMOVED***
		User:     user,
		Service:  serviceSSH,
		Method:   "publickey",
		HasSig:   false,
		Algoname: key.Type(),
		PubKey:   pubKey,
	***REMOVED***
	if err := c.writePacket(Marshal(&msg)); err != nil ***REMOVED***
		return false, err
	***REMOVED***

	return confirmKeyAck(key, c)
***REMOVED***

func confirmKeyAck(key PublicKey, c packetConn) (bool, error) ***REMOVED***
	pubKey := key.Marshal()
	algoname := key.Type()

	for ***REMOVED***
		packet, err := c.readPacket()
		if err != nil ***REMOVED***
			return false, err
		***REMOVED***
		switch packet[0] ***REMOVED***
		case msgUserAuthBanner:
			if err := handleBannerResponse(c, packet); err != nil ***REMOVED***
				return false, err
			***REMOVED***
		case msgUserAuthPubKeyOk:
			var msg userAuthPubKeyOkMsg
			if err := Unmarshal(packet, &msg); err != nil ***REMOVED***
				return false, err
			***REMOVED***
			if msg.Algo != algoname || !bytes.Equal(msg.PubKey, pubKey) ***REMOVED***
				return false, nil
			***REMOVED***
			return true, nil
		case msgUserAuthFailure:
			return false, nil
		default:
			return false, unexpectedMessageError(msgUserAuthSuccess, packet[0])
		***REMOVED***
	***REMOVED***
***REMOVED***

// PublicKeys returns an AuthMethod that uses the given key
// pairs.
func PublicKeys(signers ...Signer) AuthMethod ***REMOVED***
	return publicKeyCallback(func() ([]Signer, error) ***REMOVED*** return signers, nil ***REMOVED***)
***REMOVED***

// PublicKeysCallback returns an AuthMethod that runs the given
// function to obtain a list of key pairs.
func PublicKeysCallback(getSigners func() (signers []Signer, err error)) AuthMethod ***REMOVED***
	return publicKeyCallback(getSigners)
***REMOVED***

// handleAuthResponse returns whether the preceding authentication request succeeded
// along with a list of remaining authentication methods to try next and
// an error if an unexpected response was received.
func handleAuthResponse(c packetConn) (bool, []string, error) ***REMOVED***
	for ***REMOVED***
		packet, err := c.readPacket()
		if err != nil ***REMOVED***
			return false, nil, err
		***REMOVED***

		switch packet[0] ***REMOVED***
		case msgUserAuthBanner:
			if err := handleBannerResponse(c, packet); err != nil ***REMOVED***
				return false, nil, err
			***REMOVED***
		case msgUserAuthFailure:
			var msg userAuthFailureMsg
			if err := Unmarshal(packet, &msg); err != nil ***REMOVED***
				return false, nil, err
			***REMOVED***
			return false, msg.Methods, nil
		case msgUserAuthSuccess:
			return true, nil, nil
		default:
			return false, nil, unexpectedMessageError(msgUserAuthSuccess, packet[0])
		***REMOVED***
	***REMOVED***
***REMOVED***

func handleBannerResponse(c packetConn, packet []byte) error ***REMOVED***
	var msg userAuthBannerMsg
	if err := Unmarshal(packet, &msg); err != nil ***REMOVED***
		return err
	***REMOVED***

	transport, ok := c.(*handshakeTransport)
	if !ok ***REMOVED***
		return nil
	***REMOVED***

	if transport.bannerCallback != nil ***REMOVED***
		return transport.bannerCallback(msg.Message)
	***REMOVED***

	return nil
***REMOVED***

// KeyboardInteractiveChallenge should print questions, optionally
// disabling echoing (e.g. for passwords), and return all the answers.
// Challenge may be called multiple times in a single session. After
// successful authentication, the server may send a challenge with no
// questions, for which the user and instruction messages should be
// printed.  RFC 4256 section 3.3 details how the UI should behave for
// both CLI and GUI environments.
type KeyboardInteractiveChallenge func(user, instruction string, questions []string, echos []bool) (answers []string, err error)

// KeyboardInteractive returns an AuthMethod using a prompt/response
// sequence controlled by the server.
func KeyboardInteractive(challenge KeyboardInteractiveChallenge) AuthMethod ***REMOVED***
	return challenge
***REMOVED***

func (cb KeyboardInteractiveChallenge) method() string ***REMOVED***
	return "keyboard-interactive"
***REMOVED***

func (cb KeyboardInteractiveChallenge) auth(session []byte, user string, c packetConn, rand io.Reader) (bool, []string, error) ***REMOVED***
	type initiateMsg struct ***REMOVED***
		User       string `sshtype:"50"`
		Service    string
		Method     string
		Language   string
		Submethods string
	***REMOVED***

	if err := c.writePacket(Marshal(&initiateMsg***REMOVED***
		User:    user,
		Service: serviceSSH,
		Method:  "keyboard-interactive",
	***REMOVED***)); err != nil ***REMOVED***
		return false, nil, err
	***REMOVED***

	for ***REMOVED***
		packet, err := c.readPacket()
		if err != nil ***REMOVED***
			return false, nil, err
		***REMOVED***

		// like handleAuthResponse, but with less options.
		switch packet[0] ***REMOVED***
		case msgUserAuthBanner:
			if err := handleBannerResponse(c, packet); err != nil ***REMOVED***
				return false, nil, err
			***REMOVED***
			continue
		case msgUserAuthInfoRequest:
			// OK
		case msgUserAuthFailure:
			var msg userAuthFailureMsg
			if err := Unmarshal(packet, &msg); err != nil ***REMOVED***
				return false, nil, err
			***REMOVED***
			return false, msg.Methods, nil
		case msgUserAuthSuccess:
			return true, nil, nil
		default:
			return false, nil, unexpectedMessageError(msgUserAuthInfoRequest, packet[0])
		***REMOVED***

		var msg userAuthInfoRequestMsg
		if err := Unmarshal(packet, &msg); err != nil ***REMOVED***
			return false, nil, err
		***REMOVED***

		// Manually unpack the prompt/echo pairs.
		rest := msg.Prompts
		var prompts []string
		var echos []bool
		for i := 0; i < int(msg.NumPrompts); i++ ***REMOVED***
			prompt, r, ok := parseString(rest)
			if !ok || len(r) == 0 ***REMOVED***
				return false, nil, errors.New("ssh: prompt format error")
			***REMOVED***
			prompts = append(prompts, string(prompt))
			echos = append(echos, r[0] != 0)
			rest = r[1:]
		***REMOVED***

		if len(rest) != 0 ***REMOVED***
			return false, nil, errors.New("ssh: extra data following keyboard-interactive pairs")
		***REMOVED***

		answers, err := cb(msg.User, msg.Instruction, prompts, echos)
		if err != nil ***REMOVED***
			return false, nil, err
		***REMOVED***

		if len(answers) != len(prompts) ***REMOVED***
			return false, nil, errors.New("ssh: not enough answers from keyboard-interactive callback")
		***REMOVED***
		responseLength := 1 + 4
		for _, a := range answers ***REMOVED***
			responseLength += stringLength(len(a))
		***REMOVED***
		serialized := make([]byte, responseLength)
		p := serialized
		p[0] = msgUserAuthInfoResponse
		p = p[1:]
		p = marshalUint32(p, uint32(len(answers)))
		for _, a := range answers ***REMOVED***
			p = marshalString(p, []byte(a))
		***REMOVED***

		if err := c.writePacket(serialized); err != nil ***REMOVED***
			return false, nil, err
		***REMOVED***
	***REMOVED***
***REMOVED***

type retryableAuthMethod struct ***REMOVED***
	authMethod AuthMethod
	maxTries   int
***REMOVED***

func (r *retryableAuthMethod) auth(session []byte, user string, c packetConn, rand io.Reader) (ok bool, methods []string, err error) ***REMOVED***
	for i := 0; r.maxTries <= 0 || i < r.maxTries; i++ ***REMOVED***
		ok, methods, err = r.authMethod.auth(session, user, c, rand)
		if ok || err != nil ***REMOVED*** // either success or error terminate
			return ok, methods, err
		***REMOVED***
	***REMOVED***
	return ok, methods, err
***REMOVED***

func (r *retryableAuthMethod) method() string ***REMOVED***
	return r.authMethod.method()
***REMOVED***

// RetryableAuthMethod is a decorator for other auth methods enabling them to
// be retried up to maxTries before considering that AuthMethod itself failed.
// If maxTries is <= 0, will retry indefinitely
//
// This is useful for interactive clients using challenge/response type
// authentication (e.g. Keyboard-Interactive, Password, etc) where the user
// could mistype their response resulting in the server issuing a
// SSH_MSG_USERAUTH_FAILURE (rfc4252 #8 [password] and rfc4256 #3.4
// [keyboard-interactive]); Without this decorator, the non-retryable
// AuthMethod would be removed from future consideration, and never tried again
// (and so the user would never be able to retry their entry).
func RetryableAuthMethod(auth AuthMethod, maxTries int) AuthMethod ***REMOVED***
	return &retryableAuthMethod***REMOVED***authMethod: auth, maxTries: maxTries***REMOVED***
***REMOVED***
