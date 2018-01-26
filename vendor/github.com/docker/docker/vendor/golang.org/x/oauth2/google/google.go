// Copyright 2014 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package google provides support for making OAuth2 authorized and
// authenticated HTTP requests to Google APIs.
// It supports the Web server flow, client-side credentials, service accounts,
// Google Compute Engine service accounts, and Google App Engine service
// accounts.
//
// For more information, please read
// https://developers.google.com/accounts/docs/OAuth2
// and
// https://developers.google.com/accounts/docs/application-default-credentials.
package google // import "golang.org/x/oauth2/google"

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	"cloud.google.com/go/compute/metadata"
	"golang.org/x/net/context"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/jwt"
)

// Endpoint is Google's OAuth 2.0 endpoint.
var Endpoint = oauth2.Endpoint***REMOVED***
	AuthURL:  "https://accounts.google.com/o/oauth2/auth",
	TokenURL: "https://accounts.google.com/o/oauth2/token",
***REMOVED***

// JWTTokenURL is Google's OAuth 2.0 token URL to use with the JWT flow.
const JWTTokenURL = "https://accounts.google.com/o/oauth2/token"

// ConfigFromJSON uses a Google Developers Console client_credentials.json
// file to construct a config.
// client_credentials.json can be downloaded from
// https://console.developers.google.com, under "Credentials". Download the Web
// application credentials in the JSON format and provide the contents of the
// file as jsonKey.
func ConfigFromJSON(jsonKey []byte, scope ...string) (*oauth2.Config, error) ***REMOVED***
	type cred struct ***REMOVED***
		ClientID     string   `json:"client_id"`
		ClientSecret string   `json:"client_secret"`
		RedirectURIs []string `json:"redirect_uris"`
		AuthURI      string   `json:"auth_uri"`
		TokenURI     string   `json:"token_uri"`
	***REMOVED***
	var j struct ***REMOVED***
		Web       *cred `json:"web"`
		Installed *cred `json:"installed"`
	***REMOVED***
	if err := json.Unmarshal(jsonKey, &j); err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	var c *cred
	switch ***REMOVED***
	case j.Web != nil:
		c = j.Web
	case j.Installed != nil:
		c = j.Installed
	default:
		return nil, fmt.Errorf("oauth2/google: no credentials found")
	***REMOVED***
	if len(c.RedirectURIs) < 1 ***REMOVED***
		return nil, errors.New("oauth2/google: missing redirect URL in the client_credentials.json")
	***REMOVED***
	return &oauth2.Config***REMOVED***
		ClientID:     c.ClientID,
		ClientSecret: c.ClientSecret,
		RedirectURL:  c.RedirectURIs[0],
		Scopes:       scope,
		Endpoint: oauth2.Endpoint***REMOVED***
			AuthURL:  c.AuthURI,
			TokenURL: c.TokenURI,
		***REMOVED***,
	***REMOVED***, nil
***REMOVED***

// JWTConfigFromJSON uses a Google Developers service account JSON key file to read
// the credentials that authorize and authenticate the requests.
// Create a service account on "Credentials" for your project at
// https://console.developers.google.com to download a JSON key file.
func JWTConfigFromJSON(jsonKey []byte, scope ...string) (*jwt.Config, error) ***REMOVED***
	var f credentialsFile
	if err := json.Unmarshal(jsonKey, &f); err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	if f.Type != serviceAccountKey ***REMOVED***
		return nil, fmt.Errorf("google: read JWT from JSON credentials: 'type' field is %q (expected %q)", f.Type, serviceAccountKey)
	***REMOVED***
	scope = append([]string(nil), scope...) // copy
	return f.jwtConfig(scope), nil
***REMOVED***

// JSON key file types.
const (
	serviceAccountKey  = "service_account"
	userCredentialsKey = "authorized_user"
)

// credentialsFile is the unmarshalled representation of a credentials file.
type credentialsFile struct ***REMOVED***
	Type string `json:"type"` // serviceAccountKey or userCredentialsKey

	// Service Account fields
	ClientEmail  string `json:"client_email"`
	PrivateKeyID string `json:"private_key_id"`
	PrivateKey   string `json:"private_key"`
	TokenURL     string `json:"token_uri"`
	ProjectID    string `json:"project_id"`

	// User Credential fields
	// (These typically come from gcloud auth.)
	ClientSecret string `json:"client_secret"`
	ClientID     string `json:"client_id"`
	RefreshToken string `json:"refresh_token"`
***REMOVED***

func (f *credentialsFile) jwtConfig(scopes []string) *jwt.Config ***REMOVED***
	cfg := &jwt.Config***REMOVED***
		Email:        f.ClientEmail,
		PrivateKey:   []byte(f.PrivateKey),
		PrivateKeyID: f.PrivateKeyID,
		Scopes:       scopes,
		TokenURL:     f.TokenURL,
	***REMOVED***
	if cfg.TokenURL == "" ***REMOVED***
		cfg.TokenURL = JWTTokenURL
	***REMOVED***
	return cfg
***REMOVED***

func (f *credentialsFile) tokenSource(ctx context.Context, scopes []string) (oauth2.TokenSource, error) ***REMOVED***
	switch f.Type ***REMOVED***
	case serviceAccountKey:
		cfg := f.jwtConfig(scopes)
		return cfg.TokenSource(ctx), nil
	case userCredentialsKey:
		cfg := &oauth2.Config***REMOVED***
			ClientID:     f.ClientID,
			ClientSecret: f.ClientSecret,
			Scopes:       scopes,
			Endpoint:     Endpoint,
		***REMOVED***
		tok := &oauth2.Token***REMOVED***RefreshToken: f.RefreshToken***REMOVED***
		return cfg.TokenSource(ctx, tok), nil
	case "":
		return nil, errors.New("missing 'type' field in credentials")
	default:
		return nil, fmt.Errorf("unknown credential type: %q", f.Type)
	***REMOVED***
***REMOVED***

// ComputeTokenSource returns a token source that fetches access tokens
// from Google Compute Engine (GCE)'s metadata server. It's only valid to use
// this token source if your program is running on a GCE instance.
// If no account is specified, "default" is used.
// Further information about retrieving access tokens from the GCE metadata
// server can be found at https://cloud.google.com/compute/docs/authentication.
func ComputeTokenSource(account string) oauth2.TokenSource ***REMOVED***
	return oauth2.ReuseTokenSource(nil, computeSource***REMOVED***account: account***REMOVED***)
***REMOVED***

type computeSource struct ***REMOVED***
	account string
***REMOVED***

func (cs computeSource) Token() (*oauth2.Token, error) ***REMOVED***
	if !metadata.OnGCE() ***REMOVED***
		return nil, errors.New("oauth2/google: can't get a token from the metadata service; not running on GCE")
	***REMOVED***
	acct := cs.account
	if acct == "" ***REMOVED***
		acct = "default"
	***REMOVED***
	tokenJSON, err := metadata.Get("instance/service-accounts/" + acct + "/token")
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	var res struct ***REMOVED***
		AccessToken  string `json:"access_token"`
		ExpiresInSec int    `json:"expires_in"`
		TokenType    string `json:"token_type"`
	***REMOVED***
	err = json.NewDecoder(strings.NewReader(tokenJSON)).Decode(&res)
	if err != nil ***REMOVED***
		return nil, fmt.Errorf("oauth2/google: invalid token JSON from metadata: %v", err)
	***REMOVED***
	if res.ExpiresInSec == 0 || res.AccessToken == "" ***REMOVED***
		return nil, fmt.Errorf("oauth2/google: incomplete token received from metadata")
	***REMOVED***
	return &oauth2.Token***REMOVED***
		AccessToken: res.AccessToken,
		TokenType:   res.TokenType,
		Expiry:      time.Now().Add(time.Duration(res.ExpiresInSec) * time.Second),
	***REMOVED***, nil
***REMOVED***
