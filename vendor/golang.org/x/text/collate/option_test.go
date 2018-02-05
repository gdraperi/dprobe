// Copyright 2014 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
package collate

import (
	"reflect"
	"strings"
	"testing"

	"golang.org/x/text/internal/colltab"
	"golang.org/x/text/language"
)

var (
	defaultIgnore = ignore(colltab.Tertiary)
	defaultTable  = getTable(locales[0])
)

func TestOptions(t *testing.T) ***REMOVED***
	for i, tt := range []struct ***REMOVED***
		in  []Option
		out options
	***REMOVED******REMOVED***
		0: ***REMOVED***
			out: options***REMOVED***
				ignore: defaultIgnore,
			***REMOVED***,
		***REMOVED***,
		1: ***REMOVED***
			in: []Option***REMOVED***IgnoreDiacritics***REMOVED***,
			out: options***REMOVED***
				ignore: [colltab.NumLevels]bool***REMOVED***false, true, false, true, true***REMOVED***,
			***REMOVED***,
		***REMOVED***,
		2: ***REMOVED***
			in: []Option***REMOVED***IgnoreCase, IgnoreDiacritics***REMOVED***,
			out: options***REMOVED***
				ignore: ignore(colltab.Primary),
			***REMOVED***,
		***REMOVED***,
		3: ***REMOVED***
			in: []Option***REMOVED***ignoreDiacritics, IgnoreWidth***REMOVED***,
			out: options***REMOVED***
				ignore:    ignore(colltab.Primary),
				caseLevel: true,
			***REMOVED***,
		***REMOVED***,
		4: ***REMOVED***
			in: []Option***REMOVED***IgnoreWidth, ignoreDiacritics***REMOVED***,
			out: options***REMOVED***
				ignore:    ignore(colltab.Primary),
				caseLevel: true,
			***REMOVED***,
		***REMOVED***,
		5: ***REMOVED***
			in: []Option***REMOVED***IgnoreCase, IgnoreWidth***REMOVED***,
			out: options***REMOVED***
				ignore: ignore(colltab.Secondary),
			***REMOVED***,
		***REMOVED***,
		6: ***REMOVED***
			in: []Option***REMOVED***IgnoreCase, IgnoreWidth, Loose***REMOVED***,
			out: options***REMOVED***
				ignore: ignore(colltab.Primary),
			***REMOVED***,
		***REMOVED***,
		7: ***REMOVED***
			in: []Option***REMOVED***Force, IgnoreCase, IgnoreWidth, Loose***REMOVED***,
			out: options***REMOVED***
				ignore: [colltab.NumLevels]bool***REMOVED***false, true, true, true, false***REMOVED***,
			***REMOVED***,
		***REMOVED***,
		8: ***REMOVED***
			in: []Option***REMOVED***IgnoreDiacritics, IgnoreCase***REMOVED***,
			out: options***REMOVED***
				ignore: ignore(colltab.Primary),
			***REMOVED***,
		***REMOVED***,
		9: ***REMOVED***
			in: []Option***REMOVED***Numeric***REMOVED***,
			out: options***REMOVED***
				ignore:  defaultIgnore,
				numeric: true,
			***REMOVED***,
		***REMOVED***,
		10: ***REMOVED***
			in: []Option***REMOVED***OptionsFromTag(language.MustParse("und-u-ks-level1"))***REMOVED***,
			out: options***REMOVED***
				ignore: ignore(colltab.Primary),
			***REMOVED***,
		***REMOVED***,
		11: ***REMOVED***
			in: []Option***REMOVED***OptionsFromTag(language.MustParse("und-u-ks-level4"))***REMOVED***,
			out: options***REMOVED***
				ignore: ignore(colltab.Quaternary),
			***REMOVED***,
		***REMOVED***,
		12: ***REMOVED***
			in:  []Option***REMOVED***OptionsFromTag(language.MustParse("und-u-ks-identic"))***REMOVED***,
			out: options***REMOVED******REMOVED***,
		***REMOVED***,
		13: ***REMOVED***
			in: []Option***REMOVED***
				OptionsFromTag(language.MustParse("und-u-kn-true-kb-true-kc-true")),
			***REMOVED***,
			out: options***REMOVED***
				ignore:    defaultIgnore,
				caseLevel: true,
				backwards: true,
				numeric:   true,
			***REMOVED***,
		***REMOVED***,
		14: ***REMOVED***
			in: []Option***REMOVED***
				OptionsFromTag(language.MustParse("und-u-kn-true-kb-true-kc-true")),
				OptionsFromTag(language.MustParse("und-u-kn-false-kb-false-kc-false")),
			***REMOVED***,
			out: options***REMOVED***
				ignore: defaultIgnore,
			***REMOVED***,
		***REMOVED***,
		15: ***REMOVED***
			in: []Option***REMOVED***
				OptionsFromTag(language.MustParse("und-u-kn-true-kb-true-kc-true")),
				OptionsFromTag(language.MustParse("und-u-kn-foo-kb-foo-kc-foo")),
			***REMOVED***,
			out: options***REMOVED***
				ignore:    defaultIgnore,
				caseLevel: true,
				backwards: true,
				numeric:   true,
			***REMOVED***,
		***REMOVED***,
		16: ***REMOVED*** // Normal options take precedence over tag options.
			in: []Option***REMOVED***
				Numeric, IgnoreCase,
				OptionsFromTag(language.MustParse("und-u-kn-false-kc-true")),
			***REMOVED***,
			out: options***REMOVED***
				ignore:    ignore(colltab.Secondary),
				caseLevel: false,
				numeric:   true,
			***REMOVED***,
		***REMOVED***,
		17: ***REMOVED***
			in: []Option***REMOVED***
				OptionsFromTag(language.MustParse("und-u-ka-shifted")),
			***REMOVED***,
			out: options***REMOVED***
				ignore:    defaultIgnore,
				alternate: altShifted,
			***REMOVED***,
		***REMOVED***,
		18: ***REMOVED***
			in: []Option***REMOVED***
				OptionsFromTag(language.MustParse("und-u-ka-blanked")),
			***REMOVED***,
			out: options***REMOVED***
				ignore:    defaultIgnore,
				alternate: altBlanked,
			***REMOVED***,
		***REMOVED***,
		19: ***REMOVED***
			in: []Option***REMOVED***
				OptionsFromTag(language.MustParse("und-u-ka-posix")),
			***REMOVED***,
			out: options***REMOVED***
				ignore:    defaultIgnore,
				alternate: altShiftTrimmed,
			***REMOVED***,
		***REMOVED***,
	***REMOVED*** ***REMOVED***
		c := newCollator(defaultTable)
		c.t = nil
		c.variableTop = 0
		c.f = 0

		c.setOptions(tt.in)
		if !reflect.DeepEqual(c.options, tt.out) ***REMOVED***
			t.Errorf("%d: got %v; want %v", i, c.options, tt.out)
		***REMOVED***
	***REMOVED***
***REMOVED***

func TestAlternateSortTypes(t *testing.T) ***REMOVED***
	testCases := []struct ***REMOVED***
		lang string
		in   []string
		want []string
	***REMOVED******REMOVED******REMOVED***
		lang: "zh,cmn,zh-Hant-u-co-pinyin,zh-HK-u-co-pinyin,zh-pinyin",
		in:   []string***REMOVED***"爸爸", "妈妈", "儿子", "女儿"***REMOVED***,
		want: []string***REMOVED***"爸爸", "儿子", "妈妈", "女儿"***REMOVED***,
	***REMOVED***, ***REMOVED***
		lang: "zh-Hant,zh-u-co-stroke,zh-Hant-u-co-stroke",
		in:   []string***REMOVED***"爸爸", "妈妈", "儿子", "女儿"***REMOVED***,
		want: []string***REMOVED***"儿子", "女儿", "妈妈", "爸爸"***REMOVED***,
	***REMOVED******REMOVED***
	for _, tc := range testCases ***REMOVED***
		for _, tag := range strings.Split(tc.lang, ",") ***REMOVED***
			got := append([]string***REMOVED******REMOVED***, tc.in...)
			New(language.MustParse(tag)).SortStrings(got)
			if !reflect.DeepEqual(got, tc.want) ***REMOVED***
				t.Errorf("New(%s).SortStrings(%v) = %v; want %v", tag, tc.in, got, tc.want)
			***REMOVED***
		***REMOVED***
	***REMOVED***
***REMOVED***
