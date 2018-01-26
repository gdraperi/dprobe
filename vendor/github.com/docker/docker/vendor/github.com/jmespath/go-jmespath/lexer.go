package jmespath

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"unicode/utf8"
)

type token struct ***REMOVED***
	tokenType tokType
	value     string
	position  int
	length    int
***REMOVED***

type tokType int

const eof = -1

// Lexer contains information about the expression being tokenized.
type Lexer struct ***REMOVED***
	expression string       // The expression provided by the user.
	currentPos int          // The current position in the string.
	lastWidth  int          // The width of the current rune.  This
	buf        bytes.Buffer // Internal buffer used for building up values.
***REMOVED***

// SyntaxError is the main error used whenever a lexing or parsing error occurs.
type SyntaxError struct ***REMOVED***
	msg        string // Error message displayed to user
	Expression string // Expression that generated a SyntaxError
	Offset     int    // The location in the string where the error occurred
***REMOVED***

func (e SyntaxError) Error() string ***REMOVED***
	// In the future, it would be good to underline the specific
	// location where the error occurred.
	return "SyntaxError: " + e.msg
***REMOVED***

// HighlightLocation will show where the syntax error occurred.
// It will place a "^" character on a line below the expression
// at the point where the syntax error occurred.
func (e SyntaxError) HighlightLocation() string ***REMOVED***
	return e.Expression + "\n" + strings.Repeat(" ", e.Offset) + "^"
***REMOVED***

//go:generate stringer -type=tokType
const (
	tUnknown tokType = iota
	tStar
	tDot
	tFilter
	tFlatten
	tLparen
	tRparen
	tLbracket
	tRbracket
	tLbrace
	tRbrace
	tOr
	tPipe
	tNumber
	tUnquotedIdentifier
	tQuotedIdentifier
	tComma
	tColon
	tLT
	tLTE
	tGT
	tGTE
	tEQ
	tNE
	tJSONLiteral
	tStringLiteral
	tCurrent
	tExpref
	tAnd
	tNot
	tEOF
)

var basicTokens = map[rune]tokType***REMOVED***
	'.': tDot,
	'*': tStar,
	',': tComma,
	':': tColon,
	'***REMOVED***': tLbrace,
	'***REMOVED***': tRbrace,
	']': tRbracket, // tLbracket not included because it could be "[]"
	'(': tLparen,
	')': tRparen,
	'@': tCurrent,
***REMOVED***

// Bit mask for [a-zA-Z_] shifted down 64 bits to fit in a single uint64.
// When using this bitmask just be sure to shift the rune down 64 bits
// before checking against identifierStartBits.
const identifierStartBits uint64 = 576460745995190270

// Bit mask for [a-zA-Z0-9], 128 bits -> 2 uint64s.
var identifierTrailingBits = [2]uint64***REMOVED***287948901175001088, 576460745995190270***REMOVED***

var whiteSpace = map[rune]bool***REMOVED***
	' ': true, '\t': true, '\n': true, '\r': true,
***REMOVED***

func (t token) String() string ***REMOVED***
	return fmt.Sprintf("Token***REMOVED***%+v, %s, %d, %d***REMOVED***",
		t.tokenType, t.value, t.position, t.length)
***REMOVED***

// NewLexer creates a new JMESPath lexer.
func NewLexer() *Lexer ***REMOVED***
	lexer := Lexer***REMOVED******REMOVED***
	return &lexer
***REMOVED***

func (lexer *Lexer) next() rune ***REMOVED***
	if lexer.currentPos >= len(lexer.expression) ***REMOVED***
		lexer.lastWidth = 0
		return eof
	***REMOVED***
	r, w := utf8.DecodeRuneInString(lexer.expression[lexer.currentPos:])
	lexer.lastWidth = w
	lexer.currentPos += w
	return r
***REMOVED***

func (lexer *Lexer) back() ***REMOVED***
	lexer.currentPos -= lexer.lastWidth
***REMOVED***

func (lexer *Lexer) peek() rune ***REMOVED***
	t := lexer.next()
	lexer.back()
	return t
***REMOVED***

// tokenize takes an expression and returns corresponding tokens.
func (lexer *Lexer) tokenize(expression string) ([]token, error) ***REMOVED***
	var tokens []token
	lexer.expression = expression
	lexer.currentPos = 0
	lexer.lastWidth = 0
loop:
	for ***REMOVED***
		r := lexer.next()
		if identifierStartBits&(1<<(uint64(r)-64)) > 0 ***REMOVED***
			t := lexer.consumeUnquotedIdentifier()
			tokens = append(tokens, t)
		***REMOVED*** else if val, ok := basicTokens[r]; ok ***REMOVED***
			// Basic single char token.
			t := token***REMOVED***
				tokenType: val,
				value:     string(r),
				position:  lexer.currentPos - lexer.lastWidth,
				length:    1,
			***REMOVED***
			tokens = append(tokens, t)
		***REMOVED*** else if r == '-' || (r >= '0' && r <= '9') ***REMOVED***
			t := lexer.consumeNumber()
			tokens = append(tokens, t)
		***REMOVED*** else if r == '[' ***REMOVED***
			t := lexer.consumeLBracket()
			tokens = append(tokens, t)
		***REMOVED*** else if r == '"' ***REMOVED***
			t, err := lexer.consumeQuotedIdentifier()
			if err != nil ***REMOVED***
				return tokens, err
			***REMOVED***
			tokens = append(tokens, t)
		***REMOVED*** else if r == '\'' ***REMOVED***
			t, err := lexer.consumeRawStringLiteral()
			if err != nil ***REMOVED***
				return tokens, err
			***REMOVED***
			tokens = append(tokens, t)
		***REMOVED*** else if r == '`' ***REMOVED***
			t, err := lexer.consumeLiteral()
			if err != nil ***REMOVED***
				return tokens, err
			***REMOVED***
			tokens = append(tokens, t)
		***REMOVED*** else if r == '|' ***REMOVED***
			t := lexer.matchOrElse(r, '|', tOr, tPipe)
			tokens = append(tokens, t)
		***REMOVED*** else if r == '<' ***REMOVED***
			t := lexer.matchOrElse(r, '=', tLTE, tLT)
			tokens = append(tokens, t)
		***REMOVED*** else if r == '>' ***REMOVED***
			t := lexer.matchOrElse(r, '=', tGTE, tGT)
			tokens = append(tokens, t)
		***REMOVED*** else if r == '!' ***REMOVED***
			t := lexer.matchOrElse(r, '=', tNE, tNot)
			tokens = append(tokens, t)
		***REMOVED*** else if r == '=' ***REMOVED***
			t := lexer.matchOrElse(r, '=', tEQ, tUnknown)
			tokens = append(tokens, t)
		***REMOVED*** else if r == '&' ***REMOVED***
			t := lexer.matchOrElse(r, '&', tAnd, tExpref)
			tokens = append(tokens, t)
		***REMOVED*** else if r == eof ***REMOVED***
			break loop
		***REMOVED*** else if _, ok := whiteSpace[r]; ok ***REMOVED***
			// Ignore whitespace
		***REMOVED*** else ***REMOVED***
			return tokens, lexer.syntaxError(fmt.Sprintf("Unknown char: %s", strconv.QuoteRuneToASCII(r)))
		***REMOVED***
	***REMOVED***
	tokens = append(tokens, token***REMOVED***tEOF, "", len(lexer.expression), 0***REMOVED***)
	return tokens, nil
***REMOVED***

// Consume characters until the ending rune "r" is reached.
// If the end of the expression is reached before seeing the
// terminating rune "r", then an error is returned.
// If no error occurs then the matching substring is returned.
// The returned string will not include the ending rune.
func (lexer *Lexer) consumeUntil(end rune) (string, error) ***REMOVED***
	start := lexer.currentPos
	current := lexer.next()
	for current != end && current != eof ***REMOVED***
		if current == '\\' && lexer.peek() != eof ***REMOVED***
			lexer.next()
		***REMOVED***
		current = lexer.next()
	***REMOVED***
	if lexer.lastWidth == 0 ***REMOVED***
		// Then we hit an EOF so we never reached the closing
		// delimiter.
		return "", SyntaxError***REMOVED***
			msg:        "Unclosed delimiter: " + string(end),
			Expression: lexer.expression,
			Offset:     len(lexer.expression),
		***REMOVED***
	***REMOVED***
	return lexer.expression[start : lexer.currentPos-lexer.lastWidth], nil
***REMOVED***

func (lexer *Lexer) consumeLiteral() (token, error) ***REMOVED***
	start := lexer.currentPos
	value, err := lexer.consumeUntil('`')
	if err != nil ***REMOVED***
		return token***REMOVED******REMOVED***, err
	***REMOVED***
	value = strings.Replace(value, "\\`", "`", -1)
	return token***REMOVED***
		tokenType: tJSONLiteral,
		value:     value,
		position:  start,
		length:    len(value),
	***REMOVED***, nil
***REMOVED***

func (lexer *Lexer) consumeRawStringLiteral() (token, error) ***REMOVED***
	start := lexer.currentPos
	currentIndex := start
	current := lexer.next()
	for current != '\'' && lexer.peek() != eof ***REMOVED***
		if current == '\\' && lexer.peek() == '\'' ***REMOVED***
			chunk := lexer.expression[currentIndex : lexer.currentPos-1]
			lexer.buf.WriteString(chunk)
			lexer.buf.WriteString("'")
			lexer.next()
			currentIndex = lexer.currentPos
		***REMOVED***
		current = lexer.next()
	***REMOVED***
	if lexer.lastWidth == 0 ***REMOVED***
		// Then we hit an EOF so we never reached the closing
		// delimiter.
		return token***REMOVED******REMOVED***, SyntaxError***REMOVED***
			msg:        "Unclosed delimiter: '",
			Expression: lexer.expression,
			Offset:     len(lexer.expression),
		***REMOVED***
	***REMOVED***
	if currentIndex < lexer.currentPos ***REMOVED***
		lexer.buf.WriteString(lexer.expression[currentIndex : lexer.currentPos-1])
	***REMOVED***
	value := lexer.buf.String()
	// Reset the buffer so it can reused again.
	lexer.buf.Reset()
	return token***REMOVED***
		tokenType: tStringLiteral,
		value:     value,
		position:  start,
		length:    len(value),
	***REMOVED***, nil
***REMOVED***

func (lexer *Lexer) syntaxError(msg string) SyntaxError ***REMOVED***
	return SyntaxError***REMOVED***
		msg:        msg,
		Expression: lexer.expression,
		Offset:     lexer.currentPos - 1,
	***REMOVED***
***REMOVED***

// Checks for a two char token, otherwise matches a single character
// token. This is used whenever a two char token overlaps a single
// char token, e.g. "||" -> tPipe, "|" -> tOr.
func (lexer *Lexer) matchOrElse(first rune, second rune, matchedType tokType, singleCharType tokType) token ***REMOVED***
	start := lexer.currentPos - lexer.lastWidth
	nextRune := lexer.next()
	var t token
	if nextRune == second ***REMOVED***
		t = token***REMOVED***
			tokenType: matchedType,
			value:     string(first) + string(second),
			position:  start,
			length:    2,
		***REMOVED***
	***REMOVED*** else ***REMOVED***
		lexer.back()
		t = token***REMOVED***
			tokenType: singleCharType,
			value:     string(first),
			position:  start,
			length:    1,
		***REMOVED***
	***REMOVED***
	return t
***REMOVED***

func (lexer *Lexer) consumeLBracket() token ***REMOVED***
	// There's three options here:
	// 1. A filter expression "[?"
	// 2. A flatten operator "[]"
	// 3. A bare rbracket "["
	start := lexer.currentPos - lexer.lastWidth
	nextRune := lexer.next()
	var t token
	if nextRune == '?' ***REMOVED***
		t = token***REMOVED***
			tokenType: tFilter,
			value:     "[?",
			position:  start,
			length:    2,
		***REMOVED***
	***REMOVED*** else if nextRune == ']' ***REMOVED***
		t = token***REMOVED***
			tokenType: tFlatten,
			value:     "[]",
			position:  start,
			length:    2,
		***REMOVED***
	***REMOVED*** else ***REMOVED***
		t = token***REMOVED***
			tokenType: tLbracket,
			value:     "[",
			position:  start,
			length:    1,
		***REMOVED***
		lexer.back()
	***REMOVED***
	return t
***REMOVED***

func (lexer *Lexer) consumeQuotedIdentifier() (token, error) ***REMOVED***
	start := lexer.currentPos
	value, err := lexer.consumeUntil('"')
	if err != nil ***REMOVED***
		return token***REMOVED******REMOVED***, err
	***REMOVED***
	var decoded string
	asJSON := []byte("\"" + value + "\"")
	if err := json.Unmarshal([]byte(asJSON), &decoded); err != nil ***REMOVED***
		return token***REMOVED******REMOVED***, err
	***REMOVED***
	return token***REMOVED***
		tokenType: tQuotedIdentifier,
		value:     decoded,
		position:  start - 1,
		length:    len(decoded),
	***REMOVED***, nil
***REMOVED***

func (lexer *Lexer) consumeUnquotedIdentifier() token ***REMOVED***
	// Consume runes until we reach the end of an unquoted
	// identifier.
	start := lexer.currentPos - lexer.lastWidth
	for ***REMOVED***
		r := lexer.next()
		if r < 0 || r > 128 || identifierTrailingBits[uint64(r)/64]&(1<<(uint64(r)%64)) == 0 ***REMOVED***
			lexer.back()
			break
		***REMOVED***
	***REMOVED***
	value := lexer.expression[start:lexer.currentPos]
	return token***REMOVED***
		tokenType: tUnquotedIdentifier,
		value:     value,
		position:  start,
		length:    lexer.currentPos - start,
	***REMOVED***
***REMOVED***

func (lexer *Lexer) consumeNumber() token ***REMOVED***
	// Consume runes until we reach something that's not a number.
	start := lexer.currentPos - lexer.lastWidth
	for ***REMOVED***
		r := lexer.next()
		if r < '0' || r > '9' ***REMOVED***
			lexer.back()
			break
		***REMOVED***
	***REMOVED***
	value := lexer.expression[start:lexer.currentPos]
	return token***REMOVED***
		tokenType: tNumber,
		value:     value,
		position:  start,
		length:    lexer.currentPos - start,
	***REMOVED***
***REMOVED***
