// Copyright 2015 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package message // import "golang.org/x/text/message"

import (
	"io"
	"os"

	// Include features to facilitate generated catalogs.
	_ "golang.org/x/text/feature/plural"

	"golang.org/x/text/internal/number"
	"golang.org/x/text/language"
	"golang.org/x/text/message/catalog"
)

// A Printer implements language-specific formatted I/O analogous to the fmt
// package.
type Printer struct ***REMOVED***
	// the language
	tag language.Tag

	toDecimal    number.Formatter
	toScientific number.Formatter

	cat catalog.Catalog
***REMOVED***

type options struct ***REMOVED***
	cat catalog.Catalog
	// TODO:
	// - allow %s to print integers in written form (tables are likely too large
	//   to enable this by default).
	// - list behavior
	//
***REMOVED***

// An Option defines an option of a Printer.
type Option func(o *options)

// Catalog defines the catalog to be used.
func Catalog(c catalog.Catalog) Option ***REMOVED***
	return func(o *options) ***REMOVED*** o.cat = c ***REMOVED***
***REMOVED***

// NewPrinter returns a Printer that formats messages tailored to language t.
func NewPrinter(t language.Tag, opts ...Option) *Printer ***REMOVED***
	options := &options***REMOVED***
		cat: DefaultCatalog,
	***REMOVED***
	for _, o := range opts ***REMOVED***
		o(options)
	***REMOVED***
	p := &Printer***REMOVED***
		tag: t,
		cat: options.cat,
	***REMOVED***
	p.toDecimal.InitDecimal(t)
	p.toScientific.InitScientific(t)
	return p
***REMOVED***

// Sprint is like fmt.Sprint, but using language-specific formatting.
func (p *Printer) Sprint(a ...interface***REMOVED******REMOVED***) string ***REMOVED***
	pp := newPrinter(p)
	pp.doPrint(a)
	s := pp.String()
	pp.free()
	return s
***REMOVED***

// Fprint is like fmt.Fprint, but using language-specific formatting.
func (p *Printer) Fprint(w io.Writer, a ...interface***REMOVED******REMOVED***) (n int, err error) ***REMOVED***
	pp := newPrinter(p)
	pp.doPrint(a)
	n64, err := io.Copy(w, &pp.Buffer)
	pp.free()
	return int(n64), err
***REMOVED***

// Print is like fmt.Print, but using language-specific formatting.
func (p *Printer) Print(a ...interface***REMOVED******REMOVED***) (n int, err error) ***REMOVED***
	return p.Fprint(os.Stdout, a...)
***REMOVED***

// Sprintln is like fmt.Sprintln, but using language-specific formatting.
func (p *Printer) Sprintln(a ...interface***REMOVED******REMOVED***) string ***REMOVED***
	pp := newPrinter(p)
	pp.doPrintln(a)
	s := pp.String()
	pp.free()
	return s
***REMOVED***

// Fprintln is like fmt.Fprintln, but using language-specific formatting.
func (p *Printer) Fprintln(w io.Writer, a ...interface***REMOVED******REMOVED***) (n int, err error) ***REMOVED***
	pp := newPrinter(p)
	pp.doPrintln(a)
	n64, err := io.Copy(w, &pp.Buffer)
	pp.free()
	return int(n64), err
***REMOVED***

// Println is like fmt.Println, but using language-specific formatting.
func (p *Printer) Println(a ...interface***REMOVED******REMOVED***) (n int, err error) ***REMOVED***
	return p.Fprintln(os.Stdout, a...)
***REMOVED***

// Sprintf is like fmt.Sprintf, but using language-specific formatting.
func (p *Printer) Sprintf(key Reference, a ...interface***REMOVED******REMOVED***) string ***REMOVED***
	pp := newPrinter(p)
	lookupAndFormat(pp, key, a)
	s := pp.String()
	pp.free()
	return s
***REMOVED***

// Fprintf is like fmt.Fprintf, but using language-specific formatting.
func (p *Printer) Fprintf(w io.Writer, key Reference, a ...interface***REMOVED******REMOVED***) (n int, err error) ***REMOVED***
	pp := newPrinter(p)
	lookupAndFormat(pp, key, a)
	n, err = w.Write(pp.Bytes())
	pp.free()
	return n, err

***REMOVED***

// Printf is like fmt.Printf, but using language-specific formatting.
func (p *Printer) Printf(key Reference, a ...interface***REMOVED******REMOVED***) (n int, err error) ***REMOVED***
	pp := newPrinter(p)
	lookupAndFormat(pp, key, a)
	n, err = os.Stdout.Write(pp.Bytes())
	pp.free()
	return n, err
***REMOVED***

func lookupAndFormat(p *printer, r Reference, a []interface***REMOVED******REMOVED***) ***REMOVED***
	p.fmt.Reset(a)
	var id, msg string
	switch v := r.(type) ***REMOVED***
	case string:
		id, msg = v, v
	case key:
		id, msg = v.id, v.fallback
	default:
		panic("key argument is not a Reference")
	***REMOVED***

	if p.catContext.Execute(id) == catalog.ErrNotFound ***REMOVED***
		if p.catContext.Execute(msg) == catalog.ErrNotFound ***REMOVED***
			p.Render(msg)
			return
		***REMOVED***
	***REMOVED***
***REMOVED***

// Arg implements catmsg.Renderer.
func (p *printer) Arg(i int) interface***REMOVED******REMOVED*** ***REMOVED*** // TODO, also return "ok" bool
	i--
	if uint(i) < uint(len(p.fmt.Args)) ***REMOVED***
		return p.fmt.Args[i]
	***REMOVED***
	return nil
***REMOVED***

// Render implements catmsg.Renderer.
func (p *printer) Render(msg string) ***REMOVED***
	p.doPrintf(msg)
***REMOVED***

// A Reference is a string or a message reference.
type Reference interface ***REMOVED***
	// TODO: also allow []string
***REMOVED***

// Key creates a message Reference for a message where the given id is used for
// message lookup and the fallback is returned when no matches are found.
func Key(id string, fallback string) Reference ***REMOVED***
	return key***REMOVED***id, fallback***REMOVED***
***REMOVED***

type key struct ***REMOVED***
	id, fallback string
***REMOVED***
