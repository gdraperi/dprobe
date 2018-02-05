// Copyright 2017 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package number

import (
	"strconv"
	"unicode/utf8"

	"golang.org/x/text/language"
)

// TODO:
// - grouping of fractions
// - allow user-defined superscript notation (such as <sup>4</sup>)
// - same for non-breaking spaces, like &nbsp;

// A VisibleDigits computes digits, comma placement and trailing zeros as they
// will be shown to the user.
type VisibleDigits interface ***REMOVED***
	Digits(buf []byte, t language.Tag, scale int) Digits
	// TODO: Do we also need to add the verb or pass a format.State?
***REMOVED***

// Formatting proceeds along the following lines:
// 0) Compose rounding information from format and context.
// 1) Convert a number into a Decimal.
// 2) Sanitize Decimal by adding trailing zeros, removing leading digits, and
//    (non-increment) rounding. The Decimal that results from this is suitable
//    for determining the plural form.
// 3) Render the Decimal in the localized form.

// Formatter contains all the information needed to render a number.
type Formatter struct ***REMOVED***
	Pattern
	Info
***REMOVED***

func (f *Formatter) init(t language.Tag, index []uint8) ***REMOVED***
	f.Info = InfoFromTag(t)
	for ; ; t = t.Parent() ***REMOVED***
		if ci, ok := language.CompactIndex(t); ok ***REMOVED***
			f.Pattern = formats[index[ci]]
			break
		***REMOVED***
	***REMOVED***
***REMOVED***

// InitPattern initializes a Formatter for the given Pattern.
func (f *Formatter) InitPattern(t language.Tag, pat *Pattern) ***REMOVED***
	f.Info = InfoFromTag(t)
	f.Pattern = *pat
***REMOVED***

// InitDecimal initializes a Formatter using the default Pattern for the given
// language.
func (f *Formatter) InitDecimal(t language.Tag) ***REMOVED***
	f.init(t, tagToDecimal)
***REMOVED***

// InitScientific initializes a Formatter using the default Pattern for the
// given language.
func (f *Formatter) InitScientific(t language.Tag) ***REMOVED***
	f.init(t, tagToScientific)
	f.Pattern.MinFractionDigits = 0
	f.Pattern.MaxFractionDigits = -1
***REMOVED***

// InitEngineering initializes a Formatter using the default Pattern for the
// given language.
func (f *Formatter) InitEngineering(t language.Tag) ***REMOVED***
	f.init(t, tagToScientific)
	f.Pattern.MinFractionDigits = 0
	f.Pattern.MaxFractionDigits = -1
	f.Pattern.MaxIntegerDigits = 3
	f.Pattern.MinIntegerDigits = 1
***REMOVED***

// InitPercent initializes a Formatter using the default Pattern for the given
// language.
func (f *Formatter) InitPercent(t language.Tag) ***REMOVED***
	f.init(t, tagToPercent)
***REMOVED***

// InitPerMille initializes a Formatter using the default Pattern for the given
// language.
func (f *Formatter) InitPerMille(t language.Tag) ***REMOVED***
	f.init(t, tagToPercent)
	f.Pattern.DigitShift = 3
***REMOVED***

func (f *Formatter) Append(dst []byte, x interface***REMOVED******REMOVED***) []byte ***REMOVED***
	var d Decimal
	r := f.RoundingContext
	d.Convert(r, x)
	return f.Render(dst, FormatDigits(&d, r))
***REMOVED***

func FormatDigits(d *Decimal, r RoundingContext) Digits ***REMOVED***
	if r.isScientific() ***REMOVED***
		return scientificVisibleDigits(r, d)
	***REMOVED***
	return decimalVisibleDigits(r, d)
***REMOVED***

func (f *Formatter) Format(dst []byte, d *Decimal) []byte ***REMOVED***
	return f.Render(dst, FormatDigits(d, f.RoundingContext))
***REMOVED***

func (f *Formatter) Render(dst []byte, d Digits) []byte ***REMOVED***
	var result []byte
	var postPrefix, preSuffix int
	if d.IsScientific ***REMOVED***
		result, postPrefix, preSuffix = appendScientific(dst, f, &d)
	***REMOVED*** else ***REMOVED***
		result, postPrefix, preSuffix = appendDecimal(dst, f, &d)
	***REMOVED***
	if f.PadRune == 0 ***REMOVED***
		return result
	***REMOVED***
	width := int(f.FormatWidth)
	if count := utf8.RuneCount(result); count < width ***REMOVED***
		insertPos := 0
		switch f.Flags & PadMask ***REMOVED***
		case PadAfterPrefix:
			insertPos = postPrefix
		case PadBeforeSuffix:
			insertPos = preSuffix
		case PadAfterSuffix:
			insertPos = len(result)
		***REMOVED***
		num := width - count
		pad := [utf8.UTFMax]byte***REMOVED***' '***REMOVED***
		sz := 1
		if r := f.PadRune; r != 0 ***REMOVED***
			sz = utf8.EncodeRune(pad[:], r)
		***REMOVED***
		extra := sz * num
		if n := len(result) + extra; n < cap(result) ***REMOVED***
			result = result[:n]
			copy(result[insertPos+extra:], result[insertPos:])
		***REMOVED*** else ***REMOVED***
			buf := make([]byte, n)
			copy(buf, result[:insertPos])
			copy(buf[insertPos+extra:], result[insertPos:])
			result = buf
		***REMOVED***
		for ; num > 0; num-- ***REMOVED***
			insertPos += copy(result[insertPos:], pad[:sz])
		***REMOVED***
	***REMOVED***
	return result
***REMOVED***

// decimalVisibleDigits converts d according to the RoundingContext. Note that
// the exponent may change as a result of this operation.
func decimalVisibleDigits(r RoundingContext, d *Decimal) Digits ***REMOVED***
	if d.NaN || d.Inf ***REMOVED***
		return Digits***REMOVED***digits: digits***REMOVED***Neg: d.Neg, NaN: d.NaN, Inf: d.Inf***REMOVED******REMOVED***
	***REMOVED***
	n := Digits***REMOVED***digits: d.normalize().digits***REMOVED***

	exp := n.Exp
	exp += int32(r.DigitShift)

	// Cap integer digits. Remove *most-significant* digits.
	if r.MaxIntegerDigits > 0 ***REMOVED***
		if p := int(exp) - int(r.MaxIntegerDigits); p > 0 ***REMOVED***
			if p > len(n.Digits) ***REMOVED***
				p = len(n.Digits)
			***REMOVED***
			if n.Digits = n.Digits[p:]; len(n.Digits) == 0 ***REMOVED***
				exp = 0
			***REMOVED*** else ***REMOVED***
				exp -= int32(p)
			***REMOVED***
			// Strip leading zeros.
			for len(n.Digits) > 0 && n.Digits[0] == 0 ***REMOVED***
				n.Digits = n.Digits[1:]
				exp--
			***REMOVED***
		***REMOVED***
	***REMOVED***

	// Rounding if not already done by Convert.
	p := len(n.Digits)
	if maxSig := int(r.MaxSignificantDigits); maxSig > 0 ***REMOVED***
		p = maxSig
	***REMOVED***
	if maxFrac := int(r.MaxFractionDigits); maxFrac >= 0 ***REMOVED***
		if cap := int(exp) + maxFrac; cap < p ***REMOVED***
			p = int(exp) + maxFrac
		***REMOVED***
		if p < 0 ***REMOVED***
			p = 0
		***REMOVED***
	***REMOVED***
	n.round(r.Mode, p)

	// set End (trailing zeros)
	n.End = int32(len(n.Digits))
	if n.End == 0 ***REMOVED***
		exp = 0
		if r.MinFractionDigits > 0 ***REMOVED***
			n.End = int32(r.MinFractionDigits)
		***REMOVED***
		if p := int32(r.MinSignificantDigits) - 1; p > n.End ***REMOVED***
			n.End = p
		***REMOVED***
	***REMOVED*** else ***REMOVED***
		if end := exp + int32(r.MinFractionDigits); end > n.End ***REMOVED***
			n.End = end
		***REMOVED***
		if n.End < int32(r.MinSignificantDigits) ***REMOVED***
			n.End = int32(r.MinSignificantDigits)
		***REMOVED***
	***REMOVED***
	n.Exp = exp
	return n
***REMOVED***

// appendDecimal appends a formatted number to dst. It returns two possible
// insertion points for padding.
func appendDecimal(dst []byte, f *Formatter, n *Digits) (b []byte, postPre, preSuf int) ***REMOVED***
	if dst, ok := f.renderSpecial(dst, n); ok ***REMOVED***
		return dst, 0, len(dst)
	***REMOVED***
	digits := n.Digits
	exp := n.Exp

	// Split in integer and fraction part.
	var intDigits, fracDigits []byte
	numInt := 0
	numFrac := int(n.End - n.Exp)
	if exp > 0 ***REMOVED***
		numInt = int(exp)
		if int(exp) >= len(digits) ***REMOVED*** // ddddd | ddddd00
			intDigits = digits
		***REMOVED*** else ***REMOVED*** // ddd.dd
			intDigits = digits[:exp]
			fracDigits = digits[exp:]
		***REMOVED***
	***REMOVED*** else ***REMOVED***
		fracDigits = digits
	***REMOVED***

	neg := n.Neg
	affix, suffix := f.getAffixes(neg)
	dst = appendAffix(dst, f, affix, neg)
	savedLen := len(dst)

	minInt := int(f.MinIntegerDigits)
	if minInt == 0 && f.MinSignificantDigits > 0 ***REMOVED***
		minInt = 1
	***REMOVED***
	// add leading zeros
	for i := minInt; i > numInt; i-- ***REMOVED***
		dst = f.AppendDigit(dst, 0)
		if f.needsSep(i) ***REMOVED***
			dst = append(dst, f.Symbol(SymGroup)...)
		***REMOVED***
	***REMOVED***
	i := 0
	for ; i < len(intDigits); i++ ***REMOVED***
		dst = f.AppendDigit(dst, intDigits[i])
		if f.needsSep(numInt - i) ***REMOVED***
			dst = append(dst, f.Symbol(SymGroup)...)
		***REMOVED***
	***REMOVED***
	for ; i < numInt; i++ ***REMOVED***
		dst = f.AppendDigit(dst, 0)
		if f.needsSep(numInt - i) ***REMOVED***
			dst = append(dst, f.Symbol(SymGroup)...)
		***REMOVED***
	***REMOVED***

	if numFrac > 0 || f.Flags&AlwaysDecimalSeparator != 0 ***REMOVED***
		dst = append(dst, f.Symbol(SymDecimal)...)
	***REMOVED***
	// Add trailing zeros
	i = 0
	for n := -int(n.Exp); i < n; i++ ***REMOVED***
		dst = f.AppendDigit(dst, 0)
	***REMOVED***
	for _, d := range fracDigits ***REMOVED***
		i++
		dst = f.AppendDigit(dst, d)
	***REMOVED***
	for ; i < numFrac; i++ ***REMOVED***
		dst = f.AppendDigit(dst, 0)
	***REMOVED***
	return appendAffix(dst, f, suffix, neg), savedLen, len(dst)
***REMOVED***

func scientificVisibleDigits(r RoundingContext, d *Decimal) Digits ***REMOVED***
	if d.NaN || d.Inf ***REMOVED***
		return Digits***REMOVED***digits: digits***REMOVED***Neg: d.Neg, NaN: d.NaN, Inf: d.Inf***REMOVED******REMOVED***
	***REMOVED***
	n := Digits***REMOVED***digits: d.normalize().digits, IsScientific: true***REMOVED***

	// Normalize to have at least one digit. This simplifies engineering
	// notation.
	if len(n.Digits) == 0 ***REMOVED***
		n.Digits = append(n.Digits, 0)
		n.Exp = 1
	***REMOVED***

	// Significant digits are transformed by the parser for scientific notation
	// and do not need to be handled here.
	maxInt, numInt := int(r.MaxIntegerDigits), int(r.MinIntegerDigits)
	if numInt == 0 ***REMOVED***
		numInt = 1
	***REMOVED***

	// If a maximum number of integers is specified, the minimum must be 1
	// and the exponent is grouped by this number (e.g. for engineering)
	if maxInt > numInt ***REMOVED***
		// Correct the exponent to reflect a single integer digit.
		numInt = 1
		// engineering
		// 0.01234 ([12345]e-1) -> 1.2345e-2  12.345e-3
		// 12345   ([12345]e+5) -> 1.2345e4  12.345e3
		d := int(n.Exp-1) % maxInt
		if d < 0 ***REMOVED***
			d += maxInt
		***REMOVED***
		numInt += d
	***REMOVED***

	p := len(n.Digits)
	if maxSig := int(r.MaxSignificantDigits); maxSig > 0 ***REMOVED***
		p = maxSig
	***REMOVED***
	if maxFrac := int(r.MaxFractionDigits); maxFrac >= 0 && numInt+maxFrac < p ***REMOVED***
		p = numInt + maxFrac
	***REMOVED***
	n.round(r.Mode, p)

	n.Comma = uint8(numInt)
	n.End = int32(len(n.Digits))
	if minSig := int32(r.MinFractionDigits) + int32(numInt); n.End < minSig ***REMOVED***
		n.End = minSig
	***REMOVED***
	return n
***REMOVED***

// appendScientific appends a formatted number to dst. It returns two possible
// insertion points for padding.
func appendScientific(dst []byte, f *Formatter, n *Digits) (b []byte, postPre, preSuf int) ***REMOVED***
	if dst, ok := f.renderSpecial(dst, n); ok ***REMOVED***
		return dst, 0, 0
	***REMOVED***
	digits := n.Digits
	numInt := int(n.Comma)
	numFrac := int(n.End) - int(n.Comma)

	var intDigits, fracDigits []byte
	if numInt <= len(digits) ***REMOVED***
		intDigits = digits[:numInt]
		fracDigits = digits[numInt:]
	***REMOVED*** else ***REMOVED***
		intDigits = digits
	***REMOVED***
	neg := n.Neg
	affix, suffix := f.getAffixes(neg)
	dst = appendAffix(dst, f, affix, neg)
	savedLen := len(dst)

	i := 0
	for ; i < len(intDigits); i++ ***REMOVED***
		dst = f.AppendDigit(dst, intDigits[i])
		if f.needsSep(numInt - i) ***REMOVED***
			dst = append(dst, f.Symbol(SymGroup)...)
		***REMOVED***
	***REMOVED***
	for ; i < numInt; i++ ***REMOVED***
		dst = f.AppendDigit(dst, 0)
		if f.needsSep(numInt - i) ***REMOVED***
			dst = append(dst, f.Symbol(SymGroup)...)
		***REMOVED***
	***REMOVED***

	if numFrac > 0 || f.Flags&AlwaysDecimalSeparator != 0 ***REMOVED***
		dst = append(dst, f.Symbol(SymDecimal)...)
	***REMOVED***
	i = 0
	for ; i < len(fracDigits); i++ ***REMOVED***
		dst = f.AppendDigit(dst, fracDigits[i])
	***REMOVED***
	for ; i < numFrac; i++ ***REMOVED***
		dst = f.AppendDigit(dst, 0)
	***REMOVED***

	// exp
	buf := [12]byte***REMOVED******REMOVED***
	// TODO: use exponential if superscripting is not available (no Latin
	// numbers or no tags) and use exponential in all other cases.
	exp := n.Exp - int32(n.Comma)
	exponential := f.Symbol(SymExponential)
	if exponential == "E" ***REMOVED***
		dst = append(dst, "\u202f"...) // NARROW NO-BREAK SPACE
		dst = append(dst, f.Symbol(SymSuperscriptingExponent)...)
		dst = append(dst, "\u202f"...) // NARROW NO-BREAK SPACE
		dst = f.AppendDigit(dst, 1)
		dst = f.AppendDigit(dst, 0)
		switch ***REMOVED***
		case exp < 0:
			dst = append(dst, superMinus...)
			exp = -exp
		case f.Flags&AlwaysExpSign != 0:
			dst = append(dst, superPlus...)
		***REMOVED***
		b = strconv.AppendUint(buf[:0], uint64(exp), 10)
		for i := len(b); i < int(f.MinExponentDigits); i++ ***REMOVED***
			dst = append(dst, superDigits[0]...)
		***REMOVED***
		for _, c := range b ***REMOVED***
			dst = append(dst, superDigits[c-'0']...)
		***REMOVED***
	***REMOVED*** else ***REMOVED***
		dst = append(dst, exponential...)
		switch ***REMOVED***
		case exp < 0:
			dst = append(dst, f.Symbol(SymMinusSign)...)
			exp = -exp
		case f.Flags&AlwaysExpSign != 0:
			dst = append(dst, f.Symbol(SymPlusSign)...)
		***REMOVED***
		b = strconv.AppendUint(buf[:0], uint64(exp), 10)
		for i := len(b); i < int(f.MinExponentDigits); i++ ***REMOVED***
			dst = f.AppendDigit(dst, 0)
		***REMOVED***
		for _, c := range b ***REMOVED***
			dst = f.AppendDigit(dst, c-'0')
		***REMOVED***
	***REMOVED***
	return appendAffix(dst, f, suffix, neg), savedLen, len(dst)
***REMOVED***

const (
	superMinus = "\u207B" // SUPERSCRIPT HYPHEN-MINUS
	superPlus  = "\u207A" // SUPERSCRIPT PLUS SIGN
)

var (
	// Note: the digits are not sequential!!!
	superDigits = []string***REMOVED***
		"\u2070", // SUPERSCRIPT DIGIT ZERO
		"\u00B9", // SUPERSCRIPT DIGIT ONE
		"\u00B2", // SUPERSCRIPT DIGIT TWO
		"\u00B3", // SUPERSCRIPT DIGIT THREE
		"\u2074", // SUPERSCRIPT DIGIT FOUR
		"\u2075", // SUPERSCRIPT DIGIT FIVE
		"\u2076", // SUPERSCRIPT DIGIT SIX
		"\u2077", // SUPERSCRIPT DIGIT SEVEN
		"\u2078", // SUPERSCRIPT DIGIT EIGHT
		"\u2079", // SUPERSCRIPT DIGIT NINE
	***REMOVED***
)

func (f *Formatter) getAffixes(neg bool) (affix, suffix string) ***REMOVED***
	str := f.Affix
	if str != "" ***REMOVED***
		if f.NegOffset > 0 ***REMOVED***
			if neg ***REMOVED***
				str = str[f.NegOffset:]
			***REMOVED*** else ***REMOVED***
				str = str[:f.NegOffset]
			***REMOVED***
		***REMOVED***
		sufStart := 1 + str[0]
		affix = str[1:sufStart]
		suffix = str[sufStart+1:]
	***REMOVED***
	// TODO: introduce a NeedNeg sign to indicate if the left pattern already
	// has a sign marked?
	if f.NegOffset == 0 && (neg || f.Flags&AlwaysSign != 0) ***REMOVED***
		affix = "-" + affix
	***REMOVED***
	return affix, suffix
***REMOVED***

func (f *Formatter) renderSpecial(dst []byte, d *Digits) (b []byte, ok bool) ***REMOVED***
	if d.NaN ***REMOVED***
		return fmtNaN(dst, f), true
	***REMOVED***
	if d.Inf ***REMOVED***
		return fmtInfinite(dst, f, d), true
	***REMOVED***
	return dst, false
***REMOVED***

func fmtNaN(dst []byte, f *Formatter) []byte ***REMOVED***
	return append(dst, f.Symbol(SymNan)...)
***REMOVED***

func fmtInfinite(dst []byte, f *Formatter, d *Digits) []byte ***REMOVED***
	affix, suffix := f.getAffixes(d.Neg)
	dst = appendAffix(dst, f, affix, d.Neg)
	dst = append(dst, f.Symbol(SymInfinity)...)
	dst = appendAffix(dst, f, suffix, d.Neg)
	return dst
***REMOVED***

func appendAffix(dst []byte, f *Formatter, affix string, neg bool) []byte ***REMOVED***
	quoting := false
	escaping := false
	for _, r := range affix ***REMOVED***
		switch ***REMOVED***
		case escaping:
			// escaping occurs both inside and outside of quotes
			dst = append(dst, string(r)...)
			escaping = false
		case r == '\\':
			escaping = true
		case r == '\'':
			quoting = !quoting
		case quoting:
			dst = append(dst, string(r)...)
		case r == '%':
			if f.DigitShift == 3 ***REMOVED***
				dst = append(dst, f.Symbol(SymPerMille)...)
			***REMOVED*** else ***REMOVED***
				dst = append(dst, f.Symbol(SymPercentSign)...)
			***REMOVED***
		case r == '-' || r == '+':
			if neg ***REMOVED***
				dst = append(dst, f.Symbol(SymMinusSign)...)
			***REMOVED*** else if f.Flags&ElideSign == 0 ***REMOVED***
				dst = append(dst, f.Symbol(SymPlusSign)...)
			***REMOVED*** else ***REMOVED***
				dst = append(dst, ' ')
			***REMOVED***
		default:
			dst = append(dst, string(r)...)
		***REMOVED***
	***REMOVED***
	return dst
***REMOVED***
