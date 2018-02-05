// Copyright 2015 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package search

import (
	"golang.org/x/text/internal/colltab"
)

// TODO: handle variable primary weights?

func (p *Pattern) deleteEmptyElements() ***REMOVED***
	k := 0
	for _, e := range p.ce ***REMOVED***
		if !isIgnorable(p.m, e) ***REMOVED***
			p.ce[k] = e
			k++
		***REMOVED***
	***REMOVED***
	p.ce = p.ce[:k]
***REMOVED***

func isIgnorable(m *Matcher, e colltab.Elem) bool ***REMOVED***
	if e.Primary() > 0 ***REMOVED***
		return false
	***REMOVED***
	if e.Secondary() > 0 ***REMOVED***
		if !m.ignoreDiacritics ***REMOVED***
			return false
		***REMOVED***
		// Primary value is 0 and ignoreDiacritics is true. In this case we
		// ignore the tertiary element, as it only pertains to the modifier.
		return true
	***REMOVED***
	// TODO: further distinguish once we have the new implementation.
	if !(m.ignoreWidth || m.ignoreCase) && e.Tertiary() > 0 ***REMOVED***
		return false
	***REMOVED***
	// TODO: we ignore the Quaternary level for now.
	return true
***REMOVED***

// TODO: Use a Boyer-Moore-like algorithm (probably Sunday) for searching.

func (p *Pattern) forwardSearch(it *colltab.Iter) (start, end int) ***REMOVED***
	for start := 0; it.Next(); it.Reset(start) ***REMOVED***
		nextStart := it.End()
		if end := p.searchOnce(it); end != -1 ***REMOVED***
			return start, end
		***REMOVED***
		start = nextStart
	***REMOVED***
	return -1, -1
***REMOVED***

func (p *Pattern) anchoredForwardSearch(it *colltab.Iter) (start, end int) ***REMOVED***
	if it.Next() ***REMOVED***
		if end := p.searchOnce(it); end != -1 ***REMOVED***
			return 0, end
		***REMOVED***
	***REMOVED***
	return -1, -1
***REMOVED***

// next advances to the next weight in a pattern. f must return one of the
// weights of a collation element. next will advance to the first non-zero
// weight and return this weight and true if it exists, or 0, false otherwise.
func (p *Pattern) next(i *int, f func(colltab.Elem) int) (weight int, ok bool) ***REMOVED***
	for *i < len(p.ce) ***REMOVED***
		v := f(p.ce[*i])
		*i++
		if v != 0 ***REMOVED***
			// Skip successive ignorable values.
			for ; *i < len(p.ce) && f(p.ce[*i]) == 0; *i++ ***REMOVED***
			***REMOVED***
			return v, true
		***REMOVED***
	***REMOVED***
	return 0, false
***REMOVED***

// TODO: remove this function once Elem is internal and Tertiary returns int.
func tertiary(e colltab.Elem) int ***REMOVED***
	return int(e.Tertiary())
***REMOVED***

// searchOnce tries to match the pattern s.p at the text position i. s.buf needs
// to be filled with collation elements of the first segment, where n is the
// number of source bytes consumed for this segment. It will return the end
// position of the match or -1.
func (p *Pattern) searchOnce(it *colltab.Iter) (end int) ***REMOVED***
	var pLevel [4]int

	m := p.m
	for ***REMOVED***
		k := 0
		for ; k < it.N; k++ ***REMOVED***
			if v := it.Elems[k].Primary(); v > 0 ***REMOVED***
				if w, ok := p.next(&pLevel[0], colltab.Elem.Primary); !ok || v != w ***REMOVED***
					return -1
				***REMOVED***
			***REMOVED***

			if !m.ignoreDiacritics ***REMOVED***
				if v := it.Elems[k].Secondary(); v > 0 ***REMOVED***
					if w, ok := p.next(&pLevel[1], colltab.Elem.Secondary); !ok || v != w ***REMOVED***
						return -1
					***REMOVED***
				***REMOVED***
			***REMOVED*** else if it.Elems[k].Primary() == 0 ***REMOVED***
				// We ignore tertiary values of collation elements of the
				// secondary level.
				continue
			***REMOVED***

			// TODO: distinguish between case and width. This will be easier to
			// implement after we moved to the new collation implementation.
			if !m.ignoreWidth && !m.ignoreCase ***REMOVED***
				if v := it.Elems[k].Tertiary(); v > 0 ***REMOVED***
					if w, ok := p.next(&pLevel[2], tertiary); !ok || int(v) != w ***REMOVED***
						return -1
					***REMOVED***
				***REMOVED***
			***REMOVED***
			// TODO: check quaternary weight
		***REMOVED***
		it.Discard() // Remove the current segment from the buffer.

		// Check for completion.
		switch ***REMOVED***
		// If any of these cases match, we are not at the end.
		case pLevel[0] < len(p.ce):
		case !m.ignoreDiacritics && pLevel[1] < len(p.ce):
		case !(m.ignoreWidth || m.ignoreCase) && pLevel[2] < len(p.ce):
		default:
			// At this point, both the segment and pattern has matched fully.
			// However, the segment may still be have trailing modifiers.
			// This can be verified by another call to next.
			end = it.End()
			if it.Next() && it.Elems[0].Primary() == 0 ***REMOVED***
				if !m.ignoreDiacritics ***REMOVED***
					return -1
				***REMOVED***
				end = it.End()
			***REMOVED***
			return end
		***REMOVED***

		// Fill the buffer with the next batch of collation elements.
		if !it.Next() ***REMOVED***
			return -1
		***REMOVED***
	***REMOVED***
***REMOVED***
