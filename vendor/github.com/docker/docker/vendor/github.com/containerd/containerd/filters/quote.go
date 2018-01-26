package filters

import (
	"unicode/utf8"

	"github.com/pkg/errors"
)

// NOTE(stevvooe): Most of this code in this file is copied from the stdlib
// strconv package and modified to be able to handle quoting with `/` and `|`
// as delimiters.  The copyright is held by the Go authors.

var errQuoteSyntax = errors.New("quote syntax error")

// UnquoteChar decodes the first character or byte in the escaped string
// or character literal represented by the string s.
// It returns four values:
//
//	1) value, the decoded Unicode code point or byte value;
//	2) multibyte, a boolean indicating whether the decoded character requires a multibyte UTF-8 representation;
//	3) tail, the remainder of the string after the character; and
//	4) an error that will be nil if the character is syntactically valid.
//
// The second argument, quote, specifies the type of literal being parsed
// and therefore which escaped quote character is permitted.
// If set to a single quote, it permits the sequence \' and disallows unescaped '.
// If set to a double quote, it permits \" and disallows unescaped ".
// If set to zero, it does not permit either escape and allows both quote characters to appear unescaped.
//
// This is from Go strconv package, modified to support `|` and `/` as double
// quotes for use with regular expressions.
func unquoteChar(s string, quote byte) (value rune, multibyte bool, tail string, err error) ***REMOVED***
	// easy cases
	switch c := s[0]; ***REMOVED***
	case c == quote && (quote == '\'' || quote == '"' || quote == '/' || quote == '|'):
		err = errQuoteSyntax
		return
	case c >= utf8.RuneSelf:
		r, size := utf8.DecodeRuneInString(s)
		return r, true, s[size:], nil
	case c != '\\':
		return rune(s[0]), false, s[1:], nil
	***REMOVED***

	// hard case: c is backslash
	if len(s) <= 1 ***REMOVED***
		err = errQuoteSyntax
		return
	***REMOVED***
	c := s[1]
	s = s[2:]

	switch c ***REMOVED***
	case 'a':
		value = '\a'
	case 'b':
		value = '\b'
	case 'f':
		value = '\f'
	case 'n':
		value = '\n'
	case 'r':
		value = '\r'
	case 't':
		value = '\t'
	case 'v':
		value = '\v'
	case 'x', 'u', 'U':
		n := 0
		switch c ***REMOVED***
		case 'x':
			n = 2
		case 'u':
			n = 4
		case 'U':
			n = 8
		***REMOVED***
		var v rune
		if len(s) < n ***REMOVED***
			err = errQuoteSyntax
			return
		***REMOVED***
		for j := 0; j < n; j++ ***REMOVED***
			x, ok := unhex(s[j])
			if !ok ***REMOVED***
				err = errQuoteSyntax
				return
			***REMOVED***
			v = v<<4 | x
		***REMOVED***
		s = s[n:]
		if c == 'x' ***REMOVED***
			// single-byte string, possibly not UTF-8
			value = v
			break
		***REMOVED***
		if v > utf8.MaxRune ***REMOVED***
			err = errQuoteSyntax
			return
		***REMOVED***
		value = v
		multibyte = true
	case '0', '1', '2', '3', '4', '5', '6', '7':
		v := rune(c) - '0'
		if len(s) < 2 ***REMOVED***
			err = errQuoteSyntax
			return
		***REMOVED***
		for j := 0; j < 2; j++ ***REMOVED*** // one digit already; two more
			x := rune(s[j]) - '0'
			if x < 0 || x > 7 ***REMOVED***
				err = errQuoteSyntax
				return
			***REMOVED***
			v = (v << 3) | x
		***REMOVED***
		s = s[2:]
		if v > 255 ***REMOVED***
			err = errQuoteSyntax
			return
		***REMOVED***
		value = v
	case '\\':
		value = '\\'
	case '\'', '"', '|', '/':
		if c != quote ***REMOVED***
			err = errQuoteSyntax
			return
		***REMOVED***
		value = rune(c)
	default:
		err = errQuoteSyntax
		return
	***REMOVED***
	tail = s
	return
***REMOVED***

// unquote interprets s as a single-quoted, double-quoted,
// or backquoted Go string literal, returning the string value
// that s quotes.  (If s is single-quoted, it would be a Go
// character literal; Unquote returns the corresponding
// one-character string.)
//
// This is modified from the standard library to support `|` and `/` as quote
// characters for use with regular expressions.
func unquote(s string) (string, error) ***REMOVED***
	n := len(s)
	if n < 2 ***REMOVED***
		return "", errQuoteSyntax
	***REMOVED***
	quote := s[0]
	if quote != s[n-1] ***REMOVED***
		return "", errQuoteSyntax
	***REMOVED***
	s = s[1 : n-1]

	if quote == '`' ***REMOVED***
		if contains(s, '`') ***REMOVED***
			return "", errQuoteSyntax
		***REMOVED***
		if contains(s, '\r') ***REMOVED***
			// -1 because we know there is at least one \r to remove.
			buf := make([]byte, 0, len(s)-1)
			for i := 0; i < len(s); i++ ***REMOVED***
				if s[i] != '\r' ***REMOVED***
					buf = append(buf, s[i])
				***REMOVED***
			***REMOVED***
			return string(buf), nil
		***REMOVED***
		return s, nil
	***REMOVED***
	if quote != '"' && quote != '\'' && quote != '|' && quote != '/' ***REMOVED***
		return "", errQuoteSyntax
	***REMOVED***
	if contains(s, '\n') ***REMOVED***
		return "", errQuoteSyntax
	***REMOVED***

	// Is it trivial?  Avoid allocation.
	if !contains(s, '\\') && !contains(s, quote) ***REMOVED***
		switch quote ***REMOVED***
		case '"', '/', '|': // pipe and slash are treated like double quote
			return s, nil
		case '\'':
			r, size := utf8.DecodeRuneInString(s)
			if size == len(s) && (r != utf8.RuneError || size != 1) ***REMOVED***
				return s, nil
			***REMOVED***
		***REMOVED***
	***REMOVED***

	var runeTmp [utf8.UTFMax]byte
	buf := make([]byte, 0, 3*len(s)/2) // Try to avoid more allocations.
	for len(s) > 0 ***REMOVED***
		c, multibyte, ss, err := unquoteChar(s, quote)
		if err != nil ***REMOVED***
			return "", err
		***REMOVED***
		s = ss
		if c < utf8.RuneSelf || !multibyte ***REMOVED***
			buf = append(buf, byte(c))
		***REMOVED*** else ***REMOVED***
			n := utf8.EncodeRune(runeTmp[:], c)
			buf = append(buf, runeTmp[:n]...)
		***REMOVED***
		if quote == '\'' && len(s) != 0 ***REMOVED***
			// single-quoted must be single character
			return "", errQuoteSyntax
		***REMOVED***
	***REMOVED***
	return string(buf), nil
***REMOVED***

// contains reports whether the string contains the byte c.
func contains(s string, c byte) bool ***REMOVED***
	for i := 0; i < len(s); i++ ***REMOVED***
		if s[i] == c ***REMOVED***
			return true
		***REMOVED***
	***REMOVED***
	return false
***REMOVED***

func unhex(b byte) (v rune, ok bool) ***REMOVED***
	c := rune(b)
	switch ***REMOVED***
	case '0' <= c && c <= '9':
		return c - '0', true
	case 'a' <= c && c <= 'f':
		return c - 'a' + 10, true
	case 'A' <= c && c <= 'F':
		return c - 'A' + 10, true
	***REMOVED***
	return
***REMOVED***
