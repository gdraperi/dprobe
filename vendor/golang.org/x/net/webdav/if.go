// Copyright 2014 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package webdav

// The If header is covered by Section 10.4.
// http://www.webdav.org/specs/rfc4918.html#HEADER_If

import (
	"strings"
)

// ifHeader is a disjunction (OR) of ifLists.
type ifHeader struct ***REMOVED***
	lists []ifList
***REMOVED***

// ifList is a conjunction (AND) of Conditions, and an optional resource tag.
type ifList struct ***REMOVED***
	resourceTag string
	conditions  []Condition
***REMOVED***

// parseIfHeader parses the "If: foo bar" HTTP header. The httpHeader string
// should omit the "If:" prefix and have any "\r\n"s collapsed to a " ", as is
// returned by req.Header.Get("If") for a http.Request req.
func parseIfHeader(httpHeader string) (h ifHeader, ok bool) ***REMOVED***
	s := strings.TrimSpace(httpHeader)
	switch tokenType, _, _ := lex(s); tokenType ***REMOVED***
	case '(':
		return parseNoTagLists(s)
	case angleTokenType:
		return parseTaggedLists(s)
	default:
		return ifHeader***REMOVED******REMOVED***, false
	***REMOVED***
***REMOVED***

func parseNoTagLists(s string) (h ifHeader, ok bool) ***REMOVED***
	for ***REMOVED***
		l, remaining, ok := parseList(s)
		if !ok ***REMOVED***
			return ifHeader***REMOVED******REMOVED***, false
		***REMOVED***
		h.lists = append(h.lists, l)
		if remaining == "" ***REMOVED***
			return h, true
		***REMOVED***
		s = remaining
	***REMOVED***
***REMOVED***

func parseTaggedLists(s string) (h ifHeader, ok bool) ***REMOVED***
	resourceTag, n := "", 0
	for first := true; ; first = false ***REMOVED***
		tokenType, tokenStr, remaining := lex(s)
		switch tokenType ***REMOVED***
		case angleTokenType:
			if !first && n == 0 ***REMOVED***
				return ifHeader***REMOVED******REMOVED***, false
			***REMOVED***
			resourceTag, n = tokenStr, 0
			s = remaining
		case '(':
			n++
			l, remaining, ok := parseList(s)
			if !ok ***REMOVED***
				return ifHeader***REMOVED******REMOVED***, false
			***REMOVED***
			l.resourceTag = resourceTag
			h.lists = append(h.lists, l)
			if remaining == "" ***REMOVED***
				return h, true
			***REMOVED***
			s = remaining
		default:
			return ifHeader***REMOVED******REMOVED***, false
		***REMOVED***
	***REMOVED***
***REMOVED***

func parseList(s string) (l ifList, remaining string, ok bool) ***REMOVED***
	tokenType, _, s := lex(s)
	if tokenType != '(' ***REMOVED***
		return ifList***REMOVED******REMOVED***, "", false
	***REMOVED***
	for ***REMOVED***
		tokenType, _, remaining = lex(s)
		if tokenType == ')' ***REMOVED***
			if len(l.conditions) == 0 ***REMOVED***
				return ifList***REMOVED******REMOVED***, "", false
			***REMOVED***
			return l, remaining, true
		***REMOVED***
		c, remaining, ok := parseCondition(s)
		if !ok ***REMOVED***
			return ifList***REMOVED******REMOVED***, "", false
		***REMOVED***
		l.conditions = append(l.conditions, c)
		s = remaining
	***REMOVED***
***REMOVED***

func parseCondition(s string) (c Condition, remaining string, ok bool) ***REMOVED***
	tokenType, tokenStr, s := lex(s)
	if tokenType == notTokenType ***REMOVED***
		c.Not = true
		tokenType, tokenStr, s = lex(s)
	***REMOVED***
	switch tokenType ***REMOVED***
	case strTokenType, angleTokenType:
		c.Token = tokenStr
	case squareTokenType:
		c.ETag = tokenStr
	default:
		return Condition***REMOVED******REMOVED***, "", false
	***REMOVED***
	return c, s, true
***REMOVED***

// Single-rune tokens like '(' or ')' have a token type equal to their rune.
// All other tokens have a negative token type.
const (
	errTokenType    = rune(-1)
	eofTokenType    = rune(-2)
	strTokenType    = rune(-3)
	notTokenType    = rune(-4)
	angleTokenType  = rune(-5)
	squareTokenType = rune(-6)
)

func lex(s string) (tokenType rune, tokenStr string, remaining string) ***REMOVED***
	// The net/textproto Reader that parses the HTTP header will collapse
	// Linear White Space that spans multiple "\r\n" lines to a single " ",
	// so we don't need to look for '\r' or '\n'.
	for len(s) > 0 && (s[0] == '\t' || s[0] == ' ') ***REMOVED***
		s = s[1:]
	***REMOVED***
	if len(s) == 0 ***REMOVED***
		return eofTokenType, "", ""
	***REMOVED***
	i := 0
loop:
	for ; i < len(s); i++ ***REMOVED***
		switch s[i] ***REMOVED***
		case '\t', ' ', '(', ')', '<', '>', '[', ']':
			break loop
		***REMOVED***
	***REMOVED***

	if i != 0 ***REMOVED***
		tokenStr, remaining = s[:i], s[i:]
		if tokenStr == "Not" ***REMOVED***
			return notTokenType, "", remaining
		***REMOVED***
		return strTokenType, tokenStr, remaining
	***REMOVED***

	j := 0
	switch s[0] ***REMOVED***
	case '<':
		j, tokenType = strings.IndexByte(s, '>'), angleTokenType
	case '[':
		j, tokenType = strings.IndexByte(s, ']'), squareTokenType
	default:
		return rune(s[0]), "", s[1:]
	***REMOVED***
	if j < 0 ***REMOVED***
		return errTokenType, "", ""
	***REMOVED***
	return tokenType, s[1:j], s[j+1:]
***REMOVED***
