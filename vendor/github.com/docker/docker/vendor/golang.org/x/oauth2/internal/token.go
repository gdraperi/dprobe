// Copyright 2014 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package internal contains support packages for oauth2 package.
package internal

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"mime"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"golang.org/x/net/context"
)

// Token represents the crendentials used to authorize
// the requests to access protected resources on the OAuth 2.0
// provider's backend.
//
// This type is a mirror of oauth2.Token and exists to break
// an otherwise-circular dependency. Other internal packages
// should convert this Token into an oauth2.Token before use.
type Token struct ***REMOVED***
	// AccessToken is the token that authorizes and authenticates
	// the requests.
	AccessToken string

	// TokenType is the type of token.
	// The Type method returns either this or "Bearer", the default.
	TokenType string

	// RefreshToken is a token that's used by the application
	// (as opposed to the user) to refresh the access token
	// if it expires.
	RefreshToken string

	// Expiry is the optional expiration time of the access token.
	//
	// If zero, TokenSource implementations will reuse the same
	// token forever and RefreshToken or equivalent
	// mechanisms for that TokenSource will not be used.
	Expiry time.Time

	// Raw optionally contains extra metadata from the server
	// when updating a token.
	Raw interface***REMOVED******REMOVED***
***REMOVED***

// tokenJSON is the struct representing the HTTP response from OAuth2
// providers returning a token in JSON form.
type tokenJSON struct ***REMOVED***
	AccessToken  string         `json:"access_token"`
	TokenType    string         `json:"token_type"`
	RefreshToken string         `json:"refresh_token"`
	ExpiresIn    expirationTime `json:"expires_in"` // at least PayPal returns string, while most return number
	Expires      expirationTime `json:"expires"`    // broken Facebook spelling of expires_in
***REMOVED***

func (e *tokenJSON) expiry() (t time.Time) ***REMOVED***
	if v := e.ExpiresIn; v != 0 ***REMOVED***
		return time.Now().Add(time.Duration(v) * time.Second)
	***REMOVED***
	if v := e.Expires; v != 0 ***REMOVED***
		return time.Now().Add(time.Duration(v) * time.Second)
	***REMOVED***
	return
***REMOVED***

type expirationTime int32

func (e *expirationTime) UnmarshalJSON(b []byte) error ***REMOVED***
	var n json.Number
	err := json.Unmarshal(b, &n)
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	i, err := n.Int64()
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	*e = expirationTime(i)
	return nil
***REMOVED***

var brokenAuthHeaderProviders = []string***REMOVED***
	"https://accounts.google.com/",
	"https://api.dropbox.com/",
	"https://api.dropboxapi.com/",
	"https://api.instagram.com/",
	"https://api.netatmo.net/",
	"https://api.odnoklassniki.ru/",
	"https://api.pushbullet.com/",
	"https://api.soundcloud.com/",
	"https://api.twitch.tv/",
	"https://app.box.com/",
	"https://connect.stripe.com/",
	"https://login.microsoftonline.com/",
	"https://login.salesforce.com/",
	"https://oauth.sandbox.trainingpeaks.com/",
	"https://oauth.trainingpeaks.com/",
	"https://oauth.vk.com/",
	"https://openapi.baidu.com/",
	"https://slack.com/",
	"https://test-sandbox.auth.corp.google.com",
	"https://test.salesforce.com/",
	"https://user.gini.net/",
	"https://www.douban.com/",
	"https://www.googleapis.com/",
	"https://www.linkedin.com/",
	"https://www.strava.com/oauth/",
	"https://www.wunderlist.com/oauth/",
	"https://api.patreon.com/",
	"https://sandbox.codeswholesale.com/oauth/token",
	"https://api.codeswholesale.com/oauth/token",
***REMOVED***

func RegisterBrokenAuthHeaderProvider(tokenURL string) ***REMOVED***
	brokenAuthHeaderProviders = append(brokenAuthHeaderProviders, tokenURL)
***REMOVED***

// providerAuthHeaderWorks reports whether the OAuth2 server identified by the tokenURL
// implements the OAuth2 spec correctly
// See https://code.google.com/p/goauth2/issues/detail?id=31 for background.
// In summary:
// - Reddit only accepts client secret in the Authorization header
// - Dropbox accepts either it in URL param or Auth header, but not both.
// - Google only accepts URL param (not spec compliant?), not Auth header
// - Stripe only accepts client secret in Auth header with Bearer method, not Basic
func providerAuthHeaderWorks(tokenURL string) bool ***REMOVED***
	for _, s := range brokenAuthHeaderProviders ***REMOVED***
		if strings.HasPrefix(tokenURL, s) ***REMOVED***
			// Some sites fail to implement the OAuth2 spec fully.
			return false
		***REMOVED***
	***REMOVED***

	// Assume the provider implements the spec properly
	// otherwise. We can add more exceptions as they're
	// discovered. We will _not_ be adding configurable hooks
	// to this package to let users select server bugs.
	return true
***REMOVED***

func RetrieveToken(ctx context.Context, clientID, clientSecret, tokenURL string, v url.Values) (*Token, error) ***REMOVED***
	hc, err := ContextClient(ctx)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	v.Set("client_id", clientID)
	bustedAuth := !providerAuthHeaderWorks(tokenURL)
	if bustedAuth && clientSecret != "" ***REMOVED***
		v.Set("client_secret", clientSecret)
	***REMOVED***
	req, err := http.NewRequest("POST", tokenURL, strings.NewReader(v.Encode()))
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	if !bustedAuth ***REMOVED***
		req.SetBasicAuth(clientID, clientSecret)
	***REMOVED***
	r, err := hc.Do(req)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	defer r.Body.Close()
	body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1<<20))
	if err != nil ***REMOVED***
		return nil, fmt.Errorf("oauth2: cannot fetch token: %v", err)
	***REMOVED***
	if code := r.StatusCode; code < 200 || code > 299 ***REMOVED***
		return nil, fmt.Errorf("oauth2: cannot fetch token: %v\nResponse: %s", r.Status, body)
	***REMOVED***

	var token *Token
	content, _, _ := mime.ParseMediaType(r.Header.Get("Content-Type"))
	switch content ***REMOVED***
	case "application/x-www-form-urlencoded", "text/plain":
		vals, err := url.ParseQuery(string(body))
		if err != nil ***REMOVED***
			return nil, err
		***REMOVED***
		token = &Token***REMOVED***
			AccessToken:  vals.Get("access_token"),
			TokenType:    vals.Get("token_type"),
			RefreshToken: vals.Get("refresh_token"),
			Raw:          vals,
		***REMOVED***
		e := vals.Get("expires_in")
		if e == "" ***REMOVED***
			// TODO(jbd): Facebook's OAuth2 implementation is broken and
			// returns expires_in field in expires. Remove the fallback to expires,
			// when Facebook fixes their implementation.
			e = vals.Get("expires")
		***REMOVED***
		expires, _ := strconv.Atoi(e)
		if expires != 0 ***REMOVED***
			token.Expiry = time.Now().Add(time.Duration(expires) * time.Second)
		***REMOVED***
	default:
		var tj tokenJSON
		if err = json.Unmarshal(body, &tj); err != nil ***REMOVED***
			return nil, err
		***REMOVED***
		token = &Token***REMOVED***
			AccessToken:  tj.AccessToken,
			TokenType:    tj.TokenType,
			RefreshToken: tj.RefreshToken,
			Expiry:       tj.expiry(),
			Raw:          make(map[string]interface***REMOVED******REMOVED***),
		***REMOVED***
		json.Unmarshal(body, &token.Raw) // no error checks for optional fields
	***REMOVED***
	// Don't overwrite `RefreshToken` with an empty value
	// if this was a token refreshing request.
	if token.RefreshToken == "" ***REMOVED***
		token.RefreshToken = v.Get("refresh_token")
	***REMOVED***
	return token, nil
***REMOVED***
