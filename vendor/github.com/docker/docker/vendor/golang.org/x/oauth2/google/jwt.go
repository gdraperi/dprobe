// Copyright 2015 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package google

import (
	"crypto/rsa"
	"fmt"
	"time"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/internal"
	"golang.org/x/oauth2/jws"
)

// JWTAccessTokenSourceFromJSON uses a Google Developers service account JSON
// key file to read the credentials that authorize and authenticate the
// requests, and returns a TokenSource that does not use any OAuth2 flow but
// instead creates a JWT and sends that as the access token.
// The audience is typically a URL that specifies the scope of the credentials.
//
// Note that this is not a standard OAuth flow, but rather an
// optimization supported by a few Google services.
// Unless you know otherwise, you should use JWTConfigFromJSON instead.
func JWTAccessTokenSourceFromJSON(jsonKey []byte, audience string) (oauth2.TokenSource, error) ***REMOVED***
	cfg, err := JWTConfigFromJSON(jsonKey)
	if err != nil ***REMOVED***
		return nil, fmt.Errorf("google: could not parse JSON key: %v", err)
	***REMOVED***
	pk, err := internal.ParseKey(cfg.PrivateKey)
	if err != nil ***REMOVED***
		return nil, fmt.Errorf("google: could not parse key: %v", err)
	***REMOVED***
	ts := &jwtAccessTokenSource***REMOVED***
		email:    cfg.Email,
		audience: audience,
		pk:       pk,
		pkID:     cfg.PrivateKeyID,
	***REMOVED***
	tok, err := ts.Token()
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	return oauth2.ReuseTokenSource(tok, ts), nil
***REMOVED***

type jwtAccessTokenSource struct ***REMOVED***
	email, audience string
	pk              *rsa.PrivateKey
	pkID            string
***REMOVED***

func (ts *jwtAccessTokenSource) Token() (*oauth2.Token, error) ***REMOVED***
	iat := time.Now()
	exp := iat.Add(time.Hour)
	cs := &jws.ClaimSet***REMOVED***
		Iss: ts.email,
		Sub: ts.email,
		Aud: ts.audience,
		Iat: iat.Unix(),
		Exp: exp.Unix(),
	***REMOVED***
	hdr := &jws.Header***REMOVED***
		Algorithm: "RS256",
		Typ:       "JWT",
		KeyID:     string(ts.pkID),
	***REMOVED***
	msg, err := jws.Encode(hdr, cs, ts.pk)
	if err != nil ***REMOVED***
		return nil, fmt.Errorf("google: could not encode JWT: %v", err)
	***REMOVED***
	return &oauth2.Token***REMOVED***AccessToken: msg, TokenType: "Bearer", Expiry: exp***REMOVED***, nil
***REMOVED***
