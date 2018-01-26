// Copyright 2014 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package jws provides a partial implementation
// of JSON Web Signature encoding and decoding.
// It exists to support the golang.org/x/oauth2 package.
//
// See RFC 7515.
//
// Deprecated: this package is not intended for public use and might be
// removed in the future. It exists for internal use only.
// Please switch to another JWS package or copy this package into your own
// source tree.
package jws // import "golang.org/x/oauth2/jws"

import (
	"bytes"
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"
)

// ClaimSet contains information about the JWT signature including the
// permissions being requested (scopes), the target of the token, the issuer,
// the time the token was issued, and the lifetime of the token.
type ClaimSet struct ***REMOVED***
	Iss   string `json:"iss"`             // email address of the client_id of the application making the access token request
	Scope string `json:"scope,omitempty"` // space-delimited list of the permissions the application requests
	Aud   string `json:"aud"`             // descriptor of the intended target of the assertion (Optional).
	Exp   int64  `json:"exp"`             // the expiration time of the assertion (seconds since Unix epoch)
	Iat   int64  `json:"iat"`             // the time the assertion was issued (seconds since Unix epoch)
	Typ   string `json:"typ,omitempty"`   // token type (Optional).

	// Email for which the application is requesting delegated access (Optional).
	Sub string `json:"sub,omitempty"`

	// The old name of Sub. Client keeps setting Prn to be
	// complaint with legacy OAuth 2.0 providers. (Optional)
	Prn string `json:"prn,omitempty"`

	// See http://tools.ietf.org/html/draft-jones-json-web-token-10#section-4.3
	// This array is marshalled using custom code (see (c *ClaimSet) encode()).
	PrivateClaims map[string]interface***REMOVED******REMOVED*** `json:"-"`
***REMOVED***

func (c *ClaimSet) encode() (string, error) ***REMOVED***
	// Reverting time back for machines whose time is not perfectly in sync.
	// If client machine's time is in the future according
	// to Google servers, an access token will not be issued.
	now := time.Now().Add(-10 * time.Second)
	if c.Iat == 0 ***REMOVED***
		c.Iat = now.Unix()
	***REMOVED***
	if c.Exp == 0 ***REMOVED***
		c.Exp = now.Add(time.Hour).Unix()
	***REMOVED***
	if c.Exp < c.Iat ***REMOVED***
		return "", fmt.Errorf("jws: invalid Exp = %v; must be later than Iat = %v", c.Exp, c.Iat)
	***REMOVED***

	b, err := json.Marshal(c)
	if err != nil ***REMOVED***
		return "", err
	***REMOVED***

	if len(c.PrivateClaims) == 0 ***REMOVED***
		return base64.RawURLEncoding.EncodeToString(b), nil
	***REMOVED***

	// Marshal private claim set and then append it to b.
	prv, err := json.Marshal(c.PrivateClaims)
	if err != nil ***REMOVED***
		return "", fmt.Errorf("jws: invalid map of private claims %v", c.PrivateClaims)
	***REMOVED***

	// Concatenate public and private claim JSON objects.
	if !bytes.HasSuffix(b, []byte***REMOVED***'***REMOVED***'***REMOVED***) ***REMOVED***
		return "", fmt.Errorf("jws: invalid JSON %s", b)
	***REMOVED***
	if !bytes.HasPrefix(prv, []byte***REMOVED***'***REMOVED***'***REMOVED***) ***REMOVED***
		return "", fmt.Errorf("jws: invalid JSON %s", prv)
	***REMOVED***
	b[len(b)-1] = ','         // Replace closing curly brace with a comma.
	b = append(b, prv[1:]...) // Append private claims.
	return base64.RawURLEncoding.EncodeToString(b), nil
***REMOVED***

// Header represents the header for the signed JWS payloads.
type Header struct ***REMOVED***
	// The algorithm used for signature.
	Algorithm string `json:"alg"`

	// Represents the token type.
	Typ string `json:"typ"`

	// The optional hint of which key is being used.
	KeyID string `json:"kid,omitempty"`
***REMOVED***

func (h *Header) encode() (string, error) ***REMOVED***
	b, err := json.Marshal(h)
	if err != nil ***REMOVED***
		return "", err
	***REMOVED***
	return base64.RawURLEncoding.EncodeToString(b), nil
***REMOVED***

// Decode decodes a claim set from a JWS payload.
func Decode(payload string) (*ClaimSet, error) ***REMOVED***
	// decode returned id token to get expiry
	s := strings.Split(payload, ".")
	if len(s) < 2 ***REMOVED***
		// TODO(jbd): Provide more context about the error.
		return nil, errors.New("jws: invalid token received")
	***REMOVED***
	decoded, err := base64.RawURLEncoding.DecodeString(s[1])
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	c := &ClaimSet***REMOVED******REMOVED***
	err = json.NewDecoder(bytes.NewBuffer(decoded)).Decode(c)
	return c, err
***REMOVED***

// Signer returns a signature for the given data.
type Signer func(data []byte) (sig []byte, err error)

// EncodeWithSigner encodes a header and claim set with the provided signer.
func EncodeWithSigner(header *Header, c *ClaimSet, sg Signer) (string, error) ***REMOVED***
	head, err := header.encode()
	if err != nil ***REMOVED***
		return "", err
	***REMOVED***
	cs, err := c.encode()
	if err != nil ***REMOVED***
		return "", err
	***REMOVED***
	ss := fmt.Sprintf("%s.%s", head, cs)
	sig, err := sg([]byte(ss))
	if err != nil ***REMOVED***
		return "", err
	***REMOVED***
	return fmt.Sprintf("%s.%s", ss, base64.RawURLEncoding.EncodeToString(sig)), nil
***REMOVED***

// Encode encodes a signed JWS with provided header and claim set.
// This invokes EncodeWithSigner using crypto/rsa.SignPKCS1v15 with the given RSA private key.
func Encode(header *Header, c *ClaimSet, key *rsa.PrivateKey) (string, error) ***REMOVED***
	sg := func(data []byte) (sig []byte, err error) ***REMOVED***
		h := sha256.New()
		h.Write(data)
		return rsa.SignPKCS1v15(rand.Reader, key, crypto.SHA256, h.Sum(nil))
	***REMOVED***
	return EncodeWithSigner(header, c, sg)
***REMOVED***

// Verify tests whether the provided JWT token's signature was produced by the private key
// associated with the supplied public key.
func Verify(token string, key *rsa.PublicKey) error ***REMOVED***
	parts := strings.Split(token, ".")
	if len(parts) != 3 ***REMOVED***
		return errors.New("jws: invalid token received, token must have 3 parts")
	***REMOVED***

	signedContent := parts[0] + "." + parts[1]
	signatureString, err := base64.RawURLEncoding.DecodeString(parts[2])
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	h := sha256.New()
	h.Write([]byte(signedContent))
	return rsa.VerifyPKCS1v15(key, crypto.SHA256, h.Sum(nil), []byte(signatureString))
***REMOVED***
