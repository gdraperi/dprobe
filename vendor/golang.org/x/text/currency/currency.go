// Copyright 2015 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:generate go run gen.go gen_common.go -output tables.go

// Package currency contains currency-related functionality.
//
// NOTE: the formatting functionality is currently under development and may
// change without notice.
package currency // import "golang.org/x/text/currency"

import (
	"errors"
	"sort"

	"golang.org/x/text/internal/tag"
	"golang.org/x/text/language"
)

// TODO:
// - language-specific currency names.
// - currency formatting.
// - currency information per region
// - register currency code (there are no private use area)

// TODO: remove Currency type from package language.

// Kind determines the rounding and rendering properties of a currency value.
type Kind struct ***REMOVED***
	rounding rounding
	// TODO: formatting type: standard, accounting. See CLDR.
***REMOVED***

type rounding byte

const (
	standard rounding = iota
	cash
)

var (
	// Standard defines standard rounding and formatting for currencies.
	Standard Kind = Kind***REMOVED***rounding: standard***REMOVED***

	// Cash defines rounding and formatting standards for cash transactions.
	Cash Kind = Kind***REMOVED***rounding: cash***REMOVED***

	// Accounting defines rounding and formatting standards for accounting.
	Accounting Kind = Kind***REMOVED***rounding: standard***REMOVED***
)

// Rounding reports the rounding characteristics for the given currency, where
// scale is the number of fractional decimals and increment is the number of
// units in terms of 10^(-scale) to which to round to.
func (k Kind) Rounding(cur Unit) (scale, increment int) ***REMOVED***
	info := currency.Elem(int(cur.index))[3]
	switch k.rounding ***REMOVED***
	case standard:
		info &= roundMask
	case cash:
		info >>= cashShift
	***REMOVED***
	return int(roundings[info].scale), int(roundings[info].increment)
***REMOVED***

// Unit is an ISO 4217 currency designator.
type Unit struct ***REMOVED***
	index uint16
***REMOVED***

// String returns the ISO code of u.
func (u Unit) String() string ***REMOVED***
	if u.index == 0 ***REMOVED***
		return "XXX"
	***REMOVED***
	return currency.Elem(int(u.index))[:3]
***REMOVED***

// Amount creates an Amount for the given currency unit and amount.
func (u Unit) Amount(amount interface***REMOVED******REMOVED***) Amount ***REMOVED***
	// TODO: verify amount is a supported number type
	return Amount***REMOVED***amount: amount, currency: u***REMOVED***
***REMOVED***

var (
	errSyntax = errors.New("currency: tag is not well-formed")
	errValue  = errors.New("currency: tag is not a recognized currency")
)

// ParseISO parses a 3-letter ISO 4217 currency code. It returns an error if s
// is not well-formed or not a recognized currency code.
func ParseISO(s string) (Unit, error) ***REMOVED***
	var buf [4]byte // Take one byte more to detect oversize keys.
	key := buf[:copy(buf[:], s)]
	if !tag.FixCase("XXX", key) ***REMOVED***
		return Unit***REMOVED******REMOVED***, errSyntax
	***REMOVED***
	if i := currency.Index(key); i >= 0 ***REMOVED***
		if i == xxx ***REMOVED***
			return Unit***REMOVED******REMOVED***, nil
		***REMOVED***
		return Unit***REMOVED***uint16(i)***REMOVED***, nil
	***REMOVED***
	return Unit***REMOVED******REMOVED***, errValue
***REMOVED***

// MustParseISO is like ParseISO, but panics if the given currency unit
// cannot be parsed. It simplifies safe initialization of Unit values.
func MustParseISO(s string) Unit ***REMOVED***
	c, err := ParseISO(s)
	if err != nil ***REMOVED***
		panic(err)
	***REMOVED***
	return c
***REMOVED***

// FromRegion reports the currency unit that is currently legal tender in the
// given region according to CLDR. It will return false if region currently does
// not have a legal tender.
func FromRegion(r language.Region) (currency Unit, ok bool) ***REMOVED***
	x := regionToCode(r)
	i := sort.Search(len(regionToCurrency), func(i int) bool ***REMOVED***
		return regionToCurrency[i].region >= x
	***REMOVED***)
	if i < len(regionToCurrency) && regionToCurrency[i].region == x ***REMOVED***
		return Unit***REMOVED***regionToCurrency[i].code***REMOVED***, true
	***REMOVED***
	return Unit***REMOVED******REMOVED***, false
***REMOVED***

// FromTag reports the most likely currency for the given tag. It considers the
// currency defined in the -u extension and infers the region if necessary.
func FromTag(t language.Tag) (Unit, language.Confidence) ***REMOVED***
	if cur := t.TypeForKey("cu"); len(cur) == 3 ***REMOVED***
		c, _ := ParseISO(cur)
		return c, language.Exact
	***REMOVED***
	r, conf := t.Region()
	if cur, ok := FromRegion(r); ok ***REMOVED***
		return cur, conf
	***REMOVED***
	return Unit***REMOVED******REMOVED***, language.No
***REMOVED***

var (
	// Undefined and testing.
	XXX Unit = Unit***REMOVED******REMOVED***
	XTS Unit = Unit***REMOVED***xts***REMOVED***

	// G10 currencies https://en.wikipedia.org/wiki/G10_currencies.
	USD Unit = Unit***REMOVED***usd***REMOVED***
	EUR Unit = Unit***REMOVED***eur***REMOVED***
	JPY Unit = Unit***REMOVED***jpy***REMOVED***
	GBP Unit = Unit***REMOVED***gbp***REMOVED***
	CHF Unit = Unit***REMOVED***chf***REMOVED***
	AUD Unit = Unit***REMOVED***aud***REMOVED***
	NZD Unit = Unit***REMOVED***nzd***REMOVED***
	CAD Unit = Unit***REMOVED***cad***REMOVED***
	SEK Unit = Unit***REMOVED***sek***REMOVED***
	NOK Unit = Unit***REMOVED***nok***REMOVED***

	// Additional common currencies as defined by CLDR.
	BRL Unit = Unit***REMOVED***brl***REMOVED***
	CNY Unit = Unit***REMOVED***cny***REMOVED***
	DKK Unit = Unit***REMOVED***dkk***REMOVED***
	INR Unit = Unit***REMOVED***inr***REMOVED***
	RUB Unit = Unit***REMOVED***rub***REMOVED***
	HKD Unit = Unit***REMOVED***hkd***REMOVED***
	IDR Unit = Unit***REMOVED***idr***REMOVED***
	KRW Unit = Unit***REMOVED***krw***REMOVED***
	MXN Unit = Unit***REMOVED***mxn***REMOVED***
	PLN Unit = Unit***REMOVED***pln***REMOVED***
	SAR Unit = Unit***REMOVED***sar***REMOVED***
	THB Unit = Unit***REMOVED***thb***REMOVED***
	TRY Unit = Unit***REMOVED***try***REMOVED***
	TWD Unit = Unit***REMOVED***twd***REMOVED***
	ZAR Unit = Unit***REMOVED***zar***REMOVED***

	// Precious metals.
	XAG Unit = Unit***REMOVED***xag***REMOVED***
	XAU Unit = Unit***REMOVED***xau***REMOVED***
	XPT Unit = Unit***REMOVED***xpt***REMOVED***
	XPD Unit = Unit***REMOVED***xpd***REMOVED***
)
