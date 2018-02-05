// Package scanner implements a scanner for HCL (HashiCorp Configuration
// Language) source text.
package scanner

import (
	"bytes"
	"fmt"
	"os"
	"regexp"
	"unicode"
	"unicode/utf8"

	"github.com/hashicorp/hcl/hcl/token"
)

// eof represents a marker rune for the end of the reader.
const eof = rune(0)

// Scanner defines a lexical scanner
type Scanner struct ***REMOVED***
	buf *bytes.Buffer // Source buffer for advancing and scanning
	src []byte        // Source buffer for immutable access

	// Source Position
	srcPos  token.Pos // current position
	prevPos token.Pos // previous position, used for peek() method

	lastCharLen int // length of last character in bytes
	lastLineLen int // length of last line in characters (for correct column reporting)

	tokStart int // token text start position
	tokEnd   int // token text end  position

	// Error is called for each error encountered. If no Error
	// function is set, the error is reported to os.Stderr.
	Error func(pos token.Pos, msg string)

	// ErrorCount is incremented by one for each error encountered.
	ErrorCount int

	// tokPos is the start position of most recently scanned token; set by
	// Scan. The Filename field is always left untouched by the Scanner.  If
	// an error is reported (via Error) and Position is invalid, the scanner is
	// not inside a token.
	tokPos token.Pos
***REMOVED***

// New creates and initializes a new instance of Scanner using src as
// its source content.
func New(src []byte) *Scanner ***REMOVED***
	// even though we accept a src, we read from a io.Reader compatible type
	// (*bytes.Buffer). So in the future we might easily change it to streaming
	// read.
	b := bytes.NewBuffer(src)
	s := &Scanner***REMOVED***
		buf: b,
		src: src,
	***REMOVED***

	// srcPosition always starts with 1
	s.srcPos.Line = 1
	return s
***REMOVED***

// next reads the next rune from the bufferred reader. Returns the rune(0) if
// an error occurs (or io.EOF is returned).
func (s *Scanner) next() rune ***REMOVED***
	ch, size, err := s.buf.ReadRune()
	if err != nil ***REMOVED***
		// advance for error reporting
		s.srcPos.Column++
		s.srcPos.Offset += size
		s.lastCharLen = size
		return eof
	***REMOVED***

	if ch == utf8.RuneError && size == 1 ***REMOVED***
		s.srcPos.Column++
		s.srcPos.Offset += size
		s.lastCharLen = size
		s.err("illegal UTF-8 encoding")
		return ch
	***REMOVED***

	// remember last position
	s.prevPos = s.srcPos

	s.srcPos.Column++
	s.lastCharLen = size
	s.srcPos.Offset += size

	if ch == '\n' ***REMOVED***
		s.srcPos.Line++
		s.lastLineLen = s.srcPos.Column
		s.srcPos.Column = 0
	***REMOVED***

	// If we see a null character with data left, then that is an error
	if ch == '\x00' && s.buf.Len() > 0 ***REMOVED***
		s.err("unexpected null character (0x00)")
		return eof
	***REMOVED***

	// debug
	// fmt.Printf("ch: %q, offset:column: %d:%d\n", ch, s.srcPos.Offset, s.srcPos.Column)
	return ch
***REMOVED***

// unread unreads the previous read Rune and updates the source position
func (s *Scanner) unread() ***REMOVED***
	if err := s.buf.UnreadRune(); err != nil ***REMOVED***
		panic(err) // this is user fault, we should catch it
	***REMOVED***
	s.srcPos = s.prevPos // put back last position
***REMOVED***

// peek returns the next rune without advancing the reader.
func (s *Scanner) peek() rune ***REMOVED***
	peek, _, err := s.buf.ReadRune()
	if err != nil ***REMOVED***
		return eof
	***REMOVED***

	s.buf.UnreadRune()
	return peek
***REMOVED***

// Scan scans the next token and returns the token.
func (s *Scanner) Scan() token.Token ***REMOVED***
	ch := s.next()

	// skip white space
	for isWhitespace(ch) ***REMOVED***
		ch = s.next()
	***REMOVED***

	var tok token.Type

	// token text markings
	s.tokStart = s.srcPos.Offset - s.lastCharLen

	// token position, initial next() is moving the offset by one(size of rune
	// actually), though we are interested with the starting point
	s.tokPos.Offset = s.srcPos.Offset - s.lastCharLen
	if s.srcPos.Column > 0 ***REMOVED***
		// common case: last character was not a '\n'
		s.tokPos.Line = s.srcPos.Line
		s.tokPos.Column = s.srcPos.Column
	***REMOVED*** else ***REMOVED***
		// last character was a '\n'
		// (we cannot be at the beginning of the source
		// since we have called next() at least once)
		s.tokPos.Line = s.srcPos.Line - 1
		s.tokPos.Column = s.lastLineLen
	***REMOVED***

	switch ***REMOVED***
	case isLetter(ch):
		tok = token.IDENT
		lit := s.scanIdentifier()
		if lit == "true" || lit == "false" ***REMOVED***
			tok = token.BOOL
		***REMOVED***
	case isDecimal(ch):
		tok = s.scanNumber(ch)
	default:
		switch ch ***REMOVED***
		case eof:
			tok = token.EOF
		case '"':
			tok = token.STRING
			s.scanString()
		case '#', '/':
			tok = token.COMMENT
			s.scanComment(ch)
		case '.':
			tok = token.PERIOD
			ch = s.peek()
			if isDecimal(ch) ***REMOVED***
				tok = token.FLOAT
				ch = s.scanMantissa(ch)
				ch = s.scanExponent(ch)
			***REMOVED***
		case '<':
			tok = token.HEREDOC
			s.scanHeredoc()
		case '[':
			tok = token.LBRACK
		case ']':
			tok = token.RBRACK
		case '***REMOVED***':
			tok = token.LBRACE
		case '***REMOVED***':
			tok = token.RBRACE
		case ',':
			tok = token.COMMA
		case '=':
			tok = token.ASSIGN
		case '+':
			tok = token.ADD
		case '-':
			if isDecimal(s.peek()) ***REMOVED***
				ch := s.next()
				tok = s.scanNumber(ch)
			***REMOVED*** else ***REMOVED***
				tok = token.SUB
			***REMOVED***
		default:
			s.err("illegal char")
		***REMOVED***
	***REMOVED***

	// finish token ending
	s.tokEnd = s.srcPos.Offset

	// create token literal
	var tokenText string
	if s.tokStart >= 0 ***REMOVED***
		tokenText = string(s.src[s.tokStart:s.tokEnd])
	***REMOVED***
	s.tokStart = s.tokEnd // ensure idempotency of tokenText() call

	return token.Token***REMOVED***
		Type: tok,
		Pos:  s.tokPos,
		Text: tokenText,
	***REMOVED***
***REMOVED***

func (s *Scanner) scanComment(ch rune) ***REMOVED***
	// single line comments
	if ch == '#' || (ch == '/' && s.peek() != '*') ***REMOVED***
		if ch == '/' && s.peek() != '/' ***REMOVED***
			s.err("expected '/' for comment")
			return
		***REMOVED***

		ch = s.next()
		for ch != '\n' && ch >= 0 && ch != eof ***REMOVED***
			ch = s.next()
		***REMOVED***
		if ch != eof && ch >= 0 ***REMOVED***
			s.unread()
		***REMOVED***
		return
	***REMOVED***

	// be sure we get the character after /* This allows us to find comment's
	// that are not erminated
	if ch == '/' ***REMOVED***
		s.next()
		ch = s.next() // read character after "/*"
	***REMOVED***

	// look for /* - style comments
	for ***REMOVED***
		if ch < 0 || ch == eof ***REMOVED***
			s.err("comment not terminated")
			break
		***REMOVED***

		ch0 := ch
		ch = s.next()
		if ch0 == '*' && ch == '/' ***REMOVED***
			break
		***REMOVED***
	***REMOVED***
***REMOVED***

// scanNumber scans a HCL number definition starting with the given rune
func (s *Scanner) scanNumber(ch rune) token.Type ***REMOVED***
	if ch == '0' ***REMOVED***
		// check for hexadecimal, octal or float
		ch = s.next()
		if ch == 'x' || ch == 'X' ***REMOVED***
			// hexadecimal
			ch = s.next()
			found := false
			for isHexadecimal(ch) ***REMOVED***
				ch = s.next()
				found = true
			***REMOVED***

			if !found ***REMOVED***
				s.err("illegal hexadecimal number")
			***REMOVED***

			if ch != eof ***REMOVED***
				s.unread()
			***REMOVED***

			return token.NUMBER
		***REMOVED***

		// now it's either something like: 0421(octal) or 0.1231(float)
		illegalOctal := false
		for isDecimal(ch) ***REMOVED***
			ch = s.next()
			if ch == '8' || ch == '9' ***REMOVED***
				// this is just a possibility. For example 0159 is illegal, but
				// 0159.23 is valid. So we mark a possible illegal octal. If
				// the next character is not a period, we'll print the error.
				illegalOctal = true
			***REMOVED***
		***REMOVED***

		if ch == 'e' || ch == 'E' ***REMOVED***
			ch = s.scanExponent(ch)
			return token.FLOAT
		***REMOVED***

		if ch == '.' ***REMOVED***
			ch = s.scanFraction(ch)

			if ch == 'e' || ch == 'E' ***REMOVED***
				ch = s.next()
				ch = s.scanExponent(ch)
			***REMOVED***
			return token.FLOAT
		***REMOVED***

		if illegalOctal ***REMOVED***
			s.err("illegal octal number")
		***REMOVED***

		if ch != eof ***REMOVED***
			s.unread()
		***REMOVED***
		return token.NUMBER
	***REMOVED***

	s.scanMantissa(ch)
	ch = s.next() // seek forward
	if ch == 'e' || ch == 'E' ***REMOVED***
		ch = s.scanExponent(ch)
		return token.FLOAT
	***REMOVED***

	if ch == '.' ***REMOVED***
		ch = s.scanFraction(ch)
		if ch == 'e' || ch == 'E' ***REMOVED***
			ch = s.next()
			ch = s.scanExponent(ch)
		***REMOVED***
		return token.FLOAT
	***REMOVED***

	if ch != eof ***REMOVED***
		s.unread()
	***REMOVED***
	return token.NUMBER
***REMOVED***

// scanMantissa scans the mantissa beginning from the rune. It returns the next
// non decimal rune. It's used to determine wheter it's a fraction or exponent.
func (s *Scanner) scanMantissa(ch rune) rune ***REMOVED***
	scanned := false
	for isDecimal(ch) ***REMOVED***
		ch = s.next()
		scanned = true
	***REMOVED***

	if scanned && ch != eof ***REMOVED***
		s.unread()
	***REMOVED***
	return ch
***REMOVED***

// scanFraction scans the fraction after the '.' rune
func (s *Scanner) scanFraction(ch rune) rune ***REMOVED***
	if ch == '.' ***REMOVED***
		ch = s.peek() // we peek just to see if we can move forward
		ch = s.scanMantissa(ch)
	***REMOVED***
	return ch
***REMOVED***

// scanExponent scans the remaining parts of an exponent after the 'e' or 'E'
// rune.
func (s *Scanner) scanExponent(ch rune) rune ***REMOVED***
	if ch == 'e' || ch == 'E' ***REMOVED***
		ch = s.next()
		if ch == '-' || ch == '+' ***REMOVED***
			ch = s.next()
		***REMOVED***
		ch = s.scanMantissa(ch)
	***REMOVED***
	return ch
***REMOVED***

// scanHeredoc scans a heredoc string
func (s *Scanner) scanHeredoc() ***REMOVED***
	// Scan the second '<' in example: '<<EOF'
	if s.next() != '<' ***REMOVED***
		s.err("heredoc expected second '<', didn't see it")
		return
	***REMOVED***

	// Get the original offset so we can read just the heredoc ident
	offs := s.srcPos.Offset

	// Scan the identifier
	ch := s.next()

	// Indented heredoc syntax
	if ch == '-' ***REMOVED***
		ch = s.next()
	***REMOVED***

	for isLetter(ch) || isDigit(ch) ***REMOVED***
		ch = s.next()
	***REMOVED***

	// If we reached an EOF then that is not good
	if ch == eof ***REMOVED***
		s.err("heredoc not terminated")
		return
	***REMOVED***

	// Ignore the '\r' in Windows line endings
	if ch == '\r' ***REMOVED***
		if s.peek() == '\n' ***REMOVED***
			ch = s.next()
		***REMOVED***
	***REMOVED***

	// If we didn't reach a newline then that is also not good
	if ch != '\n' ***REMOVED***
		s.err("invalid characters in heredoc anchor")
		return
	***REMOVED***

	// Read the identifier
	identBytes := s.src[offs : s.srcPos.Offset-s.lastCharLen]
	if len(identBytes) == 0 ***REMOVED***
		s.err("zero-length heredoc anchor")
		return
	***REMOVED***

	var identRegexp *regexp.Regexp
	if identBytes[0] == '-' ***REMOVED***
		identRegexp = regexp.MustCompile(fmt.Sprintf(`[[:space:]]*%s\z`, identBytes[1:]))
	***REMOVED*** else ***REMOVED***
		identRegexp = regexp.MustCompile(fmt.Sprintf(`[[:space:]]*%s\z`, identBytes))
	***REMOVED***

	// Read the actual string value
	lineStart := s.srcPos.Offset
	for ***REMOVED***
		ch := s.next()

		// Special newline handling.
		if ch == '\n' ***REMOVED***
			// Math is fast, so we first compare the byte counts to see if we have a chance
			// of seeing the same identifier - if the length is less than the number of bytes
			// in the identifier, this cannot be a valid terminator.
			lineBytesLen := s.srcPos.Offset - s.lastCharLen - lineStart
			if lineBytesLen >= len(identBytes) && identRegexp.Match(s.src[lineStart:s.srcPos.Offset-s.lastCharLen]) ***REMOVED***
				break
			***REMOVED***

			// Not an anchor match, record the start of a new line
			lineStart = s.srcPos.Offset
		***REMOVED***

		if ch == eof ***REMOVED***
			s.err("heredoc not terminated")
			return
		***REMOVED***
	***REMOVED***

	return
***REMOVED***

// scanString scans a quoted string
func (s *Scanner) scanString() ***REMOVED***
	braces := 0
	for ***REMOVED***
		// '"' opening already consumed
		// read character after quote
		ch := s.next()

		if (ch == '\n' && braces == 0) || ch < 0 || ch == eof ***REMOVED***
			s.err("literal not terminated")
			return
		***REMOVED***

		if ch == '"' && braces == 0 ***REMOVED***
			break
		***REMOVED***

		// If we're going into a $***REMOVED******REMOVED*** then we can ignore quotes for awhile
		if braces == 0 && ch == '$' && s.peek() == '***REMOVED***' ***REMOVED***
			braces++
			s.next()
		***REMOVED*** else if braces > 0 && ch == '***REMOVED***' ***REMOVED***
			braces++
		***REMOVED***
		if braces > 0 && ch == '***REMOVED***' ***REMOVED***
			braces--
		***REMOVED***

		if ch == '\\' ***REMOVED***
			s.scanEscape()
		***REMOVED***
	***REMOVED***

	return
***REMOVED***

// scanEscape scans an escape sequence
func (s *Scanner) scanEscape() rune ***REMOVED***
	// http://en.cppreference.com/w/cpp/language/escape
	ch := s.next() // read character after '/'
	switch ch ***REMOVED***
	case 'a', 'b', 'f', 'n', 'r', 't', 'v', '\\', '"':
		// nothing to do
	case '0', '1', '2', '3', '4', '5', '6', '7':
		// octal notation
		ch = s.scanDigits(ch, 8, 3)
	case 'x':
		// hexademical notation
		ch = s.scanDigits(s.next(), 16, 2)
	case 'u':
		// universal character name
		ch = s.scanDigits(s.next(), 16, 4)
	case 'U':
		// universal character name
		ch = s.scanDigits(s.next(), 16, 8)
	default:
		s.err("illegal char escape")
	***REMOVED***
	return ch
***REMOVED***

// scanDigits scans a rune with the given base for n times. For example an
// octal notation \184 would yield in scanDigits(ch, 8, 3)
func (s *Scanner) scanDigits(ch rune, base, n int) rune ***REMOVED***
	start := n
	for n > 0 && digitVal(ch) < base ***REMOVED***
		ch = s.next()
		if ch == eof ***REMOVED***
			// If we see an EOF, we halt any more scanning of digits
			// immediately.
			break
		***REMOVED***

		n--
	***REMOVED***
	if n > 0 ***REMOVED***
		s.err("illegal char escape")
	***REMOVED***

	if n != start ***REMOVED***
		// we scanned all digits, put the last non digit char back,
		// only if we read anything at all
		s.unread()
	***REMOVED***

	return ch
***REMOVED***

// scanIdentifier scans an identifier and returns the literal string
func (s *Scanner) scanIdentifier() string ***REMOVED***
	offs := s.srcPos.Offset - s.lastCharLen
	ch := s.next()
	for isLetter(ch) || isDigit(ch) || ch == '-' || ch == '.' ***REMOVED***
		ch = s.next()
	***REMOVED***

	if ch != eof ***REMOVED***
		s.unread() // we got identifier, put back latest char
	***REMOVED***

	return string(s.src[offs:s.srcPos.Offset])
***REMOVED***

// recentPosition returns the position of the character immediately after the
// character or token returned by the last call to Scan.
func (s *Scanner) recentPosition() (pos token.Pos) ***REMOVED***
	pos.Offset = s.srcPos.Offset - s.lastCharLen
	switch ***REMOVED***
	case s.srcPos.Column > 0:
		// common case: last character was not a '\n'
		pos.Line = s.srcPos.Line
		pos.Column = s.srcPos.Column
	case s.lastLineLen > 0:
		// last character was a '\n'
		// (we cannot be at the beginning of the source
		// since we have called next() at least once)
		pos.Line = s.srcPos.Line - 1
		pos.Column = s.lastLineLen
	default:
		// at the beginning of the source
		pos.Line = 1
		pos.Column = 1
	***REMOVED***
	return
***REMOVED***

// err prints the error of any scanning to s.Error function. If the function is
// not defined, by default it prints them to os.Stderr
func (s *Scanner) err(msg string) ***REMOVED***
	s.ErrorCount++
	pos := s.recentPosition()

	if s.Error != nil ***REMOVED***
		s.Error(pos, msg)
		return
	***REMOVED***

	fmt.Fprintf(os.Stderr, "%s: %s\n", pos, msg)
***REMOVED***

// isHexadecimal returns true if the given rune is a letter
func isLetter(ch rune) bool ***REMOVED***
	return 'a' <= ch && ch <= 'z' || 'A' <= ch && ch <= 'Z' || ch == '_' || ch >= 0x80 && unicode.IsLetter(ch)
***REMOVED***

// isDigit returns true if the given rune is a decimal digit
func isDigit(ch rune) bool ***REMOVED***
	return '0' <= ch && ch <= '9' || ch >= 0x80 && unicode.IsDigit(ch)
***REMOVED***

// isDecimal returns true if the given rune is a decimal number
func isDecimal(ch rune) bool ***REMOVED***
	return '0' <= ch && ch <= '9'
***REMOVED***

// isHexadecimal returns true if the given rune is an hexadecimal number
func isHexadecimal(ch rune) bool ***REMOVED***
	return '0' <= ch && ch <= '9' || 'a' <= ch && ch <= 'f' || 'A' <= ch && ch <= 'F'
***REMOVED***

// isWhitespace returns true if the rune is a space, tab, newline or carriage return
func isWhitespace(ch rune) bool ***REMOVED***
	return ch == ' ' || ch == '\t' || ch == '\n' || ch == '\r'
***REMOVED***

// digitVal returns the integer value of a given octal,decimal or hexadecimal rune
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
