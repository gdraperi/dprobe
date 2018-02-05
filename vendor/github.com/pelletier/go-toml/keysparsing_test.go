package toml

import (
	"fmt"
	"testing"
)

func testResult(t *testing.T, key string, expected []string) ***REMOVED***
	parsed, err := parseKey(key)
	t.Logf("key=%s expected=%s parsed=%s", key, expected, parsed)
	if err != nil ***REMOVED***
		t.Fatal("Unexpected error:", err)
	***REMOVED***
	if len(expected) != len(parsed) ***REMOVED***
		t.Fatal("Expected length", len(expected), "but", len(parsed), "parsed")
	***REMOVED***
	for index, expectedKey := range expected ***REMOVED***
		if expectedKey != parsed[index] ***REMOVED***
			t.Fatal("Expected", expectedKey, "at index", index, "but found", parsed[index])
		***REMOVED***
	***REMOVED***
***REMOVED***

func testError(t *testing.T, key string, expectedError string) ***REMOVED***
	res, err := parseKey(key)
	if err == nil ***REMOVED***
		t.Fatalf("Expected error, but succesfully parsed key %s", res)
	***REMOVED***
	if fmt.Sprintf("%s", err) != expectedError ***REMOVED***
		t.Fatalf("Expected error \"%s\", but got \"%s\".", expectedError, err)
	***REMOVED***
***REMOVED***

func TestBareKeyBasic(t *testing.T) ***REMOVED***
	testResult(t, "test", []string***REMOVED***"test"***REMOVED***)
***REMOVED***

func TestBareKeyDotted(t *testing.T) ***REMOVED***
	testResult(t, "this.is.a.key", []string***REMOVED***"this", "is", "a", "key"***REMOVED***)
***REMOVED***

func TestDottedKeyBasic(t *testing.T) ***REMOVED***
	testResult(t, "\"a.dotted.key\"", []string***REMOVED***"a.dotted.key"***REMOVED***)
***REMOVED***

func TestBaseKeyPound(t *testing.T) ***REMOVED***
	testError(t, "hello#world", "invalid bare character: #")
***REMOVED***

func TestQuotedKeys(t *testing.T) ***REMOVED***
	testResult(t, `hello."foo".bar`, []string***REMOVED***"hello", "foo", "bar"***REMOVED***)
	testResult(t, `"hello!"`, []string***REMOVED***"hello!"***REMOVED***)
	testResult(t, `foo."ba.r".baz`, []string***REMOVED***"foo", "ba.r", "baz"***REMOVED***)

	// escape sequences must not be converted
	testResult(t, `"hello\tworld"`, []string***REMOVED***`hello\tworld`***REMOVED***)
***REMOVED***

func TestEmptyKey(t *testing.T) ***REMOVED***
	testError(t, "", "empty key")
	testError(t, " ", "empty key")
	testResult(t, `""`, []string***REMOVED***""***REMOVED***)
***REMOVED***
