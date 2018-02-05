// TOML JSONPath lexer.
//
// Written using the principles developed by Rob Pike in
// http://www.youtube.com/watch?v=HxaD_trXwRE

package query

import (
	"fmt"
	"github.com/pelletier/go-toml"
	"strconv"
	"strings"
	"unicode/utf8"
)

// Lexer state function
type queryLexStateFn func() queryLexStateFn

// Lexer definition
type queryLexer struct ***REMOVED***
	input      string
	start      int
	pos        int
	width      int
	tokens     chan token
	depth      int
	line       int
	col        int
	stringTerm string
***REMOVED***

func (l *queryLexer) run() ***REMOVED***
	for state := l.lexVoid; state != nil; ***REMOVED***
		state = state()
	***REMOVED***
	close(l.tokens)
***REMOVED***

func (l *queryLexer) nextStart() ***REMOVED***
	// iterate by runes (utf8 characters)
	// search for newlines and advance line/col counts
	for i := l.start; i < l.pos; ***REMOVED***
		r, width := utf8.DecodeRuneInString(l.input[i:])
		if r == '\n' ***REMOVED***
			l.line++
			l.col = 1
		***REMOVED*** else ***REMOVED***
			l.col++
		***REMOVED***
		i += width
	***REMOVED***
	// advance start position to next token
	l.start = l.pos
***REMOVED***

func (l *queryLexer) emit(t tokenType) ***REMOVED***
	l.tokens <- token***REMOVED***
		Position: toml.Position***REMOVED***Line: l.line, Col: l.col***REMOVED***,
		typ:      t,
		val:      l.input[l.start:l.pos],
	***REMOVED***
	l.nextStart()
***REMOVED***

func (l *queryLexer) emitWithValue(t tokenType, value string) ***REMOVED***
	l.tokens <- token***REMOVED***
		Position: toml.Position***REMOVED***Line: l.line, Col: l.col***REMOVED***,
		typ:      t,
		val:      value,
	***REMOVED***
	l.nextStart()
***REMOVED***

func (l *queryLexer) next() rune ***REMOVED***
	if l.pos >= len(l.input) ***REMOVED***
		l.width = 0
		return eof
	***REMOVED***
	var r rune
	r, l.width = utf8.DecodeRuneInString(l.input[l.pos:])
	l.pos += l.width
	return r
***REMOVED***

func (l *queryLexer) ignore() ***REMOVED***
	l.nextStart()
***REMOVED***

func (l *queryLexer) backup() ***REMOVED***
	l.pos -= l.width
***REMOVED***

func (l *queryLexer) errorf(format string, args ...interface***REMOVED******REMOVED***) queryLexStateFn ***REMOVED***
	l.tokens <- token***REMOVED***
		Position: toml.Position***REMOVED***Line: l.line, Col: l.col***REMOVED***,
		typ:      tokenError,
		val:      fmt.Sprintf(format, args...),
	***REMOVED***
	return nil
***REMOVED***

func (l *queryLexer) peek() rune ***REMOVED***
	r := l.next()
	l.backup()
	return r
***REMOVED***

func (l *queryLexer) accept(valid string) bool ***REMOVED***
	if strings.ContainsRune(valid, l.next()) ***REMOVED***
		return true
	***REMOVED***
	l.backup()
	return false
***REMOVED***

func (l *queryLexer) follow(next string) bool ***REMOVED***
	return strings.HasPrefix(l.input[l.pos:], next)
***REMOVED***

func (l *queryLexer) lexVoid() queryLexStateFn ***REMOVED***
	for ***REMOVED***
		next := l.peek()
		switch next ***REMOVED***
		case '$':
			l.pos++
			l.emit(tokenDollar)
			continue
		case '.':
			if l.follow("..") ***REMOVED***
				l.pos += 2
				l.emit(tokenDotDot)
			***REMOVED*** else ***REMOVED***
				l.pos++
				l.emit(tokenDot)
			***REMOVED***
			continue
		case '[':
			l.pos++
			l.emit(tokenLeftBracket)
			continue
		case ']':
			l.pos++
			l.emit(tokenRightBracket)
			continue
		case ',':
			l.pos++
			l.emit(tokenComma)
			continue
		case '*':
			l.pos++
			l.emit(tokenStar)
			continue
		case '(':
			l.pos++
			l.emit(tokenLeftParen)
			continue
		case ')':
			l.pos++
			l.emit(tokenRightParen)
			continue
		case '?':
			l.pos++
			l.emit(tokenQuestion)
			continue
		case ':':
			l.pos++
			l.emit(tokenColon)
			continue
		case '\'':
			l.ignore()
			l.stringTerm = string(next)
			return l.lexString
		case '"':
			l.ignore()
			l.stringTerm = string(next)
			return l.lexString
		***REMOVED***

		if isSpace(next) ***REMOVED***
			l.next()
			l.ignore()
			continue
		***REMOVED***

		if isAlphanumeric(next) ***REMOVED***
			return l.lexKey
		***REMOVED***

		if next == '+' || next == '-' || isDigit(next) ***REMOVED***
			return l.lexNumber
		***REMOVED***

		if l.next() == eof ***REMOVED***
			break
		***REMOVED***

		return l.errorf("unexpected char: '%v'", next)
	***REMOVED***
	l.emit(tokenEOF)
	return nil
***REMOVED***

func (l *queryLexer) lexKey() queryLexStateFn ***REMOVED***
	for ***REMOVED***
		next := l.peek()
		if !isAlphanumeric(next) ***REMOVED***
			l.emit(tokenKey)
			return l.lexVoid
		***REMOVED***

		if l.next() == eof ***REMOVED***
			break
		***REMOVED***
	***REMOVED***
	l.emit(tokenEOF)
	return nil
***REMOVED***

func (l *queryLexer) lexString() queryLexStateFn ***REMOVED***
	l.pos++
	l.ignore()
	growingString := ""

	for ***REMOVED***
		if l.follow(l.stringTerm) ***REMOVED***
			l.emitWithValue(tokenString, growingString)
			l.pos++
			l.ignore()
			return l.lexVoid
		***REMOVED***

		if l.follow("\\\"") ***REMOVED***
			l.pos++
			growingString += "\""
		***REMOVED*** else if l.follow("\\'") ***REMOVED***
			l.pos++
			growingString += "'"
		***REMOVED*** else if l.follow("\\n") ***REMOVED***
			l.pos++
			growingString += "\n"
		***REMOVED*** else if l.follow("\\b") ***REMOVED***
			l.pos++
			growingString += "\b"
		***REMOVED*** else if l.follow("\\f") ***REMOVED***
			l.pos++
			growingString += "\f"
		***REMOVED*** else if l.follow("\\/") ***REMOVED***
			l.pos++
			growingString += "/"
		***REMOVED*** else if l.follow("\\t") ***REMOVED***
			l.pos++
			growingString += "\t"
		***REMOVED*** else if l.follow("\\r") ***REMOVED***
			l.pos++
			growingString += "\r"
		***REMOVED*** else if l.follow("\\\\") ***REMOVED***
			l.pos++
			growingString += "\\"
		***REMOVED*** else if l.follow("\\u") ***REMOVED***
			l.pos += 2
			code := ""
			for i := 0; i < 4; i++ ***REMOVED***
				c := l.peek()
				l.pos++
				if !isHexDigit(c) ***REMOVED***
					return l.errorf("unfinished unicode escape")
				***REMOVED***
				code = code + string(c)
			***REMOVED***
			l.pos--
			intcode, err := strconv.ParseInt(code, 16, 32)
			if err != nil ***REMOVED***
				return l.errorf("invalid unicode escape: \\u" + code)
			***REMOVED***
			growingString += string(rune(intcode))
		***REMOVED*** else if l.follow("\\U") ***REMOVED***
			l.pos += 2
			code := ""
			for i := 0; i < 8; i++ ***REMOVED***
				c := l.peek()
				l.pos++
				if !isHexDigit(c) ***REMOVED***
					return l.errorf("unfinished unicode escape")
				***REMOVED***
				code = code + string(c)
			***REMOVED***
			l.pos--
			intcode, err := strconv.ParseInt(code, 16, 32)
			if err != nil ***REMOVED***
				return l.errorf("invalid unicode escape: \\u" + code)
			***REMOVED***
			growingString += string(rune(intcode))
		***REMOVED*** else if l.follow("\\") ***REMOVED***
			l.pos++
			return l.errorf("invalid escape sequence: \\" + string(l.peek()))
		***REMOVED*** else ***REMOVED***
			growingString += string(l.peek())
		***REMOVED***

		if l.next() == eof ***REMOVED***
			break
		***REMOVED***
	***REMOVED***

	return l.errorf("unclosed string")
***REMOVED***

func (l *queryLexer) lexNumber() queryLexStateFn ***REMOVED***
	l.ignore()
	if !l.accept("+") ***REMOVED***
		l.accept("-")
	***REMOVED***
	pointSeen := false
	digitSeen := false
	for ***REMOVED***
		next := l.next()
		if next == '.' ***REMOVED***
			if pointSeen ***REMOVED***
				return l.errorf("cannot have two dots in one float")
			***REMOVED***
			if !isDigit(l.peek()) ***REMOVED***
				return l.errorf("float cannot end with a dot")
			***REMOVED***
			pointSeen = true
		***REMOVED*** else if isDigit(next) ***REMOVED***
			digitSeen = true
		***REMOVED*** else ***REMOVED***
			l.backup()
			break
		***REMOVED***
		if pointSeen && !digitSeen ***REMOVED***
			return l.errorf("cannot start float with a dot")
		***REMOVED***
	***REMOVED***

	if !digitSeen ***REMOVED***
		return l.errorf("no digit in that number")
	***REMOVED***
	if pointSeen ***REMOVED***
		l.emit(tokenFloat)
	***REMOVED*** else ***REMOVED***
		l.emit(tokenInteger)
	***REMOVED***
	return l.lexVoid
***REMOVED***

// Entry point
func lexQuery(input string) chan token ***REMOVED***
	l := &queryLexer***REMOVED***
		input:  input,
		tokens: make(chan token),
		line:   1,
		col:    1,
	***REMOVED***
	go l.run()
	return l.tokens
***REMOVED***
