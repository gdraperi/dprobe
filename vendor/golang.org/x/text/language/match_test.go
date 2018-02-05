// Copyright 2013 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package language

import (
	"bytes"
	"flag"
	"fmt"
	"os"
	"path"
	"path/filepath"
	"strings"
	"testing"

	"golang.org/x/text/internal/testtext"
	"golang.org/x/text/internal/ucd"
)

var verbose = flag.Bool("verbose", false, "set to true to print the internal tables of matchers")

func TestCompliance(t *testing.T) ***REMOVED***
	filepath.Walk("testdata", func(file string, info os.FileInfo, err error) error ***REMOVED***
		if info.IsDir() ***REMOVED***
			return nil
		***REMOVED***
		r, err := os.Open(file)
		if err != nil ***REMOVED***
			t.Fatal(err)
		***REMOVED***
		ucd.Parse(r, func(p *ucd.Parser) ***REMOVED***
			name := strings.Replace(path.Join(p.String(0), p.String(1)), " ", "", -1)
			if skip[name] ***REMOVED***
				return
			***REMOVED***
			t.Run(info.Name()+"/"+name, func(t *testing.T) ***REMOVED***
				supported := makeTagList(p.String(0))
				desired := makeTagList(p.String(1))
				gotCombined, index, conf := NewMatcher(supported).Match(desired...)

				gotMatch := supported[index]
				wantMatch := mk(p.String(2))
				if gotMatch != wantMatch ***REMOVED***
					t.Fatalf("match: got %q; want %q (%v)", gotMatch, wantMatch, conf)
				***REMOVED***
				wantCombined, err := Raw.Parse(p.String(3))
				if err == nil && gotCombined != wantCombined ***REMOVED***
					t.Errorf("combined: got %q; want %q (%v)", gotCombined, wantCombined, conf)
				***REMOVED***
			***REMOVED***)
		***REMOVED***)
		return nil
	***REMOVED***)
***REMOVED***

var skip = map[string]bool***REMOVED***
	// TODO: bugs
	// Honor the wildcard match. This may only be useful to select non-exact
	// stuff.
	"mul,af/nl": true, // match: got "af"; want "mul"

	// TODO: include other extensions.
	// combined: got "en-GB-u-ca-buddhist-nu-arab"; want "en-GB-fonipa-t-m0-iso-i0-pinyin-u-ca-buddhist-nu-arab"
	"und,en-GB-u-sd-gbsct/en-fonipa-u-nu-Arab-ca-buddhist-t-m0-iso-i0-pinyin": true,

	// Inconsistencies with Mark Davis' implementation where it is not clear
	// which is better.

	// Inconsistencies in combined. I think the Go approach is more appropriate.
	// We could use -u-rg- and -u-va- as alternative.
	"und,fr/fr-BE-fonipa":              true, // combined: got "fr"; want "fr-BE-fonipa"
	"und,fr-CA/fr-BE-fonipa":           true, // combined: got "fr-CA"; want "fr-BE-fonipa"
	"und,fr-fonupa/fr-BE-fonipa":       true, // combined: got "fr-fonupa"; want "fr-BE-fonipa"
	"und,no/nn-BE-fonipa":              true, // combined: got "no"; want "no-BE-fonipa"
	"50,und,fr-CA-fonupa/fr-BE-fonipa": true, // combined: got "fr-CA-fonupa"; want "fr-BE-fonipa"

	// The initial number is a threshold. As we don't use scoring, we will not
	// implement this.
	"50,und,fr-Cyrl-CA-fonupa/fr-BE-fonipa": true,
	// match: got "und"; want "fr-Cyrl-CA-fonupa"
	// combined: got "und"; want "fr-Cyrl-BE-fonipa"

	// Other interesting cases to test:
	// - Should same language or same script have the preference if there is
	//   usually no understanding of the other script?
	// - More specific region in desired may replace enclosing supported.
***REMOVED***

func makeTagList(s string) (tags []Tag) ***REMOVED***
	for _, s := range strings.Split(s, ",") ***REMOVED***
		tags = append(tags, mk(strings.TrimSpace(s)))
	***REMOVED***
	return tags
***REMOVED***

func TestMatchStrings(t *testing.T) ***REMOVED***
	testCases := []struct ***REMOVED***
		supported string
		desired   string // strings separted by |
		tag       string
		index     int
	***REMOVED******REMOVED******REMOVED***
		supported: "en",
		desired:   "",
		tag:       "en",
		index:     0,
	***REMOVED***, ***REMOVED***
		supported: "en",
		desired:   "nl",
		tag:       "en",
		index:     0,
	***REMOVED***, ***REMOVED***
		supported: "en,nl",
		desired:   "nl",
		tag:       "nl",
		index:     1,
	***REMOVED***, ***REMOVED***
		supported: "en,nl",
		desired:   "nl|en",
		tag:       "nl",
		index:     1,
	***REMOVED***, ***REMOVED***
		supported: "en-GB,nl",
		desired:   "en ; q=0.1,nl",
		tag:       "nl",
		index:     1,
	***REMOVED***, ***REMOVED***
		supported: "en-GB,nl",
		desired:   "en;q=0.005 | dk; q=0.1,nl ",
		tag:       "en-GB",
		index:     0,
	***REMOVED***, ***REMOVED***
		// do not match faulty tags with und
		supported: "en,und",
		desired:   "|en",
		tag:       "en",
		index:     0,
	***REMOVED******REMOVED***
	for _, tc := range testCases ***REMOVED***
		t.Run(path.Join(tc.supported, tc.desired), func(t *testing.T) ***REMOVED***
			m := NewMatcher(makeTagList(tc.supported))
			tag, index := MatchStrings(m, strings.Split(tc.desired, "|")...)
			if tag.String() != tc.tag || index != tc.index ***REMOVED***
				t.Errorf("got %v, %d; want %v, %d", tag, index, tc.tag, tc.index)
			***REMOVED***
		***REMOVED***)
	***REMOVED***
***REMOVED***

func TestAddLikelySubtags(t *testing.T) ***REMOVED***
	tests := []struct***REMOVED*** in, out string ***REMOVED******REMOVED***
		***REMOVED***"aa", "aa-Latn-ET"***REMOVED***,
		***REMOVED***"aa-Latn", "aa-Latn-ET"***REMOVED***,
		***REMOVED***"aa-Arab", "aa-Arab-ET"***REMOVED***,
		***REMOVED***"aa-Arab-ER", "aa-Arab-ER"***REMOVED***,
		***REMOVED***"kk", "kk-Cyrl-KZ"***REMOVED***,
		***REMOVED***"kk-CN", "kk-Arab-CN"***REMOVED***,
		***REMOVED***"cmn", "cmn"***REMOVED***,
		***REMOVED***"zh-AU", "zh-Hant-AU"***REMOVED***,
		***REMOVED***"zh-VN", "zh-Hant-VN"***REMOVED***,
		***REMOVED***"zh-SG", "zh-Hans-SG"***REMOVED***,
		***REMOVED***"zh-Hant", "zh-Hant-TW"***REMOVED***,
		***REMOVED***"zh-Hani", "zh-Hani-CN"***REMOVED***,
		***REMOVED***"und-Hani", "zh-Hani-CN"***REMOVED***,
		***REMOVED***"und", "en-Latn-US"***REMOVED***,
		***REMOVED***"und-GB", "en-Latn-GB"***REMOVED***,
		***REMOVED***"und-CW", "pap-Latn-CW"***REMOVED***,
		***REMOVED***"und-YT", "fr-Latn-YT"***REMOVED***,
		***REMOVED***"und-Arab", "ar-Arab-EG"***REMOVED***,
		***REMOVED***"und-AM", "hy-Armn-AM"***REMOVED***,
		***REMOVED***"und-TW", "zh-Hant-TW"***REMOVED***,
		***REMOVED***"und-002", "en-Latn-NG"***REMOVED***,
		***REMOVED***"und-Latn-002", "en-Latn-NG"***REMOVED***,
		***REMOVED***"en-Latn-002", "en-Latn-NG"***REMOVED***,
		***REMOVED***"en-002", "en-Latn-NG"***REMOVED***,
		***REMOVED***"en-001", "en-Latn-US"***REMOVED***,
		***REMOVED***"und-003", "en-Latn-US"***REMOVED***,
		***REMOVED***"und-GB", "en-Latn-GB"***REMOVED***,
		***REMOVED***"Latn-001", "en-Latn-US"***REMOVED***,
		***REMOVED***"en-001", "en-Latn-US"***REMOVED***,
		***REMOVED***"es-419", "es-Latn-419"***REMOVED***,
		***REMOVED***"he-145", "he-Hebr-IL"***REMOVED***,
		***REMOVED***"ky-145", "ky-Latn-TR"***REMOVED***,
		***REMOVED***"kk", "kk-Cyrl-KZ"***REMOVED***,
		// Don't specialize duplicate and ambiguous matches.
		***REMOVED***"kk-034", "kk-Arab-034"***REMOVED***, // Matches IR and AF. Both are Arab.
		***REMOVED***"ku-145", "ku-Latn-TR"***REMOVED***,  // Matches IQ, TR, and LB, but kk -> TR.
		***REMOVED***"und-Arab-CC", "ms-Arab-CC"***REMOVED***,
		***REMOVED***"und-Arab-GB", "ks-Arab-GB"***REMOVED***,
		***REMOVED***"und-Hans-CC", "zh-Hans-CC"***REMOVED***,
		***REMOVED***"und-CC", "en-Latn-CC"***REMOVED***,
		***REMOVED***"sr", "sr-Cyrl-RS"***REMOVED***,
		***REMOVED***"sr-151", "sr-Latn-151"***REMOVED***, // Matches RO and RU.
		// We would like addLikelySubtags to generate the same results if the input
		// only changes by adding tags that would otherwise have been added
		// by the expansion.
		// In other words:
		//     und-AA -> xx-Scrp-AA   implies und-Scrp-AA -> xx-Scrp-AA
		//     und-AA -> xx-Scrp-AA   implies xx-AA -> xx-Scrp-AA
		//     und-Scrp -> xx-Scrp-AA implies und-Scrp-AA -> xx-Scrp-AA
		//     und-Scrp -> xx-Scrp-AA implies xx-Scrp -> xx-Scrp-AA
		//     xx -> xx-Scrp-AA       implies xx-Scrp -> xx-Scrp-AA
		//     xx -> xx-Scrp-AA       implies xx-AA -> xx-Scrp-AA
		//
		// The algorithm specified in
		//   http://unicode.org/reports/tr35/tr35-9.html#Supplemental_Data,
		// Section C.10, does not handle the first case. For example,
		// the CLDR data contains an entry und-BJ -> fr-Latn-BJ, but not
		// there is no rule for und-Latn-BJ.  According to spec, und-Latn-BJ
		// would expand to en-Latn-BJ, violating the aforementioned principle.
		// We deviate from the spec by letting und-Scrp-AA expand to xx-Scrp-AA
		// if a rule of the form und-AA -> xx-Scrp-AA is defined.
		// Note that as of version 23, CLDR has some explicitly specified
		// entries that do not conform to these rules. The implementation
		// will not correct these explicit inconsistencies. A later versions of CLDR
		// is supposed to fix this.
		***REMOVED***"und-Latn-BJ", "fr-Latn-BJ"***REMOVED***,
		***REMOVED***"und-Bugi-ID", "bug-Bugi-ID"***REMOVED***,
		// regions, scripts and languages without definitions
		***REMOVED***"und-Arab-AA", "ar-Arab-AA"***REMOVED***,
		***REMOVED***"und-Afak-RE", "fr-Afak-RE"***REMOVED***,
		***REMOVED***"und-Arab-GB", "ks-Arab-GB"***REMOVED***,
		***REMOVED***"abp-Arab-GB", "abp-Arab-GB"***REMOVED***,
		// script has preference over region
		***REMOVED***"und-Arab-NL", "ar-Arab-NL"***REMOVED***,
		***REMOVED***"zza", "zza-Latn-TR"***REMOVED***,
		// preserve variants and extensions
		***REMOVED***"de-1901", "de-Latn-DE-1901"***REMOVED***,
		***REMOVED***"de-x-abc", "de-Latn-DE-x-abc"***REMOVED***,
		***REMOVED***"de-1901-x-abc", "de-Latn-DE-1901-x-abc"***REMOVED***,
		***REMOVED***"x-abc", "x-abc"***REMOVED***, // TODO: is this the desired behavior?
	***REMOVED***
	for i, tt := range tests ***REMOVED***
		in, _ := Parse(tt.in)
		out, _ := Parse(tt.out)
		in, _ = in.addLikelySubtags()
		if in.String() != out.String() ***REMOVED***
			t.Errorf("%d: add(%s) was %s; want %s", i, tt.in, in, tt.out)
		***REMOVED***
	***REMOVED***
***REMOVED***
func TestMinimize(t *testing.T) ***REMOVED***
	tests := []struct***REMOVED*** in, out string ***REMOVED******REMOVED***
		***REMOVED***"aa", "aa"***REMOVED***,
		***REMOVED***"aa-Latn", "aa"***REMOVED***,
		***REMOVED***"aa-Latn-ET", "aa"***REMOVED***,
		***REMOVED***"aa-ET", "aa"***REMOVED***,
		***REMOVED***"aa-Arab", "aa-Arab"***REMOVED***,
		***REMOVED***"aa-Arab-ER", "aa-Arab-ER"***REMOVED***,
		***REMOVED***"aa-Arab-ET", "aa-Arab"***REMOVED***,
		***REMOVED***"und", "und"***REMOVED***,
		***REMOVED***"und-Latn", "und"***REMOVED***,
		***REMOVED***"und-Latn-US", "und"***REMOVED***,
		***REMOVED***"en-Latn-US", "en"***REMOVED***,
		***REMOVED***"cmn", "cmn"***REMOVED***,
		***REMOVED***"cmn-Hans", "cmn-Hans"***REMOVED***,
		***REMOVED***"cmn-Hant", "cmn-Hant"***REMOVED***,
		***REMOVED***"zh-AU", "zh-AU"***REMOVED***,
		***REMOVED***"zh-VN", "zh-VN"***REMOVED***,
		***REMOVED***"zh-SG", "zh-SG"***REMOVED***,
		***REMOVED***"zh-Hant", "zh-Hant"***REMOVED***,
		***REMOVED***"zh-Hant-TW", "zh-TW"***REMOVED***,
		***REMOVED***"zh-Hans", "zh"***REMOVED***,
		***REMOVED***"zh-Hani", "zh-Hani"***REMOVED***,
		***REMOVED***"und-Hans", "und-Hans"***REMOVED***,
		***REMOVED***"und-Hani", "und-Hani"***REMOVED***,

		***REMOVED***"und-CW", "und-CW"***REMOVED***,
		***REMOVED***"und-YT", "und-YT"***REMOVED***,
		***REMOVED***"und-Arab", "und-Arab"***REMOVED***,
		***REMOVED***"und-AM", "und-AM"***REMOVED***,
		***REMOVED***"und-Arab-CC", "und-Arab-CC"***REMOVED***,
		***REMOVED***"und-CC", "und-CC"***REMOVED***,
		***REMOVED***"und-Latn-BJ", "und-BJ"***REMOVED***,
		***REMOVED***"und-Bugi-ID", "und-Bugi"***REMOVED***,
		***REMOVED***"bug-Bugi-ID", "bug-Bugi"***REMOVED***,
		// regions, scripts and languages without definitions
		***REMOVED***"und-Arab-AA", "und-Arab-AA"***REMOVED***,
		// preserve variants and extensions
		***REMOVED***"de-Latn-1901", "de-1901"***REMOVED***,
		***REMOVED***"de-Latn-x-abc", "de-x-abc"***REMOVED***,
		***REMOVED***"de-DE-1901-x-abc", "de-1901-x-abc"***REMOVED***,
		***REMOVED***"x-abc", "x-abc"***REMOVED***, // TODO: is this the desired behavior?
	***REMOVED***
	for i, tt := range tests ***REMOVED***
		in, _ := Parse(tt.in)
		out, _ := Parse(tt.out)
		min, _ := in.minimize()
		if min.String() != out.String() ***REMOVED***
			t.Errorf("%d: min(%s) was %s; want %s", i, tt.in, min, tt.out)
		***REMOVED***
		max, _ := min.addLikelySubtags()
		if x, _ := in.addLikelySubtags(); x.String() != max.String() ***REMOVED***
			t.Errorf("%d: max(min(%s)) = %s; want %s", i, tt.in, max, x)
		***REMOVED***
	***REMOVED***
***REMOVED***

func TestRegionGroups(t *testing.T) ***REMOVED***
	testCases := []struct ***REMOVED***
		a, b     string
		distance uint8
	***REMOVED******REMOVED***
		***REMOVED***"zh-TW", "zh-HK", 5***REMOVED***,
		***REMOVED***"zh-MO", "zh-HK", 4***REMOVED***,
		***REMOVED***"es-ES", "es-AR", 5***REMOVED***,
		***REMOVED***"es-ES", "es", 4***REMOVED***,
		***REMOVED***"es-419", "es-MX", 4***REMOVED***,
		***REMOVED***"es-AR", "es-MX", 4***REMOVED***,
		***REMOVED***"es-ES", "es-MX", 5***REMOVED***,
		***REMOVED***"es-PT", "es-MX", 5***REMOVED***,
	***REMOVED***
	for _, tc := range testCases ***REMOVED***
		a := MustParse(tc.a)
		aScript, _ := a.Script()
		b := MustParse(tc.b)
		bScript, _ := b.Script()

		if aScript != bScript ***REMOVED***
			t.Errorf("scripts differ: %q vs %q", aScript, bScript)
			continue
		***REMOVED***
		d, _ := regionGroupDist(a.region, b.region, aScript.scriptID, a.lang)
		if d != tc.distance ***REMOVED***
			t.Errorf("got %q; want %q", d, tc.distance)
		***REMOVED***
	***REMOVED***
***REMOVED***

func TestIsParadigmLocale(t *testing.T) ***REMOVED***
	testCases := map[string]bool***REMOVED***
		"en-US":  true,
		"en-GB":  true,
		"en-VI":  false,
		"es-GB":  false,
		"es-ES":  true,
		"es-419": true,
	***REMOVED***
	for str, want := range testCases ***REMOVED***
		tag := Make(str)
		got := isParadigmLocale(tag.lang, tag.region)
		if got != want ***REMOVED***
			t.Errorf("isPL(%q) = %v; want %v", str, got, want)
		***REMOVED***
	***REMOVED***
***REMOVED***

// Implementation of String methods for various types for debugging purposes.

func (m *matcher) String() string ***REMOVED***
	w := &bytes.Buffer***REMOVED******REMOVED***
	fmt.Fprintln(w, "Default:", m.default_)
	for tag, h := range m.index ***REMOVED***
		fmt.Fprintf(w, "  %s: %v\n", tag, h)
	***REMOVED***
	return w.String()
***REMOVED***

func (h *matchHeader) String() string ***REMOVED***
	w := &bytes.Buffer***REMOVED******REMOVED***
	fmt.Fprint(w, "haveTag: ")
	for _, h := range h.haveTags ***REMOVED***
		fmt.Fprintf(w, "%v, ", h)
	***REMOVED***
	return w.String()
***REMOVED***

func (t haveTag) String() string ***REMOVED***
	return fmt.Sprintf("%v:%d:%v:%v-%v|%v", t.tag, t.index, t.conf, t.maxRegion, t.maxScript, t.altScript)
***REMOVED***

func TestBestMatchAlloc(t *testing.T) ***REMOVED***
	m := NewMatcher(makeTagList("en sr nl"))
	// Go allocates when creating a list of tags from a single tag!
	list := []Tag***REMOVED***English***REMOVED***
	avg := testtext.AllocsPerRun(1, func() ***REMOVED***
		m.Match(list...)
	***REMOVED***)
	if avg > 0 ***REMOVED***
		t.Errorf("got %f; want 0", avg)
	***REMOVED***
***REMOVED***

var benchHave = []Tag***REMOVED***
	mk("en"),
	mk("en-GB"),
	mk("za"),
	mk("zh-Hant"),
	mk("zh-Hans-CN"),
	mk("zh"),
	mk("zh-HK"),
	mk("ar-MK"),
	mk("en-CA"),
	mk("fr-CA"),
	mk("fr-US"),
	mk("fr-CH"),
	mk("fr"),
	mk("lt"),
	mk("lv"),
	mk("iw"),
	mk("iw-NL"),
	mk("he"),
	mk("he-IT"),
	mk("tlh"),
	mk("ja"),
	mk("ja-Jpan"),
	mk("ja-Jpan-JP"),
	mk("de"),
	mk("de-CH"),
	mk("de-AT"),
	mk("de-DE"),
	mk("sr"),
	mk("sr-Latn"),
	mk("sr-Cyrl"),
	mk("sr-ME"),
***REMOVED***

var benchWant = [][]Tag***REMOVED***
	[]Tag***REMOVED***
		mk("en"),
	***REMOVED***,
	[]Tag***REMOVED***
		mk("en-AU"),
		mk("de-HK"),
		mk("nl"),
		mk("fy"),
		mk("lv"),
	***REMOVED***,
	[]Tag***REMOVED***
		mk("en-AU"),
		mk("de-HK"),
		mk("nl"),
		mk("fy"),
	***REMOVED***,
	[]Tag***REMOVED***
		mk("ja-Hant"),
		mk("da-HK"),
		mk("nl"),
		mk("zh-TW"),
	***REMOVED***,
	[]Tag***REMOVED***
		mk("ja-Hant"),
		mk("da-HK"),
		mk("nl"),
		mk("hr"),
	***REMOVED***,
***REMOVED***

func BenchmarkMatch(b *testing.B) ***REMOVED***
	m := newMatcher(benchHave, nil)
	for i := 0; i < b.N; i++ ***REMOVED***
		for _, want := range benchWant ***REMOVED***
			m.getBest(want...)
		***REMOVED***
	***REMOVED***
***REMOVED***

func BenchmarkMatchExact(b *testing.B) ***REMOVED***
	want := mk("en")
	m := newMatcher(benchHave, nil)
	for i := 0; i < b.N; i++ ***REMOVED***
		m.getBest(want)
	***REMOVED***
***REMOVED***

func BenchmarkMatchAltLanguagePresent(b *testing.B) ***REMOVED***
	want := mk("hr")
	m := newMatcher(benchHave, nil)
	for i := 0; i < b.N; i++ ***REMOVED***
		m.getBest(want)
	***REMOVED***
***REMOVED***

func BenchmarkMatchAltLanguageNotPresent(b *testing.B) ***REMOVED***
	want := mk("nn")
	m := newMatcher(benchHave, nil)
	for i := 0; i < b.N; i++ ***REMOVED***
		m.getBest(want)
	***REMOVED***
***REMOVED***

func BenchmarkMatchAltScriptPresent(b *testing.B) ***REMOVED***
	want := mk("zh-Hant-CN")
	m := newMatcher(benchHave, nil)
	for i := 0; i < b.N; i++ ***REMOVED***
		m.getBest(want)
	***REMOVED***
***REMOVED***

func BenchmarkMatchAltScriptNotPresent(b *testing.B) ***REMOVED***
	want := mk("fr-Cyrl")
	m := newMatcher(benchHave, nil)
	for i := 0; i < b.N; i++ ***REMOVED***
		m.getBest(want)
	***REMOVED***
***REMOVED***

func BenchmarkMatchLimitedExact(b *testing.B) ***REMOVED***
	want := []Tag***REMOVED***mk("he-NL"), mk("iw-NL")***REMOVED***
	m := newMatcher(benchHave, nil)
	for i := 0; i < b.N; i++ ***REMOVED***
		m.getBest(want...)
	***REMOVED***
***REMOVED***
