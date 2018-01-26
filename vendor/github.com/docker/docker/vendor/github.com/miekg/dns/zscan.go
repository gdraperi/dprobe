package dns

import (
	"io"
	"log"
	"os"
	"strconv"
	"strings"
)

type debugging bool

const debug debugging = false

func (d debugging) Printf(format string, args ...interface***REMOVED******REMOVED***) ***REMOVED***
	if d ***REMOVED***
		log.Printf(format, args...)
	***REMOVED***
***REMOVED***

const maxTok = 2048 // Largest token we can return.
const maxUint16 = 1<<16 - 1

// Tokinize a RFC 1035 zone file. The tokenizer will normalize it:
// * Add ownernames if they are left blank;
// * Suppress sequences of spaces;
// * Make each RR fit on one line (_NEWLINE is send as last)
// * Handle comments: ;
// * Handle braces - anywhere.
const (
	// Zonefile
	zEOF = iota
	zString
	zBlank
	zQuote
	zNewline
	zRrtpe
	zOwner
	zClass
	zDirOrigin   // $ORIGIN
	zDirTtl      // $TTL
	zDirInclude  // $INCLUDE
	zDirGenerate // $GENERATE

	// Privatekey file
	zValue
	zKey

	zExpectOwnerDir      // Ownername
	zExpectOwnerBl       // Whitespace after the ownername
	zExpectAny           // Expect rrtype, ttl or class
	zExpectAnyNoClass    // Expect rrtype or ttl
	zExpectAnyNoClassBl  // The whitespace after _EXPECT_ANY_NOCLASS
	zExpectAnyNoTtl      // Expect rrtype or class
	zExpectAnyNoTtlBl    // Whitespace after _EXPECT_ANY_NOTTL
	zExpectRrtype        // Expect rrtype
	zExpectRrtypeBl      // Whitespace BEFORE rrtype
	zExpectRdata         // The first element of the rdata
	zExpectDirTtlBl      // Space after directive $TTL
	zExpectDirTtl        // Directive $TTL
	zExpectDirOriginBl   // Space after directive $ORIGIN
	zExpectDirOrigin     // Directive $ORIGIN
	zExpectDirIncludeBl  // Space after directive $INCLUDE
	zExpectDirInclude    // Directive $INCLUDE
	zExpectDirGenerate   // Directive $GENERATE
	zExpectDirGenerateBl // Space after directive $GENERATE
)

// ParseError is a parsing error. It contains the parse error and the location in the io.Reader
// where the error occured.
type ParseError struct ***REMOVED***
	file string
	err  string
	lex  lex
***REMOVED***

func (e *ParseError) Error() (s string) ***REMOVED***
	if e.file != "" ***REMOVED***
		s = e.file + ": "
	***REMOVED***
	s += "dns: " + e.err + ": " + strconv.QuoteToASCII(e.lex.token) + " at line: " +
		strconv.Itoa(e.lex.line) + ":" + strconv.Itoa(e.lex.column)
	return
***REMOVED***

type lex struct ***REMOVED***
	token      string // text of the token
	tokenUpper string // uppercase text of the token
	length     int    // lenght of the token
	err        bool   // when true, token text has lexer error
	value      uint8  // value: zString, _BLANK, etc.
	line       int    // line in the file
	column     int    // column in the file
	torc       uint16 // type or class as parsed in the lexer, we only need to look this up in the grammar
	comment    string // any comment text seen
***REMOVED***

// Token holds the token that are returned when a zone file is parsed.
type Token struct ***REMOVED***
	// The scanned resource record when error is not nil.
	RR
	// When an error occured, this has the error specifics.
	Error *ParseError
	// A potential comment positioned after the RR and on the same line.
	Comment string
***REMOVED***

// NewRR reads the RR contained in the string s. Only the first RR is
// returned. If s contains no RR, return nil with no error. The class
// defaults to IN and TTL defaults to 3600. The full zone file syntax
// like $TTL, $ORIGIN, etc. is supported. All fields of the returned
// RR are set, except RR.Header().Rdlength which is set to 0.
func NewRR(s string) (RR, error) ***REMOVED***
	if len(s) > 0 && s[len(s)-1] != '\n' ***REMOVED*** // We need a closing newline
		return ReadRR(strings.NewReader(s+"\n"), "")
	***REMOVED***
	return ReadRR(strings.NewReader(s), "")
***REMOVED***

// ReadRR reads the RR contained in q.
// See NewRR for more documentation.
func ReadRR(q io.Reader, filename string) (RR, error) ***REMOVED***
	r := <-parseZoneHelper(q, ".", filename, 1)
	if r == nil ***REMOVED***
		return nil, nil
	***REMOVED***

	if r.Error != nil ***REMOVED***
		return nil, r.Error
	***REMOVED***
	return r.RR, nil
***REMOVED***

// ParseZone reads a RFC 1035 style zonefile from r. It returns *Tokens on the
// returned channel, which consist out the parsed RR, a potential comment or an error.
// If there is an error the RR is nil. The string file is only used
// in error reporting. The string origin is used as the initial origin, as
// if the file would start with: $ORIGIN origin .
// The directives $INCLUDE, $ORIGIN, $TTL and $GENERATE are supported.
// The channel t is closed by ParseZone when the end of r is reached.
//
// Basic usage pattern when reading from a string (z) containing the
// zone data:
//
//	for x := range dns.ParseZone(strings.NewReader(z), "", "") ***REMOVED***
//		if x.Error != nil ***REMOVED***
//                  // log.Println(x.Error)
//          ***REMOVED*** else ***REMOVED***
//                  // Do something with x.RR
//          ***REMOVED***
//	***REMOVED***
//
// Comments specified after an RR (and on the same line!) are returned too:
//
//	foo. IN A 10.0.0.1 ; this is a comment
//
// The text "; this is comment" is returned in Token.Comment. Comments inside the
// RR are discarded. Comments on a line by themselves are discarded too.
func ParseZone(r io.Reader, origin, file string) chan *Token ***REMOVED***
	return parseZoneHelper(r, origin, file, 10000)
***REMOVED***

func parseZoneHelper(r io.Reader, origin, file string, chansize int) chan *Token ***REMOVED***
	t := make(chan *Token, chansize)
	go parseZone(r, origin, file, t, 0)
	return t
***REMOVED***

func parseZone(r io.Reader, origin, f string, t chan *Token, include int) ***REMOVED***
	defer func() ***REMOVED***
		if include == 0 ***REMOVED***
			close(t)
		***REMOVED***
	***REMOVED***()
	s := scanInit(r)
	c := make(chan lex)
	// Start the lexer
	go zlexer(s, c)
	// 6 possible beginnings of a line, _ is a space
	// 0. zRRTYPE                              -> all omitted until the rrtype
	// 1. zOwner _ zRrtype                     -> class/ttl omitted
	// 2. zOwner _ zString _ zRrtype           -> class omitted
	// 3. zOwner _ zString _ zClass  _ zRrtype -> ttl/class
	// 4. zOwner _ zClass  _ zRrtype           -> ttl omitted
	// 5. zOwner _ zClass  _ zString _ zRrtype -> class/ttl (reversed)
	// After detecting these, we know the zRrtype so we can jump to functions
	// handling the rdata for each of these types.

	if origin == "" ***REMOVED***
		origin = "."
	***REMOVED***
	origin = Fqdn(origin)
	if _, ok := IsDomainName(origin); !ok ***REMOVED***
		t <- &Token***REMOVED***Error: &ParseError***REMOVED***f, "bad initial origin name", lex***REMOVED******REMOVED******REMOVED******REMOVED***
		return
	***REMOVED***

	st := zExpectOwnerDir // initial state
	var h RR_Header
	var defttl uint32 = defaultTtl
	var prevName string
	for l := range c ***REMOVED***
		// Lexer spotted an error already
		if l.err == true ***REMOVED***
			t <- &Token***REMOVED***Error: &ParseError***REMOVED***f, l.token, l***REMOVED******REMOVED***
			return

		***REMOVED***
		switch st ***REMOVED***
		case zExpectOwnerDir:
			// We can also expect a directive, like $TTL or $ORIGIN
			h.Ttl = defttl
			h.Class = ClassINET
			switch l.value ***REMOVED***
			case zNewline:
				st = zExpectOwnerDir
			case zOwner:
				h.Name = l.token
				if l.token[0] == '@' ***REMOVED***
					h.Name = origin
					prevName = h.Name
					st = zExpectOwnerBl
					break
				***REMOVED***
				if h.Name[l.length-1] != '.' ***REMOVED***
					h.Name = appendOrigin(h.Name, origin)
				***REMOVED***
				_, ok := IsDomainName(l.token)
				if !ok ***REMOVED***
					t <- &Token***REMOVED***Error: &ParseError***REMOVED***f, "bad owner name", l***REMOVED******REMOVED***
					return
				***REMOVED***
				prevName = h.Name
				st = zExpectOwnerBl
			case zDirTtl:
				st = zExpectDirTtlBl
			case zDirOrigin:
				st = zExpectDirOriginBl
			case zDirInclude:
				st = zExpectDirIncludeBl
			case zDirGenerate:
				st = zExpectDirGenerateBl
			case zRrtpe:
				h.Name = prevName
				h.Rrtype = l.torc
				st = zExpectRdata
			case zClass:
				h.Name = prevName
				h.Class = l.torc
				st = zExpectAnyNoClassBl
			case zBlank:
				// Discard, can happen when there is nothing on the
				// line except the RR type
			case zString:
				ttl, ok := stringToTtl(l.token)
				if !ok ***REMOVED***
					t <- &Token***REMOVED***Error: &ParseError***REMOVED***f, "not a TTL", l***REMOVED******REMOVED***
					return
				***REMOVED***
				h.Ttl = ttl
				// Don't about the defttl, we should take the $TTL value
				// defttl = ttl
				st = zExpectAnyNoTtlBl

			default:
				t <- &Token***REMOVED***Error: &ParseError***REMOVED***f, "syntax error at beginning", l***REMOVED******REMOVED***
				return
			***REMOVED***
		case zExpectDirIncludeBl:
			if l.value != zBlank ***REMOVED***
				t <- &Token***REMOVED***Error: &ParseError***REMOVED***f, "no blank after $INCLUDE-directive", l***REMOVED******REMOVED***
				return
			***REMOVED***
			st = zExpectDirInclude
		case zExpectDirInclude:
			if l.value != zString ***REMOVED***
				t <- &Token***REMOVED***Error: &ParseError***REMOVED***f, "expecting $INCLUDE value, not this...", l***REMOVED******REMOVED***
				return
			***REMOVED***
			neworigin := origin // There may be optionally a new origin set after the filename, if not use current one
			l := <-c
			switch l.value ***REMOVED***
			case zBlank:
				l := <-c
				if l.value == zString ***REMOVED***
					if _, ok := IsDomainName(l.token); !ok || l.length == 0 || l.err ***REMOVED***
						t <- &Token***REMOVED***Error: &ParseError***REMOVED***f, "bad origin name", l***REMOVED******REMOVED***
						return
					***REMOVED***
					// a new origin is specified.
					if l.token[l.length-1] != '.' ***REMOVED***
						if origin != "." ***REMOVED*** // Prevent .. endings
							neworigin = l.token + "." + origin
						***REMOVED*** else ***REMOVED***
							neworigin = l.token + origin
						***REMOVED***
					***REMOVED*** else ***REMOVED***
						neworigin = l.token
					***REMOVED***
				***REMOVED***
			case zNewline, zEOF:
				// Ok
			default:
				t <- &Token***REMOVED***Error: &ParseError***REMOVED***f, "garbage after $INCLUDE", l***REMOVED******REMOVED***
				return
			***REMOVED***
			// Start with the new file
			r1, e1 := os.Open(l.token)
			if e1 != nil ***REMOVED***
				t <- &Token***REMOVED***Error: &ParseError***REMOVED***f, "failed to open `" + l.token + "'", l***REMOVED******REMOVED***
				return
			***REMOVED***
			if include+1 > 7 ***REMOVED***
				t <- &Token***REMOVED***Error: &ParseError***REMOVED***f, "too deeply nested $INCLUDE", l***REMOVED******REMOVED***
				return
			***REMOVED***
			parseZone(r1, l.token, neworigin, t, include+1)
			st = zExpectOwnerDir
		case zExpectDirTtlBl:
			if l.value != zBlank ***REMOVED***
				t <- &Token***REMOVED***Error: &ParseError***REMOVED***f, "no blank after $TTL-directive", l***REMOVED******REMOVED***
				return
			***REMOVED***
			st = zExpectDirTtl
		case zExpectDirTtl:
			if l.value != zString ***REMOVED***
				t <- &Token***REMOVED***Error: &ParseError***REMOVED***f, "expecting $TTL value, not this...", l***REMOVED******REMOVED***
				return
			***REMOVED***
			if e, _ := slurpRemainder(c, f); e != nil ***REMOVED***
				t <- &Token***REMOVED***Error: e***REMOVED***
				return
			***REMOVED***
			ttl, ok := stringToTtl(l.token)
			if !ok ***REMOVED***
				t <- &Token***REMOVED***Error: &ParseError***REMOVED***f, "expecting $TTL value, not this...", l***REMOVED******REMOVED***
				return
			***REMOVED***
			defttl = ttl
			st = zExpectOwnerDir
		case zExpectDirOriginBl:
			if l.value != zBlank ***REMOVED***
				t <- &Token***REMOVED***Error: &ParseError***REMOVED***f, "no blank after $ORIGIN-directive", l***REMOVED******REMOVED***
				return
			***REMOVED***
			st = zExpectDirOrigin
		case zExpectDirOrigin:
			if l.value != zString ***REMOVED***
				t <- &Token***REMOVED***Error: &ParseError***REMOVED***f, "expecting $ORIGIN value, not this...", l***REMOVED******REMOVED***
				return
			***REMOVED***
			if e, _ := slurpRemainder(c, f); e != nil ***REMOVED***
				t <- &Token***REMOVED***Error: e***REMOVED***
			***REMOVED***
			if _, ok := IsDomainName(l.token); !ok ***REMOVED***
				t <- &Token***REMOVED***Error: &ParseError***REMOVED***f, "bad origin name", l***REMOVED******REMOVED***
				return
			***REMOVED***
			if l.token[l.length-1] != '.' ***REMOVED***
				if origin != "." ***REMOVED*** // Prevent .. endings
					origin = l.token + "." + origin
				***REMOVED*** else ***REMOVED***
					origin = l.token + origin
				***REMOVED***
			***REMOVED*** else ***REMOVED***
				origin = l.token
			***REMOVED***
			st = zExpectOwnerDir
		case zExpectDirGenerateBl:
			if l.value != zBlank ***REMOVED***
				t <- &Token***REMOVED***Error: &ParseError***REMOVED***f, "no blank after $GENERATE-directive", l***REMOVED******REMOVED***
				return
			***REMOVED***
			st = zExpectDirGenerate
		case zExpectDirGenerate:
			if l.value != zString ***REMOVED***
				t <- &Token***REMOVED***Error: &ParseError***REMOVED***f, "expecting $GENERATE value, not this...", l***REMOVED******REMOVED***
				return
			***REMOVED***
			if e := generate(l, c, t, origin); e != "" ***REMOVED***
				t <- &Token***REMOVED***Error: &ParseError***REMOVED***f, e, l***REMOVED******REMOVED***
				return
			***REMOVED***
			st = zExpectOwnerDir
		case zExpectOwnerBl:
			if l.value != zBlank ***REMOVED***
				t <- &Token***REMOVED***Error: &ParseError***REMOVED***f, "no blank after owner", l***REMOVED******REMOVED***
				return
			***REMOVED***
			st = zExpectAny
		case zExpectAny:
			switch l.value ***REMOVED***
			case zRrtpe:
				h.Rrtype = l.torc
				st = zExpectRdata
			case zClass:
				h.Class = l.torc
				st = zExpectAnyNoClassBl
			case zString:
				ttl, ok := stringToTtl(l.token)
				if !ok ***REMOVED***
					t <- &Token***REMOVED***Error: &ParseError***REMOVED***f, "not a TTL", l***REMOVED******REMOVED***
					return
				***REMOVED***
				h.Ttl = ttl
				// defttl = ttl // don't set the defttl here
				st = zExpectAnyNoTtlBl
			default:
				t <- &Token***REMOVED***Error: &ParseError***REMOVED***f, "expecting RR type, TTL or class, not this...", l***REMOVED******REMOVED***
				return
			***REMOVED***
		case zExpectAnyNoClassBl:
			if l.value != zBlank ***REMOVED***
				t <- &Token***REMOVED***Error: &ParseError***REMOVED***f, "no blank before class", l***REMOVED******REMOVED***
				return
			***REMOVED***
			st = zExpectAnyNoClass
		case zExpectAnyNoTtlBl:
			if l.value != zBlank ***REMOVED***
				t <- &Token***REMOVED***Error: &ParseError***REMOVED***f, "no blank before TTL", l***REMOVED******REMOVED***
				return
			***REMOVED***
			st = zExpectAnyNoTtl
		case zExpectAnyNoTtl:
			switch l.value ***REMOVED***
			case zClass:
				h.Class = l.torc
				st = zExpectRrtypeBl
			case zRrtpe:
				h.Rrtype = l.torc
				st = zExpectRdata
			default:
				t <- &Token***REMOVED***Error: &ParseError***REMOVED***f, "expecting RR type or class, not this...", l***REMOVED******REMOVED***
				return
			***REMOVED***
		case zExpectAnyNoClass:
			switch l.value ***REMOVED***
			case zString:
				ttl, ok := stringToTtl(l.token)
				if !ok ***REMOVED***
					t <- &Token***REMOVED***Error: &ParseError***REMOVED***f, "not a TTL", l***REMOVED******REMOVED***
					return
				***REMOVED***
				h.Ttl = ttl
				// defttl = ttl // don't set the def ttl anymore
				st = zExpectRrtypeBl
			case zRrtpe:
				h.Rrtype = l.torc
				st = zExpectRdata
			default:
				t <- &Token***REMOVED***Error: &ParseError***REMOVED***f, "expecting RR type or TTL, not this...", l***REMOVED******REMOVED***
				return
			***REMOVED***
		case zExpectRrtypeBl:
			if l.value != zBlank ***REMOVED***
				t <- &Token***REMOVED***Error: &ParseError***REMOVED***f, "no blank before RR type", l***REMOVED******REMOVED***
				return
			***REMOVED***
			st = zExpectRrtype
		case zExpectRrtype:
			if l.value != zRrtpe ***REMOVED***
				t <- &Token***REMOVED***Error: &ParseError***REMOVED***f, "unknown RR type", l***REMOVED******REMOVED***
				return
			***REMOVED***
			h.Rrtype = l.torc
			st = zExpectRdata
		case zExpectRdata:
			r, e, c1 := setRR(h, c, origin, f)
			if e != nil ***REMOVED***
				// If e.lex is nil than we have encounter a unknown RR type
				// in that case we substitute our current lex token
				if e.lex.token == "" && e.lex.value == 0 ***REMOVED***
					e.lex = l // Uh, dirty
				***REMOVED***
				t <- &Token***REMOVED***Error: e***REMOVED***
				return
			***REMOVED***
			t <- &Token***REMOVED***RR: r, Comment: c1***REMOVED***
			st = zExpectOwnerDir
		***REMOVED***
	***REMOVED***
	// If we get here, we and the h.Rrtype is still zero, we haven't parsed anything, this
	// is not an error, because an empty zone file is still a zone file.
***REMOVED***

// zlexer scans the sourcefile and returns tokens on the channel c.
func zlexer(s *scan, c chan lex) ***REMOVED***
	var l lex
	str := make([]byte, maxTok) // Should be enough for any token
	stri := 0                   // Offset in str (0 means empty)
	com := make([]byte, maxTok) // Hold comment text
	comi := 0
	quote := false
	escape := false
	space := false
	commt := false
	rrtype := false
	owner := true
	brace := 0
	x, err := s.tokenText()
	defer close(c)
	for err == nil ***REMOVED***
		l.column = s.position.Column
		l.line = s.position.Line
		if stri >= maxTok ***REMOVED***
			l.token = "token length insufficient for parsing"
			l.err = true
			debug.Printf("[%+v]", l.token)
			c <- l
			return
		***REMOVED***
		if comi >= maxTok ***REMOVED***
			l.token = "comment length insufficient for parsing"
			l.err = true
			debug.Printf("[%+v]", l.token)
			c <- l
			return
		***REMOVED***

		switch x ***REMOVED***
		case ' ', '\t':
			if escape ***REMOVED***
				escape = false
				str[stri] = x
				stri++
				break
			***REMOVED***
			if quote ***REMOVED***
				// Inside quotes this is legal
				str[stri] = x
				stri++
				break
			***REMOVED***
			if commt ***REMOVED***
				com[comi] = x
				comi++
				break
			***REMOVED***
			if stri == 0 ***REMOVED***
				// Space directly in the beginning, handled in the grammar
			***REMOVED*** else if owner ***REMOVED***
				// If we have a string and its the first, make it an owner
				l.value = zOwner
				l.token = string(str[:stri])
				l.tokenUpper = strings.ToUpper(l.token)
				l.length = stri
				// escape $... start with a \ not a $, so this will work
				switch l.tokenUpper ***REMOVED***
				case "$TTL":
					l.value = zDirTtl
				case "$ORIGIN":
					l.value = zDirOrigin
				case "$INCLUDE":
					l.value = zDirInclude
				case "$GENERATE":
					l.value = zDirGenerate
				***REMOVED***
				debug.Printf("[7 %+v]", l.token)
				c <- l
			***REMOVED*** else ***REMOVED***
				l.value = zString
				l.token = string(str[:stri])
				l.tokenUpper = strings.ToUpper(l.token)
				l.length = stri
				if !rrtype ***REMOVED***
					if t, ok := StringToType[l.tokenUpper]; ok ***REMOVED***
						l.value = zRrtpe
						l.torc = t
						rrtype = true
					***REMOVED*** else ***REMOVED***
						if strings.HasPrefix(l.tokenUpper, "TYPE") ***REMOVED***
							t, ok := typeToInt(l.token)
							if !ok ***REMOVED***
								l.token = "unknown RR type"
								l.err = true
								c <- l
								return
							***REMOVED***
							l.value = zRrtpe
							l.torc = t
						***REMOVED***
					***REMOVED***
					if t, ok := StringToClass[l.tokenUpper]; ok ***REMOVED***
						l.value = zClass
						l.torc = t
					***REMOVED*** else ***REMOVED***
						if strings.HasPrefix(l.tokenUpper, "CLASS") ***REMOVED***
							t, ok := classToInt(l.token)
							if !ok ***REMOVED***
								l.token = "unknown class"
								l.err = true
								c <- l
								return
							***REMOVED***
							l.value = zClass
							l.torc = t
						***REMOVED***
					***REMOVED***
				***REMOVED***
				debug.Printf("[6 %+v]", l.token)
				c <- l
			***REMOVED***
			stri = 0
			// I reverse space stuff here
			if !space && !commt ***REMOVED***
				l.value = zBlank
				l.token = " "
				l.length = 1
				debug.Printf("[5 %+v]", l.token)
				c <- l
			***REMOVED***
			owner = false
			space = true
		case ';':
			if escape ***REMOVED***
				escape = false
				str[stri] = x
				stri++
				break
			***REMOVED***
			if quote ***REMOVED***
				// Inside quotes this is legal
				str[stri] = x
				stri++
				break
			***REMOVED***
			if stri > 0 ***REMOVED***
				l.value = zString
				l.token = string(str[:stri])
				l.length = stri
				debug.Printf("[4 %+v]", l.token)
				c <- l
				stri = 0
			***REMOVED***
			commt = true
			com[comi] = ';'
			comi++
		case '\r':
			escape = false
			if quote ***REMOVED***
				str[stri] = x
				stri++
				break
			***REMOVED***
			// discard if outside of quotes
		case '\n':
			escape = false
			// Escaped newline
			if quote ***REMOVED***
				str[stri] = x
				stri++
				break
			***REMOVED***
			// inside quotes this is legal
			if commt ***REMOVED***
				// Reset a comment
				commt = false
				rrtype = false
				stri = 0
				// If not in a brace this ends the comment AND the RR
				if brace == 0 ***REMOVED***
					owner = true
					owner = true
					l.value = zNewline
					l.token = "\n"
					l.length = 1
					l.comment = string(com[:comi])
					debug.Printf("[3 %+v %+v]", l.token, l.comment)
					c <- l
					l.comment = ""
					comi = 0
					break
				***REMOVED***
				com[comi] = ' ' // convert newline to space
				comi++
				break
			***REMOVED***

			if brace == 0 ***REMOVED***
				// If there is previous text, we should output it here
				if stri != 0 ***REMOVED***
					l.value = zString
					l.token = string(str[:stri])
					l.tokenUpper = strings.ToUpper(l.token)

					l.length = stri
					if !rrtype ***REMOVED***
						if t, ok := StringToType[l.tokenUpper]; ok ***REMOVED***
							l.value = zRrtpe
							l.torc = t
							rrtype = true
						***REMOVED***
					***REMOVED***
					debug.Printf("[2 %+v]", l.token)
					c <- l
				***REMOVED***
				l.value = zNewline
				l.token = "\n"
				l.length = 1
				debug.Printf("[1 %+v]", l.token)
				c <- l
				stri = 0
				commt = false
				rrtype = false
				owner = true
				comi = 0
			***REMOVED***
		case '\\':
			// comments do not get escaped chars, everything is copied
			if commt ***REMOVED***
				com[comi] = x
				comi++
				break
			***REMOVED***
			// something already escaped must be in string
			if escape ***REMOVED***
				str[stri] = x
				stri++
				escape = false
				break
			***REMOVED***
			// something escaped outside of string gets added to string
			str[stri] = x
			stri++
			escape = true
		case '"':
			if commt ***REMOVED***
				com[comi] = x
				comi++
				break
			***REMOVED***
			if escape ***REMOVED***
				str[stri] = x
				stri++
				escape = false
				break
			***REMOVED***
			space = false
			// send previous gathered text and the quote
			if stri != 0 ***REMOVED***
				l.value = zString
				l.token = string(str[:stri])
				l.length = stri

				debug.Printf("[%+v]", l.token)
				c <- l
				stri = 0
			***REMOVED***

			// send quote itself as separate token
			l.value = zQuote
			l.token = "\""
			l.length = 1
			c <- l
			quote = !quote
		case '(', ')':
			if commt ***REMOVED***
				com[comi] = x
				comi++
				break
			***REMOVED***
			if escape ***REMOVED***
				str[stri] = x
				stri++
				escape = false
				break
			***REMOVED***
			if quote ***REMOVED***
				str[stri] = x
				stri++
				break
			***REMOVED***
			switch x ***REMOVED***
			case ')':
				brace--
				if brace < 0 ***REMOVED***
					l.token = "extra closing brace"
					l.err = true
					debug.Printf("[%+v]", l.token)
					c <- l
					return
				***REMOVED***
			case '(':
				brace++
			***REMOVED***
		default:
			escape = false
			if commt ***REMOVED***
				com[comi] = x
				comi++
				break
			***REMOVED***
			str[stri] = x
			stri++
			space = false
		***REMOVED***
		x, err = s.tokenText()
	***REMOVED***
	if stri > 0 ***REMOVED***
		// Send remainder
		l.token = string(str[:stri])
		l.length = stri
		l.value = zString
		debug.Printf("[%+v]", l.token)
		c <- l
	***REMOVED***
***REMOVED***

// Extract the class number from CLASSxx
func classToInt(token string) (uint16, bool) ***REMOVED***
	offset := 5
	if len(token) < offset+1 ***REMOVED***
		return 0, false
	***REMOVED***
	class, ok := strconv.Atoi(token[offset:])
	if ok != nil || class > maxUint16 ***REMOVED***
		return 0, false
	***REMOVED***
	return uint16(class), true
***REMOVED***

// Extract the rr number from TYPExxx
func typeToInt(token string) (uint16, bool) ***REMOVED***
	offset := 4
	if len(token) < offset+1 ***REMOVED***
		return 0, false
	***REMOVED***
	typ, ok := strconv.Atoi(token[offset:])
	if ok != nil || typ > maxUint16 ***REMOVED***
		return 0, false
	***REMOVED***
	return uint16(typ), true
***REMOVED***

// Parse things like 2w, 2m, etc, Return the time in seconds.
func stringToTtl(token string) (uint32, bool) ***REMOVED***
	s := uint32(0)
	i := uint32(0)
	for _, c := range token ***REMOVED***
		switch c ***REMOVED***
		case 's', 'S':
			s += i
			i = 0
		case 'm', 'M':
			s += i * 60
			i = 0
		case 'h', 'H':
			s += i * 60 * 60
			i = 0
		case 'd', 'D':
			s += i * 60 * 60 * 24
			i = 0
		case 'w', 'W':
			s += i * 60 * 60 * 24 * 7
			i = 0
		case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
			i *= 10
			i += uint32(c) - '0'
		default:
			return 0, false
		***REMOVED***
	***REMOVED***
	return s + i, true
***REMOVED***

// Parse LOC records' <digits>[.<digits>][mM] into a
// mantissa exponent format. Token should contain the entire
// string (i.e. no spaces allowed)
func stringToCm(token string) (e, m uint8, ok bool) ***REMOVED***
	if token[len(token)-1] == 'M' || token[len(token)-1] == 'm' ***REMOVED***
		token = token[0 : len(token)-1]
	***REMOVED***
	s := strings.SplitN(token, ".", 2)
	var meters, cmeters, val int
	var err error
	switch len(s) ***REMOVED***
	case 2:
		if cmeters, err = strconv.Atoi(s[1]); err != nil ***REMOVED***
			return
		***REMOVED***
		fallthrough
	case 1:
		if meters, err = strconv.Atoi(s[0]); err != nil ***REMOVED***
			return
		***REMOVED***
	case 0:
		// huh?
		return 0, 0, false
	***REMOVED***
	ok = true
	if meters > 0 ***REMOVED***
		e = 2
		val = meters
	***REMOVED*** else ***REMOVED***
		e = 0
		val = cmeters
	***REMOVED***
	for val > 10 ***REMOVED***
		e++
		val /= 10
	***REMOVED***
	if e > 9 ***REMOVED***
		ok = false
	***REMOVED***
	m = uint8(val)
	return
***REMOVED***

func appendOrigin(name, origin string) string ***REMOVED***
	if origin == "." ***REMOVED***
		return name + origin
	***REMOVED***
	return name + "." + origin
***REMOVED***

// LOC record helper function
func locCheckNorth(token string, latitude uint32) (uint32, bool) ***REMOVED***
	switch token ***REMOVED***
	case "n", "N":
		return LOC_EQUATOR + latitude, true
	case "s", "S":
		return LOC_EQUATOR - latitude, true
	***REMOVED***
	return latitude, false
***REMOVED***

// LOC record helper function
func locCheckEast(token string, longitude uint32) (uint32, bool) ***REMOVED***
	switch token ***REMOVED***
	case "e", "E":
		return LOC_EQUATOR + longitude, true
	case "w", "W":
		return LOC_EQUATOR - longitude, true
	***REMOVED***
	return longitude, false
***REMOVED***

// "Eat" the rest of the "line". Return potential comments
func slurpRemainder(c chan lex, f string) (*ParseError, string) ***REMOVED***
	l := <-c
	com := ""
	switch l.value ***REMOVED***
	case zBlank:
		l = <-c
		com = l.comment
		if l.value != zNewline && l.value != zEOF ***REMOVED***
			return &ParseError***REMOVED***f, "garbage after rdata", l***REMOVED***, ""
		***REMOVED***
	case zNewline:
		com = l.comment
	case zEOF:
	default:
		return &ParseError***REMOVED***f, "garbage after rdata", l***REMOVED***, ""
	***REMOVED***
	return nil, com
***REMOVED***

// Parse a 64 bit-like ipv6 address: "0014:4fff:ff20:ee64"
// Used for NID and L64 record.
func stringToNodeID(l lex) (uint64, *ParseError) ***REMOVED***
	if len(l.token) < 19 ***REMOVED***
		return 0, &ParseError***REMOVED***l.token, "bad NID/L64 NodeID/Locator64", l***REMOVED***
	***REMOVED***
	// There must be three colons at fixes postitions, if not its a parse error
	if l.token[4] != ':' && l.token[9] != ':' && l.token[14] != ':' ***REMOVED***
		return 0, &ParseError***REMOVED***l.token, "bad NID/L64 NodeID/Locator64", l***REMOVED***
	***REMOVED***
	s := l.token[0:4] + l.token[5:9] + l.token[10:14] + l.token[15:19]
	u, e := strconv.ParseUint(s, 16, 64)
	if e != nil ***REMOVED***
		return 0, &ParseError***REMOVED***l.token, "bad NID/L64 NodeID/Locator64", l***REMOVED***
	***REMOVED***
	return u, nil
***REMOVED***
