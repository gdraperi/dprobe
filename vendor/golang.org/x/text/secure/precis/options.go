// Copyright 2015 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package precis

import (
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
	"golang.org/x/text/runes"
	"golang.org/x/text/transform"
	"golang.org/x/text/unicode/norm"
)

// An Option is used to define the behavior and rules of a Profile.
type Option func(*options)

type options struct ***REMOVED***
	// Preparation options
	foldWidth bool

	// Enforcement options
	asciiLower    bool
	cases         transform.SpanningTransformer
	disallow      runes.Set
	norm          transform.SpanningTransformer
	additional    []func() transform.SpanningTransformer
	width         transform.SpanningTransformer
	disallowEmpty bool
	bidiRule      bool
	repeat        bool

	// Comparison options
	ignorecase bool
***REMOVED***

func getOpts(o ...Option) (res options) ***REMOVED***
	for _, f := range o ***REMOVED***
		f(&res)
	***REMOVED***
	// Using a SpanningTransformer, instead of norm.Form prevents an allocation
	// down the road.
	if res.norm == nil ***REMOVED***
		res.norm = norm.NFC
	***REMOVED***
	return
***REMOVED***

var (
	// The IgnoreCase option causes the profile to perform a case insensitive
	// comparison during the PRECIS comparison step.
	IgnoreCase Option = ignoreCase

	// The FoldWidth option causes the profile to map non-canonical wide and
	// narrow variants to their decomposition mapping. This is useful for
	// profiles that are based on the identifier class which would otherwise
	// disallow such characters.
	FoldWidth Option = foldWidth

	// The DisallowEmpty option causes the enforcement step to return an error if
	// the resulting string would be empty.
	DisallowEmpty Option = disallowEmpty

	// The BidiRule option causes the Bidi Rule defined in RFC 5893 to be
	// applied.
	BidiRule Option = bidiRule
)

var (
	ignoreCase = func(o *options) ***REMOVED***
		o.ignorecase = true
	***REMOVED***
	foldWidth = func(o *options) ***REMOVED***
		o.foldWidth = true
	***REMOVED***
	disallowEmpty = func(o *options) ***REMOVED***
		o.disallowEmpty = true
	***REMOVED***
	bidiRule = func(o *options) ***REMOVED***
		o.bidiRule = true
	***REMOVED***
	repeat = func(o *options) ***REMOVED***
		o.repeat = true
	***REMOVED***
)

// TODO: move this logic to package transform

type spanWrap struct***REMOVED*** transform.Transformer ***REMOVED***

func (s spanWrap) Span(src []byte, atEOF bool) (n int, err error) ***REMOVED***
	return 0, transform.ErrEndOfSpan
***REMOVED***

// TODO: allow different types? For instance:
//     func() transform.Transformer
//     func() transform.SpanningTransformer
//     func([]byte) bool  // validation only
//
// Also, would be great if we could detect if a transformer is reentrant.

// The AdditionalMapping option defines the additional mapping rule for the
// Profile by applying Transformer's in sequence.
func AdditionalMapping(t ...func() transform.Transformer) Option ***REMOVED***
	return func(o *options) ***REMOVED***
		for _, f := range t ***REMOVED***
			sf := func() transform.SpanningTransformer ***REMOVED***
				return f().(transform.SpanningTransformer)
			***REMOVED***
			if _, ok := f().(transform.SpanningTransformer); !ok ***REMOVED***
				sf = func() transform.SpanningTransformer ***REMOVED***
					return spanWrap***REMOVED***f()***REMOVED***
				***REMOVED***
			***REMOVED***
			o.additional = append(o.additional, sf)
		***REMOVED***
	***REMOVED***
***REMOVED***

// The Norm option defines a Profile's normalization rule. Defaults to NFC.
func Norm(f norm.Form) Option ***REMOVED***
	return func(o *options) ***REMOVED***
		o.norm = f
	***REMOVED***
***REMOVED***

// The FoldCase option defines a Profile's case mapping rule. Options can be
// provided to determine the type of case folding used.
func FoldCase(opts ...cases.Option) Option ***REMOVED***
	return func(o *options) ***REMOVED***
		o.asciiLower = true
		o.cases = cases.Fold(opts...)
	***REMOVED***
***REMOVED***

// The LowerCase option defines a Profile's case mapping rule. Options can be
// provided to determine the type of case folding used.
func LowerCase(opts ...cases.Option) Option ***REMOVED***
	return func(o *options) ***REMOVED***
		o.asciiLower = true
		if len(opts) == 0 ***REMOVED***
			o.cases = cases.Lower(language.Und, cases.HandleFinalSigma(false))
			return
		***REMOVED***

		opts = append([]cases.Option***REMOVED***cases.HandleFinalSigma(false)***REMOVED***, opts...)
		o.cases = cases.Lower(language.Und, opts...)
	***REMOVED***
***REMOVED***

// The Disallow option further restricts a Profile's allowed characters beyond
// what is disallowed by the underlying string class.
func Disallow(set runes.Set) Option ***REMOVED***
	return func(o *options) ***REMOVED***
		o.disallow = set
	***REMOVED***
***REMOVED***
