// TOML lexer.
//
// Written using the principles developed by Rob Pike in
// http://www.youtube.com/watch?v=HxaD_trXwRE

package toml

import (
	"bytes"
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

var dateRegexp *regexp.Regexp

// Define state functions
type tomlLexStateFn func() tomlLexStateFn

// Define lexer
type tomlLexer struct ***REMOVED***
	inputIdx          int
	input             []rune // Textual source
	currentTokenStart int
	currentTokenStop  int
	tokens            []token
	depth             int
	line              int
	col               int
	endbufferLine     int
	endbufferCol      int
***REMOVED***

// Basic read operations on input

func (l *tomlLexer) read() rune ***REMOVED***
	r := l.peek()
	if r == '\n' ***REMOVED***
		l.endbufferLine++
		l.endbufferCol = 1
	***REMOVED*** else ***REMOVED***
		l.endbufferCol++
	***REMOVED***
	l.inputIdx++
	return r
***REMOVED***

func (l *tomlLexer) next() rune ***REMOVED***
	r := l.read()

	if r != eof ***REMOVED***
		l.currentTokenStop++
	***REMOVED***
	return r
***REMOVED***

func (l *tomlLexer) ignore() ***REMOVED***
	l.currentTokenStart = l.currentTokenStop
	l.line = l.endbufferLine
	l.col = l.endbufferCol
***REMOVED***

func (l *tomlLexer) skip() ***REMOVED***
	l.next()
	l.ignore()
***REMOVED***

func (l *tomlLexer) fastForward(n int) ***REMOVED***
	for i := 0; i < n; i++ ***REMOVED***
		l.next()
	***REMOVED***
***REMOVED***

func (l *tomlLexer) emitWithValue(t tokenType, value string) ***REMOVED***
	l.tokens = append(l.tokens, token***REMOVED***
		Position: Position***REMOVED***l.line, l.col***REMOVED***,
		typ:      t,
		val:      value,
	***REMOVED***)
	l.ignore()
***REMOVED***

func (l *tomlLexer) emit(t tokenType) ***REMOVED***
	l.emitWithValue(t, string(l.input[l.currentTokenStart:l.currentTokenStop]))
***REMOVED***

func (l *tomlLexer) peek() rune ***REMOVED***
	if l.inputIdx >= len(l.input) ***REMOVED***
		return eof
	***REMOVED***
	return l.input[l.inputIdx]
***REMOVED***

func (l *tomlLexer) peekString(size int) string ***REMOVED***
	maxIdx := len(l.input)
	upperIdx := l.inputIdx + size // FIXME: potential overflow
	if upperIdx > maxIdx ***REMOVED***
		upperIdx = maxIdx
	***REMOVED***
	return string(l.input[l.inputIdx:upperIdx])
***REMOVED***

func (l *tomlLexer) follow(next string) bool ***REMOVED***
	return next == l.peekString(len(next))
***REMOVED***

// Error management

func (l *tomlLexer) errorf(format string, args ...interface***REMOVED******REMOVED***) tomlLexStateFn ***REMOVED***
	l.tokens = append(l.tokens, token***REMOVED***
		Position: Position***REMOVED***l.line, l.col***REMOVED***,
		typ:      tokenError,
		val:      fmt.Sprintf(format, args...),
	***REMOVED***)
	return nil
***REMOVED***

// State functions

func (l *tomlLexer) lexVoid() tomlLexStateFn ***REMOVED***
	for ***REMOVED***
		next := l.peek()
		switch next ***REMOVED***
		case '[':
			return l.lexTableKey
		case '#':
			return l.lexComment(l.lexVoid)
		case '=':
			return l.lexEqual
		case '\r':
			fallthrough
		case '\n':
			l.skip()
			continue
		***REMOVED***

		if isSpace(next) ***REMOVED***
			l.skip()
		***REMOVED***

		if l.depth > 0 ***REMOVED***
			return l.lexRvalue
		***REMOVED***

		if isKeyStartChar(next) ***REMOVED***
			return l.lexKey
		***REMOVED***

		if next == eof ***REMOVED***
			l.next()
			break
		***REMOVED***
	***REMOVED***

	l.emit(tokenEOF)
	return nil
***REMOVED***

func (l *tomlLexer) lexRvalue() tomlLexStateFn ***REMOVED***
	for ***REMOVED***
		next := l.peek()
		switch next ***REMOVED***
		case '.':
			return l.errorf("cannot start float with a dot")
		case '=':
			return l.lexEqual
		case '[':
			l.depth++
			return l.lexLeftBracket
		case ']':
			l.depth--
			return l.lexRightBracket
		case '***REMOVED***':
			return l.lexLeftCurlyBrace
		case '***REMOVED***':
			return l.lexRightCurlyBrace
		case '#':
			return l.lexComment(l.lexRvalue)
		case '"':
			return l.lexString
		case '\'':
			return l.lexLiteralString
		case ',':
			return l.lexComma
		case '\r':
			fallthrough
		case '\n':
			l.skip()
			if l.depth == 0 ***REMOVED***
				return l.lexVoid
			***REMOVED***
			return l.lexRvalue
		case '_':
			return l.errorf("cannot start number with underscore")
		***REMOVED***

		if l.follow("true") ***REMOVED***
			return l.lexTrue
		***REMOVED***

		if l.follow("false") ***REMOVED***
			return l.lexFalse
		***REMOVED***

		if l.follow("inf") ***REMOVED***
			return l.lexInf
		***REMOVED***

		if l.follow("nan") ***REMOVED***
			return l.lexNan
		***REMOVED***

		if isSpace(next) ***REMOVED***
			l.skip()
			continue
		***REMOVED***

		if next == eof ***REMOVED***
			l.next()
			break
		***REMOVED***

		possibleDate := l.peekString(35)
		dateMatch := dateRegexp.FindString(possibleDate)
		if dateMatch != "" ***REMOVED***
			l.fastForward(len(dateMatch))
			return l.lexDate
		***REMOVED***

		if next == '+' || next == '-' || isDigit(next) ***REMOVED***
			return l.lexNumber
		***REMOVED***

		if isAlphanumeric(next) ***REMOVED***
			return l.lexKey
		***REMOVED***

		return l.errorf("no value can start with %c", next)
	***REMOVED***

	l.emit(tokenEOF)
	return nil
***REMOVED***

func (l *tomlLexer) lexLeftCurlyBrace() tomlLexStateFn ***REMOVED***
	l.next()
	l.emit(tokenLeftCurlyBrace)
	return l.lexRvalue
***REMOVED***

func (l *tomlLexer) lexRightCurlyBrace() tomlLexStateFn ***REMOVED***
	l.next()
	l.emit(tokenRightCurlyBrace)
	return l.lexRvalue
***REMOVED***

func (l *tomlLexer) lexDate() tomlLexStateFn ***REMOVED***
	l.emit(tokenDate)
	return l.lexRvalue
***REMOVED***

func (l *tomlLexer) lexTrue() tomlLexStateFn ***REMOVED***
	l.fastForward(4)
	l.emit(tokenTrue)
	return l.lexRvalue
***REMOVED***

func (l *tomlLexer) lexFalse() tomlLexStateFn ***REMOVED***
	l.fastForward(5)
	l.emit(tokenFalse)
	return l.lexRvalue
***REMOVED***

func (l *tomlLexer) lexInf() tomlLexStateFn ***REMOVED***
	l.fastForward(3)
	l.emit(tokenInf)
	return l.lexRvalue
***REMOVED***

func (l *tomlLexer) lexNan() tomlLexStateFn ***REMOVED***
	l.fastForward(3)
	l.emit(tokenNan)
	return l.lexRvalue
***REMOVED***

func (l *tomlLexer) lexEqual() tomlLexStateFn ***REMOVED***
	l.next()
	l.emit(tokenEqual)
	return l.lexRvalue
***REMOVED***

func (l *tomlLexer) lexComma() tomlLexStateFn ***REMOVED***
	l.next()
	l.emit(tokenComma)
	return l.lexRvalue
***REMOVED***

// Parse the key and emits its value without escape sequences.
// bare keys, basic string keys and literal string keys are supported.
func (l *tomlLexer) lexKey() tomlLexStateFn ***REMOVED***
	growingString := ""

	for r := l.peek(); isKeyChar(r) || r == '\n' || r == '\r'; r = l.peek() ***REMOVED***
		if r == '"' ***REMOVED***
			l.next()
			str, err := l.lexStringAsString(`"`, false, true)
			if err != nil ***REMOVED***
				return l.errorf(err.Error())
			***REMOVED***
			growingString += str
			l.next()
			continue
		***REMOVED*** else if r == '\'' ***REMOVED***
			l.next()
			str, err := l.lexLiteralStringAsString(`'`, false)
			if err != nil ***REMOVED***
				return l.errorf(err.Error())
			***REMOVED***
			growingString += str
			l.next()
			continue
		***REMOVED*** else if r == '\n' ***REMOVED***
			return l.errorf("keys cannot contain new lines")
		***REMOVED*** else if isSpace(r) ***REMOVED***
			break
		***REMOVED*** else if !isValidBareChar(r) ***REMOVED***
			return l.errorf("keys cannot contain %c character", r)
		***REMOVED***
		growingString += string(r)
		l.next()
	***REMOVED***
	l.emitWithValue(tokenKey, growingString)
	return l.lexVoid
***REMOVED***

func (l *tomlLexer) lexComment(previousState tomlLexStateFn) tomlLexStateFn ***REMOVED***
	return func() tomlLexStateFn ***REMOVED***
		for next := l.peek(); next != '\n' && next != eof; next = l.peek() ***REMOVED***
			if next == '\r' && l.follow("\r\n") ***REMOVED***
				break
			***REMOVED***
			l.next()
		***REMOVED***
		l.ignore()
		return previousState
	***REMOVED***
***REMOVED***

func (l *tomlLexer) lexLeftBracket() tomlLexStateFn ***REMOVED***
	l.next()
	l.emit(tokenLeftBracket)
	return l.lexRvalue
***REMOVED***

func (l *tomlLexer) lexLiteralStringAsString(terminator string, discardLeadingNewLine bool) (string, error) ***REMOVED***
	growingString := ""

	if discardLeadingNewLine ***REMOVED***
		if l.follow("\r\n") ***REMOVED***
			l.skip()
			l.skip()
		***REMOVED*** else if l.peek() == '\n' ***REMOVED***
			l.skip()
		***REMOVED***
	***REMOVED***

	// find end of string
	for ***REMOVED***
		if l.follow(terminator) ***REMOVED***
			return growingString, nil
		***REMOVED***

		next := l.peek()
		if next == eof ***REMOVED***
			break
		***REMOVED***
		growingString += string(l.next())
	***REMOVED***

	return "", errors.New("unclosed string")
***REMOVED***

func (l *tomlLexer) lexLiteralString() tomlLexStateFn ***REMOVED***
	l.skip()

	// handle special case for triple-quote
	terminator := "'"
	discardLeadingNewLine := false
	if l.follow("''") ***REMOVED***
		l.skip()
		l.skip()
		terminator = "'''"
		discardLeadingNewLine = true
	***REMOVED***

	str, err := l.lexLiteralStringAsString(terminator, discardLeadingNewLine)
	if err != nil ***REMOVED***
		return l.errorf(err.Error())
	***REMOVED***

	l.emitWithValue(tokenString, str)
	l.fastForward(len(terminator))
	l.ignore()
	return l.lexRvalue
***REMOVED***

// Lex a string and return the results as a string.
// Terminator is the substring indicating the end of the token.
// The resulting string does not include the terminator.
func (l *tomlLexer) lexStringAsString(terminator string, discardLeadingNewLine, acceptNewLines bool) (string, error) ***REMOVED***
	growingString := ""

	if discardLeadingNewLine ***REMOVED***
		if l.follow("\r\n") ***REMOVED***
			l.skip()
			l.skip()
		***REMOVED*** else if l.peek() == '\n' ***REMOVED***
			l.skip()
		***REMOVED***
	***REMOVED***

	for ***REMOVED***
		if l.follow(terminator) ***REMOVED***
			return growingString, nil
		***REMOVED***

		if l.follow("\\") ***REMOVED***
			l.next()
			switch l.peek() ***REMOVED***
			case '\r':
				fallthrough
			case '\n':
				fallthrough
			case '\t':
				fallthrough
			case ' ':
				// skip all whitespace chars following backslash
				for strings.ContainsRune("\r\n\t ", l.peek()) ***REMOVED***
					l.next()
				***REMOVED***
			case '"':
				growingString += "\""
				l.next()
			case 'n':
				growingString += "\n"
				l.next()
			case 'b':
				growingString += "\b"
				l.next()
			case 'f':
				growingString += "\f"
				l.next()
			case '/':
				growingString += "/"
				l.next()
			case 't':
				growingString += "\t"
				l.next()
			case 'r':
				growingString += "\r"
				l.next()
			case '\\':
				growingString += "\\"
				l.next()
			case 'u':
				l.next()
				code := ""
				for i := 0; i < 4; i++ ***REMOVED***
					c := l.peek()
					if !isHexDigit(c) ***REMOVED***
						return "", errors.New("unfinished unicode escape")
					***REMOVED***
					l.next()
					code = code + string(c)
				***REMOVED***
				intcode, err := strconv.ParseInt(code, 16, 32)
				if err != nil ***REMOVED***
					return "", errors.New("invalid unicode escape: \\u" + code)
				***REMOVED***
				growingString += string(rune(intcode))
			case 'U':
				l.next()
				code := ""
				for i := 0; i < 8; i++ ***REMOVED***
					c := l.peek()
					if !isHexDigit(c) ***REMOVED***
						return "", errors.New("unfinished unicode escape")
					***REMOVED***
					l.next()
					code = code + string(c)
				***REMOVED***
				intcode, err := strconv.ParseInt(code, 16, 64)
				if err != nil ***REMOVED***
					return "", errors.New("invalid unicode escape: \\U" + code)
				***REMOVED***
				growingString += string(rune(intcode))
			default:
				return "", errors.New("invalid escape sequence: \\" + string(l.peek()))
			***REMOVED***
		***REMOVED*** else ***REMOVED***
			r := l.peek()

			if 0x00 <= r && r <= 0x1F && !(acceptNewLines && (r == '\n' || r == '\r')) ***REMOVED***
				return "", fmt.Errorf("unescaped control character %U", r)
			***REMOVED***
			l.next()
			growingString += string(r)
		***REMOVED***

		if l.peek() == eof ***REMOVED***
			break
		***REMOVED***
	***REMOVED***

	return "", errors.New("unclosed string")
***REMOVED***

func (l *tomlLexer) lexString() tomlLexStateFn ***REMOVED***
	l.skip()

	// handle special case for triple-quote
	terminator := `"`
	discardLeadingNewLine := false
	acceptNewLines := false
	if l.follow(`""`) ***REMOVED***
		l.skip()
		l.skip()
		terminator = `"""`
		discardLeadingNewLine = true
		acceptNewLines = true
	***REMOVED***

	str, err := l.lexStringAsString(terminator, discardLeadingNewLine, acceptNewLines)

	if err != nil ***REMOVED***
		return l.errorf(err.Error())
	***REMOVED***

	l.emitWithValue(tokenString, str)
	l.fastForward(len(terminator))
	l.ignore()
	return l.lexRvalue
***REMOVED***

func (l *tomlLexer) lexTableKey() tomlLexStateFn ***REMOVED***
	l.next()

	if l.peek() == '[' ***REMOVED***
		// token '[[' signifies an array of tables
		l.next()
		l.emit(tokenDoubleLeftBracket)
		return l.lexInsideTableArrayKey
	***REMOVED***
	// vanilla table key
	l.emit(tokenLeftBracket)
	return l.lexInsideTableKey
***REMOVED***

// Parse the key till "]]", but only bare keys are supported
func (l *tomlLexer) lexInsideTableArrayKey() tomlLexStateFn ***REMOVED***
	for r := l.peek(); r != eof; r = l.peek() ***REMOVED***
		switch r ***REMOVED***
		case ']':
			if l.currentTokenStop > l.currentTokenStart ***REMOVED***
				l.emit(tokenKeyGroupArray)
			***REMOVED***
			l.next()
			if l.peek() != ']' ***REMOVED***
				break
			***REMOVED***
			l.next()
			l.emit(tokenDoubleRightBracket)
			return l.lexVoid
		case '[':
			return l.errorf("table array key cannot contain ']'")
		default:
			l.next()
		***REMOVED***
	***REMOVED***
	return l.errorf("unclosed table array key")
***REMOVED***

// Parse the key till "]" but only bare keys are supported
func (l *tomlLexer) lexInsideTableKey() tomlLexStateFn ***REMOVED***
	for r := l.peek(); r != eof; r = l.peek() ***REMOVED***
		switch r ***REMOVED***
		case ']':
			if l.currentTokenStop > l.currentTokenStart ***REMOVED***
				l.emit(tokenKeyGroup)
			***REMOVED***
			l.next()
			l.emit(tokenRightBracket)
			return l.lexVoid
		case '[':
			return l.errorf("table key cannot contain ']'")
		default:
			l.next()
		***REMOVED***
	***REMOVED***
	return l.errorf("unclosed table key")
***REMOVED***

func (l *tomlLexer) lexRightBracket() tomlLexStateFn ***REMOVED***
	l.next()
	l.emit(tokenRightBracket)
	return l.lexRvalue
***REMOVED***

type validRuneFn func(r rune) bool

func isValidHexRune(r rune) bool ***REMOVED***
	return r >= 'a' && r <= 'f' ||
		r >= 'A' && r <= 'F' ||
		r >= '0' && r <= '9' ||
		r == '_'
***REMOVED***

func isValidOctalRune(r rune) bool ***REMOVED***
	return r >= '0' && r <= '7' || r == '_'
***REMOVED***

func isValidBinaryRune(r rune) bool ***REMOVED***
	return r == '0' || r == '1' || r == '_'
***REMOVED***

func (l *tomlLexer) lexNumber() tomlLexStateFn ***REMOVED***
	r := l.peek()

	if r == '0' ***REMOVED***
		follow := l.peekString(2)
		if len(follow) == 2 ***REMOVED***
			var isValidRune validRuneFn
			switch follow[1] ***REMOVED***
			case 'x':
				isValidRune = isValidHexRune
			case 'o':
				isValidRune = isValidOctalRune
			case 'b':
				isValidRune = isValidBinaryRune
			default:
				if follow[1] >= 'a' && follow[1] <= 'z' || follow[1] >= 'A' && follow[1] <= 'Z' ***REMOVED***
					return l.errorf("unknown number base: %s. possible options are x (hex) o (octal) b (binary)", string(follow[1]))
				***REMOVED***
			***REMOVED***

			if isValidRune != nil ***REMOVED***
				l.next()
				l.next()
				digitSeen := false
				for ***REMOVED***
					next := l.peek()
					if !isValidRune(next) ***REMOVED***
						break
					***REMOVED***
					digitSeen = true
					l.next()
				***REMOVED***

				if !digitSeen ***REMOVED***
					return l.errorf("number needs at least one digit")
				***REMOVED***

				l.emit(tokenInteger)

				return l.lexRvalue
			***REMOVED***
		***REMOVED***
	***REMOVED***

	if r == '+' || r == '-' ***REMOVED***
		l.next()
		if l.follow("inf") ***REMOVED***
			return l.lexInf
		***REMOVED***
		if l.follow("nan") ***REMOVED***
			return l.lexNan
		***REMOVED***
	***REMOVED***

	pointSeen := false
	expSeen := false
	digitSeen := false
	for ***REMOVED***
		next := l.peek()
		if next == '.' ***REMOVED***
			if pointSeen ***REMOVED***
				return l.errorf("cannot have two dots in one float")
			***REMOVED***
			l.next()
			if !isDigit(l.peek()) ***REMOVED***
				return l.errorf("float cannot end with a dot")
			***REMOVED***
			pointSeen = true
		***REMOVED*** else if next == 'e' || next == 'E' ***REMOVED***
			expSeen = true
			l.next()
			r := l.peek()
			if r == '+' || r == '-' ***REMOVED***
				l.next()
			***REMOVED***
		***REMOVED*** else if isDigit(next) ***REMOVED***
			digitSeen = true
			l.next()
		***REMOVED*** else if next == '_' ***REMOVED***
			l.next()
		***REMOVED*** else ***REMOVED***
			break
		***REMOVED***
		if pointSeen && !digitSeen ***REMOVED***
			return l.errorf("cannot start float with a dot")
		***REMOVED***
	***REMOVED***

	if !digitSeen ***REMOVED***
		return l.errorf("no digit in that number")
	***REMOVED***
	if pointSeen || expSeen ***REMOVED***
		l.emit(tokenFloat)
	***REMOVED*** else ***REMOVED***
		l.emit(tokenInteger)
	***REMOVED***
	return l.lexRvalue
***REMOVED***

func (l *tomlLexer) run() ***REMOVED***
	for state := l.lexVoid; state != nil; ***REMOVED***
		state = state()
	***REMOVED***
***REMOVED***

func init() ***REMOVED***
	dateRegexp = regexp.MustCompile(`^\d***REMOVED***1,4***REMOVED***-\d***REMOVED***2***REMOVED***-\d***REMOVED***2***REMOVED***T\d***REMOVED***2***REMOVED***:\d***REMOVED***2***REMOVED***:\d***REMOVED***2***REMOVED***(\.\d***REMOVED***1,9***REMOVED***)?(Z|[+-]\d***REMOVED***2***REMOVED***:\d***REMOVED***2***REMOVED***)`)
***REMOVED***

// Entry point
func lexToml(inputBytes []byte) []token ***REMOVED***
	runes := bytes.Runes(inputBytes)
	l := &tomlLexer***REMOVED***
		input:         runes,
		tokens:        make([]token, 0, 256),
		line:          1,
		col:           1,
		endbufferLine: 1,
		endbufferCol:  1,
	***REMOVED***
	l.run()
	return l.tokens
***REMOVED***
