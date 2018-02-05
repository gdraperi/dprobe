// Copyright 2017 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package number

import (
	"fmt"
	"strings"

	"golang.org/x/text/feature/plural"
	"golang.org/x/text/internal/format"
	"golang.org/x/text/internal/number"
	"golang.org/x/text/language"
)

// A FormatFunc formates a number.
type FormatFunc func(x interface***REMOVED******REMOVED***, opts ...Option) Formatter

// NewFormat creates a FormatFunc based on another FormatFunc and new options.
// Use NewFormat to cash the creation of formatters.
func NewFormat(format FormatFunc, opts ...Option) FormatFunc ***REMOVED***
	o := *format(nil).options
	n := len(o.options)
	o.options = append(o.options[:n:n], opts...)
	return func(x interface***REMOVED******REMOVED***, opts ...Option) Formatter ***REMOVED***
		return newFormatter(&o, opts, x)
	***REMOVED***
***REMOVED***

type options struct ***REMOVED***
	verbs      string
	initFunc   initFunc
	options    []Option
	pluralFunc func(t language.Tag, scale int) (f plural.Form, n int)
***REMOVED***

type optionFlag uint16

const (
	hasScale optionFlag = 1 << iota
	hasPrecision
	noSeparator
	exact
)

type initFunc func(f *number.Formatter, t language.Tag)

func newFormatter(o *options, opts []Option, value interface***REMOVED******REMOVED***) Formatter ***REMOVED***
	if len(opts) > 0 ***REMOVED***
		n := *o
		n.options = opts
		o = &n
	***REMOVED***
	return Formatter***REMOVED***o, value***REMOVED***
***REMOVED***

func newOptions(verbs string, f initFunc) *options ***REMOVED***
	return &options***REMOVED***verbs: verbs, initFunc: f***REMOVED***
***REMOVED***

type Formatter struct ***REMOVED***
	*options
	value interface***REMOVED******REMOVED***
***REMOVED***

// Format implements format.Formatter. It is for internal use only for now.
func (f Formatter) Format(state format.State, verb rune) ***REMOVED***
	// TODO: consider implementing fmt.Formatter instead and using the following
	// piece of code. This allows numbers to be rendered mostly as expected
	// when using fmt. But it may get weird with the spellout options and we
	// may need more of format.State over time.
	// lang := language.Und
	// if s, ok := state.(format.State); ok ***REMOVED***
	// 	lang = s.Language()
	// ***REMOVED***

	lang := state.Language()
	if !strings.Contains(f.verbs, string(verb)) ***REMOVED***
		fmt.Fprintf(state, "%%!%s(%T=%v)", string(verb), f.value, f.value)
		return
	***REMOVED***
	var p number.Formatter
	f.initFunc(&p, lang)
	for _, o := range f.options.options ***REMOVED***
		o(lang, &p)
	***REMOVED***
	if w, ok := state.Width(); ok ***REMOVED***
		p.FormatWidth = uint16(w)
	***REMOVED***
	if prec, ok := state.Precision(); ok ***REMOVED***
		switch verb ***REMOVED***
		case 'd':
			p.SetScale(0)
		case 'f':
			p.SetScale(prec)
		case 'e':
			p.SetPrecision(prec + 1)
		case 'g':
			p.SetPrecision(prec)
		***REMOVED***
	***REMOVED***
	var d number.Decimal
	d.Convert(p.RoundingContext, f.value)
	state.Write(p.Format(nil, &d))
***REMOVED***

// Digits returns information about which logical digits will be presented to
// the user. This information is relevant, for instance, to determine plural
// forms.
func (f Formatter) Digits(buf []byte, tag language.Tag, scale int) number.Digits ***REMOVED***
	var p number.Formatter
	f.initFunc(&p, tag)
	if scale >= 0 ***REMOVED***
		// TODO: this only works well for decimal numbers, which is generally
		// fine.
		p.SetScale(scale)
	***REMOVED***
	var d number.Decimal
	d.Convert(p.RoundingContext, f.value)
	return number.FormatDigits(&d, p.RoundingContext)
***REMOVED***
