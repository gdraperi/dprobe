package parsers

import (
	"reflect"
	"testing"
)

func TestParseKeyValueOpt(t *testing.T) ***REMOVED***
	invalids := map[string]string***REMOVED***
		"":    "Unable to parse key/value option: ",
		"key": "Unable to parse key/value option: key",
	***REMOVED***
	for invalid, expectedError := range invalids ***REMOVED***
		if _, _, err := ParseKeyValueOpt(invalid); err == nil || err.Error() != expectedError ***REMOVED***
			t.Fatalf("Expected error %v for %v, got %v", expectedError, invalid, err)
		***REMOVED***
	***REMOVED***
	valids := map[string][]string***REMOVED***
		"key=value":               ***REMOVED***"key", "value"***REMOVED***,
		" key = value ":           ***REMOVED***"key", "value"***REMOVED***,
		"key=value1=value2":       ***REMOVED***"key", "value1=value2"***REMOVED***,
		" key = value1 = value2 ": ***REMOVED***"key", "value1 = value2"***REMOVED***,
	***REMOVED***
	for valid, expectedKeyValue := range valids ***REMOVED***
		key, value, err := ParseKeyValueOpt(valid)
		if err != nil ***REMOVED***
			t.Fatal(err)
		***REMOVED***
		if key != expectedKeyValue[0] || value != expectedKeyValue[1] ***REMOVED***
			t.Fatalf("Expected ***REMOVED***%v: %v***REMOVED*** got ***REMOVED***%v: %v***REMOVED***", expectedKeyValue[0], expectedKeyValue[1], key, value)
		***REMOVED***
	***REMOVED***
***REMOVED***

func TestParseUintList(t *testing.T) ***REMOVED***
	valids := map[string]map[int]bool***REMOVED***
		"":             ***REMOVED******REMOVED***,
		"7":            ***REMOVED***7: true***REMOVED***,
		"1-6":          ***REMOVED***1: true, 2: true, 3: true, 4: true, 5: true, 6: true***REMOVED***,
		"0-7":          ***REMOVED***0: true, 1: true, 2: true, 3: true, 4: true, 5: true, 6: true, 7: true***REMOVED***,
		"0,3-4,7,8-10": ***REMOVED***0: true, 3: true, 4: true, 7: true, 8: true, 9: true, 10: true***REMOVED***,
		"0-0,0,1-4":    ***REMOVED***0: true, 1: true, 2: true, 3: true, 4: true***REMOVED***,
		"03,1-3":       ***REMOVED***1: true, 2: true, 3: true***REMOVED***,
		"3,2,1":        ***REMOVED***1: true, 2: true, 3: true***REMOVED***,
		"0-2,3,1":      ***REMOVED***0: true, 1: true, 2: true, 3: true***REMOVED***,
	***REMOVED***
	for k, v := range valids ***REMOVED***
		out, err := ParseUintList(k)
		if err != nil ***REMOVED***
			t.Fatalf("Expected not to fail, got %v", err)
		***REMOVED***
		if !reflect.DeepEqual(out, v) ***REMOVED***
			t.Fatalf("Expected %v, got %v", v, out)
		***REMOVED***
	***REMOVED***

	invalids := []string***REMOVED***
		"this",
		"1--",
		"1-10,,10",
		"10-1",
		"-1",
		"-1,0",
	***REMOVED***
	for _, v := range invalids ***REMOVED***
		if out, err := ParseUintList(v); err == nil ***REMOVED***
			t.Fatalf("Expected failure with %s but got %v", v, out)
		***REMOVED***
	***REMOVED***
***REMOVED***
