// Copyright 2012 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"math"
	"math/rand"
	"strings"
	"unicode"
	"unicode/utf16"
	"unicode/utf8"

	"golang.org/x/text/language"
	"golang.org/x/text/unicode/norm"
)

// TODO: replace with functionality in language package.
// parent computes the parent language for the given language.
// It returns false if the parent is already root.
func parent(locale string) (parent string, ok bool) ***REMOVED***
	if locale == "und" ***REMOVED***
		return "", false
	***REMOVED***
	if i := strings.LastIndex(locale, "-"); i != -1 ***REMOVED***
		return locale[:i], true
	***REMOVED***
	return "und", true
***REMOVED***

// rewriter is used to both unique strings and create variants of strings
// to add to the test set.
type rewriter struct ***REMOVED***
	seen     map[string]bool
	addCases bool
***REMOVED***

func newRewriter() *rewriter ***REMOVED***
	return &rewriter***REMOVED***
		seen: make(map[string]bool),
	***REMOVED***
***REMOVED***

func (r *rewriter) insert(a []string, s string) []string ***REMOVED***
	if !r.seen[s] ***REMOVED***
		r.seen[s] = true
		a = append(a, s)
	***REMOVED***
	return a
***REMOVED***

// rewrite takes a sequence of strings in, adds variants of the these strings
// based on options and removes duplicates.
func (r *rewriter) rewrite(ss []string) []string ***REMOVED***
	ns := []string***REMOVED******REMOVED***
	for _, s := range ss ***REMOVED***
		ns = r.insert(ns, s)
		if r.addCases ***REMOVED***
			rs := []rune(s)
			rn := rs[0]
			for c := unicode.SimpleFold(rn); c != rn; c = unicode.SimpleFold(c) ***REMOVED***
				rs[0] = c
				ns = r.insert(ns, string(rs))
			***REMOVED***
		***REMOVED***
	***REMOVED***
	return ns
***REMOVED***

// exemplarySet holds a parsed set of characters from the exemplarCharacters table.
type exemplarySet struct ***REMOVED***
	typ       exemplarType
	set       []string
	charIndex int // cumulative total of phrases, including this set
***REMOVED***

type phraseGenerator struct ***REMOVED***
	sets [exN]exemplarySet
	n    int
***REMOVED***

func (g *phraseGenerator) init(id string) ***REMOVED***
	ec := exemplarCharacters
	loc := language.Make(id).String()
	// get sets for locale or parent locale if the set is not defined.
	for i := range g.sets ***REMOVED***
		for p, ok := loc, true; ok; p, ok = parent(p) ***REMOVED***
			if set, ok := ec[p]; ok && set[i] != "" ***REMOVED***
				g.sets[i].set = strings.Split(set[i], " ")
				break
			***REMOVED***
		***REMOVED***
	***REMOVED***
	r := newRewriter()
	r.addCases = *cases
	for i := range g.sets ***REMOVED***
		g.sets[i].set = r.rewrite(g.sets[i].set)
	***REMOVED***
	// compute indexes
	for i, set := range g.sets ***REMOVED***
		g.n += len(set.set)
		g.sets[i].charIndex = g.n
	***REMOVED***
***REMOVED***

// phrase returns the ith phrase, where i < g.n.
func (g *phraseGenerator) phrase(i int) string ***REMOVED***
	for _, set := range g.sets ***REMOVED***
		if i < set.charIndex ***REMOVED***
			return set.set[i-(set.charIndex-len(set.set))]
		***REMOVED***
	***REMOVED***
	panic("index out of range")
***REMOVED***

// generate generates inputs by combining all pairs of examplar strings.
// If doNorm is true, all input strings are normalized to NFC.
// TODO: allow other variations, statistical models, and random
// trailing sequences.
func (g *phraseGenerator) generate(doNorm bool) []Input ***REMOVED***
	const (
		M         = 1024 * 1024
		buf8Size  = 30 * M
		buf16Size = 10 * M
	)
	// TODO: use a better way to limit the input size.
	if sq := int(math.Sqrt(float64(*limit))); g.n > sq ***REMOVED***
		g.n = sq
	***REMOVED***
	size := g.n * g.n
	a := make([]Input, 0, size)
	buf8 := make([]byte, 0, buf8Size)
	buf16 := make([]uint16, 0, buf16Size)

	addInput := func(str string) ***REMOVED***
		buf8 = buf8[len(buf8):]
		buf16 = buf16[len(buf16):]
		if len(str) > cap(buf8) ***REMOVED***
			buf8 = make([]byte, 0, buf8Size)
		***REMOVED***
		if len(str) > cap(buf16) ***REMOVED***
			buf16 = make([]uint16, 0, buf16Size)
		***REMOVED***
		if doNorm ***REMOVED***
			buf8 = norm.NFD.AppendString(buf8, str)
		***REMOVED*** else ***REMOVED***
			buf8 = append(buf8, str...)
		***REMOVED***
		buf16 = appendUTF16(buf16, buf8)
		a = append(a, makeInput(buf8, buf16))
	***REMOVED***
	for i := 0; i < g.n; i++ ***REMOVED***
		p1 := g.phrase(i)
		addInput(p1)
		for j := 0; j < g.n; j++ ***REMOVED***
			p2 := g.phrase(j)
			addInput(p1 + p2)
		***REMOVED***
	***REMOVED***
	// permutate
	rnd := rand.New(rand.NewSource(int64(rand.Int())))
	for i := range a ***REMOVED***
		j := i + rnd.Intn(len(a)-i)
		a[i], a[j] = a[j], a[i]
		a[i].index = i // allow restoring this order if input is used multiple times.
	***REMOVED***
	return a
***REMOVED***

func appendUTF16(buf []uint16, s []byte) []uint16 ***REMOVED***
	for len(s) > 0 ***REMOVED***
		r, sz := utf8.DecodeRune(s)
		s = s[sz:]
		r1, r2 := utf16.EncodeRune(r)
		if r1 != 0xFFFD ***REMOVED***
			buf = append(buf, uint16(r1), uint16(r2))
		***REMOVED*** else ***REMOVED***
			buf = append(buf, uint16(r))
		***REMOVED***
	***REMOVED***
	return buf
***REMOVED***
