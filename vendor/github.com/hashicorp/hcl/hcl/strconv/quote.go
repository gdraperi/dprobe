package strconv

import (
	"errors"
	"unicode/utf8"
)

// ErrSyntax indicates that a value does not have the right syntax for the target type.
var ErrSyntax = errors.New("invalid syntax")

// Unquote interprets s as a single-quoted, double-quoted,
// or backquoted Go string literal, returning the string value
// that s quotes.  (If s is single-quoted, it would be a Go
// character literal; Unquote returns the corresponding
// one-character string.)
func Unquote(s string) (t string, err error) ***REMOVED***
	n := len(s)
	if n < 2 ***REMOVED***
		return "", ErrSyntax
	***REMOVED***
	quote := s[0]
	if quote != s[n-1] ***REMOVED***
		return "", ErrSyntax
	***REMOVED***
	s = s[1 : n-1]

	if quote != '"' ***REMOVED***
		return "", ErrSyntax
	***REMOVED***
	if !contains(s, '$') && !contains(s, '***REMOVED***') && contains(s, '\n') ***REMOVED***
		return "", ErrSyntax
	***REMOVED***

	// Is it trivial?  Avoid allocation.
	if !contains(s, '\\') && !contains(s, quote) && !contains(s, '$') ***REMOVED***
		switch quote ***REMOVED***
		case '"':
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
		// If we're starting a '$***REMOVED******REMOVED***' then let it through un-unquoted.
		// Specifically: we don't unquote any characters within the `$***REMOVED******REMOVED***`
		// section.
		if s[0] == '$' && len(s) > 1 && s[1] == '***REMOVED***' ***REMOVED***
			buf = append(buf, '$', '***REMOVED***')
			s = s[2:]

			// Continue reading until we find the closing brace, copying as-is
			braces := 1
			for len(s) > 0 && braces > 0 ***REMOVED***
				r, size := utf8.DecodeRuneInString(s)
				if r == utf8.RuneError ***REMOVED***
					return "", ErrSyntax
				***REMOVED***

				s = s[size:]

				n := utf8.EncodeRune(runeTmp[:], r)
				buf = append(buf, runeTmp[:n]...)

				switch r ***REMOVED***
				case '***REMOVED***':
					braces++
				case '***REMOVED***':
					braces--
				***REMOVED***
			***REMOVED***
			if braces != 0 ***REMOVED***
				return "", ErrSyntax
			***REMOVED***
			if len(s) == 0 ***REMOVED***
				// If there's no string left, we're done!
				break
			***REMOVED*** else ***REMOVED***
				// If there's more left, we need to pop back up to the top of the loop
				// in case there's another interpolation in this string.
				continue
			***REMOVED***
		***REMOVED***

		if s[0] == '\n' ***REMOVED***
			return "", ErrSyntax
		***REMOVED***

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
			return "", ErrSyntax
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

func unquoteChar(s string, quote byte) (value rune, multibyte bool, tail string, err error) ***REMOVED***
	// easy cases
	switch c := s[0]; ***REMOVED***
	case c == quote && (quote == '\'' || quote == '"'):
		err = ErrSyntax
		return
	case c >= utf8.RuneSelf:
		r, size := utf8.DecodeRuneInString(s)
		return r, true, s[size:], nil
	case c != '\\':
		return rune(s[0]), false, s[1:], nil
	***REMOVED***

	// hard case: c is backslash
	if len(s) <= 1 ***REMOVED***
		err = ErrSyntax
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
			err = ErrSyntax
			return
		***REMOVED***
		for j := 0; j < n; j++ ***REMOVED***
			x, ok := unhex(s[j])
			if !ok ***REMOVED***
				err = ErrSyntax
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
			err = ErrSyntax
			return
		***REMOVED***
		value = v
		multibyte = true
	case '0', '1', '2', '3', '4', '5', '6', '7':
		v := rune(c) - '0'
		if len(s) < 2 ***REMOVED***
			err = ErrSyntax
			return
		***REMOVED***
		for j := 0; j < 2; j++ ***REMOVED*** // one digit already; two more
			x := rune(s[j]) - '0'
			if x < 0 || x > 7 ***REMOVED***
				err = ErrSyntax
				return
			***REMOVED***
			v = (v << 3) | x
		***REMOVED***
		s = s[2:]
		if v > 255 ***REMOVED***
			err = ErrSyntax
			return
		***REMOVED***
		value = v
	case '\\':
		value = '\\'
	case '\'', '"':
		if c != quote ***REMOVED***
			err = ErrSyntax
			return
		***REMOVED***
		value = rune(c)
	default:
		err = ErrSyntax
		return
	***REMOVED***
	tail = s
	return
***REMOVED***
