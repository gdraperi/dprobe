package dbus

import (
	"fmt"
	"strings"
	"unicode"
	"unicode/utf8"
)

// Heavily inspired by the lexer from text/template.

type varToken struct ***REMOVED***
	typ varTokenType
	val string
***REMOVED***

type varTokenType byte

const (
	tokEOF varTokenType = iota
	tokError
	tokNumber
	tokString
	tokBool
	tokArrayStart
	tokArrayEnd
	tokDictStart
	tokDictEnd
	tokVariantStart
	tokVariantEnd
	tokComma
	tokColon
	tokType
	tokByteString
)

type varLexer struct ***REMOVED***
	input  string
	start  int
	pos    int
	width  int
	tokens []varToken
***REMOVED***

type lexState func(*varLexer) lexState

func varLex(s string) []varToken ***REMOVED***
	l := &varLexer***REMOVED***input: s***REMOVED***
	l.run()
	return l.tokens
***REMOVED***

func (l *varLexer) accept(valid string) bool ***REMOVED***
	if strings.IndexRune(valid, l.next()) >= 0 ***REMOVED***
		return true
	***REMOVED***
	l.backup()
	return false
***REMOVED***

func (l *varLexer) backup() ***REMOVED***
	l.pos -= l.width
***REMOVED***

func (l *varLexer) emit(t varTokenType) ***REMOVED***
	l.tokens = append(l.tokens, varToken***REMOVED***t, l.input[l.start:l.pos]***REMOVED***)
	l.start = l.pos
***REMOVED***

func (l *varLexer) errorf(format string, v ...interface***REMOVED******REMOVED***) lexState ***REMOVED***
	l.tokens = append(l.tokens, varToken***REMOVED***
		tokError,
		fmt.Sprintf(format, v...),
	***REMOVED***)
	return nil
***REMOVED***

func (l *varLexer) ignore() ***REMOVED***
	l.start = l.pos
***REMOVED***

func (l *varLexer) next() rune ***REMOVED***
	var r rune

	if l.pos >= len(l.input) ***REMOVED***
		l.width = 0
		return -1
	***REMOVED***
	r, l.width = utf8.DecodeRuneInString(l.input[l.pos:])
	l.pos += l.width
	return r
***REMOVED***

func (l *varLexer) run() ***REMOVED***
	for state := varLexNormal; state != nil; ***REMOVED***
		state = state(l)
	***REMOVED***
***REMOVED***

func (l *varLexer) peek() rune ***REMOVED***
	r := l.next()
	l.backup()
	return r
***REMOVED***

func varLexNormal(l *varLexer) lexState ***REMOVED***
	for ***REMOVED***
		r := l.next()
		switch ***REMOVED***
		case r == -1:
			l.emit(tokEOF)
			return nil
		case r == '[':
			l.emit(tokArrayStart)
		case r == ']':
			l.emit(tokArrayEnd)
		case r == '***REMOVED***':
			l.emit(tokDictStart)
		case r == '***REMOVED***':
			l.emit(tokDictEnd)
		case r == '<':
			l.emit(tokVariantStart)
		case r == '>':
			l.emit(tokVariantEnd)
		case r == ':':
			l.emit(tokColon)
		case r == ',':
			l.emit(tokComma)
		case r == '\'' || r == '"':
			l.backup()
			return varLexString
		case r == '@':
			l.backup()
			return varLexType
		case unicode.IsSpace(r):
			l.ignore()
		case unicode.IsNumber(r) || r == '+' || r == '-':
			l.backup()
			return varLexNumber
		case r == 'b':
			pos := l.start
			if n := l.peek(); n == '"' || n == '\'' ***REMOVED***
				return varLexByteString
			***REMOVED***
			// not a byte string; try to parse it as a type or bool below
			l.pos = pos + 1
			l.width = 1
			fallthrough
		default:
			// either a bool or a type. Try bools first.
			l.backup()
			if l.pos+4 <= len(l.input) ***REMOVED***
				if l.input[l.pos:l.pos+4] == "true" ***REMOVED***
					l.pos += 4
					l.emit(tokBool)
					continue
				***REMOVED***
			***REMOVED***
			if l.pos+5 <= len(l.input) ***REMOVED***
				if l.input[l.pos:l.pos+5] == "false" ***REMOVED***
					l.pos += 5
					l.emit(tokBool)
					continue
				***REMOVED***
			***REMOVED***
			// must be a type.
			return varLexType
		***REMOVED***
	***REMOVED***
***REMOVED***

var varTypeMap = map[string]string***REMOVED***
	"boolean":    "b",
	"byte":       "y",
	"int16":      "n",
	"uint16":     "q",
	"int32":      "i",
	"uint32":     "u",
	"int64":      "x",
	"uint64":     "t",
	"double":     "f",
	"string":     "s",
	"objectpath": "o",
	"signature":  "g",
***REMOVED***

func varLexByteString(l *varLexer) lexState ***REMOVED***
	q := l.next()
Loop:
	for ***REMOVED***
		switch l.next() ***REMOVED***
		case '\\':
			if r := l.next(); r != -1 ***REMOVED***
				break
			***REMOVED***
			fallthrough
		case -1:
			return l.errorf("unterminated bytestring")
		case q:
			break Loop
		***REMOVED***
	***REMOVED***
	l.emit(tokByteString)
	return varLexNormal
***REMOVED***

func varLexNumber(l *varLexer) lexState ***REMOVED***
	l.accept("+-")
	digits := "0123456789"
	if l.accept("0") ***REMOVED***
		if l.accept("x") ***REMOVED***
			digits = "0123456789abcdefABCDEF"
		***REMOVED*** else ***REMOVED***
			digits = "01234567"
		***REMOVED***
	***REMOVED***
	for strings.IndexRune(digits, l.next()) >= 0 ***REMOVED***
	***REMOVED***
	l.backup()
	if l.accept(".") ***REMOVED***
		for strings.IndexRune(digits, l.next()) >= 0 ***REMOVED***
		***REMOVED***
		l.backup()
	***REMOVED***
	if l.accept("eE") ***REMOVED***
		l.accept("+-")
		for strings.IndexRune("0123456789", l.next()) >= 0 ***REMOVED***
		***REMOVED***
		l.backup()
	***REMOVED***
	if r := l.peek(); unicode.IsLetter(r) ***REMOVED***
		l.next()
		return l.errorf("bad number syntax: %q", l.input[l.start:l.pos])
	***REMOVED***
	l.emit(tokNumber)
	return varLexNormal
***REMOVED***

func varLexString(l *varLexer) lexState ***REMOVED***
	q := l.next()
Loop:
	for ***REMOVED***
		switch l.next() ***REMOVED***
		case '\\':
			if r := l.next(); r != -1 ***REMOVED***
				break
			***REMOVED***
			fallthrough
		case -1:
			return l.errorf("unterminated string")
		case q:
			break Loop
		***REMOVED***
	***REMOVED***
	l.emit(tokString)
	return varLexNormal
***REMOVED***

func varLexType(l *varLexer) lexState ***REMOVED***
	at := l.accept("@")
	for ***REMOVED***
		r := l.next()
		if r == -1 ***REMOVED***
			break
		***REMOVED***
		if unicode.IsSpace(r) ***REMOVED***
			l.backup()
			break
		***REMOVED***
	***REMOVED***
	if at ***REMOVED***
		if _, err := ParseSignature(l.input[l.start+1 : l.pos]); err != nil ***REMOVED***
			return l.errorf("%s", err)
		***REMOVED***
	***REMOVED*** else ***REMOVED***
		if _, ok := varTypeMap[l.input[l.start:l.pos]]; ok ***REMOVED***
			l.emit(tokType)
			return varLexNormal
		***REMOVED***
		return l.errorf("unrecognized type %q", l.input[l.start:l.pos])
	***REMOVED***
	l.emit(tokType)
	return varLexNormal
***REMOVED***
