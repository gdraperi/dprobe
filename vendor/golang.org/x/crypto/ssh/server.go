// Copyright 2011 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package ssh

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"net"
	"strings"
)

// The Permissions type holds fine-grained permissions that are
// specific to a user or a specific authentication method for a user.
// The Permissions value for a successful authentication attempt is
// available in ServerConn, so it can be used to pass information from
// the user-authentication phase to the application layer.
type Permissions struct ***REMOVED***
	// CriticalOptions indicate restrictions to the default
	// permissions, and are typically used in conjunction with
	// user certificates. The standard for SSH certificates
	// defines "force-command" (only allow the given command to
	// execute) and "source-address" (only allow connections from
	// the given address). The SSH package currently only enforces
	// the "source-address" critical option. It is up to server
	// implementations to enforce other critical options, such as
	// "force-command", by checking them after the SSH handshake
	// is successful. In general, SSH servers should reject
	// connections that specify critical options that are unknown
	// or not supported.
	CriticalOptions map[string]string

	// Extensions are extra functionality that the server may
	// offer on authenticated connections. Lack of support for an
	// extension does not preclude authenticating a user. Common
	// extensions are "permit-agent-forwarding",
	// "permit-X11-forwarding". The Go SSH library currently does
	// not act on any extension, and it is up to server
	// implementations to honor them. Extensions can be used to
	// pass data from the authentication callbacks to the server
	// application layer.
	Extensions map[string]string
***REMOVED***

// ServerConfig holds server specific configuration data.
type ServerConfig struct ***REMOVED***
	// Config contains configuration shared between client and server.
	Config

	hostKeys []Signer

	// NoClientAuth is true if clients are allowed to connect without
	// authenticating.
	NoClientAuth bool

	// MaxAuthTries specifies the maximum number of authentication attempts
	// permitted per connection. If set to a negative number, the number of
	// attempts are unlimited. If set to zero, the number of attempts are limited
	// to 6.
	MaxAuthTries int

	// PasswordCallback, if non-nil, is called when a user
	// attempts to authenticate using a password.
	PasswordCallback func(conn ConnMetadata, password []byte) (*Permissions, error)

	// PublicKeyCallback, if non-nil, is called when a client
	// offers a public key for authentication. It must return a nil error
	// if the given public key can be used to authenticate the
	// given user. For example, see CertChecker.Authenticate. A
	// call to this function does not guarantee that the key
	// offered is in fact used to authenticate. To record any data
	// depending on the public key, store it inside a
	// Permissions.Extensions entry.
	PublicKeyCallback func(conn ConnMetadata, key PublicKey) (*Permissions, error)

	// KeyboardInteractiveCallback, if non-nil, is called when
	// keyboard-interactive authentication is selected (RFC
	// 4256). The client object's Challenge function should be
	// used to query the user. The callback may offer multiple
	// Challenge rounds. To avoid information leaks, the client
	// should be presented a challenge even if the user is
	// unknown.
	KeyboardInteractiveCallback func(conn ConnMetadata, client KeyboardInteractiveChallenge) (*Permissions, error)

	// AuthLogCallback, if non-nil, is called to log all authentication
	// attempts.
	AuthLogCallback func(conn ConnMetadata, method string, err error)

	// ServerVersion is the version identification string to announce in
	// the public handshake.
	// If empty, a reasonable default is used.
	// Note that RFC 4253 section 4.2 requires that this string start with
	// "SSH-2.0-".
	ServerVersion string

	// BannerCallback, if present, is called and the return string is sent to
	// the client after key exchange completed but before authentication.
	BannerCallback func(conn ConnMetadata) string
***REMOVED***

// AddHostKey adds a private key as a host key. If an existing host
// key exists with the same algorithm, it is overwritten. Each server
// config must have at least one host key.
func (s *ServerConfig) AddHostKey(key Signer) ***REMOVED***
	for i, k := range s.hostKeys ***REMOVED***
		if k.PublicKey().Type() == key.PublicKey().Type() ***REMOVED***
			s.hostKeys[i] = key
			return
		***REMOVED***
	***REMOVED***

	s.hostKeys = append(s.hostKeys, key)
***REMOVED***

// cachedPubKey contains the results of querying whether a public key is
// acceptable for a user.
type cachedPubKey struct ***REMOVED***
	user       string
	pubKeyData []byte
	result     error
	perms      *Permissions
***REMOVED***

const maxCachedPubKeys = 16

// pubKeyCache caches tests for public keys.  Since SSH clients
// will query whether a public key is acceptable before attempting to
// authenticate with it, we end up with duplicate queries for public
// key validity.  The cache only applies to a single ServerConn.
type pubKeyCache struct ***REMOVED***
	keys []cachedPubKey
***REMOVED***

// get returns the result for a given user/algo/key tuple.
func (c *pubKeyCache) get(user string, pubKeyData []byte) (cachedPubKey, bool) ***REMOVED***
	for _, k := range c.keys ***REMOVED***
		if k.user == user && bytes.Equal(k.pubKeyData, pubKeyData) ***REMOVED***
			return k, true
		***REMOVED***
	***REMOVED***
	return cachedPubKey***REMOVED******REMOVED***, false
***REMOVED***

// add adds the given tuple to the cache.
func (c *pubKeyCache) add(candidate cachedPubKey) ***REMOVED***
	if len(c.keys) < maxCachedPubKeys ***REMOVED***
		c.keys = append(c.keys, candidate)
	***REMOVED***
***REMOVED***

// ServerConn is an authenticated SSH connection, as seen from the
// server
type ServerConn struct ***REMOVED***
	Conn

	// If the succeeding authentication callback returned a
	// non-nil Permissions pointer, it is stored here.
	Permissions *Permissions
***REMOVED***

// NewServerConn starts a new SSH server with c as the underlying
// transport.  It starts with a handshake and, if the handshake is
// unsuccessful, it closes the connection and returns an error.  The
// Request and NewChannel channels must be serviced, or the connection
// will hang.
func NewServerConn(c net.Conn, config *ServerConfig) (*ServerConn, <-chan NewChannel, <-chan *Request, error) ***REMOVED***
	fullConf := *config
	fullConf.SetDefaults()
	if fullConf.MaxAuthTries == 0 ***REMOVED***
		fullConf.MaxAuthTries = 6
	***REMOVED***

	s := &connection***REMOVED***
		sshConn: sshConn***REMOVED***conn: c***REMOVED***,
	***REMOVED***
	perms, err := s.serverHandshake(&fullConf)
	if err != nil ***REMOVED***
		c.Close()
		return nil, nil, nil, err
	***REMOVED***
	return &ServerConn***REMOVED***s, perms***REMOVED***, s.mux.incomingChannels, s.mux.incomingRequests, nil
***REMOVED***

// signAndMarshal signs the data with the appropriate algorithm,
// and serializes the result in SSH wire format.
func signAndMarshal(k Signer, rand io.Reader, data []byte) ([]byte, error) ***REMOVED***
	sig, err := k.Sign(rand, data)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	return Marshal(sig), nil
***REMOVED***

// handshake performs key exchange and user authentication.
func (s *connection) serverHandshake(config *ServerConfig) (*Permissions, error) ***REMOVED***
	if len(config.hostKeys) == 0 ***REMOVED***
		return nil, errors.New("ssh: server has no host keys")
	***REMOVED***

	if !config.NoClientAuth && config.PasswordCallback == nil && config.PublicKeyCallback == nil && config.KeyboardInteractiveCallback == nil ***REMOVED***
		return nil, errors.New("ssh: no authentication methods configured but NoClientAuth is also false")
	***REMOVED***

	if config.ServerVersion != "" ***REMOVED***
		s.serverVersion = []byte(config.ServerVersion)
	***REMOVED*** else ***REMOVED***
		s.serverVersion = []byte(packageVersion)
	***REMOVED***
	var err error
	s.clientVersion, err = exchangeVersions(s.sshConn.conn, s.serverVersion)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	tr := newTransport(s.sshConn.conn, config.Rand, false /* not client */)
	s.transport = newServerTransport(tr, s.clientVersion, s.serverVersion, config)

	if err := s.transport.waitSession(); err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	// We just did the key change, so the session ID is established.
	s.sessionID = s.transport.getSessionID()

	var packet []byte
	if packet, err = s.transport.readPacket(); err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	var serviceRequest serviceRequestMsg
	if err = Unmarshal(packet, &serviceRequest); err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	if serviceRequest.Service != serviceUserAuth ***REMOVED***
		return nil, errors.New("ssh: requested service '" + serviceRequest.Service + "' before authenticating")
	***REMOVED***
	serviceAccept := serviceAcceptMsg***REMOVED***
		Service: serviceUserAuth,
	***REMOVED***
	if err := s.transport.writePacket(Marshal(&serviceAccept)); err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	perms, err := s.serverAuthenticate(config)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	s.mux = newMux(s.transport)
	return perms, err
***REMOVED***

func isAcceptableAlgo(algo string) bool ***REMOVED***
	switch algo ***REMOVED***
	case KeyAlgoRSA, KeyAlgoDSA, KeyAlgoECDSA256, KeyAlgoECDSA384, KeyAlgoECDSA521, KeyAlgoED25519,
		CertAlgoRSAv01, CertAlgoDSAv01, CertAlgoECDSA256v01, CertAlgoECDSA384v01, CertAlgoECDSA521v01, CertAlgoED25519v01:
		return true
	***REMOVED***
	return false
***REMOVED***

func checkSourceAddress(addr net.Addr, sourceAddrs string) error ***REMOVED***
	if addr == nil ***REMOVED***
		return errors.New("ssh: no address known for client, but source-address match required")
	***REMOVED***

	tcpAddr, ok := addr.(*net.TCPAddr)
	if !ok ***REMOVED***
		return fmt.Errorf("ssh: remote address %v is not an TCP address when checking source-address match", addr)
	***REMOVED***

	for _, sourceAddr := range strings.Split(sourceAddrs, ",") ***REMOVED***
		if allowedIP := net.ParseIP(sourceAddr); allowedIP != nil ***REMOVED***
			if allowedIP.Equal(tcpAddr.IP) ***REMOVED***
				return nil
			***REMOVED***
		***REMOVED*** else ***REMOVED***
			_, ipNet, err := net.ParseCIDR(sourceAddr)
			if err != nil ***REMOVED***
				return fmt.Errorf("ssh: error parsing source-address restriction %q: %v", sourceAddr, err)
			***REMOVED***

			if ipNet.Contains(tcpAddr.IP) ***REMOVED***
				return nil
			***REMOVED***
		***REMOVED***
	***REMOVED***

	return fmt.Errorf("ssh: remote address %v is not allowed because of source-address restriction", addr)
***REMOVED***

// ServerAuthError implements the error interface. It appends any authentication
// errors that may occur, and is returned if all of the authentication methods
// provided by the user failed to authenticate.
type ServerAuthError struct ***REMOVED***
	// Errors contains authentication errors returned by the authentication
	// callback methods.
	Errors []error
***REMOVED***

func (l ServerAuthError) Error() string ***REMOVED***
	var errs []string
	for _, err := range l.Errors ***REMOVED***
		errs = append(errs, err.Error())
	***REMOVED***
	return "[" + strings.Join(errs, ", ") + "]"
***REMOVED***

func (s *connection) serverAuthenticate(config *ServerConfig) (*Permissions, error) ***REMOVED***
	sessionID := s.transport.getSessionID()
	var cache pubKeyCache
	var perms *Permissions

	authFailures := 0
	var authErrs []error
	var displayedBanner bool

userAuthLoop:
	for ***REMOVED***
		if authFailures >= config.MaxAuthTries && config.MaxAuthTries > 0 ***REMOVED***
			discMsg := &disconnectMsg***REMOVED***
				Reason:  2,
				Message: "too many authentication failures",
			***REMOVED***

			if err := s.transport.writePacket(Marshal(discMsg)); err != nil ***REMOVED***
				return nil, err
			***REMOVED***

			return nil, discMsg
		***REMOVED***

		var userAuthReq userAuthRequestMsg
		if packet, err := s.transport.readPacket(); err != nil ***REMOVED***
			if err == io.EOF ***REMOVED***
				return nil, &ServerAuthError***REMOVED***Errors: authErrs***REMOVED***
			***REMOVED***
			return nil, err
		***REMOVED*** else if err = Unmarshal(packet, &userAuthReq); err != nil ***REMOVED***
			return nil, err
		***REMOVED***

		if userAuthReq.Service != serviceSSH ***REMOVED***
			return nil, errors.New("ssh: client attempted to negotiate for unknown service: " + userAuthReq.Service)
		***REMOVED***

		s.user = userAuthReq.User

		if !displayedBanner && config.BannerCallback != nil ***REMOVED***
			displayedBanner = true
			msg := config.BannerCallback(s)
			if msg != "" ***REMOVED***
				bannerMsg := &userAuthBannerMsg***REMOVED***
					Message: msg,
				***REMOVED***
				if err := s.transport.writePacket(Marshal(bannerMsg)); err != nil ***REMOVED***
					return nil, err
				***REMOVED***
			***REMOVED***
		***REMOVED***

		perms = nil
		authErr := errors.New("no auth passed yet")

		switch userAuthReq.Method ***REMOVED***
		case "none":
			if config.NoClientAuth ***REMOVED***
				authErr = nil
			***REMOVED***

			// allow initial attempt of 'none' without penalty
			if authFailures == 0 ***REMOVED***
				authFailures--
			***REMOVED***
		case "password":
			if config.PasswordCallback == nil ***REMOVED***
				authErr = errors.New("ssh: password auth not configured")
				break
			***REMOVED***
			payload := userAuthReq.Payload
			if len(payload) < 1 || payload[0] != 0 ***REMOVED***
				return nil, parseError(msgUserAuthRequest)
			***REMOVED***
			payload = payload[1:]
			password, payload, ok := parseString(payload)
			if !ok || len(payload) > 0 ***REMOVED***
				return nil, parseError(msgUserAuthRequest)
			***REMOVED***

			perms, authErr = config.PasswordCallback(s, password)
		case "keyboard-interactive":
			if config.KeyboardInteractiveCallback == nil ***REMOVED***
				authErr = errors.New("ssh: keyboard-interactive auth not configubred")
				break
			***REMOVED***

			prompter := &sshClientKeyboardInteractive***REMOVED***s***REMOVED***
			perms, authErr = config.KeyboardInteractiveCallback(s, prompter.Challenge)
		case "publickey":
			if config.PublicKeyCallback == nil ***REMOVED***
				authErr = errors.New("ssh: publickey auth not configured")
				break
			***REMOVED***
			payload := userAuthReq.Payload
			if len(payload) < 1 ***REMOVED***
				return nil, parseError(msgUserAuthRequest)
			***REMOVED***
			isQuery := payload[0] == 0
			payload = payload[1:]
			algoBytes, payload, ok := parseString(payload)
			if !ok ***REMOVED***
				return nil, parseError(msgUserAuthRequest)
			***REMOVED***
			algo := string(algoBytes)
			if !isAcceptableAlgo(algo) ***REMOVED***
				authErr = fmt.Errorf("ssh: algorithm %q not accepted", algo)
				break
			***REMOVED***

			pubKeyData, payload, ok := parseString(payload)
			if !ok ***REMOVED***
				return nil, parseError(msgUserAuthRequest)
			***REMOVED***

			pubKey, err := ParsePublicKey(pubKeyData)
			if err != nil ***REMOVED***
				return nil, err
			***REMOVED***

			candidate, ok := cache.get(s.user, pubKeyData)
			if !ok ***REMOVED***
				candidate.user = s.user
				candidate.pubKeyData = pubKeyData
				candidate.perms, candidate.result = config.PublicKeyCallback(s, pubKey)
				if candidate.result == nil && candidate.perms != nil && candidate.perms.CriticalOptions != nil && candidate.perms.CriticalOptions[sourceAddressCriticalOption] != "" ***REMOVED***
					candidate.result = checkSourceAddress(
						s.RemoteAddr(),
						candidate.perms.CriticalOptions[sourceAddressCriticalOption])
				***REMOVED***
				cache.add(candidate)
			***REMOVED***

			if isQuery ***REMOVED***
				// The client can query if the given public key
				// would be okay.

				if len(payload) > 0 ***REMOVED***
					return nil, parseError(msgUserAuthRequest)
				***REMOVED***

				if candidate.result == nil ***REMOVED***
					okMsg := userAuthPubKeyOkMsg***REMOVED***
						Algo:   algo,
						PubKey: pubKeyData,
					***REMOVED***
					if err = s.transport.writePacket(Marshal(&okMsg)); err != nil ***REMOVED***
						return nil, err
					***REMOVED***
					continue userAuthLoop
				***REMOVED***
				authErr = candidate.result
			***REMOVED*** else ***REMOVED***
				sig, payload, ok := parseSignature(payload)
				if !ok || len(payload) > 0 ***REMOVED***
					return nil, parseError(msgUserAuthRequest)
				***REMOVED***
				// Ensure the public key algo and signature algo
				// are supported.  Compare the private key
				// algorithm name that corresponds to algo with
				// sig.Format.  This is usually the same, but
				// for certs, the names differ.
				if !isAcceptableAlgo(sig.Format) ***REMOVED***
					break
				***REMOVED***
				signedData := buildDataSignedForAuth(sessionID, userAuthReq, algoBytes, pubKeyData)

				if err := pubKey.Verify(signedData, sig); err != nil ***REMOVED***
					return nil, err
				***REMOVED***

				authErr = candidate.result
				perms = candidate.perms
			***REMOVED***
		default:
			authErr = fmt.Errorf("ssh: unknown method %q", userAuthReq.Method)
		***REMOVED***

		authErrs = append(authErrs, authErr)

		if config.AuthLogCallback != nil ***REMOVED***
			config.AuthLogCallback(s, userAuthReq.Method, authErr)
		***REMOVED***

		if authErr == nil ***REMOVED***
			break userAuthLoop
		***REMOVED***

		authFailures++

		var failureMsg userAuthFailureMsg
		if config.PasswordCallback != nil ***REMOVED***
			failureMsg.Methods = append(failureMsg.Methods, "password")
		***REMOVED***
		if config.PublicKeyCallback != nil ***REMOVED***
			failureMsg.Methods = append(failureMsg.Methods, "publickey")
		***REMOVED***
		if config.KeyboardInteractiveCallback != nil ***REMOVED***
			failureMsg.Methods = append(failureMsg.Methods, "keyboard-interactive")
		***REMOVED***

		if len(failureMsg.Methods) == 0 ***REMOVED***
			return nil, errors.New("ssh: no authentication methods configured but NoClientAuth is also false")
		***REMOVED***

		if err := s.transport.writePacket(Marshal(&failureMsg)); err != nil ***REMOVED***
			return nil, err
		***REMOVED***
	***REMOVED***

	if err := s.transport.writePacket([]byte***REMOVED***msgUserAuthSuccess***REMOVED***); err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	return perms, nil
***REMOVED***

// sshClientKeyboardInteractive implements a ClientKeyboardInteractive by
// asking the client on the other side of a ServerConn.
type sshClientKeyboardInteractive struct ***REMOVED***
	*connection
***REMOVED***

func (c *sshClientKeyboardInteractive) Challenge(user, instruction string, questions []string, echos []bool) (answers []string, err error) ***REMOVED***
	if len(questions) != len(echos) ***REMOVED***
		return nil, errors.New("ssh: echos and questions must have equal length")
	***REMOVED***

	var prompts []byte
	for i := range questions ***REMOVED***
		prompts = appendString(prompts, questions[i])
		prompts = appendBool(prompts, echos[i])
	***REMOVED***

	if err := c.transport.writePacket(Marshal(&userAuthInfoRequestMsg***REMOVED***
		Instruction: instruction,
		NumPrompts:  uint32(len(questions)),
		Prompts:     prompts,
	***REMOVED***)); err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	packet, err := c.transport.readPacket()
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	if packet[0] != msgUserAuthInfoResponse ***REMOVED***
		return nil, unexpectedMessageError(msgUserAuthInfoResponse, packet[0])
	***REMOVED***
	packet = packet[1:]

	n, packet, ok := parseUint32(packet)
	if !ok || int(n) != len(questions) ***REMOVED***
		return nil, parseError(msgUserAuthInfoResponse)
	***REMOVED***

	for i := uint32(0); i < n; i++ ***REMOVED***
		ans, rest, ok := parseString(packet)
		if !ok ***REMOVED***
			return nil, parseError(msgUserAuthInfoResponse)
		***REMOVED***

		answers = append(answers, string(ans))
		packet = rest
	***REMOVED***
	if len(packet) != 0 ***REMOVED***
		return nil, errors.New("ssh: junk at end of message")
	***REMOVED***

	return answers, nil
***REMOVED***
