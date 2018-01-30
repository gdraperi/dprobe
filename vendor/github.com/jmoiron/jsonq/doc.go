/*
Package jsonq simplify your json usage with a simple hierarchical query.

Given some json data like:

	***REMOVED***
		"foo": 1,
		"bar": 2,
		"test": "Hello, world!",
		"baz": 123.1,
		"array": [
			***REMOVED***"foo": 1***REMOVED***,
			***REMOVED***"bar": 2***REMOVED***,
			***REMOVED***"baz": 3***REMOVED***
		],
		"subobj": ***REMOVED***
			"foo": 1,
			"subarray": [1,2,3],
			"subsubobj": ***REMOVED***
				"bar": 2,
				"baz": 3,
				"array": ["hello", "world"]
			***REMOVED***
		***REMOVED***,
		"bool": true
	***REMOVED***

Decode it into a map[string]interrface***REMOVED******REMOVED***:

	import (
		"strings"
		"encoding/json"
		"github.com/jmoiron/jsonq"
	)

	data := map[string]interface***REMOVED******REMOVED******REMOVED******REMOVED***
	dec := json.NewDecoder(strings.NewReader(jsonstring))
	dec.Decode(&data)
	jq := jsonq.NewQuery(data)

From here, you can query along different keys and indexes:

	// data["foo"] -> 1
	jq.Int("foo")

	// data["subobj"]["subarray"][1] -> 2
	jq.Int("subobj", "subarray", "1")

	// data["subobj"]["subarray"]["array"][0] -> "hello"
	jq.String("subobj", "subsubobj", "array", "0")

	// data["subobj"] -> map[string]interface***REMOVED******REMOVED******REMOVED***"subobj": ...***REMOVED***
	obj, err := jq.Object("subobj")

	Notes:

Missing keys, out of bounds indexes, and type failures will return errors.
For simplicity, integer keys (ie, ***REMOVED***"0": "zero"***REMOVED***) are inaccessible by `jsonq`
as integer strings are assumed to be array indexes.

*/
package jsonq
