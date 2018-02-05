package cldr

import (
	"fmt"
	"log"
	"reflect"
	"testing"
)

func failOnError(err error) ***REMOVED***
	if err != nil ***REMOVED***
		log.Panic(err)
	***REMOVED***
***REMOVED***

func data() *CLDR ***REMOVED***
	d := Decoder***REMOVED******REMOVED***
	data, err := d.Decode(testLoader***REMOVED******REMOVED***)
	failOnError(err)
	return data
***REMOVED***

type h struct ***REMOVED***
	A string `xml:"ha,attr"`
	E string `xml:"he"`
	D string `xml:",chardata"`
	X string
***REMOVED***

type fieldTest struct ***REMOVED***
	Common
	To  string `xml:"to,attr"`
	Key string `xml:"key,attr"`
	E   string `xml:"e"`
	D   string `xml:",chardata"`
	X   string
	h
***REMOVED***

var testStruct = fieldTest***REMOVED***
	Common: Common***REMOVED***
		name: "mapping", // exclude "type" as distinguishing attribute
		Type: "foo",
		Alt:  "foo",
	***REMOVED***,
	To:  "nyc",
	Key: "k",
	E:   "E",
	D:   "D",
	h: h***REMOVED***
		A: "A",
		E: "E",
		D: "D",
	***REMOVED***,
***REMOVED***

func TestIter(t *testing.T) ***REMOVED***
	tests := map[string]string***REMOVED***
		"Type":  "foo",
		"Alt":   "foo",
		"To":    "nyc",
		"A":     "A",
		"Alias": "<nil>",
	***REMOVED***
	k := 0
	for i := iter(reflect.ValueOf(testStruct)); !i.done(); i.next() ***REMOVED***
		v := i.value()
		if v.Kind() == reflect.Ptr && v.Elem().Kind() == reflect.String ***REMOVED***
			v = v.Elem()
		***REMOVED***
		name := i.field().Name
		if w, ok := tests[name]; ok ***REMOVED***
			s := fmt.Sprint(v.Interface())
			if w != s ***REMOVED***
				t.Errorf("value: found %q; want %q", w, s)
			***REMOVED***
			delete(tests, name)
		***REMOVED***
		k++
	***REMOVED***
	if len(tests) != 0 ***REMOVED***
		t.Errorf("missing fields: %v", tests)
	***REMOVED***
***REMOVED***

func TestFindField(t *testing.T) ***REMOVED***
	tests := []struct ***REMOVED***
		name, val string
		exist     bool
	***REMOVED******REMOVED***
		***REMOVED***"type", "foo", true***REMOVED***,
		***REMOVED***"alt", "foo", true***REMOVED***,
		***REMOVED***"to", "nyc", true***REMOVED***,
		***REMOVED***"he", "E", true***REMOVED***,
		***REMOVED***"q", "", false***REMOVED***,
	***REMOVED***
	vf := reflect.ValueOf(testStruct)
	for i, tt := range tests ***REMOVED***
		v, err := findField(vf, tt.name)
		if (err == nil) != tt.exist ***REMOVED***
			t.Errorf("%d: field %q present is %v; want %v", i, tt.name, err == nil, tt.exist)
		***REMOVED*** else if tt.exist ***REMOVED***
			if v.Kind() == reflect.Ptr ***REMOVED***
				if v.IsNil() ***REMOVED***
					continue
				***REMOVED***
				v = v.Elem()
			***REMOVED***
			if v.String() != tt.val ***REMOVED***
				t.Errorf("%d: found value %q; want %q", i, v.String(), tt.val)
			***REMOVED***
		***REMOVED***
	***REMOVED***
***REMOVED***

var keyTests = []struct ***REMOVED***
	exclude []string
	key     string
***REMOVED******REMOVED***
	***REMOVED***[]string***REMOVED******REMOVED***, "alt=foo;key=k;to=nyc"***REMOVED***,
	***REMOVED***[]string***REMOVED***"type"***REMOVED***, "alt=foo;key=k;to=nyc"***REMOVED***,
	***REMOVED***[]string***REMOVED***"choice"***REMOVED***, "alt=foo;key=k;to=nyc"***REMOVED***,
	***REMOVED***[]string***REMOVED***"alt"***REMOVED***, "key=k;to=nyc"***REMOVED***,
	***REMOVED***[]string***REMOVED***"a"***REMOVED***, "alt=foo;key=k;to=nyc"***REMOVED***,
	***REMOVED***[]string***REMOVED***"to"***REMOVED***, "alt=foo;key=k"***REMOVED***,
	***REMOVED***[]string***REMOVED***"alt", "to"***REMOVED***, "key=k"***REMOVED***,
	***REMOVED***[]string***REMOVED***"alt", "to", "key"***REMOVED***, ""***REMOVED***,
***REMOVED***

func TestAttrKey(t *testing.T) ***REMOVED***
	v := reflect.ValueOf(&testStruct)
	for i, tt := range keyTests ***REMOVED***
		key := attrKey(v, tt.exclude...)
		if key != tt.key ***REMOVED***
			t.Errorf("%d: found %q, want %q", i, key, tt.key)
		***REMOVED***
	***REMOVED***
***REMOVED***

func TestKey(t *testing.T) ***REMOVED***
	for i, tt := range keyTests ***REMOVED***
		key := Key(&testStruct, tt.exclude...)
		if key != tt.key ***REMOVED***
			t.Errorf("%d: found %q, want %q", i, key, tt.key)
		***REMOVED***
	***REMOVED***
***REMOVED***

func testEnclosing(t *testing.T, x *LDML, name string) ***REMOVED***
	eq := func(a, b Elem, i int) ***REMOVED***
		for ; i > 0; i-- ***REMOVED***
			b = b.enclosing()
		***REMOVED***
		if a != b ***REMOVED***
			t.Errorf("%s: found path %q, want %q", name, getPath(a), getPath(b))
		***REMOVED***
	***REMOVED***
	eq(x, x, 0)
	eq(x, x.Identity, 1)
	eq(x, x.Dates.Calendars, 2)
	eq(x, x.Dates.Calendars.Calendar[0], 3)
	eq(x, x.Dates.Calendars.Calendar[1], 3)
	//eq(x, x.Dates.Calendars.Calendar[0].Months, 4)
	eq(x, x.Dates.Calendars.Calendar[1].Months, 4)
***REMOVED***

func TestEnclosing(t *testing.T) ***REMOVED***
	testEnclosing(t, data().RawLDML("de"), "enclosing-raw")
	de, _ := data().LDML("de")
	testEnclosing(t, de, "enclosing")
***REMOVED***

func TestDeepCopy(t *testing.T) ***REMOVED***
	eq := func(have, want string) ***REMOVED***
		if have != want ***REMOVED***
			t.Errorf("found %q; want %q", have, want)
		***REMOVED***
	***REMOVED***
	x, _ := data().LDML("de")
	vc := deepCopy(reflect.ValueOf(x))
	c := vc.Interface().(*LDML)
	linkEnclosing(nil, c)
	if x == c ***REMOVED***
		t.Errorf("did not copy")
	***REMOVED***

	eq(c.name, "ldml")
	eq(c.Dates.name, "dates")
	testEnclosing(t, c, "deepCopy")
***REMOVED***

type getTest struct ***REMOVED***
	loc     string
	path    string
	field   string // used in combination with length
	data    string
	altData string // used for buddhist calendar if value != ""
	typ     string
	length  int
	missing bool
***REMOVED***

const (
	budMon = "dates/calendars/calendar[@type='buddhist']/months/"
	chnMon = "dates/calendars/calendar[@type='chinese']/months/"
	greMon = "dates/calendars/calendar[@type='gregorian']/months/"
)

func monthVal(path, context, width string, month int) string ***REMOVED***
	const format = "%s/monthContext[@type='%s']/monthWidth[@type='%s']/month[@type='%d']"
	return fmt.Sprintf(format, path, context, width, month)
***REMOVED***

var rootGetTests = []getTest***REMOVED***
	***REMOVED***loc: "root", path: "identity/language", typ: "root"***REMOVED***,
	***REMOVED***loc: "root", path: "characters/moreInformation", data: "?"***REMOVED***,
	***REMOVED***loc: "root", path: "characters", field: "exemplarCharacters", length: 3***REMOVED***,
	***REMOVED***loc: "root", path: greMon, field: "monthContext", length: 2***REMOVED***,
	***REMOVED***loc: "root", path: greMon + "monthContext[@type='format']/monthWidth[@type='narrow']", field: "month", length: 4***REMOVED***,
	***REMOVED***loc: "root", path: greMon + "monthContext[@type='stand-alone']/monthWidth[@type='wide']", field: "month", length: 4***REMOVED***,
	// unescaping character data
	***REMOVED***loc: "root", path: "characters/exemplarCharacters[@type='punctuation']", data: `[\- ‐ – — … ' ‘ ‚ " “ „ \& #]`***REMOVED***,
	// default resolution
	***REMOVED***loc: "root", path: "dates/calendars/calendar", typ: "gregorian"***REMOVED***,
	// alias resolution
	***REMOVED***loc: "root", path: budMon, field: "monthContext", length: 2***REMOVED***,
	// crossing but non-circular alias resolution
	***REMOVED***loc: "root", path: budMon + "monthContext[@type='format']/monthWidth[@type='narrow']", field: "month", length: 4***REMOVED***,
	***REMOVED***loc: "root", path: budMon + "monthContext[@type='stand-alone']/monthWidth[@type='wide']", field: "month", length: 4***REMOVED***,
	***REMOVED***loc: "root", path: monthVal(greMon, "format", "wide", 1), data: "11"***REMOVED***,
	***REMOVED***loc: "root", path: monthVal(greMon, "format", "narrow", 2), data: "2"***REMOVED***,
	***REMOVED***loc: "root", path: monthVal(greMon, "stand-alone", "wide", 3), data: "33"***REMOVED***,
	***REMOVED***loc: "root", path: monthVal(greMon, "stand-alone", "narrow", 4), data: "4"***REMOVED***,
	***REMOVED***loc: "root", path: monthVal(budMon, "format", "wide", 1), data: "11"***REMOVED***,
	***REMOVED***loc: "root", path: monthVal(budMon, "format", "narrow", 2), data: "2"***REMOVED***,
	***REMOVED***loc: "root", path: monthVal(budMon, "stand-alone", "wide", 3), data: "33"***REMOVED***,
	***REMOVED***loc: "root", path: monthVal(budMon, "stand-alone", "narrow", 4), data: "4"***REMOVED***,
***REMOVED***

// 19
var deGetTests = []getTest***REMOVED***
	***REMOVED***loc: "de", path: "identity/language", typ: "de"***REMOVED***,
	***REMOVED***loc: "de", path: "posix", length: 2***REMOVED***,
	***REMOVED***loc: "de", path: "characters", field: "exemplarCharacters", length: 4***REMOVED***,
	***REMOVED***loc: "de", path: "characters/exemplarCharacters[@type='auxiliary']", data: `[á à ă]`***REMOVED***,
	// identity is a blocking element, so de should not inherit generation from root.
	***REMOVED***loc: "de", path: "identity/generation", missing: true***REMOVED***,
	// default resolution
	***REMOVED***loc: "root", path: "dates/calendars/calendar", typ: "gregorian"***REMOVED***,

	// absolute path alias resolution
	***REMOVED***loc: "gsw", path: "posix", field: "messages", length: 1***REMOVED***,
	***REMOVED***loc: "gsw", path: "posix/messages/yesstr", data: "yes:y"***REMOVED***,
***REMOVED***

// 27(greMon) - 52(budMon) - 77(chnMon)
func calGetTests(s string) []getTest ***REMOVED***
	tests := []getTest***REMOVED***
		***REMOVED***loc: "de", path: s, length: 2***REMOVED***,
		***REMOVED***loc: "de", path: s + "monthContext[@type='format']/monthWidth[@type='wide']", field: "month", length: 5***REMOVED***,
		***REMOVED***loc: "de", path: monthVal(s, "format", "wide", 1), data: "11"***REMOVED***,
		***REMOVED***loc: "de", path: monthVal(s, "format", "wide", 2), data: "22"***REMOVED***,
		***REMOVED***loc: "de", path: monthVal(s, "format", "wide", 3), data: "Maerz", altData: "bbb"***REMOVED***,
		***REMOVED***loc: "de", path: monthVal(s, "format", "wide", 4), data: "April"***REMOVED***,
		***REMOVED***loc: "de", path: monthVal(s, "format", "wide", 5), data: "Mai"***REMOVED***,

		***REMOVED***loc: "de", path: s + "monthContext[@type='format']/monthWidth[@type='narrow']", field: "month", length: 5***REMOVED***,
		***REMOVED***loc: "de", path: monthVal(s, "format", "narrow", 1), data: "1"***REMOVED***,
		***REMOVED***loc: "de", path: monthVal(s, "format", "narrow", 2), data: "2"***REMOVED***,
		***REMOVED***loc: "de", path: monthVal(s, "format", "narrow", 3), data: "M", altData: "BBB"***REMOVED***,
		***REMOVED***loc: "de", path: monthVal(s, "format", "narrow", 4), data: "A"***REMOVED***,
		***REMOVED***loc: "de", path: monthVal(s, "format", "narrow", 5), data: "m"***REMOVED***,

		***REMOVED***loc: "de", path: s + "monthContext[@type='stand-alone']/monthWidth[@type='wide']", field: "month", length: 5***REMOVED***,
		***REMOVED***loc: "de", path: monthVal(s, "stand-alone", "wide", 1), data: "11"***REMOVED***,
		***REMOVED***loc: "de", path: monthVal(s, "stand-alone", "wide", 2), data: "22"***REMOVED***,
		***REMOVED***loc: "de", path: monthVal(s, "stand-alone", "wide", 3), data: "Maerz", altData: "bbb"***REMOVED***,
		***REMOVED***loc: "de", path: monthVal(s, "stand-alone", "wide", 4), data: "april"***REMOVED***,
		***REMOVED***loc: "de", path: monthVal(s, "stand-alone", "wide", 5), data: "mai"***REMOVED***,

		***REMOVED***loc: "de", path: s + "monthContext[@type='stand-alone']/monthWidth[@type='narrow']", field: "month", length: 5***REMOVED***,
		***REMOVED***loc: "de", path: monthVal(s, "stand-alone", "narrow", 1), data: "1"***REMOVED***,
		***REMOVED***loc: "de", path: monthVal(s, "stand-alone", "narrow", 2), data: "2"***REMOVED***,
		***REMOVED***loc: "de", path: monthVal(s, "stand-alone", "narrow", 3), data: "m"***REMOVED***,
		***REMOVED***loc: "de", path: monthVal(s, "stand-alone", "narrow", 4), data: "4"***REMOVED***,
		***REMOVED***loc: "de", path: monthVal(s, "stand-alone", "narrow", 5), data: "m"***REMOVED***,
	***REMOVED***
	if s == budMon ***REMOVED***
		for i, t := range tests ***REMOVED***
			if t.altData != "" ***REMOVED***
				tests[i].data = t.altData
			***REMOVED***
		***REMOVED***
	***REMOVED***
	return tests
***REMOVED***

var getTests = append(rootGetTests,
	append(deGetTests,
		append(calGetTests(greMon),
			append(calGetTests(budMon),
				calGetTests(chnMon)...)...)...)...)

func TestPath(t *testing.T) ***REMOVED***
	d := data()
	for i, tt := range getTests ***REMOVED***
		x, _ := d.LDML(tt.loc)
		e, err := walkXPath(x, tt.path)
		if err != nil ***REMOVED***
			if !tt.missing ***REMOVED***
				t.Errorf("%d:error: %v %v", i, err, tt.missing)
			***REMOVED***
			continue
		***REMOVED***
		if tt.missing ***REMOVED***
			t.Errorf("%d: missing is %v; want %v", i, e == nil, tt.missing)
			continue
		***REMOVED***
		if tt.data != "" && e.GetCommon().Data() != tt.data ***REMOVED***
			t.Errorf("%d: data is %v; want %v", i, e.GetCommon().Data(), tt.data)
			continue
		***REMOVED***
		if tt.typ != "" && e.GetCommon().Type != tt.typ ***REMOVED***
			t.Errorf("%d: type is %v; want %v", i, e.GetCommon().Type, tt.typ)
			continue
		***REMOVED***
		if tt.field != "" ***REMOVED***
			slice, _ := findField(reflect.ValueOf(e), tt.field)
			if slice.Len() != tt.length ***REMOVED***
				t.Errorf("%d: length is %v; want %v", i, slice.Len(), tt.length)
				continue
			***REMOVED***
		***REMOVED***
	***REMOVED***
***REMOVED***

func TestGet(t *testing.T) ***REMOVED***
	d := data()
	for i, tt := range getTests ***REMOVED***
		x, _ := d.LDML(tt.loc)
		e, err := Get(x, tt.path)
		if err != nil ***REMOVED***
			if !tt.missing ***REMOVED***
				t.Errorf("%d:error: %v %v", i, err, tt.missing)
			***REMOVED***
			continue
		***REMOVED***
		if tt.missing ***REMOVED***
			t.Errorf("%d: missing is %v; want %v", i, e == nil, tt.missing)
			continue
		***REMOVED***
		if tt.data != "" && e.GetCommon().Data() != tt.data ***REMOVED***
			t.Errorf("%d: data is %v; want %v", i, e.GetCommon().Data(), tt.data)
			continue
		***REMOVED***
		if tt.typ != "" && e.GetCommon().Type != tt.typ ***REMOVED***
			t.Errorf("%d: type is %v; want %v", i, e.GetCommon().Type, tt.typ)
			continue
		***REMOVED***
		if tt.field != "" ***REMOVED***
			slice, _ := findField(reflect.ValueOf(e), tt.field)
			if slice.Len() != tt.length ***REMOVED***
				t.Errorf("%d: length is %v; want %v", i, slice.Len(), tt.length)
				continue
			***REMOVED***
		***REMOVED***
	***REMOVED***
***REMOVED***
