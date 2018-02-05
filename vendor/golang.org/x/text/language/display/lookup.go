// Copyright 2014 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package display

// This file contains common lookup code that is shared between the various
// implementations of Namer and Dictionaries.

import (
	"fmt"
	"sort"
	"strings"

	"golang.org/x/text/language"
)

type namer interface ***REMOVED***
	// name gets the string for the given index. It should walk the
	// inheritance chain if a value is not present in the base index.
	name(idx int) string
***REMOVED***

func nameLanguage(n namer, x interface***REMOVED******REMOVED***) string ***REMOVED***
	t, _ := language.All.Compose(x)
	for ***REMOVED***
		i, _, _ := langTagSet.index(t.Raw())
		if s := n.name(i); s != "" ***REMOVED***
			return s
		***REMOVED***
		if t = t.Parent(); t == language.Und ***REMOVED***
			return ""
		***REMOVED***
	***REMOVED***
***REMOVED***

func nameScript(n namer, x interface***REMOVED******REMOVED***) string ***REMOVED***
	t, _ := language.DeprecatedScript.Compose(x)
	_, s, _ := t.Raw()
	return n.name(scriptIndex.index(s.String()))
***REMOVED***

func nameRegion(n namer, x interface***REMOVED******REMOVED***) string ***REMOVED***
	t, _ := language.DeprecatedRegion.Compose(x)
	_, _, r := t.Raw()
	return n.name(regionIndex.index(r.String()))
***REMOVED***

func nameTag(langN, scrN, regN namer, x interface***REMOVED******REMOVED***) string ***REMOVED***
	t, ok := x.(language.Tag)
	if !ok ***REMOVED***
		return ""
	***REMOVED***
	const form = language.All &^ language.SuppressScript
	if c, err := form.Canonicalize(t); err == nil ***REMOVED***
		t = c
	***REMOVED***
	_, sRaw, rRaw := t.Raw()
	i, scr, reg := langTagSet.index(t.Raw())
	for i != -1 ***REMOVED***
		if str := langN.name(i); str != "" ***REMOVED***
			if hasS, hasR := (scr != language.Script***REMOVED******REMOVED***), (reg != language.Region***REMOVED******REMOVED***); hasS || hasR ***REMOVED***
				ss, sr := "", ""
				if hasS ***REMOVED***
					ss = scrN.name(scriptIndex.index(scr.String()))
				***REMOVED***
				if hasR ***REMOVED***
					sr = regN.name(regionIndex.index(reg.String()))
				***REMOVED***
				// TODO: use patterns in CLDR or at least confirm they are the
				// same for all languages.
				if ss != "" && sr != "" ***REMOVED***
					return fmt.Sprintf("%s (%s, %s)", str, ss, sr)
				***REMOVED***
				if ss != "" || sr != "" ***REMOVED***
					return fmt.Sprintf("%s (%s%s)", str, ss, sr)
				***REMOVED***
			***REMOVED***
			return str
		***REMOVED***
		scr, reg = sRaw, rRaw
		if t = t.Parent(); t == language.Und ***REMOVED***
			return ""
		***REMOVED***
		i, _, _ = langTagSet.index(t.Raw())
	***REMOVED***
	return ""
***REMOVED***

// header contains the data and indexes for a single namer.
// data contains a series of strings concatenated into one. index contains the
// offsets for a string in data. For example, consider a header that defines
// strings for the languages de, el, en, fi, and nl:
//
// 		header***REMOVED***
// 			data: "GermanGreekEnglishDutch",
//  		index: []uint16***REMOVED*** 0, 6, 11, 18, 18, 23 ***REMOVED***,
// 		***REMOVED***
//
// For a language with index i, the string is defined by
// data[index[i]:index[i+1]]. So the number of elements in index is always one
// greater than the number of languages for which header defines a value.
// A string for a language may be empty, which means the name is undefined. In
// the above example, the name for fi (Finnish) is undefined.
type header struct ***REMOVED***
	data  string
	index []uint16
***REMOVED***

// name looks up the name for a tag in the dictionary, given its index.
func (h *header) name(i int) string ***REMOVED***
	if 0 <= i && i < len(h.index)-1 ***REMOVED***
		return h.data[h.index[i]:h.index[i+1]]
	***REMOVED***
	return ""
***REMOVED***

// tagSet is used to find the index of a language in a set of tags.
type tagSet struct ***REMOVED***
	single tagIndex
	long   []string
***REMOVED***

var (
	langTagSet = tagSet***REMOVED***
		single: langIndex,
		long:   langTagsLong,
	***REMOVED***

	// selfTagSet is used for indexing the language strings in their own
	// language.
	selfTagSet = tagSet***REMOVED***
		single: selfIndex,
		long:   selfTagsLong,
	***REMOVED***

	zzzz = language.MustParseScript("Zzzz")
	zz   = language.MustParseRegion("ZZ")
)

// index returns the index of the tag for the given base, script and region or
// its parent if the tag is not available. If the match is for a parent entry,
// the excess script and region are returned.
func (ts *tagSet) index(base language.Base, scr language.Script, reg language.Region) (int, language.Script, language.Region) ***REMOVED***
	lang := base.String()
	index := -1
	if (scr != language.Script***REMOVED******REMOVED*** || reg != language.Region***REMOVED******REMOVED***) ***REMOVED***
		if scr == zzzz ***REMOVED***
			scr = language.Script***REMOVED******REMOVED***
		***REMOVED***
		if reg == zz ***REMOVED***
			reg = language.Region***REMOVED******REMOVED***
		***REMOVED***

		i := sort.SearchStrings(ts.long, lang)
		// All entries have either a script or a region and not both.
		scrStr, regStr := scr.String(), reg.String()
		for ; i < len(ts.long) && strings.HasPrefix(ts.long[i], lang); i++ ***REMOVED***
			if s := ts.long[i][len(lang)+1:]; s == scrStr ***REMOVED***
				scr = language.Script***REMOVED******REMOVED***
				index = i + ts.single.len()
				break
			***REMOVED*** else if s == regStr ***REMOVED***
				reg = language.Region***REMOVED******REMOVED***
				index = i + ts.single.len()
				break
			***REMOVED***
		***REMOVED***
	***REMOVED***
	if index == -1 ***REMOVED***
		index = ts.single.index(lang)
	***REMOVED***
	return index, scr, reg
***REMOVED***

func (ts *tagSet) Tags() []language.Tag ***REMOVED***
	tags := make([]language.Tag, 0, ts.single.len()+len(ts.long))
	ts.single.keys(func(s string) ***REMOVED***
		tags = append(tags, language.Raw.MustParse(s))
	***REMOVED***)
	for _, s := range ts.long ***REMOVED***
		tags = append(tags, language.Raw.MustParse(s))
	***REMOVED***
	return tags
***REMOVED***

func supportedScripts() []language.Script ***REMOVED***
	scr := make([]language.Script, 0, scriptIndex.len())
	scriptIndex.keys(func(s string) ***REMOVED***
		scr = append(scr, language.MustParseScript(s))
	***REMOVED***)
	return scr
***REMOVED***

func supportedRegions() []language.Region ***REMOVED***
	reg := make([]language.Region, 0, regionIndex.len())
	regionIndex.keys(func(s string) ***REMOVED***
		reg = append(reg, language.MustParseRegion(s))
	***REMOVED***)
	return reg
***REMOVED***

// tagIndex holds a concatenated lists of subtags of length 2 to 4, one string
// for each length, which can be used in combination with binary search to get
// the index associated with a tag.
// For example, a tagIndex***REMOVED***
//   "arenesfrruzh",  // 6 2-byte tags.
//   "barwae",        // 2 3-byte tags.
//   "",
// ***REMOVED***
// would mean that the 2-byte tag "fr" had an index of 3, and the 3-byte tag
// "wae" had an index of 7.
type tagIndex [3]string

func (t *tagIndex) index(s string) int ***REMOVED***
	sz := len(s)
	if sz < 2 || 4 < sz ***REMOVED***
		return -1
	***REMOVED***
	a := t[sz-2]
	index := sort.Search(len(a)/sz, func(i int) bool ***REMOVED***
		p := i * sz
		return a[p:p+sz] >= s
	***REMOVED***)
	p := index * sz
	if end := p + sz; end > len(a) || a[p:end] != s ***REMOVED***
		return -1
	***REMOVED***
	// Add the number of tags for smaller sizes.
	for i := 0; i < sz-2; i++ ***REMOVED***
		index += len(t[i]) / (i + 2)
	***REMOVED***
	return index
***REMOVED***

// len returns the number of tags that are contained in the tagIndex.
func (t *tagIndex) len() (n int) ***REMOVED***
	for i, s := range t ***REMOVED***
		n += len(s) / (i + 2)
	***REMOVED***
	return n
***REMOVED***

// keys calls f for each tag.
func (t *tagIndex) keys(f func(key string)) ***REMOVED***
	for i, s := range *t ***REMOVED***
		for ; s != ""; s = s[i+2:] ***REMOVED***
			f(s[:i+2])
		***REMOVED***
	***REMOVED***
***REMOVED***
