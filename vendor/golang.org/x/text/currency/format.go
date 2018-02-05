// Copyright 2015 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package currency

import (
	"fmt"
	"io"
	"sort"

	"golang.org/x/text/internal"
	"golang.org/x/text/internal/format"
	"golang.org/x/text/language"
)

// Amount is an amount-currency unit pair.
type Amount struct ***REMOVED***
	amount   interface***REMOVED******REMOVED*** // Change to decimal(64|128).
	currency Unit
***REMOVED***

// Currency reports the currency unit of this amount.
func (a Amount) Currency() Unit ***REMOVED*** return a.currency ***REMOVED***

// TODO: based on decimal type, but may make sense to customize a bit.
// func (a Amount) Decimal()
// func (a Amount) Int() (int64, error)
// func (a Amount) Fraction() (int64, error)
// func (a Amount) Rat() *big.Rat
// func (a Amount) Float() (float64, error)
// func (a Amount) Scale() uint
// func (a Amount) Precision() uint
// func (a Amount) Sign() int
//
// Add/Sub/Div/Mul/Round.

var space = []byte(" ")

// Format implements fmt.Formatter. It accepts format.State for
// language-specific rendering.
func (a Amount) Format(s fmt.State, verb rune) ***REMOVED***
	v := formattedValue***REMOVED***
		currency: a.currency,
		amount:   a.amount,
		format:   defaultFormat,
	***REMOVED***
	v.Format(s, verb)
***REMOVED***

// formattedValue is currency amount or unit that implements language-sensitive
// formatting.
type formattedValue struct ***REMOVED***
	currency Unit
	amount   interface***REMOVED******REMOVED*** // Amount, Unit, or number.
	format   *options
***REMOVED***

// Format implements fmt.Formatter. It accepts format.State for
// language-specific rendering.
func (v formattedValue) Format(s fmt.State, verb rune) ***REMOVED***
	var lang int
	if state, ok := s.(format.State); ok ***REMOVED***
		lang, _ = language.CompactIndex(state.Language())
	***REMOVED***

	// Get the options. Use DefaultFormat if not present.
	opt := v.format
	if opt == nil ***REMOVED***
		opt = defaultFormat
	***REMOVED***
	cur := v.currency
	if cur.index == 0 ***REMOVED***
		cur = opt.currency
	***REMOVED***

	// TODO: use pattern.
	io.WriteString(s, opt.symbol(lang, cur))
	if v.amount != nil ***REMOVED***
		s.Write(space)

		// TODO: apply currency-specific rounding
		scale, _ := opt.kind.Rounding(cur)
		if _, ok := s.Precision(); !ok ***REMOVED***
			fmt.Fprintf(s, "%.*f", scale, v.amount)
		***REMOVED*** else ***REMOVED***
			fmt.Fprint(s, v.amount)
		***REMOVED***
	***REMOVED***
***REMOVED***

// Formatter decorates a given number, Unit or Amount with formatting options.
type Formatter func(amount interface***REMOVED******REMOVED***) formattedValue

// func (f Formatter) Options(opts ...Option) Formatter

// TODO: call this a Formatter or FormatFunc?

var dummy = USD.Amount(0)

// adjust creates a new Formatter based on the adjustments of fn on f.
func (f Formatter) adjust(fn func(*options)) Formatter ***REMOVED***
	var o options = *(f(dummy).format)
	fn(&o)
	return o.format
***REMOVED***

// Default creates a new Formatter that defaults to currency unit c if a numeric
// value is passed that is not associated with a currency.
func (f Formatter) Default(currency Unit) Formatter ***REMOVED***
	return f.adjust(func(o *options) ***REMOVED*** o.currency = currency ***REMOVED***)
***REMOVED***

// Kind sets the kind of the underlying currency unit.
func (f Formatter) Kind(k Kind) Formatter ***REMOVED***
	return f.adjust(func(o *options) ***REMOVED*** o.kind = k ***REMOVED***)
***REMOVED***

var defaultFormat *options = ISO(dummy).format

var (
	// Uses Narrow symbols. Overrides Symbol, if present.
	NarrowSymbol Formatter = Formatter(formNarrow)

	// Use Symbols instead of ISO codes, when available.
	Symbol Formatter = Formatter(formSymbol)

	// Use ISO code as symbol.
	ISO Formatter = Formatter(formISO)

	// TODO:
	// // Use full name as symbol.
	// Name Formatter
)

// options configures rendering and rounding options for an Amount.
type options struct ***REMOVED***
	currency Unit
	kind     Kind

	symbol func(compactIndex int, c Unit) string
***REMOVED***

func (o *options) format(amount interface***REMOVED******REMOVED***) formattedValue ***REMOVED***
	v := formattedValue***REMOVED***format: o***REMOVED***
	switch x := amount.(type) ***REMOVED***
	case Amount:
		v.amount = x.amount
		v.currency = x.currency
	case *Amount:
		v.amount = x.amount
		v.currency = x.currency
	case Unit:
		v.currency = x
	case *Unit:
		v.currency = *x
	default:
		if o.currency.index == 0 ***REMOVED***
			panic("cannot format number without a currency being set")
		***REMOVED***
		// TODO: Must be a number.
		v.amount = x
		v.currency = o.currency
	***REMOVED***
	return v
***REMOVED***

var (
	optISO    = options***REMOVED***symbol: lookupISO***REMOVED***
	optSymbol = options***REMOVED***symbol: lookupSymbol***REMOVED***
	optNarrow = options***REMOVED***symbol: lookupNarrow***REMOVED***
)

// These need to be functions, rather than curried methods, as curried methods
// are evaluated at init time, causing tables to be included unconditionally.
func formISO(x interface***REMOVED******REMOVED***) formattedValue    ***REMOVED*** return optISO.format(x) ***REMOVED***
func formSymbol(x interface***REMOVED******REMOVED***) formattedValue ***REMOVED*** return optSymbol.format(x) ***REMOVED***
func formNarrow(x interface***REMOVED******REMOVED***) formattedValue ***REMOVED*** return optNarrow.format(x) ***REMOVED***

func lookupISO(x int, c Unit) string    ***REMOVED*** return c.String() ***REMOVED***
func lookupSymbol(x int, c Unit) string ***REMOVED*** return normalSymbol.lookup(x, c) ***REMOVED***
func lookupNarrow(x int, c Unit) string ***REMOVED*** return narrowSymbol.lookup(x, c) ***REMOVED***

type symbolIndex struct ***REMOVED***
	index []uint16 // position corresponds with compact index of language.
	data  []curToIndex
***REMOVED***

var (
	normalSymbol = symbolIndex***REMOVED***normalLangIndex, normalSymIndex***REMOVED***
	narrowSymbol = symbolIndex***REMOVED***narrowLangIndex, narrowSymIndex***REMOVED***
)

func (x *symbolIndex) lookup(lang int, c Unit) string ***REMOVED***
	for ***REMOVED***
		index := x.data[x.index[lang]:x.index[lang+1]]
		i := sort.Search(len(index), func(i int) bool ***REMOVED***
			return index[i].cur >= c.index
		***REMOVED***)
		if i < len(index) && index[i].cur == c.index ***REMOVED***
			x := index[i].idx
			start := x + 1
			end := start + uint16(symbols[x])
			if start == end ***REMOVED***
				return c.String()
			***REMOVED***
			return symbols[start:end]
		***REMOVED***
		if lang == 0 ***REMOVED***
			break
		***REMOVED***
		lang = int(internal.Parent[lang])
	***REMOVED***
	return c.String()
***REMOVED***
