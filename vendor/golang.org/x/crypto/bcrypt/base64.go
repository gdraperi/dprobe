// Copyright 2011 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package bcrypt

import "encoding/base64"

const alphabet = "./ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789"

var bcEncoding = base64.NewEncoding(alphabet)

func base64Encode(src []byte) []byte ***REMOVED***
	n := bcEncoding.EncodedLen(len(src))
	dst := make([]byte, n)
	bcEncoding.Encode(dst, src)
	for dst[n-1] == '=' ***REMOVED***
		n--
	***REMOVED***
	return dst[:n]
***REMOVED***

func base64Decode(src []byte) ([]byte, error) ***REMOVED***
	numOfEquals := 4 - (len(src) % 4)
	for i := 0; i < numOfEquals; i++ ***REMOVED***
		src = append(src, '=')
	***REMOVED***

	dst := make([]byte, bcEncoding.DecodedLen(len(src)))
	n, err := bcEncoding.Decode(dst, src)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	return dst[:n], nil
***REMOVED***
