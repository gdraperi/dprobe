// Copyright 2017 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:generate stringer -type RoundingMode

package number

import (
	"math"
	"strconv"
)

// RoundingMode determines how a number is rounded to the desired precision.
type RoundingMode byte

const (
	ToNearestEven RoundingMode = iota // towards the nearest integer, or towards an even number if equidistant.
	ToNearestZero                     // towards the nearest integer, or towards zero if equidistant.
	ToNearestAway                     // towards the nearest integer, or away from zero if equidistant.
	ToPositiveInf                     // towards infinity
	ToNegativeInf                     // towards negative infinity
	ToZero                            // towards zero
	AwayFromZero                      // away from zero
	numModes
)

const maxIntDigits = 20

// A Decimal represents a floating point number in decimal format.
// Digits represents a number [0, 1.0), and the absolute value represented by
// Decimal is Digits * 10^Exp. Leading and trailing zeros may be omitted and Exp
// may point outside a valid position in Digits.
//
// Examples:
//      Number     Decimal
//      12345      Digits: [1, 2, 3, 4, 5], Exp: 5
//      12.345     Digits: [1, 2, 3, 4, 5], Exp: 2
//      12000      Digits: [1, 2],          Exp: 5
//      12000.00   Digits: [1, 2],          Exp: 5
//      0.00123    Digits: [1, 2, 3],       Exp: -2
//      0          Digits: [],              Exp: 0
type Decimal struct ***REMOVED***
	digits

	buf [maxIntDigits]byte
***REMOVED***

type digits struct ***REMOVED***
	Digits []byte // mantissa digits, big-endian
	Exp    int32  // exponent
	Neg    bool
	Inf    bool // Takes precedence over Digits and Exp.
	NaN    bool // Takes precedence over Inf.
***REMOVED***

// Digits represents a floating point number represented in digits of the
// base in which a number is to be displayed. It is similar to Decimal, but
// keeps track of trailing fraction zeros and the comma placement for
// engineering notation. Digits must have at least one digit.
//
// Examples:
//      Number     Decimal
//    decimal
//      12345      Digits: [1, 2, 3, 4, 5], Exp: 5  End: 5
//      12.345     Digits: [1, 2, 3, 4, 5], Exp: 2  End: 5
//      12000      Digits: [1, 2],          Exp: 5  End: 5
//      12000.00   Digits: [1, 2],          Exp: 5  End: 7
//      0.00123    Digits: [1, 2, 3],       Exp: -2 End: 3
//      0          Digits: [],              Exp: 0  End: 1
//    scientific (actual exp is Exp - Comma)
//      0e0        Digits: [0],             Exp: 1, End: 1, Comma: 1
//      .0e0       Digits: [0],             Exp: 0, End: 1, Comma: 0
//      0.0e0      Digits: [0],             Exp: 1, End: 2, Comma: 1
//      1.23e4     Digits: [1, 2, 3],       Exp: 5, End: 3, Comma: 1
//      .123e5     Digits: [1, 2, 3],       Exp: 5, End: 3, Comma: 0
//    engineering
//      12.3e3     Digits: [1, 2, 3],       Exp: 5, End: 3, Comma: 2
type Digits struct ***REMOVED***
	digits
	// End indicates the end position of the number.
	End int32 // For decimals Exp <= End. For scientific len(Digits) <= End.
	// Comma is used for the comma position for scientific (always 0 or 1) and
	// engineering notation (always 0, 1, 2, or 3).
	Comma uint8
	// IsScientific indicates whether this number is to be rendered as a
	// scientific number.
	IsScientific bool
***REMOVED***

func (d *Digits) NumFracDigits() int ***REMOVED***
	if d.Exp >= d.End ***REMOVED***
		return 0
	***REMOVED***
	return int(d.End - d.Exp)
***REMOVED***

// normalize returns a new Decimal with leading and trailing zeros removed.
func (d *Decimal) normalize() (n Decimal) ***REMOVED***
	n = *d
	b := n.Digits
	// Strip leading zeros. Resulting number of digits is significant digits.
	for len(b) > 0 && b[0] == 0 ***REMOVED***
		b = b[1:]
		n.Exp--
	***REMOVED***
	// Strip trailing zeros
	for len(b) > 0 && b[len(b)-1] == 0 ***REMOVED***
		b = b[:len(b)-1]
	***REMOVED***
	if len(b) == 0 ***REMOVED***
		n.Exp = 0
	***REMOVED***
	n.Digits = b
	return n
***REMOVED***

func (d *Decimal) clear() ***REMOVED***
	b := d.Digits
	if b == nil ***REMOVED***
		b = d.buf[:0]
	***REMOVED***
	*d = Decimal***REMOVED******REMOVED***
	d.Digits = b[:0]
***REMOVED***

func (x *Decimal) String() string ***REMOVED***
	if x.NaN ***REMOVED***
		return "NaN"
	***REMOVED***
	var buf []byte
	if x.Neg ***REMOVED***
		buf = append(buf, '-')
	***REMOVED***
	if x.Inf ***REMOVED***
		buf = append(buf, "Inf"...)
		return string(buf)
	***REMOVED***
	switch ***REMOVED***
	case len(x.Digits) == 0:
		buf = append(buf, '0')
	case x.Exp <= 0:
		// 0.00ddd
		buf = append(buf, "0."...)
		buf = appendZeros(buf, -int(x.Exp))
		buf = appendDigits(buf, x.Digits)

	case /* 0 < */ int(x.Exp) < len(x.Digits):
		// dd.ddd
		buf = appendDigits(buf, x.Digits[:x.Exp])
		buf = append(buf, '.')
		buf = appendDigits(buf, x.Digits[x.Exp:])

	default: // len(x.Digits) <= x.Exp
		// ddd00
		buf = appendDigits(buf, x.Digits)
		buf = appendZeros(buf, int(x.Exp)-len(x.Digits))
	***REMOVED***
	return string(buf)
***REMOVED***

func appendDigits(buf []byte, digits []byte) []byte ***REMOVED***
	for _, c := range digits ***REMOVED***
		buf = append(buf, c+'0')
	***REMOVED***
	return buf
***REMOVED***

// appendZeros appends n 0 digits to buf and returns buf.
func appendZeros(buf []byte, n int) []byte ***REMOVED***
	for ; n > 0; n-- ***REMOVED***
		buf = append(buf, '0')
	***REMOVED***
	return buf
***REMOVED***

func (d *digits) round(mode RoundingMode, n int) ***REMOVED***
	if n >= len(d.Digits) ***REMOVED***
		return
	***REMOVED***
	// Make rounding decision: The result mantissa is truncated ("rounded down")
	// by default. Decide if we need to increment, or "round up", the (unsigned)
	// mantissa.
	inc := false
	switch mode ***REMOVED***
	case ToNegativeInf:
		inc = d.Neg
	case ToPositiveInf:
		inc = !d.Neg
	case ToZero:
		// nothing to do
	case AwayFromZero:
		inc = true
	case ToNearestEven:
		inc = d.Digits[n] > 5 || d.Digits[n] == 5 &&
			(len(d.Digits) > n+1 || n == 0 || d.Digits[n-1]&1 != 0)
	case ToNearestAway:
		inc = d.Digits[n] >= 5
	case ToNearestZero:
		inc = d.Digits[n] > 5 || d.Digits[n] == 5 && len(d.Digits) > n+1
	default:
		panic("unreachable")
	***REMOVED***
	if inc ***REMOVED***
		d.roundUp(n)
	***REMOVED*** else ***REMOVED***
		d.roundDown(n)
	***REMOVED***
***REMOVED***

// roundFloat rounds a floating point number.
func (r RoundingMode) roundFloat(x float64) float64 ***REMOVED***
	// Make rounding decision: The result mantissa is truncated ("rounded down")
	// by default. Decide if we need to increment, or "round up", the (unsigned)
	// mantissa.
	abs := x
	if x < 0 ***REMOVED***
		abs = -x
	***REMOVED***
	i, f := math.Modf(abs)
	if f == 0.0 ***REMOVED***
		return x
	***REMOVED***
	inc := false
	switch r ***REMOVED***
	case ToNegativeInf:
		inc = x < 0
	case ToPositiveInf:
		inc = x >= 0
	case ToZero:
		// nothing to do
	case AwayFromZero:
		inc = true
	case ToNearestEven:
		// TODO: check overflow
		inc = f > 0.5 || f == 0.5 && int64(i)&1 != 0
	case ToNearestAway:
		inc = f >= 0.5
	case ToNearestZero:
		inc = f > 0.5
	default:
		panic("unreachable")
	***REMOVED***
	if inc ***REMOVED***
		i += 1
	***REMOVED***
	if abs != x ***REMOVED***
		i = -i
	***REMOVED***
	return i
***REMOVED***

func (x *digits) roundUp(n int) ***REMOVED***
	if n < 0 || n >= len(x.Digits) ***REMOVED***
		return // nothing to do
	***REMOVED***
	// find first digit < 9
	for n > 0 && x.Digits[n-1] >= 9 ***REMOVED***
		n--
	***REMOVED***

	if n == 0 ***REMOVED***
		// all digits are 9s => round up to 1 and update exponent
		x.Digits[0] = 1 // ok since len(x.Digits) > n
		x.Digits = x.Digits[:1]
		x.Exp++
		return
	***REMOVED***
	x.Digits[n-1]++
	x.Digits = x.Digits[:n]
	// x already trimmed
***REMOVED***

func (x *digits) roundDown(n int) ***REMOVED***
	if n < 0 || n >= len(x.Digits) ***REMOVED***
		return // nothing to do
	***REMOVED***
	x.Digits = x.Digits[:n]
	trim(x)
***REMOVED***

// trim cuts off any trailing zeros from x's mantissa;
// they are meaningless for the value of x.
func trim(x *digits) ***REMOVED***
	i := len(x.Digits)
	for i > 0 && x.Digits[i-1] == 0 ***REMOVED***
		i--
	***REMOVED***
	x.Digits = x.Digits[:i]
	if i == 0 ***REMOVED***
		x.Exp = 0
	***REMOVED***
***REMOVED***

// A Converter converts a number into decimals according to the given rounding
// criteria.
type Converter interface ***REMOVED***
	Convert(d *Decimal, r RoundingContext)
***REMOVED***

const (
	signed   = true
	unsigned = false
)

// Convert converts the given number to the decimal representation using the
// supplied RoundingContext.
func (d *Decimal) Convert(r RoundingContext, number interface***REMOVED******REMOVED***) ***REMOVED***
	switch f := number.(type) ***REMOVED***
	case Converter:
		d.clear()
		f.Convert(d, r)
	case float32:
		d.ConvertFloat(r, float64(f), 32)
	case float64:
		d.ConvertFloat(r, f, 64)
	case int:
		d.ConvertInt(r, signed, uint64(f))
	case int8:
		d.ConvertInt(r, signed, uint64(f))
	case int16:
		d.ConvertInt(r, signed, uint64(f))
	case int32:
		d.ConvertInt(r, signed, uint64(f))
	case int64:
		d.ConvertInt(r, signed, uint64(f))
	case uint:
		d.ConvertInt(r, unsigned, uint64(f))
	case uint8:
		d.ConvertInt(r, unsigned, uint64(f))
	case uint16:
		d.ConvertInt(r, unsigned, uint64(f))
	case uint32:
		d.ConvertInt(r, unsigned, uint64(f))
	case uint64:
		d.ConvertInt(r, unsigned, f)

	default:
		d.NaN = true
		// TODO:
		// case string: if produced by strconv, allows for easy arbitrary pos.
		// case reflect.Value:
		// case big.Float
		// case big.Int
		// case big.Rat?
		// catch underlyings using reflect or will this already be done by the
		//    message package?
	***REMOVED***
***REMOVED***

// ConvertInt converts an integer to decimals.
func (d *Decimal) ConvertInt(r RoundingContext, signed bool, x uint64) ***REMOVED***
	if r.Increment > 0 ***REMOVED***
		// TODO: if uint64 is too large, fall back to float64
		if signed ***REMOVED***
			d.ConvertFloat(r, float64(int64(x)), 64)
		***REMOVED*** else ***REMOVED***
			d.ConvertFloat(r, float64(x), 64)
		***REMOVED***
		return
	***REMOVED***
	d.clear()
	if signed && int64(x) < 0 ***REMOVED***
		x = uint64(-int64(x))
		d.Neg = true
	***REMOVED***
	d.fillIntDigits(x)
	d.Exp = int32(len(d.Digits))
***REMOVED***

// ConvertFloat converts a floating point number to decimals.
func (d *Decimal) ConvertFloat(r RoundingContext, x float64, size int) ***REMOVED***
	d.clear()
	if math.IsNaN(x) ***REMOVED***
		d.NaN = true
		return
	***REMOVED***
	// Simple case: decimal notation
	if r.Increment > 0 ***REMOVED***
		scale := int(r.IncrementScale)
		mult := 1.0
		if scale > len(scales) ***REMOVED***
			mult = math.Pow(10, float64(scale))
		***REMOVED*** else ***REMOVED***
			mult = scales[scale]
		***REMOVED***
		// We multiply x instead of dividing inc as it gives less rounding
		// issues.
		x *= mult
		x /= float64(r.Increment)
		x = r.Mode.roundFloat(x)
		x *= float64(r.Increment)
		x /= mult
	***REMOVED***

	abs := x
	if x < 0 ***REMOVED***
		d.Neg = true
		abs = -x
	***REMOVED***
	if math.IsInf(abs, 1) ***REMOVED***
		d.Inf = true
		return
	***REMOVED***

	// By default we get the exact decimal representation.
	verb := byte('g')
	prec := -1
	// As the strconv API does not return the rounding accuracy, we can only
	// round using ToNearestEven.
	if r.Mode == ToNearestEven ***REMOVED***
		if n := r.RoundSignificantDigits(); n >= 0 ***REMOVED***
			prec = n
		***REMOVED*** else if n = r.RoundFractionDigits(); n >= 0 ***REMOVED***
			prec = n
			verb = 'f'
		***REMOVED***
	***REMOVED*** else ***REMOVED***
		// TODO: At this point strconv's rounding is imprecise to the point that
		// it is not useable for this purpose.
		// See https://github.com/golang/go/issues/21714
		// If rounding is requested, we ask for a large number of digits and
		// round from there to simulate rounding only once.
		// Ideally we would have strconv export an AppendDigits that would take
		// a rounding mode and/or return an accuracy. Something like this would
		// work:
		// AppendDigits(dst []byte, x float64, base, size, prec int) (digits []byte, exp, accuracy int)
		hasPrec := r.RoundSignificantDigits() >= 0
		hasScale := r.RoundFractionDigits() >= 0
		if hasPrec || hasScale ***REMOVED***
			// prec is the number of mantissa bits plus some extra for safety.
			// We need at least the number of mantissa bits as decimals to
			// accurately represent the floating point without rounding, as each
			// bit requires one more decimal to represent: 0.5, 0.25, 0.125, ...
			prec = 60
		***REMOVED***
	***REMOVED***

	b := strconv.AppendFloat(d.Digits[:0], abs, verb, prec, size)
	i := 0
	k := 0
	beforeDot := 1
	for i < len(b) ***REMOVED***
		if c := b[i]; '0' <= c && c <= '9' ***REMOVED***
			b[k] = c - '0'
			k++
			d.Exp += int32(beforeDot)
		***REMOVED*** else if c == '.' ***REMOVED***
			beforeDot = 0
			d.Exp = int32(k)
		***REMOVED*** else ***REMOVED***
			break
		***REMOVED***
		i++
	***REMOVED***
	d.Digits = b[:k]
	if i != len(b) ***REMOVED***
		i += len("e")
		pSign := i
		exp := 0
		for i++; i < len(b); i++ ***REMOVED***
			exp *= 10
			exp += int(b[i] - '0')
		***REMOVED***
		if b[pSign] == '-' ***REMOVED***
			exp = -exp
		***REMOVED***
		d.Exp = int32(exp) + 1
	***REMOVED***
***REMOVED***

func (d *Decimal) fillIntDigits(x uint64) ***REMOVED***
	if cap(d.Digits) < maxIntDigits ***REMOVED***
		d.Digits = d.buf[:]
	***REMOVED*** else ***REMOVED***
		d.Digits = d.buf[:maxIntDigits]
	***REMOVED***
	i := 0
	for ; x > 0; x /= 10 ***REMOVED***
		d.Digits[i] = byte(x % 10)
		i++
	***REMOVED***
	d.Digits = d.Digits[:i]
	for p := 0; p < i; p++ ***REMOVED***
		i--
		d.Digits[p], d.Digits[i] = d.Digits[i], d.Digits[p]
	***REMOVED***
***REMOVED***

var scales [70]float64

func init() ***REMOVED***
	x := 1.0
	for i := range scales ***REMOVED***
		scales[i] = x
		x *= 10
	***REMOVED***
***REMOVED***
