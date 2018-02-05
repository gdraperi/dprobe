// Copyright 2017 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package plural

import (
	"fmt"
	"io/ioutil"
	"reflect"
	"strconv"

	"golang.org/x/text/internal/catmsg"
	"golang.org/x/text/internal/number"
	"golang.org/x/text/language"
	"golang.org/x/text/message/catalog"
)

// TODO: consider deleting this interface. Maybe VisibleDigits is always
// sufficient and practical.

// Interface is used for types that can determine their own plural form.
type Interface interface ***REMOVED***
	// PluralForm reports the plural form for the given language of the
	// underlying value. It also returns the integer value. If the integer value
	// is larger than fits in n, PluralForm may return a value modulo
	// 10,000,000.
	PluralForm(t language.Tag, scale int) (f Form, n int)
***REMOVED***

// Selectf returns the first case for which its selector is a match for the
// arg-th substitution argument to a formatting call, formatting it as indicated
// by format.
//
// The cases argument are pairs of selectors and messages. Selectors are of type
// string or Form. Messages are of type string or catalog.Message. A selector
// matches an argument if:
//    - it is "other" or Other
//    - it matches the plural form of the argument: "zero", "one", "two", "few",
//      or "many", or the equivalent Form
//    - it is of the form "=x" where x is an integer that matches the value of
//      the argument.
//    - it is of the form "<x" where x is an integer that is larger than the
//      argument.
//
// The format argument determines the formatting parameters for which to
// determine the plural form. This is especially relevant for non-integer
// values.
//
// The format string may be "", in which case a best-effort attempt is made to
// find a reasonable representation on which to base the plural form. Examples
// of format strings are:
//   - %.2f   decimal with scale 2
//   - %.2e   scientific notation with precision 3 (scale + 1)
//   - %d     integer
func Selectf(arg int, format string, cases ...interface***REMOVED******REMOVED***) catalog.Message ***REMOVED***
	var p parser
	// Intercept the formatting parameters of format by doing a dummy print.
	fmt.Fprintf(ioutil.Discard, format, &p)
	m := &message***REMOVED***arg, kindDefault, 0, cases***REMOVED***
	switch p.verb ***REMOVED***
	case 'g':
		m.kind = kindPrecision
		m.scale = p.scale
	case 'f':
		m.kind = kindScale
		m.scale = p.scale
	case 'e':
		m.kind = kindScientific
		m.scale = p.scale
	case 'd':
		m.kind = kindScale
		m.scale = 0
	default:
		// TODO: do we need to handle errors?
	***REMOVED***
	return m
***REMOVED***

type parser struct ***REMOVED***
	verb  rune
	scale int
***REMOVED***

func (p *parser) Format(s fmt.State, verb rune) ***REMOVED***
	p.verb = verb
	p.scale = -1
	if prec, ok := s.Precision(); ok ***REMOVED***
		p.scale = prec
	***REMOVED***
***REMOVED***

type message struct ***REMOVED***
	arg   int
	kind  int
	scale int
	cases []interface***REMOVED******REMOVED***
***REMOVED***

const (
	// Start with non-ASCII to allow skipping values.
	kindDefault    = 0x80 + iota
	kindScale      // verb f, number of fraction digits follows
	kindScientific // verb e, number of fraction digits follows
	kindPrecision  // verb g, number of significant digits follows
)

var handle = catmsg.Register("golang.org/x/text/feature/plural:plural", execute)

func (m *message) Compile(e *catmsg.Encoder) error ***REMOVED***
	e.EncodeMessageType(handle)

	e.EncodeUint(uint64(m.arg))

	e.EncodeUint(uint64(m.kind))
	if m.kind > kindDefault ***REMOVED***
		e.EncodeUint(uint64(m.scale))
	***REMOVED***

	forms := validForms(cardinal, e.Language())

	for i := 0; i < len(m.cases); ***REMOVED***
		if err := compileSelector(e, forms, m.cases[i]); err != nil ***REMOVED***
			return err
		***REMOVED***
		if i++; i >= len(m.cases) ***REMOVED***
			return fmt.Errorf("plural: no message defined for selector %v", m.cases[i-1])
		***REMOVED***
		var msg catalog.Message
		switch x := m.cases[i].(type) ***REMOVED***
		case string:
			msg = catalog.String(x)
		case catalog.Message:
			msg = x
		default:
			return fmt.Errorf("plural: message of type %T; must be string or catalog.Message", x)
		***REMOVED***
		if err := e.EncodeMessage(msg); err != nil ***REMOVED***
			return err
		***REMOVED***
		i++
	***REMOVED***
	return nil
***REMOVED***

func compileSelector(e *catmsg.Encoder, valid []Form, selector interface***REMOVED******REMOVED***) error ***REMOVED***
	form := Other
	switch x := selector.(type) ***REMOVED***
	case string:
		if x == "" ***REMOVED***
			return fmt.Errorf("plural: empty selector")
		***REMOVED***
		if c := x[0]; c == '=' || c == '<' ***REMOVED***
			val, err := strconv.ParseUint(x[1:], 10, 16)
			if err != nil ***REMOVED***
				return fmt.Errorf("plural: invalid number in selector %q: %v", selector, err)
			***REMOVED***
			e.EncodeUint(uint64(c))
			e.EncodeUint(val)
			return nil
		***REMOVED***
		var ok bool
		form, ok = countMap[x]
		if !ok ***REMOVED***
			return fmt.Errorf("plural: invalid plural form %q", selector)
		***REMOVED***
	case Form:
		form = x
	default:
		return fmt.Errorf("plural: selector of type %T; want string or Form", selector)
	***REMOVED***

	ok := false
	for _, f := range valid ***REMOVED***
		if f == form ***REMOVED***
			ok = true
			break
		***REMOVED***
	***REMOVED***
	if !ok ***REMOVED***
		return fmt.Errorf("plural: form %q not supported for language %q", selector, e.Language())
	***REMOVED***
	e.EncodeUint(uint64(form))
	return nil
***REMOVED***

func execute(d *catmsg.Decoder) bool ***REMOVED***
	lang := d.Language()
	argN := int(d.DecodeUint())
	kind := int(d.DecodeUint())
	scale := -1 // default
	if kind > kindDefault ***REMOVED***
		scale = int(d.DecodeUint())
	***REMOVED***
	form := Other
	n := -1
	if arg := d.Arg(argN); arg == nil ***REMOVED***
		// Default to Other.
	***REMOVED*** else if x, ok := arg.(number.VisibleDigits); ok ***REMOVED***
		d := x.Digits(nil, lang, scale)
		form, n = cardinal.matchDisplayDigits(lang, &d)
	***REMOVED*** else if x, ok := arg.(Interface); ok ***REMOVED***
		// This covers lists and formatters from the number package.
		form, n = x.PluralForm(lang, scale)
	***REMOVED*** else ***REMOVED***
		var f number.Formatter
		switch kind ***REMOVED***
		case kindScale:
			f.InitDecimal(lang)
			f.SetScale(scale)
		case kindScientific:
			f.InitScientific(lang)
			f.SetScale(scale)
		case kindPrecision:
			f.InitDecimal(lang)
			f.SetPrecision(scale)
		case kindDefault:
			// sensible default
			f.InitDecimal(lang)
			if k := reflect.TypeOf(arg).Kind(); reflect.Int <= k && k <= reflect.Uintptr ***REMOVED***
				f.SetScale(0)
			***REMOVED*** else ***REMOVED***
				f.SetScale(2)
			***REMOVED***
		***REMOVED***
		var dec number.Decimal // TODO: buffer in Printer
		dec.Convert(f.RoundingContext, arg)
		v := number.FormatDigits(&dec, f.RoundingContext)
		if !v.NaN && !v.Inf ***REMOVED***
			form, n = cardinal.matchDisplayDigits(d.Language(), &v)
		***REMOVED***
	***REMOVED***
	for !d.Done() ***REMOVED***
		f := d.DecodeUint()
		if (f == '=' && n == int(d.DecodeUint())) ||
			(f == '<' && 0 <= n && n < int(d.DecodeUint())) ||
			form == Form(f) ||
			Other == Form(f) ***REMOVED***
			return d.ExecuteMessage()
		***REMOVED***
		d.SkipMessage()
	***REMOVED***
	return false
***REMOVED***
