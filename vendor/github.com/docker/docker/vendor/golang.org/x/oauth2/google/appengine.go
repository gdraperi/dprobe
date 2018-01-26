// Copyright 2014 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package google

import (
	"sort"
	"strings"
	"sync"
	"time"

	"golang.org/x/net/context"
	"golang.org/x/oauth2"
)

// Set at init time by appenginevm_hook.go. If true, we are on App Engine Managed VMs.
var appengineVM bool

// Set at init time by appengine_hook.go. If nil, we're not on App Engine.
var appengineTokenFunc func(c context.Context, scopes ...string) (token string, expiry time.Time, err error)

// Set at init time by appengine_hook.go. If nil, we're not on App Engine.
var appengineAppIDFunc func(c context.Context) string

// AppEngineTokenSource returns a token source that fetches tokens
// issued to the current App Engine application's service account.
// If you are implementing a 3-legged OAuth 2.0 flow on App Engine
// that involves user accounts, see oauth2.Config instead.
//
// The provided context must have come from appengine.NewContext.
func AppEngineTokenSource(ctx context.Context, scope ...string) oauth2.TokenSource ***REMOVED***
	if appengineTokenFunc == nil ***REMOVED***
		panic("google: AppEngineTokenSource can only be used on App Engine.")
	***REMOVED***
	scopes := append([]string***REMOVED******REMOVED***, scope...)
	sort.Strings(scopes)
	return &appEngineTokenSource***REMOVED***
		ctx:    ctx,
		scopes: scopes,
		key:    strings.Join(scopes, " "),
	***REMOVED***
***REMOVED***

// aeTokens helps the fetched tokens to be reused until their expiration.
var (
	aeTokensMu sync.Mutex
	aeTokens   = make(map[string]*tokenLock) // key is space-separated scopes
)

type tokenLock struct ***REMOVED***
	mu sync.Mutex // guards t; held while fetching or updating t
	t  *oauth2.Token
***REMOVED***

type appEngineTokenSource struct ***REMOVED***
	ctx    context.Context
	scopes []string
	key    string // to aeTokens map; space-separated scopes
***REMOVED***

func (ts *appEngineTokenSource) Token() (*oauth2.Token, error) ***REMOVED***
	if appengineTokenFunc == nil ***REMOVED***
		panic("google: AppEngineTokenSource can only be used on App Engine.")
	***REMOVED***

	aeTokensMu.Lock()
	tok, ok := aeTokens[ts.key]
	if !ok ***REMOVED***
		tok = &tokenLock***REMOVED******REMOVED***
		aeTokens[ts.key] = tok
	***REMOVED***
	aeTokensMu.Unlock()

	tok.mu.Lock()
	defer tok.mu.Unlock()
	if tok.t.Valid() ***REMOVED***
		return tok.t, nil
	***REMOVED***
	access, exp, err := appengineTokenFunc(ts.ctx, ts.scopes...)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	tok.t = &oauth2.Token***REMOVED***
		AccessToken: access,
		Expiry:      exp,
	***REMOVED***
	return tok.t, nil
***REMOVED***
