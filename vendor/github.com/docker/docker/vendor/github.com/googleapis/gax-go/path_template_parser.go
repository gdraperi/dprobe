// Copyright 2016, Google Inc.
// All rights reserved.
//
// Redistribution and use in source and binary forms, with or without
// modification, are permitted provided that the following conditions are
// met:
//
//     * Redistributions of source code must retain the above copyright
// notice, this list of conditions and the following disclaimer.
//     * Redistributions in binary form must reproduce the above
// copyright notice, this list of conditions and the following disclaimer
// in the documentation and/or other materials provided with the
// distribution.
//     * Neither the name of Google Inc. nor the names of its
// contributors may be used to endorse or promote products derived from
// this software without specific prior written permission.
//
// THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS
// "AS IS" AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT
// LIMITED TO, THE IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR
// A PARTICULAR PURPOSE ARE DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT
// OWNER OR CONTRIBUTORS BE LIABLE FOR ANY DIRECT, INDIRECT, INCIDENTAL,
// SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES (INCLUDING, BUT NOT
// LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES; LOSS OF USE,
// DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND ON ANY
// THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT
// (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE
// OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.

package gax

import (
	"fmt"
	"io"
	"strings"
)

// This parser follows the syntax of path templates, from
// https://github.com/googleapis/googleapis/blob/master/google/api/http.proto.
// The differences are that there is no custom verb, we allow the initial slash
// to be absent, and that we are not strict as
// https://tools.ietf.org/html/rfc6570 about the characters in identifiers and
// literals.

type pathTemplateParser struct ***REMOVED***
	r                *strings.Reader
	runeCount        int             // the number of the current rune in the original string
	nextVar          int             // the number to use for the next unnamed variable
	seenName         map[string]bool // names we've seen already
	seenPathWildcard bool            // have we seen "**" already?
***REMOVED***

func parsePathTemplate(template string) (pt *PathTemplate, err error) ***REMOVED***
	p := &pathTemplateParser***REMOVED***
		r:        strings.NewReader(template),
		seenName: map[string]bool***REMOVED******REMOVED***,
	***REMOVED***

	// Handle panics with strings like errors.
	// See pathTemplateParser.error, below.
	defer func() ***REMOVED***
		if x := recover(); x != nil ***REMOVED***
			errmsg, ok := x.(errString)
			if !ok ***REMOVED***
				panic(x)
			***REMOVED***
			pt = nil
			err = ParseError***REMOVED***p.runeCount, template, string(errmsg)***REMOVED***
		***REMOVED***
	***REMOVED***()

	segs := p.template()
	// If there is a path wildcard, set its length. We can't do this
	// until we know how many segments we've got all together.
	for i, seg := range segs ***REMOVED***
		if _, ok := seg.matcher.(pathWildcardMatcher); ok ***REMOVED***
			segs[i].matcher = pathWildcardMatcher(len(segs) - i - 1)
			break
		***REMOVED***
	***REMOVED***
	return &PathTemplate***REMOVED***segments: segs***REMOVED***, nil

***REMOVED***

// Used to indicate errors "thrown" by this parser. We don't use string because
// many parts of the standard library panic with strings.
type errString string

// Terminates parsing immediately with an error.
func (p *pathTemplateParser) error(msg string) ***REMOVED***
	panic(errString(msg))
***REMOVED***

// Template = [ "/" ] Segments
func (p *pathTemplateParser) template() []segment ***REMOVED***
	var segs []segment
	if p.consume('/') ***REMOVED***
		// Initial '/' needs an initial empty matcher.
		segs = append(segs, segment***REMOVED***matcher: labelMatcher("")***REMOVED***)
	***REMOVED***
	return append(segs, p.segments("")...)
***REMOVED***

// Segments = Segment ***REMOVED*** "/" Segment ***REMOVED***
func (p *pathTemplateParser) segments(name string) []segment ***REMOVED***
	var segs []segment
	for ***REMOVED***
		subsegs := p.segment(name)
		segs = append(segs, subsegs...)
		if !p.consume('/') ***REMOVED***
			break
		***REMOVED***
	***REMOVED***
	return segs
***REMOVED***

// Segment  = "*" | "**" | LITERAL | Variable
func (p *pathTemplateParser) segment(name string) []segment ***REMOVED***
	if p.consume('*') ***REMOVED***
		if name == "" ***REMOVED***
			name = fmt.Sprintf("$%d", p.nextVar)
			p.nextVar++
		***REMOVED***
		if p.consume('*') ***REMOVED***
			if p.seenPathWildcard ***REMOVED***
				p.error("multiple '**' disallowed")
			***REMOVED***
			p.seenPathWildcard = true
			// We'll change 0 to the right number at the end.
			return []segment***REMOVED******REMOVED***name: name, matcher: pathWildcardMatcher(0)***REMOVED******REMOVED***
		***REMOVED***
		return []segment***REMOVED******REMOVED***name: name, matcher: wildcardMatcher(0)***REMOVED******REMOVED***
	***REMOVED***
	if p.consume('***REMOVED***') ***REMOVED***
		if name != "" ***REMOVED***
			p.error("recursive named bindings are not allowed")
		***REMOVED***
		return p.variable()
	***REMOVED***
	return []segment***REMOVED******REMOVED***name: name, matcher: labelMatcher(p.literal())***REMOVED******REMOVED***
***REMOVED***

// Variable = "***REMOVED***" FieldPath [ "=" Segments ] "***REMOVED***"
// "***REMOVED***" is already consumed.
func (p *pathTemplateParser) variable() []segment ***REMOVED***
	// Simplification: treat FieldPath as LITERAL, instead of IDENT ***REMOVED*** '.' IDENT ***REMOVED***
	name := p.literal()
	if p.seenName[name] ***REMOVED***
		p.error(name + " appears multiple times")
	***REMOVED***
	p.seenName[name] = true
	var segs []segment
	if p.consume('=') ***REMOVED***
		segs = p.segments(name)
	***REMOVED*** else ***REMOVED***
		// "***REMOVED***var***REMOVED***" is equivalent to "***REMOVED***var=****REMOVED***"
		segs = []segment***REMOVED******REMOVED***name: name, matcher: wildcardMatcher(0)***REMOVED******REMOVED***
	***REMOVED***
	if !p.consume('***REMOVED***') ***REMOVED***
		p.error("expected '***REMOVED***'")
	***REMOVED***
	return segs
***REMOVED***

// A literal is any sequence of characters other than a few special ones.
// The list of stop characters is not quite the same as in the template RFC.
func (p *pathTemplateParser) literal() string ***REMOVED***
	lit := p.consumeUntil("/****REMOVED******REMOVED***=")
	if lit == "" ***REMOVED***
		p.error("empty literal")
	***REMOVED***
	return lit
***REMOVED***

// Read runes until EOF or one of the runes in stopRunes is encountered.
// If the latter, unread the stop rune. Return the accumulated runes as a string.
func (p *pathTemplateParser) consumeUntil(stopRunes string) string ***REMOVED***
	var runes []rune
	for ***REMOVED***
		r, ok := p.readRune()
		if !ok ***REMOVED***
			break
		***REMOVED***
		if strings.IndexRune(stopRunes, r) >= 0 ***REMOVED***
			p.unreadRune()
			break
		***REMOVED***
		runes = append(runes, r)
	***REMOVED***
	return string(runes)
***REMOVED***

// If the next rune is r, consume it and return true.
// Otherwise, leave the input unchanged and return false.
func (p *pathTemplateParser) consume(r rune) bool ***REMOVED***
	rr, ok := p.readRune()
	if !ok ***REMOVED***
		return false
	***REMOVED***
	if r == rr ***REMOVED***
		return true
	***REMOVED***
	p.unreadRune()
	return false
***REMOVED***

// Read the next rune from the input. Return it.
// The second return value is false at EOF.
func (p *pathTemplateParser) readRune() (rune, bool) ***REMOVED***
	r, _, err := p.r.ReadRune()
	if err == io.EOF ***REMOVED***
		return r, false
	***REMOVED***
	if err != nil ***REMOVED***
		p.error(err.Error())
	***REMOVED***
	p.runeCount++
	return r, true
***REMOVED***

// Put the last rune that was read back on the input.
func (p *pathTemplateParser) unreadRune() ***REMOVED***
	if err := p.r.UnreadRune(); err != nil ***REMOVED***
		p.error(err.Error())
	***REMOVED***
	p.runeCount--
***REMOVED***
