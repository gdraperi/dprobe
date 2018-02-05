// Copyright 2017 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package catmsg

import (
	"errors"
	"fmt"

	"golang.org/x/text/language"
)

// A Renderer renders a Message.
type Renderer interface ***REMOVED***
	// Render renders the given string. The given string may be interpreted as a
	// format string, such as the one used by the fmt package or a template.
	Render(s string)

	// Arg returns the i-th argument passed to format a message. This method
	// should return nil if there is no such argument. Messages need access to
	// arguments to allow selecting a message based on linguistic features of
	// those arguments.
	Arg(i int) interface***REMOVED******REMOVED***
***REMOVED***

// A Dictionary specifies a source of messages, including variables or macros.
type Dictionary interface ***REMOVED***
	// Lookup returns the message for the given key. It returns false for ok if
	// such a message could not be found.
	Lookup(key string) (data string, ok bool)

	// TODO: consider returning an interface, instead of a string. This will
	// allow implementations to do their own message type decoding.
***REMOVED***

// An Encoder serializes a Message to a string.
type Encoder struct ***REMOVED***
	// The root encoder is used for storing encoded variables.
	root *Encoder
	// The parent encoder provides the surrounding scopes for resolving variable
	// names.
	parent *Encoder

	tag language.Tag

	// buf holds the encoded message so far. After a message completes encoding,
	// the contents of buf, prefixed by the encoded length, are flushed to the
	// parent buffer.
	buf []byte

	// vars is the lookup table of variables in the current scope.
	vars []keyVal

	err    error
	inBody bool // if false next call must be EncodeMessageType
***REMOVED***

type keyVal struct ***REMOVED***
	key    string
	offset int
***REMOVED***

// Language reports the language for which the encoded message will be stored
// in the Catalog.
func (e *Encoder) Language() language.Tag ***REMOVED*** return e.tag ***REMOVED***

func (e *Encoder) setError(err error) ***REMOVED***
	if e.root.err == nil ***REMOVED***
		e.root.err = err
	***REMOVED***
***REMOVED***

// EncodeUint encodes x.
func (e *Encoder) EncodeUint(x uint64) ***REMOVED***
	e.checkInBody()
	var buf [maxVarintBytes]byte
	n := encodeUint(buf[:], x)
	e.buf = append(e.buf, buf[:n]...)
***REMOVED***

// EncodeString encodes s.
func (e *Encoder) EncodeString(s string) ***REMOVED***
	e.checkInBody()
	e.EncodeUint(uint64(len(s)))
	e.buf = append(e.buf, s...)
***REMOVED***

// EncodeMessageType marks the current message to be of type h.
//
// It must be the first call of a Message's Compile method.
func (e *Encoder) EncodeMessageType(h Handle) ***REMOVED***
	if e.inBody ***REMOVED***
		panic("catmsg: EncodeMessageType not the first method called")
	***REMOVED***
	e.inBody = true
	e.EncodeUint(uint64(h))
***REMOVED***

// EncodeMessage serializes the given message inline at the current position.
func (e *Encoder) EncodeMessage(m Message) error ***REMOVED***
	e = &Encoder***REMOVED***root: e.root, parent: e, tag: e.tag***REMOVED***
	err := m.Compile(e)
	if _, ok := m.(*Var); !ok ***REMOVED***
		e.flushTo(e.parent)
	***REMOVED***
	return err
***REMOVED***

func (e *Encoder) checkInBody() ***REMOVED***
	if !e.inBody ***REMOVED***
		panic("catmsg: expected prior call to EncodeMessageType")
	***REMOVED***
***REMOVED***

// stripPrefix indicates the number of prefix bytes that must be stripped to
// turn a single-element sequence into a message that is just this single member
// without its size prefix. If the message can be stripped, b[1:n] contains the
// size prefix.
func stripPrefix(b []byte) (n int) ***REMOVED***
	if len(b) > 0 && Handle(b[0]) == msgFirst ***REMOVED***
		x, n, _ := decodeUint(b[1:])
		if 1+n+int(x) == len(b) ***REMOVED***
			return 1 + n
		***REMOVED***
	***REMOVED***
	return 0
***REMOVED***

func (e *Encoder) flushTo(dst *Encoder) ***REMOVED***
	data := e.buf
	p := stripPrefix(data)
	if p > 0 ***REMOVED***
		data = data[1:]
	***REMOVED*** else ***REMOVED***
		// Prefix the size.
		dst.EncodeUint(uint64(len(data)))
	***REMOVED***
	dst.buf = append(dst.buf, data...)
***REMOVED***

func (e *Encoder) addVar(key string, m Message) error ***REMOVED***
	for _, v := range e.parent.vars ***REMOVED***
		if v.key == key ***REMOVED***
			err := fmt.Errorf("catmsg: duplicate variable %q", key)
			e.setError(err)
			return err
		***REMOVED***
	***REMOVED***
	scope := e.parent
	// If a variable message is Incomplete, and does not evaluate to a message
	// during execution, we fall back to the variable name. We encode this by
	// appending the variable name if the message reports it's incomplete.

	err := m.Compile(e)
	if err != ErrIncomplete ***REMOVED***
		e.setError(err)
	***REMOVED***
	switch ***REMOVED***
	case len(e.buf) == 1 && Handle(e.buf[0]) == msgFirst: // empty sequence
		e.buf = e.buf[:0]
		e.inBody = false
		fallthrough
	case len(e.buf) == 0:
		// Empty message.
		if err := String(key).Compile(e); err != nil ***REMOVED***
			e.setError(err)
		***REMOVED***
	case err == ErrIncomplete:
		if Handle(e.buf[0]) != msgFirst ***REMOVED***
			seq := &Encoder***REMOVED***root: e.root, parent: e***REMOVED***
			seq.EncodeMessageType(msgFirst)
			e.flushTo(seq)
			e = seq
		***REMOVED***
		// e contains a sequence; append the fallback string.
		e.EncodeMessage(String(key))
	***REMOVED***

	// Flush result to variable heap.
	offset := len(e.root.buf)
	e.flushTo(e.root)
	e.buf = e.buf[:0]

	// Record variable offset in current scope.
	scope.vars = append(scope.vars, keyVal***REMOVED***key: key, offset: offset***REMOVED***)
	return err
***REMOVED***

const (
	substituteVar = iota
	substituteMacro
	substituteError
)

// EncodeSubstitution inserts a resolved reference to a variable or macro.
//
// This call must be matched with a call to ExecuteSubstitution at decoding
// time.
func (e *Encoder) EncodeSubstitution(name string, arguments ...int) ***REMOVED***
	if arity := len(arguments); arity > 0 ***REMOVED***
		// TODO: also resolve macros.
		e.EncodeUint(substituteMacro)
		e.EncodeString(name)
		for _, a := range arguments ***REMOVED***
			e.EncodeUint(uint64(a))
		***REMOVED***
		return
	***REMOVED***
	for scope := e; scope != nil; scope = scope.parent ***REMOVED***
		for _, v := range scope.vars ***REMOVED***
			if v.key != name ***REMOVED***
				continue
			***REMOVED***
			e.EncodeUint(substituteVar) // TODO: support arity > 0
			e.EncodeUint(uint64(v.offset))
			return
		***REMOVED***
	***REMOVED***
	// TODO: refer to dictionary-wide scoped variables.
	e.EncodeUint(substituteError)
	e.EncodeString(name)
	e.setError(fmt.Errorf("catmsg: unknown var %q", name))
***REMOVED***

// A Decoder deserializes and evaluates messages that are encoded by an encoder.
type Decoder struct ***REMOVED***
	tag    language.Tag
	dst    Renderer
	macros Dictionary

	err  error
	vars string
	data string

	macroArg int // TODO: allow more than one argument
***REMOVED***

// NewDecoder returns a new Decoder.
//
// Decoders are designed to be reused for multiple invocations of Execute.
// Only one goroutine may call Execute concurrently.
func NewDecoder(tag language.Tag, r Renderer, macros Dictionary) *Decoder ***REMOVED***
	return &Decoder***REMOVED***
		tag:    tag,
		dst:    r,
		macros: macros,
	***REMOVED***
***REMOVED***

func (d *Decoder) setError(err error) ***REMOVED***
	if d.err == nil ***REMOVED***
		d.err = err
	***REMOVED***
***REMOVED***

// Language returns the language in which the message is being rendered.
//
// The destination language may be a child language of the language used for
// encoding. For instance, a decoding language of "pt-PT"" is consistent with an
// encoding language of "pt".
func (d *Decoder) Language() language.Tag ***REMOVED*** return d.tag ***REMOVED***

// Done reports whether there are more bytes to process in this message.
func (d *Decoder) Done() bool ***REMOVED*** return len(d.data) == 0 ***REMOVED***

// Render implements Renderer.
func (d *Decoder) Render(s string) ***REMOVED*** d.dst.Render(s) ***REMOVED***

// Arg implements Renderer.
//
// During evaluation of macros, the argument positions may be mapped to
// arguments that differ from the original call.
func (d *Decoder) Arg(i int) interface***REMOVED******REMOVED*** ***REMOVED***
	if d.macroArg != 0 ***REMOVED***
		if i != 1 ***REMOVED***
			panic("catmsg: only macros with single argument supported")
		***REMOVED***
		i = d.macroArg
	***REMOVED***
	return d.dst.Arg(i)
***REMOVED***

// DecodeUint decodes a number that was encoded with EncodeUint and advances the
// position.
func (d *Decoder) DecodeUint() uint64 ***REMOVED***
	x, n, err := decodeUintString(d.data)
	d.data = d.data[n:]
	if err != nil ***REMOVED***
		d.setError(err)
	***REMOVED***
	return x
***REMOVED***

// DecodeString decodes a string that was encoded with EncodeString and advances
// the position.
func (d *Decoder) DecodeString() string ***REMOVED***
	size := d.DecodeUint()
	s := d.data[:size]
	d.data = d.data[size:]
	return s
***REMOVED***

// SkipMessage skips the message at the current location and advances the
// position.
func (d *Decoder) SkipMessage() ***REMOVED***
	n := int(d.DecodeUint())
	d.data = d.data[n:]
***REMOVED***

// Execute decodes and evaluates msg.
//
// Only one goroutine may call execute.
func (d *Decoder) Execute(msg string) error ***REMOVED***
	d.err = nil
	if !d.execute(msg) ***REMOVED***
		return ErrNoMatch
	***REMOVED***
	return d.err
***REMOVED***

func (d *Decoder) execute(msg string) bool ***REMOVED***
	saved := d.data
	d.data = msg
	ok := d.executeMessage()
	d.data = saved
	return ok
***REMOVED***

// executeMessageFromData is like execute, but also decodes a leading message
// size and clips the given string accordingly.
//
// It reports the number of bytes consumed and whether a message was selected.
func (d *Decoder) executeMessageFromData(s string) (n int, ok bool) ***REMOVED***
	saved := d.data
	d.data = s
	size := int(d.DecodeUint())
	n = len(s) - len(d.data)
	// Sanitize the setting. This allows skipping a size argument for
	// RawString and method Done.
	d.data = d.data[:size]
	ok = d.executeMessage()
	n += size - len(d.data)
	d.data = saved
	return n, ok
***REMOVED***

var errUnknownHandler = errors.New("catmsg: string contains unsupported handler")

// executeMessage reads the handle id, initializes the decoder and executes the
// message. It is assumed that all of d.data[d.p:] is the single message.
func (d *Decoder) executeMessage() bool ***REMOVED***
	if d.Done() ***REMOVED***
		// We interpret no data as a valid empty message.
		return true
	***REMOVED***
	handle := d.DecodeUint()

	var fn Handler
	mutex.Lock()
	if int(handle) < len(handlers) ***REMOVED***
		fn = handlers[handle]
	***REMOVED***
	mutex.Unlock()
	if fn == nil ***REMOVED***
		d.setError(errUnknownHandler)
		d.execute(fmt.Sprintf("\x02$!(UNKNOWNMSGHANDLER=%#x)", handle))
		return true
	***REMOVED***
	return fn(d)
***REMOVED***

// ExecuteMessage decodes and executes the message at the current position.
func (d *Decoder) ExecuteMessage() bool ***REMOVED***
	n, ok := d.executeMessageFromData(d.data)
	d.data = d.data[n:]
	return ok
***REMOVED***

// ExecuteSubstitution executes the message corresponding to the substitution
// as encoded by EncodeSubstitution.
func (d *Decoder) ExecuteSubstitution() ***REMOVED***
	switch x := d.DecodeUint(); x ***REMOVED***
	case substituteVar:
		offset := d.DecodeUint()
		d.executeMessageFromData(d.vars[offset:])
	case substituteMacro:
		name := d.DecodeString()
		data, ok := d.macros.Lookup(name)
		old := d.macroArg
		// TODO: support macros of arity other than 1.
		d.macroArg = int(d.DecodeUint())
		switch ***REMOVED***
		case !ok:
			// TODO: detect this at creation time.
			d.setError(fmt.Errorf("catmsg: undefined macro %q", name))
			fallthrough
		case !d.execute(data):
			d.dst.Render(name) // fall back to macro name.
		***REMOVED***
		d.macroArg = old
	case substituteError:
		d.dst.Render(d.DecodeString())
	default:
		panic("catmsg: unreachable")
	***REMOVED***
***REMOVED***
