package hcl

import (
	"unicode"
	"unicode/utf8"
)

type lexModeValue byte

const (
	lexModeUnknown lexModeValue = iota
	lexModeHcl
	lexModeJson
)

// lexMode returns whether we're going to be parsing in JSON
// mode or HCL mode.
func lexMode(v []byte) lexModeValue ***REMOVED***
	var (
		r      rune
		w      int
		offset int
	)

	for ***REMOVED***
		r, w = utf8.DecodeRune(v[offset:])
		offset += w
		if unicode.IsSpace(r) ***REMOVED***
			continue
		***REMOVED***
		if r == '***REMOVED***' ***REMOVED***
			return lexModeJson
		***REMOVED***
		break
	***REMOVED***

	return lexModeHcl
***REMOVED***
