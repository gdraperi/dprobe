// Copyright 2012 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package xsrftoken

import (
	"encoding/base64"
	"testing"
	"time"
)

const (
	key      = "quay"
	userID   = "12345678"
	actionID = "POST /form"
)

var (
	now              = time.Now()
	oneMinuteFromNow = now.Add(1 * time.Minute)
)

func TestValidToken(t *testing.T) ***REMOVED***
	tok := generateTokenAtTime(key, userID, actionID, now)
	if !validTokenAtTime(tok, key, userID, actionID, oneMinuteFromNow) ***REMOVED***
		t.Error("One second later: Expected token to be valid")
	***REMOVED***
	if !validTokenAtTime(tok, key, userID, actionID, now.Add(Timeout-1*time.Nanosecond)) ***REMOVED***
		t.Error("Just before timeout: Expected token to be valid")
	***REMOVED***
	if !validTokenAtTime(tok, key, userID, actionID, now.Add(-1*time.Minute+1*time.Millisecond)) ***REMOVED***
		t.Error("One minute in the past: Expected token to be valid")
	***REMOVED***
***REMOVED***

// TestSeparatorReplacement tests that separators are being correctly substituted
func TestSeparatorReplacement(t *testing.T) ***REMOVED***
	tok := generateTokenAtTime("foo:bar", "baz", "wah", now)
	tok2 := generateTokenAtTime("foo", "bar:baz", "wah", now)
	if tok == tok2 ***REMOVED***
		t.Errorf("Expected generated tokens to be different")
	***REMOVED***
***REMOVED***

func TestInvalidToken(t *testing.T) ***REMOVED***
	invalidTokenTests := []struct ***REMOVED***
		name, key, userID, actionID string
		t                           time.Time
	***REMOVED******REMOVED***
		***REMOVED***"Bad key", "foobar", userID, actionID, oneMinuteFromNow***REMOVED***,
		***REMOVED***"Bad userID", key, "foobar", actionID, oneMinuteFromNow***REMOVED***,
		***REMOVED***"Bad actionID", key, userID, "foobar", oneMinuteFromNow***REMOVED***,
		***REMOVED***"Expired", key, userID, actionID, now.Add(Timeout + 1*time.Millisecond)***REMOVED***,
		***REMOVED***"More than 1 minute from the future", key, userID, actionID, now.Add(-1*time.Nanosecond - 1*time.Minute)***REMOVED***,
	***REMOVED***

	tok := generateTokenAtTime(key, userID, actionID, now)
	for _, itt := range invalidTokenTests ***REMOVED***
		if validTokenAtTime(tok, itt.key, itt.userID, itt.actionID, itt.t) ***REMOVED***
			t.Errorf("%v: Expected token to be invalid", itt.name)
		***REMOVED***
	***REMOVED***
***REMOVED***

// TestValidateBadData primarily tests that no unexpected panics are triggered
// during parsing
func TestValidateBadData(t *testing.T) ***REMOVED***
	badDataTests := []struct ***REMOVED***
		name, tok string
	***REMOVED******REMOVED***
		***REMOVED***"Invalid Base64", "ASDab24(@)$*=="***REMOVED***,
		***REMOVED***"No delimiter", base64.URLEncoding.EncodeToString([]byte("foobar12345678"))***REMOVED***,
		***REMOVED***"Invalid time", base64.URLEncoding.EncodeToString([]byte("foobar:foobar"))***REMOVED***,
		***REMOVED***"Wrong length", "1234" + generateTokenAtTime(key, userID, actionID, now)***REMOVED***,
	***REMOVED***

	for _, bdt := range badDataTests ***REMOVED***
		if validTokenAtTime(bdt.tok, key, userID, actionID, oneMinuteFromNow) ***REMOVED***
			t.Errorf("%v: Expected token to be invalid", bdt.name)
		***REMOVED***
	***REMOVED***
***REMOVED***
