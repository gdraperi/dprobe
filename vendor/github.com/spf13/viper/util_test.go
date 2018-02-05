// Copyright © 2016 Steve Francia <spf@spf13.com>.
//
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

// Viper is a application configuration system.
// It believes that applications can be configured a variety of ways
// via flags, ENVIRONMENT variables, configuration files retrieved
// from the file system, or a remote key/value store.

package viper

import (
	"reflect"
	"testing"
)

func TestCopyAndInsensitiviseMap(t *testing.T) ***REMOVED***
	var (
		given = map[string]interface***REMOVED******REMOVED******REMOVED***
			"Foo": 32,
			"Bar": map[interface***REMOVED******REMOVED***]interface ***REMOVED***
			***REMOVED******REMOVED***
				"ABc": "A",
				"cDE": "B"***REMOVED***,
		***REMOVED***
		expected = map[string]interface***REMOVED******REMOVED******REMOVED***
			"foo": 32,
			"bar": map[string]interface ***REMOVED***
			***REMOVED******REMOVED***
				"abc": "A",
				"cde": "B"***REMOVED***,
		***REMOVED***
	)

	got := copyAndInsensitiviseMap(given)

	if !reflect.DeepEqual(got, expected) ***REMOVED***
		t.Fatalf("Got %q\nexpected\n%q", got, expected)
	***REMOVED***

	if _, ok := given["foo"]; ok ***REMOVED***
		t.Fatal("Input map changed")
	***REMOVED***

	if _, ok := given["bar"]; ok ***REMOVED***
		t.Fatal("Input map changed")
	***REMOVED***

	m := given["Bar"].(map[interface***REMOVED******REMOVED***]interface***REMOVED******REMOVED***)
	if _, ok := m["ABc"]; !ok ***REMOVED***
		t.Fatal("Input map changed")
	***REMOVED***
***REMOVED***
