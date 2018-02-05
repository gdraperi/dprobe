// Copyright 2015 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:generate go run gen.go

// Package htmlindex maps character set encoding names to Encodings as
// recommended by the W3C for use in HTML 5. See http://www.w3.org/TR/encoding.
package htmlindex

// TODO: perhaps have a "bare" version of the index (used by this package) that
// is not pre-loaded with all encodings. Global variables in encodings prevent
// the linker from being able to purge unneeded tables. This means that
// referencing all encodings, as this package does for the default index, links
// in all encodings unconditionally.
//
// This issue can be solved by either solving the linking issue (see
// https://github.com/golang/go/issues/6330) or refactoring the encoding tables
// (e.g. moving the tables to internal packages that do not use global
// variables).

// TODO: allow canonicalizing names

import (
	"errors"
	"strings"
	"sync"

	"golang.org/x/text/encoding"
	"golang.org/x/text/encoding/internal/identifier"
	"golang.org/x/text/language"
)

var (
	errInvalidName = errors.New("htmlindex: invalid encoding name")
	errUnknown     = errors.New("htmlindex: unknown Encoding")
	errUnsupported = errors.New("htmlindex: this encoding is not supported")
)

var (
	matcherOnce sync.Once
	matcher     language.Matcher
)

// LanguageDefault returns the canonical name of the default encoding for a
// given language.
func LanguageDefault(tag language.Tag) string ***REMOVED***
	matcherOnce.Do(func() ***REMOVED***
		tags := []language.Tag***REMOVED******REMOVED***
		for _, t := range strings.Split(locales, " ") ***REMOVED***
			tags = append(tags, language.MustParse(t))
		***REMOVED***
		matcher = language.NewMatcher(tags, language.PreferSameScript(true))
	***REMOVED***)
	_, i, _ := matcher.Match(tag)
	return canonical[localeMap[i]] // Default is Windows-1252.
***REMOVED***

// Get returns an Encoding for one of the names listed in
// http://www.w3.org/TR/encoding using the Default Index. Matching is case-
// insensitive.
func Get(name string) (encoding.Encoding, error) ***REMOVED***
	x, ok := nameMap[strings.ToLower(strings.TrimSpace(name))]
	if !ok ***REMOVED***
		return nil, errInvalidName
	***REMOVED***
	return encodings[x], nil
***REMOVED***

// Name reports the canonical name of the given Encoding. It will return
// an error if e is not associated with a supported encoding scheme.
func Name(e encoding.Encoding) (string, error) ***REMOVED***
	id, ok := e.(identifier.Interface)
	if !ok ***REMOVED***
		return "", errUnknown
	***REMOVED***
	mib, _ := id.ID()
	if mib == 0 ***REMOVED***
		return "", errUnknown
	***REMOVED***
	v, ok := mibMap[mib]
	if !ok ***REMOVED***
		return "", errUnsupported
	***REMOVED***
	return canonical[v], nil
***REMOVED***
