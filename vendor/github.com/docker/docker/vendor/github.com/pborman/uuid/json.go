// Copyright 2014 Google Inc.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package uuid

import "errors"

func (u UUID) MarshalJSON() ([]byte, error) ***REMOVED***
	if len(u) != 16 ***REMOVED***
		return []byte(`""`), nil
	***REMOVED***
	var js [38]byte
	js[0] = '"'
	encodeHex(js[1:], u)
	js[37] = '"'
	return js[:], nil
***REMOVED***

func (u *UUID) UnmarshalJSON(data []byte) error ***REMOVED***
	if string(data) == `""` ***REMOVED***
		return nil
	***REMOVED***
	if data[0] != '"' ***REMOVED***
		return errors.New("invalid UUID format")
	***REMOVED***
	data = data[1 : len(data)-1]
	uu := Parse(string(data))
	if uu == nil ***REMOVED***
		return errors.New("invalid UUID format")
	***REMOVED***
	*u = uu
	return nil
***REMOVED***
