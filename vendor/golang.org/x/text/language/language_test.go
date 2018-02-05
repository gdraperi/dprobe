// Copyright 2013 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package language

import (
	"reflect"
	"testing"

	"golang.org/x/text/internal/testtext"
)

func TestTagSize(t *testing.T) ***REMOVED***
	id := Tag***REMOVED******REMOVED***
	typ := reflect.TypeOf(id)
	if typ.Size() > 24 ***REMOVED***
		t.Errorf("size of Tag was %d; want 24", typ.Size())
	***REMOVED***
***REMOVED***

func TestIsRoot(t *testing.T) ***REMOVED***
	loc := Tag***REMOVED******REMOVED***
	if !loc.IsRoot() ***REMOVED***
		t.Errorf("unspecified should be root.")
	***REMOVED***
	for i, tt := range parseTests() ***REMOVED***
		loc, _ := Parse(tt.in)
		undef := tt.lang == "und" && tt.script == "" && tt.region == "" && tt.ext == ""
		if loc.IsRoot() != undef ***REMOVED***
			t.Errorf("%d: was %v; want %v", i, loc.IsRoot(), undef)
		***REMOVED***
	***REMOVED***
***REMOVED***

func TestEquality(t *testing.T) ***REMOVED***
	for i, tt := range parseTests()[48:49] ***REMOVED***
		s := tt.in
		tag := Make(s)
		t1 := Make(tag.String())
		if tag != t1 ***REMOVED***
			t.Errorf("%d:%s: equality test 1 failed\n got: %#v\nwant: %#v)", i, s, t1, tag)
		***REMOVED***
		t2, _ := Compose(tag)
		if tag != t2 ***REMOVED***
			t.Errorf("%d:%s: equality test 2 failed\n got: %#v\nwant: %#v", i, s, t2, tag)
		***REMOVED***
	***REMOVED***
***REMOVED***

func TestMakeString(t *testing.T) ***REMOVED***
	tests := []struct***REMOVED*** in, out string ***REMOVED******REMOVED***
		***REMOVED***"und", "und"***REMOVED***,
		***REMOVED***"und", "und-CW"***REMOVED***,
		***REMOVED***"nl", "nl-NL"***REMOVED***,
		***REMOVED***"de-1901", "nl-1901"***REMOVED***,
		***REMOVED***"de-1901", "de-Arab-1901"***REMOVED***,
		***REMOVED***"x-a-b", "de-Arab-x-a-b"***REMOVED***,
		***REMOVED***"x-a-b", "x-a-b"***REMOVED***,
	***REMOVED***
	for i, tt := range tests ***REMOVED***
		id, _ := Parse(tt.in)
		mod, _ := Parse(tt.out)
		id.setTagsFrom(mod)
		for j := 0; j < 2; j++ ***REMOVED***
			id.remakeString()
			if str := id.String(); str != tt.out ***REMOVED***
				t.Errorf("%d:%d: found %s; want %s", i, j, id.String(), tt.out)
			***REMOVED***
		***REMOVED***
		// The bytes to string conversion as used in remakeString
		// occasionally measures as more than one alloc, breaking this test.
		// To alleviate this we set the number of runs to more than 1.
		if n := testtext.AllocsPerRun(8, id.remakeString); n > 1 ***REMOVED***
			t.Errorf("%d: # allocs got %.1f; want <= 1", i, n)
		***REMOVED***
	***REMOVED***
***REMOVED***

func TestCompactIndex(t *testing.T) ***REMOVED***
	tests := []struct ***REMOVED***
		tag   string
		index int
		ok    bool
	***REMOVED******REMOVED***
		// TODO: these values will change with each CLDR update. This issue
		// will be solved if we decide to fix the indexes.
		***REMOVED***"und", 0, true***REMOVED***,
		***REMOVED***"ca-ES-valencia", 1, true***REMOVED***,
		***REMOVED***"ca-ES-valencia-u-va-posix", 0, false***REMOVED***,
		***REMOVED***"ca-ES-valencia-u-co-phonebk", 1, true***REMOVED***,
		***REMOVED***"ca-ES-valencia-u-co-phonebk-va-posix", 0, false***REMOVED***,
		***REMOVED***"x-klingon", 0, false***REMOVED***,
		***REMOVED***"en-US", 232, true***REMOVED***,
		***REMOVED***"en-US-u-va-posix", 2, true***REMOVED***,
		***REMOVED***"en", 136, true***REMOVED***,
		***REMOVED***"en-u-co-phonebk", 136, true***REMOVED***,
		***REMOVED***"en-001", 137, true***REMOVED***,
		***REMOVED***"sh", 0, false***REMOVED***, // We don't normalize.
	***REMOVED***
	for _, tt := range tests ***REMOVED***
		x, ok := CompactIndex(Raw.MustParse(tt.tag))
		if x != tt.index || ok != tt.ok ***REMOVED***
			t.Errorf("%s: got %d, %v; want %d %v", tt.tag, x, ok, tt.index, tt.ok)
		***REMOVED***
	***REMOVED***
***REMOVED***

func TestMarshal(t *testing.T) ***REMOVED***
	testCases := []string***REMOVED***
		// TODO: these values will change with each CLDR update. This issue
		// will be solved if we decide to fix the indexes.
		"und",
		"ca-ES-valencia",
		"ca-ES-valencia-u-va-posix",
		"ca-ES-valencia-u-co-phonebk",
		"ca-ES-valencia-u-co-phonebk-va-posix",
		"x-klingon",
		"en-US",
		"en-US-u-va-posix",
		"en",
		"en-u-co-phonebk",
		"en-001",
		"sh",
	***REMOVED***
	for _, tc := range testCases ***REMOVED***
		var tag Tag
		err := tag.UnmarshalText([]byte(tc))
		if err != nil ***REMOVED***
			t.Errorf("UnmarshalText(%q): unexpected error: %v", tc, err)
		***REMOVED***
		b, err := tag.MarshalText()
		if err != nil ***REMOVED***
			t.Errorf("MarshalText(%q): unexpected error: %v", tc, err)
		***REMOVED***
		if got := string(b); got != tc ***REMOVED***
			t.Errorf("%s: got %q; want %q", tc, got, tc)
		***REMOVED***
	***REMOVED***
***REMOVED***

func TestBase(t *testing.T) ***REMOVED***
	tests := []struct ***REMOVED***
		loc, lang string
		conf      Confidence
	***REMOVED******REMOVED***
		***REMOVED***"und", "en", Low***REMOVED***,
		***REMOVED***"x-abc", "und", No***REMOVED***,
		***REMOVED***"en", "en", Exact***REMOVED***,
		***REMOVED***"und-Cyrl", "ru", High***REMOVED***,
		// If a region is not included, the official language should be English.
		***REMOVED***"und-US", "en", High***REMOVED***,
		// TODO: not-explicitly listed scripts should probably be und, No
		// Modify addTags to return info on how the match was derived.
		// ***REMOVED***"und-Aghb", "und", No***REMOVED***,
	***REMOVED***
	for i, tt := range tests ***REMOVED***
		loc, _ := Parse(tt.loc)
		lang, conf := loc.Base()
		if lang.String() != tt.lang ***REMOVED***
			t.Errorf("%d: language was %s; want %s", i, lang, tt.lang)
		***REMOVED***
		if conf != tt.conf ***REMOVED***
			t.Errorf("%d: confidence was %d; want %d", i, conf, tt.conf)
		***REMOVED***
	***REMOVED***
***REMOVED***

func TestParseBase(t *testing.T) ***REMOVED***
	tests := []struct ***REMOVED***
		in  string
		out string
		ok  bool
	***REMOVED******REMOVED***
		***REMOVED***"en", "en", true***REMOVED***,
		***REMOVED***"EN", "en", true***REMOVED***,
		***REMOVED***"nld", "nl", true***REMOVED***,
		***REMOVED***"dut", "dut", true***REMOVED***,  // bibliographic
		***REMOVED***"aaj", "und", false***REMOVED***, // unknown
		***REMOVED***"qaa", "qaa", true***REMOVED***,
		***REMOVED***"a", "und", false***REMOVED***,
		***REMOVED***"", "und", false***REMOVED***,
		***REMOVED***"aaaa", "und", false***REMOVED***,
	***REMOVED***
	for i, tt := range tests ***REMOVED***
		x, err := ParseBase(tt.in)
		if x.String() != tt.out || err == nil != tt.ok ***REMOVED***
			t.Errorf("%d:%s: was %s, %v; want %s, %v", i, tt.in, x, err == nil, tt.out, tt.ok)
		***REMOVED***
		if y, _, _ := Raw.Make(tt.out).Raw(); x != y ***REMOVED***
			t.Errorf("%d:%s: tag was %s; want %s", i, tt.in, x, y)
		***REMOVED***
	***REMOVED***
***REMOVED***

func TestScript(t *testing.T) ***REMOVED***
	tests := []struct ***REMOVED***
		loc, scr string
		conf     Confidence
	***REMOVED******REMOVED***
		***REMOVED***"und", "Latn", Low***REMOVED***,
		***REMOVED***"en-Latn", "Latn", Exact***REMOVED***,
		***REMOVED***"en", "Latn", High***REMOVED***,
		***REMOVED***"sr", "Cyrl", Low***REMOVED***,
		***REMOVED***"kk", "Cyrl", High***REMOVED***,
		***REMOVED***"kk-CN", "Arab", Low***REMOVED***,
		***REMOVED***"cmn", "Hans", Low***REMOVED***,
		***REMOVED***"ru", "Cyrl", High***REMOVED***,
		***REMOVED***"ru-RU", "Cyrl", High***REMOVED***,
		***REMOVED***"yue", "Hant", Low***REMOVED***,
		***REMOVED***"x-abc", "Zzzz", Low***REMOVED***,
		***REMOVED***"und-zyyy", "Zyyy", Exact***REMOVED***,
	***REMOVED***
	for i, tt := range tests ***REMOVED***
		loc, _ := Parse(tt.loc)
		sc, conf := loc.Script()
		if sc.String() != tt.scr ***REMOVED***
			t.Errorf("%d:%s: script was %s; want %s", i, tt.loc, sc, tt.scr)
		***REMOVED***
		if conf != tt.conf ***REMOVED***
			t.Errorf("%d:%s: confidence was %d; want %d", i, tt.loc, conf, tt.conf)
		***REMOVED***
	***REMOVED***
***REMOVED***

func TestParseScript(t *testing.T) ***REMOVED***
	tests := []struct ***REMOVED***
		in  string
		out string
		ok  bool
	***REMOVED******REMOVED***
		***REMOVED***"Latn", "Latn", true***REMOVED***,
		***REMOVED***"zzzz", "Zzzz", true***REMOVED***,
		***REMOVED***"zyyy", "Zyyy", true***REMOVED***,
		***REMOVED***"Latm", "Zzzz", false***REMOVED***,
		***REMOVED***"Zzz", "Zzzz", false***REMOVED***,
		***REMOVED***"", "Zzzz", false***REMOVED***,
		***REMOVED***"Zzzxx", "Zzzz", false***REMOVED***,
	***REMOVED***
	for i, tt := range tests ***REMOVED***
		x, err := ParseScript(tt.in)
		if x.String() != tt.out || err == nil != tt.ok ***REMOVED***
			t.Errorf("%d:%s: was %s, %v; want %s, %v", i, tt.in, x, err == nil, tt.out, tt.ok)
		***REMOVED***
		if err == nil ***REMOVED***
			if _, y, _ := Raw.Make("und-" + tt.out).Raw(); x != y ***REMOVED***
				t.Errorf("%d:%s: tag was %s; want %s", i, tt.in, x, y)
			***REMOVED***
		***REMOVED***
	***REMOVED***
***REMOVED***

func TestRegion(t *testing.T) ***REMOVED***
	tests := []struct ***REMOVED***
		loc, reg string
		conf     Confidence
	***REMOVED******REMOVED***
		***REMOVED***"und", "US", Low***REMOVED***,
		***REMOVED***"en", "US", Low***REMOVED***,
		***REMOVED***"zh-Hant", "TW", Low***REMOVED***,
		***REMOVED***"en-US", "US", Exact***REMOVED***,
		***REMOVED***"cmn", "CN", Low***REMOVED***,
		***REMOVED***"ru", "RU", Low***REMOVED***,
		***REMOVED***"yue", "HK", Low***REMOVED***,
		***REMOVED***"x-abc", "ZZ", Low***REMOVED***,
	***REMOVED***
	for i, tt := range tests ***REMOVED***
		loc, _ := Raw.Parse(tt.loc)
		reg, conf := loc.Region()
		if reg.String() != tt.reg ***REMOVED***
			t.Errorf("%d:%s: region was %s; want %s", i, tt.loc, reg, tt.reg)
		***REMOVED***
		if conf != tt.conf ***REMOVED***
			t.Errorf("%d:%s: confidence was %d; want %d", i, tt.loc, conf, tt.conf)
		***REMOVED***
	***REMOVED***
***REMOVED***

func TestEncodeM49(t *testing.T) ***REMOVED***
	tests := []struct ***REMOVED***
		m49  int
		code string
		ok   bool
	***REMOVED******REMOVED***
		***REMOVED***1, "001", true***REMOVED***,
		***REMOVED***840, "US", true***REMOVED***,
		***REMOVED***899, "ZZ", false***REMOVED***,
	***REMOVED***
	for i, tt := range tests ***REMOVED***
		if r, err := EncodeM49(tt.m49); r.String() != tt.code || err == nil != tt.ok ***REMOVED***
			t.Errorf("%d:%d: was %s, %v; want %s, %v", i, tt.m49, r, err == nil, tt.code, tt.ok)
		***REMOVED***
	***REMOVED***
	for i := 1; i <= 1000; i++ ***REMOVED***
		if r, err := EncodeM49(i); err == nil && r.M49() == 0 ***REMOVED***
			t.Errorf("%d has no error, but maps to undefined region", i)
		***REMOVED***
	***REMOVED***
***REMOVED***

func TestParseRegion(t *testing.T) ***REMOVED***
	tests := []struct ***REMOVED***
		in  string
		out string
		ok  bool
	***REMOVED******REMOVED***
		***REMOVED***"001", "001", true***REMOVED***,
		***REMOVED***"840", "US", true***REMOVED***,
		***REMOVED***"899", "ZZ", false***REMOVED***,
		***REMOVED***"USA", "US", true***REMOVED***,
		***REMOVED***"US", "US", true***REMOVED***,
		***REMOVED***"BC", "ZZ", false***REMOVED***,
		***REMOVED***"C", "ZZ", false***REMOVED***,
		***REMOVED***"CCCC", "ZZ", false***REMOVED***,
		***REMOVED***"01", "ZZ", false***REMOVED***,
	***REMOVED***
	for i, tt := range tests ***REMOVED***
		r, err := ParseRegion(tt.in)
		if r.String() != tt.out || err == nil != tt.ok ***REMOVED***
			t.Errorf("%d:%s: was %s, %v; want %s, %v", i, tt.in, r, err == nil, tt.out, tt.ok)
		***REMOVED***
		if err == nil ***REMOVED***
			if _, _, y := Raw.Make("und-" + tt.out).Raw(); r != y ***REMOVED***
				t.Errorf("%d:%s: tag was %s; want %s", i, tt.in, r, y)
			***REMOVED***
		***REMOVED***
	***REMOVED***
***REMOVED***

func TestIsCountry(t *testing.T) ***REMOVED***
	tests := []struct ***REMOVED***
		reg     string
		country bool
	***REMOVED******REMOVED***
		***REMOVED***"US", true***REMOVED***,
		***REMOVED***"001", false***REMOVED***,
		***REMOVED***"958", false***REMOVED***,
		***REMOVED***"419", false***REMOVED***,
		***REMOVED***"203", true***REMOVED***,
		***REMOVED***"020", true***REMOVED***,
		***REMOVED***"900", false***REMOVED***,
		***REMOVED***"999", false***REMOVED***,
		***REMOVED***"QO", false***REMOVED***,
		***REMOVED***"EU", false***REMOVED***,
		***REMOVED***"AA", false***REMOVED***,
		***REMOVED***"XK", true***REMOVED***,
	***REMOVED***
	for i, tt := range tests ***REMOVED***
		reg, _ := getRegionID([]byte(tt.reg))
		r := Region***REMOVED***reg***REMOVED***
		if r.IsCountry() != tt.country ***REMOVED***
			t.Errorf("%d: IsCountry(%s) was %v; want %v", i, tt.reg, r.IsCountry(), tt.country)
		***REMOVED***
	***REMOVED***
***REMOVED***

func TestIsGroup(t *testing.T) ***REMOVED***
	tests := []struct ***REMOVED***
		reg   string
		group bool
	***REMOVED******REMOVED***
		***REMOVED***"US", false***REMOVED***,
		***REMOVED***"001", true***REMOVED***,
		***REMOVED***"958", false***REMOVED***,
		***REMOVED***"419", true***REMOVED***,
		***REMOVED***"203", false***REMOVED***,
		***REMOVED***"020", false***REMOVED***,
		***REMOVED***"900", false***REMOVED***,
		***REMOVED***"999", false***REMOVED***,
		***REMOVED***"QO", true***REMOVED***,
		***REMOVED***"EU", true***REMOVED***,
		***REMOVED***"AA", false***REMOVED***,
		***REMOVED***"XK", false***REMOVED***,
	***REMOVED***
	for i, tt := range tests ***REMOVED***
		reg, _ := getRegionID([]byte(tt.reg))
		r := Region***REMOVED***reg***REMOVED***
		if r.IsGroup() != tt.group ***REMOVED***
			t.Errorf("%d: IsGroup(%s) was %v; want %v", i, tt.reg, r.IsGroup(), tt.group)
		***REMOVED***
	***REMOVED***
***REMOVED***

func TestContains(t *testing.T) ***REMOVED***
	tests := []struct ***REMOVED***
		enclosing, contained string
		contains             bool
	***REMOVED******REMOVED***
		// A region contains itself.
		***REMOVED***"US", "US", true***REMOVED***,
		***REMOVED***"001", "001", true***REMOVED***,

		// Direct containment.
		***REMOVED***"001", "002", true***REMOVED***,
		***REMOVED***"039", "XK", true***REMOVED***,
		***REMOVED***"150", "XK", true***REMOVED***,
		***REMOVED***"EU", "AT", true***REMOVED***,
		***REMOVED***"QO", "AQ", true***REMOVED***,

		// Indirect containemnt.
		***REMOVED***"001", "US", true***REMOVED***,
		***REMOVED***"001", "419", true***REMOVED***,
		***REMOVED***"001", "013", true***REMOVED***,

		// No containment.
		***REMOVED***"US", "001", false***REMOVED***,
		***REMOVED***"155", "EU", false***REMOVED***,
	***REMOVED***
	for i, tt := range tests ***REMOVED***
		enc, _ := getRegionID([]byte(tt.enclosing))
		con, _ := getRegionID([]byte(tt.contained))
		r := Region***REMOVED***enc***REMOVED***
		if got := r.Contains(Region***REMOVED***con***REMOVED***); got != tt.contains ***REMOVED***
			t.Errorf("%d: %s.Contains(%s) was %v; want %v", i, tt.enclosing, tt.contained, got, tt.contains)
		***REMOVED***
	***REMOVED***
***REMOVED***

func TestRegionCanonicalize(t *testing.T) ***REMOVED***
	for i, tt := range []struct***REMOVED*** in, out string ***REMOVED******REMOVED***
		***REMOVED***"UK", "GB"***REMOVED***,
		***REMOVED***"TP", "TL"***REMOVED***,
		***REMOVED***"QU", "EU"***REMOVED***,
		***REMOVED***"SU", "SU"***REMOVED***,
		***REMOVED***"VD", "VN"***REMOVED***,
		***REMOVED***"DD", "DE"***REMOVED***,
	***REMOVED*** ***REMOVED***
		r := MustParseRegion(tt.in)
		want := MustParseRegion(tt.out)
		if got := r.Canonicalize(); got != want ***REMOVED***
			t.Errorf("%d: got %v; want %v", i, got, want)
		***REMOVED***
	***REMOVED***
***REMOVED***

func TestRegionTLD(t *testing.T) ***REMOVED***
	for _, tt := range []struct ***REMOVED***
		in, out string
		ok      bool
	***REMOVED******REMOVED***
		***REMOVED***"EH", "EH", true***REMOVED***,
		***REMOVED***"FR", "FR", true***REMOVED***,
		***REMOVED***"TL", "TL", true***REMOVED***,

		// In ccTLD before in ISO.
		***REMOVED***"GG", "GG", true***REMOVED***,

		// Non-standard assignment of ccTLD to ISO code.
		***REMOVED***"GB", "UK", true***REMOVED***,

		// Exceptionally reserved in ISO and valid ccTLD.
		***REMOVED***"UK", "UK", true***REMOVED***,
		***REMOVED***"AC", "AC", true***REMOVED***,
		***REMOVED***"EU", "EU", true***REMOVED***,
		***REMOVED***"SU", "SU", true***REMOVED***,

		// Exceptionally reserved in ISO and invalid ccTLD.
		***REMOVED***"CP", "ZZ", false***REMOVED***,
		***REMOVED***"DG", "ZZ", false***REMOVED***,
		***REMOVED***"EA", "ZZ", false***REMOVED***,
		***REMOVED***"FX", "ZZ", false***REMOVED***,
		***REMOVED***"IC", "ZZ", false***REMOVED***,
		***REMOVED***"TA", "ZZ", false***REMOVED***,

		// Transitionally reserved in ISO (e.g. deprecated) but valid ccTLD as
		// it is still being phased out.
		***REMOVED***"AN", "AN", true***REMOVED***,
		***REMOVED***"TP", "TP", true***REMOVED***,

		// Transitionally reserved in ISO (e.g. deprecated) and invalid ccTLD.
		// Defined in package language as it has a mapping in CLDR.
		***REMOVED***"BU", "ZZ", false***REMOVED***,
		***REMOVED***"CS", "ZZ", false***REMOVED***,
		***REMOVED***"NT", "ZZ", false***REMOVED***,
		***REMOVED***"YU", "ZZ", false***REMOVED***,
		***REMOVED***"ZR", "ZZ", false***REMOVED***,
		// Not defined in package: SF.

		// Indeterminately reserved in ISO.
		// Defined in package language as it has a legacy mapping in CLDR.
		***REMOVED***"DY", "ZZ", false***REMOVED***,
		***REMOVED***"RH", "ZZ", false***REMOVED***,
		***REMOVED***"VD", "ZZ", false***REMOVED***,
		// Not defined in package: EW, FL, JA, LF, PI, RA, RB, RC, RI, RL, RM,
		// RN, RP, WG, WL, WV, and YV.

		// Not assigned in ISO, but legacy definitions in CLDR.
		***REMOVED***"DD", "ZZ", false***REMOVED***,
		***REMOVED***"YD", "ZZ", false***REMOVED***,

		// Normal mappings but somewhat special status in ccTLD.
		***REMOVED***"BL", "BL", true***REMOVED***,
		***REMOVED***"MF", "MF", true***REMOVED***,
		***REMOVED***"BV", "BV", true***REMOVED***,
		***REMOVED***"SJ", "SJ", true***REMOVED***,

		// Have values when normalized, but not as is.
		***REMOVED***"QU", "ZZ", false***REMOVED***,

		// ISO Private Use.
		***REMOVED***"AA", "ZZ", false***REMOVED***,
		***REMOVED***"QM", "ZZ", false***REMOVED***,
		***REMOVED***"QO", "ZZ", false***REMOVED***,
		***REMOVED***"XA", "ZZ", false***REMOVED***,
		***REMOVED***"XK", "ZZ", false***REMOVED***, // Sometimes used for Kosovo, but invalid ccTLD.
	***REMOVED*** ***REMOVED***
		if tt.in == "" ***REMOVED***
			continue
		***REMOVED***

		r := MustParseRegion(tt.in)
		var want Region
		if tt.out != "ZZ" ***REMOVED***
			want = MustParseRegion(tt.out)
		***REMOVED***
		tld, err := r.TLD()
		if got := err == nil; got != tt.ok ***REMOVED***
			t.Errorf("error(%v): got %v; want %v", r, got, tt.ok)
		***REMOVED***
		if tld != want ***REMOVED***
			t.Errorf("TLD(%v): got %v; want %v", r, tld, want)
		***REMOVED***
	***REMOVED***
***REMOVED***

func TestCanonicalize(t *testing.T) ***REMOVED***
	// TODO: do a full test using CLDR data in a separate regression test.
	tests := []struct ***REMOVED***
		in, out string
		option  CanonType
	***REMOVED******REMOVED***
		***REMOVED***"en-Latn", "en", SuppressScript***REMOVED***,
		***REMOVED***"sr-Cyrl", "sr-Cyrl", SuppressScript***REMOVED***,
		***REMOVED***"sh", "sr-Latn", Legacy***REMOVED***,
		***REMOVED***"sh-HR", "sr-Latn-HR", Legacy***REMOVED***,
		***REMOVED***"sh-Cyrl-HR", "sr-Cyrl-HR", Legacy***REMOVED***,
		***REMOVED***"tl", "fil", Legacy***REMOVED***,
		***REMOVED***"no", "no", Legacy***REMOVED***,
		***REMOVED***"no", "nb", Legacy | CLDR***REMOVED***,
		***REMOVED***"cmn", "cmn", Legacy***REMOVED***,
		***REMOVED***"cmn", "zh", Macro***REMOVED***,
		***REMOVED***"cmn-u-co-stroke", "zh-u-co-stroke", Macro***REMOVED***,
		***REMOVED***"yue", "yue", Macro***REMOVED***,
		***REMOVED***"nb", "no", Macro***REMOVED***,
		***REMOVED***"nb", "nb", Macro | CLDR***REMOVED***,
		***REMOVED***"no", "no", Macro***REMOVED***,
		***REMOVED***"no", "no", Macro | CLDR***REMOVED***,
		***REMOVED***"iw", "he", DeprecatedBase***REMOVED***,
		***REMOVED***"iw", "he", Deprecated | CLDR***REMOVED***,
		***REMOVED***"mo", "ro-MD", Deprecated***REMOVED***, // Adopted by CLDR as of version 25.
		***REMOVED***"alb", "sq", Legacy***REMOVED***,       // bibliographic
		***REMOVED***"dut", "nl", Legacy***REMOVED***,       // bibliographic
		// As of CLDR 25, mo is no longer considered a legacy mapping.
		***REMOVED***"mo", "mo", Legacy | CLDR***REMOVED***,
		***REMOVED***"und-AN", "und-AN", Deprecated***REMOVED***,
		***REMOVED***"und-YD", "und-YE", DeprecatedRegion***REMOVED***,
		***REMOVED***"und-YD", "und-YD", DeprecatedBase***REMOVED***,
		***REMOVED***"und-Qaai", "und-Zinh", DeprecatedScript***REMOVED***,
		***REMOVED***"und-Qaai", "und-Qaai", DeprecatedBase***REMOVED***,
		***REMOVED***"drh", "mn", All***REMOVED***, // drh -> khk -> mn
	***REMOVED***
	for i, tt := range tests ***REMOVED***
		in, _ := Raw.Parse(tt.in)
		in, _ = tt.option.Canonicalize(in)
		if in.String() != tt.out ***REMOVED***
			t.Errorf("%d:%s: was %s; want %s", i, tt.in, in.String(), tt.out)
		***REMOVED***
		if int(in.pVariant) > int(in.pExt) || int(in.pExt) > len(in.str) ***REMOVED***
			t.Errorf("%d:%s:offsets %d <= %d <= %d must be true", i, tt.in, in.pVariant, in.pExt, len(in.str))
		***REMOVED***
	***REMOVED***
	// Test idempotence.
	for _, base := range Supported.BaseLanguages() ***REMOVED***
		tag, _ := Raw.Compose(base)
		got, _ := All.Canonicalize(tag)
		want, _ := All.Canonicalize(got)
		if got != want ***REMOVED***
			t.Errorf("idem(%s): got %s; want %s", tag, got, want)
		***REMOVED***
	***REMOVED***
***REMOVED***

func TestTypeForKey(t *testing.T) ***REMOVED***
	tests := []struct***REMOVED*** key, in, out string ***REMOVED******REMOVED***
		***REMOVED***"co", "en", ""***REMOVED***,
		***REMOVED***"co", "en-u-abc", ""***REMOVED***,
		***REMOVED***"co", "en-u-co-phonebk", "phonebk"***REMOVED***,
		***REMOVED***"co", "en-u-co-phonebk-cu-aud", "phonebk"***REMOVED***,
		***REMOVED***"co", "x-foo-u-co-phonebk", ""***REMOVED***,
		***REMOVED***"nu", "en-u-co-phonebk-nu-arabic", "arabic"***REMOVED***,
		***REMOVED***"kc", "cmn-u-co-stroke", ""***REMOVED***,
	***REMOVED***
	for _, tt := range tests ***REMOVED***
		if v := Make(tt.in).TypeForKey(tt.key); v != tt.out ***REMOVED***
			t.Errorf("%q[%q]: was %q; want %q", tt.in, tt.key, v, tt.out)
		***REMOVED***
	***REMOVED***
***REMOVED***

func TestSetTypeForKey(t *testing.T) ***REMOVED***
	tests := []struct ***REMOVED***
		key, value, in, out string
		err                 bool
	***REMOVED******REMOVED***
		// replace existing value
		***REMOVED***"co", "pinyin", "en-u-co-phonebk", "en-u-co-pinyin", false***REMOVED***,
		***REMOVED***"co", "pinyin", "en-u-co-phonebk-cu-xau", "en-u-co-pinyin-cu-xau", false***REMOVED***,
		***REMOVED***"co", "pinyin", "en-u-co-phonebk-v-xx", "en-u-co-pinyin-v-xx", false***REMOVED***,
		***REMOVED***"co", "pinyin", "en-u-co-phonebk-x-x", "en-u-co-pinyin-x-x", false***REMOVED***,
		***REMOVED***"nu", "arabic", "en-u-co-phonebk-nu-vaai", "en-u-co-phonebk-nu-arabic", false***REMOVED***,
		// add to existing -u extension
		***REMOVED***"co", "pinyin", "en-u-ca-gregory", "en-u-ca-gregory-co-pinyin", false***REMOVED***,
		***REMOVED***"co", "pinyin", "en-u-ca-gregory-nu-vaai", "en-u-ca-gregory-co-pinyin-nu-vaai", false***REMOVED***,
		***REMOVED***"co", "pinyin", "en-u-ca-gregory-v-va", "en-u-ca-gregory-co-pinyin-v-va", false***REMOVED***,
		***REMOVED***"co", "pinyin", "en-u-ca-gregory-x-a", "en-u-ca-gregory-co-pinyin-x-a", false***REMOVED***,
		***REMOVED***"ca", "gregory", "en-u-co-pinyin", "en-u-ca-gregory-co-pinyin", false***REMOVED***,
		// remove pair
		***REMOVED***"co", "", "en-u-co-phonebk", "en", false***REMOVED***,
		***REMOVED***"co", "", "en-u-ca-gregory-co-phonebk", "en-u-ca-gregory", false***REMOVED***,
		***REMOVED***"co", "", "en-u-co-phonebk-nu-arabic", "en-u-nu-arabic", false***REMOVED***,
		***REMOVED***"co", "", "en", "en", false***REMOVED***,
		// add -u extension
		***REMOVED***"co", "pinyin", "en", "en-u-co-pinyin", false***REMOVED***,
		***REMOVED***"co", "pinyin", "und", "und-u-co-pinyin", false***REMOVED***,
		***REMOVED***"co", "pinyin", "en-a-aaa", "en-a-aaa-u-co-pinyin", false***REMOVED***,
		***REMOVED***"co", "pinyin", "en-x-aaa", "en-u-co-pinyin-x-aaa", false***REMOVED***,
		***REMOVED***"co", "pinyin", "en-v-aa", "en-u-co-pinyin-v-aa", false***REMOVED***,
		***REMOVED***"co", "pinyin", "en-a-aaa-x-x", "en-a-aaa-u-co-pinyin-x-x", false***REMOVED***,
		***REMOVED***"co", "pinyin", "en-a-aaa-v-va", "en-a-aaa-u-co-pinyin-v-va", false***REMOVED***,
		// error on invalid values
		***REMOVED***"co", "pinyinxxx", "en", "en", true***REMOVED***,
		***REMOVED***"co", "piny.n", "en", "en", true***REMOVED***,
		***REMOVED***"co", "pinyinxxx", "en-a-aaa", "en-a-aaa", true***REMOVED***,
		***REMOVED***"co", "pinyinxxx", "en-u-aaa", "en-u-aaa", true***REMOVED***,
		***REMOVED***"co", "pinyinxxx", "en-u-aaa-co-pinyin", "en-u-aaa-co-pinyin", true***REMOVED***,
		***REMOVED***"co", "pinyi.", "en-u-aaa-co-pinyin", "en-u-aaa-co-pinyin", true***REMOVED***,
		***REMOVED***"col", "pinyin", "en", "en", true***REMOVED***,
		***REMOVED***"co", "cu", "en", "en", true***REMOVED***,
		// error when setting on a private use tag
		***REMOVED***"co", "phonebook", "x-foo", "x-foo", true***REMOVED***,
	***REMOVED***
	for i, tt := range tests ***REMOVED***
		tag := Make(tt.in)
		if v, err := tag.SetTypeForKey(tt.key, tt.value); v.String() != tt.out ***REMOVED***
			t.Errorf("%d:%q[%q]=%q: was %q; want %q", i, tt.in, tt.key, tt.value, v, tt.out)
		***REMOVED*** else if (err != nil) != tt.err ***REMOVED***
			t.Errorf("%d:%q[%q]=%q: error was %v; want %v", i, tt.in, tt.key, tt.value, err != nil, tt.err)
		***REMOVED*** else if val := v.TypeForKey(tt.key); err == nil && val != tt.value ***REMOVED***
			t.Errorf("%d:%q[%q]==%q: was %v; want %v", i, tt.out, tt.key, tt.value, val, tt.value)
		***REMOVED***
		if len(tag.String()) <= 3 ***REMOVED***
			// Simulate a tag for which the string has not been set.
			tag.str, tag.pExt, tag.pVariant = "", 0, 0
			if tag, err := tag.SetTypeForKey(tt.key, tt.value); err == nil ***REMOVED***
				if val := tag.TypeForKey(tt.key); err == nil && val != tt.value ***REMOVED***
					t.Errorf("%d:%q[%q]==%q: was %v; want %v", i, tt.out, tt.key, tt.value, val, tt.value)
				***REMOVED***
			***REMOVED***
		***REMOVED***
	***REMOVED***
***REMOVED***

func TestFindKeyAndType(t *testing.T) ***REMOVED***
	// out is either the matched type in case of a match or the original
	// string up till the insertion point.
	tests := []struct ***REMOVED***
		key     string
		hasExt  bool
		in, out string
	***REMOVED******REMOVED***
		// Don't search past a private use extension.
		***REMOVED***"co", false, "en-x-foo-u-co-pinyin", "en"***REMOVED***,
		***REMOVED***"co", false, "x-foo-u-co-pinyin", ""***REMOVED***,
		***REMOVED***"co", false, "en-s-fff-x-foo", "en-s-fff"***REMOVED***,
		// Insertion points in absence of -u extension.
		***REMOVED***"cu", false, "en", ""***REMOVED***, // t.str is ""
		***REMOVED***"cu", false, "en-v-va", "en"***REMOVED***,
		***REMOVED***"cu", false, "en-a-va", "en-a-va"***REMOVED***,
		***REMOVED***"cu", false, "en-a-va-v-va", "en-a-va"***REMOVED***,
		***REMOVED***"cu", false, "en-x-a", "en"***REMOVED***,
		// Tags with the -u extension.
		***REMOVED***"co", true, "en-u-co-standard", "standard"***REMOVED***,
		***REMOVED***"co", true, "yue-u-co-pinyin", "pinyin"***REMOVED***,
		***REMOVED***"co", true, "en-u-co-abc", "abc"***REMOVED***,
		***REMOVED***"co", true, "en-u-co-abc-def", "abc-def"***REMOVED***,
		***REMOVED***"co", true, "en-u-co-abc-def-x-foo", "abc-def"***REMOVED***,
		***REMOVED***"co", true, "en-u-co-standard-nu-arab", "standard"***REMOVED***,
		***REMOVED***"co", true, "yue-u-co-pinyin-nu-arab", "pinyin"***REMOVED***,
		// Insertion points.
		***REMOVED***"cu", true, "en-u-co-standard", "en-u-co-standard"***REMOVED***,
		***REMOVED***"cu", true, "yue-u-co-pinyin-x-foo", "yue-u-co-pinyin"***REMOVED***,
		***REMOVED***"cu", true, "en-u-co-abc", "en-u-co-abc"***REMOVED***,
		***REMOVED***"cu", true, "en-u-nu-arabic", "en-u"***REMOVED***,
		***REMOVED***"cu", true, "en-u-co-abc-def-nu-arabic", "en-u-co-abc-def"***REMOVED***,
	***REMOVED***
	for i, tt := range tests ***REMOVED***
		start, end, hasExt := Make(tt.in).findTypeForKey(tt.key)
		if start != end ***REMOVED***
			res := tt.in[start:end]
			if res != tt.out ***REMOVED***
				t.Errorf("%d:%s: was %q; want %q", i, tt.in, res, tt.out)
			***REMOVED***
		***REMOVED*** else ***REMOVED***
			if hasExt != tt.hasExt ***REMOVED***
				t.Errorf("%d:%s: hasExt was %v; want %v", i, tt.in, hasExt, tt.hasExt)
				continue
			***REMOVED***
			if tt.in[:start] != tt.out ***REMOVED***
				t.Errorf("%d:%s: insertion point was %q; want %q", i, tt.in, tt.in[:start], tt.out)
			***REMOVED***
		***REMOVED***
	***REMOVED***
***REMOVED***

func TestParent(t *testing.T) ***REMOVED***
	tests := []struct***REMOVED*** in, out string ***REMOVED******REMOVED***
		// Strip variants and extensions first
		***REMOVED***"de-u-co-phonebk", "de"***REMOVED***,
		***REMOVED***"de-1994", "de"***REMOVED***,
		***REMOVED***"de-Latn-1994", "de"***REMOVED***, // remove superfluous script.

		// Ensure the canonical Tag for an entry is in the chain for base-script
		// pairs.
		***REMOVED***"zh-Hans", "zh"***REMOVED***,

		// Skip the script if it is the maximized version. CLDR files for the
		// skipped tag are always empty.
		***REMOVED***"zh-Hans-TW", "zh"***REMOVED***,
		***REMOVED***"zh-Hans-CN", "zh"***REMOVED***,

		// Insert the script if the maximized script is not the same as the
		// maximized script of the base language.
		***REMOVED***"zh-TW", "zh-Hant"***REMOVED***,
		***REMOVED***"zh-HK", "zh-Hant"***REMOVED***,
		***REMOVED***"zh-Hant-TW", "zh-Hant"***REMOVED***,
		***REMOVED***"zh-Hant-HK", "zh-Hant"***REMOVED***,

		// Non-default script skips to und.
		// CLDR
		***REMOVED***"az-Cyrl", "und"***REMOVED***,
		***REMOVED***"bs-Cyrl", "und"***REMOVED***,
		***REMOVED***"en-Dsrt", "und"***REMOVED***,
		***REMOVED***"ha-Arab", "und"***REMOVED***,
		***REMOVED***"mn-Mong", "und"***REMOVED***,
		***REMOVED***"pa-Arab", "und"***REMOVED***,
		***REMOVED***"shi-Latn", "und"***REMOVED***,
		***REMOVED***"sr-Latn", "und"***REMOVED***,
		***REMOVED***"uz-Arab", "und"***REMOVED***,
		***REMOVED***"uz-Cyrl", "und"***REMOVED***,
		***REMOVED***"vai-Latn", "und"***REMOVED***,
		***REMOVED***"zh-Hant", "und"***REMOVED***,
		// extra
		***REMOVED***"nl-Cyrl", "und"***REMOVED***,

		// World english inherits from en-001.
		***REMOVED***"en-150", "en-001"***REMOVED***,
		***REMOVED***"en-AU", "en-001"***REMOVED***,
		***REMOVED***"en-BE", "en-001"***REMOVED***,
		***REMOVED***"en-GG", "en-001"***REMOVED***,
		***REMOVED***"en-GI", "en-001"***REMOVED***,
		***REMOVED***"en-HK", "en-001"***REMOVED***,
		***REMOVED***"en-IE", "en-001"***REMOVED***,
		***REMOVED***"en-IM", "en-001"***REMOVED***,
		***REMOVED***"en-IN", "en-001"***REMOVED***,
		***REMOVED***"en-JE", "en-001"***REMOVED***,
		***REMOVED***"en-MT", "en-001"***REMOVED***,
		***REMOVED***"en-NZ", "en-001"***REMOVED***,
		***REMOVED***"en-PK", "en-001"***REMOVED***,
		***REMOVED***"en-SG", "en-001"***REMOVED***,

		// Spanish in Latin-American countries have es-419 as parent.
		***REMOVED***"es-AR", "es-419"***REMOVED***,
		***REMOVED***"es-BO", "es-419"***REMOVED***,
		***REMOVED***"es-CL", "es-419"***REMOVED***,
		***REMOVED***"es-CO", "es-419"***REMOVED***,
		***REMOVED***"es-CR", "es-419"***REMOVED***,
		***REMOVED***"es-CU", "es-419"***REMOVED***,
		***REMOVED***"es-DO", "es-419"***REMOVED***,
		***REMOVED***"es-EC", "es-419"***REMOVED***,
		***REMOVED***"es-GT", "es-419"***REMOVED***,
		***REMOVED***"es-HN", "es-419"***REMOVED***,
		***REMOVED***"es-MX", "es-419"***REMOVED***,
		***REMOVED***"es-NI", "es-419"***REMOVED***,
		***REMOVED***"es-PA", "es-419"***REMOVED***,
		***REMOVED***"es-PE", "es-419"***REMOVED***,
		***REMOVED***"es-PR", "es-419"***REMOVED***,
		***REMOVED***"es-PY", "es-419"***REMOVED***,
		***REMOVED***"es-SV", "es-419"***REMOVED***,
		***REMOVED***"es-US", "es-419"***REMOVED***,
		***REMOVED***"es-UY", "es-419"***REMOVED***,
		***REMOVED***"es-VE", "es-419"***REMOVED***,
		// exceptions (according to CLDR)
		***REMOVED***"es-CW", "es"***REMOVED***,

		// Inherit from pt-PT, instead of pt for these countries.
		***REMOVED***"pt-AO", "pt-PT"***REMOVED***,
		***REMOVED***"pt-CV", "pt-PT"***REMOVED***,
		***REMOVED***"pt-GW", "pt-PT"***REMOVED***,
		***REMOVED***"pt-MO", "pt-PT"***REMOVED***,
		***REMOVED***"pt-MZ", "pt-PT"***REMOVED***,
		***REMOVED***"pt-ST", "pt-PT"***REMOVED***,
		***REMOVED***"pt-TL", "pt-PT"***REMOVED***,
	***REMOVED***
	for _, tt := range tests ***REMOVED***
		tag := Raw.MustParse(tt.in)
		if p := Raw.MustParse(tt.out); p != tag.Parent() ***REMOVED***
			t.Errorf("%s: was %v; want %v", tt.in, tag.Parent(), p)
		***REMOVED***
	***REMOVED***
***REMOVED***

var (
	// Tags without error that don't need to be changed.
	benchBasic = []string***REMOVED***
		"en",
		"en-Latn",
		"en-GB",
		"za",
		"zh-Hant",
		"zh",
		"zh-HK",
		"ar-MK",
		"en-CA",
		"fr-CA",
		"fr-CH",
		"fr",
		"lv",
		"he-IT",
		"tlh",
		"ja",
		"ja-Jpan",
		"ja-Jpan-JP",
		"de-1996",
		"de-CH",
		"sr",
		"sr-Latn",
	***REMOVED***
	// Tags with extensions, not changes required.
	benchExt = []string***REMOVED***
		"x-a-b-c-d",
		"x-aa-bbbb-cccccccc-d",
		"en-x_cc-b-bbb-a-aaa",
		"en-c_cc-b-bbb-a-aaa-x-x",
		"en-u-co-phonebk",
		"en-Cyrl-u-co-phonebk",
		"en-US-u-co-phonebk-cu-xau",
		"en-nedix-u-co-phonebk",
		"en-t-t0-abcd",
		"en-t-nl-latn",
		"en-t-t0-abcd-x-a",
	***REMOVED***
	// Change, but not memory allocation required.
	benchSimpleChange = []string***REMOVED***
		"EN",
		"i-klingon",
		"en-latn",
		"zh-cmn-Hans-CN",
		"iw-NL",
	***REMOVED***
	// Change and memory allocation required.
	benchChangeAlloc = []string***REMOVED***
		"en-c_cc-b-bbb-a-aaa",
		"en-u-cu-xua-co-phonebk",
		"en-u-cu-xua-co-phonebk-a-cd",
		"en-u-def-abc-cu-xua-co-phonebk",
		"en-t-en-Cyrl-NL-1994",
		"en-t-en-Cyrl-NL-1994-t0-abc-def",
	***REMOVED***
	// Tags that result in errors.
	benchErr = []string***REMOVED***
		// IllFormed
		"x_A.-B-C_D",
		"en-u-cu-co-phonebk",
		"en-u-cu-xau-co",
		"en-t-nl-abcd",
		// Invalid
		"xx",
		"nl-Uuuu",
		"nl-QB",
	***REMOVED***
	benchChange = append(benchSimpleChange, benchChangeAlloc...)
	benchAll    = append(append(append(benchBasic, benchExt...), benchChange...), benchErr...)
)

func doParse(b *testing.B, tag []string) ***REMOVED***
	for i := 0; i < b.N; i++ ***REMOVED***
		// Use the modulo instead of looping over all tags so that we get a somewhat
		// meaningful ns/op.
		Parse(tag[i%len(tag)])
	***REMOVED***
***REMOVED***

func BenchmarkParse(b *testing.B) ***REMOVED***
	doParse(b, benchAll)
***REMOVED***

func BenchmarkParseBasic(b *testing.B) ***REMOVED***
	doParse(b, benchBasic)
***REMOVED***

func BenchmarkParseError(b *testing.B) ***REMOVED***
	doParse(b, benchErr)
***REMOVED***

func BenchmarkParseSimpleChange(b *testing.B) ***REMOVED***
	doParse(b, benchSimpleChange)
***REMOVED***

func BenchmarkParseChangeAlloc(b *testing.B) ***REMOVED***
	doParse(b, benchChangeAlloc)
***REMOVED***
