// Copyright 2017 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package plural

import (
	"fmt"
	"strings"
	"testing"

	"golang.org/x/text/internal/catmsg"
	"golang.org/x/text/language"
	"golang.org/x/text/message/catalog"
)

func TestSelect(t *testing.T) ***REMOVED***
	lang := language.English
	type test struct ***REMOVED***
		arg    interface***REMOVED******REMOVED***
		result string
		err    string
	***REMOVED***
	testCases := []struct ***REMOVED***
		desc  string
		msg   catalog.Message
		err   string
		tests []test
	***REMOVED******REMOVED******REMOVED***
		desc: "basic",
		msg:  Selectf(1, "%d", "one", "foo", "other", "bar"),
		tests: []test***REMOVED***
			***REMOVED***arg: 0, result: "bar"***REMOVED***,
			***REMOVED***arg: 1, result: "foo"***REMOVED***,
			***REMOVED***arg: 2, result: "bar"***REMOVED***,
			***REMOVED***arg: opposite(1), result: "bar"***REMOVED***,
			***REMOVED***arg: opposite(2), result: "foo"***REMOVED***,
			***REMOVED***arg: "unknown", result: "bar"***REMOVED***, // other
		***REMOVED***,
	***REMOVED***, ***REMOVED***
		desc: "comparisons",
		msg: Selectf(1, "%d",
			"=0", "zero",
			"=1", "one",
			"one", "cannot match", // never matches
			"<5", "<5", // never matches
			"=5", "=5",
			Other, "other"),
		tests: []test***REMOVED***
			***REMOVED***arg: 0, result: "zero"***REMOVED***,
			***REMOVED***arg: 1, result: "one"***REMOVED***,
			***REMOVED***arg: 2, result: "<5"***REMOVED***,
			***REMOVED***arg: 4, result: "<5"***REMOVED***,
			***REMOVED***arg: 5, result: "=5"***REMOVED***,
			***REMOVED***arg: 6, result: "other"***REMOVED***,
			***REMOVED***arg: "unknown", result: "other"***REMOVED***,
		***REMOVED***,
	***REMOVED***, ***REMOVED***
		desc: "fractions",
		msg:  Selectf(1, "%.2f", "one", "foo", "other", "bar"),
		tests: []test***REMOVED***
			// fractions are always plural in english
			***REMOVED***arg: 0, result: "bar"***REMOVED***,
			***REMOVED***arg: 1, result: "bar"***REMOVED***,
		***REMOVED***,
	***REMOVED***, ***REMOVED***
		desc: "decimal without fractions",
		msg:  Selectf(1, "%.0f", "one", "foo", "other", "bar"),
		tests: []test***REMOVED***
			// fractions are always plural in english
			***REMOVED***arg: 0, result: "bar"***REMOVED***,
			***REMOVED***arg: 1, result: "foo"***REMOVED***,
		***REMOVED***,
	***REMOVED***, ***REMOVED***
		desc: "scientific",
		msg:  Selectf(1, "%.0e", "one", "foo", "other", "bar"),
		tests: []test***REMOVED***
			***REMOVED***arg: 0, result: "bar"***REMOVED***,
			***REMOVED***arg: 1, result: "foo"***REMOVED***,
		***REMOVED***,
	***REMOVED***, ***REMOVED***
		desc: "variable",
		msg:  Selectf(1, "%.1g", "one", "foo", "other", "bar"),
		tests: []test***REMOVED***
			// fractions are always plural in english
			***REMOVED***arg: 0, result: "bar"***REMOVED***,
			***REMOVED***arg: 1, result: "foo"***REMOVED***,
			***REMOVED***arg: 2, result: "bar"***REMOVED***,
		***REMOVED***,
	***REMOVED***, ***REMOVED***
		desc: "default",
		msg:  Selectf(1, "", "one", "foo", "other", "bar"),
		tests: []test***REMOVED***
			***REMOVED***arg: 0, result: "bar"***REMOVED***,
			***REMOVED***arg: 1, result: "foo"***REMOVED***,
			***REMOVED***arg: 2, result: "bar"***REMOVED***,
			***REMOVED***arg: 1.0, result: "bar"***REMOVED***,
		***REMOVED***,
	***REMOVED***, ***REMOVED***
		desc: "nested",
		msg:  Selectf(1, "", "other", Selectf(2, "", "one", "foo", "other", "bar")),
		tests: []test***REMOVED***
			***REMOVED***arg: 0, result: "bar"***REMOVED***,
			***REMOVED***arg: 1, result: "foo"***REMOVED***,
			***REMOVED***arg: 2, result: "bar"***REMOVED***,
		***REMOVED***,
	***REMOVED***, ***REMOVED***
		desc:  "arg unavailable",
		msg:   Selectf(100, "%.2f", "one", "foo", "other", "bar"),
		tests: []test***REMOVED******REMOVED***arg: 1, result: "bar"***REMOVED******REMOVED***,
	***REMOVED***, ***REMOVED***
		desc:  "no match",
		msg:   Selectf(1, "%.2f", "one", "foo"),
		tests: []test***REMOVED******REMOVED***arg: 0, result: "bar", err: catmsg.ErrNoMatch.Error()***REMOVED******REMOVED***,
	***REMOVED***, ***REMOVED***
		desc: "error invalid form",
		err:  `invalid plural form "excessive"`,
		msg:  Selectf(1, "%d", "excessive", "foo"),
	***REMOVED***, ***REMOVED***
		desc: "error form not used by language",
		err:  `form "many" not supported for language "en"`,
		msg:  Selectf(1, "%d", "many", "foo"),
	***REMOVED***, ***REMOVED***
		desc: "error invalid selector",
		err:  `selector of type int; want string or Form`,
		msg:  Selectf(1, "%d", 1, "foo"),
	***REMOVED***, ***REMOVED***
		desc: "error missing message",
		err:  `no message defined for selector one`,
		msg:  Selectf(1, "%d", "one"),
	***REMOVED***, ***REMOVED***
		desc: "error invalid number",
		err:  `invalid number in selector "<1.00"`,
		msg:  Selectf(1, "%d", "<1.00"),
	***REMOVED***, ***REMOVED***
		desc: "error empty selector",
		err:  `empty selector`,
		msg:  Selectf(1, "%d", "", "foo"),
	***REMOVED***, ***REMOVED***
		desc: "error invalid message",
		err:  `message of type int; must be string or catalog.Message`,
		msg:  Selectf(1, "%d", "one", 3),
	***REMOVED***, ***REMOVED***
		desc: "nested error",
		err:  `empty selector`,
		msg:  Selectf(1, "", "other", Selectf(2, "", "")),
	***REMOVED******REMOVED***
	for _, tc := range testCases ***REMOVED***
		t.Run(tc.desc, func(t *testing.T) ***REMOVED***
			data, err := catmsg.Compile(lang, nil, tc.msg)
			chkError(t, err, tc.err)
			for _, tx := range tc.tests ***REMOVED***
				t.Run(fmt.Sprint(tx.arg), func(t *testing.T) ***REMOVED***
					r := renderer***REMOVED***arg: tx.arg***REMOVED***
					d := catmsg.NewDecoder(lang, &r, nil)
					err := d.Execute(data)
					chkError(t, err, tx.err)
					if r.result != tx.result ***REMOVED***
						t.Errorf("got %q; want %q", r.result, tx.result)
					***REMOVED***
				***REMOVED***)
			***REMOVED***
		***REMOVED***)
	***REMOVED***
***REMOVED***

func chkError(t *testing.T, got error, want string) ***REMOVED***
	if (got == nil && want != "") ||
		(got != nil && (want == "" || !strings.Contains(got.Error(), want))) ***REMOVED***
		t.Fatalf("got %v; want %v", got, want)
	***REMOVED***
	if got != nil ***REMOVED***
		t.SkipNow()
	***REMOVED***
***REMOVED***

type renderer struct ***REMOVED***
	arg    interface***REMOVED******REMOVED***
	result string
***REMOVED***

func (r *renderer) Render(s string) ***REMOVED*** r.result += s ***REMOVED***
func (r *renderer) Arg(i int) interface***REMOVED******REMOVED*** ***REMOVED***
	if i > 10 ***REMOVED*** // Allow testing "arg unavailable" path
		return nil
	***REMOVED***
	return r.arg
***REMOVED***

type opposite int

func (o opposite) PluralForm(lang language.Tag, scale int) (Form, int) ***REMOVED***
	if o == 1 ***REMOVED***
		return Other, 1
	***REMOVED***
	return One, int(o)
***REMOVED***
