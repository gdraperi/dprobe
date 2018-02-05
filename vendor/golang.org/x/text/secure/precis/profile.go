// Copyright 2015 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package precis

import (
	"bytes"
	"errors"
	"unicode/utf8"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"
	"golang.org/x/text/runes"
	"golang.org/x/text/secure/bidirule"
	"golang.org/x/text/transform"
	"golang.org/x/text/width"
)

var (
	errDisallowedRune = errors.New("precis: disallowed rune encountered")
)

var dpTrie = newDerivedPropertiesTrie(0)

// A Profile represents a set of rules for normalizing and validating strings in
// the PRECIS framework.
type Profile struct ***REMOVED***
	options
	class *class
***REMOVED***

// NewIdentifier creates a new PRECIS profile based on the Identifier string
// class. Profiles created from this class are suitable for use where safety is
// prioritized over expressiveness like network identifiers, user accounts, chat
// rooms, and file names.
func NewIdentifier(opts ...Option) *Profile ***REMOVED***
	return &Profile***REMOVED***
		options: getOpts(opts...),
		class:   identifier,
	***REMOVED***
***REMOVED***

// NewFreeform creates a new PRECIS profile based on the Freeform string class.
// Profiles created from this class are suitable for use where expressiveness is
// prioritized over safety like passwords, and display-elements such as
// nicknames in a chat room.
func NewFreeform(opts ...Option) *Profile ***REMOVED***
	return &Profile***REMOVED***
		options: getOpts(opts...),
		class:   freeform,
	***REMOVED***
***REMOVED***

// NewTransformer creates a new transform.Transformer that performs the PRECIS
// preparation and enforcement steps on the given UTF-8 encoded bytes.
func (p *Profile) NewTransformer() *Transformer ***REMOVED***
	var ts []transform.Transformer

	// These transforms are applied in the order defined in
	// https://tools.ietf.org/html/rfc7564#section-7

	// RFC 8266 §2.1:
	//
	//     Implementation experience has shown that applying the rules for the
	//     Nickname profile is not an idempotent procedure for all code points.
	//     Therefore, an implementation SHOULD apply the rules repeatedly until
	//     the output string is stable; if the output string does not stabilize
	//     after reapplying the rules three (3) additional times after the first
	//     application, the implementation SHOULD terminate application of the
	//     rules and reject the input string as invalid.
	//
	// There is no known string that will change indefinitely, so repeat 4 times
	// and rely on the Span method to keep things relatively performant.
	r := 1
	if p.options.repeat ***REMOVED***
		r = 4
	***REMOVED***
	for ; r > 0; r-- ***REMOVED***
		if p.options.foldWidth ***REMOVED***
			ts = append(ts, width.Fold)
		***REMOVED***

		for _, f := range p.options.additional ***REMOVED***
			ts = append(ts, f())
		***REMOVED***

		if p.options.cases != nil ***REMOVED***
			ts = append(ts, p.options.cases)
		***REMOVED***

		ts = append(ts, p.options.norm)

		if p.options.bidiRule ***REMOVED***
			ts = append(ts, bidirule.New())
		***REMOVED***

		ts = append(ts, &checker***REMOVED***p: p, allowed: p.Allowed()***REMOVED***)
	***REMOVED***

	// TODO: Add the disallow empty rule with a dummy transformer?

	return &Transformer***REMOVED***transform.Chain(ts...)***REMOVED***
***REMOVED***

var errEmptyString = errors.New("precis: transformation resulted in empty string")

type buffers struct ***REMOVED***
	src  []byte
	buf  [2][]byte
	next int
***REMOVED***

func (b *buffers) apply(t transform.SpanningTransformer) (err error) ***REMOVED***
	n, err := t.Span(b.src, true)
	if err != transform.ErrEndOfSpan ***REMOVED***
		return err
	***REMOVED***
	x := b.next & 1
	if b.buf[x] == nil ***REMOVED***
		b.buf[x] = make([]byte, 0, 8+len(b.src)+len(b.src)>>2)
	***REMOVED***
	span := append(b.buf[x][:0], b.src[:n]...)
	b.src, _, err = transform.Append(t, span, b.src[n:])
	b.buf[x] = b.src
	b.next++
	return err
***REMOVED***

// Pre-allocate transformers when possible. In some cases this avoids allocation.
var (
	foldWidthT transform.SpanningTransformer = width.Fold
	lowerCaseT transform.SpanningTransformer = cases.Lower(language.Und, cases.HandleFinalSigma(false))
)

// TODO: make this a method on profile.

func (b *buffers) enforce(p *Profile, src []byte, comparing bool) (str []byte, err error) ***REMOVED***
	b.src = src

	ascii := true
	for _, c := range src ***REMOVED***
		if c >= utf8.RuneSelf ***REMOVED***
			ascii = false
			break
		***REMOVED***
	***REMOVED***
	// ASCII fast path.
	if ascii ***REMOVED***
		for _, f := range p.options.additional ***REMOVED***
			if err = b.apply(f()); err != nil ***REMOVED***
				return nil, err
			***REMOVED***
		***REMOVED***
		switch ***REMOVED***
		case p.options.asciiLower || (comparing && p.options.ignorecase):
			for i, c := range b.src ***REMOVED***
				if 'A' <= c && c <= 'Z' ***REMOVED***
					b.src[i] = c ^ 1<<5
				***REMOVED***
			***REMOVED***
		case p.options.cases != nil:
			b.apply(p.options.cases)
		***REMOVED***
		c := checker***REMOVED***p: p***REMOVED***
		if _, err := c.span(b.src, true); err != nil ***REMOVED***
			return nil, err
		***REMOVED***
		if p.disallow != nil ***REMOVED***
			for _, c := range b.src ***REMOVED***
				if p.disallow.Contains(rune(c)) ***REMOVED***
					return nil, errDisallowedRune
				***REMOVED***
			***REMOVED***
		***REMOVED***
		if p.options.disallowEmpty && len(b.src) == 0 ***REMOVED***
			return nil, errEmptyString
		***REMOVED***
		return b.src, nil
	***REMOVED***

	// These transforms are applied in the order defined in
	// https://tools.ietf.org/html/rfc8264#section-7

	r := 1
	if p.options.repeat ***REMOVED***
		r = 4
	***REMOVED***
	for ; r > 0; r-- ***REMOVED***
		// TODO: allow different width transforms options.
		if p.options.foldWidth || (p.options.ignorecase && comparing) ***REMOVED***
			b.apply(foldWidthT)
		***REMOVED***
		for _, f := range p.options.additional ***REMOVED***
			if err = b.apply(f()); err != nil ***REMOVED***
				return nil, err
			***REMOVED***
		***REMOVED***
		if p.options.cases != nil ***REMOVED***
			b.apply(p.options.cases)
		***REMOVED***
		if comparing && p.options.ignorecase ***REMOVED***
			b.apply(lowerCaseT)
		***REMOVED***
		b.apply(p.norm)
		if p.options.bidiRule && !bidirule.Valid(b.src) ***REMOVED***
			return nil, bidirule.ErrInvalid
		***REMOVED***
		c := checker***REMOVED***p: p***REMOVED***
		if _, err := c.span(b.src, true); err != nil ***REMOVED***
			return nil, err
		***REMOVED***
		if p.disallow != nil ***REMOVED***
			for i := 0; i < len(b.src); ***REMOVED***
				r, size := utf8.DecodeRune(b.src[i:])
				if p.disallow.Contains(r) ***REMOVED***
					return nil, errDisallowedRune
				***REMOVED***
				i += size
			***REMOVED***
		***REMOVED***
		if p.options.disallowEmpty && len(b.src) == 0 ***REMOVED***
			return nil, errEmptyString
		***REMOVED***
	***REMOVED***
	return b.src, nil
***REMOVED***

// Append appends the result of applying p to src writing the result to dst.
// It returns an error if the input string is invalid.
func (p *Profile) Append(dst, src []byte) ([]byte, error) ***REMOVED***
	var buf buffers
	b, err := buf.enforce(p, src, false)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	return append(dst, b...), nil
***REMOVED***

func processBytes(p *Profile, b []byte, key bool) ([]byte, error) ***REMOVED***
	var buf buffers
	b, err := buf.enforce(p, b, key)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	if buf.next == 0 ***REMOVED***
		c := make([]byte, len(b))
		copy(c, b)
		return c, nil
	***REMOVED***
	return b, nil
***REMOVED***

// Bytes returns a new byte slice with the result of applying the profile to b.
func (p *Profile) Bytes(b []byte) ([]byte, error) ***REMOVED***
	return processBytes(p, b, false)
***REMOVED***

// AppendCompareKey appends the result of applying p to src (including any
// optional rules to make strings comparable or useful in a map key such as
// applying lowercasing) writing the result to dst. It returns an error if the
// input string is invalid.
func (p *Profile) AppendCompareKey(dst, src []byte) ([]byte, error) ***REMOVED***
	var buf buffers
	b, err := buf.enforce(p, src, true)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	return append(dst, b...), nil
***REMOVED***

func processString(p *Profile, s string, key bool) (string, error) ***REMOVED***
	var buf buffers
	b, err := buf.enforce(p, []byte(s), key)
	if err != nil ***REMOVED***
		return "", err
	***REMOVED***
	return string(b), nil
***REMOVED***

// String returns a string with the result of applying the profile to s.
func (p *Profile) String(s string) (string, error) ***REMOVED***
	return processString(p, s, false)
***REMOVED***

// CompareKey returns a string that can be used for comparison, hashing, or
// collation.
func (p *Profile) CompareKey(s string) (string, error) ***REMOVED***
	return processString(p, s, true)
***REMOVED***

// Compare enforces both strings, and then compares them for bit-string identity
// (byte-for-byte equality). If either string cannot be enforced, the comparison
// is false.
func (p *Profile) Compare(a, b string) bool ***REMOVED***
	var buf buffers

	akey, err := buf.enforce(p, []byte(a), true)
	if err != nil ***REMOVED***
		return false
	***REMOVED***

	buf = buffers***REMOVED******REMOVED***
	bkey, err := buf.enforce(p, []byte(b), true)
	if err != nil ***REMOVED***
		return false
	***REMOVED***

	return bytes.Compare(akey, bkey) == 0
***REMOVED***

// Allowed returns a runes.Set containing every rune that is a member of the
// underlying profile's string class and not disallowed by any profile specific
// rules.
func (p *Profile) Allowed() runes.Set ***REMOVED***
	if p.options.disallow != nil ***REMOVED***
		return runes.Predicate(func(r rune) bool ***REMOVED***
			return p.class.Contains(r) && !p.options.disallow.Contains(r)
		***REMOVED***)
	***REMOVED***
	return p.class
***REMOVED***

type checker struct ***REMOVED***
	p       *Profile
	allowed runes.Set

	beforeBits catBitmap
	termBits   catBitmap
	acceptBits catBitmap
***REMOVED***

func (c *checker) Reset() ***REMOVED***
	c.beforeBits = 0
	c.termBits = 0
	c.acceptBits = 0
***REMOVED***

func (c *checker) span(src []byte, atEOF bool) (n int, err error) ***REMOVED***
	for n < len(src) ***REMOVED***
		e, sz := dpTrie.lookup(src[n:])
		d := categoryTransitions[category(e&catMask)]
		if sz == 0 ***REMOVED***
			if !atEOF ***REMOVED***
				return n, transform.ErrShortSrc
			***REMOVED***
			return n, errDisallowedRune
		***REMOVED***
		doLookAhead := false
		if property(e) < c.p.class.validFrom ***REMOVED***
			if d.rule == nil ***REMOVED***
				return n, errDisallowedRune
			***REMOVED***
			doLookAhead, err = d.rule(c.beforeBits)
			if err != nil ***REMOVED***
				return n, err
			***REMOVED***
		***REMOVED***
		c.beforeBits &= d.keep
		c.beforeBits |= d.set
		if c.termBits != 0 ***REMOVED***
			// We are currently in an unterminated lookahead.
			if c.beforeBits&c.termBits != 0 ***REMOVED***
				c.termBits = 0
				c.acceptBits = 0
			***REMOVED*** else if c.beforeBits&c.acceptBits == 0 ***REMOVED***
				// Invalid continuation of the unterminated lookahead sequence.
				return n, errContext
			***REMOVED***
		***REMOVED***
		if doLookAhead ***REMOVED***
			if c.termBits != 0 ***REMOVED***
				// A previous lookahead run has not been terminated yet.
				return n, errContext
			***REMOVED***
			c.termBits = d.term
			c.acceptBits = d.accept
		***REMOVED***
		n += sz
	***REMOVED***
	if m := c.beforeBits >> finalShift; c.beforeBits&m != m || c.termBits != 0 ***REMOVED***
		err = errContext
	***REMOVED***
	return n, err
***REMOVED***

// TODO: we may get rid of this transform if transform.Chain understands
// something like a Spanner interface.
func (c checker) Transform(dst, src []byte, atEOF bool) (nDst, nSrc int, err error) ***REMOVED***
	short := false
	if len(dst) < len(src) ***REMOVED***
		src = src[:len(dst)]
		atEOF = false
		short = true
	***REMOVED***
	nSrc, err = c.span(src, atEOF)
	nDst = copy(dst, src[:nSrc])
	if short && (err == transform.ErrShortSrc || err == nil) ***REMOVED***
		err = transform.ErrShortDst
	***REMOVED***
	return nDst, nSrc, err
***REMOVED***
