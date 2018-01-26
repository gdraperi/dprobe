package httputils

import (
	"net/http"
	"net/url"
	"testing"
)

func TestBoolValue(t *testing.T) ***REMOVED***
	cases := map[string]bool***REMOVED***
		"":      false,
		"0":     false,
		"no":    false,
		"false": false,
		"none":  false,
		"1":     true,
		"yes":   true,
		"true":  true,
		"one":   true,
		"100":   true,
	***REMOVED***

	for c, e := range cases ***REMOVED***
		v := url.Values***REMOVED******REMOVED***
		v.Set("test", c)
		r, _ := http.NewRequest("POST", "", nil)
		r.Form = v

		a := BoolValue(r, "test")
		if a != e ***REMOVED***
			t.Fatalf("Value: %s, expected: %v, actual: %v", c, e, a)
		***REMOVED***
	***REMOVED***
***REMOVED***

func TestBoolValueOrDefault(t *testing.T) ***REMOVED***
	r, _ := http.NewRequest("GET", "", nil)
	if !BoolValueOrDefault(r, "queryparam", true) ***REMOVED***
		t.Fatal("Expected to get true default value, got false")
	***REMOVED***

	v := url.Values***REMOVED******REMOVED***
	v.Set("param", "")
	r, _ = http.NewRequest("GET", "", nil)
	r.Form = v
	if BoolValueOrDefault(r, "param", true) ***REMOVED***
		t.Fatal("Expected not to get true")
	***REMOVED***
***REMOVED***

func TestInt64ValueOrZero(t *testing.T) ***REMOVED***
	cases := map[string]int64***REMOVED***
		"":     0,
		"asdf": 0,
		"0":    0,
		"1":    1,
	***REMOVED***

	for c, e := range cases ***REMOVED***
		v := url.Values***REMOVED******REMOVED***
		v.Set("test", c)
		r, _ := http.NewRequest("POST", "", nil)
		r.Form = v

		a := Int64ValueOrZero(r, "test")
		if a != e ***REMOVED***
			t.Fatalf("Value: %s, expected: %v, actual: %v", c, e, a)
		***REMOVED***
	***REMOVED***
***REMOVED***

func TestInt64ValueOrDefault(t *testing.T) ***REMOVED***
	cases := map[string]int64***REMOVED***
		"":   -1,
		"-1": -1,
		"42": 42,
	***REMOVED***

	for c, e := range cases ***REMOVED***
		v := url.Values***REMOVED******REMOVED***
		v.Set("test", c)
		r, _ := http.NewRequest("POST", "", nil)
		r.Form = v

		a, err := Int64ValueOrDefault(r, "test", -1)
		if a != e ***REMOVED***
			t.Fatalf("Value: %s, expected: %v, actual: %v", c, e, a)
		***REMOVED***
		if err != nil ***REMOVED***
			t.Fatalf("Error should be nil, but received: %s", err)
		***REMOVED***
	***REMOVED***
***REMOVED***

func TestInt64ValueOrDefaultWithError(t *testing.T) ***REMOVED***
	v := url.Values***REMOVED******REMOVED***
	v.Set("test", "invalid")
	r, _ := http.NewRequest("POST", "", nil)
	r.Form = v

	_, err := Int64ValueOrDefault(r, "test", -1)
	if err == nil ***REMOVED***
		t.Fatal("Expected an error.")
	***REMOVED***
***REMOVED***
