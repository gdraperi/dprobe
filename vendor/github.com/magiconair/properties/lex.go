// Copyright 2017 Frank Schroeder. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//
// Parts of the lexer are from the template/text/parser package
// For these parts the following applies:
//
// Copyright 2011 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file of the go 1.2
// distribution.

package properties

import (
	"fmt"
	"strconv"
	"strings"
	"unicode/utf8"
)

// item represents a token or text string returned from the scanner.
type item struct ***REMOVED***
	typ itemType // The type of this item.
	pos int      // The starting position, in bytes, of this item in the input string.
	val string   // The value of this item.
***REMOVED***

func (i item) String() string ***REMOVED***
	switch ***REMOVED***
	case i.typ == itemEOF:
		return "EOF"
	case i.typ == itemError:
		return i.val
	case len(i.val) > 10:
		return fmt.Sprintf("%.10q...", i.val)
	***REMOVED***
	return fmt.Sprintf("%q", i.val)
***REMOVED***

// itemType identifies the type of lex items.
type itemType int

const (
	itemError itemType = iota // error occurred; value is text of error
	itemEOF
	itemKey     // a key
	itemValue   // a value
	itemComment // a comment
)

// defines a constant for EOF
const eof = -1

// permitted whitespace characters space, FF and TAB
const whitespace = " \f\t"

// stateFn represents the state of the scanner as a function that returns the next state.
type stateFn func(*lexer) stateFn

// lexer holds the state of the scanner.
type lexer struct ***REMOVED***
	input   string    // the string being scanned
	state   stateFn   // the next lexing function to enter
	pos     int       // current position in the input
	start   int       // start position of this item
	width   int       // width of last rune read from input
	lastPos int       // position of most recent item returned by nextItem
	runes   []rune    // scanned runes for this item
	items   chan item // channel of scanned items
***REMOVED***

// next returns the next rune in the input.
func (l *lexer) next() rune ***REMOVED***
	if l.pos >= len(l.input) ***REMOVED***
		l.width = 0
		return eof
	***REMOVED***
	r, w := utf8.DecodeRuneInString(l.input[l.pos:])
	l.width = w
	l.pos += l.width
	return r
***REMOVED***

// peek returns but does not consume the next rune in the input.
func (l *lexer) peek() rune ***REMOVED***
	r := l.next()
	l.backup()
	return r
***REMOVED***

// backup steps back one rune. Can only be called once per call of next.
func (l *lexer) backup() ***REMOVED***
	l.pos -= l.width
***REMOVED***

// emit passes an item back to the client.
func (l *lexer) emit(t itemType) ***REMOVED***
	i := item***REMOVED***t, l.start, string(l.runes)***REMOVED***
	l.items <- i
	l.start = l.pos
	l.runes = l.runes[:0]
***REMOVED***

// ignore skips over the pending input before this point.
func (l *lexer) ignore() ***REMOVED***
	l.start = l.pos
***REMOVED***

// appends the rune to the current value
func (l *lexer) appendRune(r rune) ***REMOVED***
	l.runes = append(l.runes, r)
***REMOVED***

// accept consumes the next rune if it's from the valid set.
func (l *lexer) accept(valid string) bool ***REMOVED***
	if strings.ContainsRune(valid, l.next()) ***REMOVED***
		return true
	***REMOVED***
	l.backup()
	return false
***REMOVED***

// acceptRun consumes a run of runes from the valid set.
func (l *lexer) acceptRun(valid string) ***REMOVED***
	for strings.ContainsRune(valid, l.next()) ***REMOVED***
	***REMOVED***
	l.backup()
***REMOVED***

// acceptRunUntil consumes a run of runes up to a terminator.
func (l *lexer) acceptRunUntil(term rune) ***REMOVED***
	for term != l.next() ***REMOVED***
	***REMOVED***
	l.backup()
***REMOVED***

// hasText returns true if the current parsed text is not empty.
func (l *lexer) isNotEmpty() bool ***REMOVED***
	return l.pos > l.start
***REMOVED***

// lineNumber reports which line we're on, based on the position of
// the previous item returned by nextItem. Doing it this way
// means we don't have to worry about peek double counting.
func (l *lexer) lineNumber() int ***REMOVED***
	return 1 + strings.Count(l.input[:l.lastPos], "\n")
***REMOVED***

// errorf returns an error token and terminates the scan by passing
// back a nil pointer that will be the next state, terminating l.nextItem.
func (l *lexer) errorf(format string, args ...interface***REMOVED******REMOVED***) stateFn ***REMOVED***
	l.items <- item***REMOVED***itemError, l.start, fmt.Sprintf(format, args...)***REMOVED***
	return nil
***REMOVED***

// nextItem returns the next item from the input.
func (l *lexer) nextItem() item ***REMOVED***
	i := <-l.items
	l.lastPos = i.pos
	return i
***REMOVED***

// lex creates a new scanner for the input string.
func lex(input string) *lexer ***REMOVED***
	l := &lexer***REMOVED***
		input: input,
		items: make(chan item),
		runes: make([]rune, 0, 32),
	***REMOVED***
	go l.run()
	return l
***REMOVED***

// run runs the state machine for the lexer.
func (l *lexer) run() ***REMOVED***
	for l.state = lexBeforeKey(l); l.state != nil; ***REMOVED***
		l.state = l.state(l)
	***REMOVED***
***REMOVED***

// state functions

// lexBeforeKey scans until a key begins.
func lexBeforeKey(l *lexer) stateFn ***REMOVED***
	switch r := l.next(); ***REMOVED***
	case isEOF(r):
		l.emit(itemEOF)
		return nil

	case isEOL(r):
		l.ignore()
		return lexBeforeKey

	case isComment(r):
		return lexComment

	case isWhitespace(r):
		l.ignore()
		return lexBeforeKey

	default:
		l.backup()
		return lexKey
	***REMOVED***
***REMOVED***

// lexComment scans a comment line. The comment character has already been scanned.
func lexComment(l *lexer) stateFn ***REMOVED***
	l.acceptRun(whitespace)
	l.ignore()
	for ***REMOVED***
		switch r := l.next(); ***REMOVED***
		case isEOF(r):
			l.ignore()
			l.emit(itemEOF)
			return nil
		case isEOL(r):
			l.emit(itemComment)
			return lexBeforeKey
		default:
			l.appendRune(r)
		***REMOVED***
	***REMOVED***
***REMOVED***

// lexKey scans the key up to a delimiter
func lexKey(l *lexer) stateFn ***REMOVED***
	var r rune

Loop:
	for ***REMOVED***
		switch r = l.next(); ***REMOVED***

		case isEscape(r):
			err := l.scanEscapeSequence()
			if err != nil ***REMOVED***
				return l.errorf(err.Error())
			***REMOVED***

		case isEndOfKey(r):
			l.backup()
			break Loop

		case isEOF(r):
			break Loop

		default:
			l.appendRune(r)
		***REMOVED***
	***REMOVED***

	if len(l.runes) > 0 ***REMOVED***
		l.emit(itemKey)
	***REMOVED***

	if isEOF(r) ***REMOVED***
		l.emit(itemEOF)
		return nil
	***REMOVED***

	return lexBeforeValue
***REMOVED***

// lexBeforeValue scans the delimiter between key and value.
// Leading and trailing whitespace is ignored.
// We expect to be just after the key.
func lexBeforeValue(l *lexer) stateFn ***REMOVED***
	l.acceptRun(whitespace)
	l.accept(":=")
	l.acceptRun(whitespace)
	l.ignore()
	return lexValue
***REMOVED***

// lexValue scans text until the end of the line. We expect to be just after the delimiter.
func lexValue(l *lexer) stateFn ***REMOVED***
	for ***REMOVED***
		switch r := l.next(); ***REMOVED***
		case isEscape(r):
			if isEOL(l.peek()) ***REMOVED***
				l.next()
				l.acceptRun(whitespace)
			***REMOVED*** else ***REMOVED***
				err := l.scanEscapeSequence()
				if err != nil ***REMOVED***
					return l.errorf(err.Error())
				***REMOVED***
			***REMOVED***

		case isEOL(r):
			l.emit(itemValue)
			l.ignore()
			return lexBeforeKey

		case isEOF(r):
			l.emit(itemValue)
			l.emit(itemEOF)
			return nil

		default:
			l.appendRune(r)
		***REMOVED***
	***REMOVED***
***REMOVED***

// scanEscapeSequence scans either one of the escaped characters
// or a unicode literal. We expect to be after the escape character.
func (l *lexer) scanEscapeSequence() error ***REMOVED***
	switch r := l.next(); ***REMOVED***

	case isEscapedCharacter(r):
		l.appendRune(decodeEscapedCharacter(r))
		return nil

	case atUnicodeLiteral(r):
		return l.scanUnicodeLiteral()

	case isEOF(r):
		return fmt.Errorf("premature EOF")

	// silently drop the escape character and append the rune as is
	default:
		l.appendRune(r)
		return nil
	***REMOVED***
***REMOVED***

// scans a unicode literal in the form \uXXXX. We expect to be after the \u.
func (l *lexer) scanUnicodeLiteral() error ***REMOVED***
	// scan the digits
	d := make([]rune, 4)
	for i := 0; i < 4; i++ ***REMOVED***
		d[i] = l.next()
		if d[i] == eof || !strings.ContainsRune("0123456789abcdefABCDEF", d[i]) ***REMOVED***
			return fmt.Errorf("invalid unicode literal")
		***REMOVED***
	***REMOVED***

	// decode the digits into a rune
	r, err := strconv.ParseInt(string(d), 16, 0)
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	l.appendRune(rune(r))
	return nil
***REMOVED***

// decodeEscapedCharacter returns the unescaped rune. We expect to be after the escape character.
func decodeEscapedCharacter(r rune) rune ***REMOVED***
	switch r ***REMOVED***
	case 'f':
		return '\f'
	case 'n':
		return '\n'
	case 'r':
		return '\r'
	case 't':
		return '\t'
	default:
		return r
	***REMOVED***
***REMOVED***

// atUnicodeLiteral reports whether we are at a unicode literal.
// The escape character has already been consumed.
func atUnicodeLiteral(r rune) bool ***REMOVED***
	return r == 'u'
***REMOVED***

// isComment reports whether we are at the start of a comment.
func isComment(r rune) bool ***REMOVED***
	return r == '#' || r == '!'
***REMOVED***

// isEndOfKey reports whether the rune terminates the current key.
func isEndOfKey(r rune) bool ***REMOVED***
	return strings.ContainsRune(" \f\t\r\n:=", r)
***REMOVED***

// isEOF reports whether we are at EOF.
func isEOF(r rune) bool ***REMOVED***
	return r == eof
***REMOVED***

// isEOL reports whether we are at a new line character.
func isEOL(r rune) bool ***REMOVED***
	return r == '\n' || r == '\r'
***REMOVED***

// isEscape reports whether the rune is the escape character which
// prefixes unicode literals and other escaped characters.
func isEscape(r rune) bool ***REMOVED***
	return r == '\\'
***REMOVED***

// isEscapedCharacter reports whether we are at one of the characters that need escaping.
// The escape character has already been consumed.
func isEscapedCharacter(r rune) bool ***REMOVED***
	return strings.ContainsRune(" :=fnrt", r)
***REMOVED***

// isWhitespace reports whether the rune is a whitespace character.
func isWhitespace(r rune) bool ***REMOVED***
	return strings.ContainsRune(whitespace, r)
***REMOVED***
