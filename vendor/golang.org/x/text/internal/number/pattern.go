// Copyright 2015 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package number

import (
	"errors"
	"unicode/utf8"
)

// This file contains a parser for the CLDR number patterns as described in
// http://unicode.org/reports/tr35/tr35-numbers.html#Number_Format_Patterns.
//
// The following BNF is derived from this standard.
//
// pattern    := subpattern (';' subpattern)?
// subpattern := affix? number exponent? affix?
// number     := decimal | sigDigits
// decimal    := '#'* '0'* ('.' fraction)? | '#' | '0'
// fraction   := '0'* '#'*
// sigDigits  := '#'* '@' '@'* '#'*
// exponent   := 'E' '+'? '0'* '0'
// padSpec    := '*' \L
//
// Notes:
// - An affix pattern may contain any runes, but runes with special meaning
//   should be escaped.
// - Sequences of digits, '#', and '@' in decimal and sigDigits may have
//   interstitial commas.

// TODO: replace special characters in affixes (-, +, ¤) with control codes.

// Pattern holds information for formatting numbers. It is designed to hold
// information from CLDR number patterns.
//
// This pattern is precompiled  for all patterns for all languages. Even though
// the number of patterns is not very large, we want to keep this small.
//
// This type is only intended for internal use.
type Pattern struct ***REMOVED***
	RoundingContext

	Affix       string // includes prefix and suffix. First byte is prefix length.
	Offset      uint16 // Offset into Affix for prefix and suffix
	NegOffset   uint16 // Offset into Affix for negative prefix and suffix or 0.
	PadRune     rune
	FormatWidth uint16

	GroupingSize [2]uint8
	Flags        PatternFlag
***REMOVED***

// A RoundingContext indicates how a number should be converted to digits.
// It contains all information needed to determine the "visible digits" as
// required by the pluralization rules.
type RoundingContext struct ***REMOVED***
	// TODO: unify these two fields so that there is a more unambiguous meaning
	// of how precision is handled.
	MaxSignificantDigits int16 // -1 is unlimited
	MaxFractionDigits    int16 // -1 is unlimited

	Increment      uint32
	IncrementScale uint8 // May differ from printed scale.

	Mode RoundingMode

	DigitShift uint8 // Number of decimals to shift. Used for % and ‰.

	// Number of digits.
	MinIntegerDigits uint8

	MaxIntegerDigits     uint8
	MinFractionDigits    uint8
	MinSignificantDigits uint8

	MinExponentDigits uint8
***REMOVED***

// RoundSignificantDigits returns the number of significant digits an
// implementation of Convert may round to or n < 0 if there is no maximum or
// a maximum is not recommended.
func (r *RoundingContext) RoundSignificantDigits() (n int) ***REMOVED***
	if r.MaxFractionDigits == 0 && r.MaxSignificantDigits > 0 ***REMOVED***
		return int(r.MaxSignificantDigits)
	***REMOVED*** else if r.isScientific() && r.MaxIntegerDigits == 1 ***REMOVED***
		if r.MaxSignificantDigits == 0 ||
			int(r.MaxFractionDigits+1) == int(r.MaxSignificantDigits) ***REMOVED***
			// Note: don't add DigitShift: it is only used for decimals.
			return int(r.MaxFractionDigits) + 1
		***REMOVED***
	***REMOVED***
	return -1
***REMOVED***

// RoundFractionDigits returns the number of fraction digits an implementation
// of Convert may round to or n < 0 if there is no maximum or a maximum is not
// recommended.
func (r *RoundingContext) RoundFractionDigits() (n int) ***REMOVED***
	if r.MinExponentDigits == 0 &&
		r.MaxSignificantDigits == 0 &&
		r.MaxFractionDigits >= 0 ***REMOVED***
		return int(r.MaxFractionDigits) + int(r.DigitShift)
	***REMOVED***
	return -1
***REMOVED***

// SetScale fixes the RoundingContext to a fixed number of fraction digits.
func (r *RoundingContext) SetScale(scale int) ***REMOVED***
	r.MinFractionDigits = uint8(scale)
	r.MaxFractionDigits = int16(scale)
***REMOVED***

func (r *RoundingContext) SetPrecision(prec int) ***REMOVED***
	r.MaxSignificantDigits = int16(prec)
***REMOVED***

func (r *RoundingContext) isScientific() bool ***REMOVED***
	return r.MinExponentDigits > 0
***REMOVED***

func (f *Pattern) needsSep(pos int) bool ***REMOVED***
	p := pos - 1
	size := int(f.GroupingSize[0])
	if size == 0 || p == 0 ***REMOVED***
		return false
	***REMOVED***
	if p == size ***REMOVED***
		return true
	***REMOVED***
	if p -= size; p < 0 ***REMOVED***
		return false
	***REMOVED***
	// TODO: make second groupingsize the same as first if 0 so that we can
	// avoid this check.
	if x := int(f.GroupingSize[1]); x != 0 ***REMOVED***
		size = x
	***REMOVED***
	return p%size == 0
***REMOVED***

// A PatternFlag is a bit mask for the flag field of a Pattern.
type PatternFlag uint8

const (
	AlwaysSign PatternFlag = 1 << iota
	ElideSign              // Use space instead of plus sign. AlwaysSign must be true.
	AlwaysExpSign
	AlwaysDecimalSeparator
	ParenthesisForNegative // Common pattern. Saves space.

	PadAfterNumber
	PadAfterAffix

	PadBeforePrefix = 0 // Default
	PadAfterPrefix  = PadAfterAffix
	PadBeforeSuffix = PadAfterNumber
	PadAfterSuffix  = PadAfterNumber | PadAfterAffix
	PadMask         = PadAfterNumber | PadAfterAffix
)

type parser struct ***REMOVED***
	*Pattern

	leadingSharps int

	pos            int
	err            error
	doNotTerminate bool
	groupingCount  uint
	hasGroup       bool
	buf            []byte
***REMOVED***

func (p *parser) setError(err error) ***REMOVED***
	if p.err == nil ***REMOVED***
		p.err = err
	***REMOVED***
***REMOVED***

func (p *parser) updateGrouping() ***REMOVED***
	if p.hasGroup &&
		0 < p.groupingCount && p.groupingCount < 255 ***REMOVED***
		p.GroupingSize[1] = p.GroupingSize[0]
		p.GroupingSize[0] = uint8(p.groupingCount)
	***REMOVED***
	p.groupingCount = 0
	p.hasGroup = true
***REMOVED***

var (
	// TODO: more sensible and localizeable error messages.
	errMultiplePadSpecifiers = errors.New("format: pattern has multiple pad specifiers")
	errInvalidPadSpecifier   = errors.New("format: invalid pad specifier")
	errInvalidQuote          = errors.New("format: invalid quote")
	errAffixTooLarge         = errors.New("format: prefix or suffix exceeds maximum UTF-8 length of 256 bytes")
	errDuplicatePercentSign  = errors.New("format: duplicate percent sign")
	errDuplicatePermilleSign = errors.New("format: duplicate permille sign")
	errUnexpectedEnd         = errors.New("format: unexpected end of pattern")
)

// ParsePattern extracts formatting information from a CLDR number pattern.
//
// See http://unicode.org/reports/tr35/tr35-numbers.html#Number_Format_Patterns.
func ParsePattern(s string) (f *Pattern, err error) ***REMOVED***
	p := parser***REMOVED***Pattern: &Pattern***REMOVED******REMOVED******REMOVED***

	s = p.parseSubPattern(s)

	if s != "" ***REMOVED***
		// Parse negative sub pattern.
		if s[0] != ';' ***REMOVED***
			p.setError(errors.New("format: error parsing first sub pattern"))
			return nil, p.err
		***REMOVED***
		neg := parser***REMOVED***Pattern: &Pattern***REMOVED******REMOVED******REMOVED*** // just for extracting the affixes.
		s = neg.parseSubPattern(s[len(";"):])
		p.NegOffset = uint16(len(p.buf))
		p.buf = append(p.buf, neg.buf...)
	***REMOVED***
	if s != "" ***REMOVED***
		p.setError(errors.New("format: spurious characters at end of pattern"))
	***REMOVED***
	if p.err != nil ***REMOVED***
		return nil, p.err
	***REMOVED***
	if affix := string(p.buf); affix == "\x00\x00" || affix == "\x00\x00\x00\x00" ***REMOVED***
		// No prefix or suffixes.
		p.NegOffset = 0
	***REMOVED*** else ***REMOVED***
		p.Affix = affix
	***REMOVED***
	if p.Increment == 0 ***REMOVED***
		p.IncrementScale = 0
	***REMOVED***
	return p.Pattern, nil
***REMOVED***

func (p *parser) parseSubPattern(s string) string ***REMOVED***
	s = p.parsePad(s, PadBeforePrefix)
	s = p.parseAffix(s)
	s = p.parsePad(s, PadAfterPrefix)

	s = p.parse(p.number, s)
	p.updateGrouping()

	s = p.parsePad(s, PadBeforeSuffix)
	s = p.parseAffix(s)
	s = p.parsePad(s, PadAfterSuffix)
	return s
***REMOVED***

func (p *parser) parsePad(s string, f PatternFlag) (tail string) ***REMOVED***
	if len(s) >= 2 && s[0] == '*' ***REMOVED***
		r, sz := utf8.DecodeRuneInString(s[1:])
		if p.PadRune != 0 ***REMOVED***
			p.err = errMultiplePadSpecifiers
		***REMOVED*** else ***REMOVED***
			p.Flags |= f
			p.PadRune = r
		***REMOVED***
		return s[1+sz:]
	***REMOVED***
	return s
***REMOVED***

func (p *parser) parseAffix(s string) string ***REMOVED***
	x := len(p.buf)
	p.buf = append(p.buf, 0) // placeholder for affix length

	s = p.parse(p.affix, s)

	n := len(p.buf) - x - 1
	if n > 0xFF ***REMOVED***
		p.setError(errAffixTooLarge)
	***REMOVED***
	p.buf[x] = uint8(n)
	return s
***REMOVED***

// state implements a state transition. It returns the new state. A state
// function may set an error on the parser or may simply return on an incorrect
// token and let the next phase fail.
type state func(r rune) state

// parse repeatedly applies a state function on the given string until a
// termination condition is reached.
func (p *parser) parse(fn state, s string) (tail string) ***REMOVED***
	for i, r := range s ***REMOVED***
		p.doNotTerminate = false
		if fn = fn(r); fn == nil || p.err != nil ***REMOVED***
			return s[i:]
		***REMOVED***
		p.FormatWidth++
	***REMOVED***
	if p.doNotTerminate ***REMOVED***
		p.setError(errUnexpectedEnd)
	***REMOVED***
	return ""
***REMOVED***

func (p *parser) affix(r rune) state ***REMOVED***
	switch r ***REMOVED***
	case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9',
		'#', '@', '.', '*', ',', ';':
		return nil
	case '\'':
		p.FormatWidth--
		return p.escapeFirst
	case '%':
		if p.DigitShift != 0 ***REMOVED***
			p.setError(errDuplicatePercentSign)
		***REMOVED***
		p.DigitShift = 2
	case '\u2030': // ‰ Per mille
		if p.DigitShift != 0 ***REMOVED***
			p.setError(errDuplicatePermilleSign)
		***REMOVED***
		p.DigitShift = 3
		// TODO: handle currency somehow: ¤, ¤¤, ¤¤¤, ¤¤¤¤
	***REMOVED***
	p.buf = append(p.buf, string(r)...)
	return p.affix
***REMOVED***

func (p *parser) escapeFirst(r rune) state ***REMOVED***
	switch r ***REMOVED***
	case '\'':
		p.buf = append(p.buf, "\\'"...)
		return p.affix
	default:
		p.buf = append(p.buf, '\'')
		p.buf = append(p.buf, string(r)...)
	***REMOVED***
	return p.escape
***REMOVED***

func (p *parser) escape(r rune) state ***REMOVED***
	switch r ***REMOVED***
	case '\'':
		p.FormatWidth--
		p.buf = append(p.buf, '\'')
		return p.affix
	default:
		p.buf = append(p.buf, string(r)...)
	***REMOVED***
	return p.escape
***REMOVED***

// number parses a number. The BNF says the integer part should always have
// a '0', but that does not appear to be the case according to the rest of the
// documentation. We will allow having only '#' numbers.
func (p *parser) number(r rune) state ***REMOVED***
	switch r ***REMOVED***
	case '#':
		p.groupingCount++
		p.leadingSharps++
	case '@':
		p.groupingCount++
		p.leadingSharps = 0
		p.MaxFractionDigits = -1
		return p.sigDigits(r)
	case ',':
		if p.leadingSharps == 0 ***REMOVED*** // no leading commas
			return nil
		***REMOVED***
		p.updateGrouping()
	case 'E':
		p.MaxIntegerDigits = uint8(p.leadingSharps)
		return p.exponent
	case '.': // allow ".##" etc.
		p.updateGrouping()
		return p.fraction
	case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
		return p.integer(r)
	default:
		return nil
	***REMOVED***
	return p.number
***REMOVED***

func (p *parser) integer(r rune) state ***REMOVED***
	if !('0' <= r && r <= '9') ***REMOVED***
		var next state
		switch r ***REMOVED***
		case 'E':
			if p.leadingSharps > 0 ***REMOVED***
				p.MaxIntegerDigits = uint8(p.leadingSharps) + p.MinIntegerDigits
			***REMOVED***
			next = p.exponent
		case '.':
			next = p.fraction
		case ',':
			next = p.integer
		***REMOVED***
		p.updateGrouping()
		return next
	***REMOVED***
	p.Increment = p.Increment*10 + uint32(r-'0')
	p.groupingCount++
	p.MinIntegerDigits++
	return p.integer
***REMOVED***

func (p *parser) sigDigits(r rune) state ***REMOVED***
	switch r ***REMOVED***
	case '@':
		p.groupingCount++
		p.MaxSignificantDigits++
		p.MinSignificantDigits++
	case '#':
		return p.sigDigitsFinal(r)
	case 'E':
		p.updateGrouping()
		return p.normalizeSigDigitsWithExponent()
	default:
		p.updateGrouping()
		return nil
	***REMOVED***
	return p.sigDigits
***REMOVED***

func (p *parser) sigDigitsFinal(r rune) state ***REMOVED***
	switch r ***REMOVED***
	case '#':
		p.groupingCount++
		p.MaxSignificantDigits++
	case 'E':
		p.updateGrouping()
		return p.normalizeSigDigitsWithExponent()
	default:
		p.updateGrouping()
		return nil
	***REMOVED***
	return p.sigDigitsFinal
***REMOVED***

func (p *parser) normalizeSigDigitsWithExponent() state ***REMOVED***
	p.MinIntegerDigits, p.MaxIntegerDigits = 1, 1
	p.MinFractionDigits = p.MinSignificantDigits - 1
	p.MaxFractionDigits = p.MaxSignificantDigits - 1
	p.MinSignificantDigits, p.MaxSignificantDigits = 0, 0
	return p.exponent
***REMOVED***

func (p *parser) fraction(r rune) state ***REMOVED***
	switch r ***REMOVED***
	case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
		p.Increment = p.Increment*10 + uint32(r-'0')
		p.IncrementScale++
		p.MinFractionDigits++
		p.MaxFractionDigits++
	case '#':
		p.MaxFractionDigits++
	case 'E':
		if p.leadingSharps > 0 ***REMOVED***
			p.MaxIntegerDigits = uint8(p.leadingSharps) + p.MinIntegerDigits
		***REMOVED***
		return p.exponent
	default:
		return nil
	***REMOVED***
	return p.fraction
***REMOVED***

func (p *parser) exponent(r rune) state ***REMOVED***
	switch r ***REMOVED***
	case '+':
		// Set mode and check it wasn't already set.
		if p.Flags&AlwaysExpSign != 0 || p.MinExponentDigits > 0 ***REMOVED***
			break
		***REMOVED***
		p.Flags |= AlwaysExpSign
		p.doNotTerminate = true
		return p.exponent
	case '0':
		p.MinExponentDigits++
		return p.exponent
	***REMOVED***
	// termination condition
	if p.MinExponentDigits == 0 ***REMOVED***
		p.setError(errors.New("format: need at least one digit"))
	***REMOVED***
	return nil
***REMOVED***
