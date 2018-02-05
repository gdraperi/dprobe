// Copyright 2017 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package format

import (
	"reflect"
	"unicode/utf8"
)

// A Parser parses a format string. The result from the parse are set in the
// struct fields.
type Parser struct ***REMOVED***
	Verb rune

	WidthPresent bool
	PrecPresent  bool
	Minus        bool
	Plus         bool
	Sharp        bool
	Space        bool
	Zero         bool

	// For the formats %+v %#v, we set the plusV/sharpV flags
	// and clear the plus/sharp flags since %+v and %#v are in effect
	// different, flagless formats set at the top level.
	PlusV  bool
	SharpV bool

	HasIndex bool

	Width int
	Prec  int // precision

	// retain arguments across calls.
	Args []interface***REMOVED******REMOVED***
	// retain current argument number across calls
	ArgNum int

	// reordered records whether the format string used argument reordering.
	Reordered bool
	// goodArgNum records whether the most recent reordering directive was valid.
	goodArgNum bool

	// position info
	format   string
	startPos int
	endPos   int
	Status   Status
***REMOVED***

// Reset initializes a parser to scan format strings for the given args.
func (p *Parser) Reset(args []interface***REMOVED******REMOVED***) ***REMOVED***
	p.Args = args
	p.ArgNum = 0
	p.startPos = 0
	p.Reordered = false
***REMOVED***

// Text returns the part of the format string that was parsed by the last call
// to Scan. It returns the original substitution clause if the current scan
// parsed a substitution.
func (p *Parser) Text() string ***REMOVED*** return p.format[p.startPos:p.endPos] ***REMOVED***

// SetFormat sets a new format string to parse. It does not reset the argument
// count.
func (p *Parser) SetFormat(format string) ***REMOVED***
	p.format = format
	p.startPos = 0
	p.endPos = 0
***REMOVED***

// Status indicates the result type of a call to Scan.
type Status int

const (
	StatusText Status = iota
	StatusSubstitution
	StatusBadWidthSubstitution
	StatusBadPrecSubstitution
	StatusNoVerb
	StatusBadArgNum
	StatusMissingArg
)

// ClearFlags reset the parser to default behavior.
func (p *Parser) ClearFlags() ***REMOVED***
	p.WidthPresent = false
	p.PrecPresent = false
	p.Minus = false
	p.Plus = false
	p.Sharp = false
	p.Space = false
	p.Zero = false

	p.PlusV = false
	p.SharpV = false

	p.HasIndex = false
***REMOVED***

// Scan scans the next part of the format string and sets the status to
// indicate whether it scanned a string literal, substitution or error.
func (p *Parser) Scan() bool ***REMOVED***
	p.Status = StatusText
	format := p.format
	end := len(format)
	if p.endPos >= end ***REMOVED***
		return false
	***REMOVED***
	afterIndex := false // previous item in format was an index like [3].

	p.startPos = p.endPos
	p.goodArgNum = true
	i := p.startPos
	for i < end && format[i] != '%' ***REMOVED***
		i++
	***REMOVED***
	if i > p.startPos ***REMOVED***
		p.endPos = i
		return true
	***REMOVED***
	// Process one verb
	i++

	p.Status = StatusSubstitution

	// Do we have flags?
	p.ClearFlags()

simpleFormat:
	for ; i < end; i++ ***REMOVED***
		c := p.format[i]
		switch c ***REMOVED***
		case '#':
			p.Sharp = true
		case '0':
			p.Zero = !p.Minus // Only allow zero padding to the left.
		case '+':
			p.Plus = true
		case '-':
			p.Minus = true
			p.Zero = false // Do not pad with zeros to the right.
		case ' ':
			p.Space = true
		default:
			// Fast path for common case of ascii lower case simple verbs
			// without precision or width or argument indices.
			if 'a' <= c && c <= 'z' && p.ArgNum < len(p.Args) ***REMOVED***
				if c == 'v' ***REMOVED***
					// Go syntax
					p.SharpV = p.Sharp
					p.Sharp = false
					// Struct-field syntax
					p.PlusV = p.Plus
					p.Plus = false
				***REMOVED***
				p.Verb = rune(c)
				p.ArgNum++
				p.endPos = i + 1
				return true
			***REMOVED***
			// Format is more complex than simple flags and a verb or is malformed.
			break simpleFormat
		***REMOVED***
	***REMOVED***

	// Do we have an explicit argument index?
	i, afterIndex = p.updateArgNumber(format, i)

	// Do we have width?
	if i < end && format[i] == '*' ***REMOVED***
		i++
		p.Width, p.WidthPresent = p.intFromArg()

		if !p.WidthPresent ***REMOVED***
			p.Status = StatusBadWidthSubstitution
		***REMOVED***

		// We have a negative width, so take its value and ensure
		// that the minus flag is set
		if p.Width < 0 ***REMOVED***
			p.Width = -p.Width
			p.Minus = true
			p.Zero = false // Do not pad with zeros to the right.
		***REMOVED***
		afterIndex = false
	***REMOVED*** else ***REMOVED***
		p.Width, p.WidthPresent, i = parsenum(format, i, end)
		if afterIndex && p.WidthPresent ***REMOVED*** // "%[3]2d"
			p.goodArgNum = false
		***REMOVED***
	***REMOVED***

	// Do we have precision?
	if i+1 < end && format[i] == '.' ***REMOVED***
		i++
		if afterIndex ***REMOVED*** // "%[3].2d"
			p.goodArgNum = false
		***REMOVED***
		i, afterIndex = p.updateArgNumber(format, i)
		if i < end && format[i] == '*' ***REMOVED***
			i++
			p.Prec, p.PrecPresent = p.intFromArg()
			// Negative precision arguments don't make sense
			if p.Prec < 0 ***REMOVED***
				p.Prec = 0
				p.PrecPresent = false
			***REMOVED***
			if !p.PrecPresent ***REMOVED***
				p.Status = StatusBadPrecSubstitution
			***REMOVED***
			afterIndex = false
		***REMOVED*** else ***REMOVED***
			p.Prec, p.PrecPresent, i = parsenum(format, i, end)
			if !p.PrecPresent ***REMOVED***
				p.Prec = 0
				p.PrecPresent = true
			***REMOVED***
		***REMOVED***
	***REMOVED***

	if !afterIndex ***REMOVED***
		i, afterIndex = p.updateArgNumber(format, i)
	***REMOVED***
	p.HasIndex = afterIndex

	if i >= end ***REMOVED***
		p.endPos = i
		p.Status = StatusNoVerb
		return true
	***REMOVED***

	verb, w := utf8.DecodeRuneInString(format[i:])
	p.endPos = i + w
	p.Verb = verb

	switch ***REMOVED***
	case verb == '%': // Percent does not absorb operands and ignores f.wid and f.prec.
		p.startPos = p.endPos - 1
		p.Status = StatusText
	case !p.goodArgNum:
		p.Status = StatusBadArgNum
	case p.ArgNum >= len(p.Args): // No argument left over to print for the current verb.
		p.Status = StatusMissingArg
	case verb == 'v':
		// Go syntax
		p.SharpV = p.Sharp
		p.Sharp = false
		// Struct-field syntax
		p.PlusV = p.Plus
		p.Plus = false
		fallthrough
	default:
		p.ArgNum++
	***REMOVED***
	return true
***REMOVED***

// intFromArg gets the ArgNumth element of Args. On return, isInt reports
// whether the argument has integer type.
func (p *Parser) intFromArg() (num int, isInt bool) ***REMOVED***
	if p.ArgNum < len(p.Args) ***REMOVED***
		arg := p.Args[p.ArgNum]
		num, isInt = arg.(int) // Almost always OK.
		if !isInt ***REMOVED***
			// Work harder.
			switch v := reflect.ValueOf(arg); v.Kind() ***REMOVED***
			case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
				n := v.Int()
				if int64(int(n)) == n ***REMOVED***
					num = int(n)
					isInt = true
				***REMOVED***
			case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
				n := v.Uint()
				if int64(n) >= 0 && uint64(int(n)) == n ***REMOVED***
					num = int(n)
					isInt = true
				***REMOVED***
			default:
				// Already 0, false.
			***REMOVED***
		***REMOVED***
		p.ArgNum++
		if tooLarge(num) ***REMOVED***
			num = 0
			isInt = false
		***REMOVED***
	***REMOVED***
	return
***REMOVED***

// parseArgNumber returns the value of the bracketed number, minus 1
// (explicit argument numbers are one-indexed but we want zero-indexed).
// The opening bracket is known to be present at format[0].
// The returned values are the index, the number of bytes to consume
// up to the closing paren, if present, and whether the number parsed
// ok. The bytes to consume will be 1 if no closing paren is present.
func parseArgNumber(format string) (index int, wid int, ok bool) ***REMOVED***
	// There must be at least 3 bytes: [n].
	if len(format) < 3 ***REMOVED***
		return 0, 1, false
	***REMOVED***

	// Find closing bracket.
	for i := 1; i < len(format); i++ ***REMOVED***
		if format[i] == ']' ***REMOVED***
			width, ok, newi := parsenum(format, 1, i)
			if !ok || newi != i ***REMOVED***
				return 0, i + 1, false
			***REMOVED***
			return width - 1, i + 1, true // arg numbers are one-indexed and skip paren.
		***REMOVED***
	***REMOVED***
	return 0, 1, false
***REMOVED***

// updateArgNumber returns the next argument to evaluate, which is either the value of the passed-in
// argNum or the value of the bracketed integer that begins format[i:]. It also returns
// the new value of i, that is, the index of the next byte of the format to process.
func (p *Parser) updateArgNumber(format string, i int) (newi int, found bool) ***REMOVED***
	if len(format) <= i || format[i] != '[' ***REMOVED***
		return i, false
	***REMOVED***
	p.Reordered = true
	index, wid, ok := parseArgNumber(format[i:])
	if ok && 0 <= index && index < len(p.Args) ***REMOVED***
		p.ArgNum = index
		return i + wid, true
	***REMOVED***
	p.goodArgNum = false
	return i + wid, ok
***REMOVED***

// tooLarge reports whether the magnitude of the integer is
// too large to be used as a formatting width or precision.
func tooLarge(x int) bool ***REMOVED***
	const max int = 1e6
	return x > max || x < -max
***REMOVED***

// parsenum converts ASCII to integer.  num is 0 (and isnum is false) if no number present.
func parsenum(s string, start, end int) (num int, isnum bool, newi int) ***REMOVED***
	if start >= end ***REMOVED***
		return 0, false, end
	***REMOVED***
	for newi = start; newi < end && '0' <= s[newi] && s[newi] <= '9'; newi++ ***REMOVED***
		if tooLarge(num) ***REMOVED***
			return 0, false, end // Overflow; crazy long number most likely.
		***REMOVED***
		num = num*10 + int(s[newi]-'0')
		isnum = true
	***REMOVED***
	return
***REMOVED***
