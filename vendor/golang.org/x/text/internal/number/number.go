// Copyright 2016 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:generate go run gen.go gen_common.go

// Package number contains tools and data for formatting numbers.
package number

import (
	"unicode/utf8"

	"golang.org/x/text/internal"
	"golang.org/x/text/language"
)

// Info holds number formatting configuration data.
type Info struct ***REMOVED***
	system   systemData // numbering system information
	symIndex symOffset  // index to symbols
***REMOVED***

// InfoFromLangID returns a Info for the given compact language identifier and
// numbering system identifier. If system is the empty string, the default
// numbering system will be taken for that language.
func InfoFromLangID(compactIndex int, numberSystem string) Info ***REMOVED***
	p := langToDefaults[compactIndex]
	// Lookup the entry for the language.
	pSymIndex := symOffset(0) // Default: Latin, default symbols
	system, ok := systemMap[numberSystem]
	if !ok ***REMOVED***
		// Take the value for the default numbering system. This is by far the
		// most common case as an alternative numbering system is hardly used.
		if p&hasNonLatnMask == 0 ***REMOVED*** // Latn digits.
			pSymIndex = p
		***REMOVED*** else ***REMOVED*** // Non-Latn or multiple numbering systems.
			// Take the first entry from the alternatives list.
			data := langToAlt[p&^hasNonLatnMask]
			pSymIndex = data.symIndex
			system = data.system
		***REMOVED***
	***REMOVED*** else ***REMOVED***
		langIndex := compactIndex
		ns := system
	outerLoop:
		for ; ; p = langToDefaults[langIndex] ***REMOVED***
			if p&hasNonLatnMask == 0 ***REMOVED***
				if ns == 0 ***REMOVED***
					// The index directly points to the symbol data.
					pSymIndex = p
					break
				***REMOVED***
				// Move to the parent and retry.
				langIndex = int(internal.Parent[langIndex])
			***REMOVED*** else ***REMOVED***
				// The index points to a list of symbol data indexes.
				for _, e := range langToAlt[p&^hasNonLatnMask:] ***REMOVED***
					if int(e.compactTag) != langIndex ***REMOVED***
						if langIndex == 0 ***REMOVED***
							// The CLDR root defines full symbol information for
							// all numbering systems (even though mostly by
							// means of aliases). Fall back to the default entry
							// for Latn if there is no data for the numbering
							// system of this language.
							if ns == 0 ***REMOVED***
								break
							***REMOVED***
							// Fall back to Latin and start from the original
							// language. See
							// http://unicode.org/reports/tr35/#Locale_Inheritance.
							ns = numLatn
							langIndex = compactIndex
							continue outerLoop
						***REMOVED***
						// Fall back to parent.
						langIndex = int(internal.Parent[langIndex])
					***REMOVED*** else if e.system == ns ***REMOVED***
						pSymIndex = e.symIndex
						break outerLoop
					***REMOVED***
				***REMOVED***
			***REMOVED***
		***REMOVED***
	***REMOVED***
	if int(system) >= len(numSysData) ***REMOVED*** // algorithmic
		// Will generate ASCII digits in case the user inadvertently calls
		// WriteDigit or Digit on it.
		d := numSysData[0]
		d.id = system
		return Info***REMOVED***
			system:   d,
			symIndex: pSymIndex,
		***REMOVED***
	***REMOVED***
	return Info***REMOVED***
		system:   numSysData[system],
		symIndex: pSymIndex,
	***REMOVED***
***REMOVED***

// InfoFromTag returns a Info for the given language tag.
func InfoFromTag(t language.Tag) Info ***REMOVED***
	for ***REMOVED***
		if index, ok := language.CompactIndex(t); ok ***REMOVED***
			return InfoFromLangID(index, t.TypeForKey("nu"))
		***REMOVED***
		t = t.Parent()
	***REMOVED***
***REMOVED***

// IsDecimal reports if the numbering system can convert decimal to native
// symbols one-to-one.
func (n Info) IsDecimal() bool ***REMOVED***
	return int(n.system.id) < len(numSysData)
***REMOVED***

// WriteDigit writes the UTF-8 sequence for n corresponding to the given ASCII
// digit to dst and reports the number of bytes written. dst must be large
// enough to hold the rune (can be up to utf8.UTFMax bytes).
func (n Info) WriteDigit(dst []byte, asciiDigit rune) int ***REMOVED***
	copy(dst, n.system.zero[:n.system.digitSize])
	dst[n.system.digitSize-1] += byte(asciiDigit - '0')
	return int(n.system.digitSize)
***REMOVED***

// AppendDigit appends the UTF-8 sequence for n corresponding to the given digit
// to dst and reports the number of bytes written. dst must be large enough to
// hold the rune (can be up to utf8.UTFMax bytes).
func (n Info) AppendDigit(dst []byte, digit byte) []byte ***REMOVED***
	dst = append(dst, n.system.zero[:n.system.digitSize]...)
	dst[len(dst)-1] += digit
	return dst
***REMOVED***

// Digit returns the digit for the numbering system for the corresponding ASCII
// value. For example, ni.Digit('3') could return '三'. Note that the argument
// is the rune constant '3', which equals 51, not the integer constant 3.
func (n Info) Digit(asciiDigit rune) rune ***REMOVED***
	var x [utf8.UTFMax]byte
	n.WriteDigit(x[:], asciiDigit)
	r, _ := utf8.DecodeRune(x[:])
	return r
***REMOVED***

// Symbol returns the string for the given symbol type.
func (n Info) Symbol(t SymbolType) string ***REMOVED***
	return symData.Elem(int(symIndex[n.symIndex][t]))
***REMOVED***

func formatForLang(t language.Tag, index []byte) *Pattern ***REMOVED***
	for ; ; t = t.Parent() ***REMOVED***
		if x, ok := language.CompactIndex(t); ok ***REMOVED***
			return &formats[index[x]]
		***REMOVED***
	***REMOVED***
***REMOVED***
