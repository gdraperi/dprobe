// Copyright 2017 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package message

import (
	"bytes"
	"fmt" // TODO: consider copying interfaces from package fmt to avoid dependency.
	"math"
	"reflect"
	"sync"
	"unicode/utf8"

	"golang.org/x/text/internal/format"
	"golang.org/x/text/internal/number"
	"golang.org/x/text/language"
	"golang.org/x/text/message/catalog"
)

// Strings for use with buffer.WriteString.
// This is less overhead than using buffer.Write with byte arrays.
const (
	commaSpaceString  = ", "
	nilAngleString    = "<nil>"
	nilParenString    = "(nil)"
	nilString         = "nil"
	mapString         = "map["
	percentBangString = "%!"
	missingString     = "(MISSING)"
	badIndexString    = "(BADINDEX)"
	panicString       = "(PANIC="
	extraString       = "%!(EXTRA "
	badWidthString    = "%!(BADWIDTH)"
	badPrecString     = "%!(BADPREC)"
	noVerbString      = "%!(NOVERB)"

	invReflectString = "<invalid reflect.Value>"
)

var printerPool = sync.Pool***REMOVED***
	New: func() interface***REMOVED******REMOVED*** ***REMOVED*** return new(printer) ***REMOVED***,
***REMOVED***

// newPrinter allocates a new printer struct or grabs a cached one.
func newPrinter(pp *Printer) *printer ***REMOVED***
	p := printerPool.Get().(*printer)
	p.Printer = *pp
	// TODO: cache most of the following call.
	p.catContext = pp.cat.Context(pp.tag, p)

	p.panicking = false
	p.erroring = false
	p.fmt.init(&p.Buffer)
	return p
***REMOVED***

// free saves used printer structs in printerFree; avoids an allocation per invocation.
func (p *printer) free() ***REMOVED***
	p.Buffer.Reset()
	p.arg = nil
	p.value = reflect.Value***REMOVED******REMOVED***
	printerPool.Put(p)
***REMOVED***

// printer is used to store a printer's state.
// It implements "golang.org/x/text/internal/format".State.
type printer struct ***REMOVED***
	Printer

	// the context for looking up message translations
	catContext *catalog.Context

	// buffer for accumulating output.
	bytes.Buffer

	// arg holds the current item, as an interface***REMOVED******REMOVED***.
	arg interface***REMOVED******REMOVED***
	// value is used instead of arg for reflect values.
	value reflect.Value

	// fmt is used to format basic items such as integers or strings.
	fmt formatInfo

	// panicking is set by catchPanic to avoid infinite panic, recover, panic, ... recursion.
	panicking bool
	// erroring is set when printing an error string to guard against calling handleMethods.
	erroring bool
***REMOVED***

// Language implements "golang.org/x/text/internal/format".State.
func (p *printer) Language() language.Tag ***REMOVED*** return p.tag ***REMOVED***

func (p *printer) Width() (wid int, ok bool) ***REMOVED*** return p.fmt.Width, p.fmt.WidthPresent ***REMOVED***

func (p *printer) Precision() (prec int, ok bool) ***REMOVED*** return p.fmt.Prec, p.fmt.PrecPresent ***REMOVED***

func (p *printer) Flag(b int) bool ***REMOVED***
	switch b ***REMOVED***
	case '-':
		return p.fmt.Minus
	case '+':
		return p.fmt.Plus || p.fmt.PlusV
	case '#':
		return p.fmt.Sharp || p.fmt.SharpV
	case ' ':
		return p.fmt.Space
	case '0':
		return p.fmt.Zero
	***REMOVED***
	return false
***REMOVED***

// getField gets the i'th field of the struct value.
// If the field is itself is an interface, return a value for
// the thing inside the interface, not the interface itself.
func getField(v reflect.Value, i int) reflect.Value ***REMOVED***
	val := v.Field(i)
	if val.Kind() == reflect.Interface && !val.IsNil() ***REMOVED***
		val = val.Elem()
	***REMOVED***
	return val
***REMOVED***

func (p *printer) unknownType(v reflect.Value) ***REMOVED***
	if !v.IsValid() ***REMOVED***
		p.WriteString(nilAngleString)
		return
	***REMOVED***
	p.WriteByte('?')
	p.WriteString(v.Type().String())
	p.WriteByte('?')
***REMOVED***

func (p *printer) badVerb(verb rune) ***REMOVED***
	p.erroring = true
	p.WriteString(percentBangString)
	p.WriteRune(verb)
	p.WriteByte('(')
	switch ***REMOVED***
	case p.arg != nil:
		p.WriteString(reflect.TypeOf(p.arg).String())
		p.WriteByte('=')
		p.printArg(p.arg, 'v')
	case p.value.IsValid():
		p.WriteString(p.value.Type().String())
		p.WriteByte('=')
		p.printValue(p.value, 'v', 0)
	default:
		p.WriteString(nilAngleString)
	***REMOVED***
	p.WriteByte(')')
	p.erroring = false
***REMOVED***

func (p *printer) fmtBool(v bool, verb rune) ***REMOVED***
	switch verb ***REMOVED***
	case 't', 'v':
		p.fmt.fmt_boolean(v)
	default:
		p.badVerb(verb)
	***REMOVED***
***REMOVED***

// fmt0x64 formats a uint64 in hexadecimal and prefixes it with 0x or
// not, as requested, by temporarily setting the sharp flag.
func (p *printer) fmt0x64(v uint64, leading0x bool) ***REMOVED***
	sharp := p.fmt.Sharp
	p.fmt.Sharp = leading0x
	p.fmt.fmt_integer(v, 16, unsigned, ldigits)
	p.fmt.Sharp = sharp
***REMOVED***

// fmtInteger formats a signed or unsigned integer.
func (p *printer) fmtInteger(v uint64, isSigned bool, verb rune) ***REMOVED***
	switch verb ***REMOVED***
	case 'v':
		if p.fmt.SharpV && !isSigned ***REMOVED***
			p.fmt0x64(v, true)
			return
		***REMOVED***
		fallthrough
	case 'd':
		if p.fmt.Sharp || p.fmt.SharpV ***REMOVED***
			p.fmt.fmt_integer(v, 10, isSigned, ldigits)
		***REMOVED*** else ***REMOVED***
			p.fmtDecimalInt(v, isSigned)
		***REMOVED***
	case 'b':
		p.fmt.fmt_integer(v, 2, isSigned, ldigits)
	case 'o':
		p.fmt.fmt_integer(v, 8, isSigned, ldigits)
	case 'x':
		p.fmt.fmt_integer(v, 16, isSigned, ldigits)
	case 'X':
		p.fmt.fmt_integer(v, 16, isSigned, udigits)
	case 'c':
		p.fmt.fmt_c(v)
	case 'q':
		if v <= utf8.MaxRune ***REMOVED***
			p.fmt.fmt_qc(v)
		***REMOVED*** else ***REMOVED***
			p.badVerb(verb)
		***REMOVED***
	case 'U':
		p.fmt.fmt_unicode(v)
	default:
		p.badVerb(verb)
	***REMOVED***
***REMOVED***

// fmtFloat formats a float. The default precision for each verb
// is specified as last argument in the call to fmt_float.
func (p *printer) fmtFloat(v float64, size int, verb rune) ***REMOVED***
	switch verb ***REMOVED***
	case 'b':
		p.fmt.fmt_float(v, size, verb, -1)
	case 'v':
		verb = 'g'
		fallthrough
	case 'g', 'G':
		if p.fmt.Sharp || p.fmt.SharpV ***REMOVED***
			p.fmt.fmt_float(v, size, verb, -1)
		***REMOVED*** else ***REMOVED***
			p.fmtVariableFloat(v, size)
		***REMOVED***
	case 'e', 'E':
		if p.fmt.Sharp || p.fmt.SharpV ***REMOVED***
			p.fmt.fmt_float(v, size, verb, 6)
		***REMOVED*** else ***REMOVED***
			p.fmtScientific(v, size, 6)
		***REMOVED***
	case 'f', 'F':
		if p.fmt.Sharp || p.fmt.SharpV ***REMOVED***
			p.fmt.fmt_float(v, size, verb, 6)
		***REMOVED*** else ***REMOVED***
			p.fmtDecimalFloat(v, size, 6)
		***REMOVED***
	default:
		p.badVerb(verb)
	***REMOVED***
***REMOVED***

func (p *printer) setFlags(f *number.Formatter) ***REMOVED***
	f.Flags &^= number.ElideSign
	if p.fmt.Plus || p.fmt.Space ***REMOVED***
		f.Flags |= number.AlwaysSign
		if !p.fmt.Plus ***REMOVED***
			f.Flags |= number.ElideSign
		***REMOVED***
	***REMOVED*** else ***REMOVED***
		f.Flags &^= number.AlwaysSign
	***REMOVED***
***REMOVED***

func (p *printer) updatePadding(f *number.Formatter) ***REMOVED***
	f.Flags &^= number.PadMask
	if p.fmt.Minus ***REMOVED***
		f.Flags |= number.PadAfterSuffix
	***REMOVED*** else ***REMOVED***
		f.Flags |= number.PadBeforePrefix
	***REMOVED***
	f.PadRune = ' '
	f.FormatWidth = uint16(p.fmt.Width)
***REMOVED***

func (p *printer) initDecimal(minFrac, maxFrac int) ***REMOVED***
	f := &p.toDecimal
	f.MinIntegerDigits = 1
	f.MaxIntegerDigits = 0
	f.MinFractionDigits = uint8(minFrac)
	f.MaxFractionDigits = int16(maxFrac)
	p.setFlags(f)
	f.PadRune = 0
	if p.fmt.WidthPresent ***REMOVED***
		if p.fmt.Zero ***REMOVED***
			wid := p.fmt.Width
			// Use significant integers for this.
			// TODO: this is not the same as width, but so be it.
			if f.MinFractionDigits > 0 ***REMOVED***
				wid -= 1 + int(f.MinFractionDigits)
			***REMOVED***
			if p.fmt.Plus || p.fmt.Space ***REMOVED***
				wid--
			***REMOVED***
			if wid > 0 && wid > int(f.MinIntegerDigits) ***REMOVED***
				f.MinIntegerDigits = uint8(wid)
			***REMOVED***
		***REMOVED***
		p.updatePadding(f)
	***REMOVED***
***REMOVED***

func (p *printer) initScientific(minFrac, maxFrac int) ***REMOVED***
	f := &p.toScientific
	if maxFrac < 0 ***REMOVED***
		f.SetPrecision(maxFrac)
	***REMOVED*** else ***REMOVED***
		f.SetPrecision(maxFrac + 1)
		f.MinFractionDigits = uint8(minFrac)
		f.MaxFractionDigits = int16(maxFrac)
	***REMOVED***
	f.MinExponentDigits = 2
	p.setFlags(f)
	f.PadRune = 0
	if p.fmt.WidthPresent ***REMOVED***
		f.Flags &^= number.PadMask
		if p.fmt.Zero ***REMOVED***
			f.PadRune = f.Digit(0)
			f.Flags |= number.PadAfterPrefix
		***REMOVED*** else ***REMOVED***
			f.PadRune = ' '
			f.Flags |= number.PadBeforePrefix
		***REMOVED***
		p.updatePadding(f)
	***REMOVED***
***REMOVED***

func (p *printer) fmtDecimalInt(v uint64, isSigned bool) ***REMOVED***
	var d number.Decimal

	f := &p.toDecimal
	if p.fmt.PrecPresent ***REMOVED***
		p.setFlags(f)
		f.MinIntegerDigits = uint8(p.fmt.Prec)
		f.MaxIntegerDigits = 0
		f.MinFractionDigits = 0
		f.MaxFractionDigits = 0
		if p.fmt.WidthPresent ***REMOVED***
			p.updatePadding(f)
		***REMOVED***
	***REMOVED*** else ***REMOVED***
		p.initDecimal(0, 0)
	***REMOVED***
	d.ConvertInt(p.toDecimal.RoundingContext, isSigned, v)

	out := p.toDecimal.Format([]byte(nil), &d)
	p.Buffer.Write(out)
***REMOVED***

func (p *printer) fmtDecimalFloat(v float64, size, prec int) ***REMOVED***
	var d number.Decimal
	if p.fmt.PrecPresent ***REMOVED***
		prec = p.fmt.Prec
	***REMOVED***
	p.initDecimal(prec, prec)
	d.ConvertFloat(p.toDecimal.RoundingContext, v, size)

	out := p.toDecimal.Format([]byte(nil), &d)
	p.Buffer.Write(out)
***REMOVED***

func (p *printer) fmtVariableFloat(v float64, size int) ***REMOVED***
	prec := -1
	if p.fmt.PrecPresent ***REMOVED***
		prec = p.fmt.Prec
	***REMOVED***
	var d number.Decimal
	p.initScientific(0, prec)
	d.ConvertFloat(p.toScientific.RoundingContext, v, size)

	// Copy logic of 'g' formatting from strconv. It is simplified a bit as
	// we don't have to mind having prec > len(d.Digits).
	shortest := prec < 0
	ePrec := prec
	if shortest ***REMOVED***
		prec = len(d.Digits)
		ePrec = 6
	***REMOVED*** else if prec == 0 ***REMOVED***
		prec = 1
		ePrec = 1
	***REMOVED***
	exp := int(d.Exp) - 1
	if exp < -4 || exp >= ePrec ***REMOVED***
		p.initScientific(0, prec)

		out := p.toScientific.Format([]byte(nil), &d)
		p.Buffer.Write(out)
	***REMOVED*** else ***REMOVED***
		if prec > int(d.Exp) ***REMOVED***
			prec = len(d.Digits)
		***REMOVED***
		if prec -= int(d.Exp); prec < 0 ***REMOVED***
			prec = 0
		***REMOVED***
		p.initDecimal(0, prec)

		out := p.toDecimal.Format([]byte(nil), &d)
		p.Buffer.Write(out)
	***REMOVED***
***REMOVED***

func (p *printer) fmtScientific(v float64, size, prec int) ***REMOVED***
	var d number.Decimal
	if p.fmt.PrecPresent ***REMOVED***
		prec = p.fmt.Prec
	***REMOVED***
	p.initScientific(prec, prec)
	rc := p.toScientific.RoundingContext
	d.ConvertFloat(rc, v, size)

	out := p.toScientific.Format([]byte(nil), &d)
	p.Buffer.Write(out)

***REMOVED***

// fmtComplex formats a complex number v with
// r = real(v) and j = imag(v) as (r+ji) using
// fmtFloat for r and j formatting.
func (p *printer) fmtComplex(v complex128, size int, verb rune) ***REMOVED***
	// Make sure any unsupported verbs are found before the
	// calls to fmtFloat to not generate an incorrect error string.
	switch verb ***REMOVED***
	case 'v', 'b', 'g', 'G', 'f', 'F', 'e', 'E':
		p.WriteByte('(')
		p.fmtFloat(real(v), size/2, verb)
		// Imaginary part always has a sign.
		if math.IsNaN(imag(v)) ***REMOVED***
			// By CLDR's rules, NaNs do not use patterns or signs. As this code
			// relies on AlwaysSign working for imaginary parts, we need to
			// manually handle NaNs.
			f := &p.toScientific
			p.setFlags(f)
			p.updatePadding(f)
			p.setFlags(f)
			nan := f.Symbol(number.SymNan)
			extra := 0
			if w, ok := p.Width(); ok ***REMOVED***
				extra = w - utf8.RuneCountInString(nan) - 1
			***REMOVED***
			if f.Flags&number.PadAfterNumber == 0 ***REMOVED***
				for ; extra > 0; extra-- ***REMOVED***
					p.WriteRune(f.PadRune)
				***REMOVED***
			***REMOVED***
			p.WriteString(f.Symbol(number.SymPlusSign))
			p.WriteString(nan)
			for ; extra > 0; extra-- ***REMOVED***
				p.WriteRune(f.PadRune)
			***REMOVED***
			p.WriteString("i)")
			return
		***REMOVED***
		oldPlus := p.fmt.Plus
		p.fmt.Plus = true
		p.fmtFloat(imag(v), size/2, verb)
		p.WriteString("i)") // TODO: use symbol?
		p.fmt.Plus = oldPlus
	default:
		p.badVerb(verb)
	***REMOVED***
***REMOVED***

func (p *printer) fmtString(v string, verb rune) ***REMOVED***
	switch verb ***REMOVED***
	case 'v':
		if p.fmt.SharpV ***REMOVED***
			p.fmt.fmt_q(v)
		***REMOVED*** else ***REMOVED***
			p.fmt.fmt_s(v)
		***REMOVED***
	case 's':
		p.fmt.fmt_s(v)
	case 'x':
		p.fmt.fmt_sx(v, ldigits)
	case 'X':
		p.fmt.fmt_sx(v, udigits)
	case 'q':
		p.fmt.fmt_q(v)
	default:
		p.badVerb(verb)
	***REMOVED***
***REMOVED***

func (p *printer) fmtBytes(v []byte, verb rune, typeString string) ***REMOVED***
	switch verb ***REMOVED***
	case 'v', 'd':
		if p.fmt.SharpV ***REMOVED***
			p.WriteString(typeString)
			if v == nil ***REMOVED***
				p.WriteString(nilParenString)
				return
			***REMOVED***
			p.WriteByte('***REMOVED***')
			for i, c := range v ***REMOVED***
				if i > 0 ***REMOVED***
					p.WriteString(commaSpaceString)
				***REMOVED***
				p.fmt0x64(uint64(c), true)
			***REMOVED***
			p.WriteByte('***REMOVED***')
		***REMOVED*** else ***REMOVED***
			p.WriteByte('[')
			for i, c := range v ***REMOVED***
				if i > 0 ***REMOVED***
					p.WriteByte(' ')
				***REMOVED***
				p.fmt.fmt_integer(uint64(c), 10, unsigned, ldigits)
			***REMOVED***
			p.WriteByte(']')
		***REMOVED***
	case 's':
		p.fmt.fmt_s(string(v))
	case 'x':
		p.fmt.fmt_bx(v, ldigits)
	case 'X':
		p.fmt.fmt_bx(v, udigits)
	case 'q':
		p.fmt.fmt_q(string(v))
	default:
		p.printValue(reflect.ValueOf(v), verb, 0)
	***REMOVED***
***REMOVED***

func (p *printer) fmtPointer(value reflect.Value, verb rune) ***REMOVED***
	var u uintptr
	switch value.Kind() ***REMOVED***
	case reflect.Chan, reflect.Func, reflect.Map, reflect.Ptr, reflect.Slice, reflect.UnsafePointer:
		u = value.Pointer()
	default:
		p.badVerb(verb)
		return
	***REMOVED***

	switch verb ***REMOVED***
	case 'v':
		if p.fmt.SharpV ***REMOVED***
			p.WriteByte('(')
			p.WriteString(value.Type().String())
			p.WriteString(")(")
			if u == 0 ***REMOVED***
				p.WriteString(nilString)
			***REMOVED*** else ***REMOVED***
				p.fmt0x64(uint64(u), true)
			***REMOVED***
			p.WriteByte(')')
		***REMOVED*** else ***REMOVED***
			if u == 0 ***REMOVED***
				p.fmt.padString(nilAngleString)
			***REMOVED*** else ***REMOVED***
				p.fmt0x64(uint64(u), !p.fmt.Sharp)
			***REMOVED***
		***REMOVED***
	case 'p':
		p.fmt0x64(uint64(u), !p.fmt.Sharp)
	case 'b', 'o', 'd', 'x', 'X':
		if verb == 'd' ***REMOVED***
			p.fmt.Sharp = true // Print as standard go. TODO: does this make sense?
		***REMOVED***
		p.fmtInteger(uint64(u), unsigned, verb)
	default:
		p.badVerb(verb)
	***REMOVED***
***REMOVED***

func (p *printer) catchPanic(arg interface***REMOVED******REMOVED***, verb rune) ***REMOVED***
	if err := recover(); err != nil ***REMOVED***
		// If it's a nil pointer, just say "<nil>". The likeliest causes are a
		// Stringer that fails to guard against nil or a nil pointer for a
		// value receiver, and in either case, "<nil>" is a nice result.
		if v := reflect.ValueOf(arg); v.Kind() == reflect.Ptr && v.IsNil() ***REMOVED***
			p.WriteString(nilAngleString)
			return
		***REMOVED***
		// Otherwise print a concise panic message. Most of the time the panic
		// value will print itself nicely.
		if p.panicking ***REMOVED***
			// Nested panics; the recursion in printArg cannot succeed.
			panic(err)
		***REMOVED***

		oldFlags := p.fmt.Parser
		// For this output we want default behavior.
		p.fmt.ClearFlags()

		p.WriteString(percentBangString)
		p.WriteRune(verb)
		p.WriteString(panicString)
		p.panicking = true
		p.printArg(err, 'v')
		p.panicking = false
		p.WriteByte(')')

		p.fmt.Parser = oldFlags
	***REMOVED***
***REMOVED***

func (p *printer) handleMethods(verb rune) (handled bool) ***REMOVED***
	if p.erroring ***REMOVED***
		return
	***REMOVED***
	// Is it a Formatter?
	if formatter, ok := p.arg.(format.Formatter); ok ***REMOVED***
		handled = true
		defer p.catchPanic(p.arg, verb)
		formatter.Format(p, verb)
		return
	***REMOVED***
	if formatter, ok := p.arg.(fmt.Formatter); ok ***REMOVED***
		handled = true
		defer p.catchPanic(p.arg, verb)
		formatter.Format(p, verb)
		return
	***REMOVED***

	// If we're doing Go syntax and the argument knows how to supply it, take care of it now.
	if p.fmt.SharpV ***REMOVED***
		if stringer, ok := p.arg.(fmt.GoStringer); ok ***REMOVED***
			handled = true
			defer p.catchPanic(p.arg, verb)
			// Print the result of GoString unadorned.
			p.fmt.fmt_s(stringer.GoString())
			return
		***REMOVED***
	***REMOVED*** else ***REMOVED***
		// If a string is acceptable according to the format, see if
		// the value satisfies one of the string-valued interfaces.
		// Println etc. set verb to %v, which is "stringable".
		switch verb ***REMOVED***
		case 'v', 's', 'x', 'X', 'q':
			// Is it an error or Stringer?
			// The duplication in the bodies is necessary:
			// setting handled and deferring catchPanic
			// must happen before calling the method.
			switch v := p.arg.(type) ***REMOVED***
			case error:
				handled = true
				defer p.catchPanic(p.arg, verb)
				p.fmtString(v.Error(), verb)
				return

			case fmt.Stringer:
				handled = true
				defer p.catchPanic(p.arg, verb)
				p.fmtString(v.String(), verb)
				return
			***REMOVED***
		***REMOVED***
	***REMOVED***
	return false
***REMOVED***

func (p *printer) printArg(arg interface***REMOVED******REMOVED***, verb rune) ***REMOVED***
	p.arg = arg
	p.value = reflect.Value***REMOVED******REMOVED***

	if arg == nil ***REMOVED***
		switch verb ***REMOVED***
		case 'T', 'v':
			p.fmt.padString(nilAngleString)
		default:
			p.badVerb(verb)
		***REMOVED***
		return
	***REMOVED***

	// Special processing considerations.
	// %T (the value's type) and %p (its address) are special; we always do them first.
	switch verb ***REMOVED***
	case 'T':
		p.fmt.fmt_s(reflect.TypeOf(arg).String())
		return
	case 'p':
		p.fmtPointer(reflect.ValueOf(arg), 'p')
		return
	***REMOVED***

	// Some types can be done without reflection.
	switch f := arg.(type) ***REMOVED***
	case bool:
		p.fmtBool(f, verb)
	case float32:
		p.fmtFloat(float64(f), 32, verb)
	case float64:
		p.fmtFloat(f, 64, verb)
	case complex64:
		p.fmtComplex(complex128(f), 64, verb)
	case complex128:
		p.fmtComplex(f, 128, verb)
	case int:
		p.fmtInteger(uint64(f), signed, verb)
	case int8:
		p.fmtInteger(uint64(f), signed, verb)
	case int16:
		p.fmtInteger(uint64(f), signed, verb)
	case int32:
		p.fmtInteger(uint64(f), signed, verb)
	case int64:
		p.fmtInteger(uint64(f), signed, verb)
	case uint:
		p.fmtInteger(uint64(f), unsigned, verb)
	case uint8:
		p.fmtInteger(uint64(f), unsigned, verb)
	case uint16:
		p.fmtInteger(uint64(f), unsigned, verb)
	case uint32:
		p.fmtInteger(uint64(f), unsigned, verb)
	case uint64:
		p.fmtInteger(f, unsigned, verb)
	case uintptr:
		p.fmtInteger(uint64(f), unsigned, verb)
	case string:
		p.fmtString(f, verb)
	case []byte:
		p.fmtBytes(f, verb, "[]byte")
	case reflect.Value:
		// Handle extractable values with special methods
		// since printValue does not handle them at depth 0.
		if f.IsValid() && f.CanInterface() ***REMOVED***
			p.arg = f.Interface()
			if p.handleMethods(verb) ***REMOVED***
				return
			***REMOVED***
		***REMOVED***
		p.printValue(f, verb, 0)
	default:
		// If the type is not simple, it might have methods.
		if !p.handleMethods(verb) ***REMOVED***
			// Need to use reflection, since the type had no
			// interface methods that could be used for formatting.
			p.printValue(reflect.ValueOf(f), verb, 0)
		***REMOVED***
	***REMOVED***
***REMOVED***

// printValue is similar to printArg but starts with a reflect value, not an interface***REMOVED******REMOVED*** value.
// It does not handle 'p' and 'T' verbs because these should have been already handled by printArg.
func (p *printer) printValue(value reflect.Value, verb rune, depth int) ***REMOVED***
	// Handle values with special methods if not already handled by printArg (depth == 0).
	if depth > 0 && value.IsValid() && value.CanInterface() ***REMOVED***
		p.arg = value.Interface()
		if p.handleMethods(verb) ***REMOVED***
			return
		***REMOVED***
	***REMOVED***
	p.arg = nil
	p.value = value

	switch f := value; value.Kind() ***REMOVED***
	case reflect.Invalid:
		if depth == 0 ***REMOVED***
			p.WriteString(invReflectString)
		***REMOVED*** else ***REMOVED***
			switch verb ***REMOVED***
			case 'v':
				p.WriteString(nilAngleString)
			default:
				p.badVerb(verb)
			***REMOVED***
		***REMOVED***
	case reflect.Bool:
		p.fmtBool(f.Bool(), verb)
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		p.fmtInteger(uint64(f.Int()), signed, verb)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		p.fmtInteger(f.Uint(), unsigned, verb)
	case reflect.Float32:
		p.fmtFloat(f.Float(), 32, verb)
	case reflect.Float64:
		p.fmtFloat(f.Float(), 64, verb)
	case reflect.Complex64:
		p.fmtComplex(f.Complex(), 64, verb)
	case reflect.Complex128:
		p.fmtComplex(f.Complex(), 128, verb)
	case reflect.String:
		p.fmtString(f.String(), verb)
	case reflect.Map:
		if p.fmt.SharpV ***REMOVED***
			p.WriteString(f.Type().String())
			if f.IsNil() ***REMOVED***
				p.WriteString(nilParenString)
				return
			***REMOVED***
			p.WriteByte('***REMOVED***')
		***REMOVED*** else ***REMOVED***
			p.WriteString(mapString)
		***REMOVED***
		keys := f.MapKeys()
		for i, key := range keys ***REMOVED***
			if i > 0 ***REMOVED***
				if p.fmt.SharpV ***REMOVED***
					p.WriteString(commaSpaceString)
				***REMOVED*** else ***REMOVED***
					p.WriteByte(' ')
				***REMOVED***
			***REMOVED***
			p.printValue(key, verb, depth+1)
			p.WriteByte(':')
			p.printValue(f.MapIndex(key), verb, depth+1)
		***REMOVED***
		if p.fmt.SharpV ***REMOVED***
			p.WriteByte('***REMOVED***')
		***REMOVED*** else ***REMOVED***
			p.WriteByte(']')
		***REMOVED***
	case reflect.Struct:
		if p.fmt.SharpV ***REMOVED***
			p.WriteString(f.Type().String())
		***REMOVED***
		p.WriteByte('***REMOVED***')
		for i := 0; i < f.NumField(); i++ ***REMOVED***
			if i > 0 ***REMOVED***
				if p.fmt.SharpV ***REMOVED***
					p.WriteString(commaSpaceString)
				***REMOVED*** else ***REMOVED***
					p.WriteByte(' ')
				***REMOVED***
			***REMOVED***
			if p.fmt.PlusV || p.fmt.SharpV ***REMOVED***
				if name := f.Type().Field(i).Name; name != "" ***REMOVED***
					p.WriteString(name)
					p.WriteByte(':')
				***REMOVED***
			***REMOVED***
			p.printValue(getField(f, i), verb, depth+1)
		***REMOVED***
		p.WriteByte('***REMOVED***')
	case reflect.Interface:
		value := f.Elem()
		if !value.IsValid() ***REMOVED***
			if p.fmt.SharpV ***REMOVED***
				p.WriteString(f.Type().String())
				p.WriteString(nilParenString)
			***REMOVED*** else ***REMOVED***
				p.WriteString(nilAngleString)
			***REMOVED***
		***REMOVED*** else ***REMOVED***
			p.printValue(value, verb, depth+1)
		***REMOVED***
	case reflect.Array, reflect.Slice:
		switch verb ***REMOVED***
		case 's', 'q', 'x', 'X':
			// Handle byte and uint8 slices and arrays special for the above verbs.
			t := f.Type()
			if t.Elem().Kind() == reflect.Uint8 ***REMOVED***
				var bytes []byte
				if f.Kind() == reflect.Slice ***REMOVED***
					bytes = f.Bytes()
				***REMOVED*** else if f.CanAddr() ***REMOVED***
					bytes = f.Slice(0, f.Len()).Bytes()
				***REMOVED*** else ***REMOVED***
					// We have an array, but we cannot Slice() a non-addressable array,
					// so we build a slice by hand. This is a rare case but it would be nice
					// if reflection could help a little more.
					bytes = make([]byte, f.Len())
					for i := range bytes ***REMOVED***
						bytes[i] = byte(f.Index(i).Uint())
					***REMOVED***
				***REMOVED***
				p.fmtBytes(bytes, verb, t.String())
				return
			***REMOVED***
		***REMOVED***
		if p.fmt.SharpV ***REMOVED***
			p.WriteString(f.Type().String())
			if f.Kind() == reflect.Slice && f.IsNil() ***REMOVED***
				p.WriteString(nilParenString)
				return
			***REMOVED***
			p.WriteByte('***REMOVED***')
			for i := 0; i < f.Len(); i++ ***REMOVED***
				if i > 0 ***REMOVED***
					p.WriteString(commaSpaceString)
				***REMOVED***
				p.printValue(f.Index(i), verb, depth+1)
			***REMOVED***
			p.WriteByte('***REMOVED***')
		***REMOVED*** else ***REMOVED***
			p.WriteByte('[')
			for i := 0; i < f.Len(); i++ ***REMOVED***
				if i > 0 ***REMOVED***
					p.WriteByte(' ')
				***REMOVED***
				p.printValue(f.Index(i), verb, depth+1)
			***REMOVED***
			p.WriteByte(']')
		***REMOVED***
	case reflect.Ptr:
		// pointer to array or slice or struct?  ok at top level
		// but not embedded (avoid loops)
		if depth == 0 && f.Pointer() != 0 ***REMOVED***
			switch a := f.Elem(); a.Kind() ***REMOVED***
			case reflect.Array, reflect.Slice, reflect.Struct, reflect.Map:
				p.WriteByte('&')
				p.printValue(a, verb, depth+1)
				return
			***REMOVED***
		***REMOVED***
		fallthrough
	case reflect.Chan, reflect.Func, reflect.UnsafePointer:
		p.fmtPointer(f, verb)
	default:
		p.unknownType(f)
	***REMOVED***
***REMOVED***

func (p *printer) badArgNum(verb rune) ***REMOVED***
	p.WriteString(percentBangString)
	p.WriteRune(verb)
	p.WriteString(badIndexString)
***REMOVED***

func (p *printer) missingArg(verb rune) ***REMOVED***
	p.WriteString(percentBangString)
	p.WriteRune(verb)
	p.WriteString(missingString)
***REMOVED***

func (p *printer) doPrintf(fmt string) ***REMOVED***
	for p.fmt.Parser.SetFormat(fmt); p.fmt.Scan(); ***REMOVED***
		switch p.fmt.Status ***REMOVED***
		case format.StatusText:
			p.WriteString(p.fmt.Text())
		case format.StatusSubstitution:
			p.printArg(p.Arg(p.fmt.ArgNum), p.fmt.Verb)
		case format.StatusBadWidthSubstitution:
			p.WriteString(badWidthString)
			p.printArg(p.Arg(p.fmt.ArgNum), p.fmt.Verb)
		case format.StatusBadPrecSubstitution:
			p.WriteString(badPrecString)
			p.printArg(p.Arg(p.fmt.ArgNum), p.fmt.Verb)
		case format.StatusNoVerb:
			p.WriteString(noVerbString)
		case format.StatusBadArgNum:
			p.badArgNum(p.fmt.Verb)
		case format.StatusMissingArg:
			p.missingArg(p.fmt.Verb)
		default:
			panic("unreachable")
		***REMOVED***
	***REMOVED***

	// Check for extra arguments, but only if there was at least one ordered
	// argument. Note that this behavior is necessarily different from fmt:
	// different variants of messages may opt to drop some or all of the
	// arguments.
	if !p.fmt.Reordered && p.fmt.ArgNum < len(p.fmt.Args) && p.fmt.ArgNum != 0 ***REMOVED***
		p.fmt.ClearFlags()
		p.WriteString(extraString)
		for i, arg := range p.fmt.Args[p.fmt.ArgNum:] ***REMOVED***
			if i > 0 ***REMOVED***
				p.WriteString(commaSpaceString)
			***REMOVED***
			if arg == nil ***REMOVED***
				p.WriteString(nilAngleString)
			***REMOVED*** else ***REMOVED***
				p.WriteString(reflect.TypeOf(arg).String())
				p.WriteString("=")
				p.printArg(arg, 'v')
			***REMOVED***
		***REMOVED***
		p.WriteByte(')')
	***REMOVED***
***REMOVED***

func (p *printer) doPrint(a []interface***REMOVED******REMOVED***) ***REMOVED***
	prevString := false
	for argNum, arg := range a ***REMOVED***
		isString := arg != nil && reflect.TypeOf(arg).Kind() == reflect.String
		// Add a space between two non-string arguments.
		if argNum > 0 && !isString && !prevString ***REMOVED***
			p.WriteByte(' ')
		***REMOVED***
		p.printArg(arg, 'v')
		prevString = isString
	***REMOVED***
***REMOVED***

// doPrintln is like doPrint but always adds a space between arguments
// and a newline after the last argument.
func (p *printer) doPrintln(a []interface***REMOVED******REMOVED***) ***REMOVED***
	for argNum, arg := range a ***REMOVED***
		if argNum > 0 ***REMOVED***
			p.WriteByte(' ')
		***REMOVED***
		p.printArg(arg, 'v')
	***REMOVED***
	p.WriteByte('\n')
***REMOVED***
