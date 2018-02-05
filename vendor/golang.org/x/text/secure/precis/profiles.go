// Copyright 2015 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package precis

import (
	"unicode"

	"golang.org/x/text/runes"
	"golang.org/x/text/transform"
	"golang.org/x/text/unicode/norm"
)

var (
	// Implements the Nickname profile specified in RFC 8266.
	Nickname *Profile = nickname

	// Implements the UsernameCaseMapped profile specified in RFC 8265.
	UsernameCaseMapped *Profile = usernameCaseMap

	// Implements the UsernameCasePreserved profile specified in RFC 8265.
	UsernameCasePreserved *Profile = usernameNoCaseMap

	// Implements the OpaqueString profile defined in RFC 8265 for passwords and
	// other secure labels.
	OpaqueString *Profile = opaquestring
)

var (
	nickname = &Profile***REMOVED***
		options: getOpts(
			AdditionalMapping(func() transform.Transformer ***REMOVED***
				return &nickAdditionalMapping***REMOVED******REMOVED***
			***REMOVED***),
			IgnoreCase,
			Norm(norm.NFKC),
			DisallowEmpty,
			repeat,
		),
		class: freeform,
	***REMOVED***
	usernameCaseMap = &Profile***REMOVED***
		options: getOpts(
			FoldWidth,
			LowerCase(),
			Norm(norm.NFC),
			BidiRule,
		),
		class: identifier,
	***REMOVED***
	usernameNoCaseMap = &Profile***REMOVED***
		options: getOpts(
			FoldWidth,
			Norm(norm.NFC),
			BidiRule,
		),
		class: identifier,
	***REMOVED***
	opaquestring = &Profile***REMOVED***
		options: getOpts(
			AdditionalMapping(func() transform.Transformer ***REMOVED***
				return mapSpaces
			***REMOVED***),
			Norm(norm.NFC),
			DisallowEmpty,
		),
		class: freeform,
	***REMOVED***
)

// mapSpaces is a shared value of a runes.Map transformer.
var mapSpaces transform.Transformer = runes.Map(func(r rune) rune ***REMOVED***
	if unicode.Is(unicode.Zs, r) ***REMOVED***
		return ' '
	***REMOVED***
	return r
***REMOVED***)
