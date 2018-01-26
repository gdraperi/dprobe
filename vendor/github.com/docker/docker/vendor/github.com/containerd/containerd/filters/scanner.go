package filters

import (
	"fmt"
	"unicode"
	"unicode/utf8"
)

const (
	tokenEOF = -(iota + 1)
	tokenQuoted
	tokenValue
	tokenField
	tokenSeparator
	tokenOperator
	tokenIllegal
)

type token rune

func (t token) String() string ***REMOVED***
	switch t ***REMOVED***
	case tokenEOF:
		return "EOF"
	case tokenQuoted:
		return "Quoted"
	case tokenValue:
		return "Value"
	case tokenField:
		return "Field"
	case tokenSeparator:
		return "Separator"
	case tokenOperator:
		return "Operator"
	case tokenIllegal:
		return "Illegal"
	***REMOVED***

	return string(t)
***REMOVED***

func (t token) GoString() string ***REMOVED***
	return "token" + t.String()
***REMOVED***

type scanner struct ***REMOVED***
	input string
	pos   int
	ppos  int // bounds the current rune in the string
	value bool
***REMOVED***

func (s *scanner) init(input string) ***REMOVED***
	s.input = input
	s.pos = 0
	s.ppos = 0
***REMOVED***

func (s *scanner) next() rune ***REMOVED***
	if s.pos >= len(s.input) ***REMOVED***
		return tokenEOF
	***REMOVED***
	s.pos = s.ppos

	r, w := utf8.DecodeRuneInString(s.input[s.ppos:])
	s.ppos += w
	if r == utf8.RuneError ***REMOVED***
		if w > 0 ***REMOVED***
			return tokenIllegal
		***REMOVED***
		return tokenEOF
	***REMOVED***

	if r == 0 ***REMOVED***
		return tokenIllegal
	***REMOVED***

	return r
***REMOVED***

func (s *scanner) peek() rune ***REMOVED***
	pos := s.pos
	ppos := s.ppos
	ch := s.next()
	s.pos = pos
	s.ppos = ppos
	return ch
***REMOVED***

func (s *scanner) scan() (nextp int, tk token, text string) ***REMOVED***
	var (
		ch  = s.next()
		pos = s.pos
	)

chomp:
	switch ***REMOVED***
	case ch == tokenEOF:
	case ch == tokenIllegal:
	case isQuoteRune(ch):
		s.scanQuoted(ch)
		return pos, tokenQuoted, s.input[pos:s.ppos]
	case isSeparatorRune(ch):
		s.value = false
		return pos, tokenSeparator, s.input[pos:s.ppos]
	case isOperatorRune(ch):
		s.scanOperator()
		s.value = true
		return pos, tokenOperator, s.input[pos:s.ppos]
	case unicode.IsSpace(ch):
		// chomp
		ch = s.next()
		pos = s.pos
		goto chomp
	case s.value:
		s.scanValue()
		s.value = false
		return pos, tokenValue, s.input[pos:s.ppos]
	case isFieldRune(ch):
		s.scanField()
		return pos, tokenField, s.input[pos:s.ppos]
	***REMOVED***

	return s.pos, token(ch), ""
***REMOVED***

func (s *scanner) scanField() ***REMOVED***
	for ***REMOVED***
		ch := s.peek()
		if !isFieldRune(ch) ***REMOVED***
			break
		***REMOVED***
		s.next()
	***REMOVED***
***REMOVED***

func (s *scanner) scanOperator() ***REMOVED***
	for ***REMOVED***
		ch := s.peek()
		switch ch ***REMOVED***
		case '=', '!', '~':
			s.next()
		default:
			return
		***REMOVED***
	***REMOVED***
***REMOVED***

func (s *scanner) scanValue() ***REMOVED***
	for ***REMOVED***
		ch := s.peek()
		if !isValueRune(ch) ***REMOVED***
			break
		***REMOVED***
		s.next()
	***REMOVED***
***REMOVED***

func (s *scanner) scanQuoted(quote rune) ***REMOVED***
	ch := s.next() // read character after quote
	for ch != quote ***REMOVED***
		if ch == '\n' || ch < 0 ***REMOVED***
			s.error("literal not terminated")
			return
		***REMOVED***
		if ch == '\\' ***REMOVED***
			ch = s.scanEscape(quote)
		***REMOVED*** else ***REMOVED***
			ch = s.next()
		***REMOVED***
	***REMOVED***
	return
***REMOVED***

func (s *scanner) scanEscape(quote rune) rune ***REMOVED***
	ch := s.next() // read character after '/'
	switch ch ***REMOVED***
	case 'a', 'b', 'f', 'n', 'r', 't', 'v', '\\', quote:
		// nothing to do
		ch = s.next()
	case '0', '1', '2', '3', '4', '5', '6', '7':
		ch = s.scanDigits(ch, 8, 3)
	case 'x':
		ch = s.scanDigits(s.next(), 16, 2)
	case 'u':
		ch = s.scanDigits(s.next(), 16, 4)
	case 'U':
		ch = s.scanDigits(s.next(), 16, 8)
	default:
		s.error("illegal char escape")
	***REMOVED***
	return ch
***REMOVED***

func (s *scanner) scanDigits(ch rune, base, n int) rune ***REMOVED***
	for n > 0 && digitVal(ch) < base ***REMOVED***
		ch = s.next()
		n--
	***REMOVED***
	if n > 0 ***REMOVED***
		s.error("illegal char escape")
	***REMOVED***
	return ch
***REMOVED***

func (s *scanner) error(msg string) ***REMOVED***
	fmt.Println("error fixme", msg)
***REMOVED***

func digitVal(ch rune) int ***REMOVED***
	switch ***REMOVED***
	case '0' <= ch && ch <= '9':
		return int(ch - '0')
	case 'a' <= ch && ch <= 'f':
		return int(ch - 'a' + 10)
	case 'A' <= ch && ch <= 'F':
		return int(ch - 'A' + 10)
	***REMOVED***
	return 16 // larger than any legal digit val
***REMOVED***

func isFieldRune(r rune) bool ***REMOVED***
	return (r == '_' || isAlphaRune(r) || isDigitRune(r))
***REMOVED***

func isAlphaRune(r rune) bool ***REMOVED***
	return r >= 'A' && r <= 'Z' || r >= 'a' && r <= 'z'
***REMOVED***

func isDigitRune(r rune) bool ***REMOVED***
	return r >= '0' && r <= '9'
***REMOVED***

func isOperatorRune(r rune) bool ***REMOVED***
	switch r ***REMOVED***
	case '=', '!', '~':
		return true
	***REMOVED***

	return false
***REMOVED***

func isQuoteRune(r rune) bool ***REMOVED***
	switch r ***REMOVED***
	case '/', '|', '"': // maybe add single quoting?
		return true
	***REMOVED***

	return false
***REMOVED***

func isSeparatorRune(r rune) bool ***REMOVED***
	switch r ***REMOVED***
	case ',', '.':
		return true
	***REMOVED***

	return false
***REMOVED***

func isValueRune(r rune) bool ***REMOVED***
	return r != ',' && !unicode.IsSpace(r) &&
		(unicode.IsLetter(r) ||
			unicode.IsDigit(r) ||
			unicode.IsNumber(r) ||
			unicode.IsGraphic(r) ||
			unicode.IsPunct(r))
***REMOVED***
