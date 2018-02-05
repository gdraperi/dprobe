// Copyright 2014 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package display

import (
	"fmt"
	"reflect"
	"strings"
	"testing"
	"unicode"

	"golang.org/x/text/internal/testtext"
	"golang.org/x/text/language"
	"golang.org/x/text/message"
)

// TODO: test that tables are properly dropped by the linker for various use
// cases.

var (
	firstLang2aa  = language.MustParseBase("aa")
	lastLang2zu   = language.MustParseBase("zu")
	firstLang3ace = language.MustParseBase("ace")
	lastLang3zza  = language.MustParseBase("zza")
	firstTagAr001 = language.MustParse("ar-001")
	lastTagZhHant = language.MustParse("zh-Hant")
)

// TestValues tests that for all languages, regions, and scripts in Values, at
// least one language has a name defined for it by checking it exists in
// English, which is assumed to be the most comprehensive. It is also tested
// that a Namer returns "" for unsupported values.
func TestValues(t *testing.T) ***REMOVED***
	type testcase struct ***REMOVED***
		kind string
		n    Namer
	***REMOVED***
	// checkDefined checks that a value exists in a Namer.
	checkDefined := func(x interface***REMOVED******REMOVED***, namers []testcase) ***REMOVED***
		for _, n := range namers ***REMOVED***
			testtext.Run(t, fmt.Sprintf("%s.Name(%s)", n.kind, x), func(t *testing.T) ***REMOVED***
				if n.n.Name(x) == "" ***REMOVED***
					// As of version 28 there is no data for az-Arab in English,
					// although there is useful data in other languages.
					if x.(fmt.Stringer).String() == "az-Arab" ***REMOVED***
						return
					***REMOVED***
					t.Errorf("supported but no result")
				***REMOVED***
			***REMOVED***)
		***REMOVED***
	***REMOVED***
	// checkUnsupported checks that a value does not exist in a Namer.
	checkUnsupported := func(x interface***REMOVED******REMOVED***, namers []testcase) ***REMOVED***
		for _, n := range namers ***REMOVED***
			if got := n.n.Name(x); got != "" ***REMOVED***
				t.Fatalf("%s.Name(%s): unsupported tag gave non-empty result: %q", n.kind, x, got)
			***REMOVED***
		***REMOVED***
	***REMOVED***

	tags := map[language.Tag]bool***REMOVED******REMOVED***
	namers := []testcase***REMOVED***
		***REMOVED***"Languages(en)", Languages(language.English)***REMOVED***,
		***REMOVED***"Tags(en)", Tags(language.English)***REMOVED***,
		***REMOVED***"English.Languages()", English.Languages()***REMOVED***,
		***REMOVED***"English.Tags()", English.Tags()***REMOVED***,
	***REMOVED***
	for _, tag := range Values.Tags() ***REMOVED***
		checkDefined(tag, namers)
		tags[tag] = true
	***REMOVED***
	for _, base := range language.Supported.BaseLanguages() ***REMOVED***
		tag, _ := language.All.Compose(base)
		if !tags[tag] ***REMOVED***
			checkUnsupported(tag, namers)
		***REMOVED***
	***REMOVED***

	regions := map[language.Region]bool***REMOVED******REMOVED***
	namers = []testcase***REMOVED***
		***REMOVED***"Regions(en)", Regions(language.English)***REMOVED***,
		***REMOVED***"English.Regions()", English.Regions()***REMOVED***,
	***REMOVED***
	for _, r := range Values.Regions() ***REMOVED***
		checkDefined(r, namers)
		regions[r] = true
	***REMOVED***
	for _, r := range language.Supported.Regions() ***REMOVED***
		if r = r.Canonicalize(); !regions[r] ***REMOVED***
			checkUnsupported(r, namers)
		***REMOVED***
	***REMOVED***

	scripts := map[language.Script]bool***REMOVED******REMOVED***
	namers = []testcase***REMOVED***
		***REMOVED***"Scripts(en)", Scripts(language.English)***REMOVED***,
		***REMOVED***"English.Scripts()", English.Scripts()***REMOVED***,
	***REMOVED***
	for _, s := range Values.Scripts() ***REMOVED***
		checkDefined(s, namers)
		scripts[s] = true
	***REMOVED***
	for _, s := range language.Supported.Scripts() ***REMOVED***
		// Canonicalize the script.
		tag, _ := language.DeprecatedScript.Compose(s)
		if _, s, _ = tag.Raw(); !scripts[s] ***REMOVED***
			checkUnsupported(s, namers)
		***REMOVED***
	***REMOVED***
***REMOVED***

// TestSupported tests that we have at least some Namers for languages that we
// claim to support. To test the claims in the documentation, it also verifies
// that if a Namer is returned, it will have at least some data.
func TestSupported(t *testing.T) ***REMOVED***
	supportedTags := Supported.Tags()
	if len(supportedTags) != numSupported ***REMOVED***
		t.Errorf("number of supported was %d; want %d", len(supportedTags), numSupported)
	***REMOVED***

	namerFuncs := []struct ***REMOVED***
		kind string
		fn   func(language.Tag) Namer
	***REMOVED******REMOVED***
		***REMOVED***"Tags", Tags***REMOVED***,
		***REMOVED***"Languages", Languages***REMOVED***,
		***REMOVED***"Regions", Regions***REMOVED***,
		***REMOVED***"Scripts", Scripts***REMOVED***,
	***REMOVED***

	// Verify that we have at least one Namer for all tags we claim to support.
	tags := make(map[language.Tag]bool)
	for _, tag := range supportedTags ***REMOVED***
		// Test we have at least one Namer for this supported Tag.
		found := false
		for _, kind := range namerFuncs ***REMOVED***
			if defined(t, kind.kind, kind.fn(tag), tag) ***REMOVED***
				found = true
			***REMOVED***
		***REMOVED***
		if !found ***REMOVED***
			t.Errorf("%s: supported, but no data available", tag)
		***REMOVED***
		if tags[tag] ***REMOVED***
			t.Errorf("%s: included in Supported.Tags more than once", tag)
		***REMOVED***
		tags[tag] = true
	***REMOVED***

	// Verify that we have no Namers for tags we don't claim to support.
	for _, base := range language.Supported.BaseLanguages() ***REMOVED***
		tag, _ := language.All.Compose(base)
		// Skip tags that are supported after matching.
		if _, _, conf := matcher.Match(tag); conf != language.No ***REMOVED***
			continue
		***REMOVED***
		// Test there are no Namers for this tag.
		for _, kind := range namerFuncs ***REMOVED***
			if defined(t, kind.kind, kind.fn(tag), tag) ***REMOVED***
				t.Errorf("%[1]s(%[2]s) returns a Namer, but %[2]s is not in the set of supported Tags.", kind.kind, tag)
			***REMOVED***
		***REMOVED***
	***REMOVED***
***REMOVED***

// defined reports whether n is a proper Namer, which means it is non-nil and
// must have at least one non-empty value.
func defined(t *testing.T, kind string, n Namer, tag language.Tag) bool ***REMOVED***
	if n == nil ***REMOVED***
		return false
	***REMOVED***
	switch kind ***REMOVED***
	case "Tags":
		for _, t := range Values.Tags() ***REMOVED***
			if n.Name(t) != "" ***REMOVED***
				return true
			***REMOVED***
		***REMOVED***
	case "Languages":
		for _, t := range Values.BaseLanguages() ***REMOVED***
			if n.Name(t) != "" ***REMOVED***
				return true
			***REMOVED***
		***REMOVED***
	case "Regions":
		for _, t := range Values.Regions() ***REMOVED***
			if n.Name(t) != "" ***REMOVED***
				return true
			***REMOVED***
		***REMOVED***
	case "Scripts":
		for _, t := range Values.Scripts() ***REMOVED***
			if n.Name(t) != "" ***REMOVED***
				return true
			***REMOVED***
		***REMOVED***
	***REMOVED***
	t.Errorf("%s(%s) returns non-nil Namer without content", kind, tag)
	return false
***REMOVED***

func TestCoverage(t *testing.T) ***REMOVED***
	en := language.English
	tests := []struct ***REMOVED***
		n Namer
		x interface***REMOVED******REMOVED***
	***REMOVED******REMOVED***
		***REMOVED***Languages(en), Values.Tags()***REMOVED***,
		***REMOVED***Scripts(en), Values.Scripts()***REMOVED***,
		***REMOVED***Regions(en), Values.Regions()***REMOVED***,
	***REMOVED***
	for i, tt := range tests ***REMOVED***
		uniq := make(map[string]interface***REMOVED******REMOVED***)

		v := reflect.ValueOf(tt.x)
		for j := 0; j < v.Len(); j++ ***REMOVED***
			x := v.Index(j).Interface()
			// As of version 28 there is no data for az-Arab in English,
			// although there is useful data in other languages.
			if x.(fmt.Stringer).String() == "az-Arab" ***REMOVED***
				continue
			***REMOVED***
			s := tt.n.Name(x)
			if s == "" ***REMOVED***
				t.Errorf("%d:%d:%s: missing content", i, j, x)
			***REMOVED*** else if uniq[s] != nil ***REMOVED***
				t.Errorf("%d:%d:%s: identical return value %q for %v and %v", i, j, x, s, x, uniq[s])
			***REMOVED***
			uniq[s] = x
		***REMOVED***
	***REMOVED***
***REMOVED***

// TestUpdate tests whether dictionary entries for certain languages need to be
// updated. For some languages, some of the headers may be empty or they may be
// identical to the parent. This code detects if such entries need to be updated
// after a table update.
func TestUpdate(t *testing.T) ***REMOVED***
	tests := []struct ***REMOVED***
		d   *Dictionary
		tag string
	***REMOVED******REMOVED***
		***REMOVED***ModernStandardArabic, "ar-001"***REMOVED***,
		***REMOVED***AmericanEnglish, "en-US"***REMOVED***,
		***REMOVED***EuropeanSpanish, "es-ES"***REMOVED***,
		***REMOVED***BrazilianPortuguese, "pt-BR"***REMOVED***,
		***REMOVED***SimplifiedChinese, "zh-Hans"***REMOVED***,
	***REMOVED***

	for _, tt := range tests ***REMOVED***
		_, i, _ := matcher.Match(language.MustParse(tt.tag))
		if !reflect.DeepEqual(tt.d.lang, langHeaders[i]) ***REMOVED***
			t.Errorf("%s: lang table update needed", tt.tag)
		***REMOVED***
		if !reflect.DeepEqual(tt.d.script, scriptHeaders[i]) ***REMOVED***
			t.Errorf("%s: script table update needed", tt.tag)
		***REMOVED***
		if !reflect.DeepEqual(tt.d.region, regionHeaders[i]) ***REMOVED***
			t.Errorf("%s: region table update needed", tt.tag)
		***REMOVED***
	***REMOVED***
***REMOVED***

func TestIndex(t *testing.T) ***REMOVED***
	notIn := []string***REMOVED***"aa", "xx", "zz", "aaa", "xxx", "zzz", "Aaaa", "Xxxx", "Zzzz"***REMOVED***
	tests := []tagIndex***REMOVED***
		***REMOVED***
			"",
			"",
			"",
		***REMOVED***,
		***REMOVED***
			"bb",
			"",
			"",
		***REMOVED***,
		***REMOVED***
			"",
			"bbb",
			"",
		***REMOVED***,
		***REMOVED***
			"",
			"",
			"Bbbb",
		***REMOVED***,
		***REMOVED***
			"bb",
			"bbb",
			"Bbbb",
		***REMOVED***,
		***REMOVED***
			"bbccddyy",
			"bbbcccdddyyy",
			"BbbbCcccDdddYyyy",
		***REMOVED***,
	***REMOVED***
	for i, tt := range tests ***REMOVED***
		// Create the test set from the tagIndex.
		cnt := 0
		for sz := 2; sz <= 4; sz++ ***REMOVED***
			a := tt[sz-2]
			for j := 0; j < len(a); j += sz ***REMOVED***
				s := a[j : j+sz]
				if idx := tt.index(s); idx != cnt ***REMOVED***
					t.Errorf("%d:%s: index was %d; want %d", i, s, idx, cnt)
				***REMOVED***
				cnt++
			***REMOVED***
		***REMOVED***
		if n := tt.len(); n != cnt ***REMOVED***
			t.Errorf("%d: len was %d; want %d", i, n, cnt)
		***REMOVED***
		for _, x := range notIn ***REMOVED***
			if idx := tt.index(x); idx != -1 ***REMOVED***
				t.Errorf("%d:%s: index was %d; want -1", i, x, idx)
			***REMOVED***
		***REMOVED***
	***REMOVED***
***REMOVED***

func TestTag(t *testing.T) ***REMOVED***
	tests := []struct ***REMOVED***
		dict string
		tag  string
		name string
	***REMOVED******REMOVED***
		// sr is in Value.Languages(), but is not supported by agq.
		***REMOVED***"agq", "sr", "|[language: sr]"***REMOVED***,
		***REMOVED***"nl", "nl", "Nederlands"***REMOVED***,
		// CLDR 30 dropped Vlaams as the word for nl-BE. It is still called
		// Flemish in English, though. TODO: check if this is a CLDR bug.
		// ***REMOVED***"nl", "nl-BE", "Vlaams"***REMOVED***,
		***REMOVED***"nl", "nl-BE", "Nederlands (België)"***REMOVED***,
		***REMOVED***"nl", "vls", "West-Vlaams"***REMOVED***,
		***REMOVED***"en", "nl-BE", "Flemish"***REMOVED***,
		***REMOVED***"en", "en", "English"***REMOVED***,
		***REMOVED***"en", "en-GB", "British English"***REMOVED***,
		***REMOVED***"en", "en-US", "American English"***REMOVED***, // American English in CLDR 24+
		***REMOVED***"ru", "ru", "русский"***REMOVED***,
		***REMOVED***"ru", "ru-RU", "русский (Россия)"***REMOVED***,
		***REMOVED***"ru", "ru-Cyrl", "русский (кириллица)"***REMOVED***,
		***REMOVED***"en", lastLang2zu.String(), "Zulu"***REMOVED***,
		***REMOVED***"en", firstLang2aa.String(), "Afar"***REMOVED***,
		***REMOVED***"en", lastLang3zza.String(), "Zaza"***REMOVED***,
		***REMOVED***"en", firstLang3ace.String(), "Achinese"***REMOVED***,
		***REMOVED***"en", firstTagAr001.String(), "Modern Standard Arabic"***REMOVED***,
		***REMOVED***"en", lastTagZhHant.String(), "Traditional Chinese"***REMOVED***,
		***REMOVED***"en", "aaa", "|Unknown language (aaa)"***REMOVED***,
		***REMOVED***"en", "zzj", "|Unknown language (zzj)"***REMOVED***,
		// If full tag doesn't match, try without script or region.
		***REMOVED***"en", "aa-Hans", "Afar (Simplified Han)"***REMOVED***,
		***REMOVED***"en", "af-Arab", "Afrikaans (Arabic)"***REMOVED***,
		***REMOVED***"en", "zu-Cyrl", "Zulu (Cyrillic)"***REMOVED***,
		***REMOVED***"en", "aa-GB", "Afar (United Kingdom)"***REMOVED***,
		***REMOVED***"en", "af-NA", "Afrikaans (Namibia)"***REMOVED***,
		***REMOVED***"en", "zu-BR", "Zulu (Brazil)"***REMOVED***,
		// Correct inheritance and language selection.
		***REMOVED***"zh", "zh-TW", "中文 (台湾)"***REMOVED***,
		***REMOVED***"zh", "zh-Hant-TW", "繁体中文 (台湾)"***REMOVED***,
		***REMOVED***"zh-Hant", "zh-TW", "中文 (台灣)"***REMOVED***,
		***REMOVED***"zh-Hant", "zh-Hant-TW", "繁體中文 (台灣)"***REMOVED***,
		// Some rather arbitrary interpretations for Serbian. This is arguably
		// correct and consistent with the way zh-[Hant-]TW is handled. It will
		// also give results more in line with the expectations if users
		// explicitly use "sh".
		***REMOVED***"sr-Latn", "sr-ME", "srpski (Crna Gora)"***REMOVED***,
		***REMOVED***"sr-Latn", "sr-Latn-ME", "srpskohrvatski (Crna Gora)"***REMOVED***,
		// Double script and region
		***REMOVED***"nl", "en-Cyrl-BE", "Engels (Cyrillisch, België)"***REMOVED***,
	***REMOVED***
	for _, tt := range tests ***REMOVED***
		t.Run(tt.dict+"/"+tt.tag, func(t *testing.T) ***REMOVED***
			name, fmtName := splitName(tt.name)
			dict := language.MustParse(tt.dict)
			tag := language.Raw.MustParse(tt.tag)
			d := Tags(dict)
			if n := d.Name(tag); n != name ***REMOVED***
				// There are inconsistencies w.r.t. capitalization in the tests
				// due to CLDR's update procedure which treats modern and other
				// languages differently.
				// See http://unicode.org/cldr/trac/ticket/8051.
				// TODO: use language capitalization to sanitize the strings.
				t.Errorf("Name(%s) = %q; want %q", tag, n, name)
			***REMOVED***

			p := message.NewPrinter(dict)
			if n := p.Sprint(Tag(tag)); n != fmtName ***REMOVED***
				t.Errorf("Tag(%s) = %q; want %q", tag, n, fmtName)
			***REMOVED***
		***REMOVED***)
	***REMOVED***
***REMOVED***

func splitName(names string) (name, formatName string) ***REMOVED***
	split := strings.Split(names, "|")
	name, formatName = split[0], split[0]
	if len(split) > 1 ***REMOVED***
		formatName = split[1]
	***REMOVED***
	return name, formatName
***REMOVED***

func TestLanguage(t *testing.T) ***REMOVED***
	tests := []struct ***REMOVED***
		dict string
		tag  string
		name string
	***REMOVED******REMOVED***
		// sr is in Value.Languages(), but is not supported by agq.
		***REMOVED***"agq", "sr", "|[language: sr]"***REMOVED***,
		// CLDR 30 dropped Vlaams as the word for nl-BE. It is still called
		// Flemish in English, though. TODO: this is probably incorrect.
		// West-Vlaams (vls) is not Vlaams. West-Vlaams could be considered its
		// own language, whereas Vlaams is generally Dutch. So expect to have
		// to change these tests back.
		***REMOVED***"nl", "nl", "Nederlands"***REMOVED***,
		***REMOVED***"nl", "vls", "West-Vlaams"***REMOVED***,
		***REMOVED***"nl", "nl-BE", "Nederlands"***REMOVED***,
		***REMOVED***"en", "pt", "Portuguese"***REMOVED***,
		***REMOVED***"en", "pt-PT", "European Portuguese"***REMOVED***,
		***REMOVED***"en", "pt-BR", "Brazilian Portuguese"***REMOVED***,
		***REMOVED***"en", "en", "English"***REMOVED***,
		***REMOVED***"en", "en-GB", "British English"***REMOVED***,
		***REMOVED***"en", "en-US", "American English"***REMOVED***, // American English in CLDR 24+
		***REMOVED***"en", lastLang2zu.String(), "Zulu"***REMOVED***,
		***REMOVED***"en", firstLang2aa.String(), "Afar"***REMOVED***,
		***REMOVED***"en", lastLang3zza.String(), "Zaza"***REMOVED***,
		***REMOVED***"en", firstLang3ace.String(), "Achinese"***REMOVED***,
		***REMOVED***"en", firstTagAr001.String(), "Modern Standard Arabic"***REMOVED***,
		***REMOVED***"en", lastTagZhHant.String(), "Traditional Chinese"***REMOVED***,
		***REMOVED***"en", "aaa", "|Unknown language (aaa)"***REMOVED***,
		***REMOVED***"en", "zzj", "|Unknown language (zzj)"***REMOVED***,
		// If full tag doesn't match, try without script or region.
		***REMOVED***"en", "aa-Hans", "Afar"***REMOVED***,
		***REMOVED***"en", "af-Arab", "Afrikaans"***REMOVED***,
		***REMOVED***"en", "zu-Cyrl", "Zulu"***REMOVED***,
		***REMOVED***"en", "aa-GB", "Afar"***REMOVED***,
		***REMOVED***"en", "af-NA", "Afrikaans"***REMOVED***,
		***REMOVED***"en", "zu-BR", "Zulu"***REMOVED***,
		***REMOVED***"agq", "zh-Hant", "|[language: zh-Hant]"***REMOVED***,
		***REMOVED***"en", "sh", "Serbo-Croatian"***REMOVED***,
		***REMOVED***"en", "sr-Latn", "Serbo-Croatian"***REMOVED***,
		***REMOVED***"en", "sr", "Serbian"***REMOVED***,
		***REMOVED***"en", "sr-ME", "Serbian"***REMOVED***,
		***REMOVED***"en", "sr-Latn-ME", "Serbo-Croatian"***REMOVED***, // See comments in TestTag.
	***REMOVED***
	for _, tt := range tests ***REMOVED***
		testtext.Run(t, tt.dict+"/"+tt.tag, func(t *testing.T) ***REMOVED***
			name, fmtName := splitName(tt.name)
			dict := language.MustParse(tt.dict)
			tag := language.Raw.MustParse(tt.tag)
			p := message.NewPrinter(dict)
			d := Languages(dict)
			if n := d.Name(tag); n != name ***REMOVED***
				t.Errorf("Name(%v) = %q; want %q", tag, n, name)
			***REMOVED***
			if n := p.Sprint(Language(tag)); n != fmtName ***REMOVED***
				t.Errorf("Language(%v) = %q; want %q", tag, n, fmtName)
			***REMOVED***
			if len(tt.tag) <= 3 ***REMOVED***
				base := language.MustParseBase(tt.tag)
				if n := d.Name(base); n != name ***REMOVED***
					t.Errorf("Name(%v) = %q; want %q", base, n, name)
				***REMOVED***
				if n := p.Sprint(Language(base)); n != fmtName ***REMOVED***
					t.Errorf("Language(%v) = %q; want %q", base, n, fmtName)
				***REMOVED***
			***REMOVED***
		***REMOVED***)
	***REMOVED***
***REMOVED***

func TestScript(t *testing.T) ***REMOVED***
	tests := []struct ***REMOVED***
		dict string
		scr  string
		name string
	***REMOVED******REMOVED***
		***REMOVED***"nl", "Arab", "Arabisch"***REMOVED***,
		***REMOVED***"en", "Arab", "Arabic"***REMOVED***,
		***REMOVED***"en", "Zzzz", "Unknown Script"***REMOVED***,
		***REMOVED***"zh-Hant", "Hang", "韓文字"***REMOVED***,
		***REMOVED***"zh-Hant-HK", "Hang", "韓文字"***REMOVED***,
		***REMOVED***"zh", "Arab", "阿拉伯文"***REMOVED***,
		***REMOVED***"zh-Hans-HK", "Arab", "阿拉伯文"***REMOVED***, // same as zh
		***REMOVED***"zh-Hant", "Arab", "阿拉伯文"***REMOVED***,
		***REMOVED***"zh-Hant-HK", "Arab", "阿拉伯文"***REMOVED***, // same as zh
		// Canonicalized form
		***REMOVED***"en", "Qaai", "Inherited"***REMOVED***,    // deprecated script, now is Zinh
		***REMOVED***"en", "sh", "Unknown Script"***REMOVED***, // sh canonicalizes to sr-Latn
		***REMOVED***"en", "en", "Unknown Script"***REMOVED***,
		// Don't introduce scripts with canonicalization.
		***REMOVED***"en", "sh", "Unknown Script"***REMOVED***, // sh canonicalizes to sr-Latn
	***REMOVED***
	for _, tt := range tests ***REMOVED***
		t.Run(tt.dict+"/"+tt.scr, func(t *testing.T) ***REMOVED***
			name, fmtName := splitName(tt.name)
			dict := language.MustParse(tt.dict)
			p := message.NewPrinter(dict)
			d := Scripts(dict)
			var tag language.Tag
			if unicode.IsUpper(rune(tt.scr[0])) ***REMOVED***
				x := language.MustParseScript(tt.scr)
				if n := d.Name(x); n != name ***REMOVED***
					t.Errorf("Name(%v) = %q; want %q", x, n, name)
				***REMOVED***
				if n := p.Sprint(Script(x)); n != fmtName ***REMOVED***
					t.Errorf("Script(%v) = %q; want %q", x, n, fmtName)
				***REMOVED***
				tag, _ = language.Raw.Compose(x)
			***REMOVED*** else ***REMOVED***
				tag = language.Raw.MustParse(tt.scr)
			***REMOVED***
			if n := d.Name(tag); n != name ***REMOVED***
				t.Errorf("Name(%v) = %q; want %q", tag, n, name)
			***REMOVED***
			if n := p.Sprint(Script(tag)); n != fmtName ***REMOVED***
				t.Errorf("Script(%v) = %q; want %q", tag, n, fmtName)
			***REMOVED***
		***REMOVED***)
	***REMOVED***
***REMOVED***

func TestRegion(t *testing.T) ***REMOVED***
	tests := []struct ***REMOVED***
		dict string
		reg  string
		name string
	***REMOVED******REMOVED***
		***REMOVED***"nl", "NL", "Nederland"***REMOVED***,
		***REMOVED***"en", "US", "United States"***REMOVED***,
		***REMOVED***"en", "ZZ", "Unknown Region"***REMOVED***,
		***REMOVED***"en-GB", "NL", "Netherlands"***REMOVED***,
		// Canonical equivalents
		***REMOVED***"en", "UK", "United Kingdom"***REMOVED***,
		// No region
		***REMOVED***"en", "pt", "Unknown Region"***REMOVED***,
		***REMOVED***"en", "und", "Unknown Region"***REMOVED***,
		// Don't introduce regions with canonicalization.
		***REMOVED***"en", "mo", "Unknown Region"***REMOVED***,
	***REMOVED***
	for _, tt := range tests ***REMOVED***
		t.Run(tt.dict+"/"+tt.reg, func(t *testing.T) ***REMOVED***
			dict := language.MustParse(tt.dict)
			p := message.NewPrinter(dict)
			d := Regions(dict)
			var tag language.Tag
			if unicode.IsUpper(rune(tt.reg[0])) ***REMOVED***
				// Region
				x := language.MustParseRegion(tt.reg)
				if n := d.Name(x); n != tt.name ***REMOVED***
					t.Errorf("Name(%v) = %q; want %q", x, n, tt.name)
				***REMOVED***
				if n := p.Sprint(Region(x)); n != tt.name ***REMOVED***
					t.Errorf("Region(%v) = %q; want %q", x, n, tt.name)
				***REMOVED***
				tag, _ = language.Raw.Compose(x)
			***REMOVED*** else ***REMOVED***
				tag = language.Raw.MustParse(tt.reg)
			***REMOVED***
			if n := d.Name(tag); n != tt.name ***REMOVED***
				t.Errorf("Name(%v) = %q; want %q", tag, n, tt.name)
			***REMOVED***
			if n := p.Sprint(Region(tag)); n != tt.name ***REMOVED***
				t.Errorf("Region(%v) = %q; want %q", tag, n, tt.name)
			***REMOVED***
		***REMOVED***)
	***REMOVED***
***REMOVED***

func TestSelf(t *testing.T) ***REMOVED***
	tests := []struct ***REMOVED***
		tag  string
		name string
	***REMOVED******REMOVED***
		***REMOVED***"nl", "Nederlands"***REMOVED***,
		// CLDR 30 dropped Vlaams as the word for nl-BE. It is still called
		// Flemish in English, though. TODO: check if this is a CLDR bug.
		// ***REMOVED***"nl-BE", "Vlaams"***REMOVED***,
		***REMOVED***"nl-BE", "Nederlands"***REMOVED***,
		***REMOVED***"en-GB", "British English"***REMOVED***,
		***REMOVED***lastLang2zu.String(), "isiZulu"***REMOVED***,
		***REMOVED***firstLang2aa.String(), ""***REMOVED***,  // not defined
		***REMOVED***lastLang3zza.String(), ""***REMOVED***,  // not defined
		***REMOVED***firstLang3ace.String(), ""***REMOVED***, // not defined
		***REMOVED***firstTagAr001.String(), "العربية الرسمية الحديثة"***REMOVED***,
		***REMOVED***"ar", "العربية"***REMOVED***,
		***REMOVED***lastTagZhHant.String(), "繁體中文"***REMOVED***,
		***REMOVED***"aaa", ""***REMOVED***,
		***REMOVED***"zzj", ""***REMOVED***,
		// Drop entries that are not in the requested script, even if there is
		// an entry for the language.
		***REMOVED***"aa-Hans", ""***REMOVED***,
		***REMOVED***"af-Arab", ""***REMOVED***,
		***REMOVED***"zu-Cyrl", ""***REMOVED***,
		// Append the country name in the language of the matching language.
		***REMOVED***"af-NA", "Afrikaans"***REMOVED***,
		***REMOVED***"zh", "中文"***REMOVED***,
		// zh-TW should match zh-Hant instead of zh!
		***REMOVED***"zh-TW", "繁體中文"***REMOVED***,
		***REMOVED***"zh-Hant", "繁體中文"***REMOVED***,
		***REMOVED***"zh-Hans", "简体中文"***REMOVED***,
		***REMOVED***"zh-Hant-TW", "繁體中文"***REMOVED***,
		***REMOVED***"zh-Hans-TW", "简体中文"***REMOVED***,
		// Take the entry for sr which has the matching script.
		// TODO: Capitalization changed as of CLDR 26, but change seems
		// arbitrary. Revisit capitalization with revision 27. See
		// http://unicode.org/cldr/trac/ticket/8051.
		***REMOVED***"sr", "српски"***REMOVED***,
		// TODO: sr-ME should show up as Serbian or Montenegrin, not Serbo-
		// Croatian. This is an artifact of the current algorithm, which is the
		// way it is to have the preferred behavior for other languages such as
		// Chinese. We can hardwire this case in the table generator or package
		// code, but we first check if CLDR can be updated.
		// ***REMOVED***"sr-ME", "Srpski"***REMOVED***, // Is Srpskohrvatski
		***REMOVED***"sr-Latn-ME", "srpskohrvatski"***REMOVED***,
		***REMOVED***"sr-Cyrl-ME", "српски"***REMOVED***,
		***REMOVED***"sr-NL", "српски"***REMOVED***,
		// NOTE: kk is defined, but in Cyrillic script. For China, Arab is the
		// dominant script. We do not have data for kk-Arab and we chose to not
		// fall back in such cases.
		***REMOVED***"kk-CN", ""***REMOVED***,
	***REMOVED***
	for i, tt := range tests ***REMOVED***
		d := Self
		if n := d.Name(language.Raw.MustParse(tt.tag)); n != tt.name ***REMOVED***
			t.Errorf("%d:%s: was %q; want %q", i, tt.tag, n, tt.name)
		***REMOVED***
	***REMOVED***
***REMOVED***

func TestEquivalence(t *testing.T) ***REMOVED***
	testCases := []struct ***REMOVED***
		desc  string
		namer Namer
	***REMOVED******REMOVED***
		***REMOVED***"Self", Self***REMOVED***,
		***REMOVED***"Tags", Tags(language.Romanian)***REMOVED***,
		***REMOVED***"Languages", Languages(language.Romanian)***REMOVED***,
		***REMOVED***"Scripts", Scripts(language.Romanian)***REMOVED***,
	***REMOVED***
	for _, tc := range testCases ***REMOVED***
		t.Run(tc.desc, func(t *testing.T) ***REMOVED***
			ro := tc.namer.Name(language.Raw.MustParse("ro-MD"))
			mo := tc.namer.Name(language.Raw.MustParse("mo"))
			if ro != mo ***REMOVED***
				t.Errorf("%q != %q", ro, mo)
			***REMOVED***
		***REMOVED***)
	***REMOVED***
***REMOVED***

func TestDictionaryLang(t *testing.T) ***REMOVED***
	tests := []struct ***REMOVED***
		d    *Dictionary
		tag  string
		name string
	***REMOVED******REMOVED***
		***REMOVED***English, "en", "English"***REMOVED***,
		***REMOVED***Portuguese, "af", "africâner"***REMOVED***,
		***REMOVED***EuropeanPortuguese, "af", "africanês"***REMOVED***,
		***REMOVED***English, "nl-BE", "Flemish"***REMOVED***,
	***REMOVED***
	for i, test := range tests ***REMOVED***
		tag := language.MustParse(test.tag)
		if got := test.d.Tags().Name(tag); got != test.name ***REMOVED***
			t.Errorf("%d:%v: got %s; want %s", i, tag, got, test.name)
		***REMOVED***
		if base, _ := language.Compose(tag.Base()); base == tag ***REMOVED***
			if got := test.d.Languages().Name(base); got != test.name ***REMOVED***
				t.Errorf("%d:%v: got %s; want %s", i, tag, got, test.name)
			***REMOVED***
		***REMOVED***
	***REMOVED***
***REMOVED***

func TestDictionaryRegion(t *testing.T) ***REMOVED***
	tests := []struct ***REMOVED***
		d      *Dictionary
		region string
		name   string
	***REMOVED******REMOVED***
		***REMOVED***English, "FR", "France"***REMOVED***,
		***REMOVED***Portuguese, "009", "Oceania"***REMOVED***,
		***REMOVED***EuropeanPortuguese, "009", "Oceânia"***REMOVED***,
	***REMOVED***
	for i, test := range tests ***REMOVED***
		tag := language.MustParseRegion(test.region)
		if got := test.d.Regions().Name(tag); got != test.name ***REMOVED***
			t.Errorf("%d:%v: got %s; want %s", i, tag, got, test.name)
		***REMOVED***
	***REMOVED***
***REMOVED***

func TestDictionaryScript(t *testing.T) ***REMOVED***
	tests := []struct ***REMOVED***
		d      *Dictionary
		script string
		name   string
	***REMOVED******REMOVED***
		***REMOVED***English, "Cyrl", "Cyrillic"***REMOVED***,
		***REMOVED***EuropeanPortuguese, "Gujr", "guzerate"***REMOVED***,
	***REMOVED***
	for i, test := range tests ***REMOVED***
		tag := language.MustParseScript(test.script)
		if got := test.d.Scripts().Name(tag); got != test.name ***REMOVED***
			t.Errorf("%d:%v: got %s; want %s", i, tag, got, test.name)
		***REMOVED***
	***REMOVED***
***REMOVED***
