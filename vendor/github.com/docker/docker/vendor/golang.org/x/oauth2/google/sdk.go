// Copyright 2015 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package google

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/user"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"golang.org/x/net/context"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/internal"
)

type sdkCredentials struct ***REMOVED***
	Data []struct ***REMOVED***
		Credential struct ***REMOVED***
			ClientID     string     `json:"client_id"`
			ClientSecret string     `json:"client_secret"`
			AccessToken  string     `json:"access_token"`
			RefreshToken string     `json:"refresh_token"`
			TokenExpiry  *time.Time `json:"token_expiry"`
		***REMOVED*** `json:"credential"`
		Key struct ***REMOVED***
			Account string `json:"account"`
			Scope   string `json:"scope"`
		***REMOVED*** `json:"key"`
	***REMOVED***
***REMOVED***

// An SDKConfig provides access to tokens from an account already
// authorized via the Google Cloud SDK.
type SDKConfig struct ***REMOVED***
	conf         oauth2.Config
	initialToken *oauth2.Token
***REMOVED***

// NewSDKConfig creates an SDKConfig for the given Google Cloud SDK
// account. If account is empty, the account currently active in
// Google Cloud SDK properties is used.
// Google Cloud SDK credentials must be created by running `gcloud auth`
// before using this function.
// The Google Cloud SDK is available at https://cloud.google.com/sdk/.
func NewSDKConfig(account string) (*SDKConfig, error) ***REMOVED***
	configPath, err := sdkConfigPath()
	if err != nil ***REMOVED***
		return nil, fmt.Errorf("oauth2/google: error getting SDK config path: %v", err)
	***REMOVED***
	credentialsPath := filepath.Join(configPath, "credentials")
	f, err := os.Open(credentialsPath)
	if err != nil ***REMOVED***
		return nil, fmt.Errorf("oauth2/google: failed to load SDK credentials: %v", err)
	***REMOVED***
	defer f.Close()

	var c sdkCredentials
	if err := json.NewDecoder(f).Decode(&c); err != nil ***REMOVED***
		return nil, fmt.Errorf("oauth2/google: failed to decode SDK credentials from %q: %v", credentialsPath, err)
	***REMOVED***
	if len(c.Data) == 0 ***REMOVED***
		return nil, fmt.Errorf("oauth2/google: no credentials found in %q, run `gcloud auth login` to create one", credentialsPath)
	***REMOVED***
	if account == "" ***REMOVED***
		propertiesPath := filepath.Join(configPath, "properties")
		f, err := os.Open(propertiesPath)
		if err != nil ***REMOVED***
			return nil, fmt.Errorf("oauth2/google: failed to load SDK properties: %v", err)
		***REMOVED***
		defer f.Close()
		ini, err := internal.ParseINI(f)
		if err != nil ***REMOVED***
			return nil, fmt.Errorf("oauth2/google: failed to parse SDK properties %q: %v", propertiesPath, err)
		***REMOVED***
		core, ok := ini["core"]
		if !ok ***REMOVED***
			return nil, fmt.Errorf("oauth2/google: failed to find [core] section in %v", ini)
		***REMOVED***
		active, ok := core["account"]
		if !ok ***REMOVED***
			return nil, fmt.Errorf("oauth2/google: failed to find %q attribute in %v", "account", core)
		***REMOVED***
		account = active
	***REMOVED***

	for _, d := range c.Data ***REMOVED***
		if account == "" || d.Key.Account == account ***REMOVED***
			if d.Credential.AccessToken == "" && d.Credential.RefreshToken == "" ***REMOVED***
				return nil, fmt.Errorf("oauth2/google: no token available for account %q", account)
			***REMOVED***
			var expiry time.Time
			if d.Credential.TokenExpiry != nil ***REMOVED***
				expiry = *d.Credential.TokenExpiry
			***REMOVED***
			return &SDKConfig***REMOVED***
				conf: oauth2.Config***REMOVED***
					ClientID:     d.Credential.ClientID,
					ClientSecret: d.Credential.ClientSecret,
					Scopes:       strings.Split(d.Key.Scope, " "),
					Endpoint:     Endpoint,
					RedirectURL:  "oob",
				***REMOVED***,
				initialToken: &oauth2.Token***REMOVED***
					AccessToken:  d.Credential.AccessToken,
					RefreshToken: d.Credential.RefreshToken,
					Expiry:       expiry,
				***REMOVED***,
			***REMOVED***, nil
		***REMOVED***
	***REMOVED***
	return nil, fmt.Errorf("oauth2/google: no such credentials for account %q", account)
***REMOVED***

// Client returns an HTTP client using Google Cloud SDK credentials to
// authorize requests. The token will auto-refresh as necessary. The
// underlying http.RoundTripper will be obtained using the provided
// context. The returned client and its Transport should not be
// modified.
func (c *SDKConfig) Client(ctx context.Context) *http.Client ***REMOVED***
	return &http.Client***REMOVED***
		Transport: &oauth2.Transport***REMOVED***
			Source: c.TokenSource(ctx),
		***REMOVED***,
	***REMOVED***
***REMOVED***

// TokenSource returns an oauth2.TokenSource that retrieve tokens from
// Google Cloud SDK credentials using the provided context.
// It will returns the current access token stored in the credentials,
// and refresh it when it expires, but it won't update the credentials
// with the new access token.
func (c *SDKConfig) TokenSource(ctx context.Context) oauth2.TokenSource ***REMOVED***
	return c.conf.TokenSource(ctx, c.initialToken)
***REMOVED***

// Scopes are the OAuth 2.0 scopes the current account is authorized for.
func (c *SDKConfig) Scopes() []string ***REMOVED***
	return c.conf.Scopes
***REMOVED***

// sdkConfigPath tries to guess where the gcloud config is located.
// It can be overridden during tests.
var sdkConfigPath = func() (string, error) ***REMOVED***
	if runtime.GOOS == "windows" ***REMOVED***
		return filepath.Join(os.Getenv("APPDATA"), "gcloud"), nil
	***REMOVED***
	homeDir := guessUnixHomeDir()
	if homeDir == "" ***REMOVED***
		return "", errors.New("unable to get current user home directory: os/user lookup failed; $HOME is empty")
	***REMOVED***
	return filepath.Join(homeDir, ".config", "gcloud"), nil
***REMOVED***

func guessUnixHomeDir() string ***REMOVED***
	// Prefer $HOME over user.Current due to glibc bug: golang.org/issue/13470
	if v := os.Getenv("HOME"); v != "" ***REMOVED***
		return v
	***REMOVED***
	// Else, fall back to user.Current:
	if u, err := user.Current(); err == nil ***REMOVED***
		return u.HomeDir
	***REMOVED***
	return ""
***REMOVED***
