// Copyright 2012 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package ssh

// Message authentication support

import (
	"crypto/hmac"
	"crypto/sha1"
	"crypto/sha256"
	"hash"
)

type macMode struct ***REMOVED***
	keySize int
	etm     bool
	new     func(key []byte) hash.Hash
***REMOVED***

// truncatingMAC wraps around a hash.Hash and truncates the output digest to
// a given size.
type truncatingMAC struct ***REMOVED***
	length int
	hmac   hash.Hash
***REMOVED***

func (t truncatingMAC) Write(data []byte) (int, error) ***REMOVED***
	return t.hmac.Write(data)
***REMOVED***

func (t truncatingMAC) Sum(in []byte) []byte ***REMOVED***
	out := t.hmac.Sum(in)
	return out[:len(in)+t.length]
***REMOVED***

func (t truncatingMAC) Reset() ***REMOVED***
	t.hmac.Reset()
***REMOVED***

func (t truncatingMAC) Size() int ***REMOVED***
	return t.length
***REMOVED***

func (t truncatingMAC) BlockSize() int ***REMOVED*** return t.hmac.BlockSize() ***REMOVED***

var macModes = map[string]*macMode***REMOVED***
	"hmac-sha2-256-etm@openssh.com": ***REMOVED***32, true, func(key []byte) hash.Hash ***REMOVED***
		return hmac.New(sha256.New, key)
	***REMOVED******REMOVED***,
	"hmac-sha2-256": ***REMOVED***32, false, func(key []byte) hash.Hash ***REMOVED***
		return hmac.New(sha256.New, key)
	***REMOVED******REMOVED***,
	"hmac-sha1": ***REMOVED***20, false, func(key []byte) hash.Hash ***REMOVED***
		return hmac.New(sha1.New, key)
	***REMOVED******REMOVED***,
	"hmac-sha1-96": ***REMOVED***20, false, func(key []byte) hash.Hash ***REMOVED***
		return truncatingMAC***REMOVED***12, hmac.New(sha1.New, key)***REMOVED***
	***REMOVED******REMOVED***,
***REMOVED***
