// Copyright 2014 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:generate go run maketables.go -output tables.go

// Package display provides display names for languages, scripts and regions in
// a requested language.
//
// The data is based on CLDR's localeDisplayNames. It includes the names of the
// draft level "contributed" or "approved". The resulting tables are quite
// large. The display package is designed so that users can reduce the linked-in
// table sizes by cherry picking the languages one wishes to support. There is a
// Dictionary defined for a selected set of common languages for this purpose.
package display // import "golang.org/x/text/language/display"

import (
	"fmt"
	"strings"

	"golang.org/x/text/internal/format"
	"golang.org/x/text/language"
)

/*
TODO:
All fairly low priority at the moment:
  - Include alternative and variants as an option (using func options).
  - Option for returning the empty string for undefined values.
  - Support variants, currencies, time zones, option names and other data
    provided in CLDR.
  - Do various optimizations:
    - Reduce size of offset tables.
    - Consider compressing infrequently used languages and decompress on demand.
*/

// A Formatter formats a tag in the current language. It is used in conjunction
// with the message package.
type Formatter struct ***REMOVED***
	lookup func(tag int, x interface***REMOVED******REMOVED***) string
	x      interface***REMOVED******REMOVED***
***REMOVED***

// Format implements "golang.org/x/text/internal/format".Formatter.
func (f Formatter) Format(state format.State, verb rune) ***REMOVED***
	// TODO: there are a lot of inefficiencies in this code. Fix it when we
	// language.Tag has embedded compact tags.
	t := state.Language()
	_, index, _ := matcher.Match(t)
	str := f.lookup(index, f.x)
	if str == "" ***REMOVED***
		// TODO: use language-specific punctuation.
		// TODO: use codePattern instead of language?
		if unknown := f.lookup(index, language.Und); unknown != "" ***REMOVED***
			fmt.Fprintf(state, "%v (%v)", unknown, f.x)
		***REMOVED*** else ***REMOVED***
			fmt.Fprintf(state, "[language: %v]", f.x)
		***REMOVED***
	***REMOVED*** else ***REMOVED***
		state.Write([]byte(str))
	***REMOVED***
***REMOVED***

// Language returns a Formatter that renders the name for lang in the
// the current language. x may be a language.Base or a language.Tag.
// It renders lang in the default language if no translation for the current
// language is supported.
func Language(lang interface***REMOVED******REMOVED***) Formatter ***REMOVED***
	return Formatter***REMOVED***langFunc, lang***REMOVED***
***REMOVED***

// Region returns a Formatter that renders the name for region in the current
// language. region may be a language.Region or a language.Tag.
// It renders region in the default language if no translation for the current
// language is supported.
func Region(region interface***REMOVED******REMOVED***) Formatter ***REMOVED***
	return Formatter***REMOVED***regionFunc, region***REMOVED***
***REMOVED***

// Script returns a Formatter that renders the name for script in the current
// language. script may be a language.Script or a language.Tag.
// It renders script in the default language if no translation for the current
// language is supported.
func Script(script interface***REMOVED******REMOVED***) Formatter ***REMOVED***
	return Formatter***REMOVED***scriptFunc, script***REMOVED***
***REMOVED***

// Script returns a Formatter that renders the name for tag in the current
// language. tag may be a language.Tag.
// It renders tag in the default language if no translation for the current
// language is supported.
func Tag(tag interface***REMOVED******REMOVED***) Formatter ***REMOVED***
	return Formatter***REMOVED***tagFunc, tag***REMOVED***
***REMOVED***

// A Namer is used to get the name for a given value, such as a Tag, Language,
// Script or Region.
type Namer interface ***REMOVED***
	// Name returns a display string for the given value. A Namer returns an
	// empty string for values it does not support. A Namer may support naming
	// an unspecified value. For example, when getting the name for a region for
	// a tag that does not have a defined Region, it may return the name for an
	// unknown region. It is up to the user to filter calls to Name for values
	// for which one does not want to have a name string.
	Name(x interface***REMOVED******REMOVED***) string
***REMOVED***

var (
	// Supported lists the languages for which names are defined.
	Supported language.Coverage

	// The set of all possible values for which names are defined. Note that not
	// all Namer implementations will cover all the values of a given type.
	// A Namer will return the empty string for unsupported values.
	Values language.Coverage

	matcher language.Matcher
)

func init() ***REMOVED***
	tags := make([]language.Tag, numSupported)
	s := supported
	for i := range tags ***REMOVED***
		p := strings.IndexByte(s, '|')
		tags[i] = language.Raw.Make(s[:p])
		s = s[p+1:]
	***REMOVED***
	matcher = language.NewMatcher(tags)
	Supported = language.NewCoverage(tags)

	Values = language.NewCoverage(langTagSet.Tags, supportedScripts, supportedRegions)
***REMOVED***

// Languages returns a Namer for naming languages. It returns nil if there is no
// data for the given tag. The type passed to Name must be either language.Base
// or language.Tag. Note that the result may differ between passing a tag or its
// base language. For example, for English, passing "nl-BE" would return Flemish
// whereas passing "nl" returns "Dutch".
func Languages(t language.Tag) Namer ***REMOVED***
	if _, index, conf := matcher.Match(t); conf != language.No ***REMOVED***
		return languageNamer(index)
	***REMOVED***
	return nil
***REMOVED***

type languageNamer int

func langFunc(i int, x interface***REMOVED******REMOVED***) string ***REMOVED***
	return nameLanguage(languageNamer(i), x)
***REMOVED***

func (n languageNamer) name(i int) string ***REMOVED***
	return lookup(langHeaders[:], int(n), i)
***REMOVED***

// Name implements the Namer interface for language names.
func (n languageNamer) Name(x interface***REMOVED******REMOVED***) string ***REMOVED***
	return nameLanguage(n, x)
***REMOVED***

// nonEmptyIndex walks up the parent chain until a non-empty header is found.
// It returns -1 if no index could be found.
func nonEmptyIndex(h []header, index int) int ***REMOVED***
	for ; index != -1 && h[index].data == ""; index = int(parents[index]) ***REMOVED***
	***REMOVED***
	return index
***REMOVED***

// Scripts returns a Namer for naming scripts. It returns nil if there is no
// data for the given tag. The type passed to Name must be either a
// language.Script or a language.Tag. It will not attempt to infer a script for
// tags with an unspecified script.
func Scripts(t language.Tag) Namer ***REMOVED***
	if _, index, conf := matcher.Match(t); conf != language.No ***REMOVED***
		if index = nonEmptyIndex(scriptHeaders[:], index); index != -1 ***REMOVED***
			return scriptNamer(index)
		***REMOVED***
	***REMOVED***
	return nil
***REMOVED***

type scriptNamer int

func scriptFunc(i int, x interface***REMOVED******REMOVED***) string ***REMOVED***
	return nameScript(scriptNamer(i), x)
***REMOVED***

func (n scriptNamer) name(i int) string ***REMOVED***
	return lookup(scriptHeaders[:], int(n), i)
***REMOVED***

// Name implements the Namer interface for script names.
func (n scriptNamer) Name(x interface***REMOVED******REMOVED***) string ***REMOVED***
	return nameScript(n, x)
***REMOVED***

// Regions returns a Namer for naming regions. It returns nil if there is no
// data for the given tag. The type passed to Name must be either a
// language.Region or a language.Tag. It will not attempt to infer a region for
// tags with an unspecified region.
func Regions(t language.Tag) Namer ***REMOVED***
	if _, index, conf := matcher.Match(t); conf != language.No ***REMOVED***
		if index = nonEmptyIndex(regionHeaders[:], index); index != -1 ***REMOVED***
			return regionNamer(index)
		***REMOVED***
	***REMOVED***
	return nil
***REMOVED***

type regionNamer int

func regionFunc(i int, x interface***REMOVED******REMOVED***) string ***REMOVED***
	return nameRegion(regionNamer(i), x)
***REMOVED***

func (n regionNamer) name(i int) string ***REMOVED***
	return lookup(regionHeaders[:], int(n), i)
***REMOVED***

// Name implements the Namer interface for region names.
func (n regionNamer) Name(x interface***REMOVED******REMOVED***) string ***REMOVED***
	return nameRegion(n, x)
***REMOVED***

// Tags returns a Namer for giving a full description of a tag. The names of
// scripts and regions that are not already implied by the language name will
// in appended within parentheses. It returns nil if there is not data for the
// given tag. The type passed to Name must be a tag.
func Tags(t language.Tag) Namer ***REMOVED***
	if _, index, conf := matcher.Match(t); conf != language.No ***REMOVED***
		return tagNamer(index)
	***REMOVED***
	return nil
***REMOVED***

type tagNamer int

func tagFunc(i int, x interface***REMOVED******REMOVED***) string ***REMOVED***
	return nameTag(languageNamer(i), scriptNamer(i), regionNamer(i), x)
***REMOVED***

// Name implements the Namer interface for tag names.
func (n tagNamer) Name(x interface***REMOVED******REMOVED***) string ***REMOVED***
	return nameTag(languageNamer(n), scriptNamer(n), regionNamer(n), x)
***REMOVED***

// lookup finds the name for an entry in a global table, traversing the
// inheritance hierarchy if needed.
func lookup(table []header, dict, want int) string ***REMOVED***
	for dict != -1 ***REMOVED***
		if s := table[dict].name(want); s != "" ***REMOVED***
			return s
		***REMOVED***
		dict = int(parents[dict])
	***REMOVED***
	return ""
***REMOVED***

// A Dictionary holds a collection of Namers for a single language. One can
// reduce the amount of data linked in to a binary by only referencing
// Dictionaries for the languages one needs to support instead of using the
// generic Namer factories.
type Dictionary struct ***REMOVED***
	parent *Dictionary
	lang   header
	script header
	region header
***REMOVED***

// Tags returns a Namer for giving a full description of a tag. The names of
// scripts and regions that are not already implied by the language name will
// in appended within parentheses. It returns nil if there is not data for the
// given tag. The type passed to Name must be a tag.
func (d *Dictionary) Tags() Namer ***REMOVED***
	return dictTags***REMOVED***d***REMOVED***
***REMOVED***

type dictTags struct ***REMOVED***
	d *Dictionary
***REMOVED***

// Name implements the Namer interface for tag names.
func (n dictTags) Name(x interface***REMOVED******REMOVED***) string ***REMOVED***
	return nameTag(dictLanguages***REMOVED***n.d***REMOVED***, dictScripts***REMOVED***n.d***REMOVED***, dictRegions***REMOVED***n.d***REMOVED***, x)
***REMOVED***

// Languages returns a Namer for naming languages. It returns nil if there is no
// data for the given tag. The type passed to Name must be either language.Base
// or language.Tag. Note that the result may differ between passing a tag or its
// base language. For example, for English, passing "nl-BE" would return Flemish
// whereas passing "nl" returns "Dutch".
func (d *Dictionary) Languages() Namer ***REMOVED***
	return dictLanguages***REMOVED***d***REMOVED***
***REMOVED***

type dictLanguages struct ***REMOVED***
	d *Dictionary
***REMOVED***

func (n dictLanguages) name(i int) string ***REMOVED***
	for d := n.d; d != nil; d = d.parent ***REMOVED***
		if s := d.lang.name(i); s != "" ***REMOVED***
			return s
		***REMOVED***
	***REMOVED***
	return ""
***REMOVED***

// Name implements the Namer interface for language names.
func (n dictLanguages) Name(x interface***REMOVED******REMOVED***) string ***REMOVED***
	return nameLanguage(n, x)
***REMOVED***

// Scripts returns a Namer for naming scripts. It returns nil if there is no
// data for the given tag. The type passed to Name must be either a
// language.Script or a language.Tag. It will not attempt to infer a script for
// tags with an unspecified script.
func (d *Dictionary) Scripts() Namer ***REMOVED***
	return dictScripts***REMOVED***d***REMOVED***
***REMOVED***

type dictScripts struct ***REMOVED***
	d *Dictionary
***REMOVED***

func (n dictScripts) name(i int) string ***REMOVED***
	for d := n.d; d != nil; d = d.parent ***REMOVED***
		if s := d.script.name(i); s != "" ***REMOVED***
			return s
		***REMOVED***
	***REMOVED***
	return ""
***REMOVED***

// Name implements the Namer interface for script names.
func (n dictScripts) Name(x interface***REMOVED******REMOVED***) string ***REMOVED***
	return nameScript(n, x)
***REMOVED***

// Regions returns a Namer for naming regions. It returns nil if there is no
// data for the given tag. The type passed to Name must be either a
// language.Region or a language.Tag. It will not attempt to infer a region for
// tags with an unspecified region.
func (d *Dictionary) Regions() Namer ***REMOVED***
	return dictRegions***REMOVED***d***REMOVED***
***REMOVED***

type dictRegions struct ***REMOVED***
	d *Dictionary
***REMOVED***

func (n dictRegions) name(i int) string ***REMOVED***
	for d := n.d; d != nil; d = d.parent ***REMOVED***
		if s := d.region.name(i); s != "" ***REMOVED***
			return s
		***REMOVED***
	***REMOVED***
	return ""
***REMOVED***

// Name implements the Namer interface for region names.
func (n dictRegions) Name(x interface***REMOVED******REMOVED***) string ***REMOVED***
	return nameRegion(n, x)
***REMOVED***

// A SelfNamer implements a Namer that returns the name of language in this same
// language. It provides a very compact mechanism to provide a comprehensive
// list of languages to users in their native language.
type SelfNamer struct ***REMOVED***
	// Supported defines the values supported by this Namer.
	Supported language.Coverage
***REMOVED***

var (
	// Self is a shared instance of a SelfNamer.
	Self *SelfNamer = &self

	self = SelfNamer***REMOVED***language.NewCoverage(selfTagSet.Tags)***REMOVED***
)

// Name returns the name of a given language tag in the language identified by
// this tag. It supports both the language.Base and language.Tag types.
func (n SelfNamer) Name(x interface***REMOVED******REMOVED***) string ***REMOVED***
	t, _ := language.All.Compose(x)
	base, scr, reg := t.Raw()
	baseScript := language.Script***REMOVED******REMOVED***
	if (scr == language.Script***REMOVED******REMOVED*** && reg != language.Region***REMOVED******REMOVED***) ***REMOVED***
		// For looking up in the self dictionary, we need to select the
		// maximized script. This is even the case if the script isn't
		// specified.
		s1, _ := t.Script()
		if baseScript = getScript(base); baseScript != s1 ***REMOVED***
			scr = s1
		***REMOVED***
	***REMOVED***

	i, scr, reg := selfTagSet.index(base, scr, reg)
	if i == -1 ***REMOVED***
		return ""
	***REMOVED***

	// Only return the display name if the script matches the expected script.
	if (scr != language.Script***REMOVED******REMOVED***) ***REMOVED***
		if (baseScript == language.Script***REMOVED******REMOVED***) ***REMOVED***
			baseScript = getScript(base)
		***REMOVED***
		if baseScript != scr ***REMOVED***
			return ""
		***REMOVED***
	***REMOVED***

	return selfHeaders[0].name(i)
***REMOVED***

// getScript returns the maximized script for a base language.
func getScript(b language.Base) language.Script ***REMOVED***
	tag, _ := language.Raw.Compose(b)
	scr, _ := tag.Script()
	return scr
***REMOVED***
