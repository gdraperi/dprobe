// Copyright 2013 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package language

import (
	"bytes"
	"strings"
	"testing"

	"golang.org/x/text/internal/tag"
)

type scanTest struct ***REMOVED***
	ok  bool // true if scanning does not result in an error
	in  string
	tok []string // the expected tokens
***REMOVED***

var tests = []scanTest***REMOVED***
	***REMOVED***true, "", []string***REMOVED******REMOVED******REMOVED***,
	***REMOVED***true, "1", []string***REMOVED***"1"***REMOVED******REMOVED***,
	***REMOVED***true, "en", []string***REMOVED***"en"***REMOVED******REMOVED***,
	***REMOVED***true, "root", []string***REMOVED***"root"***REMOVED******REMOVED***,
	***REMOVED***true, "maxchars", []string***REMOVED***"maxchars"***REMOVED******REMOVED***,
	***REMOVED***false, "bad/", []string***REMOVED******REMOVED******REMOVED***,
	***REMOVED***false, "morethan8", []string***REMOVED******REMOVED******REMOVED***,
	***REMOVED***false, "-", []string***REMOVED******REMOVED******REMOVED***,
	***REMOVED***false, "----", []string***REMOVED******REMOVED******REMOVED***,
	***REMOVED***false, "_", []string***REMOVED******REMOVED******REMOVED***,
	***REMOVED***true, "en-US", []string***REMOVED***"en", "US"***REMOVED******REMOVED***,
	***REMOVED***true, "en_US", []string***REMOVED***"en", "US"***REMOVED******REMOVED***,
	***REMOVED***false, "en-US-", []string***REMOVED***"en", "US"***REMOVED******REMOVED***,
	***REMOVED***false, "en-US--", []string***REMOVED***"en", "US"***REMOVED******REMOVED***,
	***REMOVED***false, "en-US---", []string***REMOVED***"en", "US"***REMOVED******REMOVED***,
	***REMOVED***false, "en--US", []string***REMOVED***"en", "US"***REMOVED******REMOVED***,
	***REMOVED***false, "-en-US", []string***REMOVED***"en", "US"***REMOVED******REMOVED***,
	***REMOVED***false, "-en--US-", []string***REMOVED***"en", "US"***REMOVED******REMOVED***,
	***REMOVED***false, "-en--US-", []string***REMOVED***"en", "US"***REMOVED******REMOVED***,
	***REMOVED***false, "en-.-US", []string***REMOVED***"en", "US"***REMOVED******REMOVED***,
	***REMOVED***false, ".-en--US-.", []string***REMOVED***"en", "US"***REMOVED******REMOVED***,
	***REMOVED***false, "en-u.-US", []string***REMOVED***"en", "US"***REMOVED******REMOVED***,
	***REMOVED***true, "en-u1-US", []string***REMOVED***"en", "u1", "US"***REMOVED******REMOVED***,
	***REMOVED***true, "maxchar1_maxchar2-maxchar3", []string***REMOVED***"maxchar1", "maxchar2", "maxchar3"***REMOVED******REMOVED***,
	***REMOVED***false, "moreThan8-moreThan8-e", []string***REMOVED***"e"***REMOVED******REMOVED***,
***REMOVED***

func TestScan(t *testing.T) ***REMOVED***
	for i, tt := range tests ***REMOVED***
		scan := makeScannerString(tt.in)
		for j := 0; !scan.done; j++ ***REMOVED***
			if j >= len(tt.tok) ***REMOVED***
				t.Errorf("%d: extra token %q", i, scan.token)
			***REMOVED*** else if tag.Compare(tt.tok[j], scan.token) != 0 ***REMOVED***
				t.Errorf("%d: token %d: found %q; want %q", i, j, scan.token, tt.tok[j])
				break
			***REMOVED***
			scan.scan()
		***REMOVED***
		if s := strings.Join(tt.tok, "-"); tag.Compare(s, bytes.Replace(scan.b, b("_"), b("-"), -1)) != 0 ***REMOVED***
			t.Errorf("%d: input: found %q; want %q", i, scan.b, s)
		***REMOVED***
		if (scan.err == nil) != tt.ok ***REMOVED***
			t.Errorf("%d: ok: found %v; want %v", i, scan.err == nil, tt.ok)
		***REMOVED***
	***REMOVED***
***REMOVED***

func TestAcceptMinSize(t *testing.T) ***REMOVED***
	for i, tt := range tests ***REMOVED***
		// count number of successive tokens with a minimum size.
		for sz := 1; sz <= 8; sz++ ***REMOVED***
			scan := makeScannerString(tt.in)
			scan.end, scan.next = 0, 0
			end := scan.acceptMinSize(sz)
			n := 0
			for i := 0; i < len(tt.tok) && len(tt.tok[i]) >= sz; i++ ***REMOVED***
				n += len(tt.tok[i])
				if i > 0 ***REMOVED***
					n++
				***REMOVED***
			***REMOVED***
			if end != n ***REMOVED***
				t.Errorf("%d:%d: found len %d; want %d", i, sz, end, n)
			***REMOVED***
		***REMOVED***
	***REMOVED***
***REMOVED***

type parseTest struct ***REMOVED***
	i                    int // the index of this test
	in                   string
	lang, script, region string
	variants, ext        string
	extList              []string // only used when more than one extension is present
	invalid              bool
	rewrite              bool // special rewrite not handled by parseTag
	changed              bool // string needed to be reformatted
***REMOVED***

func parseTests() []parseTest ***REMOVED***
	tests := []parseTest***REMOVED***
		***REMOVED***in: "root", lang: "und"***REMOVED***,
		***REMOVED***in: "und", lang: "und"***REMOVED***,
		***REMOVED***in: "en", lang: "en"***REMOVED***,
		***REMOVED***in: "xy", lang: "und", invalid: true***REMOVED***,
		***REMOVED***in: "en-ZY", lang: "en", invalid: true***REMOVED***,
		***REMOVED***in: "gsw", lang: "gsw"***REMOVED***,
		***REMOVED***in: "sr_Latn", lang: "sr", script: "Latn"***REMOVED***,
		***REMOVED***in: "af-Arab", lang: "af", script: "Arab"***REMOVED***,
		***REMOVED***in: "nl-BE", lang: "nl", region: "BE"***REMOVED***,
		***REMOVED***in: "es-419", lang: "es", region: "419"***REMOVED***,
		***REMOVED***in: "und-001", lang: "und", region: "001"***REMOVED***,
		***REMOVED***in: "de-latn-be", lang: "de", script: "Latn", region: "BE"***REMOVED***,
		// Variants
		***REMOVED***in: "de-1901", lang: "de", variants: "1901"***REMOVED***,
		// Accept with unsuppressed script.
		***REMOVED***in: "de-Latn-1901", lang: "de", script: "Latn", variants: "1901"***REMOVED***,
		// Specialized.
		***REMOVED***in: "sl-rozaj", lang: "sl", variants: "rozaj"***REMOVED***,
		***REMOVED***in: "sl-rozaj-lipaw", lang: "sl", variants: "rozaj-lipaw"***REMOVED***,
		***REMOVED***in: "sl-rozaj-biske", lang: "sl", variants: "rozaj-biske"***REMOVED***,
		***REMOVED***in: "sl-rozaj-biske-1994", lang: "sl", variants: "rozaj-biske-1994"***REMOVED***,
		***REMOVED***in: "sl-rozaj-1994", lang: "sl", variants: "rozaj-1994"***REMOVED***,
		// Maximum number of variants while adhering to prefix rules.
		***REMOVED***in: "sl-rozaj-biske-1994-alalc97-fonipa-fonupa-fonxsamp", lang: "sl", variants: "rozaj-biske-1994-alalc97-fonipa-fonupa-fonxsamp"***REMOVED***,

		// Sorting.
		***REMOVED***in: "sl-1994-biske-rozaj", lang: "sl", variants: "rozaj-biske-1994", changed: true***REMOVED***,
		***REMOVED***in: "sl-rozaj-biske-1994-alalc97-fonupa-fonipa-fonxsamp", lang: "sl", variants: "rozaj-biske-1994-alalc97-fonipa-fonupa-fonxsamp", changed: true***REMOVED***,
		***REMOVED***in: "nl-fonxsamp-alalc97-fonipa-fonupa", lang: "nl", variants: "alalc97-fonipa-fonupa-fonxsamp", changed: true***REMOVED***,

		// Duplicates variants are removed, but not an error.
		***REMOVED***in: "nl-fonupa-fonupa", lang: "nl", variants: "fonupa"***REMOVED***,

		// Variants that do not have correct prefixes. We still accept these.
		***REMOVED***in: "de-Cyrl-1901", lang: "de", script: "Cyrl", variants: "1901"***REMOVED***,
		***REMOVED***in: "sl-rozaj-lipaw-1994", lang: "sl", variants: "rozaj-lipaw-1994"***REMOVED***,
		***REMOVED***in: "sl-1994-biske-rozaj-1994-biske-rozaj", lang: "sl", variants: "rozaj-biske-1994", changed: true***REMOVED***,
		***REMOVED***in: "de-Cyrl-1901", lang: "de", script: "Cyrl", variants: "1901"***REMOVED***,

		// Invalid variant.
		***REMOVED***in: "de-1902", lang: "de", variants: "", invalid: true***REMOVED***,

		***REMOVED***in: "EN_CYRL", lang: "en", script: "Cyrl"***REMOVED***,
		// private use and extensions
		***REMOVED***in: "x-a-b-c-d", ext: "x-a-b-c-d"***REMOVED***,
		***REMOVED***in: "x_A.-B-C_D", ext: "x-b-c-d", invalid: true, changed: true***REMOVED***,
		***REMOVED***in: "x-aa-bbbb-cccccccc-d", ext: "x-aa-bbbb-cccccccc-d"***REMOVED***,
		***REMOVED***in: "en-c_cc-b-bbb-a-aaa", lang: "en", changed: true, extList: []string***REMOVED***"a-aaa", "b-bbb", "c-cc"***REMOVED******REMOVED***,
		***REMOVED***in: "en-x_cc-b-bbb-a-aaa", lang: "en", ext: "x-cc-b-bbb-a-aaa", changed: true***REMOVED***,
		***REMOVED***in: "en-c_cc-b-bbb-a-aaa-x-x", lang: "en", changed: true, extList: []string***REMOVED***"a-aaa", "b-bbb", "c-cc", "x-x"***REMOVED******REMOVED***,
		***REMOVED***in: "en-v-c", lang: "en", ext: "", invalid: true***REMOVED***,
		***REMOVED***in: "en-v-abcdefghi", lang: "en", ext: "", invalid: true***REMOVED***,
		***REMOVED***in: "en-v-abc-x", lang: "en", ext: "v-abc", invalid: true***REMOVED***,
		***REMOVED***in: "en-v-abc-x-", lang: "en", ext: "v-abc", invalid: true***REMOVED***,
		***REMOVED***in: "en-v-abc-w-x-xx", lang: "en", extList: []string***REMOVED***"v-abc", "x-xx"***REMOVED***, invalid: true, changed: true***REMOVED***,
		***REMOVED***in: "en-v-abc-w-y-yx", lang: "en", extList: []string***REMOVED***"v-abc", "y-yx"***REMOVED***, invalid: true, changed: true***REMOVED***,
		***REMOVED***in: "en-v-c-abc", lang: "en", ext: "c-abc", invalid: true, changed: true***REMOVED***,
		***REMOVED***in: "en-v-w-abc", lang: "en", ext: "w-abc", invalid: true, changed: true***REMOVED***,
		***REMOVED***in: "en-v-x-abc", lang: "en", ext: "x-abc", invalid: true, changed: true***REMOVED***,
		***REMOVED***in: "en-v-x-a", lang: "en", ext: "x-a", invalid: true, changed: true***REMOVED***,
		***REMOVED***in: "en-9-aa-0-aa-z-bb-x-a", lang: "en", extList: []string***REMOVED***"0-aa", "9-aa", "z-bb", "x-a"***REMOVED***, changed: true***REMOVED***,
		***REMOVED***in: "en-u-c", lang: "en", ext: "", invalid: true***REMOVED***,
		***REMOVED***in: "en-u-co-phonebk", lang: "en", ext: "u-co-phonebk"***REMOVED***,
		***REMOVED***in: "en-u-co-phonebk-ca", lang: "en", ext: "u-co-phonebk", invalid: true***REMOVED***,
		***REMOVED***in: "en-u-nu-arabic-co-phonebk-ca", lang: "en", ext: "u-co-phonebk-nu-arabic", invalid: true, changed: true***REMOVED***,
		***REMOVED***in: "en-u-nu-arabic-co-phonebk-ca-x", lang: "en", ext: "u-co-phonebk-nu-arabic", invalid: true, changed: true***REMOVED***,
		***REMOVED***in: "en-u-nu-arabic-co-phonebk-ca-s", lang: "en", ext: "u-co-phonebk-nu-arabic", invalid: true, changed: true***REMOVED***,
		***REMOVED***in: "en-u-nu-arabic-co-phonebk-ca-a12345678", lang: "en", ext: "u-co-phonebk-nu-arabic", invalid: true, changed: true***REMOVED***,
		***REMOVED***in: "en-u-co-phonebook", lang: "en", ext: "", invalid: true***REMOVED***,
		***REMOVED***in: "en-u-co-phonebook-cu-xau", lang: "en", ext: "u-cu-xau", invalid: true, changed: true***REMOVED***,
		***REMOVED***in: "en-Cyrl-u-co-phonebk", lang: "en", script: "Cyrl", ext: "u-co-phonebk"***REMOVED***,
		***REMOVED***in: "en-US-u-co-phonebk", lang: "en", region: "US", ext: "u-co-phonebk"***REMOVED***,
		***REMOVED***in: "en-US-u-co-phonebk-cu-xau", lang: "en", region: "US", ext: "u-co-phonebk-cu-xau"***REMOVED***,
		***REMOVED***in: "en-scotland-u-co-phonebk", lang: "en", variants: "scotland", ext: "u-co-phonebk"***REMOVED***,
		***REMOVED***in: "en-u-cu-xua-co-phonebk", lang: "en", ext: "u-co-phonebk-cu-xua", changed: true***REMOVED***,
		***REMOVED***in: "en-u-def-abc-cu-xua-co-phonebk", lang: "en", ext: "u-abc-def-co-phonebk-cu-xua", changed: true***REMOVED***,
		***REMOVED***in: "en-u-def-abc", lang: "en", ext: "u-abc-def", changed: true***REMOVED***,
		***REMOVED***in: "en-u-cu-xua-co-phonebk-a-cd", lang: "en", extList: []string***REMOVED***"a-cd", "u-co-phonebk-cu-xua"***REMOVED***, changed: true***REMOVED***,
		// Invalid "u" extension. Drop invalid parts.
		***REMOVED***in: "en-u-cu-co-phonebk", lang: "en", extList: []string***REMOVED***"u-co-phonebk"***REMOVED***, invalid: true, changed: true***REMOVED***,
		***REMOVED***in: "en-u-cu-xau-co", lang: "en", extList: []string***REMOVED***"u-cu-xau"***REMOVED***, invalid: true***REMOVED***,
		// We allow duplicate keys as the LDML spec does not explicitly prohibit it.
		// TODO: Consider eliminating duplicates and returning an error.
		***REMOVED***in: "en-u-cu-xau-co-phonebk-cu-xau", lang: "en", ext: "u-co-phonebk-cu-xau-cu-xau", changed: true***REMOVED***,
		***REMOVED***in: "en-t-en-Cyrl-NL-fonipa", lang: "en", ext: "t-en-cyrl-nl-fonipa", changed: true***REMOVED***,
		***REMOVED***in: "en-t-en-Cyrl-NL-fonipa-t0-abc-def", lang: "en", ext: "t-en-cyrl-nl-fonipa-t0-abc-def", changed: true***REMOVED***,
		***REMOVED***in: "en-t-t0-abcd", lang: "en", ext: "t-t0-abcd"***REMOVED***,
		// Not necessary to have changed here.
		***REMOVED***in: "en-t-nl-abcd", lang: "en", ext: "t-nl", invalid: true***REMOVED***,
		***REMOVED***in: "en-t-nl-latn", lang: "en", ext: "t-nl-latn"***REMOVED***,
		***REMOVED***in: "en-t-t0-abcd-x-a", lang: "en", extList: []string***REMOVED***"t-t0-abcd", "x-a"***REMOVED******REMOVED***,
		// invalid
		***REMOVED***in: "", lang: "und", invalid: true***REMOVED***,
		***REMOVED***in: "-", lang: "und", invalid: true***REMOVED***,
		***REMOVED***in: "x", lang: "und", invalid: true***REMOVED***,
		***REMOVED***in: "x-", lang: "und", invalid: true***REMOVED***,
		***REMOVED***in: "x--", lang: "und", invalid: true***REMOVED***,
		***REMOVED***in: "a-a-b-c-d", lang: "und", invalid: true***REMOVED***,
		***REMOVED***in: "en-", lang: "en", invalid: true***REMOVED***,
		***REMOVED***in: "enne-", lang: "und", invalid: true***REMOVED***,
		***REMOVED***in: "en.", lang: "und", invalid: true***REMOVED***,
		***REMOVED***in: "en.-latn", lang: "und", invalid: true***REMOVED***,
		***REMOVED***in: "en.-en", lang: "en", invalid: true***REMOVED***,
		***REMOVED***in: "x-a-tooManyChars-c-d", ext: "x-a-c-d", invalid: true, changed: true***REMOVED***,
		***REMOVED***in: "a-tooManyChars-c-d", lang: "und", invalid: true***REMOVED***,
		// TODO: check key-value validity
		// ***REMOVED*** in: "en-u-cu-xd", lang: "en", ext: "u-cu-xd", invalid: true ***REMOVED***,
		***REMOVED***in: "en-t-abcd", lang: "en", invalid: true***REMOVED***,
		***REMOVED***in: "en-Latn-US-en", lang: "en", script: "Latn", region: "US", invalid: true***REMOVED***,
		// rewrites (more tests in TestGrandfathered)
		***REMOVED***in: "zh-min-nan", lang: "nan"***REMOVED***,
		***REMOVED***in: "zh-yue", lang: "yue"***REMOVED***,
		***REMOVED***in: "zh-xiang", lang: "hsn", rewrite: true***REMOVED***,
		***REMOVED***in: "zh-guoyu", lang: "cmn", rewrite: true***REMOVED***,
		***REMOVED***in: "iw", lang: "iw"***REMOVED***,
		***REMOVED***in: "sgn-BE-FR", lang: "sfb", rewrite: true***REMOVED***,
		***REMOVED***in: "i-klingon", lang: "tlh", rewrite: true***REMOVED***,
	***REMOVED***
	for i, tt := range tests ***REMOVED***
		tests[i].i = i
		if tt.extList != nil ***REMOVED***
			tests[i].ext = strings.Join(tt.extList, "-")
		***REMOVED***
		if tt.ext != "" && tt.extList == nil ***REMOVED***
			tests[i].extList = []string***REMOVED***tt.ext***REMOVED***
		***REMOVED***
	***REMOVED***
	return tests
***REMOVED***

func TestParseExtensions(t *testing.T) ***REMOVED***
	for i, tt := range parseTests() ***REMOVED***
		if tt.ext == "" || tt.rewrite ***REMOVED***
			continue
		***REMOVED***
		scan := makeScannerString(tt.in)
		if len(scan.b) > 1 && scan.b[1] != '-' ***REMOVED***
			scan.end = nextExtension(string(scan.b), 0)
			scan.next = scan.end + 1
			scan.scan()
		***REMOVED***
		start := scan.start
		scan.toLower(start, len(scan.b))
		parseExtensions(&scan)
		ext := string(scan.b[start:])
		if ext != tt.ext ***REMOVED***
			t.Errorf("%d(%s): ext was %v; want %v", i, tt.in, ext, tt.ext)
		***REMOVED***
		if changed := !strings.HasPrefix(tt.in[start:], ext); changed != tt.changed ***REMOVED***
			t.Errorf("%d(%s): changed was %v; want %v", i, tt.in, changed, tt.changed)
		***REMOVED***
	***REMOVED***
***REMOVED***

// partChecks runs checks for each part by calling the function returned by f.
func partChecks(t *testing.T, f func(*parseTest) (Tag, bool)) ***REMOVED***
	for i, tt := range parseTests() ***REMOVED***
		tag, skip := f(&tt)
		if skip ***REMOVED***
			continue
		***REMOVED***
		if l, _ := getLangID(b(tt.lang)); l != tag.lang ***REMOVED***
			t.Errorf("%d: lang was %q; want %q", i, tag.lang, l)
		***REMOVED***
		if sc, _ := getScriptID(script, b(tt.script)); sc != tag.script ***REMOVED***
			t.Errorf("%d: script was %q; want %q", i, tag.script, sc)
		***REMOVED***
		if r, _ := getRegionID(b(tt.region)); r != tag.region ***REMOVED***
			t.Errorf("%d: region was %q; want %q", i, tag.region, r)
		***REMOVED***
		if tag.str == "" ***REMOVED***
			continue
		***REMOVED***
		p := int(tag.pVariant)
		if p < int(tag.pExt) ***REMOVED***
			p++
		***REMOVED***
		if s, g := tag.str[p:tag.pExt], tt.variants; s != g ***REMOVED***
			t.Errorf("%d: variants was %q; want %q", i, s, g)
		***REMOVED***
		p = int(tag.pExt)
		if p > 0 && p < len(tag.str) ***REMOVED***
			p++
		***REMOVED***
		if s, g := (tag.str)[p:], tt.ext; s != g ***REMOVED***
			t.Errorf("%d: extensions were %q; want %q", i, s, g)
		***REMOVED***
	***REMOVED***
***REMOVED***

func TestParseTag(t *testing.T) ***REMOVED***
	partChecks(t, func(tt *parseTest) (id Tag, skip bool) ***REMOVED***
		if strings.HasPrefix(tt.in, "x-") || tt.rewrite ***REMOVED***
			return Tag***REMOVED******REMOVED***, true
		***REMOVED***
		scan := makeScannerString(tt.in)
		id, end := parseTag(&scan)
		id.str = string(scan.b[:end])
		tt.ext = ""
		tt.extList = []string***REMOVED******REMOVED***
		return id, false
	***REMOVED***)
***REMOVED***

func TestParse(t *testing.T) ***REMOVED***
	partChecks(t, func(tt *parseTest) (id Tag, skip bool) ***REMOVED***
		id, err := Raw.Parse(tt.in)
		ext := ""
		if id.str != "" ***REMOVED***
			if strings.HasPrefix(id.str, "x-") ***REMOVED***
				ext = id.str
			***REMOVED*** else if int(id.pExt) < len(id.str) && id.pExt > 0 ***REMOVED***
				ext = id.str[id.pExt+1:]
			***REMOVED***
		***REMOVED***
		if tag, _ := Raw.Parse(id.String()); tag.String() != id.String() ***REMOVED***
			t.Errorf("%d:%s: reparse was %q; want %q", tt.i, tt.in, id.String(), tag.String())
		***REMOVED***
		if ext != tt.ext ***REMOVED***
			t.Errorf("%d:%s: ext was %q; want %q", tt.i, tt.in, ext, tt.ext)
		***REMOVED***
		changed := id.str != "" && !strings.HasPrefix(tt.in, id.str)
		if changed != tt.changed ***REMOVED***
			t.Errorf("%d:%s: changed was %v; want %v", tt.i, tt.in, changed, tt.changed)
		***REMOVED***
		if (err != nil) != tt.invalid ***REMOVED***
			t.Errorf("%d:%s: invalid was %v; want %v. Error: %v", tt.i, tt.in, err != nil, tt.invalid, err)
		***REMOVED***
		return id, false
	***REMOVED***)
***REMOVED***

func TestErrors(t *testing.T) ***REMOVED***
	mkInvalid := func(s string) error ***REMOVED***
		return mkErrInvalid([]byte(s))
	***REMOVED***
	tests := []struct ***REMOVED***
		in  string
		out error
	***REMOVED******REMOVED***
		// invalid subtags.
		***REMOVED***"ac", mkInvalid("ac")***REMOVED***,
		***REMOVED***"AC", mkInvalid("ac")***REMOVED***,
		***REMOVED***"aa-Uuuu", mkInvalid("Uuuu")***REMOVED***,
		***REMOVED***"aa-AB", mkInvalid("AB")***REMOVED***,
		// ill-formed wins over invalid.
		***REMOVED***"ac-u", errSyntax***REMOVED***,
		***REMOVED***"ac-u-ca", errSyntax***REMOVED***,
		***REMOVED***"ac-u-ca-co-pinyin", errSyntax***REMOVED***,
		***REMOVED***"noob", errSyntax***REMOVED***,
	***REMOVED***
	for _, tt := range tests ***REMOVED***
		_, err := Parse(tt.in)
		if err != tt.out ***REMOVED***
			t.Errorf("%s: was %q; want %q", tt.in, err, tt.out)
		***REMOVED***
	***REMOVED***
***REMOVED***

func TestCompose1(t *testing.T) ***REMOVED***
	partChecks(t, func(tt *parseTest) (id Tag, skip bool) ***REMOVED***
		l, _ := ParseBase(tt.lang)
		s, _ := ParseScript(tt.script)
		r, _ := ParseRegion(tt.region)
		v := []Variant***REMOVED******REMOVED***
		for _, x := range strings.Split(tt.variants, "-") ***REMOVED***
			p, _ := ParseVariant(x)
			v = append(v, p)
		***REMOVED***
		e := []Extension***REMOVED******REMOVED***
		for _, x := range tt.extList ***REMOVED***
			p, _ := ParseExtension(x)
			e = append(e, p)
		***REMOVED***
		id, _ = Raw.Compose(l, s, r, v, e)
		return id, false
	***REMOVED***)
***REMOVED***

func TestCompose2(t *testing.T) ***REMOVED***
	partChecks(t, func(tt *parseTest) (id Tag, skip bool) ***REMOVED***
		l, _ := ParseBase(tt.lang)
		s, _ := ParseScript(tt.script)
		r, _ := ParseRegion(tt.region)
		p := []interface***REMOVED******REMOVED******REMOVED***l, s, r, s, r, l***REMOVED***
		for _, x := range strings.Split(tt.variants, "-") ***REMOVED***
			v, _ := ParseVariant(x)
			p = append(p, v)
		***REMOVED***
		for _, x := range tt.extList ***REMOVED***
			e, _ := ParseExtension(x)
			p = append(p, e)
		***REMOVED***
		id, _ = Raw.Compose(p...)
		return id, false
	***REMOVED***)
***REMOVED***

func TestCompose3(t *testing.T) ***REMOVED***
	partChecks(t, func(tt *parseTest) (id Tag, skip bool) ***REMOVED***
		id, _ = Raw.Parse(tt.in)
		id, _ = Raw.Compose(id)
		return id, false
	***REMOVED***)
***REMOVED***

func mk(s string) Tag ***REMOVED***
	return Raw.Make(s)
***REMOVED***

func TestParseAcceptLanguage(t *testing.T) ***REMOVED***
	type res struct ***REMOVED***
		t Tag
		q float32
	***REMOVED***
	en := []res***REMOVED******REMOVED***mk("en"), 1.0***REMOVED******REMOVED***
	tests := []struct ***REMOVED***
		out []res
		in  string
		ok  bool
	***REMOVED******REMOVED***
		***REMOVED***en, "en", true***REMOVED***,
		***REMOVED***en, "   en", true***REMOVED***,
		***REMOVED***en, "en   ", true***REMOVED***,
		***REMOVED***en, "  en  ", true***REMOVED***,
		***REMOVED***en, "en,", true***REMOVED***,
		***REMOVED***en, ",en", true***REMOVED***,
		***REMOVED***en, ",,,en,,,", true***REMOVED***,
		***REMOVED***en, ",en;q=1", true***REMOVED***,

		// We allow an empty input, contrary to spec.
		***REMOVED***nil, "", true***REMOVED***,
		***REMOVED***[]res***REMOVED******REMOVED***mk("aa"), 1***REMOVED******REMOVED***, "aa;", true***REMOVED***, // allow unspecified weight

		// errors
		***REMOVED***nil, ";", false***REMOVED***,
		***REMOVED***nil, "$", false***REMOVED***,
		***REMOVED***nil, "e;", false***REMOVED***,
		***REMOVED***nil, "x;", false***REMOVED***,
		***REMOVED***nil, "x", false***REMOVED***,
		***REMOVED***nil, "ac", false***REMOVED***, // non-existing language
		***REMOVED***nil, "aa;q", false***REMOVED***,
		***REMOVED***nil, "aa;q=", false***REMOVED***,
		***REMOVED***nil, "aa;q=.", false***REMOVED***,

		// odd fallbacks
		***REMOVED***
			[]res***REMOVED******REMOVED***mk("en"), 0.1***REMOVED******REMOVED***,
			" english ;q=.1",
			true,
		***REMOVED***,
		***REMOVED***
			[]res***REMOVED******REMOVED***mk("it"), 1.0***REMOVED***, ***REMOVED***mk("de"), 1.0***REMOVED***, ***REMOVED***mk("fr"), 1.0***REMOVED******REMOVED***,
			" italian, deutsch, french",
			true,
		***REMOVED***,

		// lists
		***REMOVED***
			[]res***REMOVED******REMOVED***mk("en"), 0.1***REMOVED******REMOVED***,
			"en;q=.1",
			true,
		***REMOVED***,
		***REMOVED***
			[]res***REMOVED******REMOVED***mk("mul"), 1.0***REMOVED******REMOVED***,
			"*",
			true,
		***REMOVED***,
		***REMOVED***
			[]res***REMOVED******REMOVED***mk("en"), 1.0***REMOVED***, ***REMOVED***mk("de"), 1.0***REMOVED******REMOVED***,
			"en,de",
			true,
		***REMOVED***,
		***REMOVED***
			[]res***REMOVED******REMOVED***mk("en"), 1.0***REMOVED***, ***REMOVED***mk("de"), .5***REMOVED******REMOVED***,
			"en,de;q=0.5",
			true,
		***REMOVED***,
		***REMOVED***
			[]res***REMOVED******REMOVED***mk("de"), 0.8***REMOVED***, ***REMOVED***mk("en"), 0.5***REMOVED******REMOVED***,
			"  en ;   q    =   0.5    ,  , de;q=0.8",
			true,
		***REMOVED***,
		***REMOVED***
			[]res***REMOVED******REMOVED***mk("en"), 1.0***REMOVED***, ***REMOVED***mk("de"), 1.0***REMOVED***, ***REMOVED***mk("fr"), 1.0***REMOVED***, ***REMOVED***mk("tlh"), 1.0***REMOVED******REMOVED***,
			"en,de,fr,i-klingon",
			true,
		***REMOVED***,
		// sorting
		***REMOVED***
			[]res***REMOVED******REMOVED***mk("tlh"), 0.4***REMOVED***, ***REMOVED***mk("de"), 0.2***REMOVED***, ***REMOVED***mk("fr"), 0.2***REMOVED***, ***REMOVED***mk("en"), 0.1***REMOVED******REMOVED***,
			"en;q=0.1,de;q=0.2,fr;q=0.2,i-klingon;q=0.4",
			true,
		***REMOVED***,
		// dropping
		***REMOVED***
			[]res***REMOVED******REMOVED***mk("fr"), 0.2***REMOVED***, ***REMOVED***mk("en"), 0.1***REMOVED******REMOVED***,
			"en;q=0.1,de;q=0,fr;q=0.2,i-klingon;q=0.0",
			true,
		***REMOVED***,
	***REMOVED***
	for i, tt := range tests ***REMOVED***
		tags, qs, e := ParseAcceptLanguage(tt.in)
		if e == nil != tt.ok ***REMOVED***
			t.Errorf("%d:%s:err: was %v; want %v", i, tt.in, e == nil, tt.ok)
		***REMOVED***
		for j, tag := range tags ***REMOVED***
			if out := tt.out[j]; !tag.equalTags(out.t) || qs[j] != out.q ***REMOVED***
				t.Errorf("%d:%s: was %s, %1f; want %s, %1f", i, tt.in, tag, qs[j], out.t, out.q)
				break
			***REMOVED***
		***REMOVED***
	***REMOVED***
***REMOVED***
