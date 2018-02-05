// Copyright 2017 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package number

import (
	"fmt"
	"math"
	"strconv"
	"strings"
	"testing"
)

func mkfloat(num string) float64 ***REMOVED***
	u, _ := strconv.ParseUint(num, 10, 32)
	return float64(u)
***REMOVED***

// mkdec creates a decimal from a string. All ASCII digits are converted to
// digits in the decimal. The dot is used to indicate the scale by which the
// digits are shifted. Numbers may have an additional exponent or be the special
// value NaN, Inf, or -Inf.
func mkdec(num string) (d Decimal) ***REMOVED***
	var r RoundingContext
	d.Convert(r, dec(num))
	return
***REMOVED***

type dec string

func (s dec) Convert(d *Decimal, _ RoundingContext) ***REMOVED***
	num := string(s)
	if num[0] == '-' ***REMOVED***
		d.Neg = true
		num = num[1:]
	***REMOVED***
	switch num ***REMOVED***
	case "NaN":
		d.NaN = true
		return
	case "Inf":
		d.Inf = true
		return
	***REMOVED***
	if p := strings.IndexAny(num, "eE"); p != -1 ***REMOVED***
		i64, err := strconv.ParseInt(num[p+1:], 10, 32)
		if err != nil ***REMOVED***
			panic(err)
		***REMOVED***
		d.Exp = int32(i64)
		num = num[:p]
	***REMOVED***
	if p := strings.IndexByte(num, '.'); p != -1 ***REMOVED***
		d.Exp += int32(p)
		num = num[:p] + num[p+1:]
	***REMOVED*** else ***REMOVED***
		d.Exp += int32(len(num))
	***REMOVED***
	d.Digits = []byte(num)
	for i := range d.Digits ***REMOVED***
		d.Digits[i] -= '0'
	***REMOVED***
	*d = d.normalize()
***REMOVED***

func byteNum(s string) []byte ***REMOVED***
	b := make([]byte, len(s))
	for i := 0; i < len(s); i++ ***REMOVED***
		if c := s[i]; '0' <= c && c <= '9' ***REMOVED***
			b[i] = s[i] - '0'
		***REMOVED*** else ***REMOVED***
			b[i] = s[i] - 'a' + 10
		***REMOVED***
	***REMOVED***
	return b
***REMOVED***

func strNum(s string) string ***REMOVED***
	return string(byteNum(s))
***REMOVED***

func TestDecimalString(t *testing.T) ***REMOVED***
	for _, test := range []struct ***REMOVED***
		x    Decimal
		want string
	***REMOVED******REMOVED***
		***REMOVED***want: "0"***REMOVED***,
		***REMOVED***Decimal***REMOVED***digits: digits***REMOVED***Digits: nil, Exp: 1000***REMOVED******REMOVED***, "0"***REMOVED***, // exponent of 1000 is ignored
		***REMOVED***Decimal***REMOVED***digits: digits***REMOVED***Digits: byteNum("12345"), Exp: 0***REMOVED******REMOVED***, "0.12345"***REMOVED***,
		***REMOVED***Decimal***REMOVED***digits: digits***REMOVED***Digits: byteNum("12345"), Exp: -3***REMOVED******REMOVED***, "0.00012345"***REMOVED***,
		***REMOVED***Decimal***REMOVED***digits: digits***REMOVED***Digits: byteNum("12345"), Exp: +3***REMOVED******REMOVED***, "123.45"***REMOVED***,
		***REMOVED***Decimal***REMOVED***digits: digits***REMOVED***Digits: byteNum("12345"), Exp: +10***REMOVED******REMOVED***, "1234500000"***REMOVED***,
	***REMOVED*** ***REMOVED***
		if got := test.x.String(); got != test.want ***REMOVED***
			t.Errorf("%v == %q; want %q", test.x, got, test.want)
		***REMOVED***
	***REMOVED***
***REMOVED***

func TestRounding(t *testing.T) ***REMOVED***
	testCases := []struct ***REMOVED***
		x string
		n int
		// modes is the result for modes. Signs are left out of the result.
		// The results are stored in the following order:
		// zero, negInf
		// nearZero, nearEven, nearAway
		// away, posInf
		modes [numModes]string
	***REMOVED******REMOVED***
		***REMOVED***"0", 1, [numModes]string***REMOVED***
			"0", "0",
			"0", "0", "0",
			"0", "0"***REMOVED******REMOVED***,
		***REMOVED***"1", 1, [numModes]string***REMOVED***
			"1", "1",
			"1", "1", "1",
			"1", "1"***REMOVED******REMOVED***,
		***REMOVED***"5", 1, [numModes]string***REMOVED***
			"5", "5",
			"5", "5", "5",
			"5", "5"***REMOVED******REMOVED***,
		***REMOVED***"15", 1, [numModes]string***REMOVED***
			"10", "10",
			"10", "20", "20",
			"20", "20"***REMOVED******REMOVED***,
		***REMOVED***"45", 1, [numModes]string***REMOVED***
			"40", "40",
			"40", "40", "50",
			"50", "50"***REMOVED******REMOVED***,
		***REMOVED***"95", 1, [numModes]string***REMOVED***
			"90", "90",
			"90", "100", "100",
			"100", "100"***REMOVED******REMOVED***,

		***REMOVED***"12344999", 4, [numModes]string***REMOVED***
			"12340000", "12340000",
			"12340000", "12340000", "12340000",
			"12350000", "12350000"***REMOVED******REMOVED***,
		***REMOVED***"12345000", 4, [numModes]string***REMOVED***
			"12340000", "12340000",
			"12340000", "12340000", "12350000",
			"12350000", "12350000"***REMOVED******REMOVED***,
		***REMOVED***"12345001", 4, [numModes]string***REMOVED***
			"12340000", "12340000",
			"12350000", "12350000", "12350000",
			"12350000", "12350000"***REMOVED******REMOVED***,
		***REMOVED***"12345100", 4, [numModes]string***REMOVED***
			"12340000", "12340000",
			"12350000", "12350000", "12350000",
			"12350000", "12350000"***REMOVED******REMOVED***,
		***REMOVED***"23454999", 4, [numModes]string***REMOVED***
			"23450000", "23450000",
			"23450000", "23450000", "23450000",
			"23460000", "23460000"***REMOVED******REMOVED***,
		***REMOVED***"23455000", 4, [numModes]string***REMOVED***
			"23450000", "23450000",
			"23450000", "23460000", "23460000",
			"23460000", "23460000"***REMOVED******REMOVED***,
		***REMOVED***"23455001", 4, [numModes]string***REMOVED***
			"23450000", "23450000",
			"23460000", "23460000", "23460000",
			"23460000", "23460000"***REMOVED******REMOVED***,
		***REMOVED***"23455100", 4, [numModes]string***REMOVED***
			"23450000", "23450000",
			"23460000", "23460000", "23460000",
			"23460000", "23460000"***REMOVED******REMOVED***,

		***REMOVED***"99994999", 4, [numModes]string***REMOVED***
			"99990000", "99990000",
			"99990000", "99990000", "99990000",
			"100000000", "100000000"***REMOVED******REMOVED***,
		***REMOVED***"99995000", 4, [numModes]string***REMOVED***
			"99990000", "99990000",
			"99990000", "100000000", "100000000",
			"100000000", "100000000"***REMOVED******REMOVED***,
		***REMOVED***"99999999", 4, [numModes]string***REMOVED***
			"99990000", "99990000",
			"100000000", "100000000", "100000000",
			"100000000", "100000000"***REMOVED******REMOVED***,

		***REMOVED***"12994999", 4, [numModes]string***REMOVED***
			"12990000", "12990000",
			"12990000", "12990000", "12990000",
			"13000000", "13000000"***REMOVED******REMOVED***,
		***REMOVED***"12995000", 4, [numModes]string***REMOVED***
			"12990000", "12990000",
			"12990000", "13000000", "13000000",
			"13000000", "13000000"***REMOVED******REMOVED***,
		***REMOVED***"12999999", 4, [numModes]string***REMOVED***
			"12990000", "12990000",
			"13000000", "13000000", "13000000",
			"13000000", "13000000"***REMOVED******REMOVED***,
	***REMOVED***
	modes := []RoundingMode***REMOVED***
		ToZero, ToNegativeInf,
		ToNearestZero, ToNearestEven, ToNearestAway,
		AwayFromZero, ToPositiveInf,
	***REMOVED***
	for _, tc := range testCases ***REMOVED***
		// Create negative counterpart tests: the sign is reversed and
		// ToPositiveInf and ToNegativeInf swapped.
		negModes := tc.modes
		negModes[1], negModes[6] = negModes[6], negModes[1]
		for i, res := range negModes ***REMOVED***
			negModes[i] = "-" + res
		***REMOVED***
		for i, m := range modes ***REMOVED***
			t.Run(fmt.Sprintf("x:%s/n:%d/%s", tc.x, tc.n, m), func(t *testing.T) ***REMOVED***
				d := mkdec(tc.x)
				d.round(m, tc.n)
				if got := d.String(); got != tc.modes[i] ***REMOVED***
					t.Errorf("pos decimal: got %q; want %q", d.String(), tc.modes[i])
				***REMOVED***

				mult := math.Pow(10, float64(len(tc.x)-tc.n))
				f := mkfloat(tc.x)
				f = m.roundFloat(f/mult) * mult
				if got := fmt.Sprintf("%.0f", f); got != tc.modes[i] ***REMOVED***
					t.Errorf("pos float: got %q; want %q", got, tc.modes[i])
				***REMOVED***

				// Test the negative case. This is the same as the positive
				// case, but with ToPositiveInf and ToNegativeInf swapped.
				d = mkdec(tc.x)
				d.Neg = true
				d.round(m, tc.n)
				if got, want := d.String(), negModes[i]; got != want ***REMOVED***
					t.Errorf("neg decimal: got %q; want %q", d.String(), want)
				***REMOVED***

				f = -mkfloat(tc.x)
				f = m.roundFloat(f/mult) * mult
				if got := fmt.Sprintf("%.0f", f); got != negModes[i] ***REMOVED***
					t.Errorf("neg float: got %q; want %q", got, negModes[i])
				***REMOVED***
			***REMOVED***)
		***REMOVED***
	***REMOVED***
***REMOVED***

func TestConvert(t *testing.T) ***REMOVED***
	scale2 := RoundingContext***REMOVED******REMOVED***
	scale2.SetScale(2)
	scale2away := RoundingContext***REMOVED***Mode: AwayFromZero***REMOVED***
	scale2away.SetScale(2)
	inc0_05 := RoundingContext***REMOVED***Increment: 5, IncrementScale: 2***REMOVED***
	inc0_05.SetScale(2)
	inc50 := RoundingContext***REMOVED***Increment: 50***REMOVED***
	prec3 := RoundingContext***REMOVED******REMOVED***
	prec3.SetPrecision(3)
	roundShift := RoundingContext***REMOVED***DigitShift: 2, MaxFractionDigits: 2***REMOVED***
	testCases := []struct ***REMOVED***
		x   interface***REMOVED******REMOVED***
		rc  RoundingContext
		out string
	***REMOVED******REMOVED***
		***REMOVED***-0.001, scale2, "-0.00"***REMOVED***,
		***REMOVED***0.1234, prec3, "0.123"***REMOVED***,
		***REMOVED***1234.0, prec3, "1230"***REMOVED***,
		***REMOVED***1.2345e10, prec3, "12300000000"***REMOVED***,

		***REMOVED***int8(-34), scale2, "-34"***REMOVED***,
		***REMOVED***int16(-234), scale2, "-234"***REMOVED***,
		***REMOVED***int32(-234), scale2, "-234"***REMOVED***,
		***REMOVED***int64(-234), scale2, "-234"***REMOVED***,
		***REMOVED***int(-234), scale2, "-234"***REMOVED***,
		***REMOVED***uint8(234), scale2, "234"***REMOVED***,
		***REMOVED***uint16(234), scale2, "234"***REMOVED***,
		***REMOVED***uint32(234), scale2, "234"***REMOVED***,
		***REMOVED***uint64(234), scale2, "234"***REMOVED***,
		***REMOVED***uint(234), scale2, "234"***REMOVED***,
		***REMOVED***-1e9, scale2, "-1000000000.00"***REMOVED***,
		// The following two causes this result to have a lot of digits:
		// 1) 0.234 cannot be accurately represented as a float64, and
		// 2) as strconv does not support the rounding AwayFromZero, Convert
		//    leaves the rounding to caller.
		***REMOVED***0.234, scale2away,
			"0.2340000000000000135447209004269097931683063507080078125"***REMOVED***,

		***REMOVED***0.0249, inc0_05, "0.00"***REMOVED***,
		***REMOVED***0.025, inc0_05, "0.00"***REMOVED***,
		***REMOVED***0.0251, inc0_05, "0.05"***REMOVED***,
		***REMOVED***0.03, inc0_05, "0.05"***REMOVED***,
		***REMOVED***0.049, inc0_05, "0.05"***REMOVED***,
		***REMOVED***0.05, inc0_05, "0.05"***REMOVED***,
		***REMOVED***0.051, inc0_05, "0.05"***REMOVED***,
		***REMOVED***0.0749, inc0_05, "0.05"***REMOVED***,
		***REMOVED***0.075, inc0_05, "0.10"***REMOVED***,
		***REMOVED***0.0751, inc0_05, "0.10"***REMOVED***,
		***REMOVED***324, inc50, "300"***REMOVED***,
		***REMOVED***325, inc50, "300"***REMOVED***,
		***REMOVED***326, inc50, "350"***REMOVED***,
		***REMOVED***349, inc50, "350"***REMOVED***,
		***REMOVED***350, inc50, "350"***REMOVED***,
		***REMOVED***351, inc50, "350"***REMOVED***,
		***REMOVED***374, inc50, "350"***REMOVED***,
		***REMOVED***375, inc50, "400"***REMOVED***,
		***REMOVED***376, inc50, "400"***REMOVED***,

		// Here the scale is 2, but the digits get shifted left. As we use
		// AppendFloat to do the rounding an exta 0 gets added.
		***REMOVED***0.123, roundShift, "0.1230"***REMOVED***,

		***REMOVED***converter(3), scale2, "100"***REMOVED***,

		***REMOVED***math.Inf(1), inc50, "Inf"***REMOVED***,
		***REMOVED***math.Inf(-1), inc50, "-Inf"***REMOVED***,
		***REMOVED***math.NaN(), inc50, "NaN"***REMOVED***,
		***REMOVED***"clearly not a number", scale2, "NaN"***REMOVED***,
	***REMOVED***
	for _, tc := range testCases ***REMOVED***
		var d Decimal
		t.Run(fmt.Sprintf("%T:%v-%v", tc.x, tc.x, tc.rc), func(t *testing.T) ***REMOVED***
			d.Convert(tc.rc, tc.x)
			if got := d.String(); got != tc.out ***REMOVED***
				t.Errorf("got %q; want %q", got, tc.out)
			***REMOVED***
		***REMOVED***)
	***REMOVED***
***REMOVED***

type converter int

func (c converter) Convert(d *Decimal, r RoundingContext) ***REMOVED***
	d.Digits = append(d.Digits, 1, 0, 0)
	d.Exp = 3
***REMOVED***
