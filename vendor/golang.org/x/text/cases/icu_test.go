// Copyright 2016 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build icu

package cases

import (
	"path"
	"strings"
	"testing"

	"golang.org/x/text/internal/testtext"
	"golang.org/x/text/language"
	"golang.org/x/text/unicode/norm"
)

func TestICUConformance(t *testing.T) ***REMOVED***
	// Build test set.
	input := []string***REMOVED***
		"a.a a_a",
		"a\u05d0a",
		"\u05d0'a",
		"a\u03084a",
		"a\u0308a",
		"a3\u30a3a",
		"a\u303aa",
		"a_\u303a_a",
		"1_a..a",
		"1_a.a",
		"a..a.",
		"a--a-",
		"a-a-",
		"a\u200ba",
		"a\u200b\u200ba",
		"a\u00ad\u00ada", // Format
		"a\u00ada",
		"a''a", // SingleQuote
		"a'a",
		"a::a", // MidLetter
		"a:a",
		"a..a", // MidNumLet
		"a.a",
		"a;;a", // MidNum
		"a;a",
		"a__a", // ExtendNumlet
		"a_a",
		"ΟΣ''a",
	***REMOVED***
	add := func(x interface***REMOVED******REMOVED***) ***REMOVED***
		switch v := x.(type) ***REMOVED***
		case string:
			input = append(input, v)
		case []string:
			for _, s := range v ***REMOVED***
				input = append(input, s)
			***REMOVED***
		***REMOVED***
	***REMOVED***
	for _, tc := range testCases ***REMOVED***
		add(tc.src)
		add(tc.lower)
		add(tc.upper)
		add(tc.title)
	***REMOVED***
	for _, tc := range bufferTests ***REMOVED***
		add(tc.src)
	***REMOVED***
	for _, tc := range breakTest ***REMOVED***
		add(strings.Replace(tc, "|", "", -1))
	***REMOVED***
	for _, tc := range foldTestCases ***REMOVED***
		add(tc)
	***REMOVED***

	// Compare ICU to Go.
	for _, c := range []string***REMOVED***"lower", "upper", "title", "fold"***REMOVED*** ***REMOVED***
		for _, tag := range []string***REMOVED***
			"und", "af", "az", "el", "lt", "nl", "tr",
		***REMOVED*** ***REMOVED***
			for _, s := range input ***REMOVED***
				if exclude(c, tag, s) ***REMOVED***
					continue
				***REMOVED***
				testtext.Run(t, path.Join(c, tag, s), func(t *testing.T) ***REMOVED***
					want := doICU(tag, c, s)
					got := doGo(tag, c, s)
					if norm.NFC.String(got) != norm.NFC.String(want) ***REMOVED***
						t.Errorf("\n    in %[3]q (%+[3]q)\n   got %[1]q (%+[1]q)\n  want %[2]q (%+[2]q)", got, want, s)
					***REMOVED***
				***REMOVED***)
			***REMOVED***
		***REMOVED***
	***REMOVED***
***REMOVED***

// exclude indicates if a string should be excluded from testing.
func exclude(cm, tag, s string) bool ***REMOVED***
	list := []struct***REMOVED*** cm, tags, pattern string ***REMOVED******REMOVED***
		// TODO: Go does not handle certain esoteric breaks correctly. This will be
		// fixed once we have a real word break iterator. Alternatively, it
		// seems like we're not too far off from making it work, so we could
		// fix these last steps. But first verify that using a separate word
		// breaker does not hurt performance.
		***REMOVED***"title", "af nl", "a''a"***REMOVED***,
		***REMOVED***"", "", "א'a"***REMOVED***,

		// All the exclusions below seem to be issues with the ICU
		// implementation (at version 57) and thus are not marked as TODO.

		// ICU does not handle leading apostrophe for Dutch and
		// Afrikaans correctly. See http://unicode.org/cldr/trac/ticket/7078.
		***REMOVED***"title", "af nl", "'n"***REMOVED***,
		***REMOVED***"title", "af nl", "'N"***REMOVED***,

		// Go terminates the final sigma check after a fixed number of
		// ignorables have been found. This ensures that the algorithm can make
		// progress in a streaming scenario.
		***REMOVED***"lower title", "", "\u039f\u03a3...............................a"***REMOVED***,
		// This also applies to upper in Greek.
		// NOTE: we could fix the following two cases by adding state to elUpper
		// and aztrLower. However, considering a modifier to not belong to the
		// preceding letter after the maximum modifiers count is reached is
		// consistent with the behavior of unicode/norm.
		***REMOVED***"upper", "el", "\u03bf" + strings.Repeat("\u0321", 29) + "\u0313"***REMOVED***,
		***REMOVED***"lower", "az tr lt", "I" + strings.Repeat("\u0321", 30) + "\u0307\u0300"***REMOVED***,
		***REMOVED***"upper", "lt", "i" + strings.Repeat("\u0321", 30) + "\u0307\u0300"***REMOVED***,
		***REMOVED***"lower", "lt", "I" + strings.Repeat("\u0321", 30) + "\u0300"***REMOVED***,

		// ICU title case seems to erroneously removes \u0307 from an upper case
		// I unconditionally, instead of only when lowercasing. The ICU
		// transform algorithm transforms these cases consistently with our
		// implementation.
		***REMOVED***"title", "az tr", "\u0307"***REMOVED***,

		// The spec says to remove \u0307 after Soft-Dotted characters. ICU
		// transforms conform but ucasemap_utf8ToUpper does not.
		***REMOVED***"upper title", "lt", "i\u0307"***REMOVED***,
		***REMOVED***"upper title", "lt", "i" + strings.Repeat("\u0321", 29) + "\u0307\u0300"***REMOVED***,

		// Both Unicode and CLDR prescribe an extra explicit dot above after a
		// Soft_Dotted character if there are other modifiers.
		// ucasemap_utf8ToUpper does not do this; ICU transforms do.
		// The issue with ucasemap_utf8ToUpper seems to be that it does not
		// consider the modifiers that are part of composition in the evaluation
		// of More_Above. For instance, according to the More_Above rule for lt,
		// a dotted capital I (U+0130) becomes i\u0307\u0307 (an small i with
		// two additional dots). This seems odd, but is correct. ICU is
		// definitely not correct as it produces different results for different
		// normal forms. For instance, for an İ:
		//    \u0130  (NFC) -> i\u0307         (incorrect)
		//    I\u0307 (NFD) -> i\u0307\u0307   (correct)
		// We could argue that we should not add a \u0307 if there already is
		// one, but this may be hard to get correct and is not conform the
		// standard.
		***REMOVED***"lower title", "lt", "\u0130"***REMOVED***,
		***REMOVED***"lower title", "lt", "\u00cf"***REMOVED***,

		// We are conform ICU ucasemap_utf8ToUpper if we remove support for
		// elUpper. However, this is clearly not conform the spec. Moreover, the
		// ICU transforms _do_ implement this transform and produces results
		// consistent with our implementation. Note that we still prefer to use
		// ucasemap_utf8ToUpper instead of transforms as the latter have
		// inconsistencies in the word breaking algorithm.
		***REMOVED***"upper", "el", "\u0386"***REMOVED***, // GREEK CAPITAL LETTER ALPHA WITH TONOS
		***REMOVED***"upper", "el", "\u0389"***REMOVED***, // GREEK CAPITAL LETTER ETA WITH TONOS
		***REMOVED***"upper", "el", "\u038A"***REMOVED***, // GREEK CAPITAL LETTER IOTA WITH TONOS

		***REMOVED***"upper", "el", "\u0391"***REMOVED***, // GREEK CAPITAL LETTER ALPHA
		***REMOVED***"upper", "el", "\u0397"***REMOVED***, // GREEK CAPITAL LETTER ETA
		***REMOVED***"upper", "el", "\u0399"***REMOVED***, // GREEK CAPITAL LETTER IOTA

		***REMOVED***"upper", "el", "\u03AC"***REMOVED***, // GREEK SMALL LETTER ALPHA WITH TONOS
		***REMOVED***"upper", "el", "\u03AE"***REMOVED***, // GREEK SMALL LETTER ALPHA WITH ETA
		***REMOVED***"upper", "el", "\u03AF"***REMOVED***, // GREEK SMALL LETTER ALPHA WITH IOTA

		***REMOVED***"upper", "el", "\u03B1"***REMOVED***, // GREEK SMALL LETTER ALPHA
		***REMOVED***"upper", "el", "\u03B7"***REMOVED***, // GREEK SMALL LETTER ETA
		***REMOVED***"upper", "el", "\u03B9"***REMOVED***, // GREEK SMALL LETTER IOTA
	***REMOVED***
	for _, x := range list ***REMOVED***
		if x.cm != "" && strings.Index(x.cm, cm) == -1 ***REMOVED***
			continue
		***REMOVED***
		if x.tags != "" && strings.Index(x.tags, tag) == -1 ***REMOVED***
			continue
		***REMOVED***
		if strings.Index(s, x.pattern) != -1 ***REMOVED***
			return true
		***REMOVED***
	***REMOVED***
	return false
***REMOVED***

func doGo(tag, caser, input string) string ***REMOVED***
	var c Caser
	t := language.MustParse(tag)
	switch caser ***REMOVED***
	case "lower":
		c = Lower(t)
	case "upper":
		c = Upper(t)
	case "title":
		c = Title(t)
	case "fold":
		c = Fold()
	***REMOVED***
	return c.String(input)
***REMOVED***
