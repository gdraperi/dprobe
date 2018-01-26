package aws

// JSONValue is a representation of a grab bag type that will be marshaled
// into a json string. This type can be used just like any other map.
//
//	Example:
//
//	values := aws.JSONValue***REMOVED***
//		"Foo": "Bar",
//	***REMOVED***
//	values["Baz"] = "Qux"
type JSONValue map[string]interface***REMOVED******REMOVED***
