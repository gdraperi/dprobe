package filters

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParseArgs(t *testing.T) ***REMOVED***
	// equivalent of `docker ps -f 'created=today' -f 'image.name=ubuntu*' -f 'image.name=*untu'`
	flagArgs := []string***REMOVED***
		"created=today",
		"image.name=ubuntu*",
		"image.name=*untu",
	***REMOVED***
	var (
		args = NewArgs()
		err  error
	)

	for i := range flagArgs ***REMOVED***
		args, err = ParseFlag(flagArgs[i], args)
		require.NoError(t, err)
	***REMOVED***
	assert.Len(t, args.Get("created"), 1)
	assert.Len(t, args.Get("image.name"), 2)
***REMOVED***

func TestParseArgsEdgeCase(t *testing.T) ***REMOVED***
	var args Args
	args, err := ParseFlag("", args)
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
	if args.Len() != 0 ***REMOVED***
		t.Fatalf("Expected an empty Args (map), got %v", args)
	***REMOVED***
	if args, err = ParseFlag("anything", args); err == nil || err != ErrBadFormat ***REMOVED***
		t.Fatalf("Expected ErrBadFormat, got %v", err)
	***REMOVED***
***REMOVED***

func TestToJSON(t *testing.T) ***REMOVED***
	fields := map[string]map[string]bool***REMOVED***
		"created":    ***REMOVED***"today": true***REMOVED***,
		"image.name": ***REMOVED***"ubuntu*": true, "*untu": true***REMOVED***,
	***REMOVED***
	a := Args***REMOVED***fields: fields***REMOVED***

	_, err := ToJSON(a)
	if err != nil ***REMOVED***
		t.Errorf("failed to marshal the filters: %s", err)
	***REMOVED***
***REMOVED***

func TestToParamWithVersion(t *testing.T) ***REMOVED***
	fields := map[string]map[string]bool***REMOVED***
		"created":    ***REMOVED***"today": true***REMOVED***,
		"image.name": ***REMOVED***"ubuntu*": true, "*untu": true***REMOVED***,
	***REMOVED***
	a := Args***REMOVED***fields: fields***REMOVED***

	str1, err := ToParamWithVersion("1.21", a)
	if err != nil ***REMOVED***
		t.Errorf("failed to marshal the filters with version < 1.22: %s", err)
	***REMOVED***
	str2, err := ToParamWithVersion("1.22", a)
	if err != nil ***REMOVED***
		t.Errorf("failed to marshal the filters with version >= 1.22: %s", err)
	***REMOVED***
	if str1 != `***REMOVED***"created":["today"],"image.name":["*untu","ubuntu*"]***REMOVED***` &&
		str1 != `***REMOVED***"created":["today"],"image.name":["ubuntu*","*untu"]***REMOVED***` ***REMOVED***
		t.Errorf("incorrectly marshaled the filters: %s", str1)
	***REMOVED***
	if str2 != `***REMOVED***"created":***REMOVED***"today":true***REMOVED***,"image.name":***REMOVED***"*untu":true,"ubuntu*":true***REMOVED******REMOVED***` &&
		str2 != `***REMOVED***"created":***REMOVED***"today":true***REMOVED***,"image.name":***REMOVED***"ubuntu*":true,"*untu":true***REMOVED******REMOVED***` ***REMOVED***
		t.Errorf("incorrectly marshaled the filters: %s", str2)
	***REMOVED***
***REMOVED***

func TestFromJSON(t *testing.T) ***REMOVED***
	invalids := []string***REMOVED***
		"anything",
		"['a','list']",
		"***REMOVED***'key': 'value'***REMOVED***",
		`***REMOVED***"key": "value"***REMOVED***`,
	***REMOVED***
	valid := map[*Args][]string***REMOVED***
		***REMOVED***fields: map[string]map[string]bool***REMOVED***"key": ***REMOVED***"value": true***REMOVED******REMOVED******REMOVED***: ***REMOVED***
			`***REMOVED***"key": ["value"]***REMOVED***`,
			`***REMOVED***"key": ***REMOVED***"value": true***REMOVED******REMOVED***`,
		***REMOVED***,
		***REMOVED***fields: map[string]map[string]bool***REMOVED***"key": ***REMOVED***"value1": true, "value2": true***REMOVED******REMOVED******REMOVED***: ***REMOVED***
			`***REMOVED***"key": ["value1", "value2"]***REMOVED***`,
			`***REMOVED***"key": ***REMOVED***"value1": true, "value2": true***REMOVED******REMOVED***`,
		***REMOVED***,
		***REMOVED***fields: map[string]map[string]bool***REMOVED***"key1": ***REMOVED***"value1": true***REMOVED***, "key2": ***REMOVED***"value2": true***REMOVED******REMOVED******REMOVED***: ***REMOVED***
			`***REMOVED***"key1": ["value1"], "key2": ["value2"]***REMOVED***`,
			`***REMOVED***"key1": ***REMOVED***"value1": true***REMOVED***, "key2": ***REMOVED***"value2": true***REMOVED******REMOVED***`,
		***REMOVED***,
	***REMOVED***

	for _, invalid := range invalids ***REMOVED***
		if _, err := FromJSON(invalid); err == nil ***REMOVED***
			t.Fatalf("Expected an error with %v, got nothing", invalid)
		***REMOVED***
	***REMOVED***

	for expectedArgs, matchers := range valid ***REMOVED***
		for _, json := range matchers ***REMOVED***
			args, err := FromJSON(json)
			if err != nil ***REMOVED***
				t.Fatal(err)
			***REMOVED***
			if args.Len() != expectedArgs.Len() ***REMOVED***
				t.Fatalf("Expected %v, go %v", expectedArgs, args)
			***REMOVED***
			for key, expectedValues := range expectedArgs.fields ***REMOVED***
				values := args.Get(key)

				if len(values) != len(expectedValues) ***REMOVED***
					t.Fatalf("Expected %v, go %v", expectedArgs, args)
				***REMOVED***

				for _, v := range values ***REMOVED***
					if !expectedValues[v] ***REMOVED***
						t.Fatalf("Expected %v, go %v", expectedArgs, args)
					***REMOVED***
				***REMOVED***
			***REMOVED***
		***REMOVED***
	***REMOVED***
***REMOVED***

func TestEmpty(t *testing.T) ***REMOVED***
	a := Args***REMOVED******REMOVED***
	v, err := ToJSON(a)
	if err != nil ***REMOVED***
		t.Errorf("failed to marshal the filters: %s", err)
	***REMOVED***
	v1, err := FromJSON(v)
	if err != nil ***REMOVED***
		t.Errorf("%s", err)
	***REMOVED***
	if a.Len() != v1.Len() ***REMOVED***
		t.Error("these should both be empty sets")
	***REMOVED***
***REMOVED***

func TestArgsMatchKVListEmptySources(t *testing.T) ***REMOVED***
	args := NewArgs()
	if !args.MatchKVList("created", map[string]string***REMOVED******REMOVED***) ***REMOVED***
		t.Fatalf("Expected true for (%v,created), got true", args)
	***REMOVED***

	args = Args***REMOVED***map[string]map[string]bool***REMOVED***"created": ***REMOVED***"today": true***REMOVED******REMOVED******REMOVED***
	if args.MatchKVList("created", map[string]string***REMOVED******REMOVED***) ***REMOVED***
		t.Fatalf("Expected false for (%v,created), got true", args)
	***REMOVED***
***REMOVED***

func TestArgsMatchKVList(t *testing.T) ***REMOVED***
	// Not empty sources
	sources := map[string]string***REMOVED***
		"key1": "value1",
		"key2": "value2",
		"key3": "value3",
	***REMOVED***

	matches := map[*Args]string***REMOVED***
		***REMOVED******REMOVED***: "field",
		***REMOVED***map[string]map[string]bool***REMOVED***
			"created": ***REMOVED***"today": true***REMOVED***,
			"labels":  ***REMOVED***"key1": true***REMOVED******REMOVED***,
		***REMOVED***: "labels",
		***REMOVED***map[string]map[string]bool***REMOVED***
			"created": ***REMOVED***"today": true***REMOVED***,
			"labels":  ***REMOVED***"key1=value1": true***REMOVED******REMOVED***,
		***REMOVED***: "labels",
	***REMOVED***

	for args, field := range matches ***REMOVED***
		if !args.MatchKVList(field, sources) ***REMOVED***
			t.Fatalf("Expected true for %v on %v, got false", sources, args)
		***REMOVED***
	***REMOVED***

	differs := map[*Args]string***REMOVED***
		***REMOVED***map[string]map[string]bool***REMOVED***
			"created": ***REMOVED***"today": true***REMOVED******REMOVED***,
		***REMOVED***: "created",
		***REMOVED***map[string]map[string]bool***REMOVED***
			"created": ***REMOVED***"today": true***REMOVED***,
			"labels":  ***REMOVED***"key4": true***REMOVED******REMOVED***,
		***REMOVED***: "labels",
		***REMOVED***map[string]map[string]bool***REMOVED***
			"created": ***REMOVED***"today": true***REMOVED***,
			"labels":  ***REMOVED***"key1=value3": true***REMOVED******REMOVED***,
		***REMOVED***: "labels",
	***REMOVED***

	for args, field := range differs ***REMOVED***
		if args.MatchKVList(field, sources) ***REMOVED***
			t.Fatalf("Expected false for %v on %v, got true", sources, args)
		***REMOVED***
	***REMOVED***
***REMOVED***

func TestArgsMatch(t *testing.T) ***REMOVED***
	source := "today"

	matches := map[*Args]string***REMOVED***
		***REMOVED******REMOVED***: "field",
		***REMOVED***map[string]map[string]bool***REMOVED***
			"created": ***REMOVED***"today": true***REMOVED******REMOVED***,
		***REMOVED***: "today",
		***REMOVED***map[string]map[string]bool***REMOVED***
			"created": ***REMOVED***"to*": true***REMOVED******REMOVED***,
		***REMOVED***: "created",
		***REMOVED***map[string]map[string]bool***REMOVED***
			"created": ***REMOVED***"to(.*)": true***REMOVED******REMOVED***,
		***REMOVED***: "created",
		***REMOVED***map[string]map[string]bool***REMOVED***
			"created": ***REMOVED***"tod": true***REMOVED******REMOVED***,
		***REMOVED***: "created",
		***REMOVED***map[string]map[string]bool***REMOVED***
			"created": ***REMOVED***"anything": true, "to*": true***REMOVED******REMOVED***,
		***REMOVED***: "created",
	***REMOVED***

	for args, field := range matches ***REMOVED***
		assert.True(t, args.Match(field, source),
			"Expected field %s to match %s", field, source)
	***REMOVED***

	differs := map[*Args]string***REMOVED***
		***REMOVED***map[string]map[string]bool***REMOVED***
			"created": ***REMOVED***"tomorrow": true***REMOVED******REMOVED***,
		***REMOVED***: "created",
		***REMOVED***map[string]map[string]bool***REMOVED***
			"created": ***REMOVED***"to(day": true***REMOVED******REMOVED***,
		***REMOVED***: "created",
		***REMOVED***map[string]map[string]bool***REMOVED***
			"created": ***REMOVED***"tom(.*)": true***REMOVED******REMOVED***,
		***REMOVED***: "created",
		***REMOVED***map[string]map[string]bool***REMOVED***
			"created": ***REMOVED***"tom": true***REMOVED******REMOVED***,
		***REMOVED***: "created",
		***REMOVED***map[string]map[string]bool***REMOVED***
			"created": ***REMOVED***"today1": true***REMOVED***,
			"labels":  ***REMOVED***"today": true***REMOVED******REMOVED***,
		***REMOVED***: "created",
	***REMOVED***

	for args, field := range differs ***REMOVED***
		assert.False(t, args.Match(field, source),
			"Expected field %s to not match %s", field, source)
	***REMOVED***
***REMOVED***

func TestAdd(t *testing.T) ***REMOVED***
	f := NewArgs()
	f.Add("status", "running")
	v := f.fields["status"]
	if len(v) != 1 || !v["running"] ***REMOVED***
		t.Fatalf("Expected to include a running status, got %v", v)
	***REMOVED***

	f.Add("status", "paused")
	if len(v) != 2 || !v["paused"] ***REMOVED***
		t.Fatalf("Expected to include a paused status, got %v", v)
	***REMOVED***
***REMOVED***

func TestDel(t *testing.T) ***REMOVED***
	f := NewArgs()
	f.Add("status", "running")
	f.Del("status", "running")
	v := f.fields["status"]
	if v["running"] ***REMOVED***
		t.Fatal("Expected to not include a running status filter, got true")
	***REMOVED***
***REMOVED***

func TestLen(t *testing.T) ***REMOVED***
	f := NewArgs()
	if f.Len() != 0 ***REMOVED***
		t.Fatal("Expected to not include any field")
	***REMOVED***
	f.Add("status", "running")
	if f.Len() != 1 ***REMOVED***
		t.Fatal("Expected to include one field")
	***REMOVED***
***REMOVED***

func TestExactMatch(t *testing.T) ***REMOVED***
	f := NewArgs()

	if !f.ExactMatch("status", "running") ***REMOVED***
		t.Fatal("Expected to match `running` when there are no filters, got false")
	***REMOVED***

	f.Add("status", "running")
	f.Add("status", "pause*")

	if !f.ExactMatch("status", "running") ***REMOVED***
		t.Fatal("Expected to match `running` with one of the filters, got false")
	***REMOVED***

	if f.ExactMatch("status", "paused") ***REMOVED***
		t.Fatal("Expected to not match `paused` with one of the filters, got true")
	***REMOVED***
***REMOVED***

func TestOnlyOneExactMatch(t *testing.T) ***REMOVED***
	f := NewArgs()

	if !f.UniqueExactMatch("status", "running") ***REMOVED***
		t.Fatal("Expected to match `running` when there are no filters, got false")
	***REMOVED***

	f.Add("status", "running")

	if !f.UniqueExactMatch("status", "running") ***REMOVED***
		t.Fatal("Expected to match `running` with one of the filters, got false")
	***REMOVED***

	if f.UniqueExactMatch("status", "paused") ***REMOVED***
		t.Fatal("Expected to not match `paused` with one of the filters, got true")
	***REMOVED***

	f.Add("status", "pause")
	if f.UniqueExactMatch("status", "running") ***REMOVED***
		t.Fatal("Expected to not match only `running` with two filters, got true")
	***REMOVED***
***REMOVED***

func TestContains(t *testing.T) ***REMOVED***
	f := NewArgs()
	if f.Contains("status") ***REMOVED***
		t.Fatal("Expected to not contain a status key, got true")
	***REMOVED***
	f.Add("status", "running")
	if !f.Contains("status") ***REMOVED***
		t.Fatal("Expected to contain a status key, got false")
	***REMOVED***
***REMOVED***

func TestInclude(t *testing.T) ***REMOVED***
	f := NewArgs()
	if f.Include("status") ***REMOVED***
		t.Fatal("Expected to not include a status key, got true")
	***REMOVED***
	f.Add("status", "running")
	if !f.Include("status") ***REMOVED***
		t.Fatal("Expected to include a status key, got false")
	***REMOVED***
***REMOVED***

func TestValidate(t *testing.T) ***REMOVED***
	f := NewArgs()
	f.Add("status", "running")

	valid := map[string]bool***REMOVED***
		"status":   true,
		"dangling": true,
	***REMOVED***

	if err := f.Validate(valid); err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***

	f.Add("bogus", "running")
	if err := f.Validate(valid); err == nil ***REMOVED***
		t.Fatal("Expected to return an error, got nil")
	***REMOVED***
***REMOVED***

func TestWalkValues(t *testing.T) ***REMOVED***
	f := NewArgs()
	f.Add("status", "running")
	f.Add("status", "paused")

	f.WalkValues("status", func(value string) error ***REMOVED***
		if value != "running" && value != "paused" ***REMOVED***
			t.Fatalf("Unexpected value %s", value)
		***REMOVED***
		return nil
	***REMOVED***)

	err := f.WalkValues("status", func(value string) error ***REMOVED***
		return errors.New("return")
	***REMOVED***)
	if err == nil ***REMOVED***
		t.Fatal("Expected to get an error, got nil")
	***REMOVED***

	err = f.WalkValues("foo", func(value string) error ***REMOVED***
		return errors.New("return")
	***REMOVED***)
	if err != nil ***REMOVED***
		t.Fatalf("Expected to not iterate when the field doesn't exist, got %v", err)
	***REMOVED***
***REMOVED***

func TestFuzzyMatch(t *testing.T) ***REMOVED***
	f := NewArgs()
	f.Add("container", "foo")

	cases := map[string]bool***REMOVED***
		"foo":    true,
		"foobar": true,
		"barfoo": false,
		"bar":    false,
	***REMOVED***
	for source, match := range cases ***REMOVED***
		got := f.FuzzyMatch("container", source)
		if got != match ***REMOVED***
			t.Fatalf("Expected %v, got %v: %s", match, got, source)
		***REMOVED***
	***REMOVED***
***REMOVED***
